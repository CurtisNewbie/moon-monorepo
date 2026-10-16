package vfm

import (
	"reflect"
	"testing"

	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/atom"
	"github.com/curtisnewbie/miso/util/json"
)

// legacyRedisKey is the old per-user last-page hash key. It is read-only legacy
// state: the feature reads it once and never writes to it.
func legacyRedisKey(userNo string) string {
	return "vfm:dir:lastpage:" + userNo
}

// setLegacyDirLastPageRedis seeds one entry into the old Redis hash, cleaning it
// up after the test.
func setLegacyDirLastPageRedis(t *testing.T, userNo, dirKey, fileKey string) {
	t.Helper()
	c := redis.GetRedis()
	raw, err := json.SWriteJson(dirLastPageEntry{FileKey: fileKey, Time: atom.Now().UnixMilli()})
	if err != nil {
		t.Fatal(err)
	}
	if err := c.HSet(miso.EmptyRail().Context(), legacyRedisKey(userNo), dirKey, string(raw)).Err(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		c.HDel(miso.EmptyRail().Context(), legacyRedisKey(userNo), dirKey)
	})
}

// deleteLegacyDirLastPageRedis removes one entry from the old Redis hash so a
// later read can only be served from the table.
func deleteLegacyDirLastPageRedis(t *testing.T, userNo, dirKey string) {
	t.Helper()
	c := redis.GetRedis()
	if err := c.HDel(miso.EmptyRail().Context(), legacyRedisKey(userNo), dirKey).Err(); err != nil {
		t.Fatal(err)
	}
}

func legacyRedisHash(t *testing.T, userNo string) map[string]string {
	t.Helper()
	c := redis.GetRedis()
	m, err := c.HGetAll(miso.EmptyRail().Context(), legacyRedisKey(userNo)).Result()
	if err != nil {
		t.Fatal(err)
	}
	return m
}

// A comic directory whose position exists only in the old Redis store still
// resumes from that entry, and once read the position is served from the table
// even after the Redis entry is gone.
func TestDirResumeCursorLegacyRedisFallbackForComicResumesAndPersists(t *testing.T) {
	rail := miso.EmptyRail()
	db := historyPreTest(t)
	user := testUser()

	dirKey := insertTestComicDir(t, db, user.UserNo)
	cleanupDir(t, db, dirKey)
	setLegacyDirLastPageRedis(t, user.UserNo, dirKey, "legacy-comic-f")

	res, err := GetDirLastPage(rail, db, user, GetDirLastPageReq{DirKey: dirKey})
	if err != nil {
		t.Fatal(err)
	}
	if res.FileKey != "legacy-comic-f" {
		t.Fatalf("expected legacy cursor legacy-comic-f, got %q", res.FileKey)
	}

	// the Redis entry is gone; the cursor must now come from the table
	deleteLegacyDirLastPageRedis(t, user.UserNo, dirKey)

	res, err = GetDirLastPage(rail, db, user, GetDirLastPageReq{DirKey: dirKey})
	if err != nil {
		t.Fatal(err)
	}
	if res.FileKey != "legacy-comic-f" {
		t.Fatalf("expected persisted cursor legacy-comic-f after Redis removal, got %q", res.FileKey)
	}
}

// The legacy fallback must not add anything to the old Redis hash: existing
// entries are read, never modified, and no new entries are created.
func TestDirResumeCursorLegacyRedisFallbackWritesNothingToRedis(t *testing.T) {
	rail := miso.EmptyRail()
	db := historyPreTest(t)
	user := testUser()

	dirKey := insertTestComicDir(t, db, user.UserNo)
	cleanupDir(t, db, dirKey)
	setLegacyDirLastPageRedis(t, user.UserNo, dirKey, "legacy-comic-f")

	before := legacyRedisHash(t, user.UserNo)

	if _, err := GetDirLastPage(rail, db, user, GetDirLastPageReq{DirKey: dirKey}); err != nil {
		t.Fatal(err)
	}

	after := legacyRedisHash(t, user.UserNo)
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("legacy Redis hash changed, before: %+v, after: %+v", before, after)
	}
}

// A non-comic directory never consults the old Redis store, even when an entry
// for it happens to exist there.
func TestDirResumeCursorNonComicIgnoresLegacyRedis(t *testing.T) {
	rail := miso.EmptyRail()
	db := historyPreTest(t)
	user := testUser()

	dirKey := insertTestDir(t, db, user.UserNo)
	cleanupDir(t, db, dirKey)
	setLegacyDirLastPageRedis(t, user.UserNo, dirKey, "stale-non-comic-f")

	res, err := GetDirLastPage(rail, db, user, GetDirLastPageReq{DirKey: dirKey, OrderBy: SortByName, OrderDir: "asc"})
	if err != nil {
		t.Fatal(err)
	}
	if res.FileKey != "" {
		t.Fatalf("expected non-comic directory to ignore legacy Redis, got %q", res.FileKey)
	}
}

// A directory with neither a table row nor a Redis entry opens at the first page
// without error.
func TestDirResumeCursorMissingEverywhereOpensEmpty(t *testing.T) {
	rail := miso.EmptyRail()
	db := historyPreTest(t)
	user := testUser()

	dirKey := insertTestComicDir(t, db, user.UserNo)
	cleanupDir(t, db, dirKey)

	res, err := GetDirLastPage(rail, db, user, GetDirLastPageReq{DirKey: dirKey})
	if err != nil {
		t.Fatal(err)
	}
	if res.FileKey != "" {
		t.Fatalf("expected empty cursor, got %q", res.FileKey)
	}
}
