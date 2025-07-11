package vfm

import (
	"errors"
	"fmt"
	"strings"
	"time"

	fstore "github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/miso/encoding/json"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	vault "github.com/curtisnewbie/user-vault/api"
	"gorm.io/gorm"
)

const (
	FileTypeFile = "FILE" // file
	FileTypeDir  = "DIR"  // directory

	DelN = 0 // normal file
	DelY = 1 // for logic delete: file marked deleted; for physic delete: file may be removed from disk or move to somewhere else.

	VfolderOwner   = "OWNER"   // owner of the vfolder
	VfolderGranted = "GRANTED" // granted access to the vfolder
)

var (
	_imageSuffix = util.NewSet[string]()
	_videoSuffix = util.NewSet[string]()

	dirParentCache   = redis.NewRCache[*CachedDirTreeNode]("vfm:dir:parent", redis.RCacheConfig{Exp: 1 * time.Hour})
	dirNameCache     = redis.NewRCache[string]("vfm:dir:name", redis.RCacheConfig{Exp: 1 * time.Hour})
	userDirTreeCache = redis.NewRCache[*DirTopDownTreeNode]("vfm:dir:user:tree", redis.RCacheConfig{Exp: 12 * time.Hour})
)

func init() {
	_imageSuffix.AddAll([]string{"jpeg", "jpg", "gif", "png", "svg", "bmp", "webp", "apng", "avif"})
	_videoSuffix.AddAll([]string{"mp4", "mov", "webm", "ogg"})
}

type FileVFolder struct {
	FolderNo   string
	Uuid       string
	CreateTime util.ETime
	CreateBy   string
	UpdateTime util.ETime
	UpdateBy   string
	IsDel      bool
}

type VFolderBrief struct {
	FolderNo string `json:"folderNo"`
	Name     string `json:"name"`
}

type ListedDir struct {
	Id   int    `json:"id"`
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

type ListedFile struct {
	Id             int        `json:"id"`
	Uuid           string     `json:"uuid"`
	Name           string     `json:"name"`
	UploadTime     util.ETime `json:"uploadTime"`
	UploaderName   string     `json:"uploaderName"`
	SizeInBytes    int64      `json:"sizeInBytes"`
	FileType       string     `json:"fileType"`
	UpdateTime     util.ETime `json:"updateTime"`
	ParentFileName string     `json:"parentFileName"`
	SensitiveMode  string     `json:"sensitiveMode"`
	ThumbnailToken string     `json:"thumbnailToken"`
	Thumbnail      string     `json:"-"`
	ParentFile     string     `json:"-"`
}

type GrantAccessReq struct {
	FileId    int    `json:"fileId" validation:"positive"`
	GrantedTo string `json:"grantedTo" validation:"notEmpty"`
}

type ListedVFolder struct {
	Id         int        `json:"id"`
	FolderNo   string     `json:"folderNo"`
	Name       string     `json:"name"`
	CreateTime util.ETime `json:"createTime"`
	CreateBy   string     `json:"createBy"`
	UpdateTime util.ETime `json:"updateTime"`
	UpdateBy   string     `json:"updateBy"`
	Ownership  string     `json:"ownership"`
}

type ListVFolderRes struct {
	Page    miso.Paging     `json:"paging"`
	Payload []ListedVFolder `json:"payload"`
}

type ShareVfolderReq struct {
	FolderNo string `json:"folderNo"`
	Username string `json:"username"`
}

type ParentFileInfo struct {
	Zero     bool   `json:"-"`
	FileKey  string `json:"fileKey"`
	Filename string `json:"fileName"`
}

type FileDownloadInfo struct {
	FileId         int
	Name           string
	IsLogicDeleted int
	FileType       string
	FstoreFileId   string
	UploaderNo     string
}

func (f *FileDownloadInfo) Deleted() bool {
	return f.IsLogicDeleted == DelY
}

func (f *FileDownloadInfo) IsFile() bool {
	return f.FileType == FileTypeFile
}

type FileInfo struct {
	Id               int
	Name             string
	Uuid             string
	FstoreFileId     string
	Thumbnail        string // thumbnail is also a fstore's file_id
	IsLogicDeleted   int
	IsPhysicDeleted  int
	SizeInBytes      int64
	UploaderNo       string // uploader's user_no
	UploaderName     string
	UploadTime       util.ETime
	LogicDeleteTime  util.ETime
	PhysicDeleteTime util.ETime
	UserGroup        int
	FileType         string
	ParentFile       string
	CreateTime       util.ETime
	CreateBy         string
	UpdateTime       util.ETime
	UpdateBy         string
	IsDel            int
	Hidden           bool
}

func (f FileInfo) IsZero() bool {
	return f.Id < 1
}

type VFolderWithOwnership struct {
	Id         int
	FolderNo   string
	Name       string
	CreateTime util.ETime
	CreateBy   string
	UpdateTime util.ETime
	UpdateBy   string
	Ownership  string
}

func (f *VFolderWithOwnership) IsOwner() bool {
	return f.Ownership == VfolderOwner
}

type VFolder struct {
	Id         int
	FolderNo   string
	Name       string
	CreateTime util.ETime
	CreateBy   string
	UpdateTime util.ETime
	UpdateBy   string
}

type UserVFolder struct {
	Id         int
	UserNo     string
	Username   string
	FolderNo   string
	Ownership  string
	GrantedBy  string // grantedBy (user_no)
	CreateTime util.ETime
	CreateBy   string
	UpdateTime util.ETime
	UpdateBy   string
}

func listFilesInVFolder(rail miso.Rail, db *gorm.DB, page miso.Paging, folderNo string, user common.User) (miso.PageRes[ListedFile], error) {
	return dbquery.NewPagedQuery[ListedFile](db).
		WithSelectQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.Select(`fi.id, fi.name, fi.parent_file, fi.uuid, fi.size_in_bytes,
			fi.uploader_name, fi.upload_time, fi.file_type, fi.update_time, fi.thumbnail`).
				Order("fi.id DESC")
		}).
		WithBaseQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.Table("file_info fi").
				Joins("LEFT JOIN file_vfolder fv ON (fi.uuid = fv.uuid AND fv.is_del = 0)").
				Joins("LEFT JOIN user_vfolder uv ON (fv.folder_no = uv.folder_no AND uv.is_del = 0)").
				Where("uv.user_no = ? AND uv.folder_no = ?", user.UserNo, folderNo).
				Where("fi.hidden = 0")
		}).Scan(rail, page)
}

type FileKeyName struct {
	Name string
	Uuid string
}

func queryFilenames(tx *gorm.DB, fileKeys []string) (map[string]string, error) {
	var rec []FileKeyName
	e := tx.Select("uuid, name").
		Table("file_info").
		Where("uuid IN ?", fileKeys).
		Scan(&rec).Error
	if e != nil {
		return nil, e
	}
	return util.StrMap[FileKeyName](rec,
			func(r FileKeyName) string { return r.Uuid },
			func(r FileKeyName) string { return r.Name }),
		nil
}

type ListFileReq struct {
	Page        miso.Paging `json:"paging"`
	Filename    *string     `json:"filename"`
	FolderNo    *string     `json:"folderNo"`
	FileType    *string     `json:"fileType"`
	ParentFile  *string     `json:"parentFile"`
	Sensitive   *bool       `json:"sensitive"`
	FileKey     *string
	OrderByName bool
}

func (q ListFileReq) IsEmpty() bool {
	return (q.ParentFile == nil || *q.ParentFile == "") && (q.Filename == nil || *q.Filename == "") && (q.FileKey == nil || *q.FileKey == "")
}

