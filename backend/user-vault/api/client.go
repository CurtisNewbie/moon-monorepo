package api

import (
	"fmt"

	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
)

const (
	ServiceName = "user-vault"
)

func FindUser(rail miso.Rail, req FindUserReq) (UserInfo, error) {
	var r miso.GnResp[UserInfo]
	err := miso.NewDynClient(rail, "/remote/user/info", ServiceName).
		PostJson(req).
		Json(&r)
	if err != nil {
		return UserInfo{}, fmt.Errorf("failed to find user (user-vault), %v", err)
	}
	return r.Res()
}

func FindUserId(rail miso.Rail, username string) (int, error) {
	var r miso.GnResp[int]
	err := miso.NewDynClient(rail, "/remote/user/id", ServiceName).
		AddQueryParams("username", username).
		Get().
		Json(&r)
	if err != nil {
		return 0, fmt.Errorf("failed to findUserId (user-vault), %v", err)
	}
	return r.Res()
}

func FetchUsernames(rail miso.Rail, req FetchNameByUserNoReq) (FetchUsernamesRes, error) {
	var r miso.GnResp[FetchUsernamesRes]
	err := miso.NewDynClient(rail, "/remote/user/userno/username", ServiceName).
		PostJson(req).
		Json(&r)
	if err != nil {
		return FetchUsernamesRes{}, fmt.Errorf("failed to fetch usernames (user-vault), %v", err)
	}
	return r.Res()
}

func FetchUsersWithRole(rail miso.Rail, req FetchUsersWithRoleReq) ([]UserInfo, error) {
	var r miso.GnResp[[]UserInfo]
	err := miso.NewDynClient(rail, "/remote/user/list/with-role", ServiceName).
		PostJson(req).
		Json(&r)
	if err != nil {
		return nil, fmt.Errorf("failed to ListUsersWithRole, %w", err)
	}
	return r.Res()
}

func GetRoleInfo(rail miso.Rail, roleNo string) (RoleInfoResp, error) {
	var r miso.GnResp[RoleInfoResp]
	err := miso.NewDynClient(rail, "/open/api/role/info", ServiceName).
		PostJson(RoleInfoReq{RoleNo: roleNo}).
		Json(&r)
	if err != nil {
		return RoleInfoResp{}, fmt.Errorf("failed to GetRoleInfo, %w", err)
	}
	return r.Res()
}

func FetchUsersWithResource(rail miso.Rail, req FetchUserWithResourceReq) ([]UserInfo, error) {
	var r miso.GnResp[[]UserInfo]
	err := miso.NewDynClient(rail, "/remote/user/list/with-resource", ServiceName).
		PostJson(req).
		Json(&r)
	if err != nil {
		return nil, fmt.Errorf("failed to FetchUsersWithResource, %w", err)
	}
	return r.Res()
}

func CreateNotification(rail miso.Rail, req CreateNotificationReq) error {
	var resp miso.GnResp[any]
	err := miso.NewDynClient(rail, "/open/api/v1/notification/create", ServiceName).
		Require2xx().
		PostJson(req).
		Json(&resp)
	if err != nil {
		return fmt.Errorf("failed to create notiication, req: %+v, %v", req, err)
	}
	return resp.Err()
}

func SendCheckResAccessReq(rail miso.Rail, req CheckResAccessReq) (CheckResAccessResp, error) {
	var res miso.GnResp[CheckResAccessResp]
	err := miso.NewDynClient(rail, "/remote/path/resource/access-test", "user-vault").
		PostJson(req).
		Json(&res)
	if err != nil {
		rail.Errorf("Request failed, %v", err)
		var dat CheckResAccessResp
		return dat, err
	}
	dat, err := res.Res()
	if err != nil {
		rail.Errorf("Request failed, %v", err)
	}
	return dat, err
}

type FindUserCommonReq struct {
	UserNo string `json:"userNo"`
}

func FindUserCommon(rail miso.Rail, req FindUserCommonReq) (common.User, error) {
	var res miso.GnResp[common.User]
	err := miso.NewDynClient(rail, "/internal/v1/user/info/common", "user-vault").
		PostJson(req).
		Json(&res)
	if err != nil {
		rail.Errorf("Request failed, %v", err)
		var dat common.User
		return dat, err
	}
	dat, err := res.Res()
	if err != nil {
		rail.Errorf("Request failed, %v", err)
	}
	return dat, err
}
