package web

import (
	"fmt"
	"time"

	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/json"
	"github.com/curtisnewbie/miso/util/randutil"
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

		miso.HttpPost("/ws-ticket", miso.ResHandler(ApiWsTicket)).
			Desc("Generate a one-time websocket ticket for notification count push").
			Resource(ResourceQueryNotification),
	)

	miso.BaseRoute("/open/api/v2/notification").Group(
		miso.HttpGet("/count", miso.RawHandler(ApiV2CountNotification)).
			Desc("Count received platform notification using long polling").
			DocQueryParam("curr", "Current count (used to implement long polling)").
			Resource(ResourceQueryNotification),

		miso.HttpGet("/ws", miso.RawHandler(ApiWsNotification)).
			Desc("WebSocket endpoint for notification count push").
			Resource(ResourceQueryNotification),
	)

	miso.HttpPost("/internal/v1/notification/ws-exchange", miso.AutoHandler[api.WsExchangeReq, flow.User](ApiWsExchange)).
		Desc("Exchange a websocket ticket for user info")
}

func ApiCreateNotification(inb *miso.Inbound, req api.CreateNotificationReq) (any, error) {
	rail := inb.Rail()
	return nil, postbox.CreateNotification(rail, mysql.GetMySQL(), req, flow.GetUser(rail))
}

func ApiQueryNotification(inb *miso.Inbound, req repo.QueryNotificationReq) (miso.PageRes[repo.ListedNotification], error) {
	rail := inb.Rail()
	return repo.QueryNotification(rail, mysql.GetMySQL(), req, flow.GetUser(rail))
}

func ApiCountNotification(inb *miso.Inbound) (int, error) {
	rail := inb.Rail()
	return postbox.CachedCountNotification(rail, mysql.GetMySQL(), flow.GetUser(rail))
}

func ApiV2CountNotification(inb *miso.Inbound) {
	rail := inb.Rail()
	user := flow.GetUser(rail)
	w, _ := inb.Unwrap()
	curr := cast.ToInt(inb.Query("curr"))
	if curr < 0 {
		curr = 0
	}
	postbox.Poll(rail, user, mysql.GetMySQL(), w, curr)
}

func ApiOpenNotification(inb *miso.Inbound, req repo.OpenNotificationReq) (any, error) {
	rail := inb.Rail()
	return nil, repo.OpenNotification(rail, mysql.GetMySQL(), req, flow.GetUser(rail))
}

func ApiOpenAllNotification(inb *miso.Inbound, req repo.OpenNotificationReq) (any, error) {
	rail := inb.Rail()
	return nil, repo.OpenAllNotification(rail, mysql.GetMySQL(), req, flow.GetUser(rail))
}

func ApiWsNotification(inb *miso.Inbound) {
	rail := inb.Rail()
	user := flow.GetUser(rail)
	w, r := inb.Unwrap()
	postbox.HandleWS(rail, user, mysql.GetMySQL(), w, r)
}

func ApiWsTicket(inb *miso.Inbound) (api.WsTicketResp, error) {
	rail := inb.Rail()
	user := flow.GetUser(rail)

	ticket := randutil.ERand(32)
	key := "user-vault:ws:ticket:" + ticket

	data, err := json.WriteJson(map[string]interface{}{
		"username": user.Username,
		"userNo":   user.UserNo,
		"roleNo":   user.RoleNo,
	})
	if err != nil {
		return api.WsTicketResp{}, fmt.Errorf("failed to encode ticket data: %w", err)
	}

	r := redis.GetRedis()
	if err := r.Set(rail.Context(), key, data, 1*time.Minute).Err(); err != nil {
		return api.WsTicketResp{}, fmt.Errorf("failed to store ticket: %w", err)
	}

	rail.Infof("WS ticket generated for user %v", user.UserNo)
	return api.WsTicketResp{Ticket: ticket}, nil
}

func ApiWsExchange(inb *miso.Inbound, req api.WsExchangeReq) (flow.User, error) {
	rail := inb.Rail()
	key := "user-vault:ws:ticket:" + req.Ticket

	r := redis.GetRedis()
	val, err := r.GetDel(rail.Context(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return flow.User{}, miso.NewErrf("Invalid or expired ticket")
		}
		return flow.User{}, fmt.Errorf("failed to lookup ticket: %w", err)
	}

	var user flow.User
	if err := json.ParseJson([]byte(val), &user); err != nil {
		return flow.User{}, fmt.Errorf("failed to decode ticket data: %w", err)
	}

	rail.Infof("WS ticket exchanged for user %v", user.UserNo)
	return user, nil
}