func ListFiles(rail miso.Rail, tx *gorm.DB, req ListFileReq, user common.User) (miso.PageRes[ListedFile], error) {
	var res miso.PageRes[ListedFile]
	var e error

	if req.FolderNo != nil && *req.FolderNo != "" {
		res, e = listFilesInVFolder(rail, tx, req.Page, *req.FolderNo, user)
	} else {
		res, e = listFilesSelective(rail, tx, req, user)
	}
	if e != nil {
		return res, e
	}

	parentFileKeys := util.NewSet[string]()
	for _, f := range res.Payload {
		if f.ParentFile != "" {
			parentFileKeys.Add(f.ParentFile)
		}
	}

	if !parentFileKeys.IsEmpty() {
		keyName, e := queryFilenames(tx, parentFileKeys.CopyKeys())
		if e != nil {
			return res, e
		}
		for i, f := range res.Payload {
			if name, ok := keyName[f.ParentFile]; ok {
				res.Payload[i].ParentFileName = name
			}
		}
	}

	// generate fstore tokens for thumbnail
	thumbnailTokenReq := make([]FstoreTmpTokenReq, 0, len(res.Payload))
	for _, p := range res.Payload {
		if p.Thumbnail != "" {
			thumbnailTokenReq = append(thumbnailTokenReq, FstoreTmpTokenReq{FileId: p.Thumbnail})
		}
	}
	m := BatchGetFstoreTmpToken(rail, thumbnailTokenReq)
	for i, f := range res.Payload {
		if f.Thumbnail != "" {
			if tkn, ok := m[f.Thumbnail]; ok {
				res.Payload[i].ThumbnailToken = tkn
			}
		}
	}

	return res, e
}

func listFilesSelective(rail miso.Rail, tx *gorm.DB, req ListFileReq, user common.User) (miso.PageRes[ListedFile], error) {

	//  If parentFile is empty, and filename are not queried, then we only return the top level file or dir.
	if req.IsEmpty() {
		req.ParentFile = new(string) // top-level file/dir
	}

	return dbquery.NewPagedQuery[ListedFile](tx).
		WithSelectQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.Select(`fi.id, fi.name, fi.parent_file, fi.uuid, fi.size_in_bytes,
			fi.uploader_name, fi.upload_time, fi.file_type, fi.update_time, fi.sensitive_mode, fi.thumbnail`)
		}).
		WithBaseQuery(func(q *dbquery.Query) *dbquery.Query {
			q = q.From("file_info fi").
				Eq("fi.uploader_no", user.UserNo).
				Eq("fi.is_logic_deleted", DelN).
				Eq("fi.is_del", 0).
				Eq("fi.hidden", 0)

			if req.OrderByName {
				q = q.Order("fi.name asc")
			}

			if req.FileKey != nil {
				q = q.Eq("fi.uuid", *req.FileKey)
			}
			if req.ParentFile != nil {
				q = q.Eq("fi.parent_file", *req.ParentFile)
			}
			if req.FileType != nil && *req.FileType != "" {
				q = q.Eq("fi.file_type", *req.FileType)
			}
			if req.Sensitive != nil && *req.Sensitive {
				q = q.Eq("fi.sensitive_mode", "N")
			}

			if req.Filename != nil && *req.Filename != "" {
				q = q.Where("match(fi.name) against (? IN NATURAL LANGUAGE MODE)", req.Filename)
				if !req.OrderByName {
					q = q.Order("fi.id desc")
				}
			} else {
				if !req.OrderByName {
					q = q.Order("fi.file_type asc, fi.id desc")
				}
			}

			return q
		}).
		Scan(rail, req.Page)
}

type PreflightCheckReq struct {
	Filename      string `form:"fileName"`
	ParentFileKey string `form:"parentFileKey"`
}

func FileExists(c miso.Rail, tx *gorm.DB, req PreflightCheckReq, userNo string) (bool, error) {
	var id int
	t := tx.Table("file_info").
		Select("id").
		Where("parent_file = ?", req.ParentFileKey).
		Where("name = ?", req.Filename).
		Where("uploader_no = ?", userNo).
		Where("file_type = ?", FileTypeFile).
		Where("is_logic_deleted = ?", DelN).
		Where("is_del = ?", false).
		Limit(1).
		Scan(&id)

	if t.Error != nil {
		return false, fmt.Errorf("failed to match file, %v", t.Error)
	}

	return id > 0, nil
}

func findFile(rail miso.Rail, tx *gorm.DB, fileKey string) (*FileInfo, error) {
	var f FileInfo
	t := tx.Raw("SELECT * FROM file_info WHERE uuid = ? AND is_del = 0", fileKey).
		Scan(&f)
	if t.Error != nil {
		return nil, t.Error
	}
	if t.RowsAffected < 1 {
		return nil, nil
	}
	return &f, t.Error
}

func findFileById(rail miso.Rail, tx *gorm.DB, id int) (FileInfo, error) {
	var f FileInfo

	t := tx.Raw("SELECT * FROM file_info WHERE id = ? AND is_del = 0", id).
		Scan(&f)
	if t.Error != nil {
		return f, t.Error
	}
	return f, nil
}

type FetchParentFileReq struct {
	FileKey string `form:"fileKey"`
}

func FindParentFile(c miso.Rail, tx *gorm.DB, req FetchParentFileReq, user common.User) (ParentFileInfo, error) {
	f, e := findFile(c, tx, req.FileKey)
	if e != nil {
		return ParentFileInfo{}, e
	}
	if f == nil {
		return ParentFileInfo{}, miso.NewErrf("File not found")
	}

	// dir is only visible to the uploader for now
	if f.UploaderNo != user.UserNo {
		return ParentFileInfo{}, miso.NewErrf("Not permitted")
	}

	if f.ParentFile == "" {
		return ParentFileInfo{Zero: true}, nil
	}

	pf, e := findFile(c, tx, f.ParentFile)
	if e != nil {
		return ParentFileInfo{}, e
	}
	if pf == nil {
		return ParentFileInfo{}, miso.NewErrf("File not found", fmt.Sprintf("ParentFile %v not found", f.ParentFile))
	}

	return ParentFileInfo{FileKey: pf.Uuid, Filename: pf.Name, Zero: false}, nil
}

type MakeDirReq struct {
	ParentFile string `json:"parentFile"`                 // Key of parent file
	Name       string `json:"name" validation:"notEmpty"` // name of the directory
}

func MakeDir(rail miso.Rail, tx *gorm.DB, req MakeDirReq, user common.User) (string, error) {
	rail.Infof("Making dir, req: %+v", req)

	var dir FileInfo
	dir.Name = req.Name
	dir.Uuid = util.GenIdP("ZZZ")
	dir.SizeInBytes = 0
	dir.FileType = FileTypeDir

	if e := _saveFile(rail, tx, dir, user); e != nil {
		return "", e
	}

	if req.ParentFile != "" {
		if e := MoveFileToDir(rail, tx, MoveIntoDirReq{Uuid: dir.Uuid, ParentFileUuid: req.ParentFile}, user); e != nil {
			return dir.Uuid, e
		}
	}

	return dir.Uuid, nil
}

type MoveIntoDirReq struct {
	Uuid           string `json:"uuid" validation:"notEmpty"`
	ParentFileUuid string `json:"parentFileUuid"`
}

