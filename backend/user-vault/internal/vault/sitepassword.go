package vault

import (
	"strings"

	"github.com/curtisnewbie/miso/middleware/crypto"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"gorm.io/gorm"
)

type ListSitePasswordReq struct {
	Alias    string
	Site     string
	Username string
	Paging   miso.Paging
}

type ListSitePasswordRes struct {
	RecordId   string
	Site       string
	Alias      string
	Username   string
	CreateTime util.ETime
}

func ListSitePasswords(rail miso.Rail, req ListSitePasswordReq, user common.User, db *gorm.DB) (miso.PageRes[ListSitePasswordRes], error) {
	return dbquery.NewPagedQuery[ListSitePasswordRes](db).
		WithBaseQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.Table("site_password").
				Eq("user_no", user.UserNo).
				LikeIf(req.Alias != "", "alias", req.Alias).
				LikeIf(req.Site != "", "site", req.Site).
				LikeIf(req.Username != "", "username", req.Username)
		}).
		WithSelectQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.Select("record_id,site,alias,username,create_time")
		}).
		Scan(rail, req.Paging)
}

type AddSitePasswordReq struct {
	Site          string
	Alias         string
	Username      string `valid:"notEmpty"`
	SitePassword  string `valid:"notEmpty"`
	LoginPassword string `valid:"notEmpty"`
}

func AddSitePassword(rail miso.Rail, req AddSitePasswordReq, user common.User, db *gorm.DB) error {
	u, err := loadUser(rail, db, user.Username)
	if err != nil {
		return err
	}

	if !checkPassword(u.Password, u.Salt, req.LoginPassword) {
		return miso.NewErrf("Login password incorrect, please try again")
	}

	encrypted, err := crypto.AesEcbEncrypt(pad256([]byte(req.LoginPassword)), req.SitePassword)
	if err != nil {
		rail.Warnf("Failed to encrypt site password, %v, %v", user.Username, err)
		return miso.ErrUnknownError
	}

	recordId := util.GenIdP("sitepw_")
	_, err = dbquery.NewQueryRail(rail, db).Exec(`
		INSERT INTO site_password (record_id, site, alias, username, password, user_no, create_by)
		values (?,?,?,?,?,?,?)
	`, recordId, req.Site, req.Alias, req.Username, encrypted, user.UserNo, user.Username)
	return err
}

type RemoveSitePasswordRes struct {
	RecordId string `valid:"notEmpty"`
}

func RemoveSitePassword(rail miso.Rail, req RemoveSitePasswordRes, user common.User, db *gorm.DB) error {
	_, err := loadBasicSitePassword(rail, db, user.UserNo, req.RecordId)
	if err != nil {
		return err
	}
	_, err = dbquery.NewQueryRail(rail, db).Exec("DELETE FROM site_password where record_id = ?", req.RecordId)
	return err
}

type DecryptSitePasswordReq struct {
	LoginPassword string `valid:"notEmpty"`
	RecordId      string `valid:"notEmpty"`
}

type DecryptSitePasswordRes struct {
	Decrypted string
}

func DecryptSitePassword(rail miso.Rail, req DecryptSitePasswordReq, user common.User, db *gorm.DB) (DecryptSitePasswordRes, error) {
	bsp, err := loadBasicSitePassword(rail, db, user.UserNo, req.RecordId)
	if err != nil {
		return DecryptSitePasswordRes{}, err
	}

	u, err := loadUser(rail, db, user.Username)
	if err != nil {
		return DecryptSitePasswordRes{}, err
	}
	if !checkPassword(u.Password, u.Salt, req.LoginPassword) {
		return DecryptSitePasswordRes{}, miso.NewErrf("Login password incorrect, please try again")
	}

	decrypted, err := crypto.AesEcbDecrypt(pad256([]byte(req.LoginPassword)), bsp.Password)
	if err != nil {
		rail.Warnf("Failed to encrypt site password, %v, %v", user.Username, err)
		return DecryptSitePasswordRes{}, miso.NewErrf("Password incorrect")
	}
	return DecryptSitePasswordRes{Decrypted: decrypted}, nil
}

type BasicSitePassword struct {
	RecordId string
	UserNo   string
	Password string
}

func loadBasicSitePassword(rail miso.Rail, db *gorm.DB, userNo string, recordId string) (BasicSitePassword, error) {
	var bsp BasicSitePassword
	n, err := dbquery.NewQueryRail(rail, db).
		From("site_password").
		Eq("record_id", recordId).
		Select("password, user_no").
		Scan(&bsp)
	if err != nil {
		return bsp, err
	}
	if n < 1 {
		return bsp, miso.NewErrf("Record not found")
	}
	bsp.RecordId = recordId
	if bsp.UserNo != userNo {
		return bsp, miso.ErrNotPermitted
	}
	return bsp, nil
}

func pad256(b []byte) []byte {
	if len(b) < 32 {
		cp := make([]byte, 32)
		copy(cp, b)
		b = cp
	}
	return b
}

type EditSitePasswordReq struct {
	RecordId      string
	Site          string
	Username      string
	Alias         string
	SitePassword  string `desc:"new site password, optional"`
	LoginPassword string `desc:"only used when site password is provided"`
}

func EditSitePassword(rail miso.Rail, req EditSitePasswordReq, user common.User, db *gorm.DB) error {
	_, err := loadBasicSitePassword(rail, db, user.UserNo, req.RecordId)
	if err != nil {
		return err
	}

	var encryptedSitePwd string
	if req.SitePassword != "" {
		u, err := loadUser(rail, db, user.Username)
		if err != nil {
			return err
		}

		if !checkPassword(u.Password, u.Salt, req.LoginPassword) {
			return miso.NewErrf("Login password incorrect, please try again")
		}

		encrypted, err := crypto.AesEcbEncrypt(pad256([]byte(req.LoginPassword)), req.SitePassword)
		if err != nil {
			rail.Warnf("Failed to encrypt site password, %v, %v", user.Username, err)
			return miso.ErrUnknownError.Wrap(err)
		}
		encryptedSitePwd = encrypted
	}

	_, err = dbquery.NewQueryRail(rail, db).From("site_password").
		Eq("record_id", req.RecordId).
		Set("site", req.Site).
		Set("alias", req.Alias).
		Set("update_by", user.Username).
		SetIf(strings.TrimSpace(req.Username) != "", "username", req.Username).
		SetIf(encryptedSitePwd != "", "password", encryptedSitePwd).
		Update()
	if err != nil {
		return miso.ErrUnknownError.Wrap(err)
	}
	return nil
}
