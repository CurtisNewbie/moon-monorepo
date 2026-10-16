package vfm

import (
	"fmt"
	"sync"
	"testing"

	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/randutil"
	"gorm.io/gorm"
)

// historyPreTest initializes Redis and returns a dedicated MySQL connection to
// the vfm schema, so the service-layer cursor tests run in isolation.
func historyPreTest(t *testing.T) *gorm.DB {
	t.Helper()
	if _, e := redis.InitRedisFromProp(miso.EmptyRail()); e != nil {
		t.Fatal(e)
	}
	return testMySQLConn(t)
}

// insertTestComicDir inserts a comic directory row owned by userNo.
func insertTestComicDir(t *testing.T, db *gorm.DB, userNo string) string {
	t.Helper()
	uuid := "tdt-comic-" + randutil.ERand(8)
	if err := db.Exec("INSERT INTO file_info (name, uuid, size_in_bytes, file_type, parent_file, thumbnail, uploader_no, is_comic) VALUES (?,?,?,?,?,?,?,?)",
		"test-comic-dir", uuid, 0, "DIR", "", "", userNo, 1).Error; err != nil {
		t.Fatal(err)
	}
	return uuid
}

func cleanupDir(t *testing.T, db *gorm.DB, dirKey string) {
	t.Helper()
	t.Cleanup(func() {
		db.Exec("DELETE FROM file_info WHERE parent_file = ? OR uuid = ?", dirKey, dirKey)
		db.Exec("DELETE FROM dir_resume_cursor WHERE dir_key = ?", dirKey)
	})
}

// Reading a directory under a second sort leaves the first sort's cursor intact,
// and each sort resumes independently.
func TestDirResumeCursorPerSortIndependent(t *testing.T) {
	rail := miso.EmptyRail()
	db := historyPreTest(t)
	user := testUser()

	dirKey := insertTestDir(t, db, user.UserNo)
	cleanupDir(t, db, dirKey)

	if err := RecordDirLastPage(rail, db, user, RecordDirLastPageReq{
		DirKey: dirKey, FileKey: "f-name", OrderBy: SortByName, OrderDir: "asc",
	}); err != nil {
		t.Fatal(err)
	}
	if err := RecordDirLastPage(rail, db, user, RecordDirLastPageReq{
		DirKey: dirKey, FileKey: "f-time", OrderBy: SortByTime, OrderDir: "desc",
	}); err != nil {
		t.Fatal(err)
	}

	res, err := GetDirLastPage(rail, db, user, GetDirLastPageReq{DirKey: dirKey, OrderBy: SortByName, OrderDir: "asc"})
	if err != nil {
		t.Fatal(err)
	}
	if res.FileKey != "f-name" {
		t.Fatalf("expected name-sort cursor f-name, got %q", res.FileKey)
	}

	res, err = GetDirLastPage(rail, db, user, GetDirLastPageReq{DirKey: dirKey, OrderBy: SortByTime, OrderDir: "desc"})
	if err != nil {
		t.Fatal(err)
	}
	if res.FileKey != "f-time" {
		t.Fatalf("expected time-sort cursor f-time, got %q", res.FileKey)
	}

	// reading the time sort again must not disturb the name sort
	res, err = GetDirLastPage(rail, db, user, GetDirLastPageReq{DirKey: dirKey, OrderBy: SortByName, OrderDir: "asc"})
	if err != nil {
		t.Fatal(err)
	}
	if res.FileKey != "f-name" {
		t.Fatalf("name-sort cursor was disturbed, got %q", res.FileKey)
	}
}

// Concurrent recordings for the same (user, directory, sort) leave exactly one cursor row.
func TestDirResumeCursorConcurrentRecordingLeavesOneRow(t *testing.T) {
	db := historyPreTest(t)
	user := testUser()

	dirKey := insertTestDir(t, db, user.UserNo)
	cleanupDir(t, db, dirKey)

	const n = 8
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			rail := miso.EmptyRail()
			err := RecordDirLastPage(rail, db, user, RecordDirLastPageReq{
				DirKey: dirKey, FileKey: fmt.Sprintf("f-%d", i), OrderBy: SortByName, OrderDir: "asc",
			})
			if err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()

	if firstErr != nil {
		t.Fatal(firstErr)
	}

	var cnt int
	if err := db.Raw("SELECT COUNT(*) FROM dir_resume_cursor WHERE user_no = ? AND dir_key = ? AND order_by = ? AND order_dir = ? AND deleted = 0",
		user.UserNo, dirKey, SortByName, "asc").Scan(&cnt).Error; err != nil {
		t.Fatal(err)
	}
	if cnt != 1 {
		t.Fatalf("expected exactly one cursor row, got %d", cnt)
	}
}

