package fstore

import (
	"archive/zip"
	"sort"
	"strings"

	"github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/miso/errs"
	"github.com/curtisnewbie/miso/miso"
	"gorm.io/gorm"
)

// zipTreeEntry is one archive entry as read from the zip central directory.
type zipTreeEntry struct {
	Name string
	Dir  bool
	Size int64
}

// ZipTree returns the folder and file tree of a stored zip file.
//
// The archive's central directory is read only: no entry is decompressed, no
// temporary file is created and nothing is written. A missing or logically
// deleted file record is refused with FILE_NOT_FOUND, a file that cannot be
// read as a zip with ILLEGAL_FORMAT.
func ZipTree(rail miso.Rail, db *gorm.DB, fileId string, maxDepth, maxEntries int) (api.ZipTreeResp, error) {
	f, err := FindFile(rail, db, fileId)
	if err != nil {
		return api.ZipTreeResp{}, errs.Wrapf(err, "failed to find file to preview, fileId: %v", fileId)
	}
	if f.IsZero() || f.IsDeleted() {
		return api.ZipTreeResp{}, ErrFileNotFound.New()
	}
	return readZipTree(f.StoragePath(), maxDepth, maxEntries)
}

// readZipTree reads path's zip central directory and builds its tree. The caps
// fall back to the server defaults when non-positive.
func readZipTree(path string, maxDepth, maxEntries int) (api.ZipTreeResp, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return api.ZipTreeResp{}, errs.NewErrfCode(api.IllegalFormat, "Zip file cannot be read").
			WithInternalMsg("failed to open zip file, path: %v, %v", path, err)
	}
	defer r.Close()

	b := newZipTreeBuilder(maxDepth, maxEntries)
	for _, f := range r.File {
		// the entry is a directory when its file info says so, or when its name
		// ends with a separator, so a zip packed on Windows counts the same as
		// one packed elsewhere
		entry := zipTreeEntry{
			Name: f.Name,
			Dir:  f.FileInfo().IsDir() || strings.HasSuffix(f.Name, "\\"),
			Size: int64(f.UncompressedSize64),
		}
		// the cap is reached: stop the scan, the rest are left out
		if b.add(entry) {
			break
		}
	}
	return b.result(), nil
}

type zipTreeBuilder struct {
	root *api.ZipTreeNode
	// folders indexes every folder node by its normalized path, so inferred
	// parents and repeated directory entries merge into one node
	folders map[string]*api.ZipTreeNode

	maxDepth   int
	maxEntries int
	// included counts the nodes created, root excluded: file nodes, explicit
	// folder nodes and inferred folders all count once
	included int

	truncatedByDepth   bool
	truncatedByEntries bool
}

func newZipTreeBuilder(maxDepth, maxEntries int) *zipTreeBuilder {
	if maxDepth <= 0 {
		maxDepth = api.DefaultZipTreeMaxDepth
	}
	if maxEntries <= 0 {
		maxEntries = api.DefaultZipTreeMaxEntries
	}
	return &zipTreeBuilder{
		root:       &api.ZipTreeNode{Dir: true, Children: []*api.ZipTreeNode{}},
		folders:    map[string]*api.ZipTreeNode{},
		maxDepth:   maxDepth,
		maxEntries: maxEntries,
	}
}

// add represents one archive entry in the tree. It returns true when the node
// cap left nodes out, which means the scan should stop.
func (b *zipTreeBuilder) add(e zipTreeEntry) bool {
	segments := zipPathSegments(e.Name)
	// a name made only of separators carries no path segment
	if len(segments) == 0 {
		return false
	}
	// the depth check comes before any node is created, so an over-deep entry
	// leaves the tree untouched and does not stop the scan
	if len(segments) > b.maxDepth {
		b.truncatedByDepth = true
		return false
	}

	if e.Dir {
		_, ok := b.ensureFolders(segments)
		return !ok
	}

	last := len(segments) - 1
	parent, ok := b.ensureFolders(segments[:last])
	if !ok {
		return true
	}
	// the file node counts like every other node
	if b.included >= b.maxEntries {
		b.truncatedByEntries = true
		return true
	}
	parent.Children = append(parent.Children, &api.ZipTreeNode{
		Name: segments[last],
		Size: e.Size,
	})
	b.included++
	return false
}

// ensureFolders creates the folder chain for segments, merging into folders
// that already exist, and returns the deepest folder. An empty chain returns
// the root. It reports false when the node cap stopped it: the folder that
// would exceed the cap is not created and the scan should stop.
func (b *zipTreeBuilder) ensureFolders(segments []string) (*api.ZipTreeNode, bool) {
	node := b.root
	key := ""
	for _, seg := range segments {
		if key == "" {
			key = seg
		} else {
			key += "/" + seg
		}
		child, ok := b.folders[key]
		if !ok {
			if b.included >= b.maxEntries {
				b.truncatedByEntries = true
				return nil, false
			}
			child = &api.ZipTreeNode{Name: seg, Dir: true, Children: []*api.ZipTreeNode{}}
			b.folders[key] = child
			node.Children = append(node.Children, child)
			b.included++
		}
		node = child
	}
	return node, true
}

func (b *zipTreeBuilder) result() api.ZipTreeResp {
	sortZipTreeChildren(b.root)
	reasons := []string{}
	if b.truncatedByDepth {
		reasons = append(reasons, api.ZipTreeTruncatedByDepth)
	}
	if b.truncatedByEntries {
		reasons = append(reasons, api.ZipTreeTruncatedByEntries)
	}
	return api.ZipTreeResp{
		Tree:             *b.root,
		Truncated:        len(reasons) > 0,
		TruncatedReasons: reasons,
		IncludedEntries:  b.included,
	}
}

// zipPathSegments splits an entry name on both separators, dropping leading,
// trailing and repeated separators. Every other character is kept as text.
func zipPathSegments(name string) []string {
	return strings.FieldsFunc(name, func(r rune) bool {
		return r == '/' || r == '\\'
	})
}

// sortZipTreeChildren orders every folder's children: folders first, then
// files, each group case-insensitive by name. Equal names keep archive order.
func sortZipTreeChildren(n *api.ZipTreeNode) {
	if !n.Dir {
		return
	}
	sort.SliceStable(n.Children, func(i, j int) bool {
		a, b := n.Children[i], n.Children[j]
		if a.Dir != b.Dir {
			return a.Dir
		}
		la, lb := strings.ToLower(a.Name), strings.ToLower(b.Name)
		return la < lb
	})
	for _, c := range n.Children {
		sortZipTreeChildren(c)
	}
}
