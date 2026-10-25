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

// rootSizeUser returns a throwaway uploader whose rows these tests seed and
// remove. The user number is unique, so cleanup never touches a row another
// test (or a peer test run) seeded.
func rootSizeUser(t *testing.T, db *gorm.DB) flow.User {
	t.Helper()
	user := flow.User{UserNo: "UE07RS" + randutil.ERand(10), Username: "rs07-uploader"}
	t.Cleanup(func() {
		if err := db.Exec("DELETE FROM file_info WHERE uploader_no = ?", user.UserNo).Error; err != nil {
			t.Errorf("cleanup seeded rows for %s: %v", user.UserNo, err)
		}
	})
	return user
}

// rootSizeRow describes one file_info row a test seeds. Empty fileType means a
// file, empty sensitive means "N", and the deleted/hidden flags default to 0.
type rootSizeRow struct {
	parent       string
	name         string
	size         int64
	fileType     string
	seqKey       string
	hidden       int
	isDel        int
	logicDeleted int
	sensitive    string
}

// insertRootSizeRow inserts one seeded row and returns its uuid.
func insertRootSizeRow(t *testing.T, db *gorm.DB, userNo string, r rootSizeRow) string {
	t.Helper()
	fileType := r.fileType
	if fileType == "" {
		fileType = FileTypeFile
	}
	sensitive := r.sensitive
	if sensitive == "" {
		sensitive = "N"
	}
	uuid := "rs07-" + randutil.ERand(12)
	if err := db.Exec(`INSERT INTO file_info (name, uuid, size_in_bytes, file_type, parent_file, thumbnail, uploader_no, seq_key, hidden, is_del, is_logic_deleted, sensitive_mode)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`,
		r.name, uuid, r.size, fileType, r.parent, "", userNo, r.seqKey, r.hidden, r.isDel, r.logicDeleted, sensitive).Error; err != nil {
		t.Fatal(err)
	}
	return uuid
}

// listRootSize lists with no directory named, the way the page asks at root.
func listRootSize(t *testing.T, db *gorm.DB, req ListFileReq, user flow.User) []ListedFile {
	t.Helper()
	req.Page = miso.Paging{Page: 1, Limit: 100}
	res, err := ListFiles(miso.EmptyRail(), db, req, user)
	if err != nil {
		t.Fatal(err)
	}
	return res.Payload
}

// rootSizeNames lists with no directory named and returns the payload names.
func rootSizeNames(t *testing.T, db *gorm.DB, req ListFileReq, user flow.User) []string {
	t.Helper()
	payload := listRootSize(t, db, req, user)
	names := make([]string, 0, len(payload))
	for _, f := range payload {
		names = append(names, f.Name)
	}
	return names
}

// --- root listing -----------------------------------------------------------

func TestListRootSizeSpansEveryDirectory(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := rootSizeUser(t, db)

	// directory rows carry bytes too, so they would outrank the files if leaked
	dirA := insertRootSizeRow(t, db, user.UserNo, rootSizeRow{name: "rs07-dir-a", fileType: FileTypeDir, size: 800})
	dirB := insertRootSizeRow(t, db, user.UserNo, rootSizeRow{name: "rs07-dir-b", fileType: FileTypeDir, size: 400})

	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{parent: dirA, name: "rs07-small.bin", size: 100})
	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{parent: dirB, name: "rs07-large.bin", size: 900})
	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{name: "rs07-medium.bin", size: 500})

	got := rootSizeNames(t, db, ListFileReq{OrderBy: SortBySize}, user)
	want := []string{"rs07-large.bin", "rs07-medium.bin", "rs07-small.bin"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("root size desc: got %v, want %v", got, want)
	}
}

