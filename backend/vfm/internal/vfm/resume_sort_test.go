package vfm

import (
	"reflect"
	"testing"

	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/randutil"
	"gorm.io/gorm"
)

// sortPreTest initializes Redis, which ListFiles reaches for custom ordering.
func sortPreTest(t *testing.T) {
	t.Helper()
	if _, err := redis.InitRedisFromProp(miso.EmptyRail()); err != nil {
		t.Fatal(err)
	}
}

// --- fixtures ---------------------------------------------------------------

// insertSortTestDir inserts a directory row, optionally flagged as comic, and
// removes it (and its children) on cleanup.
func insertSortTestDir(t *testing.T, db *gorm.DB, userNo string, comic bool) string {
	t.Helper()
	uuid := "rst-dir-" + randutil.ERand(8)
	isComic := 0
	if comic {
		isComic = 1
	}
	if err := db.Exec(`INSERT INTO file_info (name, uuid, size_in_bytes, file_type, parent_file, thumbnail, uploader_no, is_comic)
		VALUES (?,?,?,?,?,?,?,?)`,
		"rst-dir", uuid, 0, FileTypeDir, "", "", userNo, isComic).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Exec("DELETE FROM file_info WHERE parent_file = ? OR uuid = ?", uuid, uuid)
	})
	return uuid
}

// insertSortTestChild inserts a file row under parent and returns its uuid.
func insertSortTestChild(t *testing.T, db *gorm.DB, parent, name, userNo, seqKey string) string {
	t.Helper()
	uuid := "rst-child-" + randutil.ERand(8)
	if err := db.Exec(`INSERT INTO file_info (name, uuid, size_in_bytes, file_type, parent_file, thumbnail, uploader_no, seq_key)
		VALUES (?,?,?,?,?,?,?,?)`,
		name, uuid, 0, FileTypeFile, parent, "", userNo, seqKey).Error; err != nil {
		t.Fatal(err)
	}
	return uuid
}

func listFileNames(t *testing.T, db *gorm.DB, dir, orderBy, orderDir string, user flow.User) []string {
	t.Helper()
	parent := dir
	res, err := ListFiles(miso.EmptyRail(), db, ListFileReq{
		Page:       miso.Paging{Page: 1, Limit: 100},
		ParentFile: &parent,
		OrderBy:    orderBy,
		OrderDir:   orderDir,
	}, user)
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0, len(res.Payload))
	for _, f := range res.Payload {
		names = append(names, f.Name)
	}
	return names
}

func calcPosition(t *testing.T, db *gorm.DB, fileKey, dir, orderBy, orderDir string, limit int, user flow.User) int {
	t.Helper()
	p, err := CalcFilePosition(miso.EmptyRail(), db, FilePositionReq{
		FileKey:    fileKey,
		ParentFile: dir,
		Limit:      limit,
		OrderBy:    orderBy,
		OrderDir:   orderDir,
	}, user)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// --- directory listing ------------------------------------------------------

func TestListFilesComicDirForcedNameAsc(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, true)

	// inserted c, a, b so id order differs from name order
	insertSortTestChild(t, db, dir, "c.jpg", user.UserNo, "")
	insertSortTestChild(t, db, dir, "a.jpg", user.UserNo, "")
	insertSortTestChild(t, db, dir, "b.jpg", user.UserNo, "")

	want := []string{"a.jpg", "b.jpg", "c.jpg"}
	cases := []struct{ orderBy, orderDir string }{
		{SortByTime, "desc"},
		{SortByTime, "asc"},
		{SortByName, "desc"},
		{SortByName, ""},
		{"", ""},
	}
	for _, c := range cases {
		got := listFileNames(t, db, dir, c.orderBy, c.orderDir, user)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("comic dir orderBy=%q orderDir=%q: got %v, want %v", c.orderBy, c.orderDir, got, want)
		}
	}

	// custom ordering is honored even on a comic dir (prior behavior preserved)
	gotCustom := listFileNames(t, db, dir, SortByCustom, "", user)
	wantCustom := []string{"b.jpg", "a.jpg", "c.jpg"} // seq_key all equal → file_type asc, id desc
	if !reflect.DeepEqual(gotCustom, wantCustom) {
		t.Fatalf("comic dir custom order: got %v, want %v", gotCustom, wantCustom)
	}
}

