package vault

import (
	"crypto/md5"
	"encoding/base64"
	"strings"
	"time"

	doublestar "github.com/bmatcuk/doublestar/v4"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/auth"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"github.com/curtisnewbie/miso/util/errs"
	"github.com/curtisnewbie/miso/util/slutil"
	"github.com/curtisnewbie/user-vault/api"
	"gorm.io/gorm"
)

var (
	permitted = api.CheckResAccessResp{Valid: true}
	forbidden = api.CheckResAccessResp{Valid: false}

	roleInfoCache = redis.NewRCache[api.RoleInfoResp]("user-vault:role:info",
		redis.RCacheConfig{Exp: 10 * time.Minute, NoSync: true})

	// cache for role's accessible resources and api url patterns
	roleAccessCache = redis.NewRCache[RoleAccess]("user-vault:role:access",
		redis.RCacheConfig{Exp: 6 * time.Hour, NoSync: true})

	// cache for publicly accessible resources and api url patterns
	publicAccessCache = redis.NewRCache[[]PathAccessInfo]("user-vault:public:access",
		redis.RCacheConfig{Exp: 6 * time.Hour, NoSync: true})
)

const (
	// default roleno for admin
	DefaultAdminRoleNo  = "role_554107924873216177918" // deprecated
	DefaultAdminRoleNo2 = "role_super_admin"

	PathTypeProtected string = "PROTECTED"
	PathTypePublic    string = "PUBLIC"
)

type ExtendedPathRes struct {
	Id         int        // id
	Pgroup     string     // path group
	PathNo     string     // path no
	ResCode    string     // resource code
	Desc       string     // description
	Url        string     // url
	Method     string     // http method
	Ptype      string     // path type: PROTECTED, PUBLIC
	CreateTime util.ETime `gorm:"column:created_at"`
	CreateBy   string     `gorm:"column:created_by"`
	UpdateTime util.ETime `gorm:"column:updated_at"`
	UpdateBy   string     `gorm:"column:updated_by"`
}

type WRole struct {
	Id         int        `json:"id"`
	RoleNo     string     `json:"roleNo"`
	Name       string     `json:"name"`
	CreateTime util.ETime `json:"createTime" gorm:"column:created_at"`
	CreateBy   string     `json:"createBy" gorm:"column:created_by"`
	UpdateTime util.ETime `json:"updateTime" gorm:"column:updated_at"`
	UpdateBy   string     `json:"updateBy" gorm:"column:updated_by"`
}

type ResBrief struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type AddRoleReq struct {
	Name string `json:"name" valid:"notEmpty,maxLen:32"` // role name
}

type ListRoleReq struct {
	Paging miso.Paging `json:"paging"`
}

type ListRoleResp struct {
	Payload []WRole     `json:"payload"`
	Paging  miso.Paging `json:"paging"`
}

type RoleBrief struct {
	RoleNo string `json:"roleNo"`
	Name   string `json:"name"`
}

type ListPathReq struct {
	ResCode string      `json:"resCode"`
	Pgroup  string      `json:"pgroup"`
	Url     string      `json:"url"`
	Ptype   string      `json:"ptype" desc:"path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible"`
	Paging  miso.Paging `json:"paging"`
}

type WPath struct {
	Id         int        `json:"id"`
	Pgroup     string     `json:"pgroup"`
	PathNo     string     `json:"pathNo"`
	Method     string     `json:"method"`
	Desc       string     `json:"desc"`
	Url        string     `json:"url"`
	Ptype      string     `json:"ptype" desc:"path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible"`
	CreateTime util.ETime `json:"createTime" gorm:"column:created_at"`
	CreateBy   string     `json:"createBy" gorm:"column:created_by"`
	UpdateTime util.ETime `json:"updateTime" gorm:"column:updated_at"`
	UpdateBy   string     `json:"updateBy" gorm:"column:updated_by"`
}

