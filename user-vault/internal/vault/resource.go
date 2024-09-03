package vault

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"github.com/curtisnewbie/user-vault/api"
	"gorm.io/gorm"
)

var (
	permitted = TestResAccessResp{Valid: true}
	forbidden = TestResAccessResp{Valid: false}

	roleInfoCache = miso.NewRCache[api.RoleInfoResp]("user-vault:role:info", miso.RCacheConfig{Exp: 10 * time.Minute, NoSync: true})

	// cache for url's resource, url -> CachedUrlRes
	urlResCache = miso.NewRCache[CachedUrlRes]("user-vault:url:res:v2", miso.RCacheConfig{Exp: 30 * time.Minute})

	// cache for role's resource, role + res -> flag ("1")
	roleResCache = miso.NewRCache[string]("user-vault:role:res", miso.RCacheConfig{Exp: 1 * time.Hour, NoSync: true})
)

const (
	// default roleno for admin
	DefaultAdminRoleNo = "role_554107924873216177918"

	PathTypeProtected string = "PROTECTED"
	PathTypePublic    string = "PUBLIC"
)

type PathRes struct {
	Id         int    // id
	PathNo     string // path no
	ResCode    string // resource code
	CreateTime util.ETime
	CreateBy   string
	UpdateTime util.ETime
	UpdateBy   string
}

type ExtendedPathRes struct {
	Id         int    // id
	Pgroup     string // path group
	PathNo     string // path no
	ResCode    string // resource code
	Desc       string // description
	Url        string // url
	Method     string // http method
	Ptype      string // path type: PROTECTED, PUBLIC
	CreateTime util.ETime
	CreateBy   string
	UpdateTime util.ETime
	UpdateBy   string
}

type EPath struct {
	Id         int    // id
	Pgroup     string // path group
	PathNo     string // path no
	Desc       string // description
	Url        string // url
	Method     string // method
	Ptype      string // path type: PROTECTED, PUBLIC
	CreateTime util.ETime
	CreateBy   string
	UpdateTime util.ETime
	UpdateBy   string
}

type ERes struct {
	Id         int    // id
	Code       string // resource code
	Name       string // resource name
	CreateTime util.ETime
	CreateBy   string
	UpdateTime util.ETime
	UpdateBy   string
}

type ERoleRes struct {
	Id         int    // id
	RoleNo     string // role no
	ResCode    string // resource code
	CreateTime util.ETime
	CreateBy   string
	UpdateTime util.ETime
	UpdateBy   string
}

type ERole struct {
	Id         int
	RoleNo     string
	Name       string
	CreateTime util.ETime
	CreateBy   string
	UpdateTime util.ETime
	UpdateBy   string
}

type WRole struct {
	Id         int        `json:"id"`
	RoleNo     string     `json:"roleNo"`
	Name       string     `json:"name"`
	CreateTime util.ETime `json:"createTime"`
	CreateBy   string     `json:"createBy"`
	UpdateTime util.ETime `json:"updateTime"`
	UpdateBy   string     `json:"updateBy"`
}

type CachedUrlRes struct {
	Id      int    // id
	Pgroup  string // path group
	PathNo  string // path no
	ResCode string // resource code
	Url     string // url
	Method  string // http method
	Ptype   string // path type: PROTECTED, PUBLIC
}

type ResBrief struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type AddRoleReq struct {
	Name string `json:"name" validation:"notEmpty,maxLen:32"` // role name
}

type TestResAccessReq struct {
	RoleNo string `json:"roleNo"`
	Url    string `json:"url"`
	Method string `json:"method"`
}

type TestResAccessResp struct {
	Valid bool `json:"valid"`
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
	CreateTime util.ETime `json:"createTime"`
	CreateBy   string     `json:"createBy"`
	UpdateTime util.ETime `json:"updateTime"`
	UpdateBy   string     `json:"updateBy"`
}

type WRes struct {
	Id         int        `json:"id"`
	Code       string     `json:"code"`
	Name       string     `json:"name"`
	CreateTime util.ETime `json:"createTime"`
	CreateBy   string     `json:"createBy"`
	UpdateTime util.ETime `json:"updateTime"`
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
	Type   string `json:"type" validation:"notEmpty" desc:"path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible"`
	PathNo string `json:"pathNo" validation:"notEmpty"`
	Group  string `json:"group" validation:"notEmpty,maxLen:20"`
}

