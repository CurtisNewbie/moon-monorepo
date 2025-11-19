package flow

import (
	"testing"
	"time"

	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/atom"
)

func TestListCashFlows(t *testing.T) {
	rail := miso.EmptyRail()
	if err := miso.LoadConfigFromFile("../../conf.yml", rail); err != nil {
		t.Fatal(err)
	}
	miso.SetLogLevel("debug")
	mysql.InitMySQLFromProp(rail)
	LoadCategoryConfs(rail)

	l, err := ListCashFlows(rail, mysql.GetMySQL(), common.User{UserNo: "test_user"}, ListCashFlowReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("1. l: %+v", l)

	l, err = ListCashFlows(rail, mysql.GetMySQL(), common.User{UserNo: "test_user"}, ListCashFlowReq{Direction: "OUT"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("2. l: %+v", l)

	l, err = ListCashFlows(rail, mysql.GetMySQL(), common.User{UserNo: "test_user"}, ListCashFlowReq{Direction: "IN"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("3. l: %+v", l)

	l, err = ListCashFlows(rail, mysql.GetMySQL(), common.User{UserNo: "test_user"}, ListCashFlowReq{TransId: "123"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("4. l: %+v", l)

	l, err = ListCashFlows(rail, mysql.GetMySQL(), common.User{UserNo: "test_user"}, ListCashFlowReq{TransId: "444"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("5. l: %+v", l)

	l, err = ListCashFlows(rail, mysql.GetMySQL(), common.User{UserNo: "test_user"}, ListCashFlowReq{Category: "WECHAT"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("6. l: %+v", l)

	tt := atom.Now().Add(-time.Hour * 24)
	l, err = ListCashFlows(rail, mysql.GetMySQL(), common.User{UserNo: "test_user"}, ListCashFlowReq{TransTimeStart: &tt})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("7. l: %+v", l)

	tt = atom.Now().Add(time.Hour * 24)
	l, err = ListCashFlows(rail, mysql.GetMySQL(), common.User{UserNo: "test_user"}, ListCashFlowReq{TransTimeStart: &tt})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("8. l: %+v", l)
}

func TestSaveCashflows(t *testing.T) {
	rail := miso.EmptyRail()
	if err := miso.LoadConfigFromFile("../../conf.yml", rail); err != nil {
		t.Fatal(err)
	}
	mysql.InitMySQLFromProp(rail)
	redis.InitRedisFromProp(rail)
	miso.SetLogLevel("debug")

	nc := []NewCashflow{
		{
			TransId:      "1234",
			Direction:    DirectionOut,
			TransTime:    atom.Now(),
			Counterparty: "Apple Corp",
			Amount:       "123.44",
			Currency:     "CNY",
			Extra:        "{}",
			Remark:       "Wechat Pay",
		},
	}
	p := SaveCashflowParams{
		Cashflows: nc,
		User:      common.User{UserNo: "UE1049787455160320075953"},
		Category:  WechatCategory,
	}
	_, err := SaveCashflows(rail, mysql.GetMySQL(), p)
	if err != nil {
		t.Fatal(err)
	}
}
