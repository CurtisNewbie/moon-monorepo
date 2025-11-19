package api

import "github.com/curtisnewbie/miso/util/atom"

type FindUserReq struct {
	UserId   *int    `json:"userId"`
	UserNo   *string `json:"userNo"`
	Username *string `json:"username"`
}

type UserInfo struct {
	Id           int
	Username     string
	RoleName     string
	RoleNo       string
	UserNo       string
	ReviewStatus string
	IsDisabled   int
	CreateTime   atom.Time
	CreateBy     string
	UpdateTime   atom.Time
	UpdateBy     string
}

type FetchNameByUserNoReq struct {
	UserNos []string `json:"userNos"`
}

type FetchUsernamesRes struct {
	UserNoToUsername map[string]string `json:"userNoToUsername"`
}

type FetchUsersWithRoleReq struct {
	RoleNo string `valid:"notEmpty" json:"roleNo"`
}

type RoleInfoReq struct {
	RoleNo string `json:"roleNo" validation:"notEmpty"`
}

type RoleInfoResp struct {
	RoleNo string `json:"roleNo"`
	Name   string `json:"name"`
}

type FetchUserWithResourceReq struct {
	ResourceCode string `json:"resourceCode"`
}

type CreateNotificationReq struct {
	Title           string   `valid:"maxLen:255" json:"title"`
	Message         string   `valid:"maxLen:1000" json:"message"`
	ReceiverUserNos []string `json:"receiverUserNos"`
}

type CheckResAccessReq struct {
	RoleNo string `json:"roleNo"`
	Url    string `json:"url"`
	Method string `json:"method"`
}

type CheckResAccessResp struct {
	Valid bool `json:"valid"`
}