type WRes struct {
	Id         int        `json:"id"`
	Code       string     `json:"code"`
	Name       string     `json:"name"`
	CreateTime util.ETime `json:"createTime" gorm:"column:created_at"`
	CreateBy   string     `json:"createBy"`
	UpdateTime util.ETime `json:"updateTime" gorm:"column:updated_at"`
	UpdateBy   string     `json:"updateBy"`
}

type ListPathResp struct {
	Paging  miso.Paging `json:"paging"`
	Payload []WPath     `json:"payload"`
}

type BindPathResReq struct {
	PathNo  string `json:"pathNo" validation:"notEmpty"`
	ResCode string `json:"resCode" validation:"notEmpty"`
}

type UnbindPathResReq struct {
	PathNo  string `json:"pathNo" validation:"notEmpty"`
	ResCode string `json:"resCode" validation:"notEmpty"`
}

type ListRoleResReq struct {
	Paging miso.Paging `json:"paging"`
	RoleNo string      `json:"roleNo" validation:"notEmpty"`
}

type RemoveRoleResReq struct {
	RoleNo  string `json:"roleNo" validation:"notEmpty"`
	ResCode string `json:"resCode" validation:"notEmpty"`
}

type AddRoleResReq struct {
	RoleNo  string `json:"roleNo" validation:"notEmpty"`
	ResCode string `json:"resCode" validation:"notEmpty"`
}

type ListRoleResResp struct {
	Paging  miso.Paging     `json:"paging"`
	Payload []ListedRoleRes `json:"payload"`
}

type ListedRoleRes struct {
	Id         int        `json:"id"`
	ResCode    string     `json:"resCode"`
	ResName    string     `json:"resName"`
	CreateTime util.ETime `json:"createTime"`
	CreateBy   string     `json:"createBy"`
}

type GenResScriptReq struct {
	ResCodes []string `json:"resCodes" validation:"notEmpty"`
}

type UpdatePathReq struct {
	Type    string `valid:"notEmpty,member:PROTECTED|PUBLIC" desc:"path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible"`
	PathNo  string `valid:"notEmpty"`
	Group   string `valid:"notEmpty,maxLen:20"`
	ResCode string
}

type CreatePathReq struct {
	Type    string `valid:"notEmpty,member:PROTECTED|PUBLIC" desc:"path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible"`
	Url     string `valid:"notEmpty,maxLen:128"`
	Group   string `valid:"notEmpty,maxLen:20"`
	Method  string `valid:"notEmpty,maxLen:10"`
	Desc    string `valid:"maxLen:255"`
	ResCode string
}

type DeletePathReq struct {
	PathNo string `json:"pathNo" validation:"notEmpty"`
}

type ListResReq struct {
	Paging miso.Paging `json:"paging"`
}

type ListResResp struct {
	Paging  miso.Paging `json:"paging"`
	Payload []WRes      `json:"payload"`
}

type CreateResReq struct {
	Name string `json:"name" validation:"notEmpty,maxLen:32"`
	Code string `json:"code" validation:"notEmpty,maxLen:32"`
}

type DeleteResourceReq struct {
	ResCode string `json:"resCode" validation:"notEmpty"`
}

func DeleteResource(rail miso.Rail, req DeleteResourceReq) error {

	_, err := lockResourceGlobal(rail, func() (any, error) {
		return nil, mysql.GetMySQL().Transaction(func(tx *gorm.DB) error {
			q := dbquery.NewQuery(rail, tx)
			if _, err := q.Exec(`delete from resource where code = ?`, req.ResCode); err != nil {
				return err
			}
			if _, err := q.Exec(`delete from role_resource where res_code = ?`, req.ResCode); err != nil {
				return err
			}
			_, err := q.Exec(`delete from path_resource where res_code = ?`, req.ResCode)
			return err
		})
	})
	return err
}

