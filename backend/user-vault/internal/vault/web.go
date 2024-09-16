package vault

import (
	"strings"

	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/user-vault/auth"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"github.com/curtisnewbie/user-vault/api"
)

const (
	passwordLoginUrl = "/user-vault/open/api/user/login"

	ResourceManagerUser     = "manage-users"
	ResourceBasicUser       = "basic-user"
	ResourceManageResources = "manage-resources"
)

var (
	fetchUserInfoHisto       = miso.NewPromHisto("user_vault_fetch_user_info_duration")
	tokenExchangeHisto       = miso.NewPromHisto("user_vault_token_exchange_duration")
	resourceAccessCheckHisto = miso.NewPromHisto("user_vault_resource_access_check_duration")
)

type LoginReq struct {
	Username      string `json:"username" valid:"notEmpty"`
	Password      string `json:"password" valid:"notEmpty"`
	XForwardedFor string `header:"x-forwarded-for"`
	UserAgent     string `header:"user-agent"`
}

type AdminAddUserReq struct {
	Username string `json:"username" valid:"notEmpty"`
	Password string `json:"password" valid:"notEmpty"`
	RoleNo   string `json:"roleNo" valid:"notEmpty"`
}

type ListUserReq struct {
	Username   *string     `json:"username"`
	RoleNo     *string     `json:"roleNo"`
	IsDisabled *int        `json:"isDisabled"`
	Paging     miso.Paging `json:"paging"`
}

type AdminUpdateUserReq struct {
	UserNo     string `valid:"notEmpty"`
	RoleNo     string `json:"roleNo"`
	IsDisabled int    `json:"isDisabled"`
}

type AdminReviewUserReq struct {
	UserId       int    `json:"userId" valid:"positive"`
	ReviewStatus string `json:"reviewStatus"`
}

type RegisterReq struct {
	Username string `json:"username" valid:"notEmpty"`
	Password string `json:"password" valid:"notEmpty"`
}

type UserInfoRes struct {
	Id           int
	Username     string
	RoleName     string
	RoleNo       string
	UserNo       string
	RegisterDate string
}

type GetTokenUserReq struct {
	Token string `form:"token" desc:"jwt token"`
}

type ListResCandidatesReq struct {
	RoleNo string `form:"roleNo" desc:"Role No"`
}

type FetchUserIdByNameReq struct {
	Username string `form:"username" desc:"Username"`
}

func RegisterInternalPathResourcesOnBootstrapped(res []auth.Resource) {

	miso.PostServerBootstrapped(func(rail miso.Rail) error {

		user := common.NilUser()

		app := miso.GetPropStr(miso.PropAppName)
		for _, res := range res {
			if res.Code == "" || res.Name == "" {
				continue
			}
			if e := CreateResourceIfNotExist(rail, CreateResReq(res), user); e != nil {
				return e
			}
		}

		routes := miso.GetHttpRoutes()
		for _, route := range routes {
			if route.Url == "" {
				continue
			}
			var routeType = PathTypeProtected
			if route.Scope == miso.ScopePublic {
				routeType = PathTypePublic
			}

			url := route.Url
			if !strings.HasPrefix(url, "/") {
				url = "/" + url
			}

			r := CreatePathReq{
				Method:  route.Method,
				Group:   app,
				Url:     "/" + app + url,
				Type:    routeType,
				Desc:    route.Desc,
				ResCode: route.Resource,
			}
			if err := CreatePath(rail, r, user); err != nil {
				return err
			}
		}
		return nil
	})
}

