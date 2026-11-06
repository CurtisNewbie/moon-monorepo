package fstore

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/miso/errs"
)

func zipDirNode(name string, children ...*api.ZipTreeNode) *api.ZipTreeNode {
	if children == nil {
		children = []*api.ZipTreeNode{}
	}
	return &api.ZipTreeNode{Name: name, Dir: true, Children: children}
}

func zipFileNode(name string, size int64) *api.ZipTreeNode {
	return &api.ZipTreeNode{Name: name, Size: size}
}

func zipTreeResp(root *api.ZipTreeNode, included int, reasons ...string) api.ZipTreeResp {
	if reasons == nil {
		reasons = []string{}
	}
	return api.ZipTreeResp{
		Tree:             *root,
		Truncated:        len(reasons) > 0,
		TruncatedReasons: reasons,
		IncludedEntries:  included,
	}
}

// buildZipTree assembles one tree from a flat entry list for the table tests.
func buildZipTree(entries []zipTreeEntry, maxDepth, maxEntries int) api.ZipTreeResp {
	b := newZipTreeBuilder(maxDepth, maxEntries)
	for _, e := range entries {
		if b.add(e) {
			break
		}
	}
	return b.result()
}

// renderZipTree prints a tree for failure messages, where %+v would only show
// pointer addresses.
func renderZipTree(n *api.ZipTreeNode) string {
	var sb strings.Builder
	var walk func(n *api.ZipTreeNode, depth int)
	walk = func(n *api.ZipTreeNode, depth int) {
		sb.WriteString(strings.Repeat("  ", depth))
		if n.Dir {
			sb.WriteString(n.Name + "/")
			sb.WriteString("\n")
			for _, c := range n.Children {
				walk(c, depth+1)
			}
		} else {
			fmt.Fprintf(&sb, "%s (%d)\n", n.Name, n.Size)
		}
	}
	walk(n, 0)
	return sb.String()
}

// countZipTreeNodes returns the node count of a tree, root excluded.
func countZipTreeNodes(n *api.ZipTreeNode) int {
	count := 0
	for _, c := range n.Children {
		count += 1 + countZipTreeNodes(c)
	}
	return count
}

