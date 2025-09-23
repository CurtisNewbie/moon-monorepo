package web

import (
	"strings"

	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"github.com/curtisnewbie/user-vault/api"
	"github.com/curtisnewbie/user-vault/internal/metrics"
	"github.com/curtisnewbie/user-vault/internal/note"
	"github.com/curtisnewbie/user-vault/internal/vault"
	"gorm.io/gorm"
)

const (
	passwordLoginUrl = "/user-vault/open/api/user/login"

	ResourceManagerUser     = "manage-users"
	ResourceBasicUser       = "basic-user"
	ResourceManageResources = "manage-resources"
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

// misoapi-http: POST /open/api/user/login
// misoapi-desc: User Login using password, a JWT token is generated and returned
// misoapi-scope: PUBLIC
func ApiUserLogin(inb *miso.Inbound, req LoginReq) (string, error) {
	rail := inb.Rail()
	token, user, err := vault.UserLogin(rail, mysql.GetMySQL(),
		vault.PasswordLoginParam{Username: req.Username, Password: req.Password})
	remoteAddr := RemoteAddr(req.XForwardedFor)
	userAgent := req.UserAgent

	if er := vault.AccessLogPipeline.Send(rail, vault.AccessLogEvent{
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
func ApiUserRegister(inb *miso.Inbound, req vault.RegisterReq) (any, error) {
	return nil, vault.UserRegister(inb.Rail(), mysql.GetMySQL(), req)
}

// misoapi-http: POST /open/api/user/add
// misoapi-desc: Admin create new user
// misoapi-resource: ref(ResourceManagerUser)
func ApiAdminAddUser(inb *miso.Inbound, req vault.AddUserParam) (any, error) {
	return nil, vault.NewUser(inb.Rail(), mysql.GetMySQL(), vault.CreateUserParam{
		Username:     req.Username,
		Password:     req.Password,
		RoleNo:       req.RoleNo,
		ReviewStatus: api.ReviewApproved,
	})
}

// misoapi-http: POST /open/api/user/list
// misoapi-desc: Admin list users
// misoapi-resource: ref(ResourceManagerUser)
func ApiAdminListUsers(inb *miso.Inbound, req vault.ListUserReq) (miso.PageRes[vault.UserInfo], error) {
	return vault.ListUsers(inb.Rail(), mysql.GetMySQL(), req)
}

// misoapi-http: POST /open/api/user/info/update
// misoapi-desc: Admin update user info
// misoapi-resource: ref(ResourceManagerUser)
func ApiAdminUpdateUser(inb *miso.Inbound, req vault.AdminUpdateUserReq) (any, error) {
	rail := inb.Rail()
	return nil, vault.AdminUpdateUser(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

// misoapi-http: POST /open/api/user/registration/review
// misoapi-desc: Admin review user registration
// misoapi-resource: ref(ResourceManagerUser)
func ApiAdminReviewUser(inb *miso.Inbound, req vault.AdminReviewUserReq) (any, error) {
	rail := inb.Rail()
	return nil, vault.ReviewUserRegistration(rail, mysql.GetMySQL(), req)
}

// misoapi-http: GET /open/api/user/info
// misoapi-desc: User get user info
// misoapi-scope: PUBLIC
func ApiUserGetUserInfo(inb *miso.Inbound) (UserInfoRes, error) {
	rail := inb.Rail()
	timer := miso.NewHistTimer(metrics.FetchUserInfoHisto)
	defer timer.ObserveDuration()
	u := common.GetUser(rail)
	if u.UserNo == "" {
		return UserInfoRes{}, nil
	}

	res, err := vault.LoadUserBriefThrCache(rail, mysql.GetMySQL(), u.Username)

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
func ApiUserUpdatePassword(inb *miso.Inbound, req vault.UpdatePasswordReq) (any, error) {
	rail := inb.Rail()
	u := common.GetUser(rail)
	return nil, vault.UpdatePassword(rail, mysql.GetMySQL(), u.Username, req)
}

// misoapi-http: POST /open/api/token/exchange
// misoapi-desc: Exchange token
// misoapi-scope: PUBLIC
func ExchangeTokenEp(inb *miso.Inbound, req vault.ExchangeTokenReq) (string, error) {
	rail := inb.Rail()
	timer := miso.NewHistTimer(metrics.TokenExchangeHisto)
	defer timer.ObserveDuration()
	return vault.ExchangeToken(rail, req)
}

// misoapi-http: GET /open/api/token/user
// misoapi-desc: Get user info by token. This endpoint is expected to be accessible publicly
// misoapi-scope: PUBLIC
func ApiGetTokenUserInfo(inb *miso.Inbound, req GetTokenUserReq) (vault.UserInfoBrief, error) {
	rail := inb.Rail()
	return vault.GetTokenUser(rail, mysql.GetMySQL(), req.Token)
}

// misoapi-http: POST /open/api/access/history
// misoapi-desc: User list access logs
// misoapi-resource: ref(ResourceBasicUser)
func ApiUserListAccessHistory(inb *miso.Inbound, req vault.ListAccessLogReq) (miso.PageRes[vault.ListedAccessLog], error) {
	rail := inb.Rail()
	return vault.ListAccessLogs(rail, mysql.GetMySQL(), common.GetUser(rail), req)
}

// misoapi-http: POST /open/api/user/key/generate
// misoapi-desc: User generate user key
// misoapi-resource: ref(ResourceBasicUser)
func ApiUserGenUserKey(inb *miso.Inbound, req vault.GenUserKeyReq) (any, error) {
	rail := inb.Rail()
	return nil, vault.GenUserKey(rail, mysql.GetMySQL(), req, common.GetUser(rail).Username)
}

// misoapi-http: POST /open/api/user/key/list
// misoapi-desc: User list user keys
// misoapi-resource: ref(ResourceBasicUser)
func ApiUserListUserKeys(inb *miso.Inbound, req vault.ListUserKeysReq) (miso.PageRes[vault.ListedUserKey], error) {
	rail := inb.Rail()
	return vault.ListUserKeys(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

// misoapi-http: POST /open/api/user/key/delete
// misoapi-desc: User delete user key
// misoapi-resource: ref(ResourceBasicUser)
func ApiUserDeleteUserKey(inb *miso.Inbound, req vault.DeleteUserKeyReq) (any, error) {
	rail := inb.Rail()
	return nil, vault.DeleteUserKey(rail, mysql.GetMySQL(), req, common.GetUser(rail).UserNo)
}

// misoapi-http: POST /open/api/resource/add
// misoapi-desc: Admin add resource
// misoapi-resource: ref(ResourceManageResources)
func ApiAdminAddResource(inb *miso.Inbound, req vault.CreateResReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return nil, vault.CreateResourceIfNotExist(rail, req, user)
}

// misoapi-http: POST /open/api/resource/remove
// misoapi-desc: Admin remove resource
// misoapi-resource: ref(ResourceManageResources)
func ApiAdminRemoveResource(inb *miso.Inbound, req vault.DeleteResourceReq) (any, error) {
	rail := inb.Rail()
	return nil, vault.DeleteResource(rail, req)
}

// misoapi-http: GET /open/api/resource/brief/candidates
// misoapi-desc: List all resource candidates for role
// misoapi-resource: ref(ResourceManageResources)
func ApiListResCandidates(inb *miso.Inbound, req ListResCandidatesReq) ([]vault.ResBrief, error) {
	rail := inb.Rail()
	return vault.ListResourceCandidatesForRole(rail, req.RoleNo)
}

// misoapi-http: POST /open/api/resource/list
// misoapi-desc: Admin list resources
// misoapi-resource: ref(ResourceManageResources)
func ApiAdminListRes(inb *miso.Inbound, req vault.ListResReq) (vault.ListResResp, error) {
	rail := inb.Rail()
	return vault.ListResources(rail, req)
}

// misoapi-http: GET /open/api/resource/brief/user
// misoapi-desc: List resources that are accessible to current user
// misoapi-scope: PUBLIC
func ApiListUserAccessibleRes(inb *miso.Inbound) ([]vault.ResBrief, error) {
	rail := inb.Rail()
	u := common.GetUser(rail)
	if u.IsNil {
		return []vault.ResBrief{}, nil
	}
	return vault.ListAllResBriefsOfRole(rail, u.RoleNo)
}

// misoapi-http: GET /open/api/resource/brief/all
// misoapi-desc: List all resource brief info
// misoapi-scope: PUBLIC
func ApiListAllResBrief(inb *miso.Inbound) ([]vault.ResBrief, error) {
	rail := inb.Rail()
	return vault.ListAllResBriefs(rail)
}

// misoapi-http: POST /open/api/role/resource/add
// misoapi-desc: Admin add resource to role
// misoapi-resource: ref(ResourceManageResources)
func ApiAdminBindRoleRes(inb *miso.Inbound, req vault.AddRoleResReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return nil, vault.AddResToRoleIfNotExist(rail, req, user)
}

// misoapi-http: POST /open/api/role/resource/remove
// misoapi-desc: Admin remove resource from role
// misoapi-resource: ref(ResourceManageResources)
func ApiAdminUnbindRoleRes(inb *miso.Inbound, req vault.RemoveRoleResReq) (any, error) {
	rail := inb.Rail()
	return nil, vault.RemoveResFromRole(rail, req)
}

// misoapi-http: POST /open/api/role/add
// misoapi-desc: Admin add role
// misoapi-resource: ref(ResourceManageResources)
func ApiAdminAddRole(inb *miso.Inbound, req vault.AddRoleReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return nil, vault.AddRole(rail, req, user)
}

// misoapi-http: POST /open/api/role/list
// misoapi-desc: Admin list roles
// misoapi-resource: ref(ResourceManageResources)
func ApiAdminListRoles(inb *miso.Inbound, req vault.ListRoleReq) (vault.ListRoleResp, error) {
	rail := inb.Rail()
	return vault.ListRoles(rail, req)
}

// misoapi-http: GET /open/api/role/brief/all
// misoapi-desc: Admin list role brief info
// misoapi-resource: ref(ResourceManageResources)
func ApiAdminListRoleBriefs(inb *miso.Inbound) ([]vault.RoleBrief, error) {
	rail := inb.Rail()
	return vault.ListAllRoleBriefs(rail)
}

// misoapi-http: POST /open/api/role/resource/list
// misoapi-desc: Admin list resources of role
// misoapi-resource: ref(ResourceManageResources)
func ApiAdminListRoleRes(inb *miso.Inbound, req vault.ListRoleResReq) (vault.ListRoleResResp, error) {
	rail := inb.Rail()
	return vault.ListRoleRes(rail, req)
}

// misoapi-http: POST /open/api/role/info
// misoapi-desc: Get role info
// misoapi-scope: PUBLIC
func ApiGetRoleInfo(inb *miso.Inbound, req api.RoleInfoReq) (api.RoleInfoResp, error) {
	rail := inb.Rail()
	return vault.GetRoleInfo(rail, req)
}

// misoapi-http: POST /open/api/path/list
// misoapi-desc: Admin list paths
// misoapi-resource: ref(ResourceManageResources)
func ApiAdminListPaths(inb *miso.Inbound, req vault.ListPathReq) (vault.ListPathResp, error) {
	rail := inb.Rail()
	return vault.ListPaths(rail, req)
}

// misoapi-http: POST /open/api/path/resource/bind
// misoapi-desc: Admin bind resource to path
// misoapi-resource: ref(ResourceManageResources)
func ApiAdminBindResPath(inb *miso.Inbound, req vault.BindPathResReq) (any, error) {
	rail := inb.Rail()
	return nil, vault.BindPathRes(rail, req)
}

// misoapi-http: POST /open/api/path/resource/unbind
// misoapi-desc: Admin unbind resource and path
// misoapi-resource: ref(ResourceManageResources)
func ApiAdminUnbindResPath(inb *miso.Inbound, req vault.UnbindPathResReq) (any, error) {
	rail := inb.Rail()
	return nil, vault.UnbindPathRes(rail, req)
}

// misoapi-http: POST /open/api/path/delete
// misoapi-desc: Admin delete path
// misoapi-resource: ref(ResourceManageResources)
func ApiAdminDeletePath(inb *miso.Inbound, req vault.DeletePathReq) (any, error) {
	rail := inb.Rail()
	return nil, vault.DeletePath(rail, req)
}

// misoapi-http: POST /open/api/path/update
// misoapi-desc: Admin update path
// misoapi-resource: ref(ResourceManageResources)
func ApiAdminUpdatePath(inb *miso.Inbound, req vault.UpdatePathReq) (any, error) {
	rail := inb.Rail()
	return nil, vault.UpdatePath(rail, req)
}

// misoapi-http: POST /remote/user/info
// misoapi-desc: Fetch user info
func ApiFetchUserInfo(inb *miso.Inbound, req api.FindUserReq) (vault.UserInfo, error) {
	rail := inb.Rail()
	return vault.ItnFindUserInfo(rail, mysql.GetMySQL(), req)
}

// misoapi-http: POST /internal/v1/user/info/common
// misoapi-desc: System fetch user info as common.User
func ApiSysFetchUserInfo(inb *miso.Inbound, req api.FindUserReq) (common.User, error) {
	rail := inb.Rail()
	v, err := vault.ItnFindUserInfo(rail, mysql.GetMySQL(), req)
	if err != nil {
		return common.User{}, err
	}
	return common.User{
		UserNo:   v.UserNo,
		Username: v.Username,
		RoleNo:   v.RoleNo,
		IsNil:    false,
	}, nil
}

// misoapi-http: GET /remote/user/id
// misoapi-desc: Fetch id of user with the username
func ApiFetchUserIdByName(inb *miso.Inbound, req FetchUserIdByNameReq) (int, error) {
	rail := inb.Rail()
	u, err := vault.LoadUserBriefThrCache(rail, mysql.GetMySQL(), req.Username)
	return u.Id, err
}

// misoapi-http: POST /remote/user/userno/username
// misoapi-desc: Fetch usernames of users with the userNos
func ApiFetchUsernamesByNosEp(inb *miso.Inbound, req api.FetchNameByUserNoReq) (api.FetchUsernamesRes, error) {
	rail := inb.Rail()
	return vault.ItnFindNameOfUserNo(rail, mysql.GetMySQL(), req)
}

// misoapi-http: POST /remote/user/list/with-role
// misoapi-desc: Fetch users with the role_no
func ApiFindUserWithRoleEp(inb *miso.Inbound, req api.FetchUsersWithRoleReq) ([]vault.UserInfo, error) {
	rail := inb.Rail()
	return vault.ItnFindUsersWithRole(rail, mysql.GetMySQL(), req)
}

// misoapi-http: POST /remote/user/list/with-resource
// misoapi-desc: Fetch users that have access to the resource
func ApiFindUserWithResourceEp(inb *miso.Inbound, req api.FetchUserWithResourceReq) ([]vault.UserInfo, error) {
	rail := inb.Rail()
	return vault.FindUserWithRes(rail, mysql.GetMySQL(), req)
}

// misoapi-http: POST /remote/resource/add
// misoapi-desc: Report resource. This endpoint should be used internally by another backend service.
func ApiReportResourceEp(inb *miso.Inbound, req vault.CreateResReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return nil, vault.CreateResourceIfNotExist(rail, req, user)
}

// misoapi-http: POST /remote/path/resource/access-test
// misoapi-desc: Validate resource access
func ApiCheckResourceAccessEp(inb *miso.Inbound, req api.CheckResAccessReq) (api.CheckResAccessResp, error) {
	rail := inb.Rail()
	timer := miso.NewHistTimer(metrics.ResourceAccessCheckHisto)
	defer timer.ObserveDuration()
	return vault.TestResourceAccess(rail, req)
}

// misoapi-http: POST /remote/path/add
// misoapi-desc: Report endpoint info
// misoapi-resource: ref(ResourceManageResources)
func ApiReportPath(inb *miso.Inbound, req vault.CreatePathReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return nil, vault.CreatePath(rail, req, user)
}

// misoapi-http: POST /open/api/password/list-site-passwords
// misoapi-desc: List site password records
// misoapi-resource: ref(ResourceBasicUser)
func ApiListSitePasswords(rail miso.Rail, req vault.ListSitePasswordReq, user common.User, db *gorm.DB) (miso.PageRes[vault.ListSitePasswordRes], error) {
	return vault.ListSitePasswords(rail, req, user, db)
}

// misoapi-http: POST /open/api/password/add-site-password
// misoapi-desc: Add site password record
// misoapi-resource: ref(ResourceBasicUser)
func ApiAddSitePassword(rail miso.Rail, req vault.AddSitePasswordReq, user common.User, db *gorm.DB) (any, error) {
	return nil, vault.AddSitePassword(rail, req, user, db)
}

// misoapi-http: POST /open/api/password/remove-site-password
// misoapi-desc: Remove site password record
// misoapi-resource: ref(ResourceBasicUser)
func ApiRemoveSitePassword(rail miso.Rail, req vault.RemoveSitePasswordRes, user common.User, db *gorm.DB) (any, error) {
	return nil, vault.RemoveSitePassword(rail, req, user, db)
}

// misoapi-http: POST /open/api/password/decrypt-site-password
// misoapi-desc: Decrypt site password
// misoapi-resource: ref(ResourceBasicUser)
func ApiDecryptSitePassword(rail miso.Rail, req vault.DecryptSitePasswordReq, user common.User, db *gorm.DB) (vault.DecryptSitePasswordRes, error) {
	return vault.DecryptSitePassword(rail, req, user, db)
}

// misoapi-http: POST /open/api/password/edit-site-password
// misoapi-desc: Edit site password
// misoapi-resource: ref(ResourceBasicUser)
func ApiEditSitePassword(rail miso.Rail, req vault.EditSitePasswordReq, user common.User, db *gorm.DB) (any, error) {
	return nil, vault.EditSitePassword(rail, req, user, db)
}

// Clear user's failed login attempts
//
//   - misoapi-http: POST /open/api/user/clear-failed-login-attempts
//   - misoapi-desc: Admin clear user's failed login attempts
//   - misoapi-resource: ref(ResourceManagerUser)
func ApiClearUserFailedLoginAttempts(rail miso.Rail, req vault.ClearUserFailedLoginAttemptsReq) error {
	return vault.ClearFailedLoginAttempts(rail, req.UserNo)
}

// List User Notes
//
//   - misoapi-http: POST /open/api/note/list-notes
//   - misoapi-resource: ref(ResourceBasicUser)
//   - misoapi-ngtable
func ApiListNotes(rail miso.Rail, db *gorm.DB, req note.ListNoteReq, user common.User) (miso.PageRes[note.Note], error) {
	return note.ListNotes(rail, db, req, user)
}

// User Save Note
//
//   - misoapi-http: POST /open/api/note/save-note
//   - misoapi-resource: ref(ResourceBasicUser)
func ApiSaveNote(rail miso.Rail, db *gorm.DB, req note.SaveNoteReq, user common.User) error {
	return note.DBSaveNote(rail, db, req, user)
}

// User Update Note
//
//   - misoapi-http: POST /open/api/note/update-note
//   - misoapi-resource: ref(ResourceBasicUser)
func ApiUpdateNote(rail miso.Rail, db *gorm.DB, req note.UpdateNoteReq, user common.User) error {
	return note.UpdateNote(rail, db, req, user)
}

type ApiDeleteNoteReq struct {
	RecordId string
}

// User Delete Note
//
//   - misoapi-http: POST /open/api/note/delete-note
//   - misoapi-resource: ref(ResourceBasicUser)
func ApiDeleteNote(rail miso.Rail, db *gorm.DB, req ApiDeleteNoteReq, user common.User) error {
	return note.DeleteNote(rail, db, req.RecordId, user)
}