type CreatePathReq struct {
	Type    string `json:"type" validation:"notEmpty" desc:"path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible"`
	Url     string `json:"url" validation:"notEmpty,maxLen:128"`
	Group   string `json:"group" validation:"notEmpty,maxLen:20"`
	Method  string `json:"method" validation:"notEmpty,maxLen:10"`
	Desc    string `json:"desc" validation:"maxLen:255"`
	ResCode string `json:"resCode"`
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
		return nil, miso.GetMySQL().Transaction(func(tx *gorm.DB) error {
			if t := tx.Exec(`delete from resource where code = ?`, req.ResCode); t != nil {
				return t.Error
			}
			if t := tx.Exec(`delete from role_resource where res_code = ?`, req.ResCode); t != nil {
				return t.Error
			}
			return tx.Exec(`delete from path_resource where res_code = ?`, req.ResCode).Error
		})
	})

	if err == nil {
		// asynchronously reload the cache of paths and resources
		commonPool.Go(func() {
			rail := rail.NextSpan()
			if err := LoadPathResCache(rail); err != nil {
				rail.Errorf("Failed to load path resource cache, %v", err)
			}
		})
		// asynchronously reload the cache of role and resources
		commonPool.Go(func() {
			rail := rail.NextSpan()
			if err := LoadRoleResCache(rail); err != nil {
				rail.Errorf("Failed to load role resource cache, %v", err)
			}
		})
	}

	return err
}

func ListResourceCandidatesForRole(ec miso.Rail, roleNo string) ([]ResBrief, error) {
	if roleNo == "" {
		return []ResBrief{}, nil
	}

	var res []ResBrief
	tx := miso.GetMySQL().
		Select("r.name, r.code").
		Table("resource r").
		Where("NOT EXISTS (SELECT * FROM role_resource WHERE role_no = ? and res_code = r.code)", roleNo).
		Scan(&res)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if res == nil {
		res = []ResBrief{}
	}
	return res, nil
}

func ListAllResBriefsOfRole(ec miso.Rail, roleNo string) ([]ResBrief, error) {
	var res []ResBrief

	if roleNo == DefaultAdminRoleNo {
		return ListAllResBriefs(ec)
	}

	tx := miso.GetMySQL().
		Select(`r.name, r.code`).
		Table(`role_resource rr`).
		Joins(`LEFT JOIN resource r ON r.code = rr.res_code`).
		Where(`rr.role_no = ?`, roleNo).
		Scan(&res)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if res == nil {
		res = []ResBrief{}
	}
	return res, nil
}

func ListAllResBriefs(rail miso.Rail) ([]ResBrief, error) {
	var res []ResBrief
	tx := miso.GetMySQL().Raw("select name, code from resource").Scan(&res)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if res == nil {
		res = []ResBrief{}
	}
	return res, nil
}

func ListResources(ec miso.Rail, req ListResReq) (ListResResp, error) {
	var resources []WRes
	tx := miso.GetMySQL().
		Raw("select * from resource order by id desc limit ?, ?", req.Paging.GetOffset(), req.Paging.GetLimit()).
		Scan(&resources)
	if tx.Error != nil {
		return ListResResp{}, tx.Error
	}
	if resources == nil {
		resources = []WRes{}
	}

	var count int
	tx = miso.GetMySQL().Raw("select count(*) from resource").Scan(&count)
	if tx.Error != nil {
		return ListResResp{}, tx.Error
	}

	return ListResResp{Paging: miso.RespPage(req.Paging, count), Payload: resources}, nil
}

func UpdatePath(ec miso.Rail, req UpdatePathReq) error {
	// TODO: validate the ptype value
	_, e := lockPath(ec, req.PathNo, func() (any, error) {
		tx := miso.GetMySQL().Exec(`update path set pgroup = ?, ptype = ? where path_no = ?`,
			req.Group, req.Type, req.PathNo)
		return nil, tx.Error
	})

	if e == nil {
		loadOnePathResCacheAsync(ec, req.PathNo)
	}
	return e
}

func loadOnePathResCacheAsync(rail miso.Rail, pathNo string) {
	commonPool.Go(func() {
		rail := rail.NextSpan()
		// ec.Infof("Refreshing path cache, pathNo: %s", pathNo)
		ep, e := findPathRes(pathNo)
		if e != nil {
			rail.Errorf("Failed to reload path cache, pathNo: %s, %v", pathNo, e)
			return
		}

		ep.Url = preprocessUrl(ep.Url)
		if e := urlResCache.Put(rail, ep.Method+":"+ep.Url, toCachedUrlRes(ep)); e != nil {
			rail.Errorf("Failed to save cached url resource, pathNo: %s, %v", pathNo, e)
			return
		}
	})
}