func RegisterRoutes(rail miso.Rail) error {

	miso.GroupRoute("/open/api",

		miso.IPost("/user/login", UserLoginEp).
			Desc("User Login using password, a JWT token is generated and returned").
			Public(),

		miso.IPost("/user/register/request", UserRegisterEp).
			Desc("User request registration, approval needed").
			Public(),

		miso.IPost("/user/add", AdminAddUserEp).
			Desc("Admin create new user").
			Resource(ResourceManagerUser),

		miso.IPost("/user/list", AdminListUsersEp).
			Desc("Admin list users").
			Resource(ResourceManagerUser),

		miso.IPost("/user/info/update", AdminUpdateUserEp).
			Desc("Admin update user info").
			Resource(ResourceManagerUser),

		miso.IPost("/user/registration/review", AdminReviewUserEp).
			Desc("Admin review user registration").
			Resource(ResourceManagerUser),

		miso.Get("/user/info", UserGetUserInfoEp).
			Desc("User get user info").
			Public(),

		miso.IPost("/user/password/update", UserUpdatePasswordEp).
			Desc("User update password").
			Resource(ResourceBasicUser),

		miso.IPost("/token/exchange", ExchangeTokenEp).
			Desc("Exchange token").
			Public(),

		miso.IGet("/token/user", GetTokenUserInfoEp).
			Desc("Get user info by token. This endpoint is expected to be accessible publicly").
			Public(),

		miso.IPost("/access/history", UserListAccessHistoryEp).
			Desc("User list access logs").
			Resource(ResourceBasicUser),

		miso.IPost("/user/key/generate", UserGenUserKeyEp).
			Desc("User generate user key").
			Resource(ResourceBasicUser),

		miso.IPost("/user/key/list", UserListUserKeysEp).
			Desc("User list user keys").
			Resource(ResourceBasicUser),

		miso.IPost("/user/key/delete", UserDeleteUserKeyEp).
			Desc("User delete user key").
			Resource(ResourceBasicUser),

		miso.IPost("/resource/add", AdminAddResourceEp).
			Desc("Admin add resource").
			Resource(ResourceManageResources),

		miso.IPost("/resource/remove", AdminRemoveResourceEp).
			Desc("Admin remove resource").
			Resource(ResourceManageResources),

		miso.IGet("/resource/brief/candidates", ListResCandidatesEp).
			Desc("List all resource candidates for role").
			Resource(ResourceManageResources),

		miso.IPost("/resource/list", AdminListResEp).
			Desc("Admin list resources").
			Resource(ResourceManageResources),

		miso.Get("/resource/brief/user", ListUserAccessibleResEp).
			Desc("List resources that are accessible to current user").
			Public(),

		miso.Get("/resource/brief/all", ListAllResBriefEp).
			Desc("List all resource brief info").
			Public(),

		miso.IPost("/role/resource/add", AdminBindRoleResEp).
			Desc("Admin add resource to role").
			Resource(ResourceManageResources),

		miso.IPost("/role/resource/remove", AdminUnbindRoleResEp).
			Desc("Admin remove resource from role").
			Resource(ResourceManageResources),

		miso.IPost("/role/add", AdminAddRoleEp).
			Desc("Admin add role").
			Resource(ResourceManageResources),

		miso.IPost("/role/list", AdminListRolesEp).
			Desc("Admin list roles").
			Resource(ResourceManageResources),

		miso.Get("/role/brief/all", AdminListRoleBriefsEp).
			Desc("Admin list role brief info").
			Resource(ResourceManageResources),

		miso.IPost("/role/resource/list", AdminListRoleResEp).
			Desc("Admin list resources of role").
			Resource(ResourceManageResources),

		miso.IPost("/role/info", GetRoleInfoEp).
			Desc("Get role info").
			Public(),

		miso.IPost("/path/list", AdminListPathsEp).
			Desc("Admin list paths").
			Resource(ResourceManageResources),

		miso.IPost("/path/resource/bind", AdminBindResPathEp).
			Desc("Admin bind resource to path").
			Resource(ResourceManageResources),

		miso.IPost("/path/resource/unbind", AdminUnbindResPathEp).
			Desc("Admin unbind resource and path").
			Resource(ResourceManageResources),

		miso.IPost("/path/delete", AdminDeletePathEp).
			Desc("Admin delete path").
			Resource(ResourceManageResources),

		miso.IPost("/path/update", AdminUpdatePathEp).
			Desc("Admin update path").
			Resource(ResourceManageResources),
	)

	// ----------------------------------------------------------------------------------------------
	//
	// Internal endpoints
	//
	// ----------------------------------------------------------------------------------------------

	miso.BaseRoute("/remote").Group(

		miso.IPost("/user/info", ItnFetchUserInfoEp).
			Desc("Fetch user info"),

		miso.IGet("/user/id", ItnFetchUserIdByNameEp).
			Desc("Fetch id of user with the username"),

		miso.IPost("/user/userno/username", ItnFetchUsernamesByNosEp).
			Desc("Fetch usernames of users with the userNos"),

		miso.IPost("/user/list/with-role", ItnFindUserWithRoleEp).
			Desc("Fetch users with the role_no"),

		miso.IPost("/user/list/with-resource", ItnFindUserWithResourceEp).
			Desc("Fetch users that have access to the resource"),

		miso.IPost("/resource/add", ItnReportResourceEp).
			Desc("Report resource. This endpoint should be used internally by another backend service."),

		miso.IPost("/path/resource/access-test", ItnCheckResourceAccessEp).
			Desc("Validate resource access"),

		miso.IPost("/path/add", ItnReportPathEp).
			Desc("Report endpoint info"),
	)
	return nil
}