func ListResourceCandidatesForRole(rail miso.Rail, roleNo string) ([]ResBrief, error) {
	if roleNo == "" {
		return []ResBrief{}, nil
	}

	var res []ResBrief
	_, err := dbquery.NewQuery(rail, mysql.GetMySQL()).
		Select("r.name, r.code").
		Table("resource r").
		Where("NOT EXISTS (SELECT * FROM role_resource WHERE role_no = ? and res_code = r.code)", roleNo).
		Scan(&res)
	if err != nil {
		return nil, err
	}
	if res == nil {
		res = []ResBrief{}
	}
	return res, nil
}

func ListAllResBriefsOfRole(rail miso.Rail, roleNo string) ([]ResBrief, error) {
	var res []ResBrief

	if isDefAdmin(roleNo) {
		return ListAllResBriefs(rail)
	}

	_, err := dbquery.NewQuery(rail, mysql.GetMySQL()).
		Select(`r.name, r.code`).
		Table(`role_resource rr`).
		Joins(`LEFT JOIN resource r ON r.code = rr.res_code`).
		Where(`rr.role_no = ?`, roleNo).
		Scan(&res)
	if err != nil {
		return nil, err
	}
	if res == nil {
		res = []ResBrief{}
	}
	return res, nil
}

func ListAllResBriefs(rail miso.Rail) ([]ResBrief, error) {
	var res []ResBrief
	_, err := dbquery.NewQuery(rail).
		Table("resource").
		SelectCols(ResBrief{}).
		Scan(&res)
	if err != nil {
		return nil, err
	}
	if res == nil {
		res = []ResBrief{}
	}
	return res, nil
}

func ListResources(rail miso.Rail, req ListResReq) (ListResResp, error) {
	var resources []WRes
	_, err := dbquery.NewQuery(rail, mysql.GetMySQL()).
		Table("resource").
		SelectCols(WRes{}).
		OrderDesc("id").
		Limit(req.Paging.GetLimit()).
		Offset(req.Paging.GetOffset()).
		Scan(&resources)
	if err != nil {
		return ListResResp{}, err
	}
	if resources == nil {
		resources = []WRes{}
	}

	count, err := dbquery.NewQuery(rail, mysql.GetMySQL()).Count()
	if err != nil {
		return ListResResp{}, err
	}

	return ListResResp{Paging: miso.RespPage(req.Paging, int(count)), Payload: resources}, nil
}

func UpdatePath(rail miso.Rail, req UpdatePathReq) error {
	_, e := lockPath(rail, req.PathNo, func() (any, error) {
		return nil, mysql.GetMySQL().Transaction(func(tx *gorm.DB) error {
			err := dbquery.NewQuery(rail, tx).
				Table("path").
				Set("pgroup", req.Group).
				Set("ptype", req.Type).
				Eq("path_no", req.PathNo).
				UpdateAny()
			if err != nil {
				return err
			}

			ok, err := dbquery.NewQuery(rail, tx).
				Table("path_resource").
				Eq("path_no", req.PathNo).
				Eq("res_code", req.ResCode).
				HasAny()
			if err != nil {
				return err
			}
			if !ok {
				err = dbquery.NewQuery(rail, tx).
					Table("path_resource").
					CreateAny(struct {
						PathNo  string
						ResCode string
					}{req.PathNo, req.ResCode})
				return err
			}
			return nil
		})
	})
	return e
}

func GetRoleInfo(rail miso.Rail, req api.RoleInfoReq) (api.RoleInfoResp, error) {
	resp, err := roleInfoCache.GetValElse(rail, req.RoleNo, func() (api.RoleInfoResp, error) {
		var resp api.RoleInfoResp
		n, err := dbquery.NewQuery(rail, mysql.GetMySQL()).
			Raw("select role_no, name from role where role_no = ?", req.RoleNo).
			Scan(&resp)
		if err != nil {
			return resp, err
		}
		if n < 1 {
			return resp, errs.NewErrf("Role not found").WithCode(ErrCodeRoleNotFound)
		}
		return resp, nil
	})
	return resp, err
}

