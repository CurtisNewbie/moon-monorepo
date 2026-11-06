package vfm

import (
	"errors"
	"reflect"
	"testing"

	fstore "github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/miso/errs"
	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/randutil"
	"gorm.io/gorm"
)

// zipTreePreviewFake replaces the mini-fstore client and the redis cache at the
// handler seam.
type zipTreePreviewFake struct {
	tree fstore.ZipTreeResp

	cacheHit      bool
	cacheReadErr  error
	cacheWriteErr error
	fetchErr      error

	fetches   []fstore.ZipTreeReq
	writes    []fstore.ZipTreeResp
	readKeys  []string
	writeKeys []string
}

func (f *zipTreePreviewFake) deps() zipTreeDeps {
	return zipTreeDeps{
		fetchTree: func(_ miso.Rail, req fstore.ZipTreeReq) (fstore.ZipTreeResp, error) {
			f.fetches = append(f.fetches, req)
			if f.fetchErr != nil {
				return fstore.ZipTreeResp{}, f.fetchErr
			}
			return f.tree, nil
		},
		cacheRead: func(_ miso.Rail, key string) (fstore.ZipTreeResp, bool, error) {
			f.readKeys = append(f.readKeys, key)
			if f.cacheReadErr != nil {
				return fstore.ZipTreeResp{}, false, f.cacheReadErr
			}
			return f.tree, f.cacheHit, nil
		},
		cacheWrite: func(_ miso.Rail, key string, tree fstore.ZipTreeResp) error {
			f.writeKeys = append(f.writeKeys, key)
			f.writes = append(f.writes, tree)
			return f.cacheWriteErr
		},
		maxDepth:   func() int { return miso.GetPropInt(PropZipTreeMaxDepth) },
		maxEntries: func() int { return miso.GetPropInt(PropZipTreeMaxEntries) },
	}
}

// callZipTreeHandler invokes the real handler with the mini-fstore client and
// cache replaced by fake.
func callZipTreeHandler(t *testing.T, db *gorm.DB, fileKey string, user flow.User, fake *zipTreePreviewFake) (fstore.ZipTreeResp, error) {
	t.Helper()
	old := newZipTreeDeps
	newZipTreeDeps = fake.deps
	t.Cleanup(func() { newZipTreeDeps = old })
	return ApiZipTreePreview(miso.EmptyRail(), db, ZipTreePreviewReq{FileKey: fileKey}, user)
}

func zipTreeSample() fstore.ZipTreeResp {
	return fstore.ZipTreeResp{
		Tree: fstore.ZipTreeNode{
			Name: "",
			Dir:  true,
			Children: []*fstore.ZipTreeNode{
				{Name: "docs", Dir: true, Children: []*fstore.ZipTreeNode{
					{Name: "readme.txt", Size: 42},
				}},
			},
		},
		Truncated:        true,
		TruncatedReasons: []string{fstore.ZipTreeTruncatedByDepth},
		IncludedEntries:  2,
	}
}

func insertZipPreviewRow(t *testing.T, db *gorm.DB, uploaderNo, fileType, fstoreFileId string, logicDeleted int) string {
	t.Helper()
	uuid := "zt-" + randutil.ERand(10)
	if err := db.Exec(`INSERT INTO file_info (name, uuid, size_in_bytes, file_type, fstore_file_id, uploader_no, is_logic_deleted, is_del)
		VALUES (?,?,?,?,?,?,?,?)`,
		"archive.zip", uuid, 0, fileType, fstoreFileId, uploaderNo, logicDeleted, 0).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Exec("DELETE FROM file_info WHERE uuid = ?", uuid)
	})
	return uuid
}

func TestZipTreePreviewServesOwner(t *testing.T) {
	db := testMySQLConn(t)
	user := testUser()
	fileKey := insertZipPreviewRow(t, db, user.UserNo, FileTypeFile, "fs-zip-owner", 0)
	want := zipTreeSample()
	fake := &zipTreePreviewFake{tree: want}

	got, err := callZipTreeHandler(t, db, fileKey, user, fake)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("tree: got %+v, want %+v", got, want)
	}
	if len(fake.fetches) != 1 {
		t.Fatalf("expected 1 mini-fstore call, got %d", len(fake.fetches))
	}
	if fake.fetches[0].FileId != "fs-zip-owner" {
		t.Fatalf("mini-fstore fileId: got %q, want %q", fake.fetches[0].FileId, "fs-zip-owner")
	}
	if len(fake.writeKeys) != 1 || fake.writeKeys[0] != "fs-zip-owner" {
		t.Fatalf("cache write key: got %v, want [fs-zip-owner]", fake.writeKeys)
	}
}

