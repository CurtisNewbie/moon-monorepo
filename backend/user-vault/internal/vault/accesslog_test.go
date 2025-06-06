package vault

import (
	"testing"

	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/rabbit"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
)

func preAccessLogTest(t *testing.T) miso.Rail {
	rail := preTest(t)
	if e := mysql.InitMySQLFromProp(rail); e != nil {
		t.Fatal(e)
	}
	if e := rabbit.StartRabbitMqClient(rail); e != nil {
		t.Fatal(e)
	}
	return rail
}

func TestSendAccessLogEvent(t *testing.T) {
	rail := preAccessLogTest(t)
	er := AccessLogPipeline.Send(rail, AccessLogEvent{
		IpAddress:  "127.0.0.1",
		UserAgent:  "Linux Ubuntu",
		UserId:     0,
		Username:   "yongj.zhuang",
		Url:        "/user-vault/open/api/user/login",
		AccessTime: util.Now(),
		Success:    true,
	})

	if er != nil {
		t.Fatal(er)
	}
}
