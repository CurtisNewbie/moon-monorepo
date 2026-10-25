package vfm

import (
	"reflect"
	"testing"

	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/randutil"
	"gorm.io/gorm"
)

// The filename-search listing tests seed their own rows under a per-run unique
// uploader_no and a name prefix, and remove exactly those rows on cleanup, so a
// peer test run against the same MySQL schema is never touched.
const (
	fsTestUserName = "fnsrch-test"
	fsNamePrefix   = "fnsrch-"
)

// fsTestUserNo is unique per test run, so two concurrent runs of this file seed
// and clean up disjoint rows.
var fsTestUserNo = "UE07FS" + randutil.ERand(10)

func fsTestUser() flow.User {
	return flow.User{UserNo: fsTestUserNo, Username: fsTestUserName}
}

// fsSeedRow describes one seeded file_info row; zero values are the normal
// visible, not-deleted, non-sensitive FILE row.
type fsSeedRow struct {
	Name      string
	Parent    string
	FileType  string
	Size      int64
	SeqKey    string
	Sensitive string
	Hidden    bool
	LogicDel  bool
	IsDel     int
}

// fsCleanup removes every row this file seeds. Register it before seeding.
func fsCleanup(t *testing.T, db *gorm.DB) {
	t.Helper()
	t.Cleanup(func() {
		if err := db.Exec("DELETE FROM file_info WHERE uploader_no = ? AND name LIKE ?",
			fsTestUserNo, fsNamePrefix+"%").Error; err != nil {
			t.Errorf("cleanup seeded rows: %v", err)
		}
	})
}

