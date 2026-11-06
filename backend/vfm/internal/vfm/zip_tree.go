package vfm

import (
	"time"

	fstore "github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/miso/errs"
	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/miso"
	"gorm.io/gorm"
)

// zipTreeCache holds successful previews for five minutes, keyed by the stored
// file id. A cache failure only degrades the preview to a direct read.
var zipTreeCache = redis.NewRCache[fstore.ZipTreeResp]("vfm:zip:tree", redis.RCacheConfig{Exp: 5 * time.Minute})

// zipTreeDeps is the handler seam: the mini-fstore call and the cache are
// replaced with fakes in tests.
type zipTreeDeps struct {
	fetchTree  func(rail miso.Rail, req fstore.ZipTreeReq) (fstore.ZipTreeResp, error)
	cacheRead  func(rail miso.Rail, key string) (fstore.ZipTreeResp, bool, error)
	cacheWrite func(rail miso.Rail, key string, tree fstore.ZipTreeResp) error
	maxDepth   func() int
	maxEntries func() int
}

func defaultZipTreeDeps() zipTreeDeps {
	return zipTreeDeps{
		fetchTree:  fstore.ZipTree,
		cacheRead:  zipTreeCache.Get,
		cacheWrite: zipTreeCache.Put,
		maxDepth:   func() int { return miso.GetPropInt(PropZipTreeMaxDepth) },
		maxEntries: func() int { return miso.GetPropInt(PropZipTreeMaxEntries) },
	}
}

// newZipTreeDeps is replaced by tests to fake the mini-fstore client and cache
// at the handler seam.
var newZipTreeDeps = defaultZipTreeDeps

// PreviewZipTree reads a zip file's tree for its uploader, caching a successful
// reply for five minutes under the stored file id. Every access check runs
// before the cache is consulted.
func PreviewZipTree(rail miso.Rail, db *gorm.DB, fileKey string, user flow.User) (fstore.ZipTreeResp, error) {
	return zipTreePreview(rail, db, fileKey, user, newZipTreeDeps())
}

func zipTreePreview(rail miso.Rail, db *gorm.DB, fileKey string, user flow.User, deps zipTreeDeps) (fstore.ZipTreeResp, error) {
	var zero fstore.ZipTreeResp

	f, ok, err := findFile(rail, db, fileKey)
	if err != nil {
		return zero, errs.Wrapf(err, "failed to find file for zip tree preview, fileKey: %v", fileKey)
	}
	if !ok {
		return zero, ErrFileNotFound.New()
	}
	// ownership comes before the deleted/type checks, like the other
	// uploader-only operations, so a stranger learns nothing about the row
	if f.UploaderNo != user.UserNo {
		return zero, errs.ErrNotPermitted.New()
	}
	if f.IsLogicDeleted == DelY {
		return zero, ErrFileDeleted.New()
	}
	if f.FileType != FileTypeFile {
		return zero, errs.ErrNotPermitted.New()
	}

	cacheKey := f.FstoreFileId
	if cacheKey != "" {
		tree, hit, err := deps.cacheRead(rail, cacheKey)
		if err != nil {
			rail.Warnf("Failed to read zip tree cache, fileKey: %v, %v", fileKey, err)
		} else if hit {
			return tree, nil
		}
	}

	tree, err := deps.fetchTree(rail, fstore.ZipTreeReq{
		FileId:     cacheKey,
		MaxDepth:   deps.maxDepth(),
		MaxEntries: deps.maxEntries(),
	})
	if err != nil {
		// keep the mini-fstore error as is, so FILE_NOT_FOUND and ILLEGAL_FORMAT
		// keep their codes on the way to the caller
		rail.Warnf("Failed to read zip tree from mini-fstore, fileKey: %v, %v", fileKey, err)
		return zero, err
	}

	// a failed preview is never cached; a cache write failure still returns the
	// freshly built tree
	if cacheKey != "" {
		if err := deps.cacheWrite(rail, cacheKey, tree); err != nil {
			rail.Warnf("Failed to cache zip tree, fileKey: %v, %v", fileKey, err)
		}
	}
	return tree, nil
}
