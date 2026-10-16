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
	// reads under one key even though its client sends the dropdown default
	orderBy, orderDir, err := resolveEffectiveSort(rail, db, req.OrderBy, req.OrderDir, req.DirKey)
	if err != nil {
		return err
	}

	lock := dirResumeCursorLock(rail, user.UserNo, req.DirKey, orderBy, orderDir)
	if err := lock.Lock(); err != nil {
		return errs.Wrap(err)
	}
	defer lock.Unlock()

	return upsertDirResumeCursor(rail, db, user, req.DirKey, orderBy, orderDir, req.FileKey)
}

// upsertDirResumeCursor reads the cursor row for the sort in force and inserts
// or updates it. Callers hold the (user, directory, sort) lock so the
// read-then-write leaves exactly one row.
func upsertDirResumeCursor(rail miso.Rail, db *gorm.DB, user flow.User, dirKey, orderBy, orderDir, fileKey string) error {
	var cur dirResumeCursor
	ok, err := dbquery.NewQuery(rail, db).
		Table("dir_resume_cursor").
		Eq("user_no", user.UserNo).
		Eq("dir_key", dirKey).
		Eq("order_by", orderBy).
		Eq("order_dir", orderDir).
		Eq("deleted", 0).
		Limit(1).
		ScanAny(&cur)
	if err != nil {
		return errs.Wrapf(err, "failed to query dir resume cursor, userNo: %v, dirKey: %v", user.UserNo, dirKey)
	}

	if ok {
		if err := dbquery.NewQuery(rail, db).
			Table("dir_resume_cursor").
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
		Table("dir_resume_cursor").
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
	orderBy, orderDir, err := resolveEffectiveSort(rail, db, req.OrderBy, req.OrderDir, req.DirKey)
	if err != nil {
		return DirLastPageRes{}, err
	}

	var cur dirResumeCursor
	ok, err := dbquery.NewQuery(rail, db).
		Table("dir_resume_cursor").
		Eq("user_no", user.UserNo).
		Eq("dir_key", req.DirKey).
		Eq("order_by", orderBy).
		Eq("order_dir", orderDir).
		Eq("deleted", 0).
		Limit(1).
		ScanAny(&cur)
	if err != nil {
		return DirLastPageRes{}, errs.Wrapf(err, "failed to query dir resume cursor, userNo: %v, dirKey: %v", user.UserNo, req.DirKey)
	}
	if !ok {
		return DirLastPageRes{FileKey: ""}, nil
	}
	return DirLastPageRes{FileKey: cur.FileKey}, nil
}

type DirBrowseRecord struct {
	DirKey         string    `json:"dirKey"`
	Name           string    `json:"name"`
	ThumbnailToken string    `json:"thumbnailToken,omitempty"`
	FileKey        string    `json:"fileKey"`
	Time           atom.Time `json:"time"`
}

// listDirBrowseHistory returns the top directory browse history entries for a user,
// sorted by last-viewed time descending, enriched with file metadata from MySQL.
func listDirBrowseHistory(rail miso.Rail, db *gorm.DB, user flow.User) ([]DirBrowseRecord, error) {
	const maxEntries = 200

	c := red.GetRedis()
	key := "vfm:dir:lastpage:" + user.UserNo

	all, err := c.HGetAll(rail.Context(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []DirBrowseRecord{}, nil
		}
		return nil, err
	}
	if len(all) == 0 {
		return []DirBrowseRecord{}, nil
	}

	// parse all entries
	type entry struct {
		DirKey string
		dirLastPageEntry
	}
	entries := make([]entry, 0, len(all))
	for dirKey, raw := range all {
		var e dirLastPageEntry
		if err := json.SParseJson(raw, &e); err != nil {
			continue // skip corrupted
		}
		entries = append(entries, entry{DirKey: dirKey, dirLastPageEntry: e})
	}

	// sort by time descending
	sort.Slice(entries, func(i, j int) bool {
		return entries[j].Time < entries[i].Time
	})

	// take top N
	if len(entries) > maxEntries {
		entries = entries[:maxEntries]
	}

	// collect dirKeys and fileKeys for enrichment
	dirKeys := make([]string, len(entries))
	var fileKeys []string
	for i, e := range entries {
		dirKeys[i] = e.DirKey
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
			DirKey:  e.DirKey,
			FileKey: e.FileKey,
			Time:    atom.WrapTime(time.UnixMilli(e.Time)),
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
