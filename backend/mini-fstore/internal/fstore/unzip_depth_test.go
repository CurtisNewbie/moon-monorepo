package fstore

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/miso/errs"
)

func TestDeepestZipEntryDepth(t *testing.T) {
	cases := []struct {
		name     string
		names    []string
		want     int
		wantName string
	}{
		{
			name:     "flat archive",
			names:    []string{"a.txt", "b.txt"},
			want:     1,
			wantName: "a.txt",
		},
		{
			name:     "nested paths keep the deepest",
			names:    []string{"top/a.txt", "top/doc/api/index.html", "top/b.txt"},
			want:     4,
			wantName: "top/doc/api/index.html",
		},
		{
			name:     "first entry wins a tie",
			names:    []string{"top/a/one.txt", "top/b/two.txt"},
			want:     3,
			wantName: "top/a/one.txt",
		},
		{
			name:     "directory entries alone have no depth",
			names:    []string{"Documents/", "Documents/doc/", "Documents/doc/api/"},
			want:     0,
			wantName: "",
		},
		{
			name:     "directory entries are ignored even when deeper than files",
			names:    []string{"Documents/", "Documents/a/b/c/d/", "Documents/a.txt"},
			want:     2,
			wantName: "Documents/a.txt",
		},
		{
			name:     "empty segments are ignored",
			names:    []string{"/Documents//doc///a.txt"},
			want:     3,
			wantName: "/Documents//doc///a.txt",
		},
		{
			name:     "parent segment counts as a segment",
			names:    []string{"../etc/passwd"},
			want:     3,
			wantName: "../etc/passwd",
		},
		{
			name:     "backslash separates segments",
			names:    []string{`a\b\c.txt`},
			want:     3,
			wantName: `a\b\c.txt`,
		},
		{
			name:     "backslash and slash separate segments",
			names:    []string{`a\b/c.txt`},
			want:     3,
			wantName: `a\b/c.txt`,
		},
		{
			name:     "empty archive",
			names:    nil,
			want:     0,
			wantName: "",
		},
		{
			name:     "empty names",
			names:    []string{"", "/"},
			want:     0,
			wantName: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, gotName := deepestZipEntryDepth(c.names)
			if got != c.want {
				t.Errorf("deepestZipEntryDepth(%v) depth = %d, want %d", c.names, got, c.want)
			}
			if gotName != c.wantName {
				t.Errorf("deepestZipEntryDepth(%v) name = %q, want %q", c.names, gotName, c.wantName)
			}
		})
	}
}

type testZipEntry struct {
	Name string
	Body []byte
}

// writeTestZipAt writes a real zip archive to path, with directory entries for
// names ending in a slash.
func writeTestZipAt(t *testing.T, path string, entries []testZipEntry) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create test zip, %v", err)
	}
	zw := zip.NewWriter(f)
	for _, e := range entries {
		w, err := zw.Create(e.Name)
		if err != nil {
			t.Fatalf("failed to create zip entry %q, %v", e.Name, err)
		}
		if strings.HasSuffix(e.Name, "/") {
			continue
		}
		if _, err := w.Write(e.Body); err != nil {
			t.Fatalf("failed to write zip entry %q, %v", e.Name, err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("failed to close test zip writer, %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("failed to close test zip, %v", err)
	}
}

func writeTestZip(t *testing.T, entries []testZipEntry) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "archive.zip")
	writeTestZipAt(t, p, entries)
	return p
}

func TestCheckZipDepth(t *testing.T) {
	deep := writeTestZip(t, []testZipEntry{
		{Name: "top/"},
		{Name: "top/a/"},
		{Name: "top/a/b/"},
		{Name: "top/a/b/c.txt", Body: []byte("c")},
	})

	// deepest file entry is top/a/b/c.txt, four segments from the target
	err := checkZipDepth(deep, 3)
	if m, ok := err.(*errs.MisoErr); !ok || m.Code() != api.ZipTooDeep {
		t.Fatalf("checkZipDepth(depth 4, limit 3) = %v, want code %q", err, api.ZipTooDeep)
	}
	if !strings.Contains(err.Error(), "3") {
		t.Errorf("refusal message %q does not name the limit", err.Error())
	}
	if !strings.Contains(err.Error(), "top/a/b/c.txt") {
		t.Errorf("refusal message %q does not name the deepest entry", err.Error())
	}

	if err := checkZipDepth(deep, 4); err != nil {
		t.Errorf("checkZipDepth(depth 4, limit 4) = %v, want nil", err)
	}

	flat := writeTestZip(t, []testZipEntry{{Name: "a.txt", Body: []byte("a")}})
	if err := checkZipDepth(flat, 1); err != nil {
		t.Errorf("checkZipDepth(depth 1, limit 1) = %v, want nil", err)
	}
}

func TestCheckZipDepthUnreadableArchive(t *testing.T) {
	p := filepath.Join(t.TempDir(), "broken.zip")
	if err := os.WriteFile(p, []byte("this is not a zip archive"), 0o644); err != nil {
		t.Fatalf("failed to write broken archive, %v", err)
	}

	err := checkZipDepth(p, api.DefaultUnzipMaxDepth)
	if m, ok := err.(*errs.MisoErr); !ok || m.Code() != api.IllegalFormat {
		t.Fatalf("checkZipDepth(broken) = %v, want code %q", err, api.IllegalFormat)
	}
}
