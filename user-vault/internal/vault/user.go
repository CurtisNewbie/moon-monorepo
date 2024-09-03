package vault

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/curtisnewbie/miso/middleware/crypto"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"github.com/curtisnewbie/user-vault/api"
	"gorm.io/gorm"
)

var (
	usernameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_\-@.]{6,50}$`)
	passwordMinLen = 8

	userInfoCache = miso.NewRCache[UserDetail]("user-vault:user:info", miso.RCacheConfig{Exp: time.Hour * 1})
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
	CreateTime   util.ETime
	CreateBy     string
	UpdateTime   util.ETime
	UpdateBy     string
	IsDel        bool
}

func (u *User) Deleted() bool {
	return u.IsDel
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

func loadUser(rail miso.Rail, tx *gorm.DB, username string) (User, error) {
	if username == "" {
		return User{}, miso.NewErrf("Username is required")
	}

	var user User
	t := tx.Raw(`
		SELECT u.*, r.name AS role_name
		FROM user u
		LEFT JOIN role r using (role_no)
		WHERE u.username = ? and u.is_del = 0
	`, username).
		Scan(&user)

	if t.Error != nil {
		rail.Errorf("Failed to find user, username: %v, %v", username, t.Error)
		return User{}, t.Error
	}

	if t.RowsAffected < 1 {
		return User{}, miso.NewErrf("User not found").WithInternalMsg("User %v is not found", username)
	}

	return user, nil
}

func UserLogin(rail miso.Rail, tx *gorm.DB, req PasswordLoginParam) (string, User, error) {
	user, err := userLogin(rail, tx, req.Username, req.Password)
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

	return crypto.JwtEncode(claims, exp)
}

func userLogin(rail miso.Rail, tx *gorm.DB, username string, password string) (User, error) {
	if util.IsBlankStr(username) {
		return User{}, miso.NewErrf("Username is required")
	}

	if util.IsBlankStr(password) {
		return User{}, miso.NewErrf("Password is required")
	}

	user, err := loadUser(rail, tx, username)
	if err != nil {
		return User{}, err
	}

	if user.ReviewStatus == api.ReviewPending {
		return User{}, miso.NewErrf("Your registration is being reviewed, please wait for approval")
	}

	if user.ReviewStatus == api.ReviewRejected {
		return User{}, miso.NewErrf("Your are not permitted to login, please contact administrator")
	}

	if user.IsDisabled == api.UserDisabled {
		return User{}, miso.NewErrf("User is disabled")
	}

	if checkPassword(user.Password, user.Salt, password) {
		return user, nil
	}

	// if the password is incorrect, maybe a user_key is used instead
	ok, err := checkUserKey(rail, tx, user.UserNo, password)
	if err != nil {
		return User{}, err
	}
	if ok {
		return user, nil
	}

	return User{}, miso.NewErrf("Password incorrect").WithInternalMsg("User %v login failed, password incorrect", username)
}

func checkUserKey(rail miso.Rail, tx *gorm.DB, userNo string, password string) (bool, error) {
	if password == "" {
		return false, nil
	}

	var id int
	t := tx.Raw(
		`SELECT id FROM user_key WHERE user_no = ? AND secret_key = ? AND expiration_time > ? AND is_del = '0' LIMIT 1`,
		userNo, password, util.Now(),
	).
		Scan(&id)
	if t.Error != nil {
		rail.Errorf("failed to checkUserKey, userNo: %v, %v", userNo, t.Error)
	}
	return id > 0, nil
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
		return miso.NewErrf("Username must have 6-50 characters, permitted characters include: 'a-z A-Z 0-9 . - _ @'").
			WithInternalMsg("Actual username: %v", username)
	}
	return nil
}

func checkNewPassword(password string) error {
	len := len([]rune(password))
	if len < passwordMinLen {
		return miso.NewErrf("Password must have at least %v characters", passwordMinLen).
			WithInternalMsg("Actual length: %v", len)
	}
	return nil
}

type CreateUserParam struct {
	Username     string
	Password     string
	RoleNo       string
	ReviewStatus string
	Operator     string
}

func NewUser(rail miso.Rail, tx *gorm.DB, req CreateUserParam) error {
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
		return miso.NewErrf("Username and password must be different")
	}

	if _, err := loadUser(rail, tx, req.Username); err == nil {
		return miso.NewErrf("User is already registered")
	}

	user := prepUserCred(req.Password)
	user.UserNo = util.GenIdP("UE")
	user.Username = req.Username
	user.RoleNo = req.RoleNo
	user.CreateBy = req.Operator
	user.CreateTime = util.Now()
	user.IsDisabled = api.UserNormal
	user.ReviewStatus = req.ReviewStatus

	if err := tx.Table("user").Create(&user).Error; err != nil {
		rail.Errorf("failed to add new user '%v', %v", req.Username, err)
		return err
	}

	rail.Infof("New user '%v' with roleNo: %v is added by %v", req.Username, req.RoleNo, req.Operator)
	return nil
}

type NewUserParam struct {
	Id           int
	UserNo       string
	Username     string
	Password     string
	Salt         string
	ReviewStatus string
	RoleNo       string
	IsDisabled   int
	CreateTime   util.ETime
	CreateBy     string
	UpdateTime   util.ETime
	UpdateBy     string
	IsDel        bool
}

func prepUserCred(pwd string) NewUserParam {
	u := NewUserParam{}
	u.Salt = util.RandStr(6)
	u.Password = encodePasswordSalt(pwd, u.Salt)
	return u
}

func ListUsers(rail miso.Rail, tx *gorm.DB, req ListUserReq) (miso.PageRes[api.UserInfo], error) {
	return miso.NewPageQuery[api.UserInfo]().
		WithPage(req.Paging).
		WithSelectQuery(func(tx *gorm.DB) *gorm.DB {
			return tx.Select("u.*, r.name as role_name").Order("u.id DESC")
		}).
		WithBaseQuery(func(tx *gorm.DB) *gorm.DB {
			tx = tx.Table("user u").Joins("LEFT JOIN role r USING(role_no)")

			if req.RoleNo != nil && *req.RoleNo != "" {
				tx = tx.Where("u.role_no = ?", *req.RoleNo)
			}
			if req.Username != nil && *req.Username != "" {
				tx = tx.Where("u.username LIKE ?", "%"+*req.Username+"%")
			}
			if req.IsDisabled != nil {
				tx = tx.Where("u.is_disabled = ?", *req.IsDisabled)
			}
			return tx.Where("u.is_del = 0")
		}).
		Exec(rail, tx)
}

func AdminUpdateUser(rail miso.Rail, tx *gorm.DB, req AdminUpdateUserReq, operator common.User) error {
	if operator.UserNo == req.UserNo {
		return miso.NewErrf("You cannot update yourself")
	}

	if req.RoleNo != "" {
		_, err := GetRoleInfo(rail, api.RoleInfoReq{RoleNo: req.RoleNo})
		if err != nil {
			return miso.NewErrf("Invalid role").WithInternalMsg("failed to get role info, roleNo may be invalid, %v", err)
		}
	}

	return tx.Exec(
		`UPDATE user SET is_disabled = ?, update_by = ?, role_no = ? WHERE user_no = ?`,
		req.IsDisabled, operator.Username, req.RoleNo, req.UserNo,
	).Error
}

func ReviewUserRegistration(rail miso.Rail, tx *gorm.DB, req AdminReviewUserReq) error {
	if req.ReviewStatus != api.ReviewRejected && req.ReviewStatus != api.ReviewApproved {
		return miso.NewErrf("Illegal Argument").
			WithInternalMsg("ReviewStatus was neither ReviewApproved nor ReviewRejected, it was %v", req.ReviewStatus)
	}

	return miso.RLockExec(rail, fmt.Sprintf("auth:user:registration:review:%v", req.UserId),
		func() error {
			var user User
			t := tx.Raw(`SELECT * FROM user WHERE id = ?`, req.UserId).
				Scan(&user)
			if t.Error != nil {
				rail.Errorf("Failed to find user, id = %v %v", req.UserId, t.Error)
				return t.Error
			}

			if t.RowsAffected < 1 {
				return miso.NewErrf("User not found").WithInternalMsg("User %v not found", req.UserId)
			}

			if user.Deleted() {
				return miso.NewErrf("User not found").WithInternalMsg("User %v is deleted", req.UserId)
			}

			if !user.CanReview() {
				return miso.NewErrf("User's registration has already been reviewed")
			}

			isDisabled := api.UserDisabled
			if req.ReviewStatus == api.ReviewApproved {
				isDisabled = api.UserNormal
			}

			err := tx.Exec(`UPDATE user SET review_status = ?, is_disabled = ? WHERE id = ?`, req.ReviewStatus, isDisabled, req.UserId).
				Error

			if err != nil {
				rail.Errorf("failed to update user for registration review, userId: %v, %v", req.UserId, err)
			}
			return err
		},
	)
}

func UserRegister(rail miso.Rail, db *gorm.DB, req RegisterReq) error {
	if err := NewUser(rail, db, CreateUserParam{
		Username:     req.Username,
		Password:     req.Password,
		ReviewStatus: api.ReviewPending,
	}); err != nil {
		return err
	}

	err := api.CreateNotifiByAccessPipeline.Send(rail, api.CreateNotifiByAccessEvent{
		Title:   fmt.Sprintf("Review user %v's registration", req.Username),
		Message: fmt.Sprintf("Please review new user %v's registration. A role should be assigned for the new user.", req.Username),
		ResCode: ResourceManagerUser,
	})
	if err != nil {
		rail.Errorf("failed to create notification for UserRegister, %v", err)
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

func FetchUserBrief(rail miso.Rail, tx *gorm.DB, username string) (UserInfoBrief, error) {
	ud, err := LoadUserBriefThrCache(rail, miso.GetMySQL(), username)
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

func LoadUserBriefThrCache(rail miso.Rail, tx *gorm.DB, username string) (UserDetail, error) {
	rail.Debugf("LoadUserBriefThrCache, username: %v", username)
	return userInfoCache.Get(rail, username, func() (UserDetail, error) {
		rail.Debugf("LoadUserInfoBrief, username: %v", username)
		return LoadUserInfoBrief(rail, miso.GetMySQL(), username)
	})
}

func InvalidateUserInfoCache(rail miso.Rail, username string) error {
	return userInfoCache.Del(rail, username)
}

func LoadUserInfoBrief(rail miso.Rail, tx *gorm.DB, username string) (UserDetail, error) {
	u, err := loadUser(rail, tx, username)
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

func UpdatePassword(rail miso.Rail, tx *gorm.DB, username string, req UpdatePasswordReq) error {
	req.NewPassword = strings.TrimSpace(req.NewPassword)
	req.PrevPassword = strings.TrimSpace(req.PrevPassword)

	if req.NewPassword == req.PrevPassword {
		return miso.NewErrf("New password must be different")
	}

	if err := checkNewPassword(req.NewPassword); err != nil {
		return err
	}

	if username == req.NewPassword {
		return miso.NewErrf("Username and password must be different")
	}

	u, err := LoadUserBriefThrCache(rail, tx, username)
	if err != nil {
		return miso.NewErrf("Failed to load user info, please try again later").
			WithInternalMsg("Failed to LoadUserBriefThrCache, %v", err)
	}

	if !checkPassword(u.Password, u.Salt, req.PrevPassword) {
		return miso.NewErrf("Password incorrect")
	}

	t := tx.Exec("update user set password = ? where username = ?", encodePasswordSalt(req.NewPassword, u.Salt), username)
	if t.Error != nil {
		return miso.NewErrf("Failed to update password, please try again laster").
			WithInternalMsg("Failed to update password, %v", t.Error)
	}
	return nil
}

type ExchangeTokenReq struct {
	Token string `json:"token" valid:"notEmpty"`
}

func DecodeTokenUser(rail miso.Rail, token string) (TokenUser, error) {
	tu := TokenUser{}
	decoded, err := crypto.JwtDecode(token)
	if err != nil || !decoded.Valid {
		return TokenUser{}, miso.NewErrf("Illegal token").WithInternalMsg("Failed to decode jwt token, %v", err)
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
	decoded, err := crypto.JwtDecode(token)
	if err != nil || !decoded.Valid {
		return "", miso.NewErrf("Illegal token").WithInternalMsg("Failed to decode jwt token, %v", err)
	}
	username := decoded.Claims["username"]
	un, ok := username.(string)
	if !ok {
		un = fmt.Sprintf("%v", username)
	}
	return un, nil
}

func ExchangeToken(rail miso.Rail, tx *gorm.DB, req ExchangeTokenReq) (string, error) {
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

func GetTokenUser(rail miso.Rail, tx *gorm.DB, token string) (UserInfoBrief, error) {
	if util.IsBlankStr(token) {
		return UserInfoBrief{}, miso.NewErrf("Invalid token").WithInternalMsg("Token is blank")
	}
	username, err := DecodeTokenUsername(rail, token)
	if err != nil {
		return UserInfoBrief{}, err
	}

	u, err := LoadUserBriefThrCache(rail, tx, username)

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

func ItnFindUserInfo(rail miso.Rail, tx *gorm.DB, req api.FindUserReq) (api.UserInfo, error) {

	var ui api.UserInfo
	tx = tx.Table("user").
		Joins("left join role on user.role_no = role.role_no").
		Select("user.*, role.name role_name")

	if req.UserId == nil && req.UserNo == nil && req.Username == nil {
		return ui, miso.NewErrf("Must provide at least one parameter")
	}

	if req.UserId != nil {
		tx = tx.Where("user.id = ?", *req.UserId)
	}
	if req.UserNo != nil {
		tx = tx.Where("user.user_no = ?", *req.UserNo)
	}
	if req.Username != nil {
		tx = tx.Where("user.username = ?", *req.Username)
	}

	t := tx.Scan(&ui)
	if t.Error != nil {
		return ui, fmt.Errorf("failed to find user %w", t.Error)
	}
	if t.RowsAffected < 1 {
		return ui, miso.NewErrf("User not found")
	}
	return ui, nil
}

func ItnFindNameOfUserNo(rail miso.Rail, tx *gorm.DB, req api.FetchNameByUserNoReq) (api.FetchUsernamesRes, error) {
	if len(req.UserNos) < 1 {
		return api.FetchUsernamesRes{UserNoToUsername: map[string]string{}}, nil
	}

	type UserNoToName struct {
		UserNo   string
		Username string
	}

	var queried []UserNoToName
	err := tx.Table("user").
		Select("username", "user_no").
		Where("user_no in ?", util.Distinct(req.UserNos)).
		Scan(&queried).
		Error
	if err != nil {
		return api.FetchUsernamesRes{}, err
	}

	mapping := util.StrMap(queried,
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
	err := db.Table("user").
		Where("role_no = ?", req.RoleNo).
		Scan(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list users with roleNo: %v, %w", req.RoleNo, err)
	}
	return users, nil
}

func FindUserWithRes(rail miso.Rail, db *gorm.DB, req api.FetchUserWithResourceReq) ([]api.UserInfo, error) {
	var users []api.UserInfo
	err := db.Raw(`
		select u.*, r.name role_name from user u
		left join role r on u.role_no = r.role_no
		left join role_resource rr on r.role_no = rr.role_no
		where rr.res_code = ? or r.role_no = ?`, req.ResourceCode, DefaultAdminRoleNo).
		Scan(&users).
		Error
	return users, err
}
