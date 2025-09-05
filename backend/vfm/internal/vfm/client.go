package vfm

import (
	"time"

	fstore "github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	vault "github.com/curtisnewbie/user-vault/api"
)

var (
	userIdInfoCache = redis.NewRCache[vault.UserInfo]("vfm:user:info:userno", redis.RCacheConfig{Exp: 5 * time.Minute, NoSync: true})
	fstorePool      = util.NewIOAsyncPool()
)

func CachedFindUser(rail miso.Rail, userNo string) (vault.UserInfo, error) {
	return userIdInfoCache.GetValElse(rail, userNo, func() (vault.UserInfo, error) {
		return vault.FindUser(rail, vault.FindUserReq{
			UserNo: &userNo,
		})
	})
}

func GetFstoreTmpToken(rail miso.Rail, fileId string, filename string) (string, error) {
	return fstore.GenTempFileKey(rail, fileId, filename)
}

func GetFstoreTmpTokenAsync(rail miso.Rail, fileId string, filename string) util.Future[string] {
	return util.SubmitAsync(fstorePool, func() (string, error) {
		return fstore.GenTempFileKey(rail, fileId, filename)
	})
}

type FstoreTmpTokenReq struct {
	FileId   string
	Filename string
}

func BatchGetFstoreTmpToken(rail miso.Rail, reqs []FstoreTmpTokenReq) map[string]string {
	if len(reqs) < 1 {
		return map[string]string{}
	}
	futures := make(map[string]util.Future[string], len(reqs))
	for _, r := range reqs {
		futures[r.FileId] = GetFstoreTmpTokenAsync(rail, r.FileId, r.Filename)
	}

	res := make(map[string]string, len(reqs))
	for k, v := range futures {
		r, err := v.Get()
		if err != nil {
			rail.Infof("Failed to get fstore tmp token for fileId: %v, %v", k, err)
		}
		res[k] = r
	}
	return res
}
