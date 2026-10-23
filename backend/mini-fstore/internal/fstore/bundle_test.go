package fstore

import (
	"archive/zip"
	"bytes"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/mini-fstore/internal/config"
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/miso"
)

// setupBundleEnv prepares the live dependencies (MySQL, Redis) and a temporary
// storage directory for the bundle archive tests.
func setupBundleEnv(t *testing.T) {
	t.Helper()
	preTest(t)
	miso.SetProp(config.PropStorageDir, t.TempDir())
}

// createBundleFile stores content in the temporary storage directory and
// creates a matching normal file record.
func createBundleFile(t *testing.T, name string, content []byte) string {
	t.Helper()
	fileId := GenFileId()
	if err := os.WriteFile(GenStoragePath(fileId), content, os.ModePerm); err != nil {
		t.Fatalf("failed to write test file, %v", err)
	}
	if err := CreateFileRec(miso.EmptyRail(), CreateFile{
		FileId: fileId,
		Name:   name,
		Size:   int64(len(content)),
	}); err != nil {
		t.Fatalf("failed to create file record, %v", err)
	}
	t.Cleanup(func() {
		mysql.GetMySQL().Exec("delete from file where file_id = ?", fileId)
	})
	return fileId
}

func mintBundle(t *testing.T, name string, items []api.BundleItem) string {
	t.Helper()
	k, err := RandBundleKey(miso.EmptyRail(), name, items)
	if err != nil {
		t.Fatalf("RandBundleKey returned error, %v", err)
	}
	if k == "" {
		t.Fatal("RandBundleKey returned an empty key")
	}
	return k
}

// downloadBundleKey drives the archive through an httptest recorder, returning
// whatever the function produced and its error.
func downloadBundleKey(key string) (*httptest.ResponseRecorder, error) {
	rr := httptest.NewRecorder()
	err := DownloadBundleKey(miso.EmptyRail(), rr, key)
	return rr, err
}

// fetchBundleArchive serves the archive through a real httptest server so the
// transport-level response (content length, chunking) can be observed.
func fetchBundleArchive(t *testing.T, key string) (*http.Response, []byte) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := DownloadBundleKey(miso.EmptyRail(), w, key); err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	res, err := http.Get(srv.URL + "/file/bundle?key=" + url.QueryEscape(key))
	if err != nil {
		t.Fatalf("failed to fetch archive, %v", err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatalf("failed to read archive body, %v", err)
	}
	return res, body
}

func readZip(t *testing.T, b []byte) *zip.Reader {
	t.Helper()
	zr, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		t.Fatalf("response body is not a readable zip, %v", err)
	}
	return zr
}

func TestRandBundleKeyResolvesItemsInOrder(t *testing.T) {
	setupBundleEnv(t)

	items := []api.BundleItem{
		{FileId: "file_a", Name: "first.txt"},
		{FileId: "file_b", Name: "second.txt"},
		{FileId: "file_c", Name: "third.txt"},
	}
	k := mintBundle(t, "my dir", items)

	ok, resolved := ResolveBundleKey(miso.EmptyRail(), k)
	if !ok {
		t.Fatal("bundle key did not resolve")
	}
	if resolved.Name != "my dir" {
		t.Errorf("archive name = %q, want %q", resolved.Name, "my dir")
	}
	if len(resolved.Items) != len(items) {
		t.Fatalf("resolved %d items, want %d", len(resolved.Items), len(items))
	}
	for i := range items {
		if resolved.Items[i] != items[i] {
			t.Errorf("item[%d] = %+v, want %+v", i, resolved.Items[i], items[i])
		}
	}
}

