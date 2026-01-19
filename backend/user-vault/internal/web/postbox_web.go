package web

import (
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/user-vault/api"
	"github.com/curtisnewbie/user-vault/internal/postbox"
	"github.com/curtisnewbie/user-vault/internal/repo"
	"github.com/spf13/cast"
)

const (
	ResourceQueryNotification  = "postbox:notification:query"
	ResourceCreateNotification = "postbox:notification:create"
)

func RegisterPostboxRoutes(rail miso.Rail) {

	miso.BaseRoute("/open/api/v1/notification").Group(

		miso.HttpPost("/create", miso.AutoHandler(ApiCreateNotification)).
			Desc("Create platform notification").
			Resource(ResourceCreateNotification),

		miso.HttpPost("/query", miso.AutoHandler(ApiQueryNotification)).
			Desc("Query platform notification").
			Resource(ResourceQueryNotification),

		miso.HttpGet("/count", miso.ResHandler(ApiCountNotification)).
			Desc("Count received platform notification").
			Resource(ResourceQueryNotification),

		miso.HttpPost("/open", miso.AutoHandler(ApiOpenNotification)).
			Desc("Record user opened platform notification").
			Resource(ResourceQueryNotification),

		miso.HttpPost("/open-all", miso.AutoHandler(ApiOpenAllNotification)).
			Desc("Mark all notifications opened").
			Resource(ResourceQueryNotification),
	)

	miso.BaseRoute("/open/api/v2/notification").Group(
		miso.HttpGet("/count", miso.RawHandler(ApiV2CountNotification)).
			Desc("Count received platform notification using long polling").
			DocQueryParam("curr", "Current count (used to implement long polling)").
			Resource(ResourceQueryNotification),
	)
}

func ApiCreateNotification(inb *miso.Inbound, req api.CreateNotificationReq) (any, error) {
	rail := inb.Rail()
	return nil, postbox.CreateNotification(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

func ApiQueryNotification(inb *miso.Inbound, req repo.QueryNotificationReq) (miso.PageRes[repo.ListedNotification], error) {
	rail := inb.Rail()
	return repo.QueryNotification(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

func ApiCountNotification(inb *miso.Inbound) (int, error) {
	rail := inb.Rail()
	return postbox.CachedCountNotification(rail, mysql.GetMySQL(), common.GetUser(rail))
}

func ApiV2CountNotification(inb *miso.Inbound) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	w, _ := inb.Unwrap()
	curr := cast.ToInt(inb.Query("curr"))
	if curr < 0 {
		curr = 0
	}
	postbox.Poll(rail, user, mysql.GetMySQL(), w, curr)
}

func ApiOpenNotification(inb *miso.Inbound, req repo.OpenNotificationReq) (any, error) {
	rail := inb.Rail()
	return nil, repo.OpenNotification(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

func ApiOpenAllNotification(inb *miso.Inbound, req repo.OpenNotificationReq) (any, error) {
	rail := inb.Rail()
	return nil, repo.OpenAllNotification(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}
