package fstore

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/mini-fstore/internal/config"
	"github.com/curtisnewbie/miso/miso"
)

// setupStreamFile writes content to the storage dir and returns a matching DFile.
func setupStreamFile(t *testing.T, name string, content []byte) DFile {
	t.Helper()
	miso.SetProp(config.PropStorageDir, t.TempDir())
	fileId := "file_serve_test"
	if err := os.WriteFile(GenStoragePath(fileId), content, os.ModePerm); err != nil {
		t.Fatalf("failed to write test file, %v", err)
	}
	return DFile{
		FileId: fileId,
		Name:   name,
		Size:   int64(len(content)),
		Status: api.FileStatusNormal,
	}
}

func serveFileContent(t *testing.T, ff DFile, method, rangeHeader string) *http.Response {
	t.Helper()
	req := httptest.NewRequest(method, "/file/stream", nil)
	if rangeHeader != "" {
		req.Header.Set("Range", rangeHeader)
	}
	rr := httptest.NewRecorder()
	if err := ServeFileContent(rr, req, ff); err != nil {
		t.Fatalf("ServeFileContent returned error, %v", err)
	}
	return rr.Result()
}

func readBody(t *testing.T, res *http.Response) []byte {
	t.Helper()
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read response body, %v", err)
	}
	return b
}

func TestServeFileContent(t *testing.T) {
	content := []byte("0123456789abcdefghij") // 20 bytes
	n := len(content)

	t.Run("no range returns whole file", func(t *testing.T) {
		ff := setupStreamFile(t, "v.mp4", content)
		res := serveFileContent(t, ff, http.MethodGet, "")

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}
		if cr := res.Header.Get("Content-Range"); cr != "" {
			t.Errorf("Content-Range = %q, want empty", cr)
		}
		if cl := res.Header.Get("Content-Length"); cl != "20" {
			t.Errorf("Content-Length = %q, want %q", cl, "20")
		}
		if ar := res.Header.Get("Accept-Ranges"); ar != "bytes" {
			t.Errorf("Accept-Ranges = %q, want %q", ar, "bytes")
		}
		if ct := res.Header.Get("Content-Type"); ct != "video/mp4" {
			t.Errorf("Content-Type = %q, want %q", ct, "video/mp4")
		}
		if lm := res.Header.Get("Last-Modified"); lm == "" {
			t.Errorf("Last-Modified is empty, want a value")
		}
		if body := readBody(t, res); !bytes.Equal(body, content) {
			t.Errorf("body = %q, want %q", body, content)
		}
	})

	t.Run("open-ended range is not clamped", func(t *testing.T) {
		ff := setupStreamFile(t, "v.mp4", content)
		res := serveFileContent(t, ff, http.MethodGet, "bytes=0-")

		if res.StatusCode != http.StatusPartialContent {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusPartialContent)
		}
		if cr := res.Header.Get("Content-Range"); cr != "bytes 0-19/20" {
			t.Errorf("Content-Range = %q, want %q", cr, "bytes 0-19/20")
		}
		if cl := res.Header.Get("Content-Length"); cl != "20" {
			t.Errorf("Content-Length = %q, want %q", cl, "20")
		}
		if body := readBody(t, res); !bytes.Equal(body, content) {
			t.Errorf("body = %q, want %q", body, content)
		}
	})

	t.Run("bounded range returns slice", func(t *testing.T) {
		ff := setupStreamFile(t, "v.mp4", content)
		res := serveFileContent(t, ff, http.MethodGet, "bytes=5-9")

		if res.StatusCode != http.StatusPartialContent {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusPartialContent)
		}
		if cr := res.Header.Get("Content-Range"); cr != "bytes 5-9/20" {
			t.Errorf("Content-Range = %q, want %q", cr, "bytes 5-9/20")
		}
		want := content[5:10]
		if body := readBody(t, res); !bytes.Equal(body, want) {
			t.Errorf("body = %q, want %q", body, want)
		}
	})

	t.Run("suffix range returns last bytes", func(t *testing.T) {
		ff := setupStreamFile(t, "v.mp4", content)
		res := serveFileContent(t, ff, http.MethodGet, "bytes=-10")

		if res.StatusCode != http.StatusPartialContent {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusPartialContent)
		}
		if cr := res.Header.Get("Content-Range"); cr != "bytes 10-19/20" {
			t.Errorf("Content-Range = %q, want %q", cr, "bytes 10-19/20")
		}
		want := content[n-10:]
		if body := readBody(t, res); !bytes.Equal(body, want) {
			t.Errorf("body = %q, want %q", body, want)
		}
	})

	t.Run("range beyond EOF is unsatisfiable", func(t *testing.T) {
		ff := setupStreamFile(t, "v.mp4", content)
		res := serveFileContent(t, ff, http.MethodGet, "bytes=20-")

		if res.StatusCode != http.StatusRequestedRangeNotSatisfiable {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusRequestedRangeNotSatisfiable)
		}
	})

	t.Run("head request has empty body and content length", func(t *testing.T) {
		ff := setupStreamFile(t, "v.mp4", content)
		res := serveFileContent(t, ff, http.MethodHead, "")

		if res.StatusCode != http.StatusOK {
			t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
		}
		if cl := res.Header.Get("Content-Length"); cl != "20" {
			t.Errorf("Content-Length = %q, want %q", cl, "20")
		}
		if body := readBody(t, res); len(body) != 0 {
			t.Errorf("body = %q, want empty", body)
		}
	})

	t.Run("mov file has quicktime content type", func(t *testing.T) {
		ff := setupStreamFile(t, "clip.mov", content)
		res := serveFileContent(t, ff, http.MethodGet, "")

		if ct := res.Header.Get("Content-Type"); ct != "video/quicktime" {
			t.Errorf("Content-Type = %q, want %q", ct, "video/quicktime")
		}
	})
}
