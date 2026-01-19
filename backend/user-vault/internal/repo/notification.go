package repo

import (
	"fmt"

	"github.com/curtisnewbie/miso/errs"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/atom"
	"github.com/curtisnewbie/miso/util/snowflake"
	"gorm.io/gorm"
)

const (
	TableNotification = "notification"
)

const (
	StatusInit   = "INIT"
	StatusOpened = "OPENED"
)

type SaveNotifiReq struct {
	UserNo  string
	Title   string
	Message string
}

func SaveNotification(rail miso.Rail, db *gorm.DB, req SaveNotifiReq, user common.User) error {
	notifiNo := NotifiNo()
	err := dbquery.NewQuery(rail, db).
		NotLogSQL().
		Table(TableNotification).
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
	return snowflake.IdPrefix("notif_")
}

type ListedNotification struct {
	Id         int       `json:"id"`
	NotifiNo   string    `json:"notifiNo"`
	Title      string    `json:"title"`
	Message    string    `json:"message"`
	Status     string    `json:"status"`
	CreateTime atom.Time `gorm:"column:created_at" json:"createTime"`
}

type QueryNotificationReq struct {
	Page   miso.Paging `json:"page"`
	Status string      `json:"status"`
}

func QueryNotification(rail miso.Rail, db *gorm.DB, req QueryNotificationReq, user common.User) (miso.PageRes[ListedNotification], error) {
	return dbquery.NewPagedQuery[ListedNotification](db).
		WithBaseQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.Table(TableNotification).
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

func CountNotification(rail miso.Rail, db *gorm.DB, user common.User) (int, error) {
	count, err := dbquery.NewQuery(rail, db).
		Table(TableNotification).
		Eq("user_no", user.UserNo).
		Eq("status", StatusInit).
		Count()
	return int(count), err
}

type OpenNotificationReq struct {
	NotifiNo string `valid:"notEmpty" json:"notifiNo"`
}

func OpenNotification(rail miso.Rail, db *gorm.DB, req OpenNotificationReq, user common.User) error {
	err := dbquery.NewQuery(rail, db).
		Table(TableNotification).
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