func UserLoginEp(inb *miso.Inbound, req LoginReq) (string, error) {
	rail := inb.Rail()
	token, user, err := UserLogin(rail, mysql.GetMySQL(),
		PasswordLoginParam{Username: req.Username, Password: req.Password})
	remoteAddr := RemoteAddr(req.XForwardedFor)
	userAgent := req.UserAgent

	if er := AccessLogPipeline.Send(rail, AccessLogEvent{
		IpAddress:  remoteAddr,
		UserAgent:  userAgent,
		UserId:     user.Id,
		Username:   req.Username,
		Url:        passwordLoginUrl,
		Success:    err == nil,
		AccessTime: util.Now(),
	}); er != nil {
		rail.Errorf("Failed to sendAccessLogEvent, username: %v, remoteAddr: %v, userAgent: %v, %v",
			req.Username, remoteAddr, userAgent, er)
	}

	if err != nil {
		return "", err
	}

	return token, err
}

func RemoteAddr(forwardedFor string) string {
	addr := "unknown"

	if forwardedFor != "" {
		tkn := strings.Split(forwardedFor, ",")
		if len(tkn) > 0 {
			addr = tkn[0]
		}
	}
	return addr
}

func UserRegisterEp(inb *miso.Inbound, req RegisterReq) (any, error) {
	return nil, UserRegister(inb.Rail(), mysql.GetMySQL(), req)
}

func AdminAddUserEp(inb *miso.Inbound, req AddUserParam) (any, error) {
	return nil, NewUser(inb.Rail(), mysql.GetMySQL(), CreateUserParam{
		Username:     req.Username,
		Password:     req.Password,
		RoleNo:       req.RoleNo,
		ReviewStatus: api.ReviewApproved,
	})
}

func AdminListUsersEp(inb *miso.Inbound, req ListUserReq) (miso.PageRes[api.UserInfo], error) {
	return ListUsers(inb.Rail(), mysql.GetMySQL(), req)
}

