package postbox

import (
	"fmt"
	"time"

	"github.com/curtisnewbie/event-pump/client"
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"github.com/curtisnewbie/user-vault/api"
	"gorm.io/gorm"
)

const (
	StatusInit   = "INIT"
	StatusOpened = "OPENED"
)

var (
	userNotifCountCache = redis.NewRCache[int]("postbox:notification:count", redis.RCacheConfig{Exp: time.Minute * 30})
)

func CreateNotification(rail miso.Rail, db *gorm.DB, req api.CreateNotificationReq, user common.User) error {
	if len(req.ReceiverUserNos) < 1 {
		return nil
	}

	// check whether the userNos are leegal
	req.ReceiverUserNos = util.Distinct(req.ReceiverUserNos)

	for _, u := range req.ReceiverUserNos {
		sr := SaveNotifiReq{
			UserNo:  u,
			Title:   req.Title,
			Message: req.Message,
		}
		if err := SaveNotification(rail, db, sr, user); err != nil {
			return fmt.Errorf("failed to save notification, %+v, %v", sr, err)
		}
	}

	return nil
}

type SaveNotifiReq struct {
	UserNo  string
	Title   string
	Message string
}

func SaveNotification(rail miso.Rail, db *gorm.DB, req SaveNotifiReq, user common.User) error {
	notifiNo := NotifiNo()
	err := db.Exec(`insert into notification (user_no, notifi_no, title, message, created_by) values (?, ?, ?, ?, ?)`,
		req.UserNo, notifiNo, req.Title, req.Message, user.Username).Error
	if err != nil {
		return fmt.Errorf("failed to save notifiication record, %+v", req)
	}
	return nil
}

func NotifiNo() string {
	return util.GenIdP("notif_")
}

type ListedNotification struct {
	Id         int
	NotifiNo   string
	Title      string
	Message    string
	Status     string
	CreateTime util.ETime
}

func QueryNotification(rail miso.Rail, db *gorm.DB, req QueryNotificationReq, user common.User) (miso.PageRes[ListedNotification], error) {
	return mysql.NewPageQuery[ListedNotification]().
		WithPage(req.Page).
		WithBaseQuery(func(tx *gorm.DB) *gorm.DB {
			tx = tx.Table("notification")
			tx = tx.Where("user_no = ?", user.UserNo)
			if req.Status != "" {
				tx = tx.Where("status = ?", req.Status)
			}
			return tx
		}).
		WithSelectQuery(func(tx *gorm.DB) *gorm.DB {
			tx = tx.Select("id, notifi_no, title, message, status, create_time").
				Order("id desc").
				Limit(req.Page.GetLimit()).
				Offset(req.Page.GetOffset())

			return tx
		}).
		Exec(rail, db)
}

func CachedCountNotification(rail miso.Rail, db *gorm.DB, user common.User) (int, error) {
	return userNotifCountCache.Get(rail, user.UserNo, func() (int, error) {
		return CountNotification(rail, db, user)
	})
}

func CountNotification(rail miso.Rail, db *gorm.DB, user common.User) (int, error) {
	var count int
	err := db.Table("notification").
		Select("count(*)").
		Where("user_no = ?", user.UserNo).
		Where("status = ?", StatusInit).
		Scan(&count).Error
	return count, err
}

func OpenNotification(rail miso.Rail, db *gorm.DB, req OpenNotificationReq, user common.User) error {
	return db.Exec(`UPDATE notification SET status = ?, updated_by = ? WHERE notifi_no = ? AND user_no = ?`,
		StatusOpened, user.Username, req.NotifiNo, user.UserNo).Error
}

func OpenAllNotification(rail miso.Rail, db *gorm.DB, req OpenNotificationReq, user common.User) error {
	var id int
	n, err := mysql.NewQuery(db).
		From("notification").
		Select("id").
		Eq("user_no", user.UserNo).
		Eq("notifi_no", req.NotifiNo).
		Scan(&id)
	if err != nil {
		return err
	}
	if n < 1 {
		return miso.NewErrf("Record not found")
	}

	return db.Exec(`UPDATE notification SET status = ?, updated_by = ? WHERE user_no = ? AND status = ? AND id <= ?`,
		StatusOpened, user.Username, user.UserNo, StatusInit, id).Error
}

func evictNotifCountCache(rail miso.Rail, t client.StreamEvent) error {
	userNo, ok := t.ColumnAfter("user_no")
	if !ok {
		return nil
	}
	rail.Infof("User notification changed, eventType: %v, %v", t.Type, userNo)
	if err := userNotifCountCache.Del(rail, userNo); err != nil {
		rail.Errorf("Failed to evict user notification count cache, %v, %v", userNo, err)
	}

	if c := redis.GetRedis().Publish(userNotifCountChangedChannel, userNo); c.Err() != nil {
		rail.Errorf("Failed to publish user notification count change to %v, %v, %v", userNotifCountChangedChannel, userNo, c.Err())
	}
	return nil
}
