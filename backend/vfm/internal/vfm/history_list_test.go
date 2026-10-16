package vfm

import (
	"testing"

	"github.com/curtisnewbie/miso/miso"
	"gorm.io/gorm"
)

// insertCursorRow seeds a dir_resume_cursor row with an explicit updated_at so a
// test can control which visit is the directory's most recent one.
func insertCursorRow(t *testing.T, db *gorm.DB, userNo, dirKey, orderBy, orderDir, fileKey, updatedAt string) {
	t.Helper()
	if err := db.Exec("INSERT INTO dir_resume_cursor (user_no, dir_key, order_by, order_dir, file_key, updated_at, deleted) VALUES (?,?,?,?,?,?,0)",
		userNo, dirKey, orderBy, orderDir, fileKey, updatedAt).Error; err != nil {
		t.Fatal(err)
	}
}

func findDirRecord(records []DirBrowseRecord, dirKey string) *DirBrowseRecord {
	for i := range records {
		if records[i].DirKey == dirKey {
			return &records[i]
		}
	}
	return nil
}

func countDirRecords(records []DirBrowseRecord, dirKey string) int {
	n := 0
	for i := range records {
		if records[i].DirKey == dirKey {
			n++
		}
	}
	return n
}

// A directory read under several sorts appears exactly once in the list.
func TestListDirBrowseHistoryOneEntryPerDirectory(t *testing.T) {
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
	if n := countDirRecords(records, dirKey); n != 1 {
		t.Fatalf("expected exactly one entry for the directory, got %d", n)
	}
}

// The single entry carries the sort and file of the directory's most recent visit.
func TestListDirBrowseHistoryCarriesMostRecentVisitSortAndFile(t *testing.T) {
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
	got := findDirRecord(records, dirKey)
	if got == nil {
		t.Fatal("expected an entry for the directory")
	}
	if got.FileKey != fileB {
		t.Fatalf("expected file of the most recent visit %q, got %q", fileB, got.FileKey)
	}
	if got.OrderBy != SortByTime || got.OrderDir != "desc" {
		t.Fatalf("expected most recent visit sort time/desc, got %q/%q", got.OrderBy, got.OrderDir)
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
	if got := findDirRecord(records, dirKey); got != nil {
		t.Fatalf("deleted directory should not appear in the list, got %+v", got)
	}
}