func TestListFilesNonComicRespectsRequestedSort(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, false)

	// ids: c=1, a=2, b=3
	insertSortTestChild(t, db, dir, "c.jpg", user.UserNo, "")
	insertSortTestChild(t, db, dir, "a.jpg", user.UserNo, "")
	insertSortTestChild(t, db, dir, "b.jpg", user.UserNo, "")

	cases := []struct {
		orderBy, orderDir string
		want              []string
	}{
		{SortByName, "asc", []string{"a.jpg", "b.jpg", "c.jpg"}},
		{SortByName, "desc", []string{"c.jpg", "b.jpg", "a.jpg"}},
		{SortByName, "", []string{"a.jpg", "b.jpg", "c.jpg"}},     // omitted dir → name asc
		{SortByTime, "asc", []string{"c.jpg", "a.jpg", "b.jpg"}},  // id asc
		{SortByTime, "desc", []string{"b.jpg", "a.jpg", "c.jpg"}}, // id desc
		{SortByTime, "", []string{"b.jpg", "a.jpg", "c.jpg"}},     // omitted dir → time desc
		{"", "", []string{"b.jpg", "a.jpg", "c.jpg"}},             // omitted field → time desc
	}
	for _, c := range cases {
		got := listFileNames(t, db, dir, c.orderBy, c.orderDir, user)
		if !reflect.DeepEqual(got, c.want) {
			t.Fatalf("orderBy=%q orderDir=%q: got %v, want %v", c.orderBy, c.orderDir, got, c.want)
		}
	}
}

func TestListFilesCustomOrderIgnoresDirection(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, false)

	insertSortTestChild(t, db, dir, "a.jpg", user.UserNo, "3")
	insertSortTestChild(t, db, dir, "b.jpg", user.UserNo, "2")
	insertSortTestChild(t, db, dir, "c.jpg", user.UserNo, "1")

	want := []string{"c.jpg", "b.jpg", "a.jpg"} // seq_key asc
	for _, d := range []string{"", "asc", "desc"} {
		got := listFileNames(t, db, dir, SortByCustom, d, user)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("custom orderDir=%q: got %v, want %v", d, got, want)
		}
	}
}

// --- file position ----------------------------------------------------------

func TestCalcFilePositionComicDirForcedNameAsc(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, true)

	// ids: c=1, a=2, b=3; name asc puts a first, time desc puts a second
	insertSortTestChild(t, db, dir, "c.jpg", user.UserNo, "")
	aKey := insertSortTestChild(t, db, dir, "a.jpg", user.UserNo, "")
	insertSortTestChild(t, db, dir, "b.jpg", user.UserNo, "")

	cases := []struct{ orderBy, orderDir string }{
		{SortByTime, "desc"},
		{SortByTime, "asc"},
		{SortByName, "desc"},
		{SortByName, ""},
		{"", ""},
	}
	for _, c := range cases {
		got := calcPosition(t, db, aKey, dir, c.orderBy, c.orderDir, 1, user)
		if got != 1 {
			t.Fatalf("comic dir orderBy=%q orderDir=%q: got page %v, want 1 (forced name asc)", c.orderBy, c.orderDir, got)
		}
	}
}

func TestCalcFilePositionDefaultDirection(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, false)

	// ids: a=1, b=2, c=3
	aKey := insertSortTestChild(t, db, dir, "a.jpg", user.UserNo, "")
	insertSortTestChild(t, db, dir, "b.jpg", user.UserNo, "")
	insertSortTestChild(t, db, dir, "c.jpg", user.UserNo, "")

	cases := []struct {
		orderBy, orderDir string
		want              int
	}{
		{SortByName, "", 1}, // omitted dir → name asc (a first)
		{SortByName, "asc", 1},
		{SortByName, "desc", 3},
		{SortByTime, "", 3}, // omitted dir → time desc (a last)
		{SortByTime, "asc", 1},
		{SortByTime, "desc", 3},
		{"", "", 3}, // omitted field → time desc
	}
	for _, c := range cases {
		got := calcPosition(t, db, aKey, dir, c.orderBy, c.orderDir, 1, user)
		if got != c.want {
			t.Fatalf("orderBy=%q orderDir=%q: got page %v, want %v", c.orderBy, c.orderDir, got, c.want)
		}
	}
}

func TestCalcFilePositionCustomIgnoresDirection(t *testing.T) {
	sortPreTest(t)
	db := testMySQLConn(t)
	user := testUser()
	dir := insertSortTestDir(t, db, user.UserNo, false)

	// seq_key asc order is c, b, a → a is third
	aKey := insertSortTestChild(t, db, dir, "a.jpg", user.UserNo, "3")
	insertSortTestChild(t, db, dir, "b.jpg", user.UserNo, "2")
	insertSortTestChild(t, db, dir, "c.jpg", user.UserNo, "1")

	names := listFileNames(t, db, dir, SortByCustom, "", user)
	if !reflect.DeepEqual(names, []string{"c.jpg", "b.jpg", "a.jpg"}) {
		t.Fatalf("custom listing: got %v", names)
	}
	// a is last → position 3 regardless of direction
	for _, d := range []string{"", "asc", "desc"} {
		got := calcPosition(t, db, aKey, dir, SortByCustom, d, 1, user)
		if got != 3 {
			t.Fatalf("custom orderDir=%q: got page %v, want 3", d, got)
		}
	}
}
