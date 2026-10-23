package fstore

import (
	"archive/zip"
	"mime"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/miso/errs"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/json"
	"github.com/curtisnewbie/miso/util/randutil"
)

const bundleKeyPrefix = "fstore:bundle:key:"

// maxBundleItems bounds one bundle. The whole bundle is a single Redis value and
// the whole entry list is validated before the response body starts, so an
// unbounded list would let one caller grow both without limit.
const maxBundleItems = 3000

// Lifetime of a bundle key. It is refreshed on each successful read, matching
// the single-file key.
var bundleKeyTTL = 30 * time.Minute

// RandBundleKey generates a temporary key for an ordered list of files.
//
// The file ids are deliberately not validated here: they are checked at
// download time, unlike the single-file key which is resolved eagerly.
func RandBundleKey(rail miso.Rail, name string, items []api.BundleItem) (string, error) {
	if len(items) > maxBundleItems {
		return "", errs.NewErrf("too many entries in bundle: %d, max: %d", len(items), maxBundleItems)
	}
	if items == nil {
		items = []api.BundleItem{}
	}
	bk := randutil.ERand(30)
	sby, em := json.WriteJson(api.GenBundleKeyReq{Name: name, Items: items})
	if em != nil {
		return "", errs.NewErrf("failed to marshal to bundle, %v", em)
	}
	if e := redis.Set(rail, bundleKeyPrefix+bk, string(sby), bundleKeyTTL); e != nil {
		return "", errs.Wrapf(e, "failed to cache bundle key")
	}
	return bk, nil
}

// RefreshBundleKeyExp refreshes the bundle key's expiration.
func RefreshBundleKeyExp(rail miso.Rail, bundleKey string) error {
	if _, err := redis.Expire(rail, bundleKeyPrefix+bundleKey, bundleKeyTTL); err != nil {
		rail.Warnf("Failed to refresh bundle key expiration, bundleKey: %v, %v", bundleKey, err)
		return errs.Wrapf(err, "failed to refresh bundle key expiration")
	}
	return nil
}

// ResolveBundleKey resolves the bundle for the given bundle key.
func ResolveBundleKey(rail miso.Rail, bundleKey string) (bool, api.GenBundleKeyReq) {
	var cb api.GenBundleKeyReq
	val, ok, err := redis.Get(rail, bundleKeyPrefix+bundleKey)
	if err != nil {
		rail.Errorf("Failed to find bundleKey, %v", err)
		return false, cb
	}
	if !ok {
		rail.Infof("BundleKey not found, %v", bundleKey)
		return false, cb
	}

	eu := json.ParseJson([]byte(val), &cb)
	if eu != nil {
		rail.Errorf("Failed to unmarshal bundleKey, %s, %v", bundleKey, eu)
		return false, cb
	}
	return true, cb
}

// bundleArchiveName returns the download filename for the archive, appending
// ".zip" only when the archive name does not already carry it.
func bundleArchiveName(name string) string {
	if strings.HasSuffix(strings.ToLower(name), ".zip") {
		return name
	}
	return name + ".zip"
}

// DownloadBundleKey streams a flat zip archive built from a bundle key.
//
// All file records are validated before the first body byte is written. Entry
// names come verbatim from the bundle key, entries are stored (not deflated),
// and the response carries no Content-Length so the transfer is chunked.
func DownloadBundleKey(rail miso.Rail, w http.ResponseWriter, bundleKey string) error {
	ok, bundle := ResolveBundleKey(rail, bundleKey)
	if !ok {
		return ErrFileNotFound
	}

	ids := make([]string, 0, len(bundle.Items))
	for _, it := range bundle.Items {
		ids = append(ids, it.FileId)
	}
	allNormal, err := CheckAllNormalFiles(ids)
	if err != nil {
		return err
	}
	if !allNormal {
		return ErrFileNotFound.WithInternalMsg("not all bundled files are normal")
	}

	// load the records in one query; only os.Stat stays per entry
	byId, err := findDFiles(ids)
	if err != nil {
		return ErrUnknownError.WithInternalMsg("findDFiles failed, %v", err)
	}
	dfs := make([]DFile, 0, len(bundle.Items))
	for _, it := range bundle.Items {
		df, found := byId[it.FileId]
		if !found {
			return ErrFileNotFound.WithInternalMsg("file record not found, fileId: %v", it.FileId)
		}
		if df.IsDeleted() {
			return ErrFileDeleted
		}
		// fail before the body starts instead of streaming a truncated archive
		if _, err := os.Stat(df.StoragePath()); err != nil {
			return ErrFileNotFound.WithInternalMsg("stored bytes missing, fileId: %v, %v", it.FileId, err)
		}
		dfs = append(dfs, df)
	}

	if err := RefreshBundleKeyExp(rail, bundleKey); err != nil {
		return err
	}

	headers := w.Header()
	headers.Set("Content-Type", "application/zip")
	headers.Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{
		"filename": bundleArchiveName(bundle.Name),
	}))

	zw := zip.NewWriter(w)
	flushed := false
	for i, it := range bundle.Items {
		ew, err := zw.CreateHeader(&zip.FileHeader{Name: it.Name, Method: zip.Store})
		if err != nil {
			return errs.NewErrf("failed to create zip entry %v, %v", it.Name, err)
		}
		// Only once the first entry header exists is there a body worth
		// committing. Until this point a failure can still change the
		// response status. Content-Length stays absent, so the response is
		// chunked.
		if !flushed {
			// the zip writer buffers, so flush it into the response first
			if err := zw.Flush(); err != nil {
				return errs.NewErrf("failed to flush zip header for %v, %v", it.Name, err)
			}
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			flushed = true
		}
		if err := TransferFile(rail, ew, dfs[i], ZeroByteRange()); err != nil {
			return errs.Wrapf(err, "failed to write zip entry %v", it.Name)
		}
	}
	if !flushed {
		// empty archive: commit the headers before the central directory
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
	return zw.Close()
}