func MoveFileToDir(rail miso.Rail, db *gorm.DB, req MoveIntoDirReq, user common.User) error {
	if req.Uuid == "" || req.Uuid == req.ParentFileUuid {
		return nil
	}

	// lock the file
	flock := fileLock(rail, req.Uuid)
	if err := flock.Lock(); err != nil {
		return miso.WrapErr(err)
	}
	defer flock.Unlock()

	fi, err := findFile(rail, db, req.Uuid)
	if err != nil {
		return miso.NewErrf("File not found").WithInternalMsg("failed to find file, uuid: %v, %v", req.Uuid, err)
	}
	if fi == nil {
		return miso.NewErrf("File not found")
	}
	if fi.ParentFile == req.ParentFileUuid {
		return nil
	}

	// prevent cycles between dir
	if fi.FileType == FileTypeDir {
		tree, err := doFetchDirTreeBottomUp(rail, db, &DirBottomUpTreeNode{
			FileKey: req.ParentFileUuid,
		})
		if err != nil {
			return err
		}
		for tree != nil {
			if tree.FileKey == req.Uuid {
				return miso.NewErrf("Found cycle between directories, invalid operation")
			}
			tree = tree.Child
		}
	}

	err = db.Transaction(func(tx *gorm.DB) error {

		// lock directory if necessary, if parentFileUuid is empty, the file is moved out of a directory
		if req.ParentFileUuid != "" {
			pflock := fileLock(rail, req.ParentFileUuid)
			if err := pflock.Lock(); err != nil {
				return err
			}
			defer pflock.Unlock()

			pf, e := findFile(rail, tx, req.ParentFileUuid)
			if e != nil {
				return fmt.Errorf("failed to find parentFile, %v", e)
			}
			if pf == nil {
				return fmt.Errorf("perentFile not found, parentFileKey: %v", req.ParentFileUuid)
			}
			rail.Debugf("parentFile: %+v", pf)

			if pf.UploaderNo != user.UserNo {
				return miso.NewErrf("You are not the owner of this directory")
			}

			if pf.FileType != FileTypeDir {
				return miso.NewErrf("Target file is not a directory")
			}

			if pf.IsLogicDeleted != DelN {
				return miso.NewErrf("Target file deleted")
			}

			newSize := pf.SizeInBytes + fi.SizeInBytes
			_, err := dbquery.NewQueryRail(rail, tx).
				Exec("UPDATE file_info SET size_in_bytes = ?, update_by = ?, update_time = ? WHERE uuid = ?",
					newSize, user.Username, time.Now(), req.ParentFileUuid)
			if err != nil {
				return fmt.Errorf("failed to updated dir's size, dir: %v, %v", req.ParentFileUuid, err)
			}
			rail.Infof("updated dir %v size to %v", req.ParentFileUuid, newSize)

		}

		_, err := dbquery.NewQueryRail(rail, tx).
			Exec("UPDATE file_info SET parent_file = ?, update_by = ?, update_time = ? WHERE uuid = ?",
				req.ParentFileUuid, user.Username, time.Now(), req.Uuid)
		if err != nil {
			return miso.WrapErr(err)
		}
		return nil
	})

	return err
}

func _saveFile(rail miso.Rail, tx *gorm.DB, f FileInfo, user common.User) error {
	uname := user.Username
	now := util.Now()

	f.IsLogicDeleted = DelN
	f.IsPhysicDeleted = DelN
	f.UploaderName = uname
	f.CreateBy = uname
	f.UploadTime = now
	f.CreateTime = now
	f.UploaderNo = user.UserNo

	_, err := dbquery.NewQueryRail(rail, tx).
		Table("file_info").
		Omit("id", "update_time", "update_by").
		Create(&f)
	if err == nil {
		rail.Infof("Saved file %+v", f)
		return nil
	}
	return miso.WrapErr(err)
}

func fileLock(rail miso.Rail, fileKey string) *redis.RLock {
	return redis.NewCustomRLock(rail, "file:uuid:"+fileKey, redis.RLockConfig{BackoffDuration: time.Second * 5})
}

type CreateVFolderReq struct {
	Name string `json:"name"`
}

func CreateVFolder(rail miso.Rail, tx *gorm.DB, r CreateVFolderReq, user common.User) (string, error) {
	userNo := user.UserNo

	v, e := redis.RLockRun(rail, "vfolder:user:"+userNo, func() (any, error) {

		var id int
		_, err := dbquery.NewQueryRail(rail, tx).
			Table("vfolder vf").
			Select("vf.id").
			Joins("LEFT JOIN user_vfolder uv ON (vf.folder_no = uv.folder_no)").
			Where("uv.user_no = ? AND uv.ownership = 'OWNER'", userNo).
			Where("vf.name = ?", r.Name).
			Where("vf.is_del = 0 AND uv.is_del = 0").
			Limit(1).
			Scan(&id)
		if err != nil {
			return "", err
		}
		if id > 0 {
			return "", miso.NewErrf(fmt.Sprintf("Found folder with same name ('%s')", r.Name))
		}

		folderNo := util.GenIdP("VFLD")

		e := tx.Transaction(func(tx *gorm.DB) error {

			ctime := util.Now()

			// for the vfolder
			vf := VFolder{Name: r.Name, FolderNo: folderNo, CreateTime: ctime, CreateBy: user.Username}
			if _, e := dbquery.NewQueryRail(rail, tx).
				Omit("id", "update_by", "update_time").Table("vfolder").Create(&vf); e != nil {
				return fmt.Errorf("failed to save VFolder, %v", e)
			}

			// for the user - vfolder relation
			uv := UserVFolder{
				FolderNo:   folderNo,
				UserNo:     userNo,
				Username:   user.Username,
				Ownership:  VfolderOwner,
				GrantedBy:  userNo,
				CreateTime: ctime,
				CreateBy:   user.Username}
			if _, e := dbquery.NewQueryRail(rail, tx).
				Omit("id", "update_by", "update_time").Table("user_vfolder").Create(&uv); e != nil {
				return fmt.Errorf("failed to save UserVFolder, %v", e)
			}
			return nil
		})
		if e != nil {
			return "", e
		}

		return folderNo, nil
	})
	if e != nil {
		return "", e
	}
	return v.(string), e
}

func ListDirs(r miso.Rail, tx *gorm.DB, user common.User) ([]ListedDir, error) {
	var dirs []ListedDir
	_, e := dbquery.NewQueryRail(r, tx).
		Table("file_info").
		Select("id, uuid, name").
		Where("uploader_no = ?", user.UserNo).
		Where("file_type = 'DIR'").
		Where("is_logic_deleted = 0").
		Where("is_del = 0").
		Scan(&dirs)
	return dirs, e
}

func findVFolder(rail miso.Rail, tx *gorm.DB, folderNo string, userNo string) (VFolderWithOwnership, error) {
	var vfo VFolderWithOwnership
	n, err := dbquery.NewQueryRail(rail, tx).
		Table("vfolder vf").
		Select("vf.*, uv.ownership").
		Joins("LEFT JOIN user_vfolder uv ON (vf.folder_no = uv.folder_no AND uv.is_del = 0)").
		Where("vf.is_del = 0").
		Where("uv.user_no = ?", userNo).
		Where("uv.folder_no = ?", folderNo).
		Limit(1).
		Scan(&vfo)
	if err != nil {
		return vfo, fmt.Errorf("failed to fetch vfolder info for current user, userNo: %v, folderNo: %v, %v", userNo, folderNo, err)
	}
	if n < 1 {
		return vfo, fmt.Errorf("vfolder not found, userNo: %v, folderNo: %v", userNo, folderNo)
	}
	return vfo, nil
}

func _lockFolderExec(c miso.Rail, folderNo string, r redis.Runnable) error {
	return redis.RLockExec(c, "vfolder:"+folderNo, r)
}

func ShareVFolder(rail miso.Rail, tx *gorm.DB, sharedTo vault.UserInfo, folderNo string, user common.User) error {
	if user.UserNo == sharedTo.UserNo {
		return nil
	}
	return _lockFolderExec(rail, folderNo, func() error {
		vfo, e := findVFolder(rail, tx, folderNo, user.UserNo)
		if e != nil {
			return e
		}
		if !vfo.IsOwner() {
			return miso.NewErrf("Operation not permitted")
		}

		var id int
		_, e = dbquery.NewQueryRail(rail, tx).
			Table("user_vfolder").
			Select("id").
			Where("folder_no = ?", folderNo).
			Where("user_no = ?", sharedTo.UserNo).
			Where("is_del = 0").
			Limit(1).
			Scan(&id)
		if e != nil {
			return fmt.Errorf("error occurred while querying user_vfolder, %v", e)
		}
		if id > 0 {
			rail.Infof("VFolder is shared already, folderNo: %s, sharedTo: %s", folderNo, sharedTo.Username)
			return nil
		}

		uv := UserVFolder{
			FolderNo:   folderNo,
			UserNo:     sharedTo.UserNo,
			Username:   sharedTo.Username,
			Ownership:  VfolderGranted,
			GrantedBy:  user.Username,
			CreateTime: util.Now(),
			CreateBy:   user.Username,
		}
		if _, e := dbquery.NewQueryRail(rail, tx).Omit("id", "update_by", "update_time").Table("user_vfolder").Create(&uv); e != nil {
			return fmt.Errorf("failed to save UserVFolder, %v", e)
		}
		rail.Infof("VFolder %s shared to %s by %s", folderNo, sharedTo.Username, user.Username)
		return nil
	})
}