func fsSeed(t *testing.T, db *gorm.DB, row fsSeedRow) string {
	t.Helper()
	if row.FileType == "" {
		row.FileType = FileTypeFile
	}
	if row.Sensitive == "" {
		row.Sensitive = "N"
	}
	hidden, logicDel := 0, 0
	if row.Hidden {
		hidden = 1
	}
	if row.LogicDel {
		logicDel = 1
	}
	uuid := fsNamePrefix + "row-" + randutil.ERand(12)
	err := db.Exec(`INSERT INTO file_info
		(name, uuid, size_in_bytes, file_type, parent_file, thumbnail, uploader_no, uploader_name,
		 seq_key, sensitive_mode, hidden, is_logic_deleted, is_del)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		row.Name, uuid, row.Size, row.FileType, row.Parent, "", fsTestUserNo, fsTestUserName,
		row.SeqKey, row.Sensitive, hidden, logicDel, row.IsDel).Error
	if err != nil {
		t.Fatal(err)
	}
	return uuid
}

// fsListNames runs the listing seam and returns the row names in listing order.
func fsListNames(t *testing.T, db *gorm.DB, req ListFileReq) []string {
	t.Helper()
	req.Page = miso.Paging{Page: 1, Limit: 100}
	res, err := ListFiles(miso.EmptyRail(), db, req, fsTestUser())
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0, len(res.Payload))
	for _, f := range res.Payload {
		names = append(names, f.Name)
	}
	return names
}

func fsPtr(s string) *string { return &s }

func fsAssertNames(t *testing.T, got, want []string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

// A fragment of a word finds the file.
func TestListFilesKeywordMatchesNameFragment(t *testing.T) {
	db := testMySQLConn(t)
	fsCleanup(t, db)
	dir := fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "dir", FileType: FileTypeDir})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "holiday.mp4", Parent: dir, Size: 10})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "other.mp4", Parent: dir, Size: 20})

	got := fsListNames(t, db, ListFileReq{Filename: fsPtr("holid"), OrderBy: SortByName})
	fsAssertNames(t, got, []string{fsNamePrefix + "holiday.mp4"})
}

// Terms match in any order, and every term is required.
func TestListFilesKeywordTermsMatchInAnyOrder(t *testing.T) {
	db := testMySQLConn(t)
	fsCleanup(t, db)
	dir := fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "dir", FileType: FileTypeDir})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "holiday-trip-2024.mp4", Parent: dir, Size: 10})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "2024-holiday-my-trip.mp4", Parent: dir, Size: 20})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "holiday-only.mp4", Parent: dir, Size: 30})

	got := fsListNames(t, db, ListFileReq{Filename: fsPtr("trip holiday"), OrderBy: SortByName})
	fsAssertNames(t, got, []string{
		fsNamePrefix + "2024-holiday-my-trip.mp4",
		fsNamePrefix + "holiday-trip-2024.mp4",
	})

	got = fsListNames(t, db, ListFileReq{Filename: fsPtr("my holiday"), OrderBy: SortByName})
	fsAssertNames(t, got, []string{fsNamePrefix + "2024-holiday-my-trip.mp4"})
}

// Typed spaces, dashes and underscores all stand for any separator, and a name
// typed without separators finds the separated name.
func TestListFilesKeywordSeparatorsAreInterchangeable(t *testing.T) {
	db := testMySQLConn(t)
	fsCleanup(t, db)
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "my_holiday.mp4", Size: 10})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "my holiday.mp4", Size: 20})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "unrelated.mp4", Size: 30})

	want := []string{fsNamePrefix + "my holiday.mp4", fsNamePrefix + "my_holiday.mp4"}
	for _, keyword := range []string{"my holiday", "my-holiday", "myholiday", "my_holiday"} {
		got := fsListNames(t, db, ListFileReq{Filename: fsPtr(keyword), OrderBy: SortByName})
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("keyword %q: got %#v, want %#v", keyword, got, want)
		}
	}
}

// Short terms narrow the list instead of being dropped.
func TestListFilesKeywordShortTermFilters(t *testing.T) {
	db := testMySQLConn(t)
	fsCleanup(t, db)
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "my_holiday.mp4", Size: 10})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "other.mp4", Size: 20})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "剧集01.mp4", Size: 30})

	got := fsListNames(t, db, ListFileReq{Filename: fsPtr("my"), OrderBy: SortByName})
	fsAssertNames(t, got, []string{fsNamePrefix + "my_holiday.mp4"})

	got = fsListNames(t, db, ListFileReq{Filename: fsPtr("01"), OrderBy: SortByName})
	fsAssertNames(t, got, []string{fsNamePrefix + "剧集01.mp4"})
}

// A Chinese fragment finds the name without any segmentation.
func TestListFilesKeywordChineseFragment(t *testing.T) {
	db := testMySQLConn(t)
	fsCleanup(t, db)
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "剧集01.mp4", Size: 10})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "other.mp4", Size: 20})

	got := fsListNames(t, db, ListFileReq{Filename: fsPtr("剧集"), OrderBy: SortByName})
	fsAssertNames(t, got, []string{fsNamePrefix + "剧集01.mp4"})
}

// Letters-and-digits names match typed as one run, in either case.
func TestListFilesKeywordLettersAndDigits(t *testing.T) {
	db := testMySQLConn(t)
	fsCleanup(t, db)
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "IMG_2024_02.jpg", Size: 10})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "S01E02.mkv", Size: 20})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "unrelated.bin", Size: 30})

	got := fsListNames(t, db, ListFileReq{Filename: fsPtr("IMG2024"), OrderBy: SortByName})
	fsAssertNames(t, got, []string{fsNamePrefix + "IMG_2024_02.jpg"})

	got = fsListNames(t, db, ListFileReq{Filename: fsPtr("img2024"), OrderBy: SortByName})
	fsAssertNames(t, got, []string{fsNamePrefix + "IMG_2024_02.jpg"})

	got = fsListNames(t, db, ListFileReq{Filename: fsPtr("S01E02"), OrderBy: SortByName})
	fsAssertNames(t, got, []string{fsNamePrefix + "S01E02.mkv"})
}

// Every term is required, so a second term narrows the result.
func TestListFilesKeywordRequiresEveryTerm(t *testing.T) {
	db := testMySQLConn(t)
	fsCleanup(t, db)
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "holiday-only.mp4", Size: 10})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "holiday-trip.mp4", Size: 20})

	got := fsListNames(t, db, ListFileReq{Filename: fsPtr("holiday"), OrderBy: SortByName})
	fsAssertNames(t, got, []string{fsNamePrefix + "holiday-only.mp4", fsNamePrefix + "holiday-trip.mp4"})

	got = fsListNames(t, db, ListFileReq{Filename: fsPtr("holiday trip"), OrderBy: SortByName})
	fsAssertNames(t, got, []string{fsNamePrefix + "holiday-trip.mp4"})
}

// A literal % or _ cannot reach LIKE: it is a separator, never a wildcard. The
// "zz---q" row proves _ was not passed through as a single-character wildcard.
func TestListFilesKeywordPercentAndUnderscoreAreNotWildcards(t *testing.T) {
	db := testMySQLConn(t)
	fsCleanup(t, db)
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "zz---q.bin", Size: 10})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "zz-queue.bin", Size: 20})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "queue.bin", Size: 30})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "zzz.bin", Size: 40})

	want := []string{fsNamePrefix + "zz---q.bin", fsNamePrefix + "zz-queue.bin"}
	for _, keyword := range []string{"zz_q", "zz%q"} {
		got := fsListNames(t, db, ListFileReq{Filename: fsPtr(keyword), OrderBy: SortByName})
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("keyword %q: got %#v, want %#v", keyword, got, want)
		}
	}
}

// A keyword of separators only is absent entirely: it leaves the root listing
// unfiltered and still root-scoped.
func TestListFilesSeparatorOnlyKeywordLeavesRootListingUnchanged(t *testing.T) {
	db := testMySQLConn(t)
	fsCleanup(t, db)
	dir := fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "dir", FileType: FileTypeDir})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "root.txt", Size: 10})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "child-holiday.txt", Parent: dir, Size: 20})

	want := []string{fsNamePrefix + "dir", fsNamePrefix + "root.txt"}

	got := fsListNames(t, db, ListFileReq{OrderBy: SortByName})
	fsAssertNames(t, got, want)

	got = fsListNames(t, db, ListFileReq{Filename: fsPtr(" --__%% "), OrderBy: SortByName})
	fsAssertNames(t, got, want)
}

// Keyword scope is unchanged: directory-scoped inside a directory, account-wide
// at root.
func TestListFilesKeywordScopeUnchanged(t *testing.T) {
	db := testMySQLConn(t)
	fsCleanup(t, db)
	dirA := fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "dir-a", FileType: FileTypeDir})
	dirB := fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "dir-b", FileType: FileTypeDir})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "scope-a.txt", Parent: dirA, Size: 10})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "scope-b.txt", Parent: dirB, Size: 20})

	got := fsListNames(t, db, ListFileReq{
		Filename:   fsPtr("scope"),
		ParentFile: fsPtr(dirA),
		OrderBy:    SortByName,
	})
	fsAssertNames(t, got, []string{fsNamePrefix + "scope-a.txt"})

	got = fsListNames(t, db, ListFileReq{Filename: fsPtr("scope"), OrderBy: SortByName})
	fsAssertNames(t, got, []string{fsNamePrefix + "scope-a.txt", fsNamePrefix + "scope-b.txt"})
}

// Hidden, logically deleted and is_del rows stay out; sensitive rows are still
// shown unless the sensitive filter asks otherwise, exactly as before.
func TestListFilesKeywordKeepsExistingRowFilters(t *testing.T) {
	db := testMySQLConn(t)
	fsCleanup(t, db)
	dir := fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "dir", FileType: FileTypeDir})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "match-visible.txt", Parent: dir, Size: 10})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "match-hidden.txt", Parent: dir, Size: 20, Hidden: true})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "match-logic-deleted.txt", Parent: dir, Size: 30, LogicDel: true})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "match-row-deleted.txt", Parent: dir, Size: 40, IsDel: 1})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "match-sensitive.txt", Parent: dir, Size: 50, Sensitive: "Y"})

	got := fsListNames(t, db, ListFileReq{
		Filename:   fsPtr("match"),
		ParentFile: fsPtr(dir),
		OrderBy:    SortByName,
	})
	fsAssertNames(t, got, []string{
		fsNamePrefix + "match-sensitive.txt",
		fsNamePrefix + "match-visible.txt",
	})

	got = fsListNames(t, db, ListFileReq{
		Filename:   fsPtr("match"),
		ParentFile: fsPtr(dir),
		Sensitive:  func() *bool { b := true; return &b }(),
		OrderBy:    SortByName,
	})
	fsAssertNames(t, got, []string{fsNamePrefix + "match-visible.txt"})
}

// A keyword does not change any ordering clause: time still sorts by row id,
// and name, custom and size keep their own order.
func TestListFilesKeywordKeepsOrderingClauses(t *testing.T) {
	sortPreTest(t) // custom ordering bootstraps through Redis
	db := testMySQLConn(t)
	fsCleanup(t, db)
	dir := fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "dir", FileType: FileTypeDir})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "order-c.bin", Parent: dir, Size: 300})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "order-a.bin", Parent: dir, Size: 100})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "order-b.bin", Parent: dir, Size: 200})
	fsSeed(t, db, fsSeedRow{Name: fsNamePrefix + "nomatch.bin", Parent: dir, Size: 999})

	cases := []struct {
		name     string
		orderBy  string
		orderDir string
		want     []string
	}{
		{"name asc", SortByName, "", []string{fsNamePrefix + "order-a.bin", fsNamePrefix + "order-b.bin", fsNamePrefix + "order-c.bin"}},
		{"name desc", SortByName, "desc", []string{fsNamePrefix + "order-c.bin", fsNamePrefix + "order-b.bin", fsNamePrefix + "order-a.bin"}},
		{"size desc", SortBySize, "", []string{fsNamePrefix + "order-c.bin", fsNamePrefix + "order-b.bin", fsNamePrefix + "order-a.bin"}},
		{"size asc", SortBySize, "asc", []string{fsNamePrefix + "order-a.bin", fsNamePrefix + "order-b.bin", fsNamePrefix + "order-c.bin"}},
		{"time desc", SortByTime, "", []string{fsNamePrefix + "order-b.bin", fsNamePrefix + "order-a.bin", fsNamePrefix + "order-c.bin"}},
		{"time asc", SortByTime, "asc", []string{fsNamePrefix + "order-c.bin", fsNamePrefix + "order-a.bin", fsNamePrefix + "order-b.bin"}},
		{"custom", SortByCustom, "", []string{fsNamePrefix + "order-b.bin", fsNamePrefix + "order-a.bin", fsNamePrefix + "order-c.bin"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := fsListNames(t, db, ListFileReq{
				Filename:   fsPtr("order"),
				ParentFile: fsPtr(dir),
				OrderBy:    c.orderBy,
				OrderDir:   c.orderDir,
			})
			fsAssertNames(t, got, c.want)
		})
	}
}
