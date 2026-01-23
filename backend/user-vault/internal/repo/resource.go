package repo

import (
	"strings"

	"github.com/curtisnewbie/miso/errs"
	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/atom"
	"github.com/curtisnewbie/miso/util/snowflake"
	"github.com/curtisnewbie/user-vault/api"
	"gorm.io/gorm"
)

var ErrCodeRoleNotFound = "GA0001"

const (
	// default roleno for admin
	DefaultAdminRoleNo  = "role_554107924873216177918" // deprecated
	DefaultAdminRoleNo2 = "role_super_admin"

	PathTypeProtected string = "PROTECTED"
	PathTypePublic    string = "PUBLIC"
)

type ExtendedPathRes struct {
	Id         int       // id
	Pgroup     string    // path group
	PathNo     string    // path no
	ResCode    string    // resource code
	Desc       string    // description
	Url        string    // url
	Method     string    // http method
	Ptype      string    // path type: PROTECTED, PUBLIC
	CreateTime atom.Time `gorm:"column:created_at"`
	CreateBy   string    `gorm:"column:created_by"`
	UpdateTime atom.Time `gorm:"column:updated_at"`
	UpdateBy   string    `gorm:"column:updated_by"`
}

type WRole struct {
	Id         int       `json:"id"`
	RoleNo     string    `json:"roleNo"`
	Name       string    `json:"name"`
	CreateTime atom.Time `json:"createTime" gorm:"column:created_at"`
	CreateBy   string    `json:"createBy" gorm:"column:created_by"`
	UpdateTime atom.Time `json:"updateTime" gorm:"column:updated_at"`
	UpdateBy   string    `json:"updateBy" gorm:"column:updated_by"`
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
	Id         int       `json:"id"`
	Pgroup     string    `json:"pgroup"`
	PathNo     string    `json:"pathNo"`
	Method     string    `json:"method"`
	Desc       string    `json:"desc"`
	Url        string    `json:"url"`
	Ptype      string    `json:"ptype" desc:"path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible"`
	CreateTime atom.Time `json:"createTime" gorm:"column:created_at"`
	CreateBy   string    `json:"createBy" gorm:"column:created_by"`
	UpdateTime atom.Time `json:"updateTime" gorm:"column:updated_at"`
	UpdateBy   string    `json:"updateBy" gorm:"column:updated_by"`
}

