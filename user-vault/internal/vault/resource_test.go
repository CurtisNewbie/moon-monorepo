package vault

import (
	"testing"

	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/user-vault/api"
)

func before(t *testing.T) {
	rail := miso.EmptyRail()
	miso.LoadConfigFromFile("../app-conf-dev.yml", rail)
	if _, e := miso.InitRedisFromProp(rail); e != nil {
		t.Fatal(e)
	}
	if e := miso.InitMySQLFromProp(rail); e != nil {
		t.Fatal(e)
	}
}

func TestUpdatePath(t *testing.T) {
	before(t)

	req := UpdatePathReq{
		PathNo: "path_578477630062592208429",
		Type:   PathTypePublic,
		Group:  "goauth",
	}
	e := UpdatePath(miso.EmptyRail(), req)
	if e != nil {
		t.Fatal(e)
	}
}

func TestGetRoleInfo(t *testing.T) {
	before(t)

	req := api.RoleInfoReq{
		RoleNo: "role_554107924873216177918",
	}
	resp, e := GetRoleInfo(miso.EmptyRail(), req)
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("%v", resp)
}

func TestCreatePathIfNotExist(t *testing.T) {
	before(t)

	req := CreatePathReq{
		Type:  PathTypeProtected,
		Url:   "/goauth/open/api/role/resource/add",
		Group: "goauth",
	}
	e := CreatePath(miso.EmptyRail(), req, common.NilUser())
	if e != nil {
		t.Fatal(e)
	}
}

func TestDeletePath(t *testing.T) {
	before(t)

	req := DeletePathReq{
		PathNo: "path_555305367076864208429",
	}

	e := DeletePath(miso.EmptyRail(), req)
	if e != nil {
		t.Fatal(e)
	}
}

func TestCreateRes(t *testing.T) {
	before(t)

	req := CreateResReq{
		Name: "GoAuth Test  ",
	}

	e := CreateResourceIfNotExist(miso.EmptyRail(), req, common.NilUser())
	if e != nil {
		t.Fatal(e)
	}
}

func TestBindPathRes(t *testing.T) {
	before(t)

	req := BindPathResReq{
		PathNo:  "path_555326806016000208429",
		ResCode: "res_555323073019904208429",
	}

	e := BindPathRes(miso.EmptyRail(), req)
	if e != nil {
		t.Fatal(e)
	}
}

func TestPreprocessUrl(t *testing.T) {
	if v := preprocessUrl(""); v != "/" {
		t.Fatal(v)
	}

	if v := preprocessUrl("/"); v != "/" {
		t.Fatal(v)
	}

	if v := preprocessUrl("///"); v != "/" {
		t.Fatal(v)
	}

	if v := preprocessUrl("/goauth/test/path"); v != "/goauth/test/path" {
		t.Fatal(v)
	}

	if v := preprocessUrl("/goauth/test/path//"); v != "/goauth/test/path" {
		t.Fatal(v)
	}

	if v := preprocessUrl("goauth/test/path//"); v != "/goauth/test/path" {
		t.Fatal(v)
	}

	if v := preprocessUrl("goauth/test/path?abc=123"); v != "/goauth/test/path" {
		t.Fatal(v)
	}
}

func TestUnbindPathRes(t *testing.T) {
	before(t)

	req := UnbindPathResReq{
		PathNo: "path_555326806016000208429",
	}

	e := UnbindPathRes(miso.EmptyRail(), req)
	if e != nil {
		t.Fatal(e)
	}
}

func TestAddRole(t *testing.T) {
	before(t)

	req := AddRoleReq{
		Name: "Guest",
	}

	e := AddRole(miso.EmptyRail(), req, common.NilUser())
	if e != nil {
		t.Fatal(e)
	}
}

func TestAddResToRole(t *testing.T) {
	before(t)

	req := AddRoleResReq{
		RoleNo:  "role_555329954676736208429",
		ResCode: "res_555323073019904208429",
	}

	e := AddResToRoleIfNotExist(miso.EmptyRail(), req, common.NilUser())
	if e != nil {
		t.Fatal(e)
	}
}

func TestGenPathNo(t *testing.T) {
	pathNo := genPathNo("test", "/core/path/is/that/okay/if/i/amy/very", "GET")
	if pathNo == "" {
		t.Error("pathNo is empty")
		return
	}
	t.Log(pathNo)
}

func TestRemoveResFromRole(t *testing.T) {
	before(t)

	req := RemoveRoleResReq{
		RoleNo:  "role_555329954676736208429",
		ResCode: "res_555323073019904208429",
	}

	e := RemoveResFromRole(miso.EmptyRail(), req)
	if e != nil {
		t.Fatal(e)
	}
}

func TestListRoleRes(t *testing.T) {
	before(t)

	p := miso.Paging{
		Limit: 5,
		Page:  1,
	}
	req := ListRoleResReq{
		RoleNo: "role_555329954676736208429",
		Paging: p,
	}

	resp, e := ListRoleRes(miso.EmptyRail(), req)
	if e != nil {
		t.Fatal(e)
	}

	if resp.Paging.Total < 1 {
		t.Fatal("total < 1")
	}

	t.Logf("%+v", resp)
}

func TestListAllRoleBriefs(t *testing.T) {
	before(t)

	resp, e := ListAllRoleBriefs(miso.EmptyRail())
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("%+v", resp)
}

func TestListRoles(t *testing.T) {
	before(t)

	p := miso.Paging{
		Limit: 5,
		Page:  1,
	}
	req := ListRoleReq{
		Paging: p,
	}

	resp, e := ListRoles(miso.EmptyRail(), req)
	if e != nil {
		t.Fatal(e)
	}

	if resp.Paging.Total < 1 {
		t.Fatal("total < 1")
	}

	t.Logf("%+v", resp)
}

func TestTestResourceAccess(t *testing.T) {
	before(t)

	ec := miso.EmptyRail()
	LoadPathResCache(ec)
	LoadRoleResCache(ec)

	req := TestResAccessReq{
		RoleNo: "role_555329954676736208429",
		Url:    "/goauth/open/api/role/resource/add",
	}

	r, e := TestResourceAccess(ec, req)
	if e != nil {
		t.Fatal(e)
	}
	if !r.Valid {
		t.Fatal("should be valid")
	}
}