func TestDownloadBundleKeyStreamsFlatStoredZip(t *testing.T) {
	setupBundleEnv(t)

	contentA := []byte("first file bytes")
	contentB := []byte("second file bytes, a bit longer")
	a := createBundleFile(t, "a.txt", contentA)
	b := createBundleFile(t, "b.txt", contentB)
	k := mintBundle(t, "my dir", []api.BundleItem{
		{FileId: a, Name: "a.txt"},
		{FileId: b, Name: "b.txt"},
	})

	res, body := fetchBundleArchive(t, k)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusOK)
	}
	if ct := res.Header.Get("Content-Type"); ct != "application/zip" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/zip")
	}
	if cl := res.Header.Get("Content-Length"); cl != "" {
		t.Errorf("Content-Length = %q, want absent", cl)
	}
	if res.ContentLength != -1 {
		t.Errorf("ContentLength = %d, want -1 (unknown)", res.ContentLength)
	}
	chunked := false
	for _, te := range res.TransferEncoding {
		if strings.EqualFold(te, "chunked") {
			chunked = true
		}
	}
	if !chunked {
		t.Errorf("TransferEncoding = %v, want chunked", res.TransferEncoding)
	}

	mediaType, params, err := mime.ParseMediaType(res.Header.Get("Content-Disposition"))
	if err != nil {
		t.Fatalf("failed to parse Content-Disposition, %v", err)
	}
	if mediaType != "attachment" {
		t.Errorf("disposition type = %q, want %q", mediaType, "attachment")
	}
	if params["filename"] != "my dir.zip" {
		t.Errorf("disposition filename = %q, want %q", params["filename"], "my dir.zip")
	}

	zr := readZip(t, body)
	if len(zr.File) != 2 {
		t.Fatalf("archive has %d entries, want 2", len(zr.File))
	}

	wantNames := []string{"a.txt", "b.txt"}
	wantContent := [][]byte{contentA, contentB}
	for i, zf := range zr.File {
		if zf.Name != wantNames[i] {
			t.Errorf("entry[%d].Name = %q, want %q", i, zf.Name, wantNames[i])
		}
		if strings.Contains(zf.Name, "/") {
			t.Errorf("entry[%d].Name = %q, want a bare filename at the archive root", i, zf.Name)
		}
		if zf.FileInfo().IsDir() {
			t.Errorf("entry[%d] is a directory entry, want only file entries", i)
		}
		if zf.Method != zip.Store {
			t.Errorf("entry[%d].Method = %d, want zip.Store (%d)", i, zf.Method, zip.Store)
		}

		rc, err := zf.Open()
		if err != nil {
			t.Fatalf("failed to open entry %q, %v", zf.Name, err)
		}
		got, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			t.Fatalf("failed to read entry %q, %v", zf.Name, err)
		}
		if !bytes.Equal(got, wantContent[i]) {
			t.Errorf("entry %q contents = %q, want %q", zf.Name, got, wantContent[i])
		}
	}
}

func TestDownloadBundleKeyNonAsciiArchiveName(t *testing.T) {
	setupBundleEnv(t)

	content := []byte("pdf bytes")
	fileId := createBundleFile(t, "报告.pdf", content)
	k := mintBundle(t, "资料", []api.BundleItem{{FileId: fileId, Name: "报告.pdf"}})

	res, body := fetchBundleArchive(t, k)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusOK)
	}
	_, params, err := mime.ParseMediaType(res.Header.Get("Content-Disposition"))
	if err != nil {
		t.Fatalf("failed to parse Content-Disposition, %v", err)
	}
	if params["filename"] != "资料.zip" {
		t.Errorf("disposition filename = %q, want %q", params["filename"], "资料.zip")
	}

	// the entry name round-trips inside the archive body as well
	zr := readZip(t, body)
	if len(zr.File) != 1 {
		t.Fatalf("archive has %d entries, want 1", len(zr.File))
	}
	if zr.File[0].Name != "报告.pdf" {
		t.Errorf("entry name = %q, want %q", zr.File[0].Name, "报告.pdf")
	}
	rc, err := zr.File[0].Open()
	if err != nil {
		t.Fatalf("failed to open entry, %v", err)
	}
	got, err := io.ReadAll(rc)
	rc.Close()
	if err != nil {
		t.Fatalf("failed to read entry, %v", err)
	}
	if !bytes.Equal(got, content) {
		t.Errorf("entry contents = %q, want %q", got, content)
	}
}