type RemoveGrantedFolderAccessReq struct {
	FolderNo string `json:"folderNo"`
	UserNo   string `json:"userNo"`
}

func RemoveVFolderAccess(rail miso.Rail, tx *gorm.DB, req RemoveGrantedFolderAccessReq, user common.User) error {
	if user.UserNo == req.UserNo {
		return nil
	}
	return _lockFolderExec(rail, req.FolderNo, func() error {
		vfo, e := findVFolder(rail, tx, req.FolderNo, user.UserNo)
		if e != nil {
			return e
		}
		if !vfo.IsOwner() {
			return miso.NewErrf("Operation not permitted")
		}
		return tx.
			Exec("UPDATE user_vfolder SET is_del = 1, update_by = ? WHERE folder_no = ? AND user_no = ? AND ownership = 'GRANTED'",
				user.Username, req.FolderNo, req.UserNo).
			Error
	})
}

func ListVFolderBrief(rail miso.Rail, tx *gorm.DB, user common.User) ([]VFolderBrief, error) {
	var vfb []VFolderBrief
	_, e := dbquery.NewQueryRail(rail, tx).
		Select("f.folder_no, f.name").
		Table("vfolder f").
		Joins("LEFT JOIN user_vfolder uv ON (f.folder_no = uv.folder_no AND uv.is_del = 0)").
		Where("f.is_del = 0 AND uv.user_no = ? AND uv.ownership = 'OWNER'", user.UserNo).
		Scan(&vfb)
	return vfb, e
}

type AddFileToVfolderReq struct {
	FolderNo string   `json:"folderNo"`
	FileKeys []string `json:"fileKeys"`
	Sync     bool     `json:"-"`
}

func NewVFolderLock(rail miso.Rail, folderNo string) *redis.RLock {
	return redis.NewRLock(rail, "vfolder:"+folderNo)
}

func HandleAddFileToVFolderEvent(rail miso.Rail, tx *gorm.DB, evt AddFileToVfolderEvent) error {
	lock := NewVFolderLock(rail, evt.FolderNo)
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	var vfo VFolderWithOwnership
	var e error
	if vfo, e = findVFolder(rail, tx, evt.FolderNo, evt.UserNo); e != nil {
		return fmt.Errorf("failed to findVFolder, folderNo: %v, userNo: %v, %v", evt.FolderNo, evt.UserNo, e)
	}
	if !vfo.IsOwner() {
		return miso.NewErrf("Operation not permitted")
	}

	distinct := util.NewSet[string]()
	for _, fk := range evt.FileKeys {
		distinct.Add(fk)
	}

	filtered := util.Distinct(evt.FileKeys)
	if len(filtered) < 1 {
		return nil
	}

	now := util.Now()
	username := evt.Username
	doAddFileToVfolder := func(rail miso.Rail, folderNo string, fk string) error {
		var id int
		var err error
		_, err = dbquery.NewQueryRail(rail, tx).Select("id").
			Table("file_vfolder").
			Where("folder_no = ? AND uuid = ?", folderNo, fk).
			Where("is_del = 0").
			Scan(&id)
		if err != nil {
			return fmt.Errorf("failed to query file_vfolder record, %v", err)
		}
		if id > 0 {
			return nil
		}

		fvf := FileVFolder{FolderNo: folderNo, Uuid: fk, CreateTime: now, CreateBy: username}
		if _, err = dbquery.NewQueryRail(rail, tx).
			Table("file_vfolder").
			Omit("id", "update_by", "update_time").
			Create(&fvf); err != nil {
			return fmt.Errorf("failed to save file_vfolder record, %v", err)
		}
		rail.Infof("added file.uuid: %v to vfolder: %v by %v", fk, folderNo, username)
		return nil
	}

	// add files to vfolder
	dirs := []FileInfo{}
	for fk := range distinct.Keys {
		var e error

		f, e := findFile(rail, tx, fk)
		if e != nil {
			return e
		}

		if f == nil || f.UploaderNo != evt.UserNo {
			continue
		}
		if f.FileType != FileTypeFile {
			dirs = append(dirs, *f)
			continue
		}
		if e = doAddFileToVfolder(rail, evt.FolderNo, fk); e != nil {
			return fmt.Errorf("failed to doAddFileToVfolder, file.uuid: %v, %v", fk, e)
		}
	}

	// add files in dir to vfolder, but we only go one layer deep
	for _, dir := range dirs {
		var filesInDir []string
		var err error
		var page int = 1

		for {
			if filesInDir, err = ListFilesInDir(rail, tx, ListFilesInDirReq{
				FileKey: dir.Uuid,
				Limit:   500,
				Page:    page,
			}); err != nil {
				return fmt.Errorf("failed to list files in dir, dir.uuid: %v, %v", dir.Uuid, err)
			}

			if len(filesInDir) < 1 {
				break
			}

			for _, fk := range filesInDir {
				if !distinct.Add(fk) {
					continue
				}
				if err = doAddFileToVfolder(rail, evt.FolderNo, fk); err != nil {
					return fmt.Errorf("failed to doAddFileToVfolder, file.uuid: %v, %v", fk, e)
				}
			}
			page += 1
		}
	}
	return nil
}

func AddFileToVFolder(rail miso.Rail, tx *gorm.DB, req AddFileToVfolderReq, user common.User) error {

	if len(req.FileKeys) < 1 {
		return nil
	}

	vfo, e := findVFolder(rail, tx, req.FolderNo, user.UserNo)
	if e != nil {
		return e
	}
	if !vfo.IsOwner() {
		return miso.NewErrf("Operation not permitted")
	}

	evt := AddFileToVfolderEvent{
		Username: user.Username,
		UserNo:   user.UserNo,
		FolderNo: req.FolderNo,
		FileKeys: req.FileKeys,
	}

	err := AddFileToVFolderPipeline.Send(rail, evt)
	if err != nil {
		return fmt.Errorf("failed to publish AddFileToVfolderEvent, %+v, %v", evt, err)
	}
	return nil
}

type RemoveFileFromVfolderReq struct {
	FolderNo string   `json:"folderNo"`
	FileKeys []string `json:"fileKeys"`
}

func RemoveFileFromVFolder(rail miso.Rail, tx *gorm.DB, req RemoveFileFromVfolderReq, user common.User) error {
	if len(req.FileKeys) < 1 {
		return nil
	}

	return _lockFolderExec(rail, req.FolderNo, func() error {

		vfo, e := findVFolder(rail, tx, req.FolderNo, user.UserNo)
		if e != nil {
			return e
		}
		if !vfo.IsOwner() {
			return miso.NewErrf("Operation not permitted")
		}

		filtered := util.Distinct(req.FileKeys)
		if len(filtered) < 1 {
			return nil
		}

		for _, fk := range filtered {
			f, e := findFile(rail, tx, fk)
			if e != nil {
				return e
			}
			if f == nil {
				continue // file not found
			}

			if f.UploaderNo != user.UserNo {
				continue // not the uploader of the file
			}
			if f.FileType != FileTypeFile {
				continue // not a file type, may be a dir
			}

			_, e = dbquery.NewQueryRail(rail, tx).
				Exec("DELETE FROM file_vfolder WHERE folder_no = ? AND uuid = ?", req.FolderNo, fk)
			if e != nil {
				return fmt.Errorf("failed to delete file_vfolder record, %v", e)
			}
		}

		return nil
	})
}

