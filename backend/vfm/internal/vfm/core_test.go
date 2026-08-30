package vfm

import (
	"bytes"
	"container/list"
	"os"
	"strings"
	"testing"

	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/rabbit"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/randutil"
	"github.com/curtisnewbie/miso/util/strutil"
	vault "github.com/curtisnewbie/user-vault/api"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func testUser() flow.User {
	return flow.User{
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
	exist, e := FileExists(c, mysql.GetMySQL(), PreflightCheckReq{Filename: fname}, testUser().UserNo)
	if e != nil {
		t.Fatal(e)
	}
	t.Logf("%s exists? %v", fname, exist)
}

func TestFindParentFile(t *testing.T) {
	corePreTest(t)
	c := miso.EmptyRail()
	pf, _, e := FindParentFile(c, mysql.GetMySQL(), FetchParentFileReq{FileKey: "ZZZ718071967023104410314"}, testUser())
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
	r := randutil.ERand(5)
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
	err = miso.NewDynClient(c, "/file", "fstore").
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

func TestFetchDirTreeBottomUp(t *testing.T) {
	corePreTest(t)
	rail := miso.EmptyRail()
	dirParentCache.DelAll(rail)
	n, err := FetchDirTreeBottomUp(rail, mysql.GetMySQL(), FetchDirTreeReq{FileKey: "ZZZ1471280777216000148288"}, flow.User{})
	if err != nil {
		t.Fatal(err)
	}
	if n == nil {
		t.Fatal("node is nil")
	}
	for n != nil {
		t.Logf("n: %#v", n)
		n = n.Child
	}
}

func TestFetchDirTreeTopDown(t *testing.T) {
	corePreTest(t)
	rail := miso.EmptyRail()
	dirParentCache.DelAll(rail)
	root, err := FetchDirTreeTopDown(rail, mysql.GetMySQL(), flow.User{UserNo: "UE1049787455160320075953"})
	if err != nil {
		t.Fatal(err)
	}
	if root == nil {
		t.Fatal("node is nil")
	}

	l := list.List{}
	l.PushFront(root)
	d := 1
	for l.Len() > 0 {
		cnt := l.Len()
		for i := 0; i < cnt; i++ {
			front := l.Front()
			l.Remove(front)
			n := front.Value.(*DirTopDownTreeNode)
			if n.FileKey == "" {
				t.Logf("%v /", strutil.Tabs(d))
			} else {
				t.Logf("%v /%v", strutil.Tabs(d), n.Name)
			}
			for i := range n.Child {
				c := n.Child[i]
				l.PushBack(c)
			}
		}
		d++
	}
}

// testMySQLConn creates a dedicated MySQL connection to the vfm schema.
// The global connection (mysql.GetMySQL) is unreliable in tests: bookmark_test.go
// initializes it against the docindexer schema and InitMySQL is idempotent.
func testMySQLConn(t *testing.T) *gorm.DB {
	t.Helper()
	rail := miso.EmptyRail()
	p := mysql.MySQLConnParam{
		User:      "root",
		Password:  "",
		Schema:    "vfm",
		Host:      "localhost",
		Port:      3306,
		ConnParam: strings.Join(miso.GetPropStrSlice(mysql.PropMySQLConnParam), "&"),
	}
	conn, err := mysql.NewMySQLConn(rail, p)
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

// insertTestDir inserts a dir row and returns its uuid
func insertTestDir(t *testing.T, db *gorm.DB, userNo string) string {
	t.Helper()
	uuid := "tdt-dir-" + randutil.ERand(8)
	if err := db.Exec("INSERT INTO file_info (name, uuid, size_in_bytes, file_type, parent_file, thumbnail, uploader_no) VALUES (?,?,?,?,?,?,?)",
		"test-dir", uuid, 0, "DIR", "", "", userNo).Error; err != nil {
		t.Fatal(err)
	}
	return uuid
}

// insertTestChild inserts a child file row under parent, returns its uuid
func insertTestChild(t *testing.T, db *gorm.DB, parent string, name string, thumbnail string, isDel int) string {
	t.Helper()
	uuid := "tdt-child-" + randutil.ERand(8)
	if err := db.Exec("INSERT INTO file_info (name, uuid, size_in_bytes, file_type, parent_file, thumbnail, uploader_no, is_del) VALUES (?,?,?,?,?,?,?,?)",
		name, uuid, 0, "FILE", parent, thumbnail, testUser().UserNo, isDel).Error; err != nil {
		t.Fatal(err)
	}
	return uuid
}

func TestFindFirstThumbnailFileId(t *testing.T) {
	rail := miso.EmptyRail()
	db := testMySQLConn(t)

	dirKey := insertTestDir(t, db, testUser().UserNo)
	t.Cleanup(func() {
		db.Exec("DELETE FROM file_info WHERE parent_file = ? OR uuid = ?", dirKey, dirKey)
	})

	// children inserted oldest → newest
	insertTestChild(t, db, dirKey, "a.jpg", "", 0)          // oldest, no thumbnail
	insertTestChild(t, db, dirKey, "b.jpg", "THUMB_OLD", 0) // has thumbnail
	insertTestChild(t, db, dirKey, "c.jpg", "", 0)          // no thumbnail
	insertTestChild(t, db, dirKey, "d.jpg", "THUMB_NEW", 0) // newest with thumbnail
	insertTestChild(t, db, dirKey, "e.jpg", "THUMB_DEL", 1) // deleted, must be skipped

	fst, ok, err := findFirstThumbnailFileId(rail, db, dirKey)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected to find a thumbnail")
	}
	if fst != "THUMB_OLD" {
		t.Fatalf("expected oldest thumbnail THUMB_OLD, got %v", fst)
	}
}

func TestFindFirstThumbnailFileIdNoThumbnail(t *testing.T) {
	rail := miso.EmptyRail()
	db := testMySQLConn(t)

	dirKey := insertTestDir(t, db, testUser().UserNo)
	t.Cleanup(func() {
		db.Exec("DELETE FROM file_info WHERE parent_file = ? OR uuid = ?", dirKey, dirKey)
	})

	insertTestChild(t, db, dirKey, "a.jpg", "", 0)
	insertTestChild(t, db, dirKey, "b.jpg", "", 0)

	fst, ok, err := findFirstThumbnailFileId(rail, db, dirKey)
	if err != nil {
		t.Fatal(err)
	}
	if ok || fst != "" {
		t.Fatalf("expected no thumbnail, got ok=%v fst=%v", ok, fst)
	}
}

func TestFindFirstThumbnailFileIds(t *testing.T) {
	rail := miso.EmptyRail()
	db := testMySQLConn(t)

	dir1 := insertTestDir(t, db, testUser().UserNo)
	dir2 := insertTestDir(t, db, testUser().UserNo)
	dir3 := insertTestDir(t, db, testUser().UserNo) // no thumbnails
	dir4 := insertTestDir(t, db, "UE_OTHER_USER")   // owned by someone else
	t.Cleanup(func() {
		db.Exec("DELETE FROM file_info WHERE parent_file IN (?,?,?,?) OR uuid IN (?,?,?,?)", dir1, dir2, dir3, dir4, dir1, dir2, dir3, dir4)
	})

	// dir1: oldest has thumbnail
	insertTestChild(t, db, dir1, "a.jpg", "", 0)
	insertTestChild(t, db, dir1, "b.jpg", "D1_OLD", 0)
	insertTestChild(t, db, dir1, "c.jpg", "D1_NEW", 0)
	// dir2: oldest two lack thumbnails
	insertTestChild(t, db, dir2, "x.jpg", "", 0)
	insertTestChild(t, db, dir2, "y.jpg", "D2_OLD", 0)
	insertTestChild(t, db, dir2, "z.jpg", "D2_NEW", 0)
	// dir4: has thumbnail but dir owned by another user
	insertTestChild(t, db, dir4, "w.jpg", "D4_THUMB", 0)

	m, err := findFirstThumbnailFileIds(rail, db, []string{dir1, dir2, dir3, dir4}, testUser().UserNo)
	if err != nil {
		t.Fatal(err)
	}
	if m[dir1] != "D1_OLD" {
		t.Fatalf("expected dir1 oldest thumbnail D1_OLD, got %v", m[dir1])
	}
	if m[dir2] != "D2_OLD" {
		t.Fatalf("expected dir2 oldest thumbnail D2_OLD, got %v", m[dir2])
	}
	if _, ok := m[dir3]; ok {
		t.Fatalf("dir3 should have no thumbnail, got %v", m[dir3])
	}
	if _, ok := m[dir4]; ok {
		t.Fatalf("dir4 owned by another user should be excluded, got %v", m[dir4])
	}
}
