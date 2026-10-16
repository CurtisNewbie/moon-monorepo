package vfm

import (
	"errors"
	"sort"
	"time"

	"github.com/curtisnewbie/miso/errs"
	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	red "github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/atom"
	"github.com/curtisnewbie/miso/util/json"
	"github.com/curtisnewbie/miso/util/slutil"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	MaxBrowseHistoryLen int64         = 200
	BrowseHistoryTTL    time.Duration = time.Hour * 24 * 7

	TableDirResumeCursor = "dir_resume_cursor"
)

type ListBrowseRecordRes struct {
	Time           atom.Time `json:"time"`
	FileKey        string    `json:"fileKey"`
	Name           string    `json:"name"`
	ThumbnailToken string    `json:"thumbnailToken"`
	Deleted        bool      `json:"deleted"`
}

type BrowseRecord struct {
	Time    atom.Time
	FileKey string
}

type BrowseHistory struct {
	lock  *red.RLock
	user  flow.User
	c     *redis.Client
	limit int64
	ttl   time.Duration
}

func (b BrowseHistory) Push(rail miso.Rail, fileKey string) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	var rbr BrowseRecord // last record
	rcmd := b.c.LRange(rail.Context(), b.key(), -1, -1)
	if rcmd.Err() != nil {
		if !errors.Is(rcmd.Err(), redis.Nil) {
			return rcmd.Err()
		}
	} else {
		rcmdl, err := rcmd.Result()
		if err != nil {
			return err
		}
		if len(rcmdl) > 0 {
			if err := json.SParseJson(rcmdl[0], &rbr); err != nil {
				return err
			}
		}
	}

	// drop it to update the browse time
	if rbr.FileKey == fileKey {
		if err := b.c.RPop(rail.Context(), b.key()).Err(); err != nil {
			return err
		}
	}

	// push into queue
	pushed, err := json.SWriteJson(BrowseRecord{Time: atom.Now(), FileKey: fileKey})
	if err != nil {
		return err
	}
	if err := b.c.RPush(rail.Context(), b.key(), pushed).Err(); err != nil {
		return err
	}

	if rbr.FileKey != fileKey {
		lcmd := b.c.LLen(rail.Context(), b.key())
		if lcmd.Err() != nil {
			return lcmd.Err()
		}

		cnt, err := lcmd.Result()
		if err != nil {
			return err
		}
		if cnt > b.limit {
			return b.c.LPop(rail.Context(), b.key()).Err()
		}
	}

	b.c.Expire(rail.Context(), b.key(), b.ttl)
	return nil
}

func (b BrowseHistory) key() string {
	return "vfm:browse:history:" + b.user.UserNo
}

func (b BrowseHistory) List(rail miso.Rail) ([]BrowseRecord, error) {
	cmd := b.c.LRange(rail.Context(), b.key(), 0, -1)
	if cmd.Err() != nil {
		return []BrowseRecord{}, cmd.Err()
	}
	rs, err := cmd.Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []BrowseRecord{}, nil
		}
		return []BrowseRecord{}, err
	}
	var br []BrowseRecord = make([]BrowseRecord, 0, len(rs))
	for _, r := range rs {
		var b BrowseRecord
		if err := json.SParseJson(r, &b); err != nil {
			return []BrowseRecord{}, err
		}
		br = append(br, b)
	}
	return br, nil
}

func NewBrowseHistory(rail miso.Rail, user flow.User) BrowseHistory {
	bh := BrowseHistory{
		user:  user,
		lock:  red.NewRLockf(rail, "vfm:browse:history:lock:%v", user.UserNo),
		c:     red.GetRedis(),
		limit: MaxBrowseHistoryLen,
		ttl:   BrowseHistoryTTL,
	}
	return bh
}