func RemoveDeletedFileFromAllVFolder(rail miso.Rail, tx *gorm.DB, fileKey string) error {
	_, err := dbquery.NewQueryRail(rail, tx).
		Exec(`UPDATE file_vfolder SET is_del = 1 WHERE uuid = ?`, fileKey)
	if err != nil {
		return miso.WrapErrf(err, "failed to update file_vfolder, uuid: %v", fileKey)
	}
	rail.Infof("Removed file %v from all vfolders", fileKey)
	return nil
}

type ListVFolderReq struct {
	Page miso.Paging `json:"paging"`
	Name string      `json:"name"`
}

func ListVFolders(rail miso.Rail, tx *gorm.DB, req ListVFolderReq, user common.User) (ListVFolderRes, error) {
	t := newListVFoldersQuery(rail, tx, req, user.UserNo).
		Select("f.id, f.create_time, f.create_by, f.update_time, f.update_by, f.folder_no, f.name, uv.ownership").
		Order("f.id DESC").
		Offset(req.Page.GetOffset()).
		Limit(req.Page.GetLimit())

	var lvf []ListedVFolder
	if e := t.Scan(&lvf).Error; e != nil {
		return ListVFolderRes{}, fmt.Errorf("failed to query vfolder, req: %+v, %v", req, e)
	}

	var total int
	e := newListVFoldersQuery(rail, tx, req, user.UserNo).
		Select("COUNT(*)").
		Scan(&total).Error
	if e != nil {
		return ListVFolderRes{}, fmt.Errorf("failed to count vfolder, req: %+v, %v", req, e)
	}

	return ListVFolderRes{Page: miso.RespPage(req.Page, total), Payload: lvf}, nil
}

func newListVFoldersQuery(rail miso.Rail, tx *gorm.DB, req ListVFolderReq, userNo string) *gorm.DB {
	t := tx.Table("vfolder f").
		Joins("LEFT JOIN user_vfolder uv ON (f.folder_no = uv.folder_no AND uv.is_del = 0)").
		Where("f.is_del = 0 AND uv.user_no = ?", userNo)

	if req.Name != "" {
		t = t.Where("f.name like ?", "%"+req.Name+"%")
	}
	return t
}

type RemoveGrantedAccessReq struct {
	FileId int `json:"fileId" validation:"positive"`
	UserId int `json:"userId" validation:"positive"`
}

type ListGrantedFolderAccessReq struct {
	Page     miso.Paging `json:"paging"`
	FolderNo string      `json:"folderNo"`
}

type ListGrantedFolderAccessRes struct {
	Page    miso.Paging          `json:"paging"`
	Payload []ListedFolderAccess `json:"payload"`
}

type ListedFolderAccess struct {
	UserNo     string     `json:"userNo"`
	Username   string     `json:"username"`
	CreateTime util.ETime `json:"createTime"`
}

func ListGrantedFolderAccess(rail miso.Rail, tx *gorm.DB, req ListGrantedFolderAccessReq, user common.User) (ListGrantedFolderAccessRes, error) {
	folderNo := req.FolderNo
	vfo, e := findVFolder(rail, tx, folderNo, user.UserNo)
	if e != nil {
		return ListGrantedFolderAccessRes{}, e
	}
	if !vfo.IsOwner() {
		return ListGrantedFolderAccessRes{}, miso.NewErrf("Operation not permitted")
	}

	var l []ListedFolderAccess
	e = newListGrantedFolderAccessQuery(rail, tx, req).
		Select("user_no", "create_time", "username").
		Offset(req.Page.GetOffset()).
		Limit(req.Page.GetLimit()).
		Scan(&l).Error
	if e != nil {
		return ListGrantedFolderAccessRes{}, fmt.Errorf("failed to list granted folder access, req: %+v, %v", req, e)
	}

	var total int
	e = newListGrantedFolderAccessQuery(rail, tx, req).
		Select("count(*)").
		Scan(&total).Error
	if e != nil {
		return ListGrantedFolderAccessRes{}, fmt.Errorf("failed to count granted folder access, req: %+v, %v", req, e)
	}
	return ListGrantedFolderAccessRes{Payload: l, Page: miso.RespPage(req.Page, total)}, nil
}

func newListGrantedFolderAccessQuery(rail miso.Rail, tx *gorm.DB, r ListGrantedFolderAccessReq) *gorm.DB {
	return tx.Table("user_vfolder").
		Where("folder_no = ? AND ownership = 'GRANTED' AND is_del = 0", r.FolderNo)
}

type UpdateFileReq struct {
	Id            int `json:"id" validation:"positive"`
	Name          string
	SensitiveMode string
}

func UpdateFile(rail miso.Rail, tx *gorm.DB, r UpdateFileReq, user common.User) error {
	f, e := findFileById(rail, tx, r.Id)
	if e != nil {
		return e
	}
	if f.IsZero() {
		return miso.NewErrf("File not found")
	}

	// dir is only visible to the uploader for now
	if f.UploaderNo != user.UserNo {
		return miso.NewErrf("Not permitted")
	}

	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		return miso.NewErrf("Name can't be empty")
	}
	if r.SensitiveMode != "Y" && r.SensitiveMode != "N" {
		r.SensitiveMode = "N"
	}

	err := tx.Exec("UPDATE file_info SET name = ?, sensitive_mode = ?, update_by = ? WHERE id = ? AND is_logic_deleted = 0 AND is_del = 0",
		r.Name, r.SensitiveMode, user.Username, r.Id).Error

	return err
}

type CreateFileReq struct {
	Filename         string `json:"filename"`
	FakeFstoreFileId string `json:"fstoreFileId"`
	ParentFile       string `json:"parentFile"`
	Hidden           bool   `json:"-"`
}

func CreateFile(rail miso.Rail, tx *gorm.DB, r CreateFileReq, user common.User) (string, error) {
	fsf, e := fstore.FetchFileInfo(rail, fstore.FetchFileInfoReq{
		UploadFileId: r.FakeFstoreFileId,
	})
	if e != nil {
		if errors.Is(e, fstore.ErrFileNotFound) || errors.Is(e, fstore.ErrFileDeleted) {
			return "", miso.NewErrf("File not found or deleted")
		}
		return "", fmt.Errorf("failed to fetch file info from fstore, %v", e)
	}
	if fsf.Status != fstore.FileStatusNormal {
		return "", miso.NewErrf("File is deleted")
	}

	return SaveFileRecord(rail, tx, SaveFileReq{
		Filename:   r.Filename,
		Size:       fsf.Size,
		FileId:     fsf.FileId,
		Hidden:     r.Hidden,
		ParentFile: r.ParentFile,
	}, user)
}

type SysCreateFileReq struct {
	Filename         string `json:"filename"`
	FakeFstoreFileId string `json:"fstoreFileId"`
	ParentFile       string `json:"parentFile"`
	UserNo           string
}

type SaveFileReq struct {
	Filename   string
	FileId     string
	Size       int64
	ParentFile string
	Hidden     bool
}

func SaveFileRecord(rail miso.Rail, tx *gorm.DB, r SaveFileReq, user common.User) (string, error) {
	var f FileInfo
	f.Name = r.Filename
	f.Uuid = util.GenIdP("ZZZ")
	f.FstoreFileId = r.FileId
	f.SizeInBytes = r.Size
	f.FileType = FileTypeFile
	f.Hidden = r.Hidden

	if e := _saveFile(rail, tx, f, user); e != nil {
		return "", e
	}

	if r.ParentFile != "" {
		if e := MoveFileToDir(rail, tx, MoveIntoDirReq{Uuid: f.Uuid, ParentFileUuid: r.ParentFile}, user); e != nil {
			return "", e
		}
	}
	return f.Uuid, nil
}