func TestDownloadBundleKeyUsesEntryNameFromBundle(t *testing.T) {
	setupBundleEnv(t)

	content := []byte("bytes of the stored file")
	fileId := createBundleFile(t, "stored-name.txt", content)

	// the bundle key, not the file record, decides the entry name
	k := mintBundle(t, "dir", []api.BundleItem{{FileId: fileId, Name: "renamed-entry.txt"}})

	rr, err := downloadBundleKey(k)
	if err != nil {
		t.Fatalf("DownloadBundleKey returned error, %v", err)
	}
	zr := readZip(t, rr.Body.Bytes())
	if len(zr.File) != 1 {
		t.Fatalf("archive has %d entries, want 1", len(zr.File))
	}
	if zr.File[0].Name != "renamed-entry.txt" {
		t.Errorf("entry name = %q, want %q", zr.File[0].Name, "renamed-entry.txt")
	}
	rc, err := zr.File[0].Open()
	if err != nil {
		t.Fatalf("failed to open entry, %v", err)
	}
	got, err := io.ReadAll(rc)
	rc.Close()
	if err != nil {
		t.Fatalf("failed to read entry, %v", err)
	}
	if !bytes.Equal(got, content) {
		t.Errorf("entry contents = %q, want %q", got, content)
	}
}

func TestBundleArchiveNameDoesNotDoubleZipSuffix(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{"my dir", "my dir.zip"},
		{"foo.zip", "foo.zip"},
		{"foo.ZIP", "foo.ZIP"},
		{"archive.zip.zip", "archive.zip.zip"},
	}
	for _, c := range cases {
		if got := bundleArchiveName(c.name); got != c.want {
			t.Errorf("bundleArchiveName(%q) = %q, want %q", c.name, got, c.want)
		}
	}
}

func TestRandBundleKeyRejectsTooManyItems(t *testing.T) {
	setupBundleEnv(t)

	items := make([]api.BundleItem, maxBundleItems+1)
	for i := range items {
		items[i] = api.BundleItem{FileId: "file_" + strconv.Itoa(i), Name: "entry_" + strconv.Itoa(i)}
	}

	k, err := RandBundleKey(miso.EmptyRail(), "too many", items)
	if err == nil {
		t.Fatal("RandBundleKey should reject an oversized bundle")
	}
	if k != "" {
		t.Errorf("key = %q, want empty", k)
	}
	if !strings.Contains(err.Error(), "too many") {
		t.Errorf("error = %q, want it to explain the cap", err)
	}
}

func TestRandBundleKeyAcceptsCap(t *testing.T) {
	setupBundleEnv(t)

	items := make([]api.BundleItem, maxBundleItems)
	for i := range items {
		items[i] = api.BundleItem{FileId: "file_" + strconv.Itoa(i), Name: "entry_" + strconv.Itoa(i)}
	}

	if _, err := RandBundleKey(miso.EmptyRail(), "at cap", items); err != nil {
		t.Fatalf("RandBundleKey should accept the cap, %v", err)
	}
}

func TestDownloadBundleKeySingleEntry(t *testing.T) {
	setupBundleEnv(t)

	content := []byte("only one")
	fileId := createBundleFile(t, "only.txt", content)
	k := mintBundle(t, "one", []api.BundleItem{{FileId: fileId, Name: "only.txt"}})

	rr, err := downloadBundleKey(k)
	if err != nil {
		t.Fatalf("DownloadBundleKey returned error, %v", err)
	}
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	zr := readZip(t, rr.Body.Bytes())
	if len(zr.File) != 1 {
		t.Fatalf("archive has %d entries, want 1", len(zr.File))
	}
	if zr.File[0].Name != "only.txt" {
		t.Errorf("entry name = %q, want %q", zr.File[0].Name, "only.txt")
	}
}

func TestDownloadBundleKeyEmptyListYieldsEmptyZip(t *testing.T) {
	setupBundleEnv(t)

	k := mintBundle(t, "empty dir", []api.BundleItem{})

	rr, err := downloadBundleKey(k)
	if err != nil {
		t.Fatalf("DownloadBundleKey returned error, %v", err)
	}
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	zr := readZip(t, rr.Body.Bytes())
	if len(zr.File) != 0 {
		t.Errorf("archive has %d entries, want 0", len(zr.File))
	}
}

