package vfm

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"

	fstore "github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/miso/errs"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/randutil"
	"gorm.io/gorm"
)

// --- fstore fake ------------------------------------------------------------

// dirZipMint captures the mint requests vfm sends to mini-fstore. The archive
// itself is built by mini-fstore, so the vfm seam ends at this HTTP call.
type dirZipMint struct {
	mu   sync.Mutex
	reqs []fstore.GenBundleKeyReq
}

func (m *dirZipMint) add(req fstore.GenBundleKeyReq) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.reqs = append(m.reqs, req)
}

func (m *dirZipMint) all() []fstore.GenBundleKeyReq {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]fstore.GenBundleKeyReq(nil), m.reqs...)
}

func (m *dirZipMint) count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.reqs)
}

// startDirZipFakeFstore points the fstore client at a fake that answers the
// bundle mint endpoint with a fixed key.
func startDirZipFakeFstore(t *testing.T) *dirZipMint {
	t.Helper()

	oldHost := miso.GetPropStr("client.addr.fstore.host")
	oldPort := miso.GetPropInt("client.addr.fstore.port")
	t.Cleanup(func() {
		miso.SetProp("client.addr.fstore.host", oldHost)
		miso.SetProp("client.addr.fstore.port", oldPort)
	})

	m := &dirZipMint{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/file/bundle/key" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var req fstore.GenBundleKeyReq
		if err := json.Unmarshal(body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		m.add(req)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(miso.GnResp[string]{
			Msg:  "ok",
			Data: "bkey-test",
		})
	}))
	t.Cleanup(srv.Close)

	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		t.Fatal(err)
	}
	miso.SetProp("client.addr.fstore.host", u.Hostname())
	miso.SetProp("client.addr.fstore.port", port)
	return m
}

// --- fixtures ---------------------------------------------------------------

func insertDirZipRow(t *testing.T, db *gorm.DB, parent, name, userNo, fileType, fstoreFileId string, logicDeleted, isDel int) string {
	t.Helper()
	uuid := "dz-" + randutil.ERand(10)
	if err := db.Exec(`INSERT INTO file_info (name, uuid, size_in_bytes, file_type, parent_file, fstore_file_id, uploader_no, is_logic_deleted, is_del)
		VALUES (?,?,?,?,?,?,?,?,?)`,
		name, uuid, 0, fileType, parent, fstoreFileId, userNo, logicDeleted, isDel).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Exec("DELETE FROM file_info WHERE uuid = ?", uuid)
	})
	return uuid
}

func insertDirZipDir(t *testing.T, db *gorm.DB, parent, name, userNo string) string {
	t.Helper()
	return insertDirZipRow(t, db, parent, name, userNo, FileTypeDir, "", 0, 0)
}

func insertDirZipChild(t *testing.T, db *gorm.DB, parent, name, userNo, fstoreFileId string) string {
	t.Helper()
	return insertDirZipRow(t, db, parent, name, userNo, FileTypeFile, fstoreFileId, 0, 0)
}

func insertDirZipHiddenChild(t *testing.T, db *gorm.DB, parent, name, userNo, fstoreFileId string) string {
	t.Helper()
	uuid := insertDirZipChild(t, db, parent, name, userNo, fstoreFileId)
	if err := db.Exec("UPDATE file_info SET hidden = 1 WHERE uuid = ?", uuid).Error; err != nil {
		t.Fatal(err)
	}
	return uuid
}

// --- happy path -------------------------------------------------------------

func TestGenDirBundleKeyMintsDirectFileChildrenInNameOrder(t *testing.T) {
	db := testMySQLConn(t)
	mint := startDirZipFakeFstore(t)
	user := testUser()

	dir := insertDirZipDir(t, db, "", "photos", user.UserNo)
	insertDirZipChild(t, db, dir, "b.jpg", user.UserNo, "FS_B")
	insertDirZipChild(t, db, dir, "a.jpg", user.UserNo, "FS_A")
	insertDirZipChild(t, db, dir, "c.jpg", user.UserNo, "FS_C")

	// a sub-directory and its contents never enter the archive
	sub := insertDirZipDir(t, db, dir, "nested", user.UserNo)
	insertDirZipChild(t, db, sub, "inside.jpg", user.UserNo, "FS_NESTED")

	// deleted children are skipped
	insertDirZipRow(t, db, dir, "deleted.jpg", user.UserNo, FileTypeFile, "FS_DELETED", 1, 0)

	key, err := GenDirBundleKey(miso.EmptyRail(), db, GenerateDirZipReq{FileKey: dir}, user)
	if err != nil {
		t.Fatal(err)
	}
	if key != "bkey-test" {
		t.Fatalf("bundle key: got %q, want %q", key, "bkey-test")
	}

	reqs := mint.all()
	if len(reqs) != 1 {
		t.Fatalf("expected 1 mint request, got %d", len(reqs))
	}
	if reqs[0].Name != "photos" {
		t.Fatalf("archive name: got %q, want %q", reqs[0].Name, "photos")
	}

	want := []fstore.BundleItem{
		{FileId: "FS_A", Name: "a.jpg"},
		{FileId: "FS_B", Name: "b.jpg"},
		{FileId: "FS_C", Name: "c.jpg"},
	}
	if !reflect.DeepEqual(reqs[0].Items, want) {
		t.Fatalf("items: got %+v, want %+v", reqs[0].Items, want)
	}
}