func isVideo(name string) bool {
	i := strings.LastIndex(name, ".")
	if i < 0 || i == len([]rune(name))-1 {
		return false
	}

	suf := string(name[i+1:])
	return _videoSuffix.Has(strings.ToLower(suf))
}

func isImage(name string) bool {
	i := strings.LastIndex(name, ".")
	if i < 0 || i == len(name)-1 {
		return false
	}

	suf := string(name[i+1:])
	return _imageSuffix.Has(strings.ToLower(suf))
}

type DeleteFileReq struct {
	Uuid string `json:"uuid"`
}

func DeleteFile(rail miso.Rail, tx *gorm.DB, req DeleteFileReq, user common.User, condition func(FileInfo) bool) error {
	lock := fileLock(rail, req.Uuid)
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	f, e := findFile(rail, tx, req.Uuid)
	if e != nil {
		return fmt.Errorf("unable to find file, uuid: %v, %v", req.Uuid, e)
	}

	if f == nil {
		return miso.NewErrf("File not found")
	}

	if f.UploaderNo != user.UserNo {
		return miso.NewErrf("Not permitted")
	}

	if f.IsLogicDeleted == DelY {
		return nil // deleted already
	}

	if condition != nil && !condition(*f) {
		return nil // skip
	}

	if f.FileType == FileTypeDir { // if it's dir make sure it's empty
		var anyId int
		e := tx.Select("id").
			Table("file_info").
			Where("parent_file = ? AND is_logic_deleted = 0 AND is_del = 0", req.Uuid).
			Limit(1).
			Scan(&anyId).Error
		if e != nil {
			return fmt.Errorf("failed to count files in dir, uuid: %v, %v", req.Uuid, e)
		}
		if anyId > 0 {
			return miso.NewErrf("Directory is not empty, unable to delete it")
		}
	}

	if f.FstoreFileId != "" {
		if err := fstore.DeleteFile(rail, f.FstoreFileId); err != nil && !errors.Is(err, fstore.ErrFileDeleted) {
			return fmt.Errorf("failed to delete fstore file, fileId: %v, %v", f.FstoreFileId, err)
		}
	}

	if f.Thumbnail != "" {
		if err := fstore.DeleteFile(rail, f.Thumbnail); err != nil && !errors.Is(err, fstore.ErrFileDeleted) {
			return fmt.Errorf("failed to delete fstore file (thumbnail), fileId: %v, %v", f.Thumbnail, err)
		}
	}

	err := tx.Exec(`
		UPDATE file_info
		SET is_logic_deleted = 1, logic_delete_time = NOW()
		WHERE id = ? AND is_logic_deleted = 0`, f.Id).Error
	if err == nil {
		rail.Infof("Deleted file %v", f.Uuid)
	}
	return err
}

func validateFileAccess(rail miso.Rail, tx *gorm.DB, fileKey string, userNo string) (FileDownloadInfo, error) {
	var f FileDownloadInfo

	t := tx.
		Select("fi.id 'file_id', fi.fstore_file_id, fi.name, fi.is_logic_deleted, fi.file_type, fi.uploader_no").
		Table("file_info fi").
		Where("fi.uuid = ? AND fi.is_del = 0", fileKey).
		Limit(1).
		Scan(&f)
	if t.Error != nil {
		return f, t.Error
	}
	if t.RowsAffected < 1 {
		return f, miso.NewErrf("File not found")
	}
	if f.Deleted() {
		return f, miso.NewErrf("File deleted")
	}

	// is uploader of the file
	permitted := f.UploaderNo == userNo

	// user may have access to the vfolder, which contains the file
	if !permitted {
		var uvid int
		e := tx.
			Select("ifnull(uv.id, 0) as id").
			Table("file_info fi").
			Joins("LEFT JOIN file_vfolder fv ON (fi.uuid = fv.uuid AND fv.is_del = 0)").
			Joins("LEFT JOIN user_vfolder uv ON (uv.user_no = ? AND uv.folder_no = fv.folder_no AND uv.is_del = 0)", userNo).
			Where("fi.id = ?", f.FileId).
			Limit(1).
			Scan(&uvid).Error
		if e != nil {
			return f, fmt.Errorf("failed to query user folder relation for file, id: %v, %v", f.FileId, e)
		}
		permitted = uvid > 0 // granted access to a folder that contains this file
	}

	if !permitted {
		return f, miso.NewErrf("You are not permitted to access this file")
	}

	return f, nil
}

type GenerateTempTokenReq struct {
	FileKey string `json:"fileKey"`
}

func GenTempToken(rail miso.Rail, tx *gorm.DB, r GenerateTempTokenReq, user common.User) (string, error) {
	f, err := validateFileAccess(rail, tx, r.FileKey, user.UserNo)
	if err != nil {
		return "", fmt.Errorf("failed to validate file access, user: %+v, %w", user, err)
	}
	if !f.IsFile() {
		return "", miso.NewErrf("Downloading a directory is not supported")
	}

	if f.FstoreFileId == "" {
		rail.Errorf("File %v doesn't have mini-fstore file_id", r.FileKey)
		return "", miso.NewErrf("File cannot be downloaded, please contact system administrator")
	}

	return GetFstoreTmpToken(rail, f.FstoreFileId, f.Name)
}

type ListFilesInDirReq struct {
	FileKey string `form:"fileKey"`
	Limit   int    `form:"limit"`
	Page    int    `form:"page"`
}

func ListFilesInDir(rail miso.Rail, tx *gorm.DB, req ListFilesInDirReq) ([]string, error) {
	if req.Limit < 0 || req.Limit > 100 {
		req.Limit = 100
	}
	if req.Page < 1 {
		req.Page = 1
	}

	var fileKeys []string
	e := tx.Table("file_info").
		Select("uuid").
		Where("parent_file = ?", req.FileKey).
		Where("file_type = 'FILE'").
		Where("is_del = 0").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Scan(&fileKeys).Error
	return fileKeys, e
}

type FetchFileInfoReq struct {
	FileKey string `form:"fileKey"`
}

type FileInfoResp struct {
	Name         string `json:"name"`
	Uuid         string `json:"uuid"`
	SizeInBytes  int64  `json:"sizeInBytes"`
	UploaderNo   string `json:"uploaderNo"`
	UploaderName string `json:"uploaderName"`
	IsDeleted    bool   `json:"isDeleted"`
	FileType     string `json:"fileType"`
	ParentFile   string `json:"parentFile"`
	LocalPath    string `json:"localPath"`
	FstoreFileId string `json:"fstoreFileId"`
	Thumbnail    string `json:"thumbnail"`
}

func FetchFileInfoInternal(rail miso.Rail, tx *gorm.DB, req FetchFileInfoReq) (FileInfoResp, error) {
	var fir FileInfoResp
	f, e := findFile(rail, tx, req.FileKey)
	if e != nil {
		return fir, e
	}
	if f == nil {
		return fir, miso.NewErrf("File not found")
	}

	fir.Name = f.Name
	fir.Uuid = f.Uuid
	fir.SizeInBytes = f.SizeInBytes
	fir.UploaderNo = f.UploaderNo
	fir.UploaderName = f.UploaderName
	fir.IsDeleted = f.IsLogicDeleted == DelY
	fir.FileType = f.FileType
	fir.ParentFile = f.ParentFile
	fir.LocalPath = "" // files are managed by the mini-fstore, this field will no longer contain any value in it
	fir.FstoreFileId = f.FstoreFileId
	fir.Thumbnail = f.Thumbnail
	return fir, nil
}

type ValidateFileOwnerReq struct {
	FileKey string `form:"fileKey"`
	UserNo  string `form:"userNo"`
}

func ValidateFileOwner(rail miso.Rail, tx *gorm.DB, q ValidateFileOwnerReq) (bool, error) {
	var id int
	e := tx.Select("id").
		Table("file_info").
		Where("uuid = ?", q.FileKey).
		Where("uploader_no = ?", q.UserNo).
		Where("is_logic_deleted = 0").
		Limit(1).
		Scan(&id).Error
	return id > 0, e
}

