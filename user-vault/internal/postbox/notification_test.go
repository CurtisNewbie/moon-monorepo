package postbox

import (
	"testing"

	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/user-vault/api"
)

func _notificationPreTest(t *testing.T) miso.Rail {
	miso.SetProp(miso.PropMySQLEnabled, true)
	miso.SetProp(miso.PropMySQLUser, "root")
	miso.SetProp(miso.PropMySQLSchema, "postbox")
	miso.SetProp(miso.PropMySQLPassword, "")
	miso.SetProp("client.addr.user-vault.host", "localhost")
	miso.SetProp("client.addr.user-vault.port", "8089")
	rail := miso.EmptyRail()
	miso.InitMySQLFromProp(rail)
	return rail
}

func TestSaveNotification(t *testing.T) {
	rail := _notificationPreTest(t)
	user := common.User{
		Username: "postbox",
	}
	err := SaveNotification(rail, miso.GetMySQL(), SaveNotifiReq{
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
	err := CreateNotification(rail, miso.GetMySQL(), api.CreateNotificationReq{
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
	cnt, err := CountNotification(rail, miso.GetMySQL(), user)
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
	res, err := QueryNotification(rail, miso.GetMySQL(), QueryNotificationReq{}, user)
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
	err := OpenNotification(rail, miso.GetMySQL(), OpenNotificationReq{
		NotifiNo: "notif_1082375466205184155465",
	}, user)
	if err != nil {
		t.Fatal(err)
	}
}