func CreateResourceIfNotExist(rail miso.Rail, req CreateResReq, user common.User) error {
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(req.Code)

	_, e := lockResourceGlobal(rail, func() (any, error) {
		ok, err := dbquery.NewQuery(rail).
			Table("resource").
			Eq("code", req.Code).
			HasAny()
		if err != nil {
			return nil, err
		}

		if ok {
			rail.Debugf("Resource '%s' (%s) already exist", req.Code, req.Name)
			return nil, nil
		}

		res := struct {
			Code string
			Name string
		}{
			Name: req.Name,
			Code: req.Code,
		}
		_, err = dbquery.NewQuery(rail, mysql.GetMySQL()).
			Table("resource").
			Create(&res)
		return nil, err
	})
	return e
}

func genPathNo(group string, url string, method string) string {
	cksum := md5.Sum([]byte(group + method + url))
	return "path_" + base64.StdEncoding.EncodeToString(cksum[:])
}

func CreatePath(rail miso.Rail, req CreatePathReq, user common.User) error {
	req.Url = preprocessUrl(req.Url)
	req.Group = strings.TrimSpace(req.Group)
	req.Method = strings.ToUpper(strings.TrimSpace(req.Method))
	pathNo := genPathNo(req.Group, req.Url, req.Method)

	type path struct {
		Id     int
		Pgroup string
		PathNo string
		Desc   string
		Url    string
		Method string
		Ptype  string
	}

	_, err := lockPath(rail, pathNo, func() (bool, error) {
		db := mysql.GetMySQL()
		var prev path
		ok, err := dbquery.NewQuery(rail, db).
			Table("path").
			Eq("path_no", pathNo).
			ScanAny(&prev)
		if err != nil {
			return false, err
		}
		if ok { // exists already
			rail.Debugf("Path '%s %s' (%s) already exists", req.Method, req.Url, pathNo)
			if prev.Ptype != req.Type {
				err := dbquery.NewQuery(rail, db).
					Table("path").
					Set("ptype", req.Type).
					Eq("path_no", pathNo).
					UpdateAny()
				if err != nil {
					rail.Errorf("Failed to update path.ptype, pathNo: %v, %v", pathNo, err)
					return false, err
				}
			}
			return false, nil
		}

		ep := path{
			Url:    req.Url,
			Desc:   req.Desc,
			Ptype:  req.Type,
			Pgroup: req.Group,
			Method: req.Method,
			PathNo: pathNo,
		}
		_, err = dbquery.NewQuery(rail, db).
			Table("path").
			Omit("Id").
			Create(&ep)
		if err != nil {
			return false, err
		}

		rail.Infof("Created path (%s) '%s {%s}'", pathNo, req.Method, req.Url)
		return true, nil
	})
	if err != nil {
		return err
	}

	if req.ResCode != "" { // rebind path and resource
		return BindPathRes(rail, BindPathResReq{PathNo: pathNo, ResCode: req.ResCode})
	}

	return nil
}

func DeletePath(rail miso.Rail, req DeletePathReq) error {
	req.PathNo = strings.TrimSpace(req.PathNo)
	_, e := lockPath(rail, req.PathNo, func() (any, error) {
		er := mysql.GetMySQL().Transaction(func(tx *gorm.DB) error {
			_, err := dbquery.NewQuery(rail, tx).Exec(`delete from path where path_no = ?`, req.PathNo)
			if err != nil {
				return err
			}

			_, err = dbquery.NewQuery(rail, tx).Exec(`delete from path_resource where path_no = ?`, req.PathNo)
			return err
		})

		return nil, er
	})
	return e
}