type RemoveVFolderReq struct {
	FolderNo string
}

func RemoveVFolder(rail miso.Rail, tx *gorm.DB, user common.User, req RemoveVFolderReq) error {
	lock := NewVFolderLock(rail, req.FolderNo)
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	var vfo VFolderWithOwnership
	var e error
	if vfo, e = findVFolder(rail, tx, req.FolderNo, user.UserNo); e != nil {
		return fmt.Errorf("failed to findVFolder, folderNo: %v, userNo: %v, %v", req.FolderNo, user.UserNo, e)
	}
	if !vfo.IsOwner() {
		return miso.NewErrf("Operation not permitted")
	}
	if err := tx.Transaction(func(tx *gorm.DB) error {
		err := tx.Exec(`UPDATE vfolder SET is_del = 1, update_by = ? WHERE folder_no = ?`, user.Username, req.FolderNo).Error
		if err != nil {
			return fmt.Errorf("failed to update vfolder, folderNo: %v, %v", req.FolderNo, err)
		}
		err = tx.Exec(`UPDATE user_vfolder SET is_del = 1, update_by = ? WHERE folder_no = ?`, user.Username, req.FolderNo).Error
		if err != nil {
			return fmt.Errorf("failed to update user_vfolder, folderNo: %v, %v", req.FolderNo, err)
		}
		err = tx.Exec(`UPDATE file_vfolder SET is_del = 1, update_by = ? WHERE folder_no = ?`, user.Username, req.FolderNo).Error
		if err != nil {
			return fmt.Errorf("failed to update file_vfolder, folderNo: %v, %v", req.FolderNo, err)
		}
		return nil
	}); err != nil {
		return err
	}

	rail.Infof("VFolder %v deleted by %v", req.FolderNo, user.Username)
	return nil
}

type UnpackZipReq struct {
	FileKey       string // file key of the zip file
	ParentFileKey string // file key of current directory (not where the zip entries will be saved)
}

type UnpackZipExtra struct {
	FileKey       string // file key of the zip file
	ParentFileKey string // file key of the target directory
	UserNo        string
	Username      string
}

func UnpackZip(rail miso.Rail, db *gorm.DB, user common.User, req UnpackZipReq) error {
	flock := fileLock(rail, req.FileKey)
	if err := flock.Lock(); err != nil {
		return err
	}
	defer flock.Unlock()

	fi, err := findFile(rail, db, req.FileKey)
	if err != nil {
		return miso.NewErrf("File not found").WithInternalMsg("failed to find file, uuid: %v, %v", req.FileKey, err)
	}
	if fi == nil {
		return miso.NewErrf("File not found")
	}

	if fi.IsLogicDeleted == DelY {
		return miso.NewErrf("File is deleted")
	}

	if !strings.HasSuffix(strings.ToLower(fi.Name), ".zip") {
		return miso.NewErrf("File is not a zip")
	}

	dir, err := MakeDir(rail, db, MakeDirReq{
		Name:       fi.Name + " unpacked " + time.Now().Format("20060102_150405"),
		ParentFile: req.ParentFileKey,
	}, user)
	if err != nil {
		return fmt.Errorf("failed to make directory before unpacking zip, %w", err)
	}

	extra, err := json.WriteJson(UnpackZipExtra{
		FileKey:       req.FileKey,
		ParentFileKey: dir,
		UserNo:        user.UserNo,
		Username:      user.Username,
	})
	if err != nil {
		return fmt.Errorf("failed to write json as extra, %w", err)
	}

	err = fstore.TriggerFileUnzip(rail, fstore.UnzipFileReq{
		FileId:          fi.FstoreFileId,
		ReplyToEventBus: UnzipResultNotifyEventBus,
		Extra:           string(extra),
	})
	if err != nil {
		return fmt.Errorf("failed to TriggerFileUnZip, %w", err)
	}
	return nil
}

func HandleZipUnpackResult(rail miso.Rail, db *gorm.DB, evt fstore.UnzipFileReplyEvent) error {
	var extra UnpackZipExtra
	if err := json.ParseJson([]byte(evt.Extra), &extra); err != nil {
		return miso.UnknownErrf(err, "failed to unmarshal from extra")
	}

	if len(evt.ZipEntries) < 1 {
		return nil
	}

	for _, ze := range evt.ZipEntries {
		_, err := SaveFileRecord(rail, db, SaveFileReq{
			Filename:   ze.Name,
			FileId:     ze.FileId,
			Size:       ze.Size,
			ParentFile: extra.ParentFileKey,
		}, common.User{
			UserNo:   extra.UserNo,
			Username: extra.Username,
		})
		if err != nil {
			return miso.UnknownErrf(err, "failed to save zip entry, entry: %#v", ze)
		}
	}
	return nil
}

func TruncateDir(rail miso.Rail, db *gorm.DB, req DeleteFileReq, user common.User, async bool) error {
	rail.Infof("Truncating dir %v", req.Uuid)

	dir, e := findFile(rail, db, req.Uuid)
	if e != nil {
		return fmt.Errorf("unable to find file, uuid: %v, %v", req.Uuid, e)
	}

	if dir == nil {
		return miso.NewErrf("File not found")
	}

	if dir.UploaderNo != user.UserNo {
		return miso.NewErrf("Not permitted")
	}

	if dir.IsLogicDeleted == DelY {
		return nil // deleted already
	}

	if dir.FileType != FileTypeDir {
		return miso.NewErrf("Not a directory")
	}

	type ListedFilesInDir struct {
		Id       int
		Uuid     string
		FileType string
	}

	doTruncate := func() {
		rail := rail
		if async {
			rail = rail.NextSpan()
		}
		listFilesInDir := func(rail miso.Rail, minId int) ([]ListedFilesInDir, error) {
			var l []ListedFilesInDir
			_, err := dbquery.NewQueryRail(rail, db).Table("file_info").
				Select("id, uuid, file_type").
				Where("parent_file = ?", dir.Uuid).
				Where("id > ?", minId).
				Order("id asc").
				Limit(50).
				Scan(&l)

			rail.Debugf("listFilesInDir, minId: %v, dir.uuid: %v, count: %d", minId, dir.Uuid, len(l))
			return l, err
		}

		stillInDir := func(fi FileInfo) bool { return fi.ParentFile == dir.Uuid }

		minId := 0
		for {
			l, err := listFilesInDir(rail, minId)
			if err != nil {
				rail.Errorf("failed to listFilesInDir, minId: %v, dir.uuid: %v, %v", minId, dir.Uuid, err)
				return
			}
			if len(l) < 1 {
				rail.Infof("Truncated dir %v", req.Uuid)
				return
			}
			minId = l[len(l)-1].Id

			for _, lf := range l {
				if lf.FileType == FileTypeFile {
					if err := DeleteFile(rail, db, DeleteFileReq{Uuid: lf.Uuid}, user, stillInDir); err != nil {
						rail.Errorf("failed to DeleteFile in dir, dir.uuid: %v, deleting file.uuid: %v, %v", dir.Uuid, lf.Uuid, err)
						return
					}
					rail.Infof("Deleted file %v in dir %v", lf.Uuid, dir.Uuid)
				} else {
					if err := TruncateDir(rail, db, DeleteFileReq{Uuid: lf.Uuid}, user, false); err != nil {
						rail.Errorf("failed to TruncateDir in dir, in dir.uuid: %v, truncating dir.uuid: %v, %v", dir.Uuid, lf.Uuid, err)
						return
					}
				}
			}
		}
	}

	if async {
		vfmPool.Go(doTruncate)
	} else {
		doTruncate()
	}

	return nil
}

type CachedDirTreeNode struct {
	FileKey string
}