func TestZipTreePreviewRefusesNonOwnerBeforeCacheOrFstore(t *testing.T) {
	db := testMySQLConn(t)
	fileKey := insertZipPreviewRow(t, db, "UE_OTHER_USER", FileTypeFile, "fs-zip-other", 0)
	fake := &zipTreePreviewFake{tree: zipTreeSample()}

	_, err := callZipTreeHandler(t, db, fileKey, testUser(), fake)
	if !errors.Is(err, errs.ErrNotPermitted) {
		t.Fatalf("error = %v, want NOT_PERMITTED", err)
	}
	if len(fake.readKeys) != 0 {
		t.Fatalf("cache read must not run for a non-owner, got %v", fake.readKeys)
	}
	if len(fake.fetches) != 0 {
		t.Fatalf("mini-fstore must not run for a non-owner, got %+v", fake.fetches)
	}
}

func TestZipTreePreviewMissingRow(t *testing.T) {
	db := testMySQLConn(t)
	fake := &zipTreePreviewFake{tree: zipTreeSample()}

	_, err := callZipTreeHandler(t, db, "zt-missing-"+randutil.ERand(8), testUser(), fake)
	if !errors.Is(err, ErrFileNotFound) {
		t.Fatalf("error = %v, want FILE_NOT_FOUND", err)
	}
	if len(fake.readKeys) != 0 || len(fake.fetches) != 0 {
		t.Fatalf("a missing row must not read the cache or call mini-fstore")
	}
}

func TestZipTreePreviewDeletedRow(t *testing.T) {
	db := testMySQLConn(t)
	user := testUser()
	fileKey := insertZipPreviewRow(t, db, user.UserNo, FileTypeFile, "fs-zip-deleted", 1)
	fake := &zipTreePreviewFake{tree: zipTreeSample()}

	_, err := callZipTreeHandler(t, db, fileKey, user, fake)
	if !errors.Is(err, ErrFileDeleted) {
		t.Fatalf("error = %v, want FILE_DELETED", err)
	}
	if len(fake.readKeys) != 0 || len(fake.fetches) != 0 {
		t.Fatalf("a deleted row must not read the cache or call mini-fstore")
	}
}

func TestZipTreePreviewFolderKey(t *testing.T) {
	db := testMySQLConn(t)
	user := testUser()
	fileKey := insertZipPreviewRow(t, db, user.UserNo, FileTypeDir, "", 0)
	fake := &zipTreePreviewFake{tree: zipTreeSample()}

	_, err := callZipTreeHandler(t, db, fileKey, user, fake)
	if !errors.Is(err, errs.ErrNotPermitted) {
		t.Fatalf("error = %v, want NOT_PERMITTED", err)
	}
	if len(fake.readKeys) != 0 || len(fake.fetches) != 0 {
		t.Fatalf("a folder key must not read the cache or call mini-fstore")
	}
}

func TestZipTreePreviewCacheHitSkipsFstore(t *testing.T) {
	db := testMySQLConn(t)
	user := testUser()
	fileKey := insertZipPreviewRow(t, db, user.UserNo, FileTypeFile, "fs-zip-hit", 0)
	want := zipTreeSample()
	fake := &zipTreePreviewFake{tree: want, cacheHit: true}

	got, err := callZipTreeHandler(t, db, fileKey, user, fake)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("tree: got %+v, want %+v", got, want)
	}
	if len(fake.readKeys) != 1 || fake.readKeys[0] != "fs-zip-hit" {
		t.Fatalf("cache read key: got %v, want [fs-zip-hit]", fake.readKeys)
	}
	if len(fake.fetches) != 0 {
		t.Fatalf("a cache hit must not call mini-fstore, got %+v", fake.fetches)
	}
	if len(fake.writes) != 0 {
		t.Fatalf("a cache hit must not write the cache again, got %+v", fake.writes)
	}
}

func TestZipTreePreviewCacheMissSendsConfiguredCaps(t *testing.T) {
	db := testMySQLConn(t)
	user := testUser()
	fileKey := insertZipPreviewRow(t, db, user.UserNo, FileTypeFile, "fs-zip-caps", 0)
	setZipTreeProp(t, PropZipTreeMaxDepth, 7)
	setZipTreeProp(t, PropZipTreeMaxEntries, 123)
	fake := &zipTreePreviewFake{tree: zipTreeSample()}

	if _, err := callZipTreeHandler(t, db, fileKey, user, fake); err != nil {
		t.Fatal(err)
	}
	if len(fake.fetches) != 1 {
		t.Fatalf("expected 1 mini-fstore call, got %d", len(fake.fetches))
	}
	req := fake.fetches[0]
	if req.MaxDepth != 7 || req.MaxEntries != 123 {
		t.Fatalf("caps: got depth=%d entries=%d, want 7/123", req.MaxDepth, req.MaxEntries)
	}
}