func AdminUpdateUserEp(inb *miso.Inbound, req AdminUpdateUserReq) (any, error) {
	rail := inb.Rail()
	return nil, AdminUpdateUser(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

func AdminReviewUserEp(inb *miso.Inbound, req AdminReviewUserReq) (any, error) {
	rail := inb.Rail()
	return nil, ReviewUserRegistration(rail, mysql.GetMySQL(), req)
}

func UserGetUserInfoEp(inb *miso.Inbound) (UserInfoRes, error) {
	rail := inb.Rail()
	timer := miso.NewHistTimer(fetchUserInfoHisto)
	defer timer.ObserveDuration()
	u := common.GetUser(rail)
	if u.UserNo == "" {
		return UserInfoRes{}, nil
	}

	res, err := LoadUserBriefThrCache(rail, mysql.GetMySQL(), u.Username)

	if err != nil {
		return UserInfoRes{}, err
	}

	return UserInfoRes{
		Id:           res.Id,
		Username:     res.Username,
		RoleName:     res.RoleName,
		RoleNo:       res.RoleNo,
		UserNo:       res.UserNo,
		RegisterDate: res.RegisterDate,
	}, nil
}

func UserUpdatePasswordEp(inb *miso.Inbound, req UpdatePasswordReq) (any, error) {
	rail := inb.Rail()
	u := common.GetUser(rail)
	return nil, UpdatePassword(rail, mysql.GetMySQL(), u.Username, req)
}

func ExchangeTokenEp(inb *miso.Inbound, req ExchangeTokenReq) (string, error) {
	rail := inb.Rail()
	timer := miso.NewHistTimer(tokenExchangeHisto)
	defer timer.ObserveDuration()
	return ExchangeToken(rail, mysql.GetMySQL(), req)
}

func GetTokenUserInfoEp(inb *miso.Inbound, req GetTokenUserReq) (UserInfoBrief, error) {
	rail := inb.Rail()
	return GetTokenUser(rail, mysql.GetMySQL(), req.Token)
}

func UserListAccessHistoryEp(inb *miso.Inbound, req ListAccessLogReq) (miso.PageRes[ListedAccessLog], error) {
	rail := inb.Rail()
	return ListAccessLogs(rail, mysql.GetMySQL(), common.GetUser(rail), req)
}

func UserGenUserKeyEp(inb *miso.Inbound, req GenUserKeyReq) (any, error) {
	rail := inb.Rail()
	return nil, GenUserKey(rail, mysql.GetMySQL(), req, common.GetUser(rail).Username)
}

func UserListUserKeysEp(inb *miso.Inbound, req ListUserKeysReq) (miso.PageRes[ListedUserKey], error) {
	rail := inb.Rail()
	return ListUserKeys(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

func UserDeleteUserKeyEp(inb *miso.Inbound, req DeleteUserKeyReq) (any, error) {
	rail := inb.Rail()
	return nil, DeleteUserKey(rail, mysql.GetMySQL(), req, common.GetUser(rail).UserNo)
}

func AdminAddResourceEp(inb *miso.Inbound, req CreateResReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return nil, CreateResourceIfNotExist(rail, req, user)
}

func AdminRemoveResourceEp(inb *miso.Inbound, req DeleteResourceReq) (any, error) {
	rail := inb.Rail()
	return nil, DeleteResource(rail, req)
}

func ListResCandidatesEp(inb *miso.Inbound, req ListResCandidatesReq) ([]ResBrief, error) {
	rail := inb.Rail()
	return ListResourceCandidatesForRole(rail, req.RoleNo)
}

func AdminListResEp(inb *miso.Inbound, req ListResReq) (ListResResp, error) {
	rail := inb.Rail()
	return ListResources(rail, req)
}

func ListUserAccessibleResEp(inb *miso.Inbound) ([]ResBrief, error) {
	rail := inb.Rail()
	u := common.GetUser(rail)
	if u.IsNil {
		return []ResBrief{}, nil
	}
	return ListAllResBriefsOfRole(rail, u.RoleNo)
}

func ListAllResBriefEp(inb *miso.Inbound) ([]ResBrief, error) {
	rail := inb.Rail()
	return ListAllResBriefs(rail)
}

func AdminBindRoleResEp(inb *miso.Inbound, req AddRoleResReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return nil, AddResToRoleIfNotExist(rail, req, user)
}

func AdminUnbindRoleResEp(inb *miso.Inbound, req RemoveRoleResReq) (any, error) {
	rail := inb.Rail()
	return nil, RemoveResFromRole(rail, req)
}

func AdminAddRoleEp(inb *miso.Inbound, req AddRoleReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return nil, AddRole(rail, req, user)
}

func AdminListRolesEp(inb *miso.Inbound, req ListRoleReq) (ListRoleResp, error) {
	rail := inb.Rail()
	return ListRoles(rail, req)
}

func AdminListRoleBriefsEp(inb *miso.Inbound) ([]RoleBrief, error) {
	rail := inb.Rail()
	return ListAllRoleBriefs(rail)
}

func AdminListRoleResEp(inb *miso.Inbound, req ListRoleResReq) (ListRoleResResp, error) {
	rail := inb.Rail()
	return ListRoleRes(rail, req)
}

func GetRoleInfoEp(inb *miso.Inbound, req api.RoleInfoReq) (api.RoleInfoResp, error) {
	rail := inb.Rail()
	return GetRoleInfo(rail, req)
}

func AdminListPathsEp(inb *miso.Inbound, req ListPathReq) (ListPathResp, error) {
	rail := inb.Rail()
	return ListPaths(rail, req)
}

func AdminBindResPathEp(inb *miso.Inbound, req BindPathResReq) (any, error) {
	rail := inb.Rail()
	return nil, BindPathRes(rail, req)
}

func AdminUnbindResPathEp(inb *miso.Inbound, req UnbindPathResReq) (any, error) {
	rail := inb.Rail()
	return nil, UnbindPathRes(rail, req)
}

func AdminDeletePathEp(inb *miso.Inbound, req DeletePathReq) (any, error) {
	rail := inb.Rail()
	return nil, DeletePath(rail, req)
}

func AdminUpdatePathEp(inb *miso.Inbound, req UpdatePathReq) (any, error) {
	rail := inb.Rail()
	return nil, UpdatePath(rail, req)
}

func ItnFetchUserInfoEp(inb *miso.Inbound, req api.FindUserReq) (api.UserInfo, error) {
	rail := inb.Rail()
	return ItnFindUserInfo(rail, mysql.GetMySQL(), req)
}

func ItnFetchUserIdByNameEp(inb *miso.Inbound, req FetchUserIdByNameReq) (int, error) {
	rail := inb.Rail()
	u, err := LoadUserBriefThrCache(rail, mysql.GetMySQL(), req.Username)
	return u.Id, err
}

func ItnFetchUsernamesByNosEp(inb *miso.Inbound, req api.FetchNameByUserNoReq) (api.FetchUsernamesRes, error) {
	rail := inb.Rail()
	return ItnFindNameOfUserNo(rail, mysql.GetMySQL(), req)
}

func ItnFindUserWithRoleEp(inb *miso.Inbound, req api.FetchUsersWithRoleReq) ([]api.UserInfo, error) {
	rail := inb.Rail()
	return ItnFindUsersWithRole(rail, mysql.GetMySQL(), req)
}

func ItnReportResourceEp(inb *miso.Inbound, req CreateResReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return nil, CreateResourceIfNotExist(rail, req, user)
}

func ItnCheckResourceAccessEp(inb *miso.Inbound, req TestResAccessReq) (TestResAccessResp, error) {
	rail := inb.Rail()
	timer := miso.NewHistTimer(resourceAccessCheckHisto)
	defer timer.ObserveDuration()
	return TestResourceAccess(rail, req)
}

func ItnReportPathEp(inb *miso.Inbound, req CreatePathReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return nil, CreatePath(rail, req, user)
}

func ItnFindUserWithResourceEp(inb *miso.Inbound, req api.FetchUserWithResourceReq) ([]api.UserInfo, error) {
	rail := inb.Rail()
	return FindUserWithRes(rail, mysql.GetMySQL(), req)
}
