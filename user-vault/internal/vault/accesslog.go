package vault

import (
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"gorm.io/gorm"
)

type AccessLog struct {
	Id         int
	UserAgent  string
	IpAddress  string
	UserId     int
	Username   string
	Url        string
	AccessTime util.ETime
	CreateTime util.ETime
	CreateBy   string
	UpdateTime util.ETime
	UpdateBy   string
	IsDel      bool
}

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
	return tx.Table("access_log").Create(&p).Error
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
	return miso.NewPageQuery[ListedAccessLog]().
		WithPage(req.Paging).
		WithSelectQuery(func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id", "access_time", "ip_address", "username", "url", "user_agent", "success").
				Order("id desc")
		}).
		WithBaseQuery(func(tx *gorm.DB) *gorm.DB {
			return tx.Table("access_log").
				Where("username = ?", user.Username)
		}).
		Exec(rail, tx)
}