func UnbindPathRes(rail miso.Rail, req UnbindPathResReq) error {
	req.PathNo = strings.TrimSpace(req.PathNo)
	_, e := lockPath(rail, req.PathNo, func() (any, error) {
		_, err := dbquery.NewQuery(rail, mysql.GetMySQL()).
			Exec(`delete from path_resource where path_no = ?`, req.PathNo)
		return nil, err
	})
	return e
}

func BindPathRes(rail miso.Rail, req BindPathResReq) error {
	req.PathNo = strings.TrimSpace(req.PathNo)
	e := lockPathExec(rail, req.PathNo, func() error { // lock for path
		return lockResourceGlobalExec(rail, func() error {
			db := mysql.GetMySQL()

			// check if resource exist
			var resId int
			_, err := dbquery.NewQuery(rail, db).
				Raw(`SELECT id FROM resource WHERE code = ?`, req.ResCode).
				Scan(&resId)
			if err != nil {
				return err
			}
			if resId < 1 {
				rail.Errorf("Resource %v not found", req.ResCode)
				return errs.NewErrf("Resource not found")
			}

			// check if the path is already bound to current resource
			var prid int
			_, err = dbquery.NewQuery(rail, db).
				Raw(`SELECT id FROM path_resource WHERE path_no = ? AND res_code = ? LIMIT 1`, req.PathNo, req.ResCode).
				Scan(&prid)

			if err != nil {
				rail.Errorf("Failed to bind path %v to resource %v, %v", req.PathNo, req.ResCode, err)
				return err
			}
			if prid > 0 {
				rail.Debugf("Path %v already bound to resource %v", req.PathNo, req.ResCode)
				return err
			}

			// bind resource to path
			return dbquery.NewQuery(rail, db).
				Table("path_resource").
				CreateAny(struct {
					PathNo  string
					ResCode string
				}(req))
		})
	})

	return e
}

func ListPaths(rail miso.Rail, req ListPathReq) (ListPathResp, error) {

	applyCond := func(t *dbquery.Query) *dbquery.Query {
		if req.Pgroup != "" {
			t = t.Where("p.pgroup = ?", req.Pgroup)
		}
		if req.ResCode != "" {
			t = t.Joins("LEFT JOIN path_resource pr ON p.path_no = pr.path_no").
				Where("pr.res_code = ?", req.ResCode)
		}
		if req.Url != "" {
			t = t.Where("p.url LIKE ?", "%"+req.Url+"%")
		}
		if req.Ptype != "" {
			t = t.Where("p.ptype = ?", req.Ptype)
		}
		return t
	}

	var paths []WPath
	err := applyCond(dbquery.NewQuery(rail).
		Table("path p").
		Select("p.*").
		Order("id DESC")).
		Offset(req.Paging.GetOffset()).
		Limit(req.Paging.GetLimit()).
		ScanVal(&paths)
	if err != nil {
		return ListPathResp{}, err
	}

	var count int
	err = applyCond(dbquery.NewQuery(rail).
		Table("path p").
		Select("COUNT(*)")).
		ScanVal(&count)
	if err != nil {
		return ListPathResp{}, err
	}

	return ListPathResp{
		Payload: paths,
		Paging:  miso.Paging{Limit: req.Paging.Limit, Page: req.Paging.Page, Total: count},
	}, nil
}

func AddRole(rail miso.Rail, req AddRoleReq, user common.User) error {
	_, e := redis.RLockRun(rail, "user-vault:role:add"+req.Name, func() (any, error) {
		return nil, dbquery.NewQuery(rail, mysql.GetMySQL()).
			Table("role").
			CreateAny(struct {
				RoleNo string
				Name   string
			}{
				RoleNo: util.GenIdP("role_"),
				Name:   req.Name,
			})
	})
	return e
}

func RemoveResFromRole(rail miso.Rail, req RemoveRoleResReq) error {
	_, e := redis.RLockRun(rail, "user-vault:role:"+req.RoleNo, func() (any, error) {
		_, err := dbquery.NewQuery(rail, dbquery.GetDB()).
			Exec(`delete from role_resource where role_no = ? and res_code = ?`, req.RoleNo, req.ResCode)
		return nil, err
	})

	return e
}

