package vfm

import (
	"strings"
	"unicode"

	"github.com/curtisnewbie/miso/errs"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/miso"
	"gorm.io/gorm"
)

const (
	// unpackMaxNameBytes is the width of the name column, so one path segment
	// longer than this can never be stored.
	unpackMaxNameBytes = 255

	// unpackNameTooLong is the reason reported when a name, raw or numbered,
	// does not fit the name column.
	unpackNameTooLong = "a name is longer than 255 bytes"
)

// unpackTargetChild is one existing direct child of the unpack target
// directory. Only that one level is read, because a folder already present
// there is merged instead of duplicated.
type unpackTargetChild struct {
	Name    string
	FileKey string // only set for a directory, it is reused as the parent of the entry
	IsDir   bool
}

// unpackDrop is one entry left out of the unpack, with the reason reported to
// the operator.
type unpackDrop struct {
	Path   string // entry name as it arrived in the reply
	Reason string
}

// unpackPlanDir is a folder the writer still has to create under the target
// directory.
type unpackPlanDir struct {
	Index  int    // index of the entry that first needs the folder
	Path   string // slash separated path relative to the target directory, e.g. "Documents/doc"
	Name   string // resolved name of the deepest segment
	Parent string // path of the parent folder, "" for the target directory itself
}

// unpackPlanFile is a file entry resolved to the folder that holds it and its
// own name, never its whole path.
type unpackPlanFile struct {
	Index  int    // position of the entry in the reply's entry list
	Path   string // slash separated path relative to the target directory, for logging
	Name   string // resolved file name, numbered when another file already holds it
	Parent string // path of the parent folder, "" for the target directory itself
}

// unpackPlan is everything an unpack run needs, resolved before any row is
// written. Dirs are ordered parents first, and each carries the entry that
// first needs it, so the writer creates a folder only when it reaches a file
// that lives in it.
type unpackPlan struct {
	Dirs  []unpackPlanDir
	Files []unpackPlanFile
	Drops []unpackDrop
}

// unpackPlanNode is one folder of the planned tree. It is keyed by the entry
// path it was first needed for, so a second entry sharing that path reuses the
// same folder even when the folder had to be numbered.
type unpackPlanNode struct {
	rawPath string
	path    string // resolved path
	parent  string // resolved parent path
}

// unpackPlanTree is the mutable state behind planUnpackTree.
type unpackPlanTree struct {
	plan         unpackPlan
	nodes        map[string]*unpackPlanNode
	failed       map[string]string              // raw folder path -> reason it cannot be created
	existingDirs map[string]struct{}            // resolved path -> existing directory of the target dir
	taken        map[string]map[string]struct{} // resolved parent path -> names already used below it
}

// planUnpackTree resolves the reply's entry names into the folders and files to
// create under the target directory.
//
// Entries are processed in reply order, so numbering follows that order.
// Nothing is written outside the target directory: every resolved path is
// relative to it and holds no "..". An entry that cannot be expressed is
// dropped with a reason instead of being rewritten.
func planUnpackTree(entryNames []string, existing []unpackTargetChild) unpackPlan {
	t := &unpackPlanTree{
		plan: unpackPlan{
			Dirs:  make([]unpackPlanDir, 0, len(entryNames)/2),
			Files: make([]unpackPlanFile, 0, len(entryNames)),
		},
		nodes:        make(map[string]*unpackPlanNode, len(entryNames)),
		failed:       make(map[string]string, 0),
		existingDirs: make(map[string]struct{}, len(existing)),
		taken:        make(map[string]map[string]struct{}, len(existing)),
	}

	for _, c := range existing {
		t.take("", c.Name)
		if c.IsDir {
			t.existingDirs[c.Name] = struct{}{}
		}
	}

	for i, raw := range entryNames {
		segments, drop := unpackEntrySegments(raw)
		if drop != nil {
			t.plan.Drops = append(t.plan.Drops, *drop)
			continue
		}

		rawParent, parent := "", ""
		dropReason := ""
		for _, seg := range segments[:len(segments)-1] {
			n, reason := t.resolveDir(rawParent, seg, parent, i)
			if reason != "" {
				dropReason = reason
				break
			}
			rawParent, parent = n.rawPath, n.path
		}
		if dropReason == "" {
			name := dirZipEntryName(segments[len(segments)-1], t.nameSet(parent))
			if len(name) > unpackMaxNameBytes {
				// numbering can push a name that just fit past the column
				dropReason = unpackNameTooLong
			}
			if dropReason == "" {
				t.plan.Files = append(t.plan.Files, unpackPlanFile{
					Index:  i,
					Path:   unpackJoinPath(parent, name),
					Name:   name,
					Parent: parent,
				})
				continue
			}
		}
		t.plan.Drops = append(t.plan.Drops, unpackDrop{Path: raw, Reason: dropReason})
	}

	return t.plan
}

