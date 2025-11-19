package vfm

import (
	"errors"
	"sort"
	"time"

	red "github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
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
	user  common.User
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

func NewBrowseHistory(rail miso.Rail, user common.User) BrowseHistory {
	bh := BrowseHistory{
		user:  user,
		lock:  red.NewRLockf(rail, "vfm:browse:history:lock:%v", user.UserNo),
		c:     red.GetRedis(),
		limit: MaxBrowseHistoryLen,
		ttl:   BrowseHistoryTTL,
	}
	return bh
}

func ListBrowseHistory(rail miso.Rail, db *gorm.DB, user common.User) ([]ListBrowseRecordRes, error) {
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

func RecordBrowseHistory(rail miso.Rail, db *gorm.DB, user common.User, req RecordBrowseHistoryReq) error {
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
