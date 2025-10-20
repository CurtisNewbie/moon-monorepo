package web

import (
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/user-vault/api"
	"github.com/curtisnewbie/user-vault/internal/postbox"
	"github.com/spf13/cast"
)

const (
	ResourceQueryNotification  = "postbox:notification:query"
	ResourceCreateNotification = "postbox:notification:create"
)

func RegisterPostboxRoutes(rail miso.Rail) {

	miso.BaseRoute("/open/api/v1/notification").Group(

		miso.HttpPost("/create", miso.AutoHandler(CreateNotificationEp)).
			Desc("Create platform notification").
			Resource(ResourceCreateNotification),

		miso.HttpPost("/query", miso.AutoHandler(QueryNotificationEp)).
			Desc("Query platform notification").
			Resource(ResourceQueryNotification),

		miso.HttpGet("/count", miso.ResHandler(CountNotificationEp)).
			Desc("Count received platform notification").
			Resource(ResourceQueryNotification),

		miso.HttpPost("/open", miso.AutoHandler(OpenNotificationEp)).
			Desc("Record user opened platform notification").
			Resource(ResourceQueryNotification),

		miso.HttpPost("/open-all", miso.AutoHandler(OpenAllNotificationEp)).
			Desc("Mark all notifications opened").
			Resource(ResourceQueryNotification),
	)

	miso.BaseRoute("/open/api/v2/notification").Group(
		miso.HttpGet("/count", miso.RawHandler(CountNotificationV2Ep)).
			Desc("Count received platform notification using long polling").
			DocQueryParam("curr", "Current count (used to implement long polling)").
			Resource(ResourceQueryNotification),
	)
}

func CreateNotificationEp(inb *miso.Inbound, req api.CreateNotificationReq) (any, error) {
	rail := inb.Rail()
	return nil, postbox.CreateNotification(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

func QueryNotificationEp(inb *miso.Inbound, req postbox.QueryNotificationReq) (any, error) {
	rail := inb.Rail()
	return postbox.QueryNotification(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

func CountNotificationEp(inb *miso.Inbound) (any, error) {
	rail := inb.Rail()
	return postbox.CachedCountNotification(rail, mysql.GetMySQL(), common.GetUser(rail))
}

func CountNotificationV2Ep(inb *miso.Inbound) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	w, _ := inb.Unwrap()
	curr := cast.ToInt(inb.Query("curr"))
	if curr < 0 {
		curr = 0
	}
	postbox.Poll(rail, user, mysql.GetMySQL(), w, curr)
}

func OpenNotificationEp(inb *miso.Inbound, req postbox.OpenNotificationReq) (any, error) {
	rail := inb.Rail()
	return nil, postbox.OpenNotification(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

func OpenAllNotificationEp(inb *miso.Inbound, req postbox.OpenNotificationReq) (any, error) {
	rail := inb.Rail()
	return nil, postbox.OpenAllNotification(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}