type WRes struct {
	Id         int       `json:"id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	CreateTime atom.Time `json:"createTime" gorm:"column:created_at"`
	CreateBy   string    `json:"createBy" gorm:"column:created_by"`
	UpdateTime atom.Time `json:"updateTime" gorm:"column:updated_at"`
	UpdateBy   string    `json:"updateBy" gorm:"column:updated_by"`
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
	Id         int       `json:"id"`
	ResCode    string    `json:"resCode"`
	ResName    string    `json:"resName"`
	CreateTime atom.Time `json:"createTime" gorm:"column:created_at"`
	CreateBy   string    `json:"createBy" gorm:"column:created_by"`
}

type GenResScriptReq struct {
	ResCodes []string `json:"resCodes" validation:"notEmpty"`
}

type UpdatePathReq struct {
	Type    string `valid:"notEmpty,member:PROTECTED|PUBLIC" desc:"path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible" json:"type"`
	PathNo  string `valid:"notEmpty" json:"pathNo"`
	Group   string `valid:"notEmpty,maxLen:20" json:"group"`
	ResCode string `json:"resCode"`
}

type CreatePathReq struct {
	Type    string `valid:"notEmpty,member:PROTECTED|PUBLIC" desc:"path type: 'PROTECTED' - authorization required, 'PUBLIC' - publicly accessible" json:"type"`
	Url     string `valid:"notEmpty,maxLen:128" json:"url"`
	Group   string `valid:"notEmpty,maxLen:20" json:"group"`
	Method  string `valid:"notEmpty,maxLen:10" json:"method"`
	Desc    string `valid:"maxLen:255" json:"desc"`
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
	return dbquery.RunTransaction(rail, dbquery.GetDB(), func(qry func() *dbquery.Query) error {
		q := qry()
		if err := q.ExecAny(`delete from resource where code = ?`, req.ResCode); err != nil {
			return err
		}
		if err := q.ExecAny(`delete from role_resource where res_code = ?`, req.ResCode); err != nil {
			return err
		}
		return q.ExecAny(`delete from path_resource where res_code = ?`, req.ResCode)
	})
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

	if IsDefAdmin(roleNo) {
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

	count, err := dbquery.NewQuery(rail).Table("resource").Count()
	if err != nil {
		return ListResResp{}, err
	}

	return ListResResp{Paging: miso.RespPage(req.Paging, int(count)), Payload: resources}, nil
}

func UpdatePath(rail miso.Rail, req UpdatePathReq) error {
	return mysql.GetMySQL().Transaction(func(tx *gorm.DB) error {
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
}

func GetRoleInfo(rail miso.Rail, req api.RoleInfoReq) (api.RoleInfoResp, error) {
	var resp api.RoleInfoResp
	n, err := dbquery.NewQuery(rail).
		Table("role").
		Eq("role_no", req.RoleNo).
		SelectCols(resp).
		Scan(&resp)
	if err != nil {
		return resp, err
	}
	if n < 1 {
		return resp, errs.NewErrf("Role not found").WithCode(ErrCodeRoleNotFound)
	}
	return resp, nil
}

func CreateResourceIfNotExist(rail miso.Rail, req CreateResReq, user flow.User) error {
	ok, err := dbquery.NewQuery(rail).
		Table("resource").
		Eq("code", req.Code).
		HasAny()
	if err != nil {
		return err
	}

	if ok {
		rail.Debugf("Resource '%s' (%s) already exist", req.Code, req.Name)
		return nil
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
	return err
}

func CreatePath(rail miso.Rail, req CreatePathReq, pathNo string, user flow.User) error {

	type path struct {
		Id     int
		Pgroup string
		PathNo string
		Desc   string
		Url    string
		Method string
		Ptype  string
	}

	db := mysql.GetMySQL()
	var prev path
	ok, err := dbquery.NewQuery(rail, db).
		Table("path").
		Eq("path_no", pathNo).
		ScanAny(&prev)
	if err != nil {
		return err
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
				return err
			}
		}
		return nil
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
		return err
	}

	rail.Infof("Created path (%s) '%s {%s}'", pathNo, req.Method, req.Url)
	return nil
}

func DeletePath(rail miso.Rail, req DeletePathReq) error {
	return mysql.GetMySQL().Transaction(func(tx *gorm.DB) error {
		_, err := dbquery.NewQuery(rail, tx).Exec(`delete from path where path_no = ?`, req.PathNo)
		if err != nil {
			return err
		}

		_, err = dbquery.NewQuery(rail, tx).Exec(`delete from path_resource where path_no = ?`, req.PathNo)
		return err
	})
}

func UnbindPathRes(rail miso.Rail, req UnbindPathResReq) error {
	req.PathNo = strings.TrimSpace(req.PathNo)
	return dbquery.NewQuery(rail).
		ExecAny(`delete from path_resource where path_no = ?`, req.PathNo)
}

func BindPathRes(rail miso.Rail, db *gorm.DB, req BindPathResReq) error {
	// check if resource exist
	ok, err := dbquery.NewQuery(rail, db).
		Table("resource").Eq("code", req.ResCode).
		HasAny()
	if err != nil {
		return err
	}
	if !ok {
		rail.Errorf("Resource %v not found", req.ResCode)
		return errs.NewErrf("Resource not found")
	}

	// check if the path is already bound to current resource
	ok, err = dbquery.NewQuery(rail, db).
		Table("path_resource").
		Eq("path_no", req.PathNo).
		Eq("res_code", req.ResCode).
		HasAny()

	if err != nil {
		rail.Errorf("Failed to bind path %v to resource %v, %v", req.PathNo, req.ResCode, err)
		return err
	}
	if ok {
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

	count, err := applyCond(dbquery.NewQuery(rail).
		Table("path p")).
		Count()
	if err != nil {
		return ListPathResp{}, err
	}

	return ListPathResp{
		Payload: paths,
		Paging:  miso.Paging{Limit: req.Paging.Limit, Page: req.Paging.Page, Total: int(count)},
	}, nil
}

func AddRole(rail miso.Rail, req AddRoleReq, user flow.User) error {
	return dbquery.NewQuery(rail).
		Table("role").
		CreateAny(struct {
			RoleNo string
			Name   string
		}{
			RoleNo: snowflake.IdPrefix("role_"),
			Name:   req.Name,
		})
}

func RemoveResFromRole(rail miso.Rail, req RemoveRoleResReq) error {
	_, err := dbquery.NewQuery(rail, dbquery.GetDB()).
		Exec(`delete from role_resource where role_no = ? and res_code = ?`, req.RoleNo, req.ResCode)
	return err
}

func AddResToRoleIfNotExist(rail miso.Rail, req AddRoleResReq, user flow.User) error {
	// check if resource exist
	var resId int
	_, err := dbquery.NewQuery(rail).
		Raw(`select id from resource where code = ?`, req.ResCode).
		Scan(&resId)
	if err != nil {
		return err
	}
	if resId < 1 {
		return errs.NewErrf("Resource not found")
	}

	// check if role-resource relation exists
	ok, err := dbquery.NewQuery(rail).
		Table("role_resource").
		Eq("role_no", req.RoleNo).
		Eq("res_code", req.ResCode).
		HasAny()
	if err != nil {
		return err
	}
	if ok { // relation exists already
		return nil
	}

	// create role-resource relation
	rr := struct {
		RoleNo  string // role no
		ResCode string // resource code
	}{
		RoleNo:  req.RoleNo,
		ResCode: req.ResCode,
	}
	err = dbquery.NewQuery(rail).
		Table("role_resource").
		CreateAny(&rr)
	return err
}

func ListRoleRes(rail miso.Rail, req ListRoleResReq) (ListRoleResResp, error) {
	var res []ListedRoleRes
	err := dbquery.NewQuery(rail).
		Table("role_resource rr").
		Joins("LEFT JOIN resource r on rr.res_code = r.code").
		Select(`rr.id, rr.res_code, rr.created_at, rr.created_by, r.name 'res_name'`).
		Eq("rr.role_no", req.RoleNo).
		OrderDesc("rr.id").
		Offset(req.Paging.GetOffset()).
		Limit(req.Paging.GetLimit()).
		ScanVal(&res)

	if err != nil {
		return ListRoleResResp{}, err
	}

	if res == nil {
		res = []ListedRoleRes{}
	}

	count, err := dbquery.NewQuery(rail).
		Table("role_resource rr").
		Joins("LEFT JOIN resource r on rr.res_code = r.code").
		Eq("rr.role_no", req.RoleNo).
		Count()

	if err != nil {
		return ListRoleResResp{}, err
	}

	return ListRoleResResp{Payload: res,
		Paging: miso.Paging{Limit: req.Paging.Limit, Page: req.Paging.Page, Total: int(count)}}, nil
}

func ListAllRoleBriefs(rail miso.Rail) ([]RoleBrief, error) {
	var roles []RoleBrief
	err := dbquery.NewQuery(rail).
		Table("role").
		SelectCols(RoleBrief{}).
		ScanVal(&roles)
	if err != nil {
		return nil, err
	}
	if roles == nil {
		roles = []RoleBrief{}
	}
	return roles, nil
}

func ListRoles(rail miso.Rail, req ListRoleReq) (ListRoleResp, error) {
	var roles []WRole
	err := dbquery.NewQuery(rail).
		Table("role").
		SelectCols(WRole{}).
		OrderDesc("id").
		AtPage(req.Paging).
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

func ListRoleNos(rail miso.Rail) ([]string, error) {
	var ern []string
	err := dbquery.NewQuery(rail).
		Table("role").
		Select("role_no").
		ScanVal(&ern)
	if err != nil {
		return nil, err
	}

	if ern == nil {
		ern = []string{}
	}
	return ern, nil
}

func IsDefAdmin(roleNo string) bool {
	return roleNo == DefaultAdminRoleNo || roleNo == DefaultAdminRoleNo2
}

func ListAllPathRes(rail flow.Rail, db *gorm.DB) ([]ExtendedPathRes, error) {
	var paths []ExtendedPathRes
	err := dbquery.NewQuery(rail, db).
		Raw(`SELECT p.*, pr.res_code
				FROM path_resource pr
				LEFT JOIN path p ON p.path_no = pr.path_no`).
		ScanVal(&paths)
	return paths, err
}

func ListRolePathRes(rail miso.Rail, db *gorm.DB, roleNo string) ([]ExtendedPathRes, error) {
	var paths []ExtendedPathRes
	err := dbquery.NewQuery(rail, db).
		Raw(`SELECT p.*, pr.res_code
		FROM role_resource rr
		LEFT JOIN path_resource pr ON rr.res_code = pr.res_code
		LEFT JOIN path p ON p.path_no = pr.path_no
		WHERE rr.role_no = ?
		`, roleNo).
		ScanVal(&paths)
	if err != nil {
		return nil, err
	}

	return paths, err
}

func ListPublicPathRes(rail miso.Rail, db *gorm.DB) ([]ExtendedPathRes, error) {
	var public []ExtendedPathRes
	err := dbquery.NewQuery(rail, db).
		Raw(`SELECT p.* FROM path p WHERE p.ptype = ?`, PathTypePublic).
		ScanVal(&public)
	return public, err
}
