package postbox

import (
	"fmt"
	"time"

	"github.com/curtisnewbie/event-pump/client"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/slutil"
	"github.com/curtisnewbie/user-vault/api"
	"github.com/curtisnewbie/user-vault/internal/repo"
	"gorm.io/gorm"
)

var (
	userNotifCountCache = redis.NewRCache[int]("postbox:notification:count", redis.RCacheConfig{Exp: time.Minute * 30})
)

func CreateNotification(rail miso.Rail, db *gorm.DB, req api.CreateNotificationReq, user common.User) error {
	if len(req.ReceiverUserNos) < 1 {
		return nil
	}

	// check whether the userNos are leegal
	req.ReceiverUserNos = slutil.Distinct(req.ReceiverUserNos)

	for _, u := range req.ReceiverUserNos {
		sr := repo.SaveNotifiReq{
			UserNo:  u,
			Title:   req.Title,
			Message: req.Message,
		}
		if err := repo.SaveNotification(rail, db, sr, user); err != nil {
			return fmt.Errorf("failed to save notification, %+v, %v", sr, err)
		}
	}

	return nil
}

func CachedCountNotification(rail miso.Rail, db *gorm.DB, user common.User) (int, error) {
	v, err := userNotifCountCache.GetValElse(rail, user.UserNo, func() (int, error) {
		return repo.CountNotification(rail, db, user)
	})
	return v, err
}

func evictNotifCountCache(rail miso.Rail, t client.StreamEvent) error {
	userNo, ok := t.ColumnAfter("user_no")
	if !ok {
		return nil
	}
	rail.Infof("User notification changed, eventType: %v, %v", t.Type, userNo)
	if err := userNotifCountCache.Del(rail, userNo); err != nil {
		rail.Errorf("Failed to evict user notification count cache, %v, %v", userNo, err)
	}

	if c := redis.GetRedis().Publish(rail.Context(), userNotifCountChangedChannel, userNo); c.Err() != nil {
		rail.Errorf("Failed to publish user notification count change to %v, %v, %v", userNotifCountChangedChannel, userNo, c.Err())
	}
	return nil
}