// resolveDir returns the folder for one segment below parent, creating it in
// the plan only when it does not exist yet. A folder whose numbered name no
// longer fits the name column is refused, and every entry under it is dropped
// with the same reason.
func (t *unpackPlanTree) resolveDir(rawParent, rawName, parent string, index int) (*unpackPlanNode, string) {
	rawPath := unpackJoinPath(rawParent, rawName)
	if n, ok := t.nodes[rawPath]; ok {
		return n, ""
	}
	if reason, ok := t.failed[rawPath]; ok {
		return nil, reason
	}

	// a folder already present in the target directory is merged into
	if _, ok := t.existingDirs[unpackJoinPath(parent, rawName)]; ok {
		n := &unpackPlanNode{rawPath: rawPath, path: unpackJoinPath(parent, rawName), parent: parent}
		t.nodes[rawPath] = n
		return n, ""
	}

	name := dirZipEntryName(rawName, t.nameSet(parent))
	if len(name) > unpackMaxNameBytes {
		t.failed[rawPath] = unpackNameTooLong
		return nil, unpackNameTooLong
	}

	n := &unpackPlanNode{rawPath: rawPath, path: unpackJoinPath(parent, name), parent: parent}
	t.nodes[rawPath] = n
	t.plan.Dirs = append(t.plan.Dirs, unpackPlanDir{Index: index, Path: n.path, Name: name, Parent: parent})
	return n, ""
}

func (t *unpackPlanTree) nameSet(parent string) map[string]struct{} {
	s, ok := t.taken[parent]
	if !ok {
		s = map[string]struct{}{}
		t.taken[parent] = s
	}
	return s
}

func (t *unpackPlanTree) take(parent, name string) {
	t.nameSet(parent)[name] = struct{}{}
}

// unpackEntrySegments sanitises one entry name into its path segments, or
// reports why the entry cannot be honoured.
//
// Leading separators are dropped, backslashes are read as separators, and empty
// and "." segments disappear. A ".." segment, a control character, a trailing
// dot or space, or a segment wider than the name column rejects the whole
// entry: the name is never rewritten.
func unpackEntrySegments(name string) ([]string, *unpackDrop) {
	if strings.HasSuffix(name, "/") || strings.HasSuffix(name, `\`) {
		return nil, &unpackDrop{Path: name, Reason: "the entry names a directory"}
	}

	segments := make([]string, 0, 4)
	for _, seg := range strings.Split(strings.ReplaceAll(name, `\`, "/"), "/") {
		if seg == "" || seg == "." {
			continue
		}
		if reason := unpackSegmentProblem(seg); reason != "" {
			return nil, &unpackDrop{Path: name, Reason: reason}
		}
		segments = append(segments, seg)
	}

	if len(segments) == 0 {
		return nil, &unpackDrop{Path: name, Reason: "the entry name has no path segments"}
	}
	return segments, nil
}

// unpackSegmentProblem reports why one path segment cannot be stored, or "".
func unpackSegmentProblem(seg string) string {
	if seg == ".." {
		return "the path contains a '..' segment"
	}
	for _, r := range seg {
		if unicode.IsControl(r) {
			return "a name contains a control character"
		}
	}
	if strings.HasSuffix(seg, ".") || strings.HasSuffix(seg, " ") {
		return "a name ends with a dot or a space"
	}
	if len(seg) > unpackMaxNameBytes {
		return unpackNameTooLong
	}
	return ""
}

func unpackJoinPath(parent, name string) string {
	if parent == "" {
		return name
	}
	return parent + "/" + name
}

// listUnpackTargetChildren reads the target directory's direct children once,
// so a folder that is already there is merged rather than duplicated and a name
// that is already held is numbered.
func listUnpackTargetChildren(rail miso.Rail, db *gorm.DB, dirKey string) ([]unpackTargetChild, error) {
	if dirKey == "" {
		return nil, nil
	}

	type childRow struct {
		Name     string
		FileKey  string
		FileType string
	}
	var rows []childRow
	if err := dbquery.NewQuery(rail, db).
		Table("file_info").
		Select("name, uuid as file_key, file_type").
		Eq("parent_file", dirKey).
		Eq("is_logic_deleted", DelN).
		Eq("is_del", 0).
		ScanVal(&rows); err != nil {
		return nil, errs.Wrapf(err, "failed to list children of directory, fileKey: %v", dirKey)
	}

	children := make([]unpackTargetChild, 0, len(rows))
	for _, r := range rows {
		children = append(children, unpackTargetChild{
			Name:    r.Name,
			FileKey: r.FileKey,
			IsDir:   r.FileType == FileTypeDir,
		})
	}
	return children, nil
}