func TestDownloadBundleKeyRejectsUnknownFileBeforeBody(t *testing.T) {
	setupBundleEnv(t)

	k := mintBundle(t, "dir", []api.BundleItem{{FileId: "file_does_not_exist", Name: "ghost.txt"}})

	rr, err := downloadBundleKey(k)
	if err == nil {
		t.Fatal("DownloadBundleKey should fail for an unknown file id")
	}
	if rr.Body.Len() != 0 {
		t.Errorf("wrote %d body bytes, want none", rr.Body.Len())
	}
	if ct := rr.Header().Get("Content-Type"); ct != "" {
		t.Errorf("Content-Type = %q, want it unset before validation passes", ct)
	}
}

func TestDownloadBundleKeyRejectsNonNormalFileBeforeBody(t *testing.T) {
	setupBundleEnv(t)

	fileId := createBundleFile(t, "deleted.txt", []byte("gone"))
	if err := mysql.GetMySQL().Exec("update file set status = ? where file_id = ?", api.FileStatusLogicDel, fileId).Error; err != nil {
		t.Fatalf("failed to mark file deleted, %v", err)
	}
	k := mintBundle(t, "dir", []api.BundleItem{{FileId: fileId, Name: "deleted.txt"}})

	rr, err := downloadBundleKey(k)
	if err == nil {
		t.Fatal("DownloadBundleKey should fail for a non-normal file record")
	}
	if rr.Body.Len() != 0 {
		t.Errorf("wrote %d body bytes, want none", rr.Body.Len())
	}
	if ct := rr.Header().Get("Content-Type"); ct != "" {
		t.Errorf("Content-Type = %q, want it unset before validation passes", ct)
	}
}

func TestDownloadBundleKeyRejectsMissingStoredBytesBeforeBody(t *testing.T) {
	setupBundleEnv(t)

	// create a normal record without writing any bytes to the storage directory
	fileId := GenFileId()
	if err := CreateFileRec(miso.EmptyRail(), CreateFile{
		FileId: fileId,
		Name:   "missing.bin",
		Size:   10,
	}); err != nil {
		t.Fatalf("failed to create file record, %v", err)
	}
	t.Cleanup(func() {
		mysql.GetMySQL().Exec("delete from file where file_id = ?", fileId)
	})

	k := mintBundle(t, "dir", []api.BundleItem{{FileId: fileId, Name: "missing.bin"}})

	rr, err := downloadBundleKey(k)
	if err == nil {
		t.Fatal("DownloadBundleKey should fail when the stored bytes are missing")
	}
	if rr.Body.Len() != 0 {
		t.Errorf("wrote %d body bytes, want none", rr.Body.Len())
	}
}

func TestDownloadBundleKeyRejectsUnknownKey(t *testing.T) {
	setupBundleEnv(t)

	rr, err := downloadBundleKey("not_a_real_bundle_key")
	if err == nil {
		t.Fatal("DownloadBundleKey should fail for an unknown bundle key")
	}
	if rr.Body.Len() != 0 {
		t.Errorf("wrote %d body bytes, want none", rr.Body.Len())
	}
}

func TestBundleKeyExpiresAndRefreshes(t *testing.T) {
	setupBundleEnv(t)

	oldTTL := bundleKeyTTL
	bundleKeyTTL = 3 * time.Second
	t.Cleanup(func() { bundleKeyTTL = oldTTL })

	step := 300 * time.Millisecond
	// more than the ttl, in small steps so a slow run only spends a little
	// extra wall clock instead of flaking
	expirySteps := 12

	fileId := createBundleFile(t, "kept.txt", []byte("still here"))

	k := mintBundle(t, "dir", []api.BundleItem{{FileId: fileId, Name: "kept.txt"}})

	// a successful read refreshes the lifetime
	for i := 0; i < 3; i++ {
		time.Sleep(step)
		if _, err := downloadBundleKey(k); err != nil {
			t.Fatalf("bundle key should still resolve before expiration, %v", err)
		}
	}

	// without a read, the key expires
	k2 := mintBundle(t, "dir", []api.BundleItem{{FileId: fileId, Name: "kept.txt"}})
	time.Sleep(time.Duration(expirySteps) * step)
	if _, err := downloadBundleKey(k2); err == nil {
		t.Fatal("bundle key should have expired after the ttl of no reads")
	}
}