func TestGenDirBundleKeyExcludesForeignAndHiddenChildren(t *testing.T) {
	db := testMySQLConn(t)
	mint := startDirZipFakeFstore(t)
	user := testUser()

	dir := insertDirZipDir(t, db, "", "mixed", user.UserNo)
	insertDirZipChild(t, db, dir, "mine.jpg", user.UserNo, "FS_MINE")
	// another user's file that was moved into this directory stays invisible in
	// the listing, so it must not enter the archive either
	insertDirZipChild(t, db, dir, "foreign.jpg", "UE_OTHER_USER", "FS_FOREIGN")
	// a hidden row (e.g. a versioned file) is likewise absent from the listing
	insertDirZipHiddenChild(t, db, dir, "hidden.jpg", user.UserNo, "FS_HIDDEN")

	if _, err := GenDirBundleKey(miso.EmptyRail(), db, GenerateDirZipReq{FileKey: dir}, user); err != nil {
		t.Fatal(err)
	}

	reqs := mint.all()
	if len(reqs) != 1 {
		t.Fatalf("expected 1 mint request, got %d", len(reqs))
	}
	want := []fstore.BundleItem{{FileId: "FS_MINE", Name: "mine.jpg"}}
	if !reflect.DeepEqual(reqs[0].Items, want) {
		t.Fatalf("items: got %+v, want %+v", reqs[0].Items, want)
	}
}

func TestGenDirBundleKeyRejectsUnsafeEntryNames(t *testing.T) {
	db := testMySQLConn(t)
	mint := startDirZipFakeFstore(t)
	user := testUser()

	cases := []struct {
		name     string
		fileName string
	}{
		{"empty", ""},
		{"slash", "a/b.txt"},
		{"backslash", `a\b.txt`},
		{"dot dot", ".."},
		{"nested dot dot", "a/../b.txt"},
		{"control character", "a\x01b.txt"},
		{"trailing dot", "report."},
		{"trailing space", "report "},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := insertDirZipDir(t, db, "", "unsafe", user.UserNo)
			insertDirZipChild(t, db, dir, c.fileName, user.UserNo, "FS_UNSAFE")

			key, err := GenDirBundleKey(miso.EmptyRail(), db, GenerateDirZipReq{FileKey: dir}, user)
			if err == nil {
				t.Fatal("expected an error")
			}
			if key != "" {
				t.Fatalf("expected no bundle key, got %q", key)
			}
			if mint.count() != 0 {
				t.Fatalf("expected no mint request, got %d", mint.count())
			}
			if !strings.Contains(err.Error(), strconv.Quote(c.fileName)) {
				t.Errorf("error %q should name the offending file %q", err, c.fileName)
			}
		})
	}
}

func TestGenDirBundleKeyIsDeterministic(t *testing.T) {
	db := testMySQLConn(t)
	mint := startDirZipFakeFstore(t)
	user := testUser()

	dir := insertDirZipDir(t, db, "", "stable", user.UserNo)
	insertDirZipChild(t, db, dir, "b.txt", user.UserNo, "FS_B")
	insertDirZipChild(t, db, dir, "a.txt", user.UserNo, "FS_A")
	insertDirZipChild(t, db, dir, "a.txt", user.UserNo, "FS_A2")

	for i := 0; i < 2; i++ {
		if _, err := GenDirBundleKey(miso.EmptyRail(), db, GenerateDirZipReq{FileKey: dir}, user); err != nil {
			t.Fatal(err)
		}
	}

	reqs := mint.all()
	if len(reqs) != 2 {
		t.Fatalf("expected 2 mint requests, got %d", len(reqs))
	}
	if reqs[0].Name != reqs[1].Name {
		t.Errorf("archive names differ: %q vs %q", reqs[0].Name, reqs[1].Name)
	}
	if !reflect.DeepEqual(reqs[0].Items, reqs[1].Items) {
		t.Errorf("items differ across mints:\n%+v\n%+v", reqs[0].Items, reqs[1].Items)
	}
}