func TestBuildZipTree(t *testing.T) {
	cases := []struct {
		name       string
		entries    []zipTreeEntry
		maxDepth   int
		maxEntries int
		want       api.ZipTreeResp
	}{
		{
			name: "nesting rebuilds the folder chain",
			entries: []zipTreeEntry{
				{Name: "a.txt", Size: 3},
				{Name: "doc/api/index.html", Size: 10},
				{Name: "doc/readme.md", Size: 5},
			},
			want: zipTreeResp(zipDirNode("",
				zipDirNode("doc",
					zipDirNode("api",
						zipFileNode("index.html", 10),
					),
					zipFileNode("readme.md", 5),
				),
				zipFileNode("a.txt", 3),
			), 5),
		},
		{
			name: "both separators build the same tree",
			entries: []zipTreeEntry{
				{Name: `a\b\c.txt`, Size: 4},
				{Name: "d/e/f.txt", Size: 5},
			},
			want: zipTreeResp(zipDirNode("",
				zipDirNode("a",
					zipDirNode("b",
						zipFileNode("c.txt", 4),
					),
				),
				zipDirNode("d",
					zipDirNode("e",
						zipFileNode("f.txt", 5),
					),
				),
			), 6),
		},
		{
			name: "leading separators and empty segments are ignored",
			entries: []zipTreeEntry{
				{Name: "/Documents//doc///a.txt", Size: 1},
				{Name: `\\windows\path\x.txt`, Size: 2},
			},
			want: zipTreeResp(zipDirNode("",
				zipDirNode("Documents",
					zipDirNode("doc",
						zipFileNode("a.txt", 1),
					),
				),
				zipDirNode("windows",
					zipDirNode("path",
						zipFileNode("x.txt", 2),
					),
				),
			), 6),
		},
		{
			name: "parent segments, control characters and drive prefixes are kept as text",
			entries: []zipTreeEntry{
				{Name: "../etc/passwd", Size: 5},
				{Name: "a\x00b.txt", Size: 2},
				{Name: `C:\temp\file.txt`, Size: 3},
			},
			want: zipTreeResp(zipDirNode("",
				zipDirNode("..",
					zipDirNode("etc",
						zipFileNode("passwd", 5),
					),
				),
				zipDirNode("C:",
					zipDirNode("temp",
						zipFileNode("file.txt", 3),
					),
				),
				zipFileNode("a\x00b.txt", 2),
			), 7),
		},
		{
			name: "explicit directory entries merge into inferred folders",
			entries: []zipTreeEntry{
				{Name: "Documents/", Dir: true},
				{Name: "Documents/empty/", Dir: true},
				{Name: "Documents/a.txt", Size: 10},
			},
			want: zipTreeResp(zipDirNode("",
				zipDirNode("Documents",
					zipDirNode("empty"),
					zipFileNode("a.txt", 10),
				),
			), 3),
		},
		{
			name: "repeated file entries appear once each",
			entries: []zipTreeEntry{
				{Name: "a.txt", Size: 1},
				{Name: "a.txt", Size: 2},
			},
			want: zipTreeResp(zipDirNode("",
				zipFileNode("a.txt", 1),
				zipFileNode("a.txt", 2),
			), 2),
		},
		{
			name: "repeated directory entries merge into one folder",
			entries: []zipTreeEntry{
				{Name: "docs/", Dir: true},
				{Name: "docs/", Dir: true},
				{Name: "docs/a.txt", Size: 4},
			},
			want: zipTreeResp(zipDirNode("",
				zipDirNode("docs",
					zipFileNode("a.txt", 4),
				),
			), 2),
		},
		{
			name: "a file and a folder at one path both appear",
			entries: []zipTreeEntry{
				{Name: "x/", Dir: true},
				{Name: "x", Size: 4},
			},
			want: zipTreeResp(zipDirNode("",
				zipDirNode("x"),
				zipFileNode("x", 4),
			), 2),
		},
		{
			name: "siblings are folders first then files, case-insensitive by name",
			entries: []zipTreeEntry{
				{Name: "zebra.txt", Size: 1},
				{Name: "Apple.txt", Size: 2},
				{Name: "banana.txt", Size: 3},
				{Name: "Zdir/", Dir: true},
				{Name: "adir/", Dir: true},
			},
			want: zipTreeResp(zipDirNode("",
				zipDirNode("adir"),
				zipDirNode("Zdir"),
				zipFileNode("Apple.txt", 2),
				zipFileNode("banana.txt", 3),
				zipFileNode("zebra.txt", 1),
			), 5),
		},
		{
			name: "entries deeper than the maximum are omitted",
			entries: []zipTreeEntry{
				{Name: "a/", Dir: true},
				{Name: "a/b.txt", Size: 5},
				{Name: "a/c/d.txt", Size: 6},
				{Name: "e/f/g", Size: 7},
			},
			maxDepth: 2,
			want: zipTreeResp(zipDirNode("",
				zipDirNode("a",
					zipFileNode("b.txt", 5),
				),
			), 2, api.ZipTreeTruncatedByDepth),
		},
		{
			name: "entries beyond the maximum are omitted",
			entries: []zipTreeEntry{
				{Name: "a.txt", Size: 1},
				{Name: "b.txt", Size: 2},
				{Name: "c.txt", Size: 3},
			},
			maxEntries: 2,
			want: zipTreeResp(zipDirNode("",
				zipFileNode("a.txt", 1),
				zipFileNode("b.txt", 2),
			), 2, api.ZipTreeTruncatedByEntries),
		},
		{
			name: "an over-deep entry is left out without stopping the scan at the node cap",
			entries: []zipTreeEntry{
				{Name: "a.txt", Size: 1},
				{Name: "b.txt", Size: 2},
				{Name: "deep/deeper/deepest.txt", Size: 3},
			},
			maxDepth:   2,
			maxEntries: 2,
			want: zipTreeResp(zipDirNode("",
				zipFileNode("a.txt", 1),
				zipFileNode("b.txt", 2),
			), 2, api.ZipTreeTruncatedByDepth),
		},
		{
			name: "a repeated folder beyond the node cap creates nothing and truncates nothing",
			entries: []zipTreeEntry{
				{Name: "a/", Dir: true},
				{Name: "a/", Dir: true},
			},
			maxEntries: 1,
			want: zipTreeResp(zipDirNode("",
				zipDirNode("a"),
			), 1),
		},
		{
			name: "a new folder beyond the node cap reports entries",
			entries: []zipTreeEntry{
				{Name: "a.txt", Size: 1},
				{Name: "b/", Dir: true},
			},
			maxEntries: 1,
			want: zipTreeResp(zipDirNode("",
				zipFileNode("a.txt", 1),
			), 1, api.ZipTreeTruncatedByEntries),
		},
		{
			name: "inferred folders count toward the node cap",
			entries: []zipTreeEntry{
				{Name: "a/b/c/d/e.txt", Size: 1},
			},
			maxEntries: 3,
			// a, b and c are created, then d would exceed the cap: it is left
			// out with the rest of the path and the scan stops
			want: zipTreeResp(zipDirNode("",
				zipDirNode("a",
					zipDirNode("b",
						zipDirNode("c"),
					),
				),
			), 3, api.ZipTreeTruncatedByEntries),
		},
		{
			name:    "empty archive",
			entries: nil,
			want:    zipTreeResp(zipDirNode(""), 0),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := buildZipTree(c.entries, c.maxDepth, c.maxEntries)
			if got.Truncated != c.want.Truncated {
				t.Errorf("Truncated = %v, want %v", got.Truncated, c.want.Truncated)
			}
			if !reflect.DeepEqual(got.TruncatedReasons, c.want.TruncatedReasons) {
				t.Errorf("TruncatedReasons = %v, want %v", got.TruncatedReasons, c.want.TruncatedReasons)
			}
			if got.IncludedEntries != c.want.IncludedEntries {
				t.Errorf("IncludedEntries = %d, want %d", got.IncludedEntries, c.want.IncludedEntries)
			}
			// the cap bounds every node, and the reported count matches the tree
			limit := c.maxEntries
			if limit <= 0 {
				limit = api.DefaultZipTreeMaxEntries
			}
			if n := countZipTreeNodes(&got.Tree); n > limit {
				t.Errorf("tree has %d nodes, cap is %d", n, limit)
			}
			if n := countZipTreeNodes(&got.Tree); n != got.IncludedEntries {
				t.Errorf("tree has %d nodes, IncludedEntries = %d", n, got.IncludedEntries)
			}
			if !reflect.DeepEqual(got.Tree, c.want.Tree) {
				t.Errorf("tree =\n%s\nwant\n%s", renderZipTree(&got.Tree), renderZipTree(&c.want.Tree))
			}
		})
	}
}

