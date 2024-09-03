package vfm

import (
	"time"

	fstore "github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/miso"
	vault "github.com/curtisnewbie/user-vault/api"
)

var (
	userIdInfoCache = redis.NewRCache[vault.UserInfo]("vfm:user:info:userno", redis.RCacheConfig{Exp: 5 * time.Minute, NoSync: true})
)

func CachedFindUser(rail miso.Rail, userNo string) (vault.UserInfo, error) {
	return userIdInfoCache.Get(rail, userNo, func() (vault.UserInfo, error) {
		fui, errFind := vault.FindUser(rail, vault.FindUserReq{
			UserNo: &userNo,
		})
		return fui, errFind
	})
}

func GetFstoreTmpToken(rail miso.Rail, fileId string, filename string) (string, error) {
	return fstore.GenTempFileKey(rail, fileId, filename)
}
