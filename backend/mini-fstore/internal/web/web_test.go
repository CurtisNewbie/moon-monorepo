package web

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/mini-fstore/internal/config"
	"github.com/curtisnewbie/mini-fstore/internal/fstore"
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/rabbit"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/miso"
)

func preTest(t *testing.T) {
	t.Helper()
	c := miso.EmptyRail()
	ag := []string{"configFile=../../conf.yml"}
	miso.DefaultReadConfig(ag, c)
	miso.ConfigureLogging()
	miso.SetProp(config.PropStorageDir, t.TempDir())
	miso.SetProp(config.PropTrashDir, t.TempDir())
	if err := mysql.InitMySQLFromProp(c); err != nil {
		t.Fatal(err)
	}
	if _, err := redis.InitRedisFromProp(c); err != nil {
		t.Fatal(err)
	}
	if err := rabbit.StartRabbitMqClient(c); err != nil {
		t.Fatal(err)
	}
}

// statusRecorder records explicit status writes and body bytes separately, so a
// status change after the body started stays visible.
type statusRecorder struct {
	header http.Header
	status []int
	body   bytes.Buffer
}

func newStatusRecorder() *statusRecorder {
	return &statusRecorder{header: http.Header{}}
}

func (r *statusRecorder) Header() http.Header { return r.header }

func (r *statusRecorder) WriteHeader(code int) { r.status = append(r.status, code) }

func (r *statusRecorder) Write(p []byte) (int, error) { return r.body.Write(p) }

func createFileRecord(t *testing.T, fileId, name string, size int64) {
	t.Helper()
	if err := fstore.CreateFileRec(miso.EmptyRail(), fstore.CreateFile{
		FileId: fileId,
		Name:   name,
		Size:   size,
	}); err != nil {
		t.Fatalf("failed to create file record, %v", err)
	}
	t.Cleanup(func() {
		mysql.GetMySQL().Exec("delete from file where file_id = ?", fileId)
	})
}

func TestBundleDownloadUnknownKeyReturnsNotFound(t *testing.T) {
	preTest(t)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/file/bundle?key=unknown-bundle-key", nil)

	downloadBundle(miso.EmptyRail(), rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
	if rr.Body.Len() != 0 {
		t.Errorf("body has %d bytes, want none", rr.Body.Len())
	}
}

func TestBundleDownloadMidStreamFailureDoesNotChangeStatus(t *testing.T) {
	preTest(t)

	fileId := fstore.GenFileId()
	// a directory at the storage path passes the stat check, then fails while
	// its bytes are read, i.e. after the archive body has started
	if err := os.MkdirAll(fstore.GenStoragePath(fileId), os.ModePerm); err != nil {
		t.Fatalf("failed to create storage path, %v", err)
	}
	createFileRecord(t, fileId, "broken.bin", 10)

	key, err := fstore.RandBundleKey(miso.EmptyRail(), "dir", []api.BundleItem{
		{FileId: fileId, Name: "broken.bin"},
	})
	if err != nil {
		t.Fatalf("RandBundleKey returned error, %v", err)
	}

	w := newStatusRecorder()
	req := httptest.NewRequest(http.MethodGet, "/file/bundle?key="+url.QueryEscape(key), nil)

	downloadBundle(miso.EmptyRail(), w, req)

	if len(w.status) != 0 {
		t.Fatalf("handler wrote status %v after the body started, want none", w.status)
	}
	if w.body.Len() == 0 {
		t.Fatal("expected body bytes before the mid-stream failure")
	}
}

func TestBundleDownloadPreBodyFailureReturnsNotFoundWithNoBody(t *testing.T) {
	preTest(t)

	// no bytes are written for this id, so the request fails before the body
	fileId := fstore.GenFileId()
	createFileRecord(t, fileId, "missing.bin", 10)

	key, err := fstore.RandBundleKey(miso.EmptyRail(), "dir", []api.BundleItem{
		{FileId: fileId, Name: "missing.bin"},
	})
	if err != nil {
		t.Fatalf("RandBundleKey returned error, %v", err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/file/bundle?key="+url.QueryEscape(key), nil)

	downloadBundle(miso.EmptyRail(), rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
	if rr.Body.Len() != 0 {
		t.Errorf("body has %d bytes, want none", rr.Body.Len())
	}
}