func ListBrowseHistory(rail miso.Rail, db *gorm.DB, user flow.User) ([]ListBrowseRecordRes, error) {
	l, err := NewBrowseHistory(rail, user).List(rail)
	if err != nil {
		return nil, err
	}
	if len(l) < 1 {
		return []ListBrowseRecordRes{}, nil
	}
	keys := slutil.FastDistinct(slutil.MapTo(l, func(br BrowseRecord) string { return br.FileKey }))
	ffi, err := queryFileFstoreInfo(rail, db, keys)
	if err != nil {
		return []ListBrowseRecordRes{}, err
	}

	res := slutil.MapTo(l, func(br BrowseRecord) ListBrowseRecordRes {
		r := ListBrowseRecordRes{
			Time:    br.Time,
			FileKey: br.FileKey,
		}
		if fi, ok := ffi[br.FileKey]; ok {
			r.Name = fi.Name
			if !fi.IsLogicDeleted && fi.Thumbnail != "" {
				tkn, err := GetFstoreTmpToken(rail, fi.Thumbnail, "")
				if err != nil {
					rail.Errorf("Failed to generate browse history thumbnail token, thumbnail_file_id: %v, %v", fi.Thumbnail, err)
				} else {
					r.ThumbnailToken = tkn
				}
			}
			r.Deleted = fi.IsLogicDeleted
		}
		return r
	})
	sort.Slice(res, func(i, j int) bool {
		return res[j].Time.UnixMilli() < res[i].Time.UnixMilli()
	})
	return res, nil
}

type RecordBrowseHistoryReq struct {
	FileKey string `valid:"notEmpty" json:"fileKey"`
}

func RecordBrowseHistory(rail miso.Rail, db *gorm.DB, user flow.User, req RecordBrowseHistoryReq) error {
	f, ok, err := findFile(rail, db, req.FileKey)
	if err != nil {
		return err
	}
	if !ok {
		return ErrUnknown.WithInternalMsg("File is not found, %v", req.FileKey)
	}

	// only record files that are directly owned by user to prevent access control issue
	if f.UploaderNo != user.UserNo {
		rail.Debugf("Ignore recording browse history operation, file (%v) not owned by user (%v)", req.FileKey, user.Username)
		return nil
	}

	return NewBrowseHistory(rail, user).Push(rail, req.FileKey)
}

type RecordDirLastPageReq struct {
	DirKey   string `valid:"notEmpty" json:"dirKey"`
	FileKey  string `valid:"notEmpty" json:"fileKey"`
	OrderBy  string `json:"orderBy" desc:"sort field, 'name', 'time' or 'custom', empty means default"`
	OrderDir string `json:"orderDir" desc:"sort direction, 'asc' or 'desc', empty means default"`
}

type GetDirLastPageReq struct {
	DirKey   string `form:"dirKey"`
	OrderBy  string `form:"orderBy"`
	OrderDir string `form:"orderDir"`
}

type DirLastPageRes struct {
	FileKey string `json:"fileKey"`
}

type dirLastPageEntry struct {
	FileKey string `json:"fileKey"`
	Time    int64  `json:"time"`
}

// dirResumeCursor maps a row of dir_resume_cursor: one resume cursor per user,
// directory and sort order.
type dirResumeCursor struct {
	Id        int
	UserNo    string
	DirKey    string
	OrderBy   string
	OrderDir  string
	FileKey   string
	CreatedBy string
	UpdatedBy string
}

func dirResumeCursorLock(rail miso.Rail, userNo, dirKey, orderBy, orderDir string) *red.RLock {
	return red.NewRLockf(rail, "vfm:dir:resume:cursor:%v:%v:%v:%v", userNo, dirKey, orderBy, orderDir)
}

// RecordDirLastPage records the file a user last read in a directory under the
// sort order in force. Every directory records; only the directory's owner and
// directories that exist are accepted. The comic-only gate is gone, so the
// cursor key is (user, directory, resolved sort).
func RecordDirLastPage(rail miso.Rail, db *gorm.DB, user flow.User, req RecordDirLastPageReq) error {
	f, ok, err := findFile(rail, db, req.DirKey)
	if err != nil {
		return err
	}
	if !ok {
		return nil // silently ignore, directory not found
	}
	if f.UploaderNo != user.UserNo {
		return nil // silently ignore, not owned by user
	}
	if f.FileType != FileTypeDir {
		return nil // silently ignore, not a directory
	}

	// resolve the effective sort server-side so a comic directory records and
	// reads under one key even though its client sends the dropdown default.
	// Reuse the directory row already loaded above instead of loading it again.
	orderBy, orderDir := resolveEffectiveSortForDir(req.OrderBy, req.OrderDir, f.IsComic)

	lock := dirResumeCursorLock(rail, user.UserNo, req.DirKey, orderBy, orderDir)
	locked, err := lock.TryLock(red.WithBackoff(3 * time.Second))
	if err != nil {
		return errs.Wrap(err)
	}
	if !locked {
		rail.Warnf("Skip recording dir resume cursor, lock is busy, userNo: %v, dirKey: %v", user.UserNo, req.DirKey)
		return nil // lock busy, the client re-records on the next preview
	}
	defer lock.Unlock()

	return upsertDirResumeCursor(rail, db, user, req.DirKey, orderBy, orderDir, req.FileKey)
}

