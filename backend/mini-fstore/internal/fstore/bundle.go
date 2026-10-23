package fstore

import (
	"archive/zip"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"os"
	"time"

	"github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/json"
	"github.com/curtisnewbie/miso/util/randutil"
)

const bundleKeyPrefix = "fstore:bundle:key:"

// Lifetime of a bundle key. It is refreshed on each successful read, matching
// the single-file key.
var bundleKeyTTL = 30 * time.Minute

// BundleItem is one file that goes into a bundle archive.
type BundleItem struct {
	FileId string `json:"fileId"`
	Name   string `json:"name"`
}

// CachedBundle is the payload held by a bundle key: the archive name plus the
// ordered list of entries.
type CachedBundle struct {
	Name  string       `json:"name"`
	Items []BundleItem `json:"items"`
}

// RandBundleKey generates a temporary key for an ordered list of files.
func RandBundleKey(rail miso.Rail, name string, items []BundleItem) (string, error) {
	if items == nil {
		items = []BundleItem{}
	}
	bk := randutil.ERand(30)
	sby, em := json.WriteJson(CachedBundle{Name: name, Items: items})
	if em != nil {
		return "", fmt.Errorf("failed to marshal to CachedBundle, %v", em)
	}
	c := redis.GetRedis().Set(rail.Context(), bundleKeyPrefix+bk, string(sby), bundleKeyTTL)
	return bk, c.Err()
}

// RefreshBundleKeyExp refreshes the bundle key's expiration.
func RefreshBundleKeyExp(rail miso.Rail, bundleKey string) error {
	c := redis.GetRedis().Expire(rail.Context(), bundleKeyPrefix+bundleKey, bundleKeyTTL)
	if c.Err() != nil {
		rail.Warnf("Failed to refresh bundle key expiration, bundleKey: %v, %v", bundleKey, c.Err())
		return fmt.Errorf("failed to refresh key expiration, %v", c.Err())
	}
	return nil
}

// ResolveBundleKey resolves the CachedBundle for the given bundle key.
func ResolveBundleKey(rail miso.Rail, bundleKey string) (bool, CachedBundle) {
	var cb CachedBundle
	c := redis.GetRedis().Get(rail.Context(), bundleKeyPrefix+bundleKey)
	if c.Err() != nil {
		if errors.Is(c.Err(), redis.Nil) {
			rail.Infof("BundleKey not found, %v", bundleKey)
		} else {
			rail.Errorf("Failed to find bundleKey, %v", c.Err())
		}
		return false, cb
	}

	eu := json.ParseJson([]byte(c.Val()), &cb)
	if eu != nil {
		rail.Errorf("Failed to unmarshal bundleKey, %s, %v", bundleKey, eu)
		return false, cb
	}
	return true, cb
}

// allFilesNormal checks in one batch that every distinct file id refers to a
// normal file record.
//
// CheckAllNormalFiles cannot be used here: it relies on slutil.Distinct, which
// in the miso version used by this service returns a slice padded with empty
// strings, so the count can never match the id list.
func allFilesNormal(ids []string) (bool, error) {
	distinct := make([]string, 0, len(ids))
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		distinct = append(distinct, id)
	}
	if len(distinct) == 0 {
		return true, nil
	}

	var cnt int64
	t := mysql.GetMySQL().
		Table("file").
		Where("file_id in ?", distinct).
		Where("status = ?", api.FileStatusNormal).
		Count(&cnt)
	if t.Error != nil {
		return false, fmt.Errorf("failed to count normal files, %w", t.Error)
	}
	return cnt == int64(len(distinct)), nil
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
	if len(ids) > 0 {
		allNormal, err := allFilesNormal(ids)
		if err != nil {
			return err
		}
		if !allNormal {
			return ErrFileNotFound.WithInternalMsg("not all bundled files are normal")
		}
	}

	dfs := make([]DFile, 0, len(bundle.Items))
	for _, it := range bundle.Items {
		df, err := findDFile(it.FileId)
		if err != nil {
			return ErrFileNotFound.WithInternalMsg("findDFile failed, fileId: %v, %v", it.FileId, err)
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
		"filename": bundle.Name + ".zip",
	}))
	// Content-Length is deliberately absent: the archive size is unknown until
	// the last entry is written, so the response is chunked. Flushing here
	// commits the headers and starts the transfer before the first entry is
	// read.
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	zw := zip.NewWriter(w)
	for i, it := range bundle.Items {
		ew, err := zw.CreateHeader(&zip.FileHeader{Name: it.Name, Method: zip.Store})
		if err != nil {
			return fmt.Errorf("failed to create zip entry %v, %v", it.Name, err)
		}
		if err := TransferFile(rail, ew, dfs[i], ZeroByteRange()); err != nil {
			return fmt.Errorf("failed to write zip entry %v, %w", it.Name, err)
		}
	}
	return zw.Close()
}
