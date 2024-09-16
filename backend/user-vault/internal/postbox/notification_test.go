package postbox

import (
	"testing"

	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/user-vault/api"
)

func _notificationPreTest(t *testing.T) miso.Rail {
	miso.SetProp(mysql.PropMySQLEnabled, true)
	miso.SetProp(mysql.PropMySQLUser, "root")
	miso.SetProp(mysql.PropMySQLSchema, "postbox")
	miso.SetProp(mysql.PropMySQLPassword, "")
	miso.SetProp("client.addr.user-vault.host", "localhost")
	miso.SetProp("client.addr.user-vault.port", "8089")
	rail := miso.EmptyRail()
	mysql.InitMySQLFromProp(rail)
	return rail
}

func TestSaveNotification(t *testing.T) {
	rail := _notificationPreTest(t)
	user := common.User{
		Username: "postbox",
	}
	err := SaveNotification(rail, mysql.GetMySQL(), SaveNotifiReq{
		UserNo:  "UE1049787455160320075953",
		Title:   "Some message",
		Message: "Notification should be saved",
	}, user)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateNotification(t *testing.T) {
	miso.SetLogLevel("debug")
	rail := _notificationPreTest(t)
	user := common.User{
		Username: "postbox",
	}
	err := CreateNotification(rail, mysql.GetMySQL(), api.CreateNotificationReq{
		ReceiverUserNos: []string{"UE1049787455160320075953"},
		Title:           "Some message",
		Message:         "Notification should be saved",
	}, user)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCountNotification(t *testing.T) {
	miso.SetLogLevel("debug")
	rail := _notificationPreTest(t)
	user := common.User{
		UserNo:   "UE1049787455160320075953",
		Username: "postbox",
	}
	cnt, err := CountNotification(rail, mysql.GetMySQL(), user)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cnt)
}

func TestQueryNotification(t *testing.T) {
	miso.SetLogLevel("debug")
	rail := _notificationPreTest(t)
	user := common.User{
		UserNo:   "UE1049787455160320075953",
		Username: "postbox",
	}
	res, err := QueryNotification(rail, mysql.GetMySQL(), QueryNotificationReq{}, user)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", res)
}

func TestOpenNotification(t *testing.T) {
	miso.SetLogLevel("debug")
	rail := _notificationPreTest(t)
	user := common.User{
		UserNo:   "UE1049787455160320075953",
		Username: "admin",
	}
	err := OpenNotification(rail, mysql.GetMySQL(), OpenNotificationReq{
		NotifiNo: "notif_1082375466205184155465",
	}, user)
	if err != nil {
		t.Fatal(err)
	}
}