func TestListRootSizeFlipsDirectionAndBreaksTiesById(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := rootSizeUser(t, db)

	// one equal-size file in each of three directories, so across directories
	// only the row id orders them; inserted a, b, c so id order is a, b, c
	dirA := insertRootSizeRow(t, db, user.UserNo, rootSizeRow{name: "rs07-tie-dir-a", fileType: FileTypeDir})
	dirB := insertRootSizeRow(t, db, user.UserNo, rootSizeRow{name: "rs07-tie-dir-b", fileType: FileTypeDir})
	dirC := insertRootSizeRow(t, db, user.UserNo, rootSizeRow{name: "rs07-tie-dir-c", fileType: FileTypeDir})

	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{parent: dirA, name: "rs07-a.bin", size: 500})
	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{parent: dirB, name: "rs07-b.bin", size: 500})
	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{parent: dirC, name: "rs07-c.bin", size: 500})

	cases := []struct {
		orderDir string
		want     []string
	}{
		{"", []string{"rs07-c.bin", "rs07-b.bin", "rs07-a.bin"}},     // omitted → desc, id desc
		{"desc", []string{"rs07-c.bin", "rs07-b.bin", "rs07-a.bin"}}, // id desc crosses directories
		{"asc", []string{"rs07-a.bin", "rs07-b.bin", "rs07-c.bin"}},  // id asc
	}
	for _, c := range cases {
		got := rootSizeNames(t, db, ListFileReq{OrderBy: SortBySize, OrderDir: c.orderDir}, user)
		if !reflect.DeepEqual(got, c.want) {
			t.Fatalf("root size orderDir=%q: got %v, want %v", c.orderDir, got, c.want)
		}

		// the same order across pages, with no repeat and no skip
		paged := make([]string, 0, len(c.want))
		for page := 1; page <= 2; page++ {
			res, err := ListFiles(miso.EmptyRail(), db, ListFileReq{
				Page:     miso.Paging{Page: page, Limit: 2},
				OrderBy:  SortBySize,
				OrderDir: c.orderDir,
			}, user)
			if err != nil {
				t.Fatal(err)
			}
			for _, f := range res.Payload {
				paged = append(paged, f.Name)
			}
		}
		if !reflect.DeepEqual(paged, c.want) {
			t.Fatalf("root size orderDir=%q across pages: got %v, want %v", c.orderDir, paged, c.want)
		}
	}
}

func TestListRootSizeRanksKeywordMatchesAcrossDirectories(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := rootSizeUser(t, db)

	dirA := insertRootSizeRow(t, db, user.UserNo, rootSizeRow{name: "rs07-kw-dir-a", fileType: FileTypeDir})
	dirB := insertRootSizeRow(t, db, user.UserNo, rootSizeRow{name: "rs07-kw-dir-b", fileType: FileTypeDir})
	dirC := insertRootSizeRow(t, db, user.UserNo, rootSizeRow{name: "rs07-kw-dir-c", fileType: FileTypeDir})

	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{parent: dirA, name: "rs07kwmatch-big.bin", size: 900})
	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{parent: dirB, name: "rs07kwmatch-small.bin", size: 100})
	// bigger than both matches, but it never matches the keyword
	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{parent: dirC, name: "rs07-unrelated.bin", size: 9999})

	keyword := "rs07kwmatch"
	got := rootSizeNames(t, db, ListFileReq{OrderBy: SortBySize, Filename: &keyword}, user)
	want := []string{"rs07kwmatch-big.bin", "rs07kwmatch-small.bin"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("root size with keyword: got %v, want %v", got, want)
	}
}

func TestListSizeInsideDirectoryStaysScoped(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := rootSizeUser(t, db)

	dirA := insertRootSizeRow(t, db, user.UserNo, rootSizeRow{name: "rs07-scope-dir-a", fileType: FileTypeDir})
	dirB := insertRootSizeRow(t, db, user.UserNo, rootSizeRow{name: "rs07-scope-dir-b", fileType: FileTypeDir})

	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{parent: dirA, name: "rs07-inside-small.bin", size: 100})
	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{parent: dirA, name: "rs07-inside-big.bin", size: 900})
	// a larger sibling outside the named directory that must not appear
	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{parent: dirB, name: "rs07-sibling-huge.bin", size: 9999})

	parent := dirA
	got := rootSizeNames(t, db, ListFileReq{OrderBy: SortBySize, ParentFile: &parent}, user)
	want := []string{"rs07-inside-big.bin", "rs07-inside-small.bin"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("size inside a directory: got %v, want %v", got, want)
	}
}

func TestListRootSizeReturnsFilesOnly(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := rootSizeUser(t, db)

	dirA := insertRootSizeRow(t, db, user.UserNo, rootSizeRow{name: "rs07-only-dir-a", fileType: FileTypeDir, size: 8000})
	dirB := insertRootSizeRow(t, db, user.UserNo, rootSizeRow{name: "rs07-only-dir-b", fileType: FileTypeDir, size: 9000})

	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{parent: dirA, name: "rs07-only-small.bin", size: 100})
	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{parent: dirB, name: "rs07-only-big.bin", size: 900})

	empty := ""
	dirFilter := FileTypeDir
	fileFilter := FileTypeFile
	cases := []struct {
		name     string
		fileType *string
	}{
		{"absent", nil},
		{"all", &empty},
		{"directory", &dirFilter},
		{"file", &fileFilter},
	}
	for _, c := range cases {
		payload := listRootSize(t, db, ListFileReq{OrderBy: SortBySize, FileType: c.fileType}, user)
		names := make([]string, 0, len(payload))
		for _, f := range payload {
			if f.FileType != FileTypeFile {
				t.Fatalf("fileType=%s: got a %s row in a size listing", c.name, f.FileType)
			}
			names = append(names, f.Name)
		}
		want := []string{"rs07-only-big.bin", "rs07-only-small.bin"}
		if !reflect.DeepEqual(names, want) {
			t.Fatalf("fileType=%s: got %v, want %v", c.name, names, want)
		}
	}
}

