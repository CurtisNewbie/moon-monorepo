package vault

import (
	"crypto/md5"
	"encoding/base64"
	"strings"
	"time"

	doublestar "github.com/bmatcuk/doublestar/v4"
	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/auth"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/atom"
	"github.com/curtisnewbie/miso/util/slutil"
	"github.com/curtisnewbie/user-vault/api"
	"github.com/curtisnewbie/user-vault/internal/repo"
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

type WRes struct {
	Id         int       `json:"id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	CreateTime atom.Time `json:"createTime" gorm:"column:created_at"`
	CreateBy   string    `json:"createBy" gorm:"column:created_by"`
	UpdateTime atom.Time `json:"updateTime" gorm:"column:updated_at"`
	UpdateBy   string    `json:"updateBy" gorm:"column:updated_by"`
}

type UnbindPathResReq struct {
	PathNo  string `json:"pathNo" validation:"notEmpty"`
	ResCode string `json:"resCode" validation:"notEmpty"`
}

type ListResReq struct {
	Paging miso.Paging `json:"paging"`
}

type ListResResp struct {
	Paging  miso.Paging `json:"paging"`
	Payload []WRes      `json:"payload"`
}

func DeleteResource(rail miso.Rail, req repo.DeleteResourceReq) error {
	_, err := lockResourceGlobal(rail, func() (any, error) {
		return nil, repo.DeleteResource(rail, req)
	})
	return err
}

func ListResourceCandidatesForRole(rail miso.Rail, roleNo string) ([]repo.ResBrief, error) {
	return repo.ListResourceCandidatesForRole(rail, roleNo)
}

func ListAllResBriefsOfRole(rail miso.Rail, roleNo string) ([]repo.ResBrief, error) {
	return repo.ListAllResBriefsOfRole(rail, roleNo)
}

func ListAllResBriefs(rail miso.Rail) ([]repo.ResBrief, error) {
	return repo.ListAllResBriefs(rail)
}

func ListResources(rail miso.Rail, req repo.ListResReq) (repo.ListResResp, error) {
	return repo.ListResources(rail, req)
}

func UpdatePath(rail miso.Rail, req repo.UpdatePathReq) error {
	_, e := lockPath(rail, req.PathNo, func() (any, error) {
		return nil, repo.UpdatePath(rail, req)
	})
	return e
}

func GetRoleInfo(rail miso.Rail, req api.RoleInfoReq) (api.RoleInfoResp, error) {
	resp, err := roleInfoCache.GetValElse(rail, req.RoleNo, func() (api.RoleInfoResp, error) {
		return repo.GetRoleInfo(rail, req)
	})
	return resp, err
}

func CreateResourceIfNotExist(rail miso.Rail, req repo.CreateResReq, user flow.User) error {
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(req.Code)

	_, e := lockResourceGlobal(rail, func() (any, error) {
		return nil, repo.CreateResourceIfNotExist(rail, req, user)
	})
	return e
}

func genPathNo(group string, url string, method string) string {
	cksum := md5.Sum([]byte(group + method + url))
	return "path_" + base64.StdEncoding.EncodeToString(cksum[:])
}

func CreatePath(rail miso.Rail, req repo.CreatePathReq, user flow.User) error {
	req.Url = preprocessUrl(req.Url)
	req.Group = strings.TrimSpace(req.Group)
	req.Method = strings.ToUpper(strings.TrimSpace(req.Method))
	pathNo := genPathNo(req.Group, req.Url, req.Method)

	_, err := lockPath(rail, pathNo, func() (any, error) {
		return nil, repo.CreatePath(rail, req, pathNo, user)
	})
	if err != nil {
		return err
	}

	if req.ResCode != "" { // rebind path and resource
		return BindPathRes(rail, repo.BindPathResReq{PathNo: pathNo, ResCode: req.ResCode})
	}

	return nil
}

func DeletePath(rail miso.Rail, req repo.DeletePathReq) error {
	req.PathNo = strings.TrimSpace(req.PathNo)
	_, e := lockPath(rail, req.PathNo, func() (any, error) {
		return nil, repo.DeletePath(rail, req)
	})
	return e
}

func UnbindPathRes(rail miso.Rail, req repo.UnbindPathResReq) error {
	req.PathNo = strings.TrimSpace(req.PathNo)
	_, e := lockPath(rail, req.PathNo, func() (any, error) {
		return nil, repo.UnbindPathRes(rail, req)
	})
	return e
}

func BindPathRes(rail miso.Rail, req repo.BindPathResReq) error {
	req.PathNo = strings.TrimSpace(req.PathNo)
	e := lockPathExec(rail, req.PathNo, func() error { // lock for path
		return lockResourceGlobalExec(rail, func() error {
			return repo.BindPathRes(rail, dbquery.GetDB(), req)
		})
	})

	return e
}

func ListPaths(rail miso.Rail, req repo.ListPathReq) (repo.ListPathResp, error) {
	return repo.ListPaths(rail, req)
}

func AddRole(rail miso.Rail, req repo.AddRoleReq, user flow.User) error {
	_, e := redis.RLockRun(rail, "user-vault:role:add"+req.Name, func() (any, error) {
		return nil, repo.AddRole(rail, req, user)
	})
	return e
}

func RemoveResFromRole(rail miso.Rail, req repo.RemoveRoleResReq) error {
	_, e := redis.RLockRun(rail, "user-vault:role:"+req.RoleNo, func() (any, error) {
		return nil, repo.RemoveResFromRole(rail, req)
	})

	return e
}

func AddResToRoleIfNotExist(rail miso.Rail, req repo.AddRoleResReq, user flow.User) error {

	_, e := redis.RLockRun(rail, "user-vault:role:"+req.RoleNo, func() (any, error) { // lock for role
		return lockResourceGlobal(rail, func() (any, error) {
			return nil, repo.AddResToRoleIfNotExist(rail, req, user)
		})
	})
	return e
}

func ListRoleRes(rail miso.Rail, req repo.ListRoleResReq) (repo.ListRoleResResp, error) {
	return repo.ListRoleRes(rail, req)
}

func ListAllRoleBriefs(rail miso.Rail) ([]repo.RoleBrief, error) {
	return repo.ListAllRoleBriefs(rail)
}

func ListRoles(rail miso.Rail, req repo.ListRoleReq) (repo.ListRoleResp, error) {
	return repo.ListRoles(rail, req)
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

		lr, e := repo.ListRoleNos(rail)
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
	var paths []repo.ExtendedPathRes
	if isDefAdmin(roleNo) {
		p, err := repo.ListAllPathRes(rail, dbquery.GetDB())
		if err != nil {
			return err
		}
		paths = p
	} else {
		p, err := repo.ListRolePathRes(rail, dbquery.GetDB(), roleNo)
		if err != nil {
			return err
		}
		paths = p
	}
	if paths == nil {
		return nil
	}

	public, err := repo.ListPublicPathRes(rail, dbquery.GetDB())
	if err != nil {
		return err
	}
	paths = append(paths, public...)

	var pai []PathAccessInfo = slutil.MapTo(paths,
		func(t repo.ExtendedPathRes) PathAccessInfo {
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
	public, err := repo.ListPublicPathRes(rail, dbquery.GetDB())
	if err != nil {
		return err
	}
	if public == nil {
		return nil
	}

	var pai []PathAccessInfo = slutil.MapTo(public,
		func(t repo.ExtendedPathRes) PathAccessInfo {
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

		user := flow.NilUser()

		app := miso.GetPropStr(miso.PropAppName)
		for _, res := range res {
			if res.Code == "" || res.Name == "" {
				continue
			}
			if e := CreateResourceIfNotExist(rail, repo.CreateResReq(res), user); e != nil {
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

			r := repo.CreatePathReq{
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
