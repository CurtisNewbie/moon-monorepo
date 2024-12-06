package postbox

import (
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/user-vault/api"
	"github.com/spf13/cast"
)

const (
	ResourceQueryNotification  = "postbox:notification:query"
	ResourceCreateNotification = "postbox:notification:create"
)

func RegisterRoutes(rail miso.Rail) error {

	miso.BaseRoute("/open/api/v1/notification").Group(

		miso.IPost("/create", CreateNotificationEp).
			Desc("Create platform notification").
			Resource(ResourceCreateNotification),

		miso.IPost("/query", QueryNotificationEp).
			Desc("Query platform notification").
			Resource(ResourceQueryNotification),

		miso.Get("/count", CountNotificationEp).
			Desc("Count received platform notification").
			Resource(ResourceQueryNotification),

		miso.IPost("/open", OpenNotificationEp).
			Desc("Record user opened platform notification").
			Resource(ResourceQueryNotification),

		miso.IPost("/open-all", OpenAllNotificationEp).
			Desc("Mark all notifications opened").
			Resource(ResourceQueryNotification),
	)

	miso.BaseRoute("/open/api/v2/notification").Group(
		miso.RawGet("/count", CountNotificationV2Ep).
			Desc("Count received platform notification using long polling").
			DocQueryParam("curr", "Current count (used to implement long polling)").
			Resource(ResourceQueryNotification),
	)
	return nil
}

func CreateNotificationEp(inb *miso.Inbound, req api.CreateNotificationReq) (any, error) {
	rail := inb.Rail()
	return nil, CreateNotification(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

type QueryNotificationReq struct {
	Page   miso.Paging
	Status string
}

func QueryNotificationEp(inb *miso.Inbound, req QueryNotificationReq) (any, error) {
	rail := inb.Rail()
	return QueryNotification(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

func CountNotificationEp(inb *miso.Inbound) (any, error) {
	rail := inb.Rail()
	return CachedCountNotification(rail, mysql.GetMySQL(), common.GetUser(rail))
}

func CountNotificationV2Ep(inb *miso.Inbound) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	w, _ := inb.Unwrap()
	curr := cast.ToInt(inb.Query("curr"))
	if curr < 0 {
		curr = 0
	}
	longPollingHandler.Poll(rail, user, mysql.GetMySQL(), w, curr)
}

type OpenNotificationReq struct {
	NotifiNo string `valid:"notEmpty"`
}

func OpenNotificationEp(inb *miso.Inbound, req OpenNotificationReq) (any, error) {
	rail := inb.Rail()
	return nil, OpenNotification(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

func OpenAllNotificationEp(inb *miso.Inbound, req OpenNotificationReq) (any, error) {
	rail := inb.Rail()
	return nil, OpenAllNotification(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}