func TestListRootOtherOrdersStayRootLevel(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := rootSizeUser(t, db)

	dirA := insertRootSizeRow(t, db, user.UserNo, rootSizeRow{name: "rs07-root-dir-a", fileType: FileTypeDir})
	dirB := insertRootSizeRow(t, db, user.UserNo, rootSizeRow{name: "rs07-root-dir-b", fileType: FileTypeDir})

	// files directly at root, plus children of different sizes in each directory
	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{name: "rs07-root-a.bin", size: 100, seqKey: "2"})
	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{name: "rs07-root-b.bin", size: 900, seqKey: "1"})
	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{parent: dirA, name: "rs07-child-a.bin", size: 500})
	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{parent: dirB, name: "rs07-child-b.bin", size: 700})

	// Every order below is deterministic: time orders by file_type then row id
	// (dirs sort before files), name by the name, custom by seq_key then
	// file_type then id desc.
	cases := []struct {
		orderBy, orderDir string
		want              []string
	}{
		{SortByTime, "desc", []string{"rs07-root-dir-b", "rs07-root-dir-a", "rs07-root-b.bin", "rs07-root-a.bin"}},
		{SortByTime, "asc", []string{"rs07-root-dir-a", "rs07-root-dir-b", "rs07-root-a.bin", "rs07-root-b.bin"}},
		{SortByName, "asc", []string{"rs07-root-a.bin", "rs07-root-b.bin", "rs07-root-dir-a", "rs07-root-dir-b"}},
		{SortByName, "desc", []string{"rs07-root-dir-b", "rs07-root-dir-a", "rs07-root-b.bin", "rs07-root-a.bin"}},
		{SortByCustom, "", []string{"rs07-root-dir-b", "rs07-root-dir-a", "rs07-root-b.bin", "rs07-root-a.bin"}},
	}
	for _, c := range cases {
		got := rootSizeNames(t, db, ListFileReq{OrderBy: c.orderBy, OrderDir: c.orderDir}, user)
		if !reflect.DeepEqual(got, c.want) {
			t.Fatalf("root orderBy=%q orderDir=%q: got %v, want %v", c.orderBy, c.orderDir, got, c.want)
		}
	}
}

func TestListRootSizeExcludesHiddenDeletedAndSensitiveRows(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := rootSizeUser(t, db)

	dirA := insertRootSizeRow(t, db, user.UserNo, rootSizeRow{name: "rs07-vis-dir-a", fileType: FileTypeDir})
	dirB := insertRootSizeRow(t, db, user.UserNo, rootSizeRow{name: "rs07-vis-dir-b", fileType: FileTypeDir})

	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{parent: dirA, name: "rs07-visible.bin", size: 900})
	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{parent: dirA, name: "rs07-hidden.bin", size: 800, hidden: 1})
	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{parent: dirA, name: "rs07-logic-deleted.bin", size: 700, logicDeleted: 1})
	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{parent: dirB, name: "rs07-isdel.bin", size: 600, isDel: 1})
	insertRootSizeRow(t, db, user.UserNo, rootSizeRow{parent: dirB, name: "rs07-sensitive.bin", size: 500, sensitive: "Y"})

	// the sensitive filter is what keeps sensitive rows out, as today
	sensitive := true
	got := rootSizeNames(t, db, ListFileReq{OrderBy: SortBySize, Sensitive: &sensitive}, user)
	if want := []string{"rs07-visible.bin"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("root size with sensitive filter: got %v, want %v", got, want)
	}

	// hidden and deleted rows never appear, whatever the sensitive filter says
	got = rootSizeNames(t, db, ListFileReq{OrderBy: SortBySize}, user)
	if want := []string{"rs07-visible.bin", "rs07-sensitive.bin"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("root size without sensitive filter: got %v, want %v", got, want)
	}
}
