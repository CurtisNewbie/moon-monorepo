package vfm

import (
	"testing"

	"github.com/curtisnewbie/miso/miso"
	"gorm.io/gorm"
)

// insertCursorRow seeds a dir_resume_cursor row with an explicit updated_at so a
// test can control the order in which the directory's visits appear.
func insertCursorRow(t *testing.T, db *gorm.DB, userNo, dirKey, orderBy, orderDir, fileKey, updatedAt string) {
	t.Helper()
	if err := db.Exec("INSERT INTO dir_resume_cursor (user_no, dir_key, order_by, order_dir, file_key, updated_at, deleted) VALUES (?,?,?,?,?,?,0)",
		userNo, dirKey, orderBy, orderDir, fileKey, updatedAt).Error; err != nil {
		t.Fatal(err)
	}
}

// dirRecords returns the list's entries for one directory, in list order.
func dirRecords(records []DirBrowseRecord, dirKey string) []DirBrowseRecord {
	var rs []DirBrowseRecord
	for i := range records {
		if records[i].DirKey == dirKey {
			rs = append(rs, records[i])
		}
	}
	return rs
}

// A directory read under several sorts appears once per sort, most recently
// updated first, with no grouping by directory.
func TestListDirBrowseHistoryOneEntryPerSort(t *testing.T) {
	rail := miso.EmptyRail()
	db := historyPreTest(t)
	user := testUser()

	dirKey := insertTestDir(t, db, user.UserNo)
	cleanupDir(t, db, dirKey)
	fileA := insertTestChild(t, db, dirKey, "a.txt", "", 0)
	fileB := insertTestChild(t, db, dirKey, "b.txt", "", 0)

	// older visit: name-ascending, then newer visit: time-descending
	insertCursorRow(t, db, user.UserNo, dirKey, SortByName, "asc", fileA, "2020-01-01 00:00:00")
	insertCursorRow(t, db, user.UserNo, dirKey, SortByTime, "desc", fileB, "2020-06-01 00:00:00")

	records, err := listDirBrowseHistory(rail, db, user)
	if err != nil {
		t.Fatal(err)
	}
	rs := dirRecords(records, dirKey)
	if len(rs) != 2 {
		t.Fatalf("expected one entry per sort, got %d entries for the directory", len(rs))
	}
	// newest visit first: time-descending, then name-ascending
	if rs[0].FileKey != fileB {
		t.Fatalf("expected the newest visit first, got file %q", rs[0].FileKey)
	}
	if rs[1].FileKey != fileA {
		t.Fatalf("expected the older visit second, got file %q", rs[1].FileKey)
	}
}

// Each entry carries the sort order and file of its own cursor row, so two
// entries for one directory differ by sort and file.
func TestListDirBrowseHistoryCarriesEachRowsSortAndFile(t *testing.T) {
	rail := miso.EmptyRail()
	db := historyPreTest(t)
	user := testUser()

	dirKey := insertTestDir(t, db, user.UserNo)
	cleanupDir(t, db, dirKey)
	fileA := insertTestChild(t, db, dirKey, "a.txt", "", 0)
	fileB := insertTestChild(t, db, dirKey, "b.txt", "", 0)

	// older visit: name-ascending, then newer visit: time-descending
	insertCursorRow(t, db, user.UserNo, dirKey, SortByName, "asc", fileA, "2020-01-01 00:00:00")
	insertCursorRow(t, db, user.UserNo, dirKey, SortByTime, "desc", fileB, "2020-06-01 00:00:00")

	records, err := listDirBrowseHistory(rail, db, user)
	if err != nil {
		t.Fatal(err)
	}
	bySort := make(map[string]DirBrowseRecord)
	for _, r := range dirRecords(records, dirKey) {
		bySort[r.OrderBy+"/"+r.OrderDir] = r
	}

	named, ok := bySort[SortByName+"/asc"]
	if !ok {
		t.Fatalf("expected an entry for the name/asc visit, got %v", bySort)
	}
	if named.FileKey != fileA {
		t.Fatalf("expected name/asc entry to carry file %q, got %q", fileA, named.FileKey)
	}

	timed, ok := bySort[SortByTime+"/desc"]
	if !ok {
		t.Fatalf("expected an entry for the time/desc visit, got %v", bySort)
	}
	if timed.FileKey != fileB {
		t.Fatalf("expected time/desc entry to carry file %q, got %q", fileB, timed.FileKey)
	}
}

// A directory marked deleted is left out of the list.
func TestListDirBrowseHistoryExcludesDeletedDirectory(t *testing.T) {
	rail := miso.EmptyRail()
	db := historyPreTest(t)
	user := testUser()

	dirKey := insertTestDir(t, db, user.UserNo)
	cleanupDir(t, db, dirKey)
	fileA := insertTestChild(t, db, dirKey, "a.txt", "", 0)
	insertCursorRow(t, db, user.UserNo, dirKey, SortByName, "asc", fileA, "2020-01-01 00:00:00")

	if err := db.Exec("UPDATE file_info SET is_logic_deleted = 1 WHERE uuid = ?", dirKey).Error; err != nil {
		t.Fatal(err)
	}

	records, err := listDirBrowseHistory(rail, db, user)
	if err != nil {
		t.Fatal(err)
	}
	if rs := dirRecords(records, dirKey); len(rs) != 0 {
		t.Fatalf("deleted directory should not appear in the list, got %+v", rs)
	}
}
