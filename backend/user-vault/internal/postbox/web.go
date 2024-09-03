package postbox

import (
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/user-vault/api"
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

		miso.Post("/open-all", OpenAllNotificationEp).
			Desc("Mark all notifications opened").
			Resource(ResourceQueryNotification),
	)
	return nil
}

func CreateNotificationEp(inb *miso.Inbound, req api.CreateNotificationReq) (any, error) {
	rail := inb.Rail()
	return nil, CreateNotification(rail, miso.GetMySQL(), req, common.GetUser(rail))
}

type QueryNotificationReq struct {
	Page   miso.Paging
	Status string
}

func QueryNotificationEp(inb *miso.Inbound, req QueryNotificationReq) (any, error) {
	rail := inb.Rail()
	return QueryNotification(rail, miso.GetMySQL(), req, common.GetUser(rail))
}

func CountNotificationEp(inb *miso.Inbound) (any, error) {
	rail := inb.Rail()
	return CountNotification(rail, miso.GetMySQL(), common.GetUser(rail))
}

type OpenNotificationReq struct {
	NotifiNo string `valid:"notEmpty"`
}

func OpenNotificationEp(inb *miso.Inbound, req OpenNotificationReq) (any, error) {
	rail := inb.Rail()
	return nil, OpenNotification(rail, miso.GetMySQL(), req, common.GetUser(rail))
}

func OpenAllNotificationEp(inb *miso.Inbound) (any, error) {
	rail := inb.Rail()
	return nil, OpenAllNotification(rail, miso.GetMySQL(), common.GetUser(rail))
}