// findDirResumeCursor loads the live cursor row for one (user, directory, sort).
// It returns ok=false when none exists. Callers own any error wrapping.
func findDirResumeCursor(rail miso.Rail, db *gorm.DB, userNo, dirKey, orderBy, orderDir string) (dirResumeCursor, bool, error) {
	var cur dirResumeCursor
	ok, err := dbquery.NewQuery(rail, db).
		Table(TableDirResumeCursor).
		Eq("user_no", userNo).
		Eq("dir_key", dirKey).
		Eq("order_by", orderBy).
		Eq("order_dir", orderDir).
		Eq("deleted", 0).
		Limit(1).
		ScanAny(&cur)
	return cur, ok, err
}

// upsertDirResumeCursor reads the cursor row for the sort in force and inserts
// or updates it. Callers hold the (user, directory, sort) lock so the
// read-then-write leaves exactly one row.
func upsertDirResumeCursor(rail miso.Rail, db *gorm.DB, user flow.User, dirKey, orderBy, orderDir, fileKey string) error {
	cur, ok, err := findDirResumeCursor(rail, db, user.UserNo, dirKey, orderBy, orderDir)
	if err != nil {
		return errs.Wrapf(err, "failed to query dir resume cursor, userNo: %v, dirKey: %v", user.UserNo, dirKey)
	}

	if ok {
		if err := dbquery.NewQuery(rail, db).
			Table(TableDirResumeCursor).
			Set("file_key", fileKey).
			Set("updated_by", user.Username).
			Eq("id", cur.Id).
			UpdateAny(); err != nil {
			return errs.Wrapf(err, "failed to update dir resume cursor, userNo: %v, dirKey: %v", user.UserNo, dirKey)
		}
		return nil
	}

	row := dirResumeCursor{
		UserNo:    user.UserNo,
		DirKey:    dirKey,
		OrderBy:   orderBy,
		OrderDir:  orderDir,
		FileKey:   fileKey,
		CreatedBy: user.Username,
		UpdatedBy: user.Username,
	}
	if err := dbquery.NewQuery(rail, db).
		Table(TableDirResumeCursor).
		Omit("id").
		CreateAny(&row); err != nil {
		return errs.Wrapf(err, "failed to create dir resume cursor, userNo: %v, dirKey: %v", user.UserNo, dirKey)
	}
	return nil
}

// GetDirLastPage returns the file a user last read in a directory under the
// requested sort, resolved the same way as recording. It returns the raw file,
// empty when no cursor exists.
func GetDirLastPage(rail miso.Rail, db *gorm.DB, user flow.User, req GetDirLastPageReq) (DirLastPageRes, error) {
	es, err := resolveEffectiveSort(rail, db, req.OrderBy, req.OrderDir, req.DirKey)
	if err != nil {
		return DirLastPageRes{}, err
	}

	cur, ok, err := findDirResumeCursor(rail, db, user.UserNo, req.DirKey, es.OrderBy, es.OrderDir)
	if err != nil {
		return DirLastPageRes{}, errs.Wrapf(err, "failed to query dir resume cursor, userNo: %v, dirKey: %v", user.UserNo, req.DirKey)
	}
	if !ok {
		fileKey, found, err := legacyDirResumeCursorFallback(rail, db, user, req.DirKey, es)
		if err != nil {
			return DirLastPageRes{}, err
		}
		if !found {
			return DirLastPageRes{FileKey: ""}, nil
		}
		return DirLastPageRes{FileKey: fileKey}, nil
	}
	return DirLastPageRes{FileKey: cur.FileKey}, nil
}