func GetRoleInfo(ec miso.Rail, req api.RoleInfoReq) (api.RoleInfoResp, error) {
	resp, err := roleInfoCache.Get(ec, req.RoleNo, func() (api.RoleInfoResp, error) {
		var resp api.RoleInfoResp
		tx := miso.GetMySQL().Raw("select role_no, name from role where role_no = ?", req.RoleNo).Scan(&resp)
		if tx.Error != nil {
			return resp, tx.Error
		}

		if tx.RowsAffected < 1 {
			return resp, miso.NewErrf("Role not found").WithCode(ErrCodeRoleNotFound)
		}
		return resp, nil
	})
	return resp, err
}

func CreateResourceIfNotExist(rail miso.Rail, req CreateResReq, user common.User) error {
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(req.Code)

	_, e := lockResourceGlobal(rail, func() (any, error) {
		var id int
		tx := miso.GetMySQL().Raw(`select id from resource where code = ? limit 1`, req.Code).Scan(&id)
		if tx.Error != nil {
			return nil, tx.Error
		}

		if id > 0 {
			rail.Debugf("Resource '%s' (%s) already exist", req.Code, req.Name)
			return nil, nil
		}

		res := ERes{
			Name:     req.Name,
			Code:     req.Code,
			CreateBy: user.Username,
			UpdateBy: user.Username,
		}

		tx = miso.GetMySQL().
			Table("resource").
			Omit("Id", "CreateTime", "UpdateTime").
			Create(&res)
		return nil, tx.Error
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

	changed, err := lockPath(rail, pathNo, func() (bool, error) {
		var prev EPath
		tx := miso.GetMySQL().Raw(`select * from path where path_no = ? limit 1`, pathNo).Scan(&prev)
		if tx.Error != nil {
			return false, tx.Error
		}
		if prev.Id > 0 { // exists already
			rail.Debugf("Path '%s %s' (%s) already exists", req.Method, req.Url, pathNo)
			if prev.Ptype != req.Type {
				err := miso.GetMySQL().Exec(`UPDATE path SET ptype = ? WHERE path_no = ?`, req.Type, pathNo).Error
				if err != nil {
					rail.Errorf("failed to update path.ptype, pathNo: %v, %v", pathNo, err)
					return false, err
				}
			}
			return false, nil
		}

		ep := EPath{
			Url:      req.Url,
			Desc:     req.Desc,
			Ptype:    req.Type,
			Pgroup:   req.Group,
			Method:   req.Method,
			PathNo:   pathNo,
			CreateBy: user.Username,
			UpdateBy: user.Username,
		}
		tx = miso.GetMySQL().
			Table("path").
			Omit("Id", "CreateTime", "UpdateTime").
			Create(&ep)
		if tx.Error != nil {
			return false, tx.Error
		}

		rail.Infof("Created path (%s) '%s {%s}'", pathNo, req.Method, req.Url)
		return true, nil
	})
	if err != nil {
		return err
	}

	if changed { // reload cache for the path
		loadOnePathResCacheAsync(rail, pathNo)
	}

	if req.ResCode != "" { // rebind path and resource
		return BindPathRes(rail, BindPathResReq{PathNo: pathNo, ResCode: req.ResCode})
	}

	return nil
}

func DeletePath(ec miso.Rail, req DeletePathReq) error {
	req.PathNo = strings.TrimSpace(req.PathNo)
	_, e := lockPath(ec, req.PathNo, func() (any, error) {
		er := miso.GetMySQL().Transaction(func(tx *gorm.DB) error {
			tx = tx.Exec(`delete from path where path_no = ?`, req.PathNo)
			if tx.Error != nil {
				return tx.Error
			}

			return tx.Exec(`delete from path_resource where path_no = ?`, req.PathNo).Error
		})

		return nil, er
	})
	return e
}

func UnbindPathRes(ec miso.Rail, req UnbindPathResReq) error {
	req.PathNo = strings.TrimSpace(req.PathNo)
	_, e := lockPath(ec, req.PathNo, func() (any, error) {
		tx := miso.GetMySQL().Exec(`delete from path_resource where path_no = ?`, req.PathNo)
		return nil, tx.Error
	})

	if e != nil {
		// asynchronously reload the cache of paths and resources
		commonPool.Go(func() {
			if e := LoadPathResCache(ec); e != nil {
				ec.Errorf("Failed to load path resource cache, %v", e)
			}
		})
	}
	return e
}

func BindPathRes(rail miso.Rail, req BindPathResReq) error {
	req.PathNo = strings.TrimSpace(req.PathNo)
	e := lockPathExec(rail, req.PathNo, func() error { // lock for path
		return lockResourceGlobalExec(rail, func() error {

			// check if resource exist
			var resId int
			tx := miso.GetMySQL().
				Raw(`SELECT id FROM resource WHERE code = ?`, req.ResCode).
				Scan(&resId)
			if tx.Error != nil {
				return tx.Error
			}
			if resId < 1 {
				rail.Errorf("Resource %v not found", req.ResCode)
				return miso.NewErrf("Resource not found")
			}

			// check if the path is already bound to current resource
			var prid int
			tx = miso.GetMySQL().
				Raw(`SELECT id FROM path_resource WHERE path_no = ? AND res_code = ? LIMIT 1`, req.PathNo, req.ResCode).
				Scan(&prid)

			if tx.Error != nil {
				rail.Errorf("Failed to bind path %v to resource %v, %v", req.PathNo, req.ResCode, tx.Error)
				return tx.Error
			}
			if prid > 0 {
				rail.Debugf("Path %v already bound to resource %v", req.PathNo, req.ResCode)
				return tx.Error
			}

			// bind resource to path
			return miso.GetMySQL().
				Exec(`INSERT INTO path_resource (path_no, res_code) VALUES (?, ?)`, req.PathNo, req.ResCode).
				Error
		})
	})

	if e == nil {
		// asynchronously reload the cache of paths and resources
		loadOnePathResCacheAsync(rail, req.PathNo)
	}
	return e
}

func ListPaths(ec miso.Rail, req ListPathReq) (ListPathResp, error) {

	applyCond := func(t *gorm.DB) *gorm.DB {
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
	tx := miso.GetMySQL().
		Table("path p").
		Select("p.*").
		Order("id DESC")

	tx = applyCond(tx).
		Offset(req.Paging.GetOffset()).
		Limit(req.Paging.GetLimit()).
		Scan(&paths)
	if tx.Error != nil {
		return ListPathResp{}, tx.Error
	}

	var count int
	tx = miso.GetMySQL().
		Table("path p").
		Select("COUNT(*)")

	tx = applyCond(tx).
		Scan(&count)

	if tx.Error != nil {
		return ListPathResp{}, tx.Error
	}

	return ListPathResp{Payload: paths, Paging: miso.Paging{Limit: req.Paging.Limit, Page: req.Paging.Page, Total: count}}, nil
}

func AddRole(ec miso.Rail, req AddRoleReq, user common.User) error {
	_, e := miso.RLockRun(ec, "user-vault:role:add"+req.Name, func() (any, error) {
		r := ERole{
			RoleNo:   util.GenIdP("role_"),
			Name:     req.Name,
			CreateBy: user.Username,
			UpdateBy: user.Username,
		}
		return nil, miso.GetMySQL().
			Table("role").
			Omit("Id", "CreateTime", "UpdateTime").
			Create(&r).Error
	})
	return e
}

func RemoveResFromRole(ec miso.Rail, req RemoveRoleResReq) error {
	_, e := miso.RLockRun(ec, "user-vault:role:"+req.RoleNo, func() (any, error) {
		tx := miso.GetMySQL().Exec(`delete from role_resource where role_no = ? and res_code = ?`, req.RoleNo, req.ResCode)
		return nil, tx.Error
	})

	if e != nil {
		e = roleResCache.Put(ec, fmt.Sprintf("role:%s:res:%s", req.RoleNo, req.ResCode), "")
	}

	return e
}

func AddResToRoleIfNotExist(rail miso.Rail, req AddRoleResReq, user common.User) error {

	res, e := miso.RLockRun(rail, "user-vault:role:"+req.RoleNo, func() (any, error) { // lock for role
		return lockResourceGlobal(rail, func() (any, error) {
			// check if resource exist
			var resId int
			tx := miso.GetMySQL().Raw(`select id from resource where code = ?`, req.ResCode).Scan(&resId)
			if tx.Error != nil {
				return false, tx.Error
			}
			if resId < 1 {
				return false, miso.NewErrf("Resource not found")
			}

			// check if role-resource relation exists
			var id int
			tx = miso.GetMySQL().Raw(`select id from role_resource where role_no = ? and res_code = ?`, req.RoleNo, req.ResCode).Scan(&id)
			if tx.Error != nil {
				return false, tx.Error
			}
			if id > 0 { // relation exists already
				return false, nil
			}

			// create role-resource relation
			rr := ERoleRes{
				RoleNo:   req.RoleNo,
				ResCode:  req.ResCode,
				CreateBy: user.Username,
				UpdateBy: user.Username,
			}

			return true, miso.GetMySQL().
				Table("role_resource").
				Omit("Id", "CreateTime", "UpdateTime").
				Create(&rr).Error
		})
	})

	if e != nil {
		return e
	}

	if isAdded := res.(bool); isAdded {
		e = _loadResOfRole(rail, req.RoleNo)
	}

	return e
}

func ListRoleRes(ec miso.Rail, req ListRoleResReq) (ListRoleResResp, error) {
	var res []ListedRoleRes
	tx := miso.GetMySQL().
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
	tx = miso.GetMySQL().
		Raw(`select count(*) from role_resource rr
			left join resource r on rr.res_code = r.code
			where rr.role_no = ?`, req.RoleNo).
		Scan(&count)

	if tx.Error != nil {
		return ListRoleResResp{}, tx.Error
	}

	return ListRoleResResp{Payload: res, Paging: miso.Paging{Limit: req.Paging.Limit, Page: req.Paging.Page, Total: count}}, nil
}

func ListAllRoleBriefs(ec miso.Rail) ([]RoleBrief, error) {
	var roles []RoleBrief
	tx := miso.GetMySQL().Raw("select role_no, name from role").Scan(&roles)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if roles == nil {
		roles = []RoleBrief{}
	}
	return roles, nil
}

func ListRoles(ec miso.Rail, req ListRoleReq) (ListRoleResp, error) {
	var roles []WRole
	tx := miso.GetMySQL().
		Raw("select * from role order by id desc limit ?, ?", req.Paging.GetOffset(), req.Paging.GetLimit()).
		Scan(&roles)
	if tx.Error != nil {
		return ListRoleResp{}, tx.Error
	}
	if roles == nil {
		roles = []WRole{}
	}

	var count int
	tx = miso.GetMySQL().Raw("select count(*) from role").Scan(&count)
	if tx.Error != nil {
		return ListRoleResp{}, tx.Error
	}

	return ListRoleResp{Payload: roles, Paging: miso.Paging{Limit: req.Paging.Limit, Page: req.Paging.Page, Total: count}}, nil
}

// Test access to resource
func TestResourceAccess(ec miso.Rail, req TestResAccessReq) (TestResAccessResp, error) {
	url := req.Url
	roleNo := req.RoleNo

	// some sanitization & standardization for the url
	url = preprocessUrl(url)
	method := strings.ToUpper(strings.TrimSpace(req.Method))

	// find resource required for the url
	cur, e := lookupUrlRes(ec, url, method)
	if e != nil {
		ec.Infof("Rejected '%s' (%s), path not found", url, method)
		return forbidden, nil
	}

	// public path type, doesn't require access to resource
	if cur.Ptype == PathTypePublic {
		return permitted, nil
	}

	// doesn't even have role
	roleNo = strings.TrimSpace(roleNo)
	if roleNo == "" {
		ec.Infof("Rejected '%s', user doesn't have roleNo", url)
		return forbidden, nil
	}

	// the requiredRes resources no
	requiredRes := cur.ResCode
	if requiredRes == "" {
		ec.Infof("Rejected '%s', path doesn't have any resource bound yet", url)
		return forbidden, nil
	}

	ok, e := checkRoleRes(ec, roleNo, requiredRes)
	if e != nil {
		return forbidden, e
	}

	// the role doesn't have access to the required resource
	if !ok {
		ec.Infof("Rejected '%s', roleNo: '%s', role doesn't have access to required resource '%s'", url, roleNo, requiredRes)
		return forbidden, nil
	}

	return permitted, nil
}

func checkRoleRes(rail miso.Rail, roleNo string, resCode string) (bool, error) {
	if roleNo == DefaultAdminRoleNo {
		return true, nil
	}

	ok, e := roleResCache.Exists(rail, fmt.Sprintf("role:%s:res:%s", roleNo, resCode))
	if e != nil {
		return false, e
	}
	return ok, nil
}

// Load cache for role -> resources
func LoadRoleResCache(ec miso.Rail) error {

	_, e := lockRoleResCache(ec, func() (any, error) {

		lr, e := listRoleNos(ec)
		if e != nil {
			return nil, e
		}

		for _, roleNo := range lr {
			e = _loadResOfRole(ec, roleNo)
			if e != nil {
				return nil, e
			}
		}
		return nil, nil
	})
	return e
}

func _loadResOfRole(ec miso.Rail, roleNo string) error {
	roleResList, e := listRoleRes(ec, roleNo)
	if e != nil {
		return e
	}

	for _, rr := range roleResList {
		roleResCache.Put(ec, fmt.Sprintf("role:%s:res:%s", rr.RoleNo, rr.ResCode), "1")
	}
	return nil
}

func listRoleNos(ec miso.Rail) ([]string, error) {
	var ern []string
	t := miso.GetMySQL().Raw("select role_no from role").Scan(&ern)
	if t.Error != nil {
		return nil, t.Error
	}

	if ern == nil {
		ern = []string{}
	}
	return ern, nil
}

func listRoleRes(ec miso.Rail, roleNo string) ([]ERoleRes, error) {
	var rr []ERoleRes
	t := miso.GetMySQL().Raw("select * from role_resource where role_no = ?", roleNo).Scan(&rr)
	if t.Error != nil {
		if errors.Is(t.Error, gorm.ErrRecordNotFound) {
			return []ERoleRes{}, nil
		}
		return nil, t.Error
	}

	return rr, nil
}

func lookupUrlRes(ec miso.Rail, url string, method string) (CachedUrlRes, error) {
	cur, e := urlResCache.Get(ec, method+":"+url, nil)
	if e != nil {
		return CachedUrlRes{}, e
	}
	return cur, nil
}

// Load cache for path -> resource
func LoadPathResCache(rail miso.Rail) error {

	_, e := miso.RLockRun(rail, "user-vault:path:res:cache", func() (any, error) {
		var paths []ExtendedPathRes
		tx := miso.GetMySQL().
			Raw("select p.*, pr.res_code from path p left join path_resource pr on p.path_no = pr.path_no").
			Scan(&paths)
		if tx.Error != nil {
			return nil, tx.Error
		}
		if paths == nil {
			return nil, nil
		}

		for _, ep := range paths {
			ep.Url = preprocessUrl(ep.Url)
			if e := urlResCache.Put(rail, ep.Method+":"+ep.Url, toCachedUrlRes(ep)); e != nil {
				return nil, fmt.Errorf("failed to store urlResCache, %w", e)
			}
		}
		return nil, nil
	})

	return e
}

func toCachedUrlRes(epath ExtendedPathRes) CachedUrlRes {
	cur := CachedUrlRes{
		Id:      epath.Id,
		Pgroup:  epath.Pgroup,
		PathNo:  epath.PathNo,
		ResCode: epath.ResCode,
		Url:     epath.Url,
		Method:  epath.Method,
		Ptype:   epath.Ptype,
	}
	return cur
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

func findPathRes(pathNo string) (ExtendedPathRes, error) {
	var ep ExtendedPathRes
	tx := miso.GetMySQL().
		Raw("select p.*, pr.res_code from path p left join path_resource pr on p.path_no = pr.path_no where p.path_no = ? limit 1", pathNo).
		Scan(&ep)
	if tx.Error != nil {
		return ep, tx.Error
	}

	if tx.RowsAffected < 1 {
		return ep, miso.NewErrf("Path not found")
	}

	return ep, nil
}

// global lock for resources
func lockResourceGlobal(ec miso.Rail, runnable miso.LRunnable[any]) (any, error) {
	return miso.RLockRun(ec, "user-vault:resource:global", runnable)
}

// global lock for resources
func lockResourceGlobalExec(ec miso.Rail, runnable miso.Runnable) error {
	return miso.RLockExec(ec, "user-vault:resource:global", runnable)
}

// lock for path
func lockPath[T any](ec miso.Rail, pathNo string, runnable miso.LRunnable[T]) (T, error) {
	return miso.RLockRun(ec, "user-vault:path:"+pathNo, runnable)
}

// lock for path
func lockPathExec(ec miso.Rail, pathNo string, runnable miso.Runnable) error {
	return miso.RLockExec(ec, "user-vault:path:"+pathNo, runnable)
}

// lock for role-resource cache
func lockRoleResCache(ec miso.Rail, runnable miso.LRunnable[any]) (any, error) {
	return miso.RLockRun(ec, "user-vault:role:res:cache", runnable)
}
