package vault

import (
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"gorm.io/gorm"
)

type SaveAccessLogParam struct {
	UserAgent  string
	IpAddress  string
	UserId     int
	Username   string
	Url        string
	Success    bool
	AccessTime util.ETime
}

func SaveAccessLogEvent(rail miso.Rail, tx *gorm.DB, p SaveAccessLogParam) error {
	_, err := dbquery.NewQuery(rail, tx).Table("access_log").Create(&p)
	return err
}

type ListedAccessLog struct {
	Id         int        `json:"id"`
	UserAgent  string     `json:"userAgent"`
	IpAddress  string     `json:"ipAddress"`
	Username   string     `json:"username"`
	Url        string     `json:"url"`
	AccessTime util.ETime `json:"accessTime"`
	Success    bool
}

type ListAccessLogReq struct {
	Paging miso.Paging `json:"paging"`
}

func ListAccessLogs(rail miso.Rail, tx *gorm.DB, user common.User, req ListAccessLogReq) (miso.PageRes[ListedAccessLog], error) {
	return dbquery.NewPagedQuery[ListedAccessLog](tx).
		WithSelectQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.SelectCols(ListedAccessLog{}).
				OrderDesc("id")
		}).
		WithBaseQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.Table("access_log").Eq("username", user.Username)
		}).
		Scan(rail, req.Paging)
}