func TestBuildZipTreeDefaultCaps(t *testing.T) {
	deepest := func(segments int) string {
		return strings.Repeat("d/", segments) + "file.txt"
	}

	t.Run("depth defaults to 32", func(t *testing.T) {
		// 32 segments is the deepest path the default admits; 33 is omitted.
		atLimit := buildZipTree([]zipTreeEntry{
			{Name: deepest(31), Size: 1},
		}, 0, 0)
		if atLimit.Truncated {
			t.Errorf("depth 32 entry was truncated: %v", atLimit.TruncatedReasons)
		}
		if atLimit.IncludedEntries != 32 {
			t.Errorf("IncludedEntries = %d, want 32", atLimit.IncludedEntries)
		}

		overLimit := buildZipTree([]zipTreeEntry{
			{Name: deepest(32), Size: 1},
		}, 0, 0)
		if overLimit.IncludedEntries != 0 {
			t.Errorf("IncludedEntries = %d, want 0", overLimit.IncludedEntries)
		}
		if !overLimit.Truncated || len(overLimit.TruncatedReasons) != 1 || overLimit.TruncatedReasons[0] != api.ZipTreeTruncatedByDepth {
			t.Errorf("truncation = %v %v, want depth", overLimit.Truncated, overLimit.TruncatedReasons)
		}
	})

	t.Run("entry count defaults to 10000", func(t *testing.T) {
		entries := make([]zipTreeEntry, 0, api.DefaultZipTreeMaxEntries+1)
		for i := 0; i <= api.DefaultZipTreeMaxEntries; i++ {
			entries = append(entries, zipTreeEntry{Name: fmt.Sprintf("f%05d.txt", i), Size: int64(i)})
		}

		got := buildZipTree(entries, 0, 0)
		if got.IncludedEntries != api.DefaultZipTreeMaxEntries {
			t.Errorf("IncludedEntries = %d, want %d", got.IncludedEntries, api.DefaultZipTreeMaxEntries)
		}
		if len(got.Tree.Children) != api.DefaultZipTreeMaxEntries {
			t.Errorf("root has %d children, want %d", len(got.Tree.Children), api.DefaultZipTreeMaxEntries)
		}
		if !got.Truncated || len(got.TruncatedReasons) != 1 || got.TruncatedReasons[0] != api.ZipTreeTruncatedByEntries {
			t.Errorf("truncation = %v %v, want entries", got.Truncated, got.TruncatedReasons)
		}
	})
}

