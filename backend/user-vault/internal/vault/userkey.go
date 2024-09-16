package vault

import (
	"time"

	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
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
	ExpirationTime util.ETime
	UserId         int
	UserNo         string
}

func GenUserKey(rail miso.Rail, tx *gorm.DB, req GenUserKeyReq, username string) error {

	user, err := loadUser(rail, tx, username)
	if err != nil {
		return err
	}

	if !checkPassword(user.Password, user.Salt, req.Password) {
		return miso.NewErrf("Password incorrect, unable to generate user secret key")
	}

	key := util.RandStr(userKeyLen)
	return tx.Table("user_key").
		Create(NewUserKey{
			Name:           req.KeyName,
			SecretKey:      key,
			ExpirationTime: util.Now().Add(userKeyExpDur),
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
	Id             int        `json:"id"`
	SecretKey      string     `json:"secretKey"`
	Name           string     `json:"name"`
	ExpirationTime util.ETime `json:"expirationTime"`
	CreateTime     util.ETime `json:"createTime"`
}

func ListUserKeys(rail miso.Rail, tx *gorm.DB, req ListUserKeysReq, user common.User) (miso.PageRes[ListedUserKey], error) {
	return mysql.NewPageQuery[ListedUserKey]().
		WithPage(req.Paging).
		WithBaseQuery(func(tx *gorm.DB) *gorm.DB {
			tx = tx.Table("user_key").
				Where("user_no = ?", user.UserNo).
				Where("expiration_time > ?", util.Now()).
				Where("is_del = 0")

			if !util.IsBlankStr(req.Name) {
				tx = tx.Where("name LIKE ?", "%"+req.Name+"%")
			}
			return tx
		}).
		WithSelectQuery(func(tx *gorm.DB) *gorm.DB {
			return tx.Select("id, secret_key, name, expiration_time, create_time").
				Order("id DESC")
		}).
		Exec(rail, tx)
}

type DeleteUserKeyReq struct {
	UserKeyId int `json:"userKeyId"`
}

func DeleteUserKey(rail miso.Rail, tx *gorm.DB, req DeleteUserKeyReq, userNo string) error {
	return tx.Exec(`UPDATE user_key SET is_del = 1 WHERE user_no = ? AND id = ? AND is_del = 0`, userNo, req.UserKeyId).Error
}
