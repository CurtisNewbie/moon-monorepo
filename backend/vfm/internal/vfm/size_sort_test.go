package vfm

import (
	"reflect"
	"testing"

	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/randutil"
	"gorm.io/gorm"
)

// --- fixtures ---------------------------------------------------------------

// insertSortTestChildSized inserts a file row of the given size under parent and
// returns its uuid.
func insertSortTestChildSized(t *testing.T, db *gorm.DB, parent, name, userNo string, size int64) string {
	t.Helper()
	uuid := "sst-child-" + randutil.ERand(8)
	if err := db.Exec(`INSERT INTO file_info (name, uuid, size_in_bytes, file_type, parent_file, thumbnail, uploader_no, seq_key)
		VALUES (?,?,?,?,?,?,?,?)`,
		name, uuid, size, FileTypeFile, parent, "", userNo, "").Error; err != nil {
		t.Fatal(err)
	}
	return uuid
}

// insertSortTestSubDir inserts a directory row of the given size under parent.
func insertSortTestSubDir(t *testing.T, db *gorm.DB, parent, name, userNo string, size int64) string {
	t.Helper()
	uuid := "sst-subdir-" + randutil.ERand(8)
	if err := db.Exec(`INSERT INTO file_info (name, uuid, size_in_bytes, file_type, parent_file, thumbnail, uploader_no)
		VALUES (?,?,?,?,?,?,?)`,
		name, uuid, size, FileTypeDir, parent, "", userNo).Error; err != nil {
		t.Fatal(err)
	}
	return uuid
}

// listSortedFiles lists dir with the given request fields (the parent and page
// are filled in) and returns the payload.
func listSortedFiles(t *testing.T, db *gorm.DB, dir string, req ListFileReq, user flow.User) []ListedFile {
	t.Helper()
	parent := dir
	req.ParentFile = &parent
	req.Page = miso.Paging{Page: 1, Limit: 100}
	res, err := ListFiles(miso.EmptyRail(), db, req, user)
	if err != nil {
		t.Fatal(err)
	}
	return res.Payload
}

// listSortedFileNames lists dir with the given request fields and returns the names.
func listSortedFileNames(t *testing.T, db *gorm.DB, dir string, req ListFileReq, user flow.User) []string {
	t.Helper()
	payload := listSortedFiles(t, db, dir, req, user)
	names := make([]string, 0, len(payload))
	for _, f := range payload {
		names = append(names, f.Name)
	}
	return names
}

// --- ordering ---------------------------------------------------------------

func TestListFilesSortBySizeDefaultsToDescending(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, false)

	// inserted in neither size nor id order
	insertSortTestChildSized(t, db, dir, "small.bin", user.UserNo, 100)
	insertSortTestChildSized(t, db, dir, "large.bin", user.UserNo, 900)
	insertSortTestChildSized(t, db, dir, "medium.bin", user.UserNo, 500)

	got := listSortedFileNames(t, db, dir, ListFileReq{OrderBy: "size"}, user)
	want := []string{"large.bin", "medium.bin", "small.bin"} // omitted direction → desc
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("size desc: got %v, want %v", got, want)
	}
}

func TestListFilesSortBySizeHonorsAscending(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, false)

	insertSortTestChildSized(t, db, dir, "small.bin", user.UserNo, 100)
	insertSortTestChildSized(t, db, dir, "large.bin", user.UserNo, 900)
	insertSortTestChildSized(t, db, dir, "medium.bin", user.UserNo, 500)

	got := listSortedFileNames(t, db, dir, ListFileReq{OrderBy: "size", OrderDir: "asc"}, user)
	want := []string{"small.bin", "medium.bin", "large.bin"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("size asc: got %v, want %v", got, want)
	}
}

func TestListFilesSortBySizeBreaksTiesById(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, false)

	// ids: c=1, a=2, b=3; all the same size
	insertSortTestChildSized(t, db, dir, "c.bin", user.UserNo, 500)
	insertSortTestChildSized(t, db, dir, "a.bin", user.UserNo, 500)
	insertSortTestChildSized(t, db, dir, "b.bin", user.UserNo, 500)

	cases := []struct {
		orderDir string
		want     []string
	}{
		{"", []string{"b.bin", "a.bin", "c.bin"}},     // omitted → desc, id desc
		{"desc", []string{"b.bin", "a.bin", "c.bin"}}, // id desc
		{"asc", []string{"c.bin", "a.bin", "b.bin"}},  // id asc
	}
	for _, c := range cases {
		got := listSortedFileNames(t, db, dir, ListFileReq{OrderBy: "size", OrderDir: c.orderDir}, user)
		if !reflect.DeepEqual(got, c.want) {
			t.Fatalf("equal sizes orderDir=%q: got %v, want %v", c.orderDir, got, c.want)
		}
	}
}

// --- type filter coercion ---------------------------------------------------

func TestListFilesSortBySizeReturnsFilesOnly(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, false)

	insertSortTestChildSized(t, db, dir, "small.bin", user.UserNo, 100)
	// a directory row large enough in bytes to top a mixed size listing
	insertSortTestSubDir(t, db, dir, "zzz-subdir", user.UserNo, 9999)
	insertSortTestChildSized(t, db, dir, "large.bin", user.UserNo, 900)

	want := []string{"large.bin", "small.bin"} // largest first, never a directory
	all := FileTypeFile
	empty := ""
	dirFilter := FileTypeDir
	cases := []struct {
		name     string
		fileType *string
	}{
		{"absent", nil},
		{"all", &empty},
		{"directory", &dirFilter},
		{"file", &all},
	}
	for _, c := range cases {
		payload := listSortedFiles(t, db, dir, ListFileReq{OrderBy: "size", FileType: c.fileType}, user)
		names := make([]string, 0, len(payload))
		for _, f := range payload {
			if f.FileType != FileTypeFile {
				t.Fatalf("fileType=%s: got a %s row in a size listing", c.name, f.FileType)
			}
			names = append(names, f.Name)
		}
		if !reflect.DeepEqual(names, want) {
			t.Fatalf("fileType=%s: got %v, want %v", c.name, names, want)
		}
	}
}

// --- keyword search ---------------------------------------------------------

func TestListFilesSortBySizeRanksKeywordMatches(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, false)

	insertSortTestChildSized(t, db, dir, "sizezzmarker-small.bin", user.UserNo, 100)
	insertSortTestChildSized(t, db, dir, "sizezzmarker-big.bin", user.UserNo, 900)
	insertSortTestChildSized(t, db, dir, "unrelated.bin", user.UserNo, 500)

	keyword := "sizezzmarker"
	got := listSortedFileNames(t, db, dir, ListFileReq{OrderBy: "size", Filename: &keyword}, user)
	want := []string{"sizezzmarker-big.bin", "sizezzmarker-small.bin"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("size with keyword: got %v, want %v", got, want)
	}
}

// --- comic directories ------------------------------------------------------

func TestListFilesComicDirSortBySizeReturnsNameAsc(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, true)

	// size order would be b, c, a; comic reading order is name asc
	insertSortTestChildSized(t, db, dir, "a.bin", user.UserNo, 100)
	insertSortTestChildSized(t, db, dir, "b.bin", user.UserNo, 900)
	insertSortTestChildSized(t, db, dir, "c.bin", user.UserNo, 500)

	want := []string{"a.bin", "b.bin", "c.bin"}
	for _, d := range []string{"", "asc", "desc"} {
		got := listSortedFileNames(t, db, dir, ListFileReq{OrderBy: "size", OrderDir: d}, user)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("comic dir size orderDir=%q: got %v, want %v", d, got, want)
		}
	}
}