func AddResToRoleIfNotExist(rail miso.Rail, req AddRoleResReq, user common.User) error {

	_, e := redis.RLockRun(rail, "user-vault:role:"+req.RoleNo, func() (any, error) { // lock for role
		return lockResourceGlobal(rail, func() (any, error) {
			// check if resource exist
			var resId int
			_, err := dbquery.NewQuery(rail, mysql.GetMySQL()).
				Raw(`select id from resource where code = ?`, req.ResCode).
				Scan(&resId)
			if err != nil {
				return false, err
			}
			if resId < 1 {
				return false, errs.NewErrf("Resource not found")
			}

			// check if role-resource relation exists
			ok, err := dbquery.NewQuery(rail, mysql.GetMySQL()).
				Table("role_resource").
				Eq("role_no", req.RoleNo).
				Eq("res_code", req.ResCode).
				HasAny()
			if err != nil {
				return false, err
			}
			if ok { // relation exists already
				return false, nil
			}

			// create role-resource relation
			rr := struct {
				RoleNo  string // role no
				ResCode string // resource code
			}{
				RoleNo:  req.RoleNo,
				ResCode: req.ResCode,
			}
			return true, mysql.GetMySQL().
				Table("role_resource").
				Create(&rr).Error
		})
	})
	return e
}

func ListRoleRes(ec miso.Rail, req ListRoleResReq) (ListRoleResResp, error) {
	var res []ListedRoleRes
	tx := mysql.GetMySQL().
		Raw(`select rr.id, rr.res_code, rr.create_time, rr.create_by, r.name 'res_name' from role_resource rr
			left join resource r on rr.res_code = r.code
			where rr.role_no = ? order by rr.id desc limit ?, ?`, req.RoleNo, req.Paging.GetOffset(), req.Paging.GetLimit()).
		Scan(&res)

	if tx.Error != nil {
		return ListRoleResResp{}, tx.Error
	}

	if res == nil {
		res = []ListedRoleRes{}
	}

	var count int
	tx = mysql.GetMySQL().
		Raw(`select count(*) from role_resource rr
			left join resource r on rr.res_code = r.code
			where rr.role_no = ?`, req.RoleNo).
		Scan(&count)

	if tx.Error != nil {
		return ListRoleResResp{}, tx.Error
	}

	return ListRoleResResp{Payload: res,
		Paging: miso.Paging{Limit: req.Paging.Limit, Page: req.Paging.Page, Total: count}}, nil
}

func ListAllRoleBriefs(rail miso.Rail) ([]RoleBrief, error) {
	var roles []RoleBrief
	tx := mysql.GetMySQL().Raw("select role_no, name from role").Scan(&roles)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if roles == nil {
		roles = []RoleBrief{}
	}
	return roles, nil
}

func ListRoles(rail miso.Rail, req ListRoleReq) (ListRoleResp, error) {
	var roles []WRole
	err := dbquery.NewQuery(rail).
		Raw("select * from role order by id desc limit ?, ?", req.Paging.GetOffset(), req.Paging.GetLimit()).
		ScanVal(&roles)
	if err != nil {
		return ListRoleResp{}, err
	}
	if roles == nil {
		roles = []WRole{}
	}

	count, err := dbquery.NewQuery(rail).
		Table("role").
		Count()
	if err != nil {
		return ListRoleResp{}, err
	}

	return ListRoleResp{
		Payload: roles,
		Paging:  miso.Paging{Limit: req.Paging.Limit, Page: req.Paging.Page, Total: int(count)},
	}, nil
}

