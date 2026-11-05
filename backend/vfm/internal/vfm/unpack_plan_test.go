package vfm

import (
	"reflect"
	"strings"
	"testing"
)

// withExistingDir is a target directory child that is a folder.
func withExistingDir(name, key string) unpackTargetChild {
	return unpackTargetChild{Name: name, FileKey: key, IsDir: true}
}

// withExistingFile is a target directory child that is a file.
func withExistingFile(name, key string) unpackTargetChild {
	return unpackTargetChild{Name: name, FileKey: key}
}

func normalizePlan(p unpackPlan) unpackPlan {
	if p.Dirs == nil {
		p.Dirs = []unpackPlanDir{}
	}
	if p.Files == nil {
		p.Files = []unpackPlanFile{}
	}
	if p.Drops == nil {
		p.Drops = []unpackDrop{}
	}
	return p
}

func TestPlanUnpackTree(t *testing.T) {
	tests := []struct {
		name     string
		entries  []string
		existing []unpackTargetChild
		want     unpackPlan
	}{
		{
			name:    "nested entry rebuilds its folder layers",
			entries: []string{"Documents/doc/jdk/docs/api/a.html"},
			want: unpackPlan{
				Dirs: []unpackPlanDir{
					{Index: 0, Path: "Documents", Name: "Documents", Parent: ""},
					{Index: 0, Path: "Documents/doc", Name: "doc", Parent: "Documents"},
					{Index: 0, Path: "Documents/doc/jdk", Name: "jdk", Parent: "Documents/doc"},
					{Index: 0, Path: "Documents/doc/jdk/docs", Name: "docs", Parent: "Documents/doc/jdk"},
					{Index: 0, Path: "Documents/doc/jdk/docs/api", Name: "api", Parent: "Documents/doc/jdk/docs"},
				},
				Files: []unpackPlanFile{
					{Index: 0, Path: "Documents/doc/jdk/docs/api/a.html", Name: "a.html", Parent: "Documents/doc/jdk/docs/api"},
				},
			},
		},
		{
			name:    "entries sharing ancestors create each folder once",
			entries: []string{"a/b/x.txt", "a/c/y.txt", "a/b/z.txt"},
			want: unpackPlan{
				Dirs: []unpackPlanDir{
					{Index: 0, Path: "a", Name: "a", Parent: ""},
					{Index: 0, Path: "a/b", Name: "b", Parent: "a"},
					{Index: 1, Path: "a/c", Name: "c", Parent: "a"},
				},
				Files: []unpackPlanFile{
					{Index: 0, Path: "a/b/x.txt", Name: "x.txt", Parent: "a/b"},
					{Index: 1, Path: "a/c/y.txt", Name: "y.txt", Parent: "a/c"},
					{Index: 2, Path: "a/b/z.txt", Name: "z.txt", Parent: "a/b"},
				},
			},
		},
		{
			name:     "an existing folder is reused",
			entries:  []string{"Documents/a.txt", "Documents/sub/b.txt"},
			existing: []unpackTargetChild{withExistingDir("Documents", "KEY_DOCS")},
			want: unpackPlan{
				Dirs: []unpackPlanDir{
					{Index: 1, Path: "Documents/sub", Name: "sub", Parent: "Documents"},
				},
				Files: []unpackPlanFile{
					{Index: 0, Path: "Documents/a.txt", Name: "a.txt", Parent: "Documents"},
					{Index: 1, Path: "Documents/sub/b.txt", Name: "b.txt", Parent: "Documents/sub"},
				},
			},
		},
		{
			name:     "a folder name held by an existing file is numbered",
			entries:  []string{"Documents/a.txt"},
			existing: []unpackTargetChild{withExistingFile("Documents", "KEY_FILE")},
			want: unpackPlan{
				Dirs: []unpackPlanDir{
					{Index: 0, Path: "Documents (1)", Name: "Documents (1)", Parent: ""},
				},
				Files: []unpackPlanFile{
					{Index: 0, Path: "Documents (1)/a.txt", Name: "a.txt", Parent: "Documents (1)"},
				},
			},
		},
		{
			name:     "two entries blocked on the same folder share its numbered name",
			entries:  []string{"x/a.txt", "x/b.txt"},
			existing: []unpackTargetChild{withExistingFile("x", "KEY_FILE")},
			want: unpackPlan{
				Dirs: []unpackPlanDir{
					{Index: 0, Path: "x (1)", Name: "x (1)", Parent: ""},
				},
				Files: []unpackPlanFile{
					{Index: 0, Path: "x (1)/a.txt", Name: "a.txt", Parent: "x (1)"},
					{Index: 1, Path: "x (1)/b.txt", Name: "b.txt", Parent: "x (1)"},
				},
			},
		},
		{
			name:    "two entries resolving to the same file path are numbered",
			entries: []string{"a/x.txt", "a/x.txt", "a/x.txt"},
			want: unpackPlan{
				Dirs: []unpackPlanDir{{Index: 0, Path: "a", Name: "a", Parent: ""}},
				Files: []unpackPlanFile{
					{Index: 0, Path: "a/x.txt", Name: "x.txt", Parent: "a"},
					{Index: 1, Path: "a/x (1).txt", Name: "x (1).txt", Parent: "a"},
					{Index: 2, Path: "a/x (2).txt", Name: "x (2).txt", Parent: "a"},
				},
			},
		},
		{
			name:    "a file first, then a folder of the same name",
			entries: []string{"a", "a/b.txt"},
			want: unpackPlan{
				Dirs: []unpackPlanDir{{Index: 1, Path: "a (1)", Name: "a (1)", Parent: ""}},
				Files: []unpackPlanFile{
					{Index: 0, Path: "a", Name: "a", Parent: ""},
					{Index: 1, Path: "a (1)/b.txt", Name: "b.txt", Parent: "a (1)"},
				},
			},
		},
		{
			name:    "leading separators and dot segments disappear",
			entries: []string{"/a/./b.txt"},
			want: unpackPlan{
				Dirs:  []unpackPlanDir{{Index: 0, Path: "a", Name: "a", Parent: ""}},
				Files: []unpackPlanFile{{Index: 0, Path: "a/b.txt", Name: "b.txt", Parent: "a"}},
			},
		},
		{
			name:    "backslashes are read as separators",
			entries: []string{`dir\sub\file.txt`},
			want: unpackPlan{
				Dirs: []unpackPlanDir{
					{Index: 0, Path: "dir", Name: "dir", Parent: ""},
					{Index: 0, Path: "dir/sub", Name: "sub", Parent: "dir"},
				},
				Files: []unpackPlanFile{{Index: 0, Path: "dir/sub/file.txt", Name: "file.txt", Parent: "dir/sub"}},
			},
		},
		{
			name:    "directory entries are skipped",
			entries: []string{"docs/", `docs\`, "docs/a.txt"},
			want: unpackPlan{
				Dirs: []unpackPlanDir{{Index: 2, Path: "docs", Name: "docs", Parent: ""}},
				Files: []unpackPlanFile{
					{Index: 2, Path: "docs/a.txt", Name: "a.txt", Parent: "docs"},
				},
				Drops: []unpackDrop{
					{Path: "docs/", Reason: "the entry names a directory"},
					{Path: `docs\`, Reason: "the entry names a directory"},
				},
			},
		},
		{
			name:    "an entry escaping the target directory is dropped",
			entries: []string{"../evil.txt", "a/../../evil.txt", "a/../b.txt", "good.txt"},
			want: unpackPlan{
				Files: []unpackPlanFile{{Index: 3, Path: "good.txt", Name: "good.txt", Parent: ""}},
				Drops: []unpackDrop{
					{Path: "../evil.txt", Reason: "the path contains a '..' segment"},
					{Path: "a/../../evil.txt", Reason: "the path contains a '..' segment"},
					{Path: "a/../b.txt", Reason: "the path contains a '..' segment"},
				},
			},
		},
		{
			name:    "a control character drops the entry",
			entries: []string{"a/b\x01c.txt"},
			want: unpackPlan{
				Drops: []unpackDrop{{Path: "a/b\x01c.txt", Reason: "a name contains a control character"}},
			},
		},
		{
			name:    "a name ending with a dot or a space drops the entry",
			entries: []string{"a/b.", "a/c "},
			want: unpackPlan{
				Drops: []unpackDrop{
					{Path: "a/b.", Reason: "a name ends with a dot or a space"},
					{Path: "a/c ", Reason: "a name ends with a dot or a space"},
				},
			},
		},
		{
			name:    "a segment longer than the name column drops the entry",
			entries: []string{"a/" + strings.Repeat("n", 252) + ".txt", "b/" + strings.Repeat("n", 251) + ".txt"},
			want: unpackPlan{
				Dirs: []unpackPlanDir{{Index: 1, Path: "b", Name: "b", Parent: ""}},
				Files: []unpackPlanFile{
					{Index: 1, Path: "b/" + strings.Repeat("n", 251) + ".txt", Name: strings.Repeat("n", 251) + ".txt", Parent: "b"},
				},
				Drops: []unpackDrop{
					{Path: "a/" + strings.Repeat("n", 252) + ".txt", Reason: "a name is longer than 255 bytes"},
				},
			},
		},
		{
			name:    "an entry with no usable segment is dropped",
			entries: []string{"/", "./"},
			want: unpackPlan{
				Drops: []unpackDrop{
					{Path: "/", Reason: "the entry names a directory"},
					{Path: "./", Reason: "the entry names a directory"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizePlan(planUnpackTree(tt.entries, tt.existing))
			want := normalizePlan(tt.want)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("plan:\n got: %+v\nwant: %+v", got, want)
			}
		})
	}
}

func TestPlanUnpackTreeKeepsFoldersOfDifferentCaseApart(t *testing.T) {
	// names are taken verbatim, so case stays as the archive wrote it
	plan := planUnpackTree([]string{"Docs/a.txt", "docs/b.txt"}, nil)
	if len(plan.Dirs) != 2 {
		t.Fatalf("dirs: got %+v, want two folders", plan.Dirs)
	}
	if plan.Dirs[0].Path != "Docs" || plan.Dirs[1].Path != "docs" {
		t.Fatalf("dirs: got %+v", plan.Dirs)
	}
}

func TestPlanUnpackTreeDropsANameThatNumberingPushesPastTheLimit(t *testing.T) {
	// 255 bytes as it arrives, so it only breaks once " (1)" is inserted
	blunt := strings.Repeat("f", 251) + ".txt"
	if len(blunt) != unpackMaxNameBytes {
		t.Fatalf("fixture is not at the limit: %d", len(blunt))
	}

	plan := planUnpackTree([]string{blunt, blunt}, []unpackTargetChild{withExistingFile(blunt, "KEY_FILE")})
	if len(plan.Files) != 0 || len(plan.Dirs) != 0 {
		t.Fatalf("plan: got %+v, want nothing created", plan)
	}
	if len(plan.Drops) != 2 {
		t.Fatalf("drops: got %+v, want both entries dropped", plan.Drops)
	}
	for _, d := range plan.Drops {
		if d.Path != blunt || d.Reason != unpackNameTooLong {
			t.Fatalf("drop: got %+v", d)
		}
	}
}

func TestPlanUnpackTreeDropsASubtreeWhenItsFolderNameGrowsTooLong(t *testing.T) {
	// a 255 byte folder name that is taken, so numbering makes it unstoreable
	folder := strings.Repeat("d", 255)
	plan := planUnpackTree(
		[]string{folder + "/a.txt", folder + "/b.txt"},
		[]unpackTargetChild{withExistingFile(folder, "KEY_FILE")},
	)

	if len(plan.Files) != 0 || len(plan.Dirs) != 0 {
		t.Fatalf("plan: got %+v, want nothing created", plan)
	}
	if len(plan.Drops) != 2 {
		t.Fatalf("drops: got %+v, want both entries dropped", plan.Drops)
	}
	for _, d := range plan.Drops {
		if d.Reason != unpackNameTooLong {
			t.Fatalf("drop: got %+v", d)
		}
	}
}

func TestPlanUnpackTreeKeepsANumberedNameAtTheLimit(t *testing.T) {
	// 251 bytes plus " (1)" is exactly the column width
	name := strings.Repeat("g", 251)
	plan := planUnpackTree([]string{name}, []unpackTargetChild{withExistingFile(name, "KEY_FILE")})

	if len(plan.Drops) != 0 {
		t.Fatalf("drops: got %+v, want none", plan.Drops)
	}
	if len(plan.Files) != 1 {
		t.Fatalf("files: got %+v, want one", plan.Files)
	}
	if got, want := plan.Files[0].Name, name+" (1)"; got != want {
		t.Fatalf("name: got %q, want %q", got, want)
	}
	if len(plan.Files[0].Name) != unpackMaxNameBytes {
		t.Fatalf("fixture is not at the limit: %d", len(plan.Files[0].Name))
	}
}