// A lookup for a directory with no cursor returns empty and no error.
func TestDirResumeCursorLookupMissingReturnsEmpty(t *testing.T) {
	rail := miso.EmptyRail()
	db := historyPreTest(t)
	user := testUser()

	dirKey := insertTestDir(t, db, user.UserNo)
	cleanupDir(t, db, dirKey)

	res, err := GetDirLastPage(rail, db, user, GetDirLastPageReq{DirKey: dirKey, OrderBy: SortByName, OrderDir: "asc"})
	if err != nil {
		t.Fatal(err)
	}
	if res.FileKey != "" {
		t.Fatalf("expected empty cursor, got %q", res.FileKey)
	}
}

// Recording is refused for a directory the user does not own.
func TestDirResumeCursorRefusedForForeignDirectory(t *testing.T) {
	rail := miso.EmptyRail()
	db := historyPreTest(t)
	user := testUser()

	dirKey := insertTestDir(t, db, "UE_OTHER_USER")
	cleanupDir(t, db, dirKey)

	if err := RecordDirLastPage(rail, db, user, RecordDirLastPageReq{
		DirKey: dirKey, FileKey: "f-x", OrderBy: SortByName, OrderDir: "asc",
	}); err != nil {
		t.Fatal(err)
	}

	res, err := GetDirLastPage(rail, db, user, GetDirLastPageReq{DirKey: dirKey, OrderBy: SortByName, OrderDir: "asc"})
	if err != nil {
		t.Fatal(err)
	}
	if res.FileKey != "" {
		t.Fatalf("expected no cursor for a directory the user does not own, got %q", res.FileKey)
	}
}

// Recording is refused for a directory that does not exist.
func TestDirResumeCursorRefusedForMissingDirectory(t *testing.T) {
	rail := miso.EmptyRail()
	db := historyPreTest(t)
	user := testUser()

	dirKey := "tdt-missing-" + randutil.ERand(8)
	cleanupDir(t, db, dirKey)

	if err := RecordDirLastPage(rail, db, user, RecordDirLastPageReq{
		DirKey: dirKey, FileKey: "f-x", OrderBy: SortByName, OrderDir: "asc",
	}); err != nil {
		t.Fatal(err)
	}

	res, err := GetDirLastPage(rail, db, user, GetDirLastPageReq{DirKey: dirKey, OrderBy: SortByName, OrderDir: "asc"})
	if err != nil {
		t.Fatal(err)
	}
	if res.FileKey != "" {
		t.Fatalf("expected no cursor for a missing directory, got %q", res.FileKey)
	}
}

// A comic directory resolves to name-ascending on both record and lookup even
// though the client sends its dropdown default, so both hit the same key.
func TestDirResumeCursorComicDirectoryUsesNameAscending(t *testing.T) {
	rail := miso.EmptyRail()
	db := historyPreTest(t)
	user := testUser()

	dirKey := insertTestComicDir(t, db, user.UserNo)
	cleanupDir(t, db, dirKey)

	// client sends its dropdown default: no OrderBy / OrderDir
	if err := RecordDirLastPage(rail, db, user, RecordDirLastPageReq{
		DirKey: dirKey, FileKey: "comic-f",
	}); err != nil {
		t.Fatal(err)
	}

	res, err := GetDirLastPage(rail, db, user, GetDirLastPageReq{DirKey: dirKey})
	if err != nil {
		t.Fatal(err)
	}
	if res.FileKey != "comic-f" {
		t.Fatalf("expected comic cursor comic-f, got %q", res.FileKey)
	}

	// a requested time sort collapses to the same name-ascending key
	res, err = GetDirLastPage(rail, db, user, GetDirLastPageReq{DirKey: dirKey, OrderBy: SortByTime, OrderDir: "desc"})
	if err != nil {
		t.Fatal(err)
	}
	if res.FileKey != "comic-f" {
		t.Fatalf("expected comic cursor comic-f under time request, got %q", res.FileKey)
	}
}
