package api

import "github.com/curtisnewbie/miso/util"

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
	CreateTime   util.ETime
	CreateBy     string
	UpdateTime   util.ETime
	UpdateBy     string
}

type FetchNameByUserNoReq struct {
	UserNos []string `json:"userNos"`
}

type FetchUsernamesRes struct {
	UserNoToUsername map[string]string `json:"userNoToUsername"`
}

type FetchUsersWithRoleReq struct {
	RoleNo string `valid:"notEmpty"`
}

type RoleInfoReq struct {
	RoleNo string `json:"roleNo" validation:"notEmpty"`
}

type RoleInfoResp struct {
	RoleNo string `json:"roleNo"`
	Name   string `json:"name"`
}

type FetchUserWithResourceReq struct {
	ResourceCode string
}

type CreateNotificationReq struct {
	Title           string `valid:"maxLen:255"`
	Message         string `valid:"maxLen:1000"`
	ReceiverUserNos []string
}
