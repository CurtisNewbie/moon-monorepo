package vault

import (
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/rabbit"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
)

var (
	AccessLogPipeline = rabbit.NewEventPipeline[AccessLogEvent]("event.bus.user-vault.access.log").
		LogPayload().
		MaxRetry(2).
		Listen(2, func(rail miso.Rail, evt AccessLogEvent) error {
			return SaveAccessLogEvent(rail, mysql.GetMySQL(), SaveAccessLogParam(evt))
		})
)

type AccessLogEvent struct {
	UserAgent  string
	IpAddress  string
	UserId     int
	Username   string
	Url        string
	Success    bool
	AccessTime util.ETime
}