// legacyDirResumeCursorFallback carries a position recorded before the cursor
// moved to MySQL into the table. It runs only when the directory is comic — the
// only kind the old Redis last-page hash ever held — and only on a table miss.
// A hit is written into the table under the resolved sort and then returned.
// The Redis hash is read-only legacy state: it receives no writes here.
//
// es carries the directory the resolver already loaded for name/time sorts. Only
// a custom-sort lookup (which loads nothing) falls back to loading it here.
func legacyDirResumeCursorFallback(rail miso.Rail, db *gorm.DB, user flow.User, dirKey string, es effectiveSort) (string, bool, error) {
	dir := es.Dir
	if !es.DirKnown {
		d, ok, err := findFile(rail, db, dirKey)
		if err != nil {
			return "", false, err
		}
		if !ok {
			return "", false, nil
		}
		dir = d
	}
	if dir.FileType != FileTypeDir || !dir.IsComic {
		return "", false, nil
	}
	// mirror the recording gate: only carry a legacy position forward for a
	// directory the user owns
	if dir.UploaderNo != user.UserNo {
		return "", false, nil
	}

	fileKey, found, err := readLegacyDirLastPage(rail, user.UserNo, dirKey)
	if err != nil || !found {
		return "", false, err
	}

	lock := dirResumeCursorLock(rail, user.UserNo, dirKey, es.OrderBy, es.OrderDir)
	locked, err := lock.TryLock(red.WithBackoff(3 * time.Second))
	if err != nil {
		return "", false, errs.Wrap(err)
	}
	if !locked {
		rail.Warnf("Skip writing back legacy dir resume cursor, lock is busy, userNo: %v, dirKey: %v", user.UserNo, dirKey)
		return fileKey, true, nil // lock busy, serve the legacy position without the write-back
	}
	defer lock.Unlock()

	if err := upsertDirResumeCursor(rail, db, user, dirKey, es.OrderBy, es.OrderDir, fileKey); err != nil {
		return "", false, err
	}
	return fileKey, true, nil
}

// readLegacyDirLastPage reads a comic directory's last-read file from the old
// Redis hash vfm:dir:lastpage:<userNo>. It never writes to the hash.
func readLegacyDirLastPage(rail miso.Rail, userNo, dirKey string) (string, bool, error) {
	c := red.GetRedis()
	key := "vfm:dir:lastpage:" + userNo
	raw, err := c.HGet(rail.Context(), key, dirKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", false, nil
		}
		return "", false, errs.Wrapf(err, "failed to read legacy dir last page, userNo: %v, dirKey: %v", userNo, dirKey)
	}

	var e dirLastPageEntry
	if err := json.SParseJson(raw, &e); err != nil {
		rail.Warnf("Failed to parse legacy dir last page entry, userNo: %v, dirKey: %v, %v", userNo, dirKey, err)
		return "", false, nil
	}
	if e.FileKey == "" {
		return "", false, nil
	}
	return e.FileKey, true, nil
}

type DirBrowseRecord struct {
	DirKey         string    `json:"dirKey"`
	Name           string    `json:"name"`
	ThumbnailToken string    `json:"thumbnailToken,omitempty"`
	FileKey        string    `json:"fileKey"`
	OrderBy        string    `json:"orderBy"`
	OrderDir       string    `json:"orderDir"`
	Time           atom.Time `json:"time"`
}

// dirBrowseCursorRow is a dir_resume_cursor row used to build the browse
// history, carrying the visit's sort and file.
type dirBrowseCursorRow struct {
	DirKey    string
	OrderBy   string
	OrderDir  string
	FileKey   string
	UpdatedAt atom.Time
}

