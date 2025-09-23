package vault

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/jwt"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"github.com/curtisnewbie/miso/util/errs"
	"github.com/curtisnewbie/miso/util/hash"
	"github.com/curtisnewbie/miso/util/slutil"
	"github.com/curtisnewbie/miso/util/strutil"
	"github.com/curtisnewbie/user-vault/api"
	"github.com/curtisnewbie/user-vault/internal/config"
	"gorm.io/gorm"
)

var (
	usernameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_\-@.]{6,50}$`)
	passwordMinLen = 8

	userInfoCache = redis.NewRCache[UserDetail]("user-vault:user:info", redis.RCacheConfig{Exp: time.Hour * 1})
)

const (
	maxFailedLoginAttempts     = 15
	failedLoginAttemptRedisKey = "user-vault:user:login-failed-count:"
)

type PasswordLoginParam struct {
	Username string
	Password string
}

type AddUserParam struct {
	Username string `json:"username" valid:"notEmpty"`
	Password string `json:"password" valid:"notEmpty"`
	RoleNo   string `json:"roleNo"`
}

type User struct {
	Id           int
	UserNo       string
	Username     string
	Password     string
	Salt         string
	ReviewStatus string
	RoleNo       string
	RoleName     string
	IsDisabled   int
	CreateTime   util.ETime `gorm:"column:created_at"`
	CreateBy     string     `gorm:"column:created_by"`
	UpdateTime   util.ETime `gorm:"column:updated_at"`
	UpdateBy     string     `gorm:"column:updated_by"`
	Deleted      bool
}

func (u *User) CanReview() bool {
	return u.ReviewStatus == api.ReviewPending
}

type UserDetail struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	RoleName     string `json:"roleName"`
	RoleNo       string `json:"roleNo"`
	UserNo       string `json:"userNo"`
	RegisterDate string `json:"registerDate"`
	Password     string `json:"password"`
	Salt         string `json:"salt"`
}

func loadUser(rail miso.Rail, db *gorm.DB, username string) (User, error) {
	if username == "" {
		return User{}, errs.NewErrf("Username is required")
	}

	var user User
	n, err := dbquery.NewQuery(rail, db).Raw(`
		SELECT u.*, r.name AS role_name
		FROM user u
		LEFT JOIN role r using (role_no)
		WHERE u.username = ? and u.deleted = 0
	`, username).
		Scan(&user)

	if err != nil {
		rail.Errorf("Failed to find user, username: %v, %v", username, err)
		return User{}, err
	}

	if n < 1 {
		return User{}, errs.NewErrf("User not found").WithInternalMsg("User %v is not found", username)
	}

	return user, nil
}

func UserLogin(rail miso.Rail, db *gorm.DB, req PasswordLoginParam) (string, User, error) {
	user, err := userLogin(rail, db, req.Username, req.Password)
	if err != nil {
		return "", User{}, err
	}

	tu := TokenUser{
		Id:       user.Id,
		UserNo:   user.UserNo,
		Username: user.Username,
		RoleNo:   user.RoleNo,
	}

	rail.Debugf("buildToken %+v", tu)
	tkn, err := buildToken(tu, 15*time.Minute)
	if err != nil {
		return "", User{}, err
	}
	return tkn, user, nil
}

type TokenUser struct {
	Id       int
	UserNo   string
	Username string
	RoleNo   string
}

func buildToken(user TokenUser, exp time.Duration) (string, error) {
	claims := map[string]any{
		"id":       user.Id,
		"username": user.Username,
		"userno":   user.UserNo,
		"roleno":   user.RoleNo,
	}

	return jwt.JwtEncode(claims, exp)
}

func userLogin(rail miso.Rail, db *gorm.DB, username string, password string) (User, error) {
	if strutil.IsBlankStr(username) {
		return User{}, errs.NewErrf("Username is required")
	}

	if strutil.IsBlankStr(password) {
		return User{}, errs.NewErrf("Password is required")
	}

	user, err := loadUser(rail, db, username)
	if err != nil {
		return User{}, err
	}

	if user.ReviewStatus == api.ReviewPending {
		return User{}, errs.NewErrf("Your registration is being reviewed, please wait for approval")
	}

	if user.ReviewStatus == api.ReviewRejected {
		return User{}, errs.NewErrf("Your are not permitted to login, please contact administrator")
	}

	if user.IsDisabled == api.UserDisabled {
		return User{}, errs.NewErrf("User is disabled")
	}

	{
		ok, er := CheckFailedLoginAttempts(rail, user.UserNo)
		if er != nil {
			rail.Errorf("Failed to check user's failed login attempts, userNo: %v, %v", user.UserNo, er)
		} else if !ok {
			rail.Infof("User's failed login attempts exceeded limit, userNo: %v, reject login request", user.UserNo)
			return User{}, errs.NewErrf("Exceeded maximum login attempts, please try again later.")
		}
	}

	if checkPassword(user.Password, user.Salt, password) {
		return user, nil
	}

	// if the password is incorrect, maybe a user_key is used instead
	ok, err := checkUserKey(rail, db, user.UserNo, password)
	if err != nil {
		return User{}, err
	}
	if ok {
		return user, nil
	}

	if er := IncrFailedLoginAttempts(rail, user.UserNo); er != nil {
		rail.Warnf("Failed to update user's failed login attempt, userNo: %v, %v", user.UserNo, er)
	}

	return User{}, errs.NewErrf("Password incorrect").WithInternalMsg("User %v login failed, password incorrect", username)
}

func checkUserKey(rail miso.Rail, db *gorm.DB, userNo string, password string) (bool, error) {
	if password == "" {
		return false, nil
	}

	var id int
	n, err := dbquery.NewQuery(rail, db).Raw(
		`SELECT id FROM user_key WHERE user_no = ? AND secret_key = ? AND expiration_time > ? AND deleted = 0 LIMIT 1`,
		userNo, password, util.Now(),
	).Scan(&id)

	if err != nil {
		rail.Errorf("failed to checkUserKey, userNo: %v, %v", userNo, err)
	}
	return n > 0, nil
}

func checkPassword(encoded string, salt string, password string) bool {
	if password == "" {
		return false
	}
	springSalt := extractSpringSalt(encoded) // for backward compatibility (auth-service)
	ep := encodePasswordSalt(password, salt)
	provided := springSalt + ep
	return provided == encoded
}

func encodePasswordSalt(pwd string, salt string) string {
	return encodePassword(pwd + salt)
}

func encodePassword(text string) string {
	sha := sha256.New()
	sha.Write([]byte(text))
	return fmt.Sprintf("%x", sha.Sum(nil))
}

// for backward compatibility, we are still using the schema used by auth-service
func extractSpringSalt(encoded string) string {
	ru := []rune(encoded)
	if len(ru) < 1 {
		return ""
	}

	if ru[0] != '{' {
		return "" // none
	}

	for i := range ru {
		if ru[i] == '}' { // end of the embedded salt
			return string(ru[0 : i+1])
		}
	}

	return "" // illegal format, or maybe none
}

func checkNewUsername(username string) error {
	if !usernameRegexp.MatchString(username) {
		return errs.NewErrf("Username must have 6-50 characters, permitted characters include: 'a-z A-Z 0-9 . - _ @'").
			WithInternalMsg("Actual username: %v", username)
	}
	return nil
}

func checkNewPassword(password string) error {
	len := len([]rune(password))
	if len < passwordMinLen {
		return errs.NewErrf("Password must have at least %v characters", passwordMinLen).
			WithInternalMsg("Actual length: %v", len)
	}
	return nil
}

type CreateUserParam struct {
	Username     string
	Password     string
	RoleNo       string
	ReviewStatus string
}

func NewUser(rail miso.Rail, db *gorm.DB, req CreateUserParam) error {
	if req.RoleNo != "" {
		_, err := GetRoleInfo(rail, api.RoleInfoReq{RoleNo: req.RoleNo})
		if err != nil {
			return err
		}
	}

	if e := checkNewUsername(req.Username); e != nil {
		return e
	}

	if e := checkNewPassword(req.Password); e != nil {
		return e
	}

	if req.Username == req.Password {
		return errs.NewErrf("Username and password must be different")
	}

	if _, err := loadUser(rail, db, req.Username); err == nil {
		return errs.NewErrf("User is already registered")
	}

	user := prepUserCred(req.Password)
	user.UserNo = util.GenIdP("UE")
	user.Username = req.Username
	user.RoleNo = req.RoleNo
	user.IsDisabled = api.UserNormal
	user.ReviewStatus = req.ReviewStatus

	if err := dbquery.NewQuery(rail, db).Table("user").CreateAny(&user); err != nil {
		rail.Errorf("Failed to add new user '%v', %v", req.Username, err)
		return err
	}

	rail.Infof("New user '%v' with roleNo: %v is created", req.Username, req.RoleNo)
	return nil
}

type NewUserParam struct {
	UserNo       string
	Username     string
	Password     string
	Salt         string
	ReviewStatus string
	RoleNo       string
	IsDisabled   int
}

func prepUserCred(pwd string) NewUserParam {
	u := NewUserParam{}
	u.Salt = util.RandStr(6)
	u.Password = encodePasswordSalt(pwd, u.Salt)
	return u
}

type ListUserReq struct {
	Username   *string     `json:"username"`
	RoleNo     *string     `json:"roleNo"`
	IsDisabled *int        `json:"isDisabled"`
	Paging     miso.Paging `json:"paging"`
}

func ListUsers(rail miso.Rail, db *gorm.DB, req ListUserReq) (miso.PageRes[api.UserInfo], error) {
	return dbquery.NewPagedQuery[api.UserInfo](db).
		WithSelectQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.Select("u.*, r.name as role_name").Order("u.id DESC")
		}).
		WithBaseQuery(func(q *dbquery.Query) *dbquery.Query {
			q = q.Table("user u").Joins("LEFT JOIN role r USING(role_no)")

			if req.RoleNo != nil && *req.RoleNo != "" {
				q = q.Eq("u.role_no", *req.RoleNo)
			}
			if req.Username != nil && *req.Username != "" {
				q = q.Like("u.username", *req.Username)
			}
			if req.IsDisabled != nil {
				q = q.Eq("u.is_disabled", *req.IsDisabled)
			}
			return q.Where("u.deleted = 0")
		}).
		Scan(rail, req.Paging)
}

type AdminUpdateUserReq struct {
	UserNo     string `valid:"notEmpty"`
	RoleNo     string `json:"roleNo"`
	IsDisabled int    `json:"isDisabled"`
}

func AdminUpdateUser(rail miso.Rail, db *gorm.DB, req AdminUpdateUserReq, operator common.User) error {
	if operator.UserNo == req.UserNo {
		return errs.NewErrf("You cannot update yourself")
	}

	if req.RoleNo != "" {
		_, err := GetRoleInfo(rail, api.RoleInfoReq{RoleNo: req.RoleNo})
		if err != nil {
			return errs.NewErrf("Invalid role").WithInternalMsg("failed to get role info, roleNo may be invalid, %v", err)
		}
	}

	return dbquery.NewQuery(rail, db).
		Table("user").
		SetCols(struct {
			RoleNo     string
			IsDisabled int
		}{
			IsDisabled: req.IsDisabled,
			RoleNo:     req.RoleNo,
		}).
		Eq("user_no", req.UserNo).
		UpdateAny()
}

type AdminReviewUserReq struct {
	UserId       int    `json:"userId" valid:"positive"`
	ReviewStatus string `json:"reviewStatus"`
}

func ReviewUserRegistration(rail miso.Rail, db *gorm.DB, req AdminReviewUserReq) error {
	if req.ReviewStatus != api.ReviewRejected && req.ReviewStatus != api.ReviewApproved {
		return errs.NewErrf("Illegal Argument").
			WithInternalMsg("ReviewStatus was neither ReviewApproved nor ReviewRejected, it was %v", req.ReviewStatus)
	}

	return redis.RLockExec(rail, fmt.Sprintf("auth:user:registration:review:%v", req.UserId),
		func() error {
			var user User
			n, err := dbquery.NewQuery(rail, db).
				Raw(`SELECT * FROM user WHERE id = ?`, req.UserId).
				Scan(&user)
			if err != nil {
				rail.Errorf("Failed to find user, id = %v %v", req.UserId, err)
				return err
			}
			if n < 1 {
				return errs.NewErrf("User not found").WithInternalMsg("User %v not found", req.UserId)
			}

			if user.Deleted {
				return errs.NewErrf("User not found").WithInternalMsg("User %v is deleted", req.UserId)
			}

			if !user.CanReview() {
				return errs.NewErrf("User's registration has already been reviewed")
			}

			var roleNo string
			isDisabled := api.UserDisabled
			if req.ReviewStatus == api.ReviewApproved {
				isDisabled = api.UserNormal
				rail.Infof("User role: %v", user.RoleNo)

				if user.RoleNo == "" {
					dr := config.DefaultUserRole()
					rail.Infof("Default role: %v", dr)
					if dr != "" {
						roleNo = dr
					}
				}
			}

			err = dbquery.NewQuery(rail, db).
				Table("user").
				Set("review_status", req.ReviewStatus).
				Set("is_disabled", isDisabled).
				SetIf(roleNo != "", "role_no", roleNo).
				Eq("id", req.UserId).
				UpdateAny()

			rail.ErrorIf(err, "Failed to update user for registration review, userId: %v", req.UserId)
			return err
		},
	)
}

type RegisterReq struct {
	Username string `json:"username" valid:"notEmpty"`
	Password string `json:"password" valid:"notEmpty"`
}

func UserRegister(rail miso.Rail, db *gorm.DB, req RegisterReq) error {
	if err := NewUser(rail, db, CreateUserParam{
		Username:     req.Username,
		Password:     req.Password,
		ReviewStatus: api.ReviewPending,
	}); err != nil {
		return err
	}
	return nil
}

type UserInfoBrief struct {
	Id           int    `json:"id"`
	Username     string `json:"username"`
	RoleName     string `json:"roleName"`
	RoleNo       string `json:"roleNo"`
	UserNo       string `json:"userNo"`
	RegisterDate string `json:"registerDate"`
}

func FetchUserBrief(rail miso.Rail, db *gorm.DB, username string) (UserInfoBrief, error) {
	ud, err := LoadUserBriefThrCache(rail, db, username)
	if err != nil {
		return UserInfoBrief{}, err
	}
	return UserInfoBrief{
		Id:           ud.Id,
		Username:     ud.Username,
		RoleName:     ud.RoleName,
		RoleNo:       ud.RoleNo,
		UserNo:       ud.UserNo,
		RegisterDate: ud.RegisterDate,
	}, nil
}

func LoadUserBriefThrCache(rail miso.Rail, db *gorm.DB, username string) (UserDetail, error) {
	rail.Debugf("LoadUserBriefThrCache, username: %v", username)
	return userInfoCache.GetValElse(rail, username, func() (UserDetail, error) {
		rail.Debugf("LoadUserInfoBrief, username: %v", username)
		return LoadUserInfoBrief(rail, db, username)
	})
}

func InvalidateUserInfoCache(rail miso.Rail, username string) error {
	return userInfoCache.Del(rail, username)
}

func LoadUserInfoBrief(rail miso.Rail, db *gorm.DB, username string) (UserDetail, error) {
	u, err := loadUser(rail, db, username)
	if err != nil {
		return UserDetail{}, err
	}

	return UserDetail{
		Id:           u.Id,
		Username:     u.Username,
		RoleName:     u.RoleName,
		RoleNo:       u.RoleNo,
		UserNo:       u.UserNo,
		RegisterDate: u.CreateTime.FormatClassic(),
		Salt:         u.Salt,
		Password:     u.Password,
	}, nil
}

type UpdatePasswordReq struct {
	PrevPassword string `json:"prevPassword" valid:"notEmpty"`
	NewPassword  string `json:"newPassword" valid:"notEmpty"`
}

func UpdatePassword(rail miso.Rail, db *gorm.DB, username string, req UpdatePasswordReq) error {
	req.NewPassword = strings.TrimSpace(req.NewPassword)
	req.PrevPassword = strings.TrimSpace(req.PrevPassword)

	if req.NewPassword == req.PrevPassword {
		return errs.NewErrf("New password must be different")
	}

	if err := checkNewPassword(req.NewPassword); err != nil {
		return err
	}

	if username == req.NewPassword {
		return errs.NewErrf("Username and password must be different")
	}

	u, err := LoadUserBriefThrCache(rail, db, username)
	if err != nil {
		return errs.NewErrf("Failed to load user info, please try again later").
			WithInternalMsg("Failed to LoadUserBriefThrCache, %v", err)
	}

	if !checkPassword(u.Password, u.Salt, req.PrevPassword) {
		return errs.NewErrf("Password incorrect")
	}

	err = dbquery.NewQuery(rail, db).
		Table("user").
		Set("password", encodePasswordSalt(req.NewPassword, u.Salt)).
		Eq("username", username).
		UpdateAny()
	if err != nil {
		return errs.NewErrf("Failed to update password, please try again laster").
			WithInternalMsg("Failed to update password, %v", err)
	}
	return nil
}

type ExchangeTokenReq struct {
	Token string `json:"token" valid:"notEmpty"`
}

func DecodeTokenUser(rail miso.Rail, token string) (TokenUser, error) {
	tu := TokenUser{}
	decoded, err := jwt.JwtDecode(token)
	if err != nil || !decoded.Valid {
		return TokenUser{}, errs.NewErrf("Illegal token").WithInternalMsg("Failed to decode jwt token, %v", err)
	}

	tu.Id, err = strconv.Atoi(fmt.Sprintf("%v", decoded.Claims["id"]))
	if err != nil {
		return tu, err
	}
	tu.Username = decoded.Claims["username"].(string)
	tu.UserNo = decoded.Claims["userno"].(string)
	tu.RoleNo = decoded.Claims["roleno"].(string)
	return tu, nil
}

func DecodeTokenUsername(rail miso.Rail, token string) (string, error) {
	decoded, err := jwt.JwtDecode(token)
	if err != nil || !decoded.Valid {
		return "", errs.NewErrf("Illegal token").WithInternalMsg("Failed to decode jwt token, %v", err)
	}
	username := decoded.Claims["username"]
	un, ok := username.(string)
	if !ok {
		un = fmt.Sprintf("%v", username)
	}
	return un, nil
}

func ExchangeToken(rail miso.Rail, req ExchangeTokenReq) (string, error) {
	u, err := DecodeTokenUser(rail, req.Token)
	if err != nil {
		return "", err
	}

	tu := TokenUser{
		Id:       u.Id,
		UserNo:   u.UserNo,
		Username: u.Username,
		RoleNo:   u.RoleNo,
	}

	rail.Debugf("buildToken %+v", tu)
	return buildToken(tu, 15*time.Minute)
}

func GetTokenUser(rail miso.Rail, db *gorm.DB, token string) (UserInfoBrief, error) {
	if strutil.IsBlankStr(token) {
		return UserInfoBrief{}, errs.NewErrf("Invalid token").WithInternalMsg("Token is blank")
	}
	username, err := DecodeTokenUsername(rail, token)
	if err != nil {
		return UserInfoBrief{}, err
	}

	u, err := LoadUserBriefThrCache(rail, db, username)
	if err != nil {
		return UserInfoBrief{}, err
	}
	return UserInfoBrief{
		Id:           u.Id,
		Username:     u.Username,
		RoleName:     u.RoleName,
		RoleNo:       u.RoleNo,
		UserNo:       u.UserNo,
		RegisterDate: u.RegisterDate,
	}, nil
}

func ItnFindUserInfo(rail miso.Rail, db *gorm.DB, req api.FindUserReq) (api.UserInfo, error) {

	var ui api.UserInfo
	q := dbquery.NewQuery(rail, db).
		Table("user").
		Joins("left join role on user.role_no = role.role_no").
		Select("user.*, role.name role_name")

	if req.UserId == nil && req.UserNo == nil && req.Username == nil {
		return ui, errs.NewErrf("Must provide at least one parameter")
	}

	if req.UserId != nil {
		q = q.Where("user.id = ?", *req.UserId)
	}
	if req.UserNo != nil {
		q = q.Where("user.user_no = ?", *req.UserNo)
	}
	if req.Username != nil {
		q = q.Where("user.username = ?", *req.Username)
	}

	n, err := q.Scan(&ui)
	if err != nil {
		return ui, fmt.Errorf("failed to find user %w", err)
	}
	if n < 1 {
		return ui, errs.NewErrf("User not found")
	}
	return ui, nil
}

func ItnFindNameOfUserNo(rail miso.Rail, db *gorm.DB, req api.FetchNameByUserNoReq) (api.FetchUsernamesRes, error) {
	if len(req.UserNos) < 1 {
		return api.FetchUsernamesRes{UserNoToUsername: map[string]string{}}, nil
	}

	type UserNoToName struct {
		UserNo   string
		Username string
	}

	var queried []UserNoToName
	err := dbquery.NewQuery(rail, db).
		Table("user").
		Select("username", "user_no").
		Where("user_no in ?", slutil.Distinct(req.UserNos)).
		ScanVal(&queried)
	if err != nil {
		return api.FetchUsernamesRes{}, err
	}

	mapping := hash.StrMap(queried,
		func(un UserNoToName) string {
			return un.UserNo
		},
		func(un UserNoToName) string {
			return un.Username
		},
	)
	return api.FetchUsernamesRes{UserNoToUsername: mapping}, nil
}

func ItnFindUsersWithRole(rail miso.Rail, db *gorm.DB, req api.FetchUsersWithRoleReq) ([]api.UserInfo, error) {
	var users []api.UserInfo
	_, err := dbquery.NewQuery(rail, db).
		Table("user").
		Where("role_no = ?", req.RoleNo).
		Scan(&users)
	if err != nil {
		return nil, fmt.Errorf("failed to list users with roleNo: %v, %w", req.RoleNo, err)
	}
	return users, nil
}

func FindUserWithRes(rail miso.Rail, db *gorm.DB, req api.FetchUserWithResourceReq) ([]api.UserInfo, error) {
	var users []api.UserInfo
	_, err := dbquery.NewQuery(rail, db).Raw(`
		select u.*, r.name role_name from user u
		left join role r on u.role_no = r.role_no
		left join role_resource rr on r.role_no = rr.role_no
		where rr.res_code = ? or r.role_no in ?`, req.ResourceCode, []string{DefaultAdminRoleNo, DefaultAdminRoleNo2}).
		Scan(&users)
	return users, err
}

type ClearUserFailedLoginAttemptsReq struct {
	UserNo string
}

func ClearFailedLoginAttempts(rail miso.Rail, userNo string) error {
	r := redis.GetRedis()
	k := failedLoginAttemptRedisKey + userNo
	if err := r.Del(rail.Context(), k).Err(); err != nil {
		return err
	}
	rail.Infof("Reset user %v failed login attempts", userNo)
	return nil
}

func IncrFailedLoginAttempts(rail miso.Rail, userNo string) error {
	r := redis.GetRedis()
	k := failedLoginAttemptRedisKey + userNo
	c := r.Incr(rail.Context(), k)
	if c.Err() != nil {
		return c.Err()
	}
	rail.Infof("User %v login failed, curr failed attempts: %v", userNo, c.Val())
	return r.Expire(rail.Context(), k, time.Minute*15).Err()
}

func CheckFailedLoginAttempts(rail miso.Rail, userNo string) (bool, error) {
	c := redis.GetRedis().Get(rail.Context(), failedLoginAttemptRedisKey+userNo)
	if c.Err() != nil {
		if redis.IsNil(c.Err()) {
			return true, nil
		}
		return false, c.Err()
	}
	n, err := c.Int()
	if err != nil {
		return false, err
	}
	return n < maxFailedLoginAttempts, nil
}
