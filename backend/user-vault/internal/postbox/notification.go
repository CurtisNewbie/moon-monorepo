package postbox

import (
	"fmt"
	"time"

	"github.com/curtisnewbie/event-pump/client"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"github.com/curtisnewbie/miso/util/errs"
	"github.com/curtisnewbie/miso/util/slutil"
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
	req.ReceiverUserNos = slutil.Distinct(req.ReceiverUserNos)

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
	err := dbquery.NewQuery(rail, db).
		NotLogSQL().
		Table("notification").
		CreateAny(struct {
			UserNo   string
			NotifiNo string
			Title    string
			Message  string
		}{
			UserNo:   req.UserNo,
			NotifiNo: notifiNo,
			Title:    req.Title,
			Message:  req.Message,
		})
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
	CreateTime util.ETime `gorm:"column:created_at"`
}

type QueryNotificationReq struct {
	Page   miso.Paging
	Status string
}

func QueryNotification(rail miso.Rail, db *gorm.DB, req QueryNotificationReq, user common.User) (miso.PageRes[ListedNotification], error) {
	return dbquery.NewPagedQuery[ListedNotification](db).
		WithBaseQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.Table("notification").
				Eq("user_no", user.UserNo).
				EqNotEmpty("status", req.Status)
		}).
		WithSelectQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.SelectCols(ListedNotification{}).
				OrderDesc("id").
				Limit(req.Page.GetLimit()).
				Offset(req.Page.GetOffset())
		}).
		Scan(rail, req.Page)
}

func CachedCountNotification(rail miso.Rail, db *gorm.DB, user common.User) (int, error) {
	v, err := userNotifCountCache.GetValElse(rail, user.UserNo, func() (int, error) {
		return CountNotification(rail, db, user)
	})
	return v, err
}

func CountNotification(rail miso.Rail, db *gorm.DB, user common.User) (int, error) {
	count, err := dbquery.NewQuery(rail, db).
		Table("notification").
		Eq("user_no", user.UserNo).
		Eq("status", StatusInit).
		Count()
	return int(count), err
}

type OpenNotificationReq struct {
	NotifiNo string `valid:"notEmpty"`
}

func OpenNotification(rail miso.Rail, db *gorm.DB, req OpenNotificationReq, user common.User) error {
	err := dbquery.NewQuery(rail, db).
		Table("notification").
		Set("status", StatusOpened).
		Eq("notifi_no", req.NotifiNo).
		Eq("user_no", user.UserNo).
		UpdateAny()
	return err
}

func OpenAllNotification(rail miso.Rail, db *gorm.DB, req OpenNotificationReq, user common.User) error {
	var id int
	n, err := dbquery.NewQuery(rail, db).
		Table("notification").
		Select("id").
		Eq("user_no", user.UserNo).
		Eq("notifi_no", req.NotifiNo).
		Scan(&id)
	if err != nil {
		return err
	}
	if n < 1 {
		return errs.NewErrf("Record not found")
	}

	err = dbquery.NewQuery(rail, db).
		Table("notification").
		Set("status", StatusOpened).
		Eq("user_no", user.UserNo).
		Eq("status", StatusInit).
		Le("id", id).
		UpdateAny()
	return err
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

	if c := redis.GetRedis().Publish(rail.Context(), userNotifCountChangedChannel, userNo); c.Err() != nil {
		rail.Errorf("Failed to publish user notification count change to %v, %v, %v", userNotifCountChangedChannel, userNo, c.Err())
	}
	return nil
}