func TestReadZipTree(t *testing.T) {
	t.Run("normal tree", func(t *testing.T) {
		p := writeTestZip(t, []testZipEntry{
			{Name: "docs/"},
			{Name: "docs/a.txt", Body: []byte("hello")},
			{Name: "b.txt", Body: []byte("xy")},
		})

		got, err := readZipTree(p, 0, 0)
		if err != nil {
			t.Fatalf("readZipTree returned error, %v", err)
		}

		want := zipTreeResp(zipDirNode("",
			zipDirNode("docs",
				zipFileNode("a.txt", 5),
			),
			zipFileNode("b.txt", 2),
		), 3)
		assertZipTreeResp(t, got, want)
	})

	t.Run("explicit directory entry", func(t *testing.T) {
		p := writeTestZip(t, []testZipEntry{
			{Name: "empty/"},
		})

		got, err := readZipTree(p, 0, 0)
		if err != nil {
			t.Fatalf("readZipTree returned error, %v", err)
		}

		assertZipTreeResp(t, got, zipTreeResp(zipDirNode("",
			zipDirNode("empty"),
		), 1))
	})

	t.Run("directory marked by entry mode without a trailing separator", func(t *testing.T) {
		p := writeTestZipWithDirMode(t, "folder", []testZipEntry{
			{Name: "folder/a.txt", Body: []byte("abc")},
		})

		got, err := readZipTree(p, 0, 0)
		if err != nil {
			t.Fatalf("readZipTree returned error, %v", err)
		}

		assertZipTreeResp(t, got, zipTreeResp(zipDirNode("",
			zipDirNode("folder",
				zipFileNode("a.txt", 3),
			),
		), 2))
	})

	t.Run("backslash directory and path", func(t *testing.T) {
		p := writeTestZip(t, []testZipEntry{
			{Name: `win\dir\`},
			{Name: `win\dir\file.txt`, Body: []byte("abc")},
		})

		got, err := readZipTree(p, 0, 0)
		if err != nil {
			t.Fatalf("readZipTree returned error, %v", err)
		}

		assertZipTreeResp(t, got, zipTreeResp(zipDirNode("",
			zipDirNode("win",
				zipDirNode("dir",
					zipFileNode("file.txt", 3),
				),
			),
		), 3))
	})

	t.Run("non-zip file", func(t *testing.T) {
		p := writeTestFile(t, "not a zip archive, just bytes")

		_, err := readZipTree(p, 0, 0)
		assertMisoErrCode(t, err, api.IllegalFormat)
	})

	t.Run("truncated zip file", func(t *testing.T) {
		p := writeTestZip(t, []testZipEntry{
			{Name: "a.txt", Body: []byte("hello")},
		})
		b, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("failed to read test zip, %v", err)
		}
		if err := os.WriteFile(p, b[:len(b)-32], 0o644); err != nil {
			t.Fatalf("failed to truncate test zip, %v", err)
		}

		_, err = readZipTree(p, 0, 0)
		assertMisoErrCode(t, err, api.IllegalFormat)
	})

	t.Run("entry cap small enough to truncate", func(t *testing.T) {
		p := writeTestZip(t, []testZipEntry{
			{Name: "a.txt", Body: []byte("a")},
			{Name: "b.txt", Body: []byte("b")},
			{Name: "c.txt", Body: []byte("c")},
		})

		got, err := readZipTree(p, 0, 2)
		if err != nil {
			t.Fatalf("readZipTree returned error, %v", err)
		}

		assertZipTreeResp(t, got, zipTreeResp(zipDirNode("",
			zipFileNode("a.txt", 1),
			zipFileNode("b.txt", 1),
		), 2, api.ZipTreeTruncatedByEntries))
	})
}

func assertZipTreeResp(t *testing.T, got, want api.ZipTreeResp) {
	t.Helper()
	if got.Truncated != want.Truncated {
		t.Errorf("Truncated = %v, want %v", got.Truncated, want.Truncated)
	}
	if !reflect.DeepEqual(got.TruncatedReasons, want.TruncatedReasons) {
		t.Errorf("TruncatedReasons = %v, want %v", got.TruncatedReasons, want.TruncatedReasons)
	}
	if got.IncludedEntries != want.IncludedEntries {
		t.Errorf("IncludedEntries = %d, want %d", got.IncludedEntries, want.IncludedEntries)
	}
	if !reflect.DeepEqual(got.Tree, want.Tree) {
		t.Errorf("tree =\n%s\nwant\n%s", renderZipTree(&got.Tree), renderZipTree(&want.Tree))
	}
}

func writeTestFile(t *testing.T, content string) string {
	t.Helper()
	p := t.TempDir() + "/file.bin"
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test file, %v", err)
	}
	return p
}

// writeTestZipWithDirMode writes an archive whose dirName node is a directory
// only through its entry mode: the name carries no trailing separator.
func writeTestZipWithDirMode(t *testing.T, dirName string, entries []testZipEntry) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "archive.zip")
	f, err := os.Create(p)
	if err != nil {
		t.Fatalf("failed to create test zip, %v", err)
	}
	zw := zip.NewWriter(f)
	h := &zip.FileHeader{Name: dirName}
	h.SetMode(os.ModeDir | 0o755)
	if _, err := zw.CreateHeader(h); err != nil {
		t.Fatalf("failed to create mode-directory entry %q, %v", dirName, err)
	}
	for _, e := range entries {
		w, err := zw.Create(e.Name)
		if err != nil {
			t.Fatalf("failed to create zip entry %q, %v", e.Name, err)
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
	return p
}

func assertMisoErrCode(t *testing.T, err error, code string) {
	t.Helper()
	m, ok := err.(*errs.MisoErr)
	if !ok || m.Code() != code {
		t.Fatalf("error = %v, want code %q", err, code)
	}
}
