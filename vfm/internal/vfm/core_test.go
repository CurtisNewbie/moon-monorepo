package vfm

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/rabbit"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	vault "github.com/curtisnewbie/user-vault/api"
	"github.com/sirupsen/logrus"
)

func testUser() common.User {
	return common.User{
		UserNo:   "UE202205142310076187414",
		Username: "zhuangyongj",
	}
}

func corePreTest(t *testing.T) {
	user := "root"
	pw := ""
	db := "vfm"
	host := "localhost"
	port := 3306
	rail := miso.EmptyRail()

	p := mysql.MySQLConnParam{
		User:      user,
		Password:  pw,
		Schema:    db,
		Host:      host,
		Port:      port,
		ConnParam: strings.Join(miso.GetPropStrSlice(mysql.PropMySQLConnParam), "&"),
	}

	if e := mysql.InitMySQL(rail, p); e != nil {
		t.Fatal(e)
	}
	if _, e := redis.InitRedisFromProp(rail); e != nil {
		t.Fatal(e)
	}

	miso.SetProp(rabbit.PropRabbitMqUsername, "guest")
	miso.SetProp(rabbit.PropRabbitMqPassword, "guest")
	if e := rabbit.StartRabbitMqClient(rail); e != nil {
		t.Fatal(e)
	}
	miso.SetProp("client.addr.fstore.host", "localhost")
	miso.SetProp("client.addr.fstore.port", "8084")

	logrus.SetLevel(logrus.DebugLevel)
}

func TestListFilesInVFolder(t *testing.T) {
	corePreTest(t)
	c := miso.EmptyRail()
	var folderNo string = "hfKh3QZSsWjKufZWflqu8jb0n"
	r, e := listFilesInVFolder(c, mysql.GetMySQL(), miso.Paging{}, folderNo, testUser())
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("%+v", r)
}

func TestListFilesSelective(t *testing.T) {
	corePreTest(t)
	c := miso.EmptyRail()
	r, e := listFilesSelective(c, mysql.GetMySQL(), ListFileReq{}, testUser())
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("%+v", r)

	var filename = "head"
	r, e = listFilesSelective(c, mysql.GetMySQL(), ListFileReq{
		Filename: &filename,
	}, testUser())
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("%+v", r)
}

func TestFileExists(t *testing.T) {
	corePreTest(t)
	c := miso.EmptyRail()
	fname := "test-files.zip"
	exist, e := FileExists(c, mysql.GetMySQL(), PreflightCheckReq{Filename: fname}, testUser())
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("%s exists? %v", fname, exist)
}

func TestFindParentFile(t *testing.T) {
	corePreTest(t)
	c := miso.EmptyRail()
	pf, e := FindParentFile(c, mysql.GetMySQL(), FetchParentFileReq{FileKey: "ZZZ718071967023104410314"}, testUser())
	if e != nil {
		t.Fatal(e)
	}
	if pf.FileKey != "ZZZ718222444658688014704" {
		t.Fatalf("Incorrent ParentFileInfo, fileKey: %v, pf: %+v", pf.FileKey, pf)
	}
	t.Logf("%+v", pf)
}

func TestMoveFileToDir(t *testing.T) {
	corePreTest(t)
	c := miso.EmptyRail()
	req := MoveIntoDirReq{
		Uuid: "eb6bc04f-15c5-4f85-a84d-be3d5a7236d8",
		// ParentFileUuid: "5ddf49ca-dec9-4ecf-962d-47b0f3eab90c",
		ParentFileUuid: "",
	}
	e := MoveFileToDir(c, mysql.GetMySQL(), req, testUser())
	if e != nil {
		t.Fatal(e)
	}
}

func TestMakeDir(t *testing.T) {
	corePreTest(t)
	c := miso.EmptyRail()
	fileKey, e := MakeDir(c, mysql.GetMySQL(), MakeDirReq{Name: "mydir"}, testUser())
	if e != nil {
		t.Fatal(e)
	}
	if fileKey == "" {
		t.Fatal("fileKey is empty")
	}
	t.Logf("fileKey: %v", fileKey)
}

func TestCreateVFolder(t *testing.T) {
	corePreTest(t)
	c := miso.EmptyRail()
	r := util.ERand(5)
	folderNo, e := CreateVFolder(c, mysql.GetMySQL(), CreateVFolderReq{"MyFolder_" + r}, testUser())
	if e != nil {
		t.Fatal(e)
	}
	if folderNo == "" {
		t.Fatal("folderNo is empty")
	}

	t.Logf("FolderNo: %v", folderNo)
}

func TestListDirs(t *testing.T) {
	corePreTest(t)
	c := miso.EmptyRail()
	dirs, e := ListDirs(c, mysql.GetMySQL(), testUser())
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("%+v", dirs)
}

func TestShareVFolder(t *testing.T) {
	corePreTest(t)
	if e := ShareVFolder(miso.EmptyRail(), mysql.GetMySQL(),
		vault.UserInfo{Id: 30, Username: "sharon", UserNo: "UE202205142310074386952"}, "hfKh3QZSsWjKufZWflqu8jb0n", testUser()); e != nil {
		t.Fatal(e)
	}
}