// Test access to resource
func TestResourceAccess(rail miso.Rail, req api.CheckResAccessReq) (api.CheckResAccessResp, error) {
	url := req.Url
	roleNo := req.RoleNo

	// some sanitization & standardization for the url
	url = preprocessUrl(url)
	method := strings.ToUpper(strings.TrimSpace(req.Method))
	match := func(p PathAccessInfo) bool {
		if p.Method != "*" && p.Method != method {
			return false
		}
		ok, err := doublestar.Match(p.Url, url)
		if err != nil {
			rail.Errorf("Path Pattern is invalid, %v, %v", p.Url, err)
			return false
		}
		if ok {
			rail.Infof("Request path matched, '%v %v', resource: %v (%v), roleNo: %v", p.Method, p.Url,
				p.ResCode, p.Ptype, roleNo)
		}
		return ok
	}

	if roleNo == "" {
		public, ok, err := publicAccessCache.Get(rail, "")
		if err != nil {
			rail.Warnf("Failed to load PublicAccessCache, %v", err)
			return forbidden, nil
		}
		if !ok {
			return forbidden, nil
		}
		for _, p := range public {
			if !match(p) {
				continue
			}
			return permitted, nil
		}
		rail.Infof("Rejected '%v %s', roleNo: '%s', role doesn't have access to required resource", method, url, roleNo)
		return forbidden, nil
	}

	rr, ok, err := roleAccessCache.Get(rail, roleNo)
	if err != nil {
		rail.Warnf("Failed to find RoleAccess for %v from cache, %v", roleNo, err)
		return forbidden, nil
	}
	if !ok {
		return forbidden, nil
	}

	for _, p := range rr.Paths {
		if !match(p) {
			continue
		}
		return permitted, nil
	}

	// doesn't even have role
	roleNo = strings.TrimSpace(roleNo)
	if roleNo == "" {
		rail.Infof("Rejected '%s', user doesn't have roleNo", url)
		return forbidden, nil
	}

	// the role doesn't have access to the required resource
	rail.Infof("Rejected '%v %s', roleNo: '%s', role doesn't have access to required resource", method, url, roleNo)
	return forbidden, nil
}

func listRoleNos(rail miso.Rail) ([]string, error) {
	var ern []string
	err := dbquery.NewQuery(rail).Raw("select role_no from role").ScanVal(&ern)
	if err != nil {
		return nil, err
	}

	if ern == nil {
		ern = []string{}
	}
	return ern, nil
}

// preprocess url, the processed url will always starts with '/' and never ends with '/'
func preprocessUrl(url string) string {
	ru := []rune(strings.TrimSpace(url))
	l := len(ru)
	if l < 1 {
		return "/"
	}

	j := strings.LastIndex(url, "?")
	if j > -1 {
		ru = ru[0:j]
		l = len(ru)
	}

	// never ends with '/'
	if ru[l-1] == '/' && l > 1 {
		lj := l - 1
		for lj > 1 && ru[lj-1] == '/' {
			lj -= 1
		}

		ru = ru[0:lj]
	}

	// always start with '/'
	if ru[0] != '/' {
		return "/" + string(ru)
	}
	return string(ru)
}

// global lock for resources
func lockResourceGlobal(ec miso.Rail, runnable redis.LRunnable[any]) (any, error) {
	return redis.RLockRun(ec, "user-vault:resource:global", runnable)
}

// global lock for resources
func lockResourceGlobalExec(ec miso.Rail, runnable redis.Runnable) error {
	return redis.RLockExec(ec, "user-vault:resource:global", runnable)
}

// lock for path
func lockPath[T any](ec miso.Rail, pathNo string, runnable redis.LRunnable[T]) (T, error) {
	return redis.RLockRun(ec, "user-vault:path:"+pathNo, runnable)
}

// lock for path
func lockPathExec(ec miso.Rail, pathNo string, runnable redis.Runnable) error {
	return redis.RLockExec(ec, "user-vault:path:"+pathNo, runnable)
}

func isDefAdmin(roleNo string) bool {
	return roleNo == DefaultAdminRoleNo || roleNo == DefaultAdminRoleNo2
}

