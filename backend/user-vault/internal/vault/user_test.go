package vault

import (
	"testing"
	"time"

	"github.com/curtisnewbie/miso/middleware/jwt"
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/user-vault/api"
	"github.com/sirupsen/logrus"
)

func preTest(t *testing.T) miso.Rail {
	logrus.SetLevel(logrus.DebugLevel)
	rail := miso.EmptyRail()
	miso.DefaultReadConfig([]string{"configFile=../../conf.yml"}, rail)
	return rail
}

func preUserTest(t *testing.T) miso.Rail {
	rail := preTest(t)
	if mysql.InitMySQLFromProp(rail) != nil {
		t.FailNow()
	}
	if e := miso.InitConsulClient(); e != nil {
		t.Log(e)
		t.FailNow()
	}
	if _, e := redis.InitRedisFromProp(rail); e != nil {
		t.Log(e)
		t.FailNow()
	}
	return rail
}

func TestExtractSpringSalt(t *testing.T) {
	salt := extractSpringSalt("{asdfasdfasdf}sadkfasdfasdf")
	if salt != "{asdfasdfasdf}" {
		t.Log(salt)
		t.FailNow()
	}
}

func TestCheckPassword(t *testing.T) {
	pw := ""
	ok := checkPassword("d7030adc17d5623265162432398c9d25dd14fd8cf3ddc9504b149e590cbacd73", "30689", pw)
	if !ok {
		t.FailNow()
	}
}

func TestUserLogin(t *testing.T) {
	rail := preUserTest(t)
	uname := "banana"
	pword := "12345678"

	usr, err := userLogin(rail, mysql.GetMySQL(), uname, pword)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Logf("user: %+v", usr)

	tkn, err := buildToken(TokenUser{
		Id:       usr.Id,
		UserNo:   usr.UserNo,
		Username: usr.Username,
		RoleNo:   usr.RoleNo,
	}, time.Minute*15)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	t.Logf("tkn: %+v", tkn)

	decoded, err := jwt.JwtDecode(tkn)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Logf("decoded: %+v", decoded)

	tu, err := DecodeTokenUser(rail, tkn)
	if err != nil {
		t.Fatal(err)
	}
	rail.Infof("tokenUser: %+v", tu)
}

func TestAdminAddUser(t *testing.T) {
	rail := preUserTest(t)
	e := NewUser(rail, mysql.GetMySQL(), CreateUserParam{
		Username:     "dummydummy2",
		Password:     "12345678",
		RoleNo:       "role_628043111874560208429",
		ReviewStatus: api.ReviewApproved,
	})
	if e != nil {
		t.Fatal(e)
	}
}

func TestListUsers(t *testing.T) {
	rail := preUserTest(t)
	users, err := ListUsers(rail, mysql.GetMySQL(), ListUserReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", users)
}

func TestAdminUpdateUser(t *testing.T) {
	rail := preUserTest(t)
	err := AdminUpdateUser(rail, mysql.GetMySQL(), AdminUpdateUserReq{
		UserNo:     "",
		RoleNo:     "role_628043111874560208429",
		IsDisabled: 0,
	}, common.NilUser())
	if err != nil {
		t.Fatal(err)
	}
}

func TestReviewUserRegistration(t *testing.T) {
	rail := preUserTest(t)
	err := ReviewUserRegistration(rail, mysql.GetMySQL(), AdminReviewUserReq{
		UserId:       1107,
		ReviewStatus: api.ReviewApproved,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadUserInfoBrief(t *testing.T) {
	rail := preUserTest(t)
	uib, err := LoadUserInfoBrief(rail, mysql.GetMySQL(), "dummydummy2")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("uib: %+v", uib)
}

func TestFindUserWithRes(t *testing.T) {
	rail := preUserTest(t)
	u, err := FindUserWithRes(rail, mysql.GetMySQL(), api.FetchUserWithResourceReq{
		ResourceCode: "manage-resources",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("u: %+v", u)
}
