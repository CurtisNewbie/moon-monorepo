package vault

import (
	"strings"

	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/user-vault/auth"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"github.com/curtisnewbie/user-vault/api"
	"gorm.io/gorm"
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

// misoapi-http: POST /open/api/user/login
// misoapi-desc: User Login using password, a JWT token is generated and returned
// misoapi-scope: PUBLIC
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

// misoapi-http: POST /open/api/user/register/request
// misoapi-desc: User request registration, approval needed
// misoapi-scope: PUBLIC
func UserRegisterEp(inb *miso.Inbound, req RegisterReq) (any, error) {
	return nil, UserRegister(inb.Rail(), mysql.GetMySQL(), req)
}

// misoapi-http: POST /open/api/user/add
// misoapi-desc: Admin create new user
// misoapi-resource: ref(ResourceManagerUser)
func AdminAddUserEp(inb *miso.Inbound, req AddUserParam) (any, error) {
	return nil, NewUser(inb.Rail(), mysql.GetMySQL(), CreateUserParam{
		Username:     req.Username,
		Password:     req.Password,
		RoleNo:       req.RoleNo,
		ReviewStatus: api.ReviewApproved,
	})
}

// misoapi-http: POST /open/api/user/list
// misoapi-desc: Admin list users
// misoapi-resource: ref(ResourceManagerUser)
func AdminListUsersEp(inb *miso.Inbound, req ListUserReq) (miso.PageRes[api.UserInfo], error) {
	return ListUsers(inb.Rail(), mysql.GetMySQL(), req)
}

// misoapi-http: POST /open/api/user/info/update
// misoapi-desc: Admin update user info
// misoapi-resource: ref(ResourceManagerUser)
func AdminUpdateUserEp(inb *miso.Inbound, req AdminUpdateUserReq) (any, error) {
	rail := inb.Rail()
	return nil, AdminUpdateUser(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

// misoapi-http: POST /open/api/user/registration/review
// misoapi-desc: Admin review user registration
// misoapi-resource: ref(ResourceManagerUser)
func AdminReviewUserEp(inb *miso.Inbound, req AdminReviewUserReq) (any, error) {
	rail := inb.Rail()
	return nil, ReviewUserRegistration(rail, mysql.GetMySQL(), req)
}

// misoapi-http: GET /open/api/user/info
// misoapi-desc: User get user info
// misoapi-scope: PUBLIC
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

// misoapi-http: POST /open/api/user/password/update
// misoapi-desc: User update password
// misoapi-resource: ref(ResourceBasicUser)
func UserUpdatePasswordEp(inb *miso.Inbound, req UpdatePasswordReq) (any, error) {
	rail := inb.Rail()
	u := common.GetUser(rail)
	return nil, UpdatePassword(rail, mysql.GetMySQL(), u.Username, req)
}

// misoapi-http: POST /open/api/token/exchange
// misoapi-desc: Exchange token
// misoapi-scope: PUBLIC
func ExchangeTokenEp(inb *miso.Inbound, req ExchangeTokenReq) (string, error) {
	rail := inb.Rail()
	timer := miso.NewHistTimer(tokenExchangeHisto)
	defer timer.ObserveDuration()
	return ExchangeToken(rail, mysql.GetMySQL(), req)
}

// misoapi-http: GET /open/api/token/user
// misoapi-desc: Get user info by token. This endpoint is expected to be accessible publicly
// misoapi-scope: PUBLIC
func GetTokenUserInfoEp(inb *miso.Inbound, req GetTokenUserReq) (UserInfoBrief, error) {
	rail := inb.Rail()
	return GetTokenUser(rail, mysql.GetMySQL(), req.Token)
}

// misoapi-http: POST /open/api/access/history
// misoapi-desc: User list access logs
// misoapi-resource: ref(ResourceBasicUser)
func UserListAccessHistoryEp(inb *miso.Inbound, req ListAccessLogReq) (miso.PageRes[ListedAccessLog], error) {
	rail := inb.Rail()
	return ListAccessLogs(rail, mysql.GetMySQL(), common.GetUser(rail), req)
}

// misoapi-http: POST /open/api/user/key/generate
// misoapi-desc: User generate user key
// misoapi-resource: ref(ResourceBasicUser)
func UserGenUserKeyEp(inb *miso.Inbound, req GenUserKeyReq) (any, error) {
	rail := inb.Rail()
	return nil, GenUserKey(rail, mysql.GetMySQL(), req, common.GetUser(rail).Username)
}

// misoapi-http: POST /open/api/user/key/list
// misoapi-desc: User list user keys
// misoapi-resource: ref(ResourceBasicUser)
func UserListUserKeysEp(inb *miso.Inbound, req ListUserKeysReq) (miso.PageRes[ListedUserKey], error) {
	rail := inb.Rail()
	return ListUserKeys(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

// misoapi-http: POST /open/api/user/key/delete
// misoapi-desc: User delete user key
// misoapi-resource: ref(ResourceBasicUser)
func UserDeleteUserKeyEp(inb *miso.Inbound, req DeleteUserKeyReq) (any, error) {
	rail := inb.Rail()
	return nil, DeleteUserKey(rail, mysql.GetMySQL(), req, common.GetUser(rail).UserNo)
}

// misoapi-http: POST /open/api/resource/add
// misoapi-desc: Admin add resource
// misoapi-resource: ref(ResourceManageResources)
func AdminAddResourceEp(inb *miso.Inbound, req CreateResReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return nil, CreateResourceIfNotExist(rail, req, user)
}

// misoapi-http: POST /open/api/resource/remove
// misoapi-desc: Admin remove resource
// misoapi-resource: ref(ResourceManageResources)
func AdminRemoveResourceEp(inb *miso.Inbound, req DeleteResourceReq) (any, error) {
	rail := inb.Rail()
	return nil, DeleteResource(rail, req)
}

// misoapi-http: GET /open/api/resource/brief/candidates
// misoapi-desc: List all resource candidates for role
// misoapi-resource: ref(ResourceManageResources)
func ListResCandidatesEp(inb *miso.Inbound, req ListResCandidatesReq) ([]ResBrief, error) {
	rail := inb.Rail()
	return ListResourceCandidatesForRole(rail, req.RoleNo)
}

// misoapi-http: POST /open/api/resource/list
// misoapi-desc: Admin list resources
// misoapi-resource: ref(ResourceManageResources)
func AdminListResEp(inb *miso.Inbound, req ListResReq) (ListResResp, error) {
	rail := inb.Rail()
	return ListResources(rail, req)
}

// misoapi-http: GET /open/api/resource/brief/user
// misoapi-desc: List resources that are accessible to current user
// misoapi-scope: PUBLIC
func ListUserAccessibleResEp(inb *miso.Inbound) ([]ResBrief, error) {
	rail := inb.Rail()
	u := common.GetUser(rail)
	if u.IsNil {
		return []ResBrief{}, nil
	}
	return ListAllResBriefsOfRole(rail, u.RoleNo)
}

// misoapi-http: GET /open/api/resource/brief/all
// misoapi-desc: List all resource brief info
// misoapi-scope: PUBLIC
func ListAllResBriefEp(inb *miso.Inbound) ([]ResBrief, error) {
	rail := inb.Rail()
	return ListAllResBriefs(rail)
}

// misoapi-http: POST /open/api/role/resource/add
// misoapi-desc: Admin add resource to role
// misoapi-resource: ref(ResourceManageResources)
func AdminBindRoleResEp(inb *miso.Inbound, req AddRoleResReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return nil, AddResToRoleIfNotExist(rail, req, user)
}

// misoapi-http: POST /open/api/role/resource/remove
// misoapi-desc: Admin remove resource from role
// misoapi-resource: ref(ResourceManageResources)
func AdminUnbindRoleResEp(inb *miso.Inbound, req RemoveRoleResReq) (any, error) {
	rail := inb.Rail()
	return nil, RemoveResFromRole(rail, req)
}

// misoapi-http: POST /open/api/role/add
// misoapi-desc: Admin add role
// misoapi-resource: ref(ResourceManageResources)
func AdminAddRoleEp(inb *miso.Inbound, req AddRoleReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return nil, AddRole(rail, req, user)
}

// misoapi-http: POST /open/api/role/list
// misoapi-desc: Admin list roles
// misoapi-resource: ref(ResourceManageResources)
func AdminListRolesEp(inb *miso.Inbound, req ListRoleReq) (ListRoleResp, error) {
	rail := inb.Rail()
	return ListRoles(rail, req)
}

// misoapi-http: GET /open/api/role/brief/all
// misoapi-desc: Admin list role brief info
// misoapi-resource: ref(ResourceManageResources)
func AdminListRoleBriefsEp(inb *miso.Inbound) ([]RoleBrief, error) {
	rail := inb.Rail()
	return ListAllRoleBriefs(rail)
}

// misoapi-http: POST /open/api/role/resource/list
// misoapi-desc: Admin list resources of role
// misoapi-resource: ref(ResourceManageResources)
func AdminListRoleResEp(inb *miso.Inbound, req ListRoleResReq) (ListRoleResResp, error) {
	rail := inb.Rail()
	return ListRoleRes(rail, req)
}

// misoapi-http: POST /open/api/role/info
// misoapi-desc: Get role info
// misoapi-scope: PUBLIC
func GetRoleInfoEp(inb *miso.Inbound, req api.RoleInfoReq) (api.RoleInfoResp, error) {
	rail := inb.Rail()
	return GetRoleInfo(rail, req)
}

// misoapi-http: POST /open/api/path/list
// misoapi-desc: Admin list paths
// misoapi-resource: ref(ResourceManageResources)
func AdminListPathsEp(inb *miso.Inbound, req ListPathReq) (ListPathResp, error) {
	rail := inb.Rail()
	return ListPaths(rail, req)
}

// misoapi-http: POST /open/api/path/resource/bind
// misoapi-desc: Admin bind resource to path
// misoapi-resource: ref(ResourceManageResources)
func AdminBindResPathEp(inb *miso.Inbound, req BindPathResReq) (any, error) {
	rail := inb.Rail()
	return nil, BindPathRes(rail, req)
}

// misoapi-http: POST /open/api/path/resource/unbind
// misoapi-desc: Admin unbind resource and path
// misoapi-resource: ref(ResourceManageResources)
func AdminUnbindResPathEp(inb *miso.Inbound, req UnbindPathResReq) (any, error) {
	rail := inb.Rail()
	return nil, UnbindPathRes(rail, req)
}

// misoapi-http: POST /open/api/path/delete
// misoapi-desc: Admin delete path
// misoapi-resource: ref(ResourceManageResources)
func AdminDeletePathEp(inb *miso.Inbound, req DeletePathReq) (any, error) {
	rail := inb.Rail()
	return nil, DeletePath(rail, req)
}

// misoapi-http: POST /open/api/path/update
// misoapi-desc: Admin update path
// misoapi-resource: ref(ResourceManageResources)
func AdminUpdatePathEp(inb *miso.Inbound, req UpdatePathReq) (any, error) {
	rail := inb.Rail()
	return nil, UpdatePath(rail, req)
}

// misoapi-http: POST /remote/user/info
// misoapi-desc: Fetch user info
func ItnFetchUserInfoEp(inb *miso.Inbound, req api.FindUserReq) (api.UserInfo, error) {
	rail := inb.Rail()
	return ItnFindUserInfo(rail, mysql.GetMySQL(), req)
}

// misoapi-http: GET /remote/user/id
// misoapi-desc: Fetch id of user with the username
func ItnFetchUserIdByNameEp(inb *miso.Inbound, req FetchUserIdByNameReq) (int, error) {
	rail := inb.Rail()
	u, err := LoadUserBriefThrCache(rail, mysql.GetMySQL(), req.Username)
	return u.Id, err
}

// misoapi-http: POST /remote/user/userno/username
// misoapi-desc: Fetch usernames of users with the userNos
func ItnFetchUsernamesByNosEp(inb *miso.Inbound, req api.FetchNameByUserNoReq) (api.FetchUsernamesRes, error) {
	rail := inb.Rail()
	return ItnFindNameOfUserNo(rail, mysql.GetMySQL(), req)
}

// misoapi-http: POST /remote/user/list/with-role
// misoapi-desc: Fetch users with the role_no
func ItnFindUserWithRoleEp(inb *miso.Inbound, req api.FetchUsersWithRoleReq) ([]api.UserInfo, error) {
	rail := inb.Rail()
	return ItnFindUsersWithRole(rail, mysql.GetMySQL(), req)
}

// misoapi-http: POST /remote/user/list/with-resource
// misoapi-desc: Fetch users that have access to the resource
func ItnFindUserWithResourceEp(inb *miso.Inbound, req api.FetchUserWithResourceReq) ([]api.UserInfo, error) {
	rail := inb.Rail()
	return FindUserWithRes(rail, mysql.GetMySQL(), req)
}

// misoapi-http: POST /remote/resource/add
// misoapi-desc: Report resource. This endpoint should be used internally by another backend service.
func ItnReportResourceEp(inb *miso.Inbound, req CreateResReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return nil, CreateResourceIfNotExist(rail, req, user)
}

// misoapi-http: POST /remote/path/resource/access-test
// misoapi-desc: Validate resource access
func ItnCheckResourceAccessEp(inb *miso.Inbound, req TestResAccessReq) (TestResAccessResp, error) {
	rail := inb.Rail()
	timer := miso.NewHistTimer(resourceAccessCheckHisto)
	defer timer.ObserveDuration()
	return TestResourceAccess(rail, req)
}

// misoapi-http: POST /remote/path/add
// misoapi-desc: Report endpoint info
func ItnReportPathEp(inb *miso.Inbound, req CreatePathReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return nil, CreatePath(rail, req, user)
}

type ListSitePasswordReq struct {
	Alias    string
	Site     string
	Username string
	Paging   miso.Paging
}

type ListSitePasswordRes struct {
	RecordId   string
	Site       string
	Alias      string
	Username   string
	CreateTime util.ETime
}

// misoapi-http: POST /open/api/password/list-site-passwords
// misoapi-desc: List site password records
// misoapi-resource: ref(ResourceBasicUser)
func ApiListSitePasswords(rail miso.Rail, req ListSitePasswordReq, user common.User, db *gorm.DB) (miso.PageRes[ListSitePasswordRes], error) {
	return ListSitePasswords(rail, req, user, db)
}

type AddSitePasswordReq struct {
	Site          string
	Alias         string
	Username      string `valid:"notEmpty"`
	SitePassword  string `valid:"notEmpty"`
	LoginPassword string `valid:"notEmpty"`
}

// misoapi-http: POST /open/api/password/add-site-password
// misoapi-desc: Add site password record
// misoapi-resource: ref(ResourceBasicUser)
func ApiAddSitePassword(rail miso.Rail, req AddSitePasswordReq, user common.User, db *gorm.DB) (any, error) {
	return nil, AddSitePassword(rail, req, user, db)
}

type RemoveSitePasswordRes struct {
	RecordId string `valid:"notEmpty"`
}

// misoapi-http: POST /open/api/password/remove-site-password
// misoapi-desc: Remove site password record
// misoapi-resource: ref(ResourceBasicUser)
func ApiRemoveSitePassword(rail miso.Rail, req RemoveSitePasswordRes, user common.User, db *gorm.DB) (any, error) {
	return nil, RemoveSitePassword(rail, req, user, db)
}

type DecryptSitePasswordReq struct {
	LoginPassword string `valid:"notEmpty"`
	RecordId      string `valid:"notEmpty"`
}

type DecryptSitePasswordRes struct {
	Decrypted string
}

// misoapi-http: POST /open/api/password/decrypt-site-password
// misoapi-desc: Decrypt site password
// misoapi-resource: ref(ResourceBasicUser)
func ApiDecryptSitePassword(rail miso.Rail, req DecryptSitePasswordReq, user common.User, db *gorm.DB) (DecryptSitePasswordRes, error) {
	return DecryptSitePassword(rail, req, user, db)
}

type EditSitePasswordReq struct {
	RecordId string
	Site     string
	Alias    string
}

// misoapi-http: POST /open/api/password/edit-site-password
// misoapi-desc: Edit site password
// misoapi-resource: ref(ResourceBasicUser)
func ApiEditSitePassword(rail miso.Rail, req EditSitePasswordReq, user common.User, db *gorm.DB) (any, error) {
	return nil, EditSitePassword(rail, req, user, db)
}