func TestZipTreePreviewCapsDefaultTo32And10000(t *testing.T) {
	db := testMySQLConn(t)
	user := testUser()
	fileKey := insertZipPreviewRow(t, db, user.UserNo, FileTypeFile, "fs-zip-defaults", 0)
	fake := &zipTreePreviewFake{tree: zipTreeSample()}

	if _, err := callZipTreeHandler(t, db, fileKey, user, fake); err != nil {
		t.Fatal(err)
	}
	req := fake.fetches[0]
	if req.MaxDepth != 32 || req.MaxEntries != 10000 {
		t.Fatalf("default caps: got depth=%d entries=%d, want 32/10000", req.MaxDepth, req.MaxEntries)
	}
}

func TestZipTreePreviewCacheReadFailureFallsBackToFstore(t *testing.T) {
	db := testMySQLConn(t)
	user := testUser()
	fileKey := insertZipPreviewRow(t, db, user.UserNo, FileTypeFile, "fs-zip-readfail", 0)
	want := zipTreeSample()
	fake := &zipTreePreviewFake{tree: want, cacheReadErr: errs.NewErrf("redis unavailable")}

	got, err := callZipTreeHandler(t, db, fileKey, user, fake)
	if err != nil {
		t.Fatalf("cache read failure must not fail the preview, got %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("tree: got %+v, want %+v", got, want)
	}
	if len(fake.fetches) != 1 {
		t.Fatalf("expected 1 mini-fstore call, got %d", len(fake.fetches))
	}
	if len(fake.writes) != 1 {
		t.Fatalf("expected the freshly built tree to be cached, got %d writes", len(fake.writes))
	}
}

func TestZipTreePreviewCacheWriteFailureStillReturnsTree(t *testing.T) {
	db := testMySQLConn(t)
	user := testUser()
	fileKey := insertZipPreviewRow(t, db, user.UserNo, FileTypeFile, "fs-zip-writefail", 0)
	want := zipTreeSample()
	fake := &zipTreePreviewFake{tree: want, cacheWriteErr: errs.NewErrf("redis unavailable")}

	got, err := callZipTreeHandler(t, db, fileKey, user, fake)
	if err != nil {
		t.Fatalf("cache write failure must not fail the preview, got %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("tree: got %+v, want %+v", got, want)
	}
	if len(fake.fetches) != 1 || len(fake.writes) != 1 {
		t.Fatalf("expected one fstore read and one cache write, got %d/%d", len(fake.fetches), len(fake.writes))
	}
}

func TestZipTreePreviewIllegalFormatIsNotCached(t *testing.T) {
	db := testMySQLConn(t)
	user := testUser()
	fileKey := insertZipPreviewRow(t, db, user.UserNo, FileTypeFile, "fs-zip-broken", 0)
	fake := &zipTreePreviewFake{fetchErr: fstore.ErrIllegalFormat.New()}

	_, err := callZipTreeHandler(t, db, fileKey, user, fake)
	if !errors.Is(err, fstore.ErrIllegalFormat) {
		t.Fatalf("error = %v, want ILLEGAL_FORMAT", err)
	}
	if len(fake.writes) != 0 {
		t.Fatalf("a failed preview must never be cached, got %+v", fake.writes)
	}
}

func TestZipTreePreviewFstoreErrorsKeepTheirCodes(t *testing.T) {
	db := testMySQLConn(t)
	user := testUser()
	fileKey := insertZipPreviewRow(t, db, user.UserNo, FileTypeFile, "fs-zip-gone", 0)

	cases := []struct {
		name     string
		wantErr  error
		wantCode string
	}{
		{"missing file", fstore.ErrFileNotFound.New(), fstore.FileNotFound},
		{"illegal format", fstore.ErrIllegalFormat.New(), fstore.IllegalFormat},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fake := &zipTreePreviewFake{fetchErr: c.wantErr}
			_, err := callZipTreeHandler(t, db, fileKey, user, fake)
			if !errors.Is(err, c.wantErr) {
				t.Fatalf("error = %v, want %v", err, c.wantErr)
			}
			var me *errs.MisoErr
			if !errors.As(err, &me) {
				t.Fatalf("error = %v, want a MisoErr", err)
			}
			if me.Code() != c.wantCode {
				t.Fatalf("code = %q, want %q", me.Code(), c.wantCode)
			}
		})
	}
}

func setZipTreeProp(t *testing.T, key string, value int) {
	t.Helper()
	old := miso.GetPropInt(key)
	miso.SetProp(key, value)
	t.Cleanup(func() {
		miso.SetProp(key, old)
	})
}