type RoleAccess struct {
	Paths []PathAccessInfo
}

type PathAccessInfo struct {
	ResCode string // resource code
	Url     string // url
	Method  string // http method
	Ptype   string // path type: PROTECTED, PUBLIC
}

func BatchLoadRoleAccessCache(rail miso.Rail) error {

	_, e := lockRoleAccessCache(rail, func() (any, error) {

		lr, e := listRoleNos(rail)
		if e != nil {
			return nil, e
		}
		lr = append(lr, DefaultAdminRoleNo, DefaultAdminRoleNo2)

		for _, roleNo := range lr {
			e = LoadOneRoleAccessCache(rail, roleNo)
			if e != nil {
				return nil, e
			}
		}
		return nil, nil
	})
	return e
}

func LoadOneRoleAccessCache(rail miso.Rail, roleNo string) error {
	var paths []ExtendedPathRes

	if isDefAdmin(roleNo) {
		err := dbquery.NewQuery(rail).
			Raw(`SELECT p.*, pr.res_code
				FROM path_resource pr
				LEFT JOIN path p ON p.path_no = pr.path_no`).
			ScanVal(&paths)
		if err != nil {
			return err
		}
	} else {
		err := dbquery.NewQuery(rail).
			Raw(`SELECT p.*, pr.res_code
		FROM role_resource rr
		LEFT JOIN path_resource pr ON rr.res_code = pr.res_code
		LEFT JOIN path p ON p.path_no = pr.path_no
		WHERE rr.role_no = ?
		`, roleNo).
			ScanVal(&paths)
		if err != nil {
			return err
		}
	}
	if paths == nil {
		return nil
	}

	var public []ExtendedPathRes
	err := dbquery.NewQuery(rail).
		Raw(`SELECT p.* FROM path p WHERE p.ptype = ?`, PathTypePublic).
		ScanVal(&public)
	if err != nil {
		return err
	}
	if len(public) > 0 {
		paths = append(paths, public...)
	}

	var pai []PathAccessInfo = slutil.MapTo(paths,
		func(t ExtendedPathRes) PathAccessInfo {
			return PathAccessInfo{
				ResCode: t.ResCode,
				Url:     preprocessUrl(t.Url),
				Method:  t.Method,
				Ptype:   t.Ptype,
			}
		})
	cached := RoleAccess{
		Paths: pai,
	}
	err = roleAccessCache.Put(rail, roleNo, cached)
	if err == nil {
		rail.Infof("Updated RoleAccessCache for %v, path counts: %v, public paths: %v", roleNo, len(pai), len(public))
		return nil
	}
	return err
}

func LoadPublicAccessCache(rail miso.Rail) error {
	var public []ExtendedPathRes
	err := dbquery.NewQuery(rail).
		Raw(`SELECT p.* FROM path p WHERE p.ptype = ?`, PathTypePublic).
		ScanVal(&public)
	if err != nil {
		return err
	}
	if public == nil {
		return nil
	}

	var pai []PathAccessInfo = slutil.MapTo(public,
		func(t ExtendedPathRes) PathAccessInfo {
			return PathAccessInfo{
				ResCode: t.ResCode,
				Url:     preprocessUrl(t.Url),
				Method:  t.Method,
				Ptype:   t.Ptype,
			}
		})
	err = publicAccessCache.Put(rail, "", pai)
	if err == nil {
		rail.Infof("Updated PublicAccessCache, public paths: %v", len(public))
		return nil
	}
	return err
}

// lock for role-access cache
func lockRoleAccessCache(ec miso.Rail, runnable redis.LRunnable[any]) (any, error) {
	return redis.RLockRun(ec, "user-vault:role:access:cache", runnable)
}

func RegisterInternalPathResourcesOnBootstrapped(res []auth.Resource) {

	miso.PostServerBootstrap(func(rail miso.Rail) error {

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
