package postbox

import (
	"time"

	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/user-vault/api"
	"github.com/curtisnewbie/user-vault/internal/vault"
)

var (
	UserWithResCache = miso.NewRCache[[]api.UserInfo]("user-vault:users-with-res:cache", miso.RCacheConfig{
		Exp: time.Second * 30,
	})
)

func InitPipeline(rail miso.Rail) error {

	api.CreateNotifiPipeline.Listen(3, func(rail miso.Rail, evt api.CreateNotifiEvent) error {
		if err := miso.Validate(evt); err != nil {
			rail.Errorf("Invalid event, %#v, %v", evt, err)
			return nil
		}
		return CreateNotification(rail, miso.GetMySQL(), api.CreateNotificationReq(evt), common.NilUser())
	})

	api.CreateNotifiByAccessPipeline.Listen(2, func(rail miso.Rail, evt api.CreateNotifiByAccessEvent) error {
		if err := miso.Validate(evt); err != nil {
			rail.Errorf("Invalid event, %#v, %v", evt, err)
			return nil
		}

		users, err := UserWithResCache.Get(rail, evt.ResCode, func() ([]api.UserInfo, error) {
			users, err := vault.FindUserWithRes(rail, miso.GetMySQL(), api.FetchUserWithResourceReq{ResourceCode: evt.ResCode})
			if err != nil {
				rail.Errorf("failed to FindUserWithRes, %v", err)
				return nil, err
			}
			return users, err
		})
		if err != nil {
			return err
		}

		if len(users) < 1 {
			return nil
		}

		un := make([]string, 0, len(users))
		for _, u := range users {
			un = append(un, u.UserNo)
		}

		return CreateNotification(rail, miso.GetMySQL(), api.CreateNotificationReq{
			Title:           evt.Title,
			Message:         evt.Message,
			ReceiverUserNos: un,
		}, common.NilUser())
	})
	return nil
}