func FetchDirTreeBottomUp(rail miso.Rail, db *gorm.DB, req FetchDirTreeReq, user common.User) (*DirBottomUpTreeNode, error) {
	if util.IsBlankStr(req.FileKey) {
		return nil, nil
	}
	fi, err := findFile(rail, db, req.FileKey)
	if err != nil || fi == nil {
		return nil, err
	}
	if fi.IsLogicDeleted == DelY {
		return nil, nil
	}
	bottom := &DirBottomUpTreeNode{
		FileKey: req.FileKey,
		Name:    fi.Name,
	}
	return doFetchDirTreeBottomUp(rail, db, bottom)
}

func doFetchDirTreeBottomUp(rail miso.Rail, db *gorm.DB, child *DirBottomUpTreeNode) (*DirBottomUpTreeNode, error) {
	p, err := dirParentCache.Get(rail, child.FileKey, func() (*CachedDirTreeNode, error) {
		pi, err := doFindParentDir(rail, dbquery.NewQueryRail(rail, db), child.FileKey)
		if err != nil {
			return nil, err
		}
		if pi == nil {
			return nil, miso.NoneErr
		}
		return &CachedDirTreeNode{
			FileKey: pi.FileKey,
		}, nil
	})
	if err != nil {
		if miso.IsNoneErr(err) {
			return child, nil
		}
		return nil, err
	}

	name, err := cachedFindDirName(rail, dbquery.NewQueryRail(rail, db), p.FileKey)
	if err != nil {
		return nil, err
	}
	n := &DirBottomUpTreeNode{
		FileKey: p.FileKey,
		Name:    name,
		Child:   child,
	}
	return doFetchDirTreeBottomUp(rail, db, n)
}

type ParentDir struct {
	FileKey string
}

func doFindParentDir(c miso.Rail, q *dbquery.Query, fileKey string) (*ParentDir, error) {
	var pd ParentDir
	n, err := q.Raw(`SELECT parent_file file_key FROM file_info WHERE uuid = ? AND is_del = 0 AND is_logic_deleted = 0 LIMIT 1`, fileKey).
		Scan(&pd)
	if err != nil {
		return nil, fmt.Errorf("failed to find parent dir, %v", err)
	}
	if n < 1 {
		return nil, nil
	}

	// root directory
	if pd.FileKey == "" {
		return nil, nil
	}

	return &pd, nil
}

func cachedFindDirName(rail miso.Rail, q *dbquery.Query, fileKey string) (string, error) {
	return dirNameCache.Get(rail, fileKey, func() (string, error) {
		return findDirName(rail, q, fileKey)
	})
}

func findDirName(rail miso.Rail, q *dbquery.Query, fileKey string) (string, error) {
	var name string
	n, err := q.Raw(`SELECT name FROM file_info WHERE uuid = ? AND is_del = 0 AND is_logic_deleted = 0 LIMIT 1`, fileKey).
		Scan(&name)
	if err != nil {
		return "", fmt.Errorf("failed to find dir name, %v", err)
	}
	if n < 1 {
		return "", nil
	}
	return name, nil
}

func FetchDirTreeTopDown(rail miso.Rail, db *gorm.DB, user common.User) (*DirTopDownTreeNode, error) {
	return userDirTreeCache.Get(rail, user.UserNo, func() (*DirTopDownTreeNode, error) {
		root := &DirTopDownTreeNode{
			FileKey: "",
			Name:    "",
			Child:   []*DirTopDownTreeNode{},
		}
		seen := util.NewSet[string]()
		seen.Add(root.FileKey)
		return root, dfsDirTree(rail, db, root, user, seen)
	})
}

type TopDownTreeNodeBrief struct {
	Uuid string
	Name string
}

func dfsDirTree(rail miso.Rail, db *gorm.DB, root *DirTopDownTreeNode, user common.User, seen util.Set[string]) error {
	var cl []TopDownTreeNodeBrief
	n, err := dbquery.NewQueryRail(rail, db).
		From("file_info").
		Select("uuid, name").
		Eq("parent_file", root.FileKey).
		Eq("uploader_no", user.UserNo).
		Eq("file_type", "DIR").
		Eq("is_del", 0).
		Eq("is_logic_deleted", 0).
		Scan(&cl)
	if err != nil {
		return err
	}
	if n < 1 {
		return nil
	}
	for _, c := range cl {
		curr := &DirTopDownTreeNode{
			FileKey: c.Uuid,
			Name:    c.Name,
			Child:   []*DirTopDownTreeNode{},
		}
		root.Child = append(root.Child, curr)
		if seen.Add(curr.FileKey) { // just in case if the dir is somehow moved; can be solved by extra lock, not really necessary tho
			dfsDirTree(rail, db, curr, user, seen)
		}
	}
	return nil
}

type FileFstoreInfo struct {
	Name           string
	Uuid           string
	FstoreFileId   string
	Thumbnail      string
	IsLogicDeleted bool
}

func queryFileFstoreInfo(tx *gorm.DB, fileKeys []string) (map[string]FileFstoreInfo, error) {
	var rec []FileFstoreInfo
	e := tx.Select("uuid, name", "fstore_file_id", "thumbnail", "is_logic_deleted").
		Table("file_info").
		Where("uuid IN ?", fileKeys).
		Scan(&rec).Error
	if e != nil {
		return nil, e
	}
	return util.StrMap[FileFstoreInfo](rec,
			func(r FileFstoreInfo) string { return r.Uuid },
			func(r FileFstoreInfo) FileFstoreInfo { return r }),
		nil
}

func ValidateFileAccess(rail miso.Rail, db *gorm.DB, fileKey string, userNo string) error {
	if fileKey == "" {
		return nil // root dir
	}
	_, err := validateFileAccess(rail, db, fileKey, userNo)
	if err != nil {
		return miso.ErrNotPermitted.Wrapf(err, "failed to validate file access, userNo: %v, fileKey: %v", userNo, fileKey)
	}
	return nil
}

func InternalFetchFileInfo(rail miso.Rail, db *gorm.DB, req InternalFetchFileInfoReq) (InternalFetchFileInfoRes, error) {
	var res InternalFetchFileInfoRes
	n, err := dbquery.NewQueryRail(rail, db).
		Table("file_info").
		Eq("uuid", req.FileKey).
		Eq("is_logic_deleted", DelN).
		Select("name,upload_time,size_in_bytes,file_type").
		Scan(&res)
	if err != nil {
		return res, miso.UnknownErrf(err, "fileKey: %v", req.FileKey)
	}
	if n < 1 {
		return res, ErrFileNotFound.WithInternalMsg("fileKey: %v", req.FileKey)
	}
	return res, nil
}

type CheckDirExistsReq struct {
	ParentFile string `valid:"notEmpty"`
	Name       string `valid:"notEmpty"`
}

func CheckDirExists(rail miso.Rail, db *gorm.DB, req CheckDirExistsReq, user common.User) (string, error) {
	dirLock := fileLock(rail, req.ParentFile)
	if err := dirLock.Lock(); err != nil {
		return "", miso.UnknownErrf(err, "Unable to lock dir: %v", req.ParentFile)
	}
	defer dirLock.Unlock()

	if req.ParentFile != "" {
		fi, err := findFile(rail, db, req.ParentFile)
		if err != nil {
			return "", ErrFileNotFound.Wrapf(err, "failed to find file, parentFile: %v", req.ParentFile)
		}
		if fi == nil {
			return "", ErrFileNotFound
		}
	}

	var dirKey string
	_, err := dbquery.NewQueryRail(rail, db).
		Table("file_info").
		Eq("parent_file", req.ParentFile).
		Eq("name", req.Name).
		Where("is_logic_deleted != ?", DelY).
		Limit(1).
		Select("uuid").
		Scan(&dirKey)
	if err != nil {
		return "", miso.UnknownErrf(err, "failed to CheckDirExists, parentFile: %v, name: %v", req.ParentFile, req.Name)
	}
	return dirKey, nil
}