func TestGenDirBundleKeySuffixesDuplicateNames(t *testing.T) {
	db := testMySQLConn(t)
	mint := startDirZipFakeFstore(t)
	user := testUser()

	dir := insertDirZipDir(t, db, "", "dups", user.UserNo)
	insertDirZipChild(t, db, dir, "dup.txt", user.UserNo, "FS_DUP_1")
	insertDirZipChild(t, db, dir, "dup.txt", user.UserNo, "FS_DUP_2")
	insertDirZipChild(t, db, dir, "dup.txt", user.UserNo, "FS_DUP_3")
	insertDirZipChild(t, db, dir, "license", user.UserNo, "FS_LICENSE_1")
	insertDirZipChild(t, db, dir, "license", user.UserNo, "FS_LICENSE_2")
	insertDirZipChild(t, db, dir, "notes", user.UserNo, "FS_NOTES_1")
	insertDirZipChild(t, db, dir, "notes", user.UserNo, "FS_NOTES_2")

	if _, err := GenDirBundleKey(miso.EmptyRail(), db, GenerateDirZipReq{FileKey: dir}, user); err != nil {
		t.Fatal(err)
	}

	reqs := mint.all()
	if len(reqs) != 1 {
		t.Fatalf("expected 1 mint request, got %d", len(reqs))
	}
	want := []fstore.BundleItem{
		{FileId: "FS_DUP_1", Name: "dup.txt"},
		{FileId: "FS_DUP_2", Name: "dup (1).txt"},
		{FileId: "FS_DUP_3", Name: "dup (2).txt"},
		{FileId: "FS_LICENSE_1", Name: "license"},
		{FileId: "FS_LICENSE_2", Name: "license (1)"},
		{FileId: "FS_NOTES_1", Name: "notes"},
		{FileId: "FS_NOTES_2", Name: "notes (1)"},
	}
	if !reflect.DeepEqual(reqs[0].Items, want) {
		t.Fatalf("items: got %+v, want %+v", reqs[0].Items, want)
	}
}

func TestGenDirBundleKeyEmptyDirectoryMintsEmptyBundle(t *testing.T) {
	db := testMySQLConn(t)
	mint := startDirZipFakeFstore(t)
	user := testUser()

	dir := insertDirZipDir(t, db, "", "empty", user.UserNo)

	key, err := GenDirBundleKey(miso.EmptyRail(), db, GenerateDirZipReq{FileKey: dir}, user)
	if err != nil {
		t.Fatal(err)
	}
	if key == "" {
		t.Fatal("expected a bundle key")
	}

	reqs := mint.all()
	if len(reqs) != 1 {
		t.Fatalf("expected 1 mint request, got %d", len(reqs))
	}
	if reqs[0].Name != "empty" {
		t.Fatalf("archive name: got %q, want %q", reqs[0].Name, "empty")
	}
	if len(reqs[0].Items) != 0 {
		t.Fatalf("expected no items, got %+v", reqs[0].Items)
	}
}

func TestGenDirBundleKeyKeepsNonAsciiArchiveName(t *testing.T) {
	db := testMySQLConn(t)
	mint := startDirZipFakeFstore(t)
	user := testUser()

	dir := insertDirZipDir(t, db, "", "我的照片", user.UserNo)
	insertDirZipChild(t, db, dir, "猫.jpg", user.UserNo, "FS_CAT")

	if _, err := GenDirBundleKey(miso.EmptyRail(), db, GenerateDirZipReq{FileKey: dir}, user); err != nil {
		t.Fatal(err)
	}

	reqs := mint.all()
	if len(reqs) != 1 {
		t.Fatalf("expected 1 mint request, got %d", len(reqs))
	}
	if reqs[0].Name != "我的照片" {
		t.Fatalf("archive name: got %q, want %q", reqs[0].Name, "我的照片")
	}
	if len(reqs[0].Items) != 1 || reqs[0].Items[0].Name != "猫.jpg" {
		t.Fatalf("items: got %+v", reqs[0].Items)
	}
}