// listDirBrowseHistory returns the top directory browse history entries for a user,
// sorted by last-viewed time descending, enriched with file metadata from MySQL.
// Each directory appears once, taken from its most recently updated cursor row,
// and directories marked deleted are left out.
func listDirBrowseHistory(rail miso.Rail, db *gorm.DB, user flow.User) ([]DirBrowseRecord, error) {
	const maxEntries = 200

	// newest cursor rows first; the user_updated_idx (user_no, updated_at) index serves this
	var rows []dirBrowseCursorRow
	if err := dbquery.NewQuery(rail, db).
		Table(TableDirResumeCursor).
		Eq("user_no", user.UserNo).
		Eq("deleted", 0).
		OrderDesc("updated_at").
		OrderDesc("id").
		Select("dir_key", "order_by", "order_dir", "file_key", "updated_at").
		ScanVal(&rows); err != nil {
		return nil, errs.Wrapf(err, "failed to query dir resume cursor, userNo: %v", user.UserNo)
	}
	if len(rows) == 0 {
		return []DirBrowseRecord{}, nil
	}

	// one entry per directory: keep the first (most recent) row seen for each dir
	seen := make(map[string]bool, len(rows))
	entries := make([]dirBrowseCursorRow, 0, len(rows))
	dirKeys := make([]string, 0, len(rows))
	for _, r := range rows {
		if seen[r.DirKey] {
			continue
		}
		seen[r.DirKey] = true
		entries = append(entries, r)
		dirKeys = append(dirKeys, r.DirKey)
	}

	// exclude directories marked deleted (the directory's own file_info row)
	var liveDirs []struct {
		Uuid string
	}
	if err := dbquery.NewQuery(rail, db).
		Table("file_info").
		Where("uuid IN ?", dirKeys).
		Eq("is_logic_deleted", DelN).
		Eq("is_del", 0).
		Select("uuid").
		ScanVal(&liveDirs); err != nil {
		return nil, errs.Wrapf(err, "failed to query live directories for browse history, userNo: %v", user.UserNo)
	}
	live := make(map[string]bool, len(liveDirs))
	for _, d := range liveDirs {
		live[d.Uuid] = true
	}
	kept := entries[:0]
	for _, e := range entries {
		if live[e.DirKey] {
			kept = append(kept, e)
		}
	}
	entries = kept

	// take top N
	if len(entries) > maxEntries {
		entries = entries[:maxEntries]
	}
	if len(entries) == 0 {
		return []DirBrowseRecord{}, nil
	}

	// collect dirKeys and fileKeys for enrichment
	dirKeys = dirKeys[:0]
	var fileKeys []string
	for _, e := range entries {
		dirKeys = append(dirKeys, e.DirKey)
		if e.FileKey != "" {
			fileKeys = append(fileKeys, e.FileKey)
		}
	}

	// enrich with file metadata from MySQL (both dirs and files)
	allKeys := append(dirKeys, fileKeys...)
	ffi, err := queryFileFstoreInfo(rail, db, allKeys)
	if err != nil {
		return nil, err
	}

	// use batch dir-thumbnail API as fallback (first image inside each dir)
	var thumbnailMap map[string]string
	if len(dirKeys) > 0 {
		thumbs, err := BatchFetchDirThumbnail(rail, db, BatchFetchDirThumbnailReq{DirFileKeys: dirKeys}, user)
		if err != nil {
			rail.Errorf("Failed to batch fetch dir thumbnails for browse history: %v", err)
		} else {
			thumbnailMap = make(map[string]string, len(thumbs))
			for _, t := range thumbs {
				if t.FstoreToken != "" {
					thumbnailMap[t.DirFileKey] = t.FstoreToken
				}
			}
		}
	}

	records := make([]DirBrowseRecord, 0, len(entries))
	for _, e := range entries {
		r := DirBrowseRecord{
			DirKey:   e.DirKey,
			FileKey:  e.FileKey,
			OrderBy:  e.OrderBy,
			OrderDir: e.OrderDir,
			Time:     e.UpdatedAt,
		}
		if fi, ok := ffi[e.DirKey]; ok {
			r.Name = fi.Name
		}
		// Prefer thumbnail of the last-viewed file over dir thumbnail
		if e.FileKey != "" {
			if fi, ok := ffi[e.FileKey]; ok && fi.Thumbnail != "" && !fi.IsLogicDeleted {
				tkn, err := GetFstoreTmpToken(rail, fi.Thumbnail, "")
				if err != nil {
					rail.Errorf("Failed to generate thumbnail token for fileKey: %v, %v", e.FileKey, err)
				} else {
					r.ThumbnailToken = tkn
				}
			}
		}
		// Fallback to dir thumbnail
		if r.ThumbnailToken == "" {
			if tkn, ok := thumbnailMap[e.DirKey]; ok {
				r.ThumbnailToken = tkn
			}
		}
		records = append(records, r)
	}

	return records, nil
}

func ListDirBrowseHistory(rail miso.Rail, db *gorm.DB, user flow.User) ([]DirBrowseRecord, error) {
	return listDirBrowseHistory(rail, db, user)
}
