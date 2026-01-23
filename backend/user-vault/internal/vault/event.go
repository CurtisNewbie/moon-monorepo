package vault

import (
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/rabbit"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/atom"
	"github.com/curtisnewbie/user-vault/internal/repo"
)

var (
	AccessLogPipeline = rabbit.NewEventPipeline[AccessLogEvent]("event.bus.user-vault.access.log").
		LogPayload().
		MaxRetry(2).
		Listen(2, func(rail miso.Rail, evt AccessLogEvent) error {
			return repo.SaveAccessLogEvent(rail, mysql.GetMySQL(), repo.SaveAccessLogParam(evt))
		})
)

type AccessLogEvent struct {
	UserAgent  string
	IpAddress  string
	UserId     int
	Username   string
	Url        string
	Success    bool
	AccessTime atom.Time
}