func TestGenDirBundleKeyRefusesNonOwner(t *testing.T) {
	db := testMySQLConn(t)
	mint := startDirZipFakeFstore(t)

	dir := insertDirZipDir(t, db, "", "someone-elses", "UE_OTHER_USER")
	insertDirZipChild(t, db, dir, "secret.jpg", "UE_OTHER_USER", "FS_SECRET")

	key, err := GenDirBundleKey(miso.EmptyRail(), db, GenerateDirZipReq{FileKey: dir}, testUser())
	if err == nil {
		t.Fatal("expected an error")
	}
	if !errors.Is(err, errs.ErrNotPermitted) {
		t.Fatalf("error = %v, want ErrNotPermitted", err)
	}
	if key != "" {
		t.Fatalf("expected no bundle key, got %q", key)
	}
	if mint.count() != 0 {
		t.Fatalf("expected no mint request, got %d", mint.count())
	}
}

func TestGenDirBundleKeyChecksOwnershipBeforeDeletedAndType(t *testing.T) {
	db := testMySQLConn(t)
	mint := startDirZipFakeFstore(t)

	deletedDir := insertDirZipRow(t, db, "", "someone-elses-deleted", "UE_OTHER_USER", FileTypeDir, "", 1, 0)
	fileRow := insertDirZipChild(t, db, "", "someone-elses-file", "UE_OTHER_USER", "FS_PLAIN")

	for _, fileKey := range []string{deletedDir, fileRow} {
		key, err := GenDirBundleKey(miso.EmptyRail(), db, GenerateDirZipReq{FileKey: fileKey}, testUser())
		if err == nil {
			t.Fatal("expected an error")
		}
		if !errors.Is(err, errs.ErrNotPermitted) {
			t.Fatalf("error = %v, want ErrNotPermitted, so a stranger learns nothing about the uuid", err)
		}
		if key != "" {
			t.Fatalf("expected no bundle key, got %q", key)
		}
	}

	if mint.count() != 0 {
		t.Fatalf("expected no mint request, got %d", mint.count())
	}
}

func TestGenDirBundleKeyRefusesMissingFileRowAndDeletedDir(t *testing.T) {
	db := testMySQLConn(t)
	mint := startDirZipFakeFstore(t)
	user := testUser()

	fileKey := insertDirZipChild(t, db, "", "plain.txt", user.UserNo, "FS_PLAIN")
	deletedDir := insertDirZipRow(t, db, "", "deleted-dir", user.UserNo, FileTypeDir, "", 1, 0)

	cases := []struct {
		name    string
		fileKey string
		wantErr error
	}{
		{"missing", "dz-does-not-exist-" + randutil.ERand(8), ErrFileNotFound},
		{"logically deleted dir", deletedDir, ErrFileDeleted},
		{"file row", fileKey, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			key, err := GenDirBundleKey(miso.EmptyRail(), db, GenerateDirZipReq{FileKey: c.fileKey}, user)
			if err == nil {
				t.Fatal("expected an error")
			}
			if c.wantErr != nil && !errors.Is(err, c.wantErr) {
				t.Fatalf("error = %v, want %v", err, c.wantErr)
			}
			if c.wantErr == nil && !strings.Contains(err.Error(), "Only a directory") {
				t.Fatalf("error = %v, want the not-a-directory error", err)
			}
			if key != "" {
				t.Fatalf("expected no bundle key, got %q", key)
			}
		})
	}

	if mint.count() != 0 {
		t.Fatalf("expected no mint request, got %d", mint.count())
	}
}

func TestGenDirBundleKeySkipsChildWithoutFstoreFileId(t *testing.T) {
	db := testMySQLConn(t)
	mint := startDirZipFakeFstore(t)
	user := testUser()

	dir := insertDirZipDir(t, db, "", "with-ghost", user.UserNo)
	insertDirZipChild(t, db, dir, "ghost.jpg", user.UserNo, "")
	insertDirZipChild(t, db, dir, "real.jpg", user.UserNo, "fs-real")

	key, err := GenDirBundleKey(miso.EmptyRail(), db, GenerateDirZipReq{FileKey: dir}, user)
	if err != nil {
		t.Fatalf("a child without an fstore file id should be skipped, got %v", err)
	}
	if key == "" {
		t.Fatal("expected a bundle key")
	}
	m := mint.all()
	if len(m) != 1 || len(m[0].Items) != 1 || m[0].Items[0].Name != "real.jpg" {
		t.Fatalf("expected only the servable child, got %+v", m)
	}
}

func TestGenTempTokenStillRefusesDirectory(t *testing.T) {
	db := testMySQLConn(t)
	user := testUser()

	dir := insertDirZipDir(t, db, "", "not-a-file", user.UserNo)

	key, err := GenTempToken(miso.EmptyRail(), db, GenerateTempTokenReq{FileKey: dir}, user)
	if err == nil {
		t.Fatal("expected an error")
	}
	if key != "" {
		t.Fatalf("expected no token, got %q", key)
	}
}