func TestRemoveVFolderAccess(t *testing.T) {
	corePreTest(t)
	req := RemoveGrantedFolderAccessReq{
		UserNo:   "UE202303190019399941339",
		FolderNo: "hfKh3QZSsWjKufZWflqu8jb0n",
	}
	if e := RemoveVFolderAccess(miso.EmptyRail(), mysql.GetMySQL(), req, testUser()); e != nil {
		t.Fatal(e)
	}
}

func TestListVFolderBrief(t *testing.T) {
	corePreTest(t)
	v, e := ListVFolderBrief(miso.EmptyRail(), mysql.GetMySQL(), testUser())
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("%+v", v)
}

func TestAddFileToVFolder(t *testing.T) {
	corePreTest(t)
	e := AddFileToVFolder(miso.EmptyRail(), mysql.GetMySQL(),
		AddFileToVfolderReq{
			FolderNo: "hfKh3QZSsWjKufZWflqu8jb0n",
			FileKeys: []string{"ZZZ687250481528832971813"},
			Sync:     true,
		}, testUser())
	if e != nil {
		t.Fatal(e)
	}
}

func TestRemoveFileFromVFolder(t *testing.T) {
	corePreTest(t)
	e := RemoveFileFromVFolder(miso.EmptyRail(), mysql.GetMySQL(),
		RemoveFileFromVfolderReq{
			FolderNo: "hfKh3QZSsWjKufZWflqu8jb0n",
			FileKeys: []string{"ZZZ687250481528832971813"},
		}, testUser())
	if e != nil {
		t.Fatal(e)
	}
}

func TestListVFolders(t *testing.T) {
	corePreTest(t)
	l, e := ListVFolders(miso.EmptyRail(), mysql.GetMySQL(), ListVFolderReq{}, testUser())
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("%+v", l)
}

func TestListGrantedFolderAccess(t *testing.T) {
	corePreTest(t)
	l, e := ListGrantedFolderAccess(miso.EmptyRail(), mysql.GetMySQL(),
		ListGrantedFolderAccessReq{FolderNo: "hfKh3QZSsWjKufZWflqu8jb0n"}, testUser())
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("%+v", l)
}

func TestUpdateFile(t *testing.T) {
	corePreTest(t)
	e := UpdateFile(miso.EmptyRail(), mysql.GetMySQL(), UpdateFileReq{Id: 301, Name: "test-files-222.zip"}, testUser())
	if e != nil {
		t.Fatal(e)
	}
}

func TestCreateFile(t *testing.T) {
	corePreTest(t)
	c := miso.EmptyRail()

	file, err := os.ReadFile("../README.md")
	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.NewBuffer(file)

	var r miso.GnResp[string]
	err = miso.NewDynTClient(c, "/file", "fstore").
		AddHeader("filename", "README.md").
		Put(buf).
		Json(&r)
	if err != nil {
		t.Fatal(err)
	}

	if err := r.Err(); err != nil {
		t.Fatal(err)
	}

	fakeFileId := r.Data
	c.Infof("fake fileId: %v", fakeFileId)

	_, e := CreateFile(c, mysql.GetMySQL(), CreateFileReq{
		Filename:         "myfile",
		FakeFstoreFileId: fakeFileId,
	}, testUser())
	if e != nil {
		t.Fatal(e)
	}
}

func TestDeleteFile(t *testing.T) {
	corePreTest(t)
	c := miso.EmptyRail()
	e := DeleteFile(c, mysql.GetMySQL(), DeleteFileReq{Uuid: "ZZZ718078073798656022858"}, testUser(), nil)
	if e != nil {
		t.Fatal(e)
	}
}

func TestGenTempToken(t *testing.T) {
	corePreTest(t)
	c := miso.EmptyRail()
	tkn, e := GenTempToken(c, mysql.GetMySQL(), GenerateTempTokenReq{"ZZZ687250496077824971813"}, testUser())
	if e != nil {
		t.Fatal(e)
	}
	if tkn == "" {
		t.Fatal("Token is empty")
	}
	t.Logf("tkn: %v", tkn)
}

func TestIsImage(t *testing.T) {
	n := "abc.jpg"
	if !isImage(n) {
		t.Fatal(n)
	}

	n = "abc.txt"
	if isImage(n) {
		t.Fatal(n)
	}
}

func TestBatchClacDirSize(t *testing.T) {
	corePreTest(t)
	c := miso.EmptyRail()
	err := ImMemBatchCalcDirSize(c, mysql.GetMySQL())
	if err != nil {
		t.Fatal(err)
	}
}

func TestUnpackZip(t *testing.T) {
	corePreTest(t)
	rail := miso.EmptyRail()
	req := UnpackZipReq{
		FileKey:       "ZZZ1065471829557248604128",
		ParentFileKey: "",
	}
	err := UnpackZip(rail, mysql.GetMySQL(), testUser(), req)
	if err != nil {
		t.Fatal(err)
	}
}
