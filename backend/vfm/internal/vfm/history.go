package vfm

import (
	"errors"
	"sort"
	"time"

	"github.com/curtisnewbie/miso/encoding/json"
	red "github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

const (
	MaxBrowseHistoryLen = 100
)

type ListBrowseRecordRes struct {
	Time           util.ETime
	FileKey        string
	Name           string
	ThumbnailToken string
}

type BrowseRecord struct {
	Time    util.ETime
	FileKey string
}

type BrowseHistory struct {
	lock  *red.RLock
	user  common.User
	c     *redis.Client
	limit int
}

func (b BrowseHistory) Push(rail miso.Rail, fileKey string) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	var rbr BrowseRecord // last record
	rcmd := b.c.LRange(b.key(), -1, -1)
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
		if err := b.c.RPop(b.key()).Err(); err != nil {
			return err
		}
	}

	// push into queue
	pushed, err := json.SWriteJson(BrowseRecord{Time: util.Now(), FileKey: fileKey})
	if err != nil {
		return err
	}
	if err := b.c.RPush(b.key(), pushed).Err(); err != nil {
		return err
	}

	if rbr.FileKey != fileKey {
		lcmd := b.c.LLen(b.key())
		if lcmd.Err() != nil {
			return lcmd.Err()
		}

		cnt, err := lcmd.Result()
		if err != nil {
			return err
		}
		if cnt > MaxBrowseHistoryLen {
			return b.c.LPop(b.key()).Err()
		}
	}

	b.c.Expire(b.key(), time.Hour*48)
	return nil
}

func (b BrowseHistory) key() string {
	return "vfm:browse:history:" + b.user.UserNo
}

func (b BrowseHistory) List(rail miso.Rail) ([]BrowseRecord, error) {
	cmd := b.c.LRange(b.key(), 0, -1)
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
	keys := util.FastDistinct(util.MapTo(l, func(br BrowseRecord) string { return br.FileKey }))
	ffi, err := queryFileFstoreInfo(db, keys)
	if err != nil {
		return []ListBrowseRecordRes{}, err
	}

	res := util.MapTo(l, func(br BrowseRecord) ListBrowseRecordRes {
		var name string = ""
		var thumbnailToken string = ""
		if fi, ok := ffi[br.FileKey]; ok {
			name = fi.Name
			if fi.Thumbnail != "" {
				tkn, err := GetFstoreTmpToken(rail, fi.Thumbnail, "")
				if err != nil {
					rail.Errorf("Failed to generate browse history thumbnail token, thumbnail_file_id: %v, %v", fi.Thumbnail, err)
				} else {
					thumbnailToken = tkn
				}
			}
		}
		return ListBrowseRecordRes{
			Time:           br.Time,
			FileKey:        br.FileKey,
			Name:           name,
			ThumbnailToken: thumbnailToken,
		}
	})
	sort.Slice(res, func(i, j int) bool {
		return res[j].Time.UnixMilli() < res[i].Time.UnixMilli()
	})
	return res, nil
}

type RecordBrowseHistoryReq struct {
	FileKey string `valid:"notEmpty"`
}

func RecordBrowseHistory(rail miso.Rail, user common.User, req RecordBrowseHistoryReq) error {
	return NewBrowseHistory(rail, user).Push(rail, req.FileKey)
}