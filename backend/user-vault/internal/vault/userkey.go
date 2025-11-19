package vault

import (
	"time"

	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/atom"
	"github.com/curtisnewbie/miso/util/errs"
	"github.com/curtisnewbie/miso/util/randutil"
	"github.com/curtisnewbie/miso/util/strutil"
	"gorm.io/gorm"
)

var (
	userKeyExpDur time.Duration = 90 * 24 * time.Hour
	userKeyLen                  = 64
)

type GenUserKeyReq struct {
	Password string `json:"password" valid:"notEmpty"`
	KeyName  string `json:"keyName" valid:"notEmpty"`
}

type NewUserKey struct {
	Name           string
	SecretKey      string
	ExpirationTime atom.Time
	UserId         int
	UserNo         string
}

func GenUserKey(rail miso.Rail, tx *gorm.DB, req GenUserKeyReq, username string) error {

	user, err := loadUser(rail, tx, username)
	if err != nil {
		return err
	}

	if !checkPassword(user.Password, user.Salt, req.Password) {
		return errs.NewErrf("Password incorrect, unable to generate user secret key")
	}

	key := randutil.RandStr(userKeyLen)
	return tx.Table("user_key").
		Create(NewUserKey{
			Name:           req.KeyName,
			SecretKey:      key,
			ExpirationTime: atom.Now().Add(userKeyExpDur),
			UserId:         user.Id,
			UserNo:         user.UserNo,
		}).
		Error
}

type ListUserKeysReq struct {
	Paging miso.Paging `json:"paging"`
	Name   string      `json:"name"`
}

type ListedUserKey struct {
	Id             int       `json:"id"`
	SecretKey      string    `json:"secretKey"`
	Name           string    `json:"name"`
	ExpirationTime atom.Time `json:"expirationTime"`
	CreateTime     atom.Time `json:"createTime" gorm:"column:created_at"`
}

func ListUserKeys(rail miso.Rail, tx *gorm.DB, req ListUserKeysReq, user common.User) (miso.PageRes[ListedUserKey], error) {
	return dbquery.NewPagedQuery[ListedUserKey](tx).
		WithBaseQuery(func(q *dbquery.Query) *dbquery.Query {
			q = q.Table("user_key").
				Where("user_no = ?", user.UserNo).
				Where("expiration_time > ?", atom.Now()).
				Where("deleted = 0")
			return q.LikeIf(!strutil.IsBlankStr(req.Name), "name", req.Name)
		}).
		WithSelectQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.SelectCols(ListedUserKey{}).
				Order("id DESC")
		}).
		Scan(rail, req.Paging)
}

type DeleteUserKeyReq struct {
	UserKeyId int `json:"userKeyId"`
}

func DeleteUserKey(rail miso.Rail, tx *gorm.DB, req DeleteUserKeyReq, userNo string) error {
	err := dbquery.NewQuery(rail, tx).
		Table("user_key").
		Set("deleted", true).
		Eq("user_no", userNo).
		Eq("id", req.UserKeyId).
		Eq("deleted", false).
		UpdateAny()
	return err
}
