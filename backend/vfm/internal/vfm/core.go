package vfm

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	fstore "github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/miso/errs"
	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/atom"
	"github.com/curtisnewbie/miso/util/hash"
	"github.com/curtisnewbie/miso/util/json"
	"github.com/curtisnewbie/miso/util/slutil"
	"github.com/curtisnewbie/miso/util/snowflake"
	"github.com/curtisnewbie/miso/util/strutil"
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
	_imageSuffix = hash.NewSet[string]()
	_videoSuffix = hash.NewSet[string]()

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
	CreateTime atom.Time
	CreateBy   string
	UpdateTime atom.Time
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
	Id             int       `json:"id"`
	Uuid           string    `json:"uuid"`
	Name           string    `json:"name"`
	UploadTime     atom.Time `json:"uploadTime"`
	UploaderName   string    `json:"uploaderName"`
	SizeInBytes    int64     `json:"sizeInBytes"`
	FileType       string    `json:"fileType"`
	UpdateTime     atom.Time `json:"updateTime"`
	ParentFileName string    `json:"parentFileName"`
	SensitiveMode  string    `json:"sensitiveMode"`
	IsComic        bool      `json:"isComic"`
	SeqKey         string    `json:"seqKey"`
	ThumbnailToken string    `json:"thumbnailToken"`
	Thumbnail      string    `json:"-"`
	ParentFile     string    `json:"-"`
}

type GrantAccessReq struct {
	FileId    int    `json:"fileId" validation:"positive"`
	GrantedTo string `json:"grantedTo" validation:"notEmpty"`
}

type ListedVFolder struct {
	Id         int       `json:"id"`
	FolderNo   string    `json:"folderNo"`
	Name       string    `json:"name"`
	CreateTime atom.Time `json:"createTime"`
	CreateBy   string    `json:"createBy"`
	UpdateTime atom.Time `json:"updateTime"`
	UpdateBy   string    `json:"updateBy"`
	Ownership  string    `json:"ownership"`
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
	UploadTime       atom.Time
	LogicDeleteTime  atom.Time
	PhysicDeleteTime atom.Time
	UserGroup        int
	FileType         string
	ParentFile       string
	SeqKey           string
	CreateTime       atom.Time
	CreateBy         string
	UpdateTime       atom.Time
	UpdateBy         string
	IsDel            int
	Hidden           bool
	IsComic          bool
}

func (f FileInfo) IsZero() bool {
	return f.Id < 1
}

type VFolderWithOwnership struct {
	Id         int
	FolderNo   string
	Name       string
	CreateTime atom.Time
	CreateBy   string
	UpdateTime atom.Time
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
	CreateTime atom.Time
	CreateBy   string
	UpdateTime atom.Time
	UpdateBy   string
}

type UserVFolder struct {
	Id         int
	UserNo     string
	Username   string
	FolderNo   string
	Ownership  string
	GrantedBy  string // grantedBy (user_no)
	CreateTime atom.Time
	CreateBy   string
	UpdateTime atom.Time
	UpdateBy   string
}

func listFilesInVFolder(rail miso.Rail, db *gorm.DB, page miso.Paging, folderNo string, user flow.User) (miso.PageRes[ListedFile], error) {
	return dbquery.NewPagedQuery[ListedFile](db).
		WithSelectQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.Select(`fi.id, fi.name, fi.parent_file, fi.uuid, fi.size_in_bytes,
			fi.uploader_name, fi.upload_time, fi.file_type, fi.update_time, fi.is_comic, fi.thumbnail`).
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

func queryFilenames(rail miso.Rail, db *gorm.DB, fileKeys []string) (map[string]string, error) {
	var rec []FileKeyName
	_, err := dbquery.NewQuery(rail, db).
		Select("uuid, name").
		Table("file_info").
		Where("uuid IN ?", fileKeys).
		Scan(&rec)
	if err != nil {
		return nil, err
	}
	return hash.StrMap[FileKeyName](rec,
			func(r FileKeyName) string { return r.Uuid },
			func(r FileKeyName) string { return r.Name }),
		nil
}

const (
	SortByTime   = "time"
	SortByName   = "name"
	SortByCustom = "custom"
)

type ListFileReq struct {
	Page        miso.Paging `json:"paging"`
	Filename    *string     `json:"filename"`
	FolderNo    *string     `json:"folderNo"`
	FileType    *string     `json:"fileType"`
	ParentFile  *string     `json:"parentFile"`
	Sensitive   *bool       `json:"sensitive"`
	FileKey     *string     `json:"fileKey"`
	OrderByName bool        `json:"orderByName" desc:"deprecated, use OrderBy instead"`
	OrderBy     string      `json:"orderBy"` // SortByTime (default), SortByName, SortByCustom
}

type FilePositionReq struct {
	FileKey     string `json:"fileKey"`
	ParentFile  string `json:"parentFile"`
	Limit       int    `json:"limit"`
	OrderByName bool   `json:"orderByName"` // deprecated, use OrderBy instead
	OrderBy     string `json:"orderBy"`     // "time" (default), "name", "custom"
	FileType    string `json:"fileType"`
}

type FilePositionRes struct {
	Page int `json:"page"`
}

func (q ListFileReq) IsEmpty() bool {
	return (q.ParentFile == nil || *q.ParentFile == "") && (q.Filename == nil || *q.Filename == "") && (q.FileKey == nil || *q.FileKey == "")
}

func ListFiles(rail miso.Rail, db *gorm.DB, req ListFileReq, user flow.User) (miso.PageRes[ListedFile], error) {
	var res miso.PageRes[ListedFile]
	var e error

	// backward compatibility: map OrderByName to OrderBy
	if req.OrderBy == "" && req.OrderByName {
		req.OrderBy = SortByName
	}

	// bootstrap custom order on first use
	if req.OrderBy == SortByCustom && req.ParentFile != nil && *req.ParentFile != "" {
		if err := bootstrapCustomOrder(rail, db, *req.ParentFile, user); err != nil {
			rail.Warnf("Failed to bootstrap custom order for dir %v: %v", *req.ParentFile, err)
		}
	}

	if req.FolderNo != nil && *req.FolderNo != "" {
		res, e = listFilesInVFolder(rail, db, req.Page, *req.FolderNo, user)
	} else {
		// force order by name for comic directories in default/time mode
		if req.OrderBy == SortByTime || req.OrderBy == "" {
			if req.ParentFile != nil && *req.ParentFile != "" {
				parent, ok, err := findFile(rail, db, *req.ParentFile)
				if err != nil {
					return res, err
				}
				if ok && parent.FileType == FileTypeDir && parent.IsComic {
					req.OrderBy = SortByName
				}
			}
		}
		res, e = listFilesSelective(rail, db, req, user)
	}
	if e != nil {
		return res, e
	}

	parentFileKeys := hash.NewSet[string]()
	for _, f := range res.Payload {
		if f.ParentFile != "" {
			parentFileKeys.Add(f.ParentFile)
		}
	}

	if !parentFileKeys.IsEmpty() {
		keyName, e := queryFilenames(rail, db, parentFileKeys.CopyKeys())
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

func listFilesSelective(rail miso.Rail, db *gorm.DB, req ListFileReq, user flow.User) (miso.PageRes[ListedFile], error) {

	//  If parentFile is empty, and filename are not queried, then we only return the top level file or dir.
	if req.IsEmpty() {
		req.ParentFile = new(string) // top-level file/dir
	}

	return dbquery.NewPagedQuery[ListedFile](db).
		WithSelectQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.Select(`fi.id, fi.name, fi.parent_file, fi.uuid, fi.size_in_bytes,
			fi.uploader_name, fi.upload_time, fi.file_type, fi.update_time, fi.sensitive_mode, fi.is_comic, fi.seq_key, fi.thumbnail`)
		}).
		WithBaseQuery(func(q *dbquery.Query) *dbquery.Query {
			q = q.From("file_info fi").
				Eq("fi.uploader_no", user.UserNo).
				Eq("fi.is_logic_deleted", DelN).
				Eq("fi.is_del", 0).
				Eq("fi.hidden", 0)

			// apply order based on OrderBy
			switch req.OrderBy {
			case SortByCustom:
				q = q.Order("fi.seq_key COLLATE utf8mb4_bin asc, fi.file_type asc, fi.id desc")
			case SortByName:
				q = q.Order("fi.name asc")
			default: // SortByTime or empty
				// ordering will be set below based on filename presence
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
				if req.OrderBy == "" || req.OrderBy == SortByTime {
					q = q.Order("fi.id desc")
				}
			} else {
				if req.OrderBy == "" || req.OrderBy == SortByTime {
					q = q.Order("fi.file_type asc, fi.id desc")
				}
			}

			return q
		}).
		Scan(rail, req.Page)
}

// CalcFilePosition returns which page a file appears on in the parent directory
// under the current sort and filter conditions (excluding fileKey filter).
func CalcFilePosition(rail miso.Rail, db *gorm.DB, req FilePositionReq, user flow.User) (int, error) {
	// 1. Look up target file to get its sort attributes
	f, ok, err := findFile(rail, db, req.FileKey)
	if err != nil || !ok || f.UploaderNo != user.UserNo {
		return 1, nil // fallback to first page
	}

	// backward compatibility: map OrderByName to OrderBy
	orderBy := req.OrderBy
	if orderBy == "" && req.OrderByName {
		orderBy = SortByName
	}

	// force name ordering for comic directories in default/time mode
	if (orderBy == SortByTime || orderBy == "") && req.ParentFile != "" {
		parent, ok, err := findFile(rail, db, req.ParentFile)
		if err == nil && ok && parent.FileType == FileTypeDir && parent.IsComic {
			orderBy = SortByName
		}
	}

	// 2. Build the same query as listFilesSelective but:
	//    - remove fileKey filter
	//    - add ordering-boundary condition to count rows before the target
	q := dbquery.NewQuery(rail, db).
		From("file_info fi").
		Eq("fi.uploader_no", user.UserNo).
		Eq("fi.is_logic_deleted", DelN).
		Eq("fi.is_del", 0).
		Eq("fi.hidden", 0)

	if req.ParentFile != "" {
		q = q.Eq("fi.parent_file", req.ParentFile)
	}
	if req.FileType != "" {
		q = q.Eq("fi.file_type", req.FileType)
	}

	// 3. Ordering-boundary condition: count rows that come BEFORE target file
	switch orderBy {
	case SortByCustom:
		// sort: fi.seq_key asc → rows before have seq_key < target seq_key
		q = q.Where("fi.seq_key COLLATE utf8mb4_bin < ?", f.SeqKey)
	case SortByName:
		// sort: fi.name asc → rows before have name < targetName
		q = q.Where("fi.name < ?", f.Name)
	default:
		// default: fi.file_type asc, fi.id desc
		q = q.Where("(fi.file_type < ? OR (fi.file_type = ? AND fi.id > ?))", f.FileType, f.FileType, f.Id)
	}

	count, err := q.Count()
	if err != nil {
		return 1, nil // fallback to first page
	}

	// 4. Calculate page: ceil((count+1) / limit)
	if req.Limit <= 0 {
		req.Limit = 10
	}
	page := (int(count) + req.Limit) / req.Limit // int ceiling: (count + 1 + limit - 1) / limit
	if page < 1 {
		page = 1
	}
	return page, nil
}

type PreflightCheckReq struct {
	Filename      string `form:"fileName"`
	ParentFileKey string `form:"parentFileKey"`
}

func FileExists(rail miso.Rail, db *gorm.DB, req PreflightCheckReq, userNo string) (bool, error) {
	ok, err := dbquery.NewQuery(rail, db).
		Table("file_info").
		Where("parent_file = ?", req.ParentFileKey).
		Where("name = ?", req.Filename).
		Where("uploader_no = ?", userNo).
		Where("file_type = ?", FileTypeFile).
		Where("is_logic_deleted = ?", DelN).
		Where("is_del = ?", false).
		Limit(1).
		HasAny()
	return ok, err
}

func findFile(rail miso.Rail, db *gorm.DB, fileKey string) (FileInfo, bool, error) {
	var f FileInfo
	ok, err := dbquery.NewQuery(rail, db).
		Table("file_info").
		Eq("uuid", fileKey).
		Eq("is_del", 0).
		ScanAny(&f)
	return f, ok, err
}

func findFileById(rail miso.Rail, tx *gorm.DB, id int) (FileInfo, bool, error) {
	var f FileInfo
	ok, err := dbquery.NewQuery(rail, tx).
		Raw("SELECT * FROM file_info WHERE id = ? AND is_del = 0", id).
		ScanAny(&f)
	if err != nil {
		return f, false, err
	}
	return f, ok, nil
}

type FetchParentFileReq struct {
	FileKey string `form:"fileKey"`
}

func FindParentFile(c miso.Rail, db *gorm.DB, req FetchParentFileReq, user flow.User) (ParentFileInfo, bool, error) {
	f, ok, err := findFile(c, db, req.FileKey)
	if err != nil {
		return ParentFileInfo{}, false, err
	}
	if !ok {
		return ParentFileInfo{}, false, ErrFileNotFound.New()
	}

	// dir is only visible to the uploader for now
	if f.UploaderNo != user.UserNo {
		return ParentFileInfo{}, false, miso.ErrNotPermitted.New()
	}

	if f.ParentFile == "" {
		return ParentFileInfo{}, false, nil
	}

	pf, ok, err := findFile(c, db, f.ParentFile)
	if err != nil {
		return ParentFileInfo{}, false, err
	}
	if !ok {
		return ParentFileInfo{}, false, ErrFileNotFound.WithInternalMsg("ParentFile %v not found", f.ParentFile)
	}

	return ParentFileInfo{FileKey: pf.Uuid, Filename: pf.Name}, true, nil
}

type MakeDirReq struct {
	ParentFile string `json:"parentFile"`                 // Key of parent file
	Name       string `json:"name" validation:"notEmpty"` // name of the directory
	Comic      bool   `json:"comic"`                      // mark directory as comic
}

func MakeDir(rail miso.Rail, tx *gorm.DB, req MakeDirReq, user flow.User) (string, error) {
	rail.Infof("Making dir, req: %+v", req)

	var dir FileInfo
	dir.Name = req.Name
	dir.Uuid = snowflake.IdPrefix("ZZZ")
	dir.SizeInBytes = 0
	dir.FileType = FileTypeDir
	dir.IsComic = req.Comic

	if e := _saveFile(rail, tx, dir, user); e != nil {
		return "", e
	}

	if req.ParentFile != "" {
		if e := MoveFileToDir(rail, tx, MoveIntoDirReq{Uuid: dir.Uuid, ParentFileUuid: req.ParentFile}, user); e != nil {
			return dir.Uuid, e
		}
		// If parent dir has custom ordering, assign seq_key at the beginning
		assignSeqKeyForNewFile(rail, tx, dir.Uuid, req.ParentFile, user)
	}

	return dir.Uuid, nil
}

type MoveIntoDirReq struct {
	Uuid           string `json:"uuid" validation:"notEmpty"`
	ParentFileUuid string `json:"parentFileUuid"`
}

func MoveFileToDir(rail miso.Rail, db *gorm.DB, req MoveIntoDirReq, user flow.User) error {
	if req.Uuid == "" || req.Uuid == req.ParentFileUuid {
		return nil
	}

	// lock the file
	flock := fileLock(rail, req.Uuid)
	if err := flock.Lock(); err != nil {
		return errs.Wrap(err)
	}
	defer flock.Unlock()

	fi, ok, err := findFile(rail, db, req.Uuid)
	if err != nil {
		return ErrFileNotFound.Wrapf(err, "failed to find file, uuid: %v", req.Uuid)
	}
	if !ok {
		return ErrFileNotFound.New()
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
				return errs.NewErrf("Found cycle between directories, invalid operation")
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

			pf, ok, e := findFile(rail, tx, req.ParentFileUuid)
			if e != nil {
				return errs.Wrapf(e, "failed to find parentFile")
			}
			if !ok {
				return errs.NewErrf("perentFile not found, parentFileKey: %v", req.ParentFileUuid)
			}
			rail.Debugf("parentFile: %+v", pf)

			if pf.UploaderNo != user.UserNo {
				return errs.NewErrf("You are not the owner of this directory")
			}

			if pf.FileType != FileTypeDir {
				return errs.NewErrf("Target file is not a directory")
			}

			if pf.IsLogicDeleted != DelN {
				return errs.NewErrf("Target file deleted")
			}

			newSize := pf.SizeInBytes + fi.SizeInBytes
			_, err := dbquery.NewQuery(rail, tx).
				Exec("UPDATE file_info SET size_in_bytes = ?, update_by = ?, update_time = ? WHERE uuid = ?",
					newSize, user.Username, time.Now(), req.ParentFileUuid)
			if err != nil {
				return errs.Wrapf(err, "failed to updated dir's size, dir: %v", req.ParentFileUuid)
			}
			rail.Infof("updated dir %v size to %v", req.ParentFileUuid, newSize)

		}

		_, err := dbquery.NewQuery(rail, tx).
			Exec("UPDATE file_info SET parent_file = ?, update_by = ?, update_time = ? WHERE uuid = ?",
				req.ParentFileUuid, user.Username, time.Now(), req.Uuid)
		return err
	})

	return err
}

func _saveFile(rail miso.Rail, db *gorm.DB, f FileInfo, user flow.User) error {
	uname := user.Username
	now := atom.Now()

	f.IsLogicDeleted = DelN
	f.IsPhysicDeleted = DelN
	f.UploaderName = uname
	f.CreateBy = uname
	f.UploadTime = now
	f.CreateTime = now
	f.UploaderNo = user.UserNo

	err := dbquery.NewQuery(rail, db).
		Table("file_info").
		Omit("id", "update_time", "update_by").
		CreateAny(&f)
	if err == nil {
		rail.Infof("Saved file %+v", f)
		return nil
	}
	return errs.Wrap(err)
}

func fileLock(rail miso.Rail, fileKey string) *redis.RLock {
	return redis.NewCustomRLock(rail, "file:uuid:"+fileKey, redis.RLockConfig{BackoffDuration: time.Second * 5})
}

type CreateVFolderReq struct {
	Name string `json:"name"`
}

func CreateVFolder(rail miso.Rail, db *gorm.DB, r CreateVFolderReq, user flow.User) (string, error) {
	userNo := user.UserNo

	return redis.RLockRun(rail, "vfolder:user:"+userNo, func() (string, error) {

		ok, err := dbquery.NewQuery(rail, db).
			Table("vfolder vf").
			Joins("LEFT JOIN user_vfolder uv ON (vf.folder_no = uv.folder_no)").
			Where("uv.user_no = ? AND uv.ownership = 'OWNER'", userNo).
			Where("vf.name = ?", r.Name).
			Where("vf.is_del = 0 AND uv.is_del = 0").
			Limit(1).
			HasAny()
		if err != nil {
			return "", err
		}
		if ok {
			return "", errs.NewErrf("Found folder with same name ('%s')", r.Name)
		}

		folderNo := snowflake.IdPrefix("VFLD")
		e := db.Transaction(func(tx *gorm.DB) error {
			ctime := atom.Now()

			// for the vfolder
			vf := VFolder{Name: r.Name, FolderNo: folderNo, CreateTime: ctime, CreateBy: user.Username}
			if _, e := dbquery.NewQuery(rail, tx).
				Omit("id", "update_by", "update_time").Table("vfolder").Create(&vf); e != nil {
				return errs.Wrapf(e, "failed to save VFolder")
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

			err := dbquery.NewQuery(rail, tx).
				Omit("id", "update_by", "update_time").
				Table("user_vfolder").
				CreateAny(&uv)
			if err != nil {
				return errs.Wrapf(err, "failed to save UserVFolder")
			}
			return nil
		})
		return folderNo, e
	})
}

func ListDirs(r miso.Rail, db *gorm.DB, user flow.User) ([]ListedDir, error) {
	var dirs []ListedDir
	_, e := dbquery.NewQuery(r, db).
		Table("file_info").
		Select("id, uuid, name").
		Where("uploader_no = ?", user.UserNo).
		Where("file_type = 'DIR'").
		Where("is_logic_deleted = 0").
		Where("is_del = 0").
		Scan(&dirs)
	return dirs, e
}

func findVFolder(rail miso.Rail, db *gorm.DB, folderNo string, userNo string) (VFolderWithOwnership, error) {
	var vfo VFolderWithOwnership
	ok, err := dbquery.NewQuery(rail, db).
		Table("vfolder vf").
		Select("vf.*, uv.ownership").
		Joins("LEFT JOIN user_vfolder uv ON (vf.folder_no = uv.folder_no AND uv.is_del = 0)").
		Where("vf.is_del = 0").
		Where("uv.user_no = ?", userNo).
		Where("uv.folder_no = ?", folderNo).
		Limit(1).
		ScanAny(&vfo)
	if err != nil {
		return vfo, errs.Wrapf(err, "failed to fetch vfolder info for current user, userNo: %v, folderNo: %v", userNo, folderNo)
	}
	if !ok {
		return vfo, errs.NewErrf("vfolder not found, userNo: %v, folderNo: %v", userNo, folderNo)
	}
	return vfo, nil
}

func _lockFolderExec(c miso.Rail, folderNo string, r redis.Runnable) error {
	return redis.RLockExec(c, "vfolder:"+folderNo, r)
}

func ShareVFolder(rail miso.Rail, db *gorm.DB, sharedTo vault.UserInfo, folderNo string, user flow.User) error {
	if user.UserNo == sharedTo.UserNo {
		return nil
	}
	return _lockFolderExec(rail, folderNo, func() error {
		vfo, err := findVFolder(rail, db, folderNo, user.UserNo)
		if err != nil {
			return err
		}
		if !vfo.IsOwner() {
			return errs.NewErrf("Operation not permitted")
		}

		ok, err := dbquery.NewQuery(rail, db).
			Table("user_vfolder").
			Where("folder_no = ?", folderNo).
			Where("user_no = ?", sharedTo.UserNo).
			Where("is_del = 0").
			Limit(1).
			HasAny()
		if err != nil {
			return fmt.Errorf("error occurred while querying user_vfolder, %v", err)
		}
		if ok {
			rail.Infof("VFolder is shared already, folderNo: %s, sharedTo: %s", folderNo, sharedTo.Username)
			return nil
		}

		uv := UserVFolder{
			FolderNo:   folderNo,
			UserNo:     sharedTo.UserNo,
			Username:   sharedTo.Username,
			Ownership:  VfolderGranted,
			GrantedBy:  user.Username,
			CreateTime: atom.Now(),
			CreateBy:   user.Username,
		}
		err = dbquery.NewQuery(rail, db).Omit("id", "update_by", "update_time").Table("user_vfolder").CreateAny(&uv)
		if err != nil {
			return errs.Wrapf(err, "failed to save UserVFolder")
		}
		rail.Infof("VFolder %s shared to %s by %s", folderNo, sharedTo.Username, user.Username)
		return nil
	})
}

type RemoveGrantedFolderAccessReq struct {
	FolderNo string `json:"folderNo"`
	UserNo   string `json:"userNo"`
}

func RemoveVFolderAccess(rail miso.Rail, db *gorm.DB, req RemoveGrantedFolderAccessReq, user flow.User) error {
	if user.UserNo == req.UserNo {
		return nil
	}
	return _lockFolderExec(rail, req.FolderNo, func() error {
		vfo, e := findVFolder(rail, db, req.FolderNo, user.UserNo)
		if e != nil {
			return e
		}
		if !vfo.IsOwner() {
			return errs.NewErrf("Operation not permitted")
		}
		_, err := dbquery.NewQuery(rail, db).
			Exec("UPDATE user_vfolder SET is_del = 1, update_by = ? WHERE folder_no = ? AND user_no = ? AND ownership = 'GRANTED'",
				user.Username, req.FolderNo, req.UserNo)
		return err
	})
}

func ListVFolderBrief(rail miso.Rail, tx *gorm.DB, user flow.User) ([]VFolderBrief, error) {
	var vfb []VFolderBrief
	err := dbquery.NewQuery(rail, tx).
		Select("f.folder_no, f.name").
		Table("vfolder f").
		Joins("LEFT JOIN user_vfolder uv ON (f.folder_no = uv.folder_no AND uv.is_del = 0)").
		Where("f.is_del = 0 AND uv.user_no = ? AND uv.ownership = 'OWNER'", user.UserNo).
		ScanVal(&vfb)
	return vfb, err
}

type AddFileToVfolderReq struct {
	FolderNo string   `json:"folderNo"`
	FileKeys []string `json:"fileKeys"`
	Sync     bool     `json:"-"`
}

func NewVFolderLock(rail miso.Rail, folderNo string) *redis.RLock {
	return redis.NewRLock(rail, "vfolder:"+folderNo)
}

func newReorderDirLock(rail miso.Rail, parentFile string) *redis.RLock {
	return redis.NewRLockf(rail, "vfm:reorder:dir:%v", parentFile)
}

func HandleAddFileToVFolderEvent(rail miso.Rail, db *gorm.DB, evt AddFileToVfolderEvent) error {
	lock := NewVFolderLock(rail, evt.FolderNo)
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	var vfo VFolderWithOwnership
	var e error
	if vfo, e = findVFolder(rail, db, evt.FolderNo, evt.UserNo); e != nil {
		return errs.Wrapf(e, "failed to findVFolder, folderNo: %v, userNo: %v", evt.FolderNo, evt.UserNo)
	}
	if !vfo.IsOwner() {
		return errs.ErrNotPermitted.New()
	}

	distinct := hash.NewSet[string](evt.FileKeys...)
	if distinct.IsEmpty() {
		return nil
	}

	now := atom.Now()
	username := evt.Username
	doAddFileToVfolder := func(rail miso.Rail, folderNo string, fk string) error {
		var err error
		ok, err := dbquery.NewQuery(rail, db).
			Table("file_vfolder").
			Where("folder_no = ? AND uuid = ?", folderNo, fk).
			Where("is_del = 0").
			HasAny()
		if err != nil {
			return errs.Wrapf(err, "failed to query file_vfolder record")
		}
		if ok {
			return nil
		}

		fvf := FileVFolder{FolderNo: folderNo, Uuid: fk, CreateTime: now, CreateBy: username}
		if _, err = dbquery.NewQuery(rail, db).
			Table("file_vfolder").
			Omit("id", "update_by", "update_time").
			Create(&fvf); err != nil {
			return errs.Wrapf(err, "failed to save file_vfolder record")
		}
		rail.Infof("added file.uuid: %v to vfolder: %v by %v", fk, folderNo, username)
		return nil
	}

	// add files to vfolder
	dirs := []FileInfo{}
	for fk := range distinct.Keys {
		var e error

		f, ok, e := findFile(rail, db, fk)
		if e != nil {
			return e
		}
		if !ok || f.UploaderNo != evt.UserNo {
			continue
		}
		if f.FileType != FileTypeFile {
			dirs = append(dirs, f)
			continue
		}
		if e = doAddFileToVfolder(rail, evt.FolderNo, fk); e != nil {
			return errs.Wrapf(e, "failed to doAddFileToVfolder, file.uuid: %v", fk)
		}
	}

	// add files in dir to vfolder, but we only go one layer deep
	for _, dir := range dirs {
		var filesInDir []string
		var err error
		var page int = 1

		for {
			if filesInDir, err = ListFilesInDir(rail, db, ListFilesInDirReq{
				FileKey: dir.Uuid,
				Limit:   500,
				Page:    page,
			}); err != nil {
				return errs.Wrapf(err, "failed to list files in dir, dir.uuid: %v", dir.Uuid)
			}

			if len(filesInDir) < 1 {
				break
			}

			for _, fk := range filesInDir {
				if !distinct.Add(fk) {
					continue
				}
				if err = doAddFileToVfolder(rail, evt.FolderNo, fk); err != nil {
					return errs.Wrapf(e, "failed to doAddFileToVfolder, file.uuid: %v", fk)
				}
			}
			page += 1
		}
	}
	return nil
}

func AddFileToVFolder(rail miso.Rail, db *gorm.DB, req AddFileToVfolderReq, user flow.User) error {

	if len(req.FileKeys) < 1 {
		return nil
	}

	vfo, e := findVFolder(rail, db, req.FolderNo, user.UserNo)
	if e != nil {
		return e
	}
	if !vfo.IsOwner() {
		return miso.ErrNotPermitted.New()
	}

	evt := AddFileToVfolderEvent{
		Username: user.Username,
		UserNo:   user.UserNo,
		FolderNo: req.FolderNo,
		FileKeys: req.FileKeys,
	}

	err := AddFileToVFolderPipeline.Send(rail, evt)
	if err != nil {
		return errs.Wrapf(err, "failed to publish AddFileToVfolderEvent, %+v", evt)
	}
	return nil
}

type RemoveFileFromVfolderReq struct {
	FolderNo string   `json:"folderNo"`
	FileKeys []string `json:"fileKeys"`
}

func RemoveFileFromVFolder(rail miso.Rail, tx *gorm.DB, req RemoveFileFromVfolderReq, user flow.User) error {
	if len(req.FileKeys) < 1 {
		return nil
	}

	return _lockFolderExec(rail, req.FolderNo, func() error {

		vfo, e := findVFolder(rail, tx, req.FolderNo, user.UserNo)
		if e != nil {
			return e
		}
		if !vfo.IsOwner() {
			return errs.ErrNotPermitted.New()
		}

		filtered := slutil.Distinct(req.FileKeys)
		if len(filtered) < 1 {
			return nil
		}

		for _, fk := range filtered {
			f, ok, err := findFile(rail, tx, fk)
			if err != nil {
				return err
			}
			if !ok {
				continue // file not found
			}

			if f.UploaderNo != user.UserNo {
				continue // not the uploader of the file
			}
			if f.FileType != FileTypeFile {
				continue // not a file type, may be a dir
			}

			if _, err = dbquery.NewQuery(rail, tx).
				Exec("DELETE FROM file_vfolder WHERE folder_no = ? AND uuid = ?", req.FolderNo, fk); err != nil {
				return errs.Wrapf(err, "failed to delete file_vfolder record")
			}
		}

		return nil
	})
}

func RemoveDeletedFileFromAllVFolder(rail miso.Rail, tx *gorm.DB, fileKey string) error {
	_, err := dbquery.NewQuery(rail, tx).
		Exec(`UPDATE file_vfolder SET is_del = 1 WHERE uuid = ?`, fileKey)
	if err != nil {
		return errs.Wrapf(err, "failed to update file_vfolder, uuid: %v", fileKey)
	}
	rail.Infof("Removed file %v from all vfolders", fileKey)
	return nil
}

type ListVFolderReq struct {
	Page miso.Paging `json:"paging"`
	Name string      `json:"name"`
}

func ListVFolders(rail miso.Rail, tx *gorm.DB, req ListVFolderReq, user flow.User) (ListVFolderRes, error) {
	var lvf []ListedVFolder
	err := newListVFoldersQuery(rail, tx, req, user.UserNo).
		Select("f.id, f.create_time, f.create_by, f.update_time, f.update_by, f.folder_no, f.name, uv.ownership").
		Order("f.id DESC").
		Offset(req.Page.GetOffset()).
		Limit(req.Page.GetLimit()).
		ScanVal(&lvf)
	if err != nil {
		return ListVFolderRes{}, errs.Wrapf(err, "failed to query vfolder, req: %+v", req)
	}

	var total int
	e := newListVFoldersQuery(rail, tx, req, user.UserNo).
		Select("COUNT(*)").
		ScanVal(&total)
	if e != nil {
		return ListVFolderRes{}, errs.Wrapf(e, "failed to count vfolder, req: %+v", req)
	}

	return ListVFolderRes{Page: miso.RespPage(req.Page, total), Payload: lvf}, nil
}

func newListVFoldersQuery(rail miso.Rail, db *gorm.DB, req ListVFolderReq, userNo string) *dbquery.Query {
	t := dbquery.NewQuery(rail, db).
		Table("vfolder f").
		Joins("LEFT JOIN user_vfolder uv ON (f.folder_no = uv.folder_no AND uv.is_del = 0)").
		Where("f.is_del = 0 AND uv.user_no = ?", userNo).
		LikeIf(req.Name != "", "f.name", req.Name)

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
	UserNo     string    `json:"userNo"`
	Username   string    `json:"username"`
	CreateTime atom.Time `json:"createTime"`
}

func ListGrantedFolderAccess(rail miso.Rail, tx *gorm.DB, req ListGrantedFolderAccessReq, user flow.User) (ListGrantedFolderAccessRes, error) {
	folderNo := req.FolderNo
	vfo, err := findVFolder(rail, tx, folderNo, user.UserNo)
	if err != nil {
		return ListGrantedFolderAccessRes{}, err
	}
	if !vfo.IsOwner() {
		return ListGrantedFolderAccessRes{}, errs.ErrNotPermitted.New()
	}

	var l []ListedFolderAccess
	err = newListGrantedFolderAccessQuery(rail, tx, req).
		SelectCols(ListedFolderAccess{}).
		Offset(req.Page.GetOffset()).
		Limit(req.Page.GetLimit()).
		ScanVal(&l)
	if err != nil {
		return ListGrantedFolderAccessRes{}, errs.Wrapf(err, "failed to list granted folder access, req: %+v", req)
	}

	var total int
	err = newListGrantedFolderAccessQuery(rail, tx, req).
		Select("COUNT(*)").
		ScanVal(&total)
	if err != nil {
		return ListGrantedFolderAccessRes{}, errs.Wrapf(err, "failed to count granted folder access, req: %+v", req)
	}
	return ListGrantedFolderAccessRes{Payload: l, Page: miso.RespPage(req.Page, total)}, nil
}

func newListGrantedFolderAccessQuery(rail miso.Rail, tx *gorm.DB, r ListGrantedFolderAccessReq) *dbquery.Query {
	return dbquery.NewQuery(rail, tx).
		Table("user_vfolder").
		Where("folder_no = ? AND ownership = 'GRANTED' AND is_del = 0", r.FolderNo)
}

type UpdateFileReq struct {
	Id            int    `json:"id" validation:"positive"`
	Name          string `json:"name"`
	SensitiveMode string `json:"sensitiveMode"`
	IsComic       *bool  `json:"isComic"`
}

func UpdateFile(rail miso.Rail, tx *gorm.DB, r UpdateFileReq, user flow.User) error {
	f, ok, e := findFileById(rail, tx, r.Id)
	if e != nil {
		return e
	}
	if !ok {
		return ErrFileNotFound.New()
	}

	// dir is only visible to the uploader for now
	if f.UploaderNo != user.UserNo {
		return errs.ErrNotPermitted.New()
	}

	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		return errs.NewErrf("Name can't be empty")
	}
	if r.SensitiveMode != "Y" && r.SensitiveMode != "N" {
		r.SensitiveMode = "N"
	}
	if r.IsComic != nil && f.FileType != FileTypeDir {
		return errs.NewErrf("is_comic can only be set on directories")
	}

	q := dbquery.NewQuery(rail, tx).
		Table("file_info").
		Set("name", r.Name).
		Set("sensitive_mode", r.SensitiveMode).
		Set("update_by", user.Username).
		Eq("id", r.Id).
		Eq("is_logic_deleted", DelN).
		Eq("is_del", 0)

	if r.IsComic != nil {
		q = q.Set("is_comic", *r.IsComic)
	}

	err := q.UpdateAny()
	return err
}

type CreateFileReq struct {
	Filename         string `json:"filename"`
	FakeFstoreFileId string `json:"fstoreFileId"`
	ParentFile       string `json:"parentFile"`
	Hidden           bool   `json:"-"`
}

func CreateFile(rail miso.Rail, tx *gorm.DB, r CreateFileReq, user flow.User) (string, error) {
	fsf, e := fstore.FetchFileInfo(rail, fstore.FetchFileInfoReq{
		UploadFileId: r.FakeFstoreFileId,
	})
	if e != nil {
		if errs.IsAny(e, fstore.ErrFileNotFound, fstore.ErrFileDeleted) {
			return "", errs.NewErrf("File not found or deleted")
		}
		return "", errs.Wrapf(e, "failed to fetch file info from fstore")
	}
	if fsf.Status != fstore.FileStatusNormal {
		return "", ErrFileDeleted.New()
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
	UserNo           string `json:"userNo"`
}

type SaveFileReq struct {
	Filename   string
	FileId     string
	Size       int64
	ParentFile string
	Hidden     bool
}

func SaveFileRecord(rail miso.Rail, tx *gorm.DB, r SaveFileReq, user flow.User) (string, error) {
	var f FileInfo
	f.Name = r.Filename
	f.Uuid = snowflake.IdPrefix("ZZZ")
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
		// If parent dir has custom ordering, assign seq_key at the beginning
		assignSeqKeyForNewFile(rail, tx, f.Uuid, r.ParentFile, user)
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

func DeleteFile(rail miso.Rail, tx *gorm.DB, req DeleteFileReq, user flow.User, condition func(FileInfo) bool) error {
	lock := fileLock(rail, req.Uuid)
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	f, ok, e := findFile(rail, tx, req.Uuid)
	if e != nil {
		return fmt.Errorf("unable to find file, uuid: %v, %v", req.Uuid, e)
	}

	if !ok {
		return ErrFileNotFound.New()
	}

	if f.UploaderNo != user.UserNo {
		return errs.ErrNotPermitted.New()
	}

	if f.IsLogicDeleted == DelY {
		return nil // deleted already
	}

	if condition != nil && !condition(f) {
		return nil // skip
	}

	if f.FileType == FileTypeDir { // if it's dir make sure it's empty
		ok, e := dbquery.NewQuery(rail, tx).
			Table("file_info").
			Where("parent_file = ? AND is_logic_deleted = 0 AND is_del = 0", req.Uuid).
			Limit(1).
			HasAny()
		if e != nil {
			return fmt.Errorf("failed to count files in dir, uuid: %v, %v", req.Uuid, e)
		}
		if ok {
			return errs.NewErrf("Directory is not empty, unable to delete it")
		}
	}

	if f.FstoreFileId != "" {
		if err := fstore.DeleteFile(rail, f.FstoreFileId); err != nil && !errors.Is(err, fstore.ErrFileDeleted) {
			return errs.Wrapf(err, "failed to delete fstore file, fileId: %v", f.FstoreFileId)
		}
	}

	if f.Thumbnail != "" {
		if err := fstore.DeleteFile(rail, f.Thumbnail); err != nil && !errors.Is(err, fstore.ErrFileDeleted) {
			return errs.Wrapf(err, "failed to delete fstore file (thumbnail), fileId: %v", f.Thumbnail)
		}
	}

	_, err := dbquery.NewQuery(rail, tx).Exec(`
		UPDATE file_info
		SET is_logic_deleted = 1, logic_delete_time = NOW()
		WHERE id = ? AND is_logic_deleted = 0`, f.Id)
	if err == nil {
		rail.Infof("Deleted file %v (File Type: %v)", f.Uuid, f.FileType)
	}
	return err
}

func validateFileAccess(rail miso.Rail, tx *gorm.DB, fileKey string, userNo string) (FileDownloadInfo, error) {
	var f FileDownloadInfo

	ok, err := dbquery.NewQuery(rail, tx).
		Select("fi.id 'file_id', fi.fstore_file_id, fi.name, fi.is_logic_deleted, fi.file_type, fi.uploader_no").
		Table("file_info fi").
		Where("fi.uuid = ? AND fi.is_del = 0", fileKey).
		Limit(1).
		ScanAny(&f)
	if err != nil {
		return f, err
	}
	if !ok {
		return f, ErrFileNotFound.New()
	}
	if f.Deleted() {
		return f, ErrFileDeleted.New()
	}

	// is uploader of the file
	permitted := f.UploaderNo == userNo

	// user may have access to the vfolder, which contains the file
	if !permitted {
		var uvid int
		ok, e := dbquery.NewQuery(rail, tx).
			Select("ifnull(uv.id, 0) as id").
			Table("file_info fi").
			Joins("LEFT JOIN file_vfolder fv ON (fi.uuid = fv.uuid AND fv.is_del = 0)").
			Joins("LEFT JOIN user_vfolder uv ON (uv.user_no = ? AND uv.folder_no = fv.folder_no AND uv.is_del = 0)", userNo).
			Where("fi.id = ?", f.FileId).
			Limit(1).
			ScanAny(&uvid)
		if e != nil {
			return f, errs.Wrapf(e, "failed to query user folder relation for file, id: %v", f.FileId)
		}
		permitted = ok // granted access to a folder that contains this file
	}

	if !permitted {
		return f, errs.NewErrf("You are not permitted to access this file")
	}

	return f, nil
}

type GenerateTempTokenReq struct {
	FileKey string `json:"fileKey"`
}

func GenTempToken(rail miso.Rail, tx *gorm.DB, r GenerateTempTokenReq, user flow.User) (string, error) {
	f, err := validateFileAccess(rail, tx, r.FileKey, user.UserNo)
	if err != nil {
		return "", errs.Wrapf(err, "failed to validate file access, user: %+v", user)
	}
	if !f.IsFile() {
		return "", errs.NewErrf("Downloading a directory is not supported")
	}

	if f.FstoreFileId == "" {
		rail.Errorf("File %v doesn't have mini-fstore file_id", r.FileKey)
		return "", errs.NewErrf("File cannot be downloaded, please contact system administrator")
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
	e := dbquery.NewQuery(rail, tx).
		Table("file_info").
		Select("uuid").
		Where("parent_file = ?", req.FileKey).
		Where("file_type = 'FILE'").
		Where("is_del = 0").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		ScanVal(&fileKeys)
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
	f, ok, e := findFile(rail, tx, req.FileKey)
	if e != nil {
		return fir, e
	}
	if !ok {
		return fir, ErrFileNotFound.New()
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
	ok, e := dbquery.NewQuery(rail, tx).
		Select("id").
		Table("file_info").
		Where("uuid = ?", q.FileKey).
		Where("uploader_no = ?", q.UserNo).
		Where("is_logic_deleted = 0").
		HasAny()
	return ok, e
}

type RemoveVFolderReq struct {
	FolderNo string `json:"folderNo"`
}

func RemoveVFolder(rail miso.Rail, db *gorm.DB, user flow.User, req RemoveVFolderReq) error {
	lock := NewVFolderLock(rail, req.FolderNo)
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	var vfo VFolderWithOwnership
	var e error
	if vfo, e = findVFolder(rail, db, req.FolderNo, user.UserNo); e != nil {
		return errs.Wrapf(e, "failed to findVFolder, folderNo: %v, userNo: %v", req.FolderNo, user.UserNo)
	}
	if !vfo.IsOwner() {
		return errs.ErrNotPermitted.New()
	}

	if err := dbquery.RunTransaction(rail, db, func(qry func() *dbquery.Query) error {
		_, err := qry().Exec(`UPDATE vfolder SET is_del = 1, update_by = ? WHERE folder_no = ?`, user.Username, req.FolderNo)
		if err != nil {
			return errs.Wrapf(err, "failed to update vfolder, folderNo: %v", req.FolderNo)
		}
		_, err = qry().Exec(`UPDATE user_vfolder SET is_del = 1, update_by = ? WHERE folder_no = ?`, user.Username, req.FolderNo)
		if err != nil {
			return errs.Wrapf(err, "failed to update user_vfolder, folderNo: %v", req.FolderNo)
		}
		_, err = qry().Exec(`UPDATE file_vfolder SET is_del = 1, update_by = ? WHERE folder_no = ?`, user.Username, req.FolderNo)
		if err != nil {
			return errs.Wrapf(err, "failed to update file_vfolder, folderNo: %v", req.FolderNo)
		}
		return nil
	}); err != nil {
		return err
	}

	rail.Infof("VFolder %v deleted by %v", req.FolderNo, user.Username)
	return nil
}

type UnpackZipReq struct {
	FileKey       string `json:"fileKey"`       // file key of the zip file
	ParentFileKey string `json:"parentFileKey"` // file key of current directory (not where the zip entries will be saved)
}

type UnpackZipExtra struct {
	FileKey       string // file key of the zip file
	ParentFileKey string // file key of the target directory
	UserNo        string
	Username      string
}

func UnpackZip(rail miso.Rail, db *gorm.DB, user flow.User, req UnpackZipReq) error {
	flock := fileLock(rail, req.FileKey)
	if err := flock.Lock(); err != nil {
		return err
	}
	defer flock.Unlock()

	fi, ok, err := findFile(rail, db, req.FileKey)
	if err != nil {
		return ErrFileNotFound.Wrapf(err, "failed to find file, uuid: %v", req.FileKey)
	}
	if !ok {
		return ErrFileNotFound.New()
	}

	if fi.IsLogicDeleted == DelY {
		return ErrFileDeleted.New()
	}

	if !strings.HasSuffix(strings.ToLower(fi.Name), ".zip") {
		return errs.NewErrf("File is not a zip")
	}

	dir, err := MakeDir(rail, db, MakeDirReq{
		Name:       fi.Name + " unpacked " + time.Now().Format("20060102_150405"),
		ParentFile: req.ParentFileKey,
	}, user)
	if err != nil {
		return errs.Wrapf(err, "failed to make directory before unpacking zip")
	}

	extra, err := json.WriteJson(UnpackZipExtra{
		FileKey:       req.FileKey,
		ParentFileKey: dir,
		UserNo:        user.UserNo,
		Username:      user.Username,
	})
	if err != nil {
		return errs.Wrapf(err, "failed to write json as extra")
	}

	err = fstore.TriggerFileUnzip(rail, fstore.UnzipFileReq{
		FileId:          fi.FstoreFileId,
		ReplyToEventBus: UnzipResultNotifyEventBus,
		Extra:           string(extra),
	})
	if err != nil {
		return errs.Wrapf(err, "failed to TriggerFileUnZip")
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
		}, flow.User{
			UserNo:   extra.UserNo,
			Username: extra.Username,
		})
		if err != nil {
			return miso.UnknownErrf(err, "failed to save zip entry, entry: %#v", ze)
		}
	}
	return nil
}

func TruncateDir(rail miso.Rail, db *gorm.DB, req DeleteFileReq, user flow.User, async bool) error {
	rail.Infof("Truncating dir %v", req.Uuid)

	dir, ok, e := findFile(rail, db, req.Uuid)
	if e != nil {
		return errs.Wrapf(e, "unable to find file, uuid: %v", req.Uuid)
	}

	if !ok {
		return ErrFileNotFound.New()
	}

	if dir.UploaderNo != user.UserNo {
		return errs.ErrNotPermitted.New()
	}

	if dir.IsLogicDeleted == DelY {
		return nil // deleted already
	}

	if dir.FileType != FileTypeDir {
		return errs.NewErrf("Not a directory")
	}

	type ListedFilesInDir struct {
		Id       int
		Uuid     string
		FileType string
	}

	doTruncate := func() {
		rail := rail
		if async {
			rail = rail.NewCtx().NextSpanId()
		}
		listFilesInDir := func(rail miso.Rail, minId int) ([]ListedFilesInDir, error) {
			var l []ListedFilesInDir
			err := dbquery.NewQuery(rail, db).
				Table("file_info").
				Select("id, uuid, file_type").
				Where("parent_file = ?", dir.Uuid).
				Where("id > ?", minId).
				Order("id asc").
				Limit(50).
				ScanVal(&l)

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

				if err := DeleteFile(rail, db, DeleteFileReq{Uuid: dir.Uuid}, user, nil); err != nil {
					rail.Errorf("Failed to delete empty dir after truncation, %v", dir.Uuid, err)
				}

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

func FetchDirTreeBottomUp(rail miso.Rail, db *gorm.DB, req FetchDirTreeReq, user flow.User) (*DirBottomUpTreeNode, error) {
	if strutil.IsBlankStr(req.FileKey) {
		return nil, nil
	}
	fi, ok, err := findFile(rail, db, req.FileKey)
	if err != nil || !ok {
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
	p, ok, err := dirParentCache.GetElse(rail, child.FileKey, func() (*CachedDirTreeNode, bool, error) {
		pi, err := doFindParentDir(rail, dbquery.NewQuery(rail, db), child.FileKey)
		if err != nil || pi == nil {
			return nil, false, err
		}
		return &CachedDirTreeNode{
			FileKey: pi.FileKey,
		}, true, nil
	})
	if err != nil {
		return nil, err
	}
	if !ok {
		return child, nil
	}

	name, err := cachedFindDirName(rail, dbquery.NewQuery(rail, db), p.FileKey)
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
	ok, err := q.Raw(`SELECT parent_file file_key FROM file_info WHERE uuid = ? AND is_del = 0 AND is_logic_deleted = 0 LIMIT 1`, fileKey).
		ScanAny(&pd)
	if err != nil {
		return nil, fmt.Errorf("failed to find parent dir, %v", err)
	}
	if !ok {
		return nil, nil
	}

	// root directory
	if pd.FileKey == "" {
		return nil, nil
	}

	return &pd, nil
}

func cachedFindDirName(rail miso.Rail, q *dbquery.Query, fileKey string) (string, error) {
	v, err := dirNameCache.GetValElse(rail, fileKey, func() (string, error) {
		return findDirName(rail, q, fileKey)
	})
	return v, err
}

func findDirName(rail miso.Rail, q *dbquery.Query, fileKey string) (string, error) {
	var name string
	err := q.Raw(`SELECT name FROM file_info WHERE uuid = ? AND is_del = 0 AND is_logic_deleted = 0 LIMIT 1`, fileKey).
		ScanVal(&name)
	if err != nil {
		return "", errs.Wrapf(err, "failed to find dir name")
	}
	return name, nil
}

func FetchDirTreeTopDown(rail miso.Rail, db *gorm.DB, user flow.User) (*DirTopDownTreeNode, error) {
	v, err := userDirTreeCache.GetValElse(rail, user.UserNo, func() (*DirTopDownTreeNode, error) {
		root := &DirTopDownTreeNode{
			FileKey: "",
			Name:    "",
			Child:   []*DirTopDownTreeNode{},
		}
		seen := hash.NewSet[string]()
		seen.Add(root.FileKey)
		return root, dfsDirTree(rail, db, root, user, seen)
	})
	return v, err
}

type TopDownTreeNodeBrief struct {
	Uuid string
	Name string
}

func dfsDirTree(rail miso.Rail, db *gorm.DB, root *DirTopDownTreeNode, user flow.User, seen hash.Set[string]) error {
	var cl []TopDownTreeNodeBrief
	n, err := dbquery.NewQuery(rail, db).
		Table("file_info").
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

func queryFileFstoreInfo(rail miso.Rail, tx *gorm.DB, fileKeys []string) (map[string]FileFstoreInfo, error) {
	var rec []FileFstoreInfo
	e := dbquery.NewQuery(rail, tx).
		Select("uuid, name", "fstore_file_id", "thumbnail", "is_logic_deleted").
		Table("file_info").
		Where("uuid IN ?", fileKeys).
		ScanVal(&rec)
	if e != nil {
		return nil, e
	}
	return hash.StrMap[FileFstoreInfo](rec,
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

func InternalFetchFileInfo(rail miso.Rail, db *gorm.DB, req ItnFetchFileInfoReq) (ItnFetchFileInfoRes, error) {
	var res ItnFetchFileInfoRes
	n, err := dbquery.NewQuery(rail, db).
		Table("file_info").
		Eq("uuid", req.FileKey).
		Eq("is_logic_deleted", DelN).
		Select("name,upload_time,size_in_bytes,file_type,uuid 'file_key'").
		Scan(&res)
	if err != nil {
		return res, miso.UnknownErrf(err, "fileKey: %v", req.FileKey)
	}
	if n < 1 {
		return res, ErrFileNotFound.WithInternalMsg("fileKey: %v", req.FileKey)
	}
	return res, nil
}

func InternalBatchFetchFileInfo(rail miso.Rail, db *gorm.DB, req InternalBatchFetchFileInfoReq) ([]ItnFetchFileInfoRes, error) {
	req.FileKey = slutil.FastDistinct(req.FileKey)
	if len(req.FileKey) < 1 {
		return []ItnFetchFileInfoRes{}, nil
	}
	var res []ItnFetchFileInfoRes
	err := dbquery.NewQuery(rail, db).
		Table("file_info").
		In("uuid", req.FileKey).
		Eq("is_logic_deleted", DelN).
		Select("name,upload_time,size_in_bytes,file_type,uuid 'file_key'").
		ScanVal(&res)
	return res, err
}

type CheckDirExistsReq struct {
	ParentFile string `valid:"notEmpty"`
	Name       string `valid:"notEmpty"`
}

func CheckDirExists(rail miso.Rail, db *gorm.DB, req CheckDirExistsReq, user flow.User) (string, error) {
	dirLock := fileLock(rail, req.ParentFile)
	if err := dirLock.Lock(); err != nil {
		return "", miso.UnknownErrf(err, "Unable to lock dir: %v", req.ParentFile)
	}
	defer dirLock.Unlock()

	if req.ParentFile != "" {
		_, ok, err := findFile(rail, db, req.ParentFile)
		if err != nil {
			return "", ErrFileNotFound.Wrapf(err, "failed to find file, parentFile: %v", req.ParentFile)
		}
		if !ok {
			return "", ErrFileNotFound.New()
		}
	}

	var dirKey string
	_, err := dbquery.NewQuery(rail, db).
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

func InternalUpdateFileInfo(rail miso.Rail, db *gorm.DB, req ApiItnUpdateFileInfoReq) error {
	flock := fileLock(rail, req.FileKey)
	if err := flock.Lock(); err != nil {
		return err
	}
	defer flock.Unlock()

	return dbquery.NewQuery(rail, db).
		Table("file_info").
		Set("name", req.Name).
		Eq("uuid", req.FileKey).
		Eq("is_del", 0).
		UpdateAny()
}

type FetchDirThumbnailReq struct {
	DirFileKey string `json:"dirFileKey" valid:"notEmpty"`
}

type FetchDirThumbnailRes struct {
	FstoreToken string `json:"fstoreToken,omitzero"`
}

type BatchFetchDirThumbnailReq struct {
	DirFileKeys []string `json:"dirFileKeys" valid:"notEmpty"`
}

type DirThumbnailWithKey struct {
	DirFileKey  string `json:"dirFileKey"`
	FstoreToken string `json:"fstoreToken,omitzero"`
}

func FetchDirThumbnail(rail miso.Rail, db *gorm.DB, req FetchDirThumbnailReq, user flow.User) (FetchDirThumbnailRes, error) {
	f, ok, err := findFile(rail, db, req.DirFileKey)
	if err != nil {
		return FetchDirThumbnailRes{}, err
	}
	if !ok {
		return FetchDirThumbnailRes{}, ErrFileNotFound.New()
	}
	if f.UploaderNo != user.UserNo {
		return FetchDirThumbnailRes{}, miso.ErrNotPermitted.New()
	}
	fstFileId, ok, err := findFirstThumbnailFileId(rail, db, req.DirFileKey)
	if err != nil {
		return FetchDirThumbnailRes{}, err
	}
	if !ok || fstFileId == "" {
		return FetchDirThumbnailRes{}, nil
	}
	tkn, err := GetFstoreTmpToken(rail, fstFileId, "")
	if err != nil {
		return FetchDirThumbnailRes{}, err
	}
	return FetchDirThumbnailRes{FstoreToken: tkn}, nil
}

func BatchFetchDirThumbnail(rail miso.Rail, db *gorm.DB, req BatchFetchDirThumbnailReq, user flow.User) ([]DirThumbnailWithKey, error) {
	if len(req.DirFileKeys) == 0 {
		return []DirThumbnailWithKey{}, nil
	}

	// Remove duplicates
	dirKeys := slutil.FastDistinct(req.DirFileKeys)

	// Single query using correlated subquery for MySQL 5.7 compatibility
	// This is much more efficient than scanning all files with thumbnails
	type DirWithThumbnail struct {
		DirFileKey   string
		FstoreFileId string
	}
	var dirs []DirWithThumbnail

	query := `
		SELECT
			d.uuid as dir_file_key,
			(SELECT f.thumbnail
			 FROM file_info f
			 WHERE f.parent_file = d.uuid
			   AND f.is_del = 0
			   AND f.thumbnail != ''
			 ORDER BY f.id DESC
			 LIMIT 1) as fstore_file_id
		FROM file_info d
		WHERE d.uuid IN ?
			AND d.is_del = 0
			AND d.uploader_no = ?
	`

	err := db.Raw(query, dirKeys, user.UserNo).Scan(&dirs).Error
	if err != nil {
		return nil, miso.UnknownErrf(err, "failed to batch fetch dir thumbnails")
	}

	// Build response with all requested keys
	res := make([]DirThumbnailWithKey, 0, len(dirKeys))
	keyMap := make(map[string]string, len(dirKeys))
	for _, d := range dirs {
		keyMap[d.DirFileKey] = d.FstoreFileId
	}

	// Collect fileIds that need tokens
	var tokenReqs []FstoreTmpTokenReq
	for _, key := range dirKeys {
		if fileId, exists := keyMap[key]; exists && fileId != "" {
			tokenReqs = append(tokenReqs, FstoreTmpTokenReq{FileId: fileId})
		}
	}

	// Batch get tokens
	tokenMap := BatchGetFstoreTmpToken(rail, tokenReqs)

	// Build final response
	for _, key := range dirKeys {
		item := DirThumbnailWithKey{DirFileKey: key}
		if fileId, exists := keyMap[key]; exists && fileId != "" {
			item.FstoreToken = tokenMap[fileId]
		}
		res = append(res, item)
	}

	return res, nil
}

func findFirstThumbnailFileId(rail miso.Rail, db *gorm.DB, dirFileKey string) (string, bool, error) {
	var fstFileId string
	ok, err := dbquery.NewQuery(rail, db).
		Table("file_info").
		Eq("parent_file", dirFileKey).
		Eq("is_del", 0).
		Ne("thumbnail", "").
		OrderDesc("id").
		Limit(1).
		Select("thumbnail").
		ScanAny(&fstFileId)
	if err != nil {
		return "", false, err
	}
	if !ok {
		return "", false, nil
	}
	return fstFileId, true, nil
}

// ---------------------------------------------------------------------------
// Custom Drag-and-Drop Ordering (fractional-indexing)
//
// Port of rocicorp/fractional-indexing v4.0.0 (CC0 license)
// https://github.com/rocicorp/fractional-indexing
// ---------------------------------------------------------------------------

const (
	base62Digits = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	base52Digits = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// getDigitIndex returns a [256]int lookup table mapping byte -> digit index.
func getDigitIndex(digits string) [256]int {
	var m [256]int
	for i := 0; i < len(digits); i++ {
		m[digits[i]] = i
	}
	return m
}

// midpoint computes the fractional part between two keys. a may be empty string,
// b may be nil (infinite upper bound) or non-nil. digits must be in ascending
// byte order. lookup is the precomputed digit index for digits.
func midpoint(a string, b *string, digits string, lookup [256]int) (string, error) {
	zero := digits[0]
	if len(a) > 0 && a[len(a)-1] == zero {
		return "", fmt.Errorf("trailing zero in a: %q", a)
	}
	if b != nil && len(*b) > 0 && (*b)[len(*b)-1] == zero {
		return "", fmt.Errorf("trailing zero in b: %q", *b)
	}
	if b != nil && a >= *b {
		return "", fmt.Errorf("a >= b: %q >= %q", a, *b)
	}

	if b != nil {
		n := 0
		for {
			var ac, bc byte
			if n < len(a) {
				ac = a[n]
			} else {
				ac = zero
			}
			if n < len(*b) {
				bc = (*b)[n]
			} else {
				break
			}
			if ac != bc {
				break
			}
			n++
		}
		if n > 0 {
			r, err := midpoint(a[n:], bStrPtr((*b)[n:]), digits, lookup)
			if err != nil {
				return "", err
			}
			return (*b)[:n] + r, nil
		}
	}

	digitA := 0
	if len(a) > 0 {
		digitA = lookup[a[0]]
	}
	digitB := len(digits)
	if b != nil {
		digitB = lookup[(*b)[0]]
	}

	if digitB-digitA > 1 {
		midDigit := int(math.Round(0.5 * float64(digitA+digitB)))
		return string(digits[midDigit]), nil
	}

	if b != nil && len(*b) > 1 {
		return (*b)[:1], nil
	}

	r, err := midpoint(a[1:], nil, digits, lookup)
	if err != nil {
		return "", err
	}
	return string(digits[digitA]) + r, nil
}

// bStrPtr converts a string to *string for use with nil (null bound).
func bStrPtr(s string) *string {
	return &s
}

func validateInteger(x string, intDigits string, intLookup [256]int) error {
	if len(x) != getIntegerLength(x[0], intDigits, intLookup) {
		return fmt.Errorf("invalid integer part of order key: %s", x)
	}
	return nil
}

func getIntegerLength(head byte, intDigits string, intLookup [256]int) int {
	i := intLookup[head]
	half := len(intDigits) / 2
	if i < half {
		return half - i + 1
	}
	return i - half + 2
}

func getIntegerPart(key string, intDigits string, intLookup [256]int) (string, error) {
	integerPartLength := getIntegerLength(key[0], intDigits, intLookup)
	if integerPartLength > len(key) {
		return "", fmt.Errorf("invalid order key: %s", key)
	}
	return key[:integerPartLength], nil
}

func validateOrderKey(key string, digits string, intDigits string, intLookup [256]int) error {
	if isSmallestInteger(key, digits, intDigits) {
		return fmt.Errorf("invalid order key: %s", key)
	}
	i, err := getIntegerPart(key, intDigits, intLookup)
	if err != nil {
		return err
	}
	f := key[len(i):]
	if len(f) > 0 && f[len(f)-1] == digits[0] {
		return fmt.Errorf("invalid order key: %s", key)
	}
	return nil
}

func incrementInteger(x string, digits string, lookup [256]int, intDigits string, intLookup [256]int) (*string, error) {
	if err := validateInteger(x, intDigits, intLookup); err != nil {
		return nil, err
	}
	head := x[0]
	zero := digits[0]
	trailing := ""
	for i := len(x) - 1; i >= 1; i-- {
		d := lookup[x[i]] + 1
		if d == len(digits) {
			trailing = string(zero) + trailing
		} else {
			r := string(head) + x[1:i] + string(digits[d]) + trailing
			return &r, nil
		}
	}
	headIndex := intLookup[head]
	if headIndex == len(intDigits)-1 {
		return nil, nil
	}
	h := intDigits[headIndex+1]
	lengthDelta := getIntegerLength(h, intDigits, intLookup) - getIntegerLength(head, intDigits, intLookup)
	var r string
	if lengthDelta > 0 {
		r = string(h) + trailing + string(zero)
	} else if lengthDelta < 0 {
		if len(trailing) > 0 {
			r = string(h) + trailing[1:]
		} else {
			r = string(h)
		}
	} else {
		r = string(h) + trailing
	}
	return &r, nil
}

func decrementInteger(x string, digits string, lookup [256]int, intDigits string, intLookup [256]int) (*string, error) {
	if err := validateInteger(x, intDigits, intLookup); err != nil {
		return nil, err
	}
	head := x[0]
	last := digits[len(digits)-1]
	trailing := ""
	for i := len(x) - 1; i >= 1; i-- {
		d := lookup[x[i]] - 1
		if d == -1 {
			trailing = string(last) + trailing
		} else {
			r := string(head) + x[1:i] + string(digits[d]) + trailing
			return &r, nil
		}
	}
	headIndex := intLookup[head]
	if headIndex == 0 {
		return nil, nil
	}
	h := intDigits[headIndex-1]
	lengthDelta := getIntegerLength(h, intDigits, intLookup) - getIntegerLength(head, intDigits, intLookup)
	var r string
	if lengthDelta > 0 {
		r = string(h) + trailing + string(last)
	} else if lengthDelta < 0 {
		if len(trailing) > 0 {
			r = string(h) + trailing[1:]
		} else {
			r = string(h)
		}
	} else {
		r = string(h) + trailing
	}
	return &r, nil
}

func isSmallestInteger(key string, digits string, intDigits string) bool {
	if len(key) != len(intDigits)/2+1 {
		return false
	}
	head := intDigits[0]
	if key[0] != head {
		return false
	}
	for i := 1; i < len(key); i++ {
		if key[i] != digits[0] {
			return false
		}
	}
	return true
}

// generateKeyBetween returns an order key between a and b.
// a and b are *string where nil means unbounded (null in JS).
// digits is the digit set (base62Digits), intDigits is the integer-digit set (base52Digits).
func generateKeyBetween(a, b *string, digits, intDigits string) (string, error) {
	if len(intDigits) < 2 || len(intDigits)%2 != 0 {
		return "", fmt.Errorf("intDigits must be even length >= 2")
	}
	if len(digits) < 2 {
		return "", fmt.Errorf("digits must be >= 2 chars")
	}

	lookup := getDigitIndex(digits)
	intLookup := getDigitIndex(intDigits)

	if a != nil {
		if err := validateOrderKey(*a, digits, intDigits, intLookup); err != nil {
			return "", err
		}
	}
	if b != nil {
		if err := validateOrderKey(*b, digits, intDigits, intLookup); err != nil {
			return "", err
		}
	}
	if a != nil && b != nil && *a > *b {
		a, b = b, a
	}

	if a == nil {
		if b == nil {
			head := intDigits[len(intDigits)/2]
			return string(head) + string(digits[0]), nil
		}

		ib, err := getIntegerPart(*b, intDigits, intLookup)
		if err != nil {
			return "", err
		}
		fb := (*b)[len(ib):]

		if isSmallestInteger(ib, digits, intDigits) {
			r, err := midpoint("", &fb, digits, lookup)
			if err != nil {
				return "", err
			}
			return ib + r, nil
		}
		if ib < *b {
			return ib, nil
		}
		res, err := decrementInteger(ib, digits, lookup, intDigits, intLookup)
		if err != nil {
			return "", err
		}
		if res == nil {
			return "", fmt.Errorf("cannot decrement any more")
		}
		return *res, nil
	}

	if b == nil {
		ia, err := getIntegerPart(*a, intDigits, intLookup)
		if err != nil {
			return "", err
		}
		fa := (*a)[len(ia):]
		i, err := incrementInteger(ia, digits, lookup, intDigits, intLookup)
		if err != nil {
			return "", err
		}
		if i == nil {
			r, err := midpoint(fa, nil, digits, lookup)
			if err != nil {
				return "", err
			}
			return ia + r, nil
		}
		return *i, nil
	}

	ia, err := getIntegerPart(*a, intDigits, intLookup)
	if err != nil {
		return "", err
	}
	fa := (*a)[len(ia):]
	ib, err := getIntegerPart(*b, intDigits, intLookup)
	if err != nil {
		return "", err
	}
	fb := (*b)[len(ib):]

	if ia == ib {
		r, err := midpoint(fa, &fb, digits, lookup)
		if err != nil {
			return "", err
		}
		return ia + r, nil
	}

	i, err := incrementInteger(ia, digits, lookup, intDigits, intLookup)
	if err != nil {
		return "", err
	}
	if i == nil {
		return "", fmt.Errorf("cannot increment any more")
	}
	if *i < *b {
		return *i, nil
	}
	r, err := midpoint(fa, nil, digits, lookup)
	if err != nil {
		return "", err
	}
	return ia + r, nil
}

// generateNKeysBetween returns n order keys between a and b.
func generateNKeysBetween(a, b *string, n int, digits, intDigits string) ([]string, error) {
	if n == 0 {
		return []string{}, nil
	}
	if n == 1 {
		k, err := generateKeyBetween(a, b, digits, intDigits)
		if err != nil {
			return nil, err
		}
		return []string{k}, nil
	}
	if b == nil {
		c, err := generateKeyBetween(a, b, digits, intDigits)
		if err != nil {
			return nil, err
		}
		result := []string{c}
		for i := 0; i < n-1; i++ {
			c, err = generateKeyBetween(&c, b, digits, intDigits)
			if err != nil {
				return nil, err
			}
			result = append(result, c)
		}
		return result, nil
	}
	if a == nil {
		c, err := generateKeyBetween(a, b, digits, intDigits)
		if err != nil {
			return nil, err
		}
		result := []string{c}
		for i := 0; i < n-1; i++ {
			c, err = generateKeyBetween(a, &c, digits, intDigits)
			if err != nil {
				return nil, err
			}
			result = append(result, c)
		}
		// reverse
		for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
			result[i], result[j] = result[j], result[i]
		}
		return result, nil
	}
	mid := n / 2
	c, err := generateKeyBetween(a, b, digits, intDigits)
	if err != nil {
		return nil, err
	}
	left, err := generateNKeysBetween(a, &c, mid, digits, intDigits)
	if err != nil {
		return nil, err
	}
	right, err := generateNKeysBetween(&c, b, n-mid-1, digits, intDigits)
	if err != nil {
		return nil, err
	}
	return append(append(left, c), right...), nil
}

// bootstrapCustomOrder assigns seq_key to all files in a directory using
// the default time-based order (file_type asc, id desc).
// Idempotent: checks Redis flag before running.
func bootstrapCustomOrder(rail miso.Rail, db *gorm.DB, parentFile string, user flow.User) error {
	customOrderKey := "vfm:custom_order:" + parentFile

	// Fast path: check Redis
	exists, _ := redis.GetRedis().Exists(rail.Context(), customOrderKey).Result()
	if exists > 0 {
		// Verify: Redis says bootstrapped, but are files actually bootstrapped?
		has, err := dbquery.NewQuery(rail, db).
			Table("file_info").
			Eq("parent_file", parentFile).
			Where("seq_key != ''").
			Eq("is_logic_deleted", DelN).
			Eq("is_del", 0).
			HasAny()
		if err != nil {
			return err
		}
		if has {
			return nil // truly bootstrapped
		}
		// Redis flag exists but no files have seq_key → bootstrap failed previously
		rail.Warnf("Redis flag exists for dir %v but no files have seq_key, re-bootstrapping", parentFile)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// Read all files in the directory, ordered by default time order
		var files []FileInfo
		if err := dbquery.NewQuery(rail, tx).
			Table("file_info").
			Where("parent_file = ? AND is_logic_deleted = ? AND is_del = ? AND uploader_no = ?",
				parentFile, DelN, 0, user.UserNo).
			Order("file_type asc, id desc").
			ScanVal(&files); err != nil {
			return err
		}

		if len(files) == 0 {
			return nil
		}

		// Distribute keys uniformly using generateNKeysBetween
		keys, err := generateNKeysBetween(nil, nil, len(files), base62Digits, base52Digits)
		if err != nil {
			return err
		}
		for i, f := range files {
			if err := dbquery.NewQuery(rail, tx).
				Table("file_info").
				Set("seq_key", keys[i]).
				Eq("id", f.Id).
				UpdateAny(); err != nil {
				return err
			}
		}

		// Set Redis flag to mark bootstrapped
		if err := redis.GetRedis().Set(rail.Context(), customOrderKey, "1", 0).Err(); err != nil {
			rail.Warnf("failed to set custom order redis flag for dir %v: %v", parentFile, err)
		}
		return nil
	})
}

// assignSeqKeyForNewFile assigns a seq_key for a newly created file/dir
// if its parent directory has custom ordering (checked via Redis flag).
// The new file is placed at the BEGINNING of the list (before the current first file).
func assignSeqKeyForNewFile(rail miso.Rail, tx *gorm.DB, fileKey string, parentFile string, user flow.User) {
	if parentFile == "" {
		return
	}

	// Check Redis flag for custom ordering
	customOrderKey := "vfm:custom_order:" + parentFile
	exists, _ := redis.GetRedis().Exists(rail.Context(), customOrderKey).Result()
	if exists == 0 {
		return // directory not in custom order
	}

	// Lock the directory to prevent concurrent seq_key modifications
	lock := newReorderDirLock(rail, parentFile)
	if err := lock.Lock(); err != nil {
		rail.Warnf("failed to acquire reorder lock for dir %v: %v", parentFile, err)
		return
	}
	defer lock.Unlock()

	// Find the current first file in custom order
	var firstFile struct {
		SeqKey string
	}
	_, err2 := dbquery.NewQuery(rail, tx).
		Table("file_info").
		Select("seq_key").
		Eq("parent_file", parentFile).
		Eq("seq_key !=", "").
		Eq("is_logic_deleted", DelN).
		Eq("is_del", 0).
		Order("seq_key COLLATE utf8mb4_bin asc").
		Limit(1).
		ScanAny(&firstFile)
	if err2 != nil || firstFile.SeqKey == "" {
		// No existing files with seq_key, use default starting key
		keys, err := generateNKeysBetween(nil, nil, 1, base62Digits, base52Digits)
		if err != nil {
			rail.Warnf("failed to generate seq_key for new file %v: %v", fileKey, err)
			return
		}
		if err := dbquery.NewQuery(rail, tx).
			Table("file_info").
			Set("seq_key", keys[0]).
			Eq("uuid", fileKey).
			UpdateAny(); err != nil {
			rail.Warnf("failed to assign seq_key for new file %v: %v", fileKey, err)
		}
		return
	}

	// Prepend: key before the current first key
	newKey, err := generateKeyBetween(nil, &firstFile.SeqKey, base62Digits, base52Digits)
	if err != nil {
		rail.Warnf("failed to generate seq_key for new file %v: %v", fileKey, err)
		return
	}
	if err := dbquery.NewQuery(rail, tx).
		Table("file_info").
		Set("seq_key", newKey).
		Eq("uuid", fileKey).
		UpdateAny(); err != nil {
		rail.Warnf("failed to assign seq_key for new file %v: %v", fileKey, err)
	}
}

// lookupNeighborSeqKey fetches the seq_key of a neighbor file and validates it
// belongs to the expected parent directory (i.e., is a sibling).
// Returns empty string if the file is not found or not a sibling.
func lookupNeighborSeqKey(rail miso.Rail, db *gorm.DB, fileKey string, parentFile string) string {
	var seqKey string
	ok, err := dbquery.NewQuery(rail, db).
		Table("file_info").
		Select("seq_key").
		Eq("uuid", fileKey).
		Eq("parent_file", parentFile).
		Eq("is_del", 0).
		ScanAny(&seqKey)
	if err != nil || !ok {
		return ""
	}
	return seqKey
}

// ReorderFileReq is the request to reorder a file within its parent directory.
type ReorderFileReq struct {
	FileKey    string `json:"fileKey"`    // file to move
	ParentFile string `json:"parentFile"` // parent directory
	BeforeKey  string `json:"beforeKey"`  // fileKey of the file that should be AFTER the moved file; empty = move to bottom
	AfterKey   string `json:"afterKey"`   // fileKey of the file that should be BEFORE the moved file; empty = move to top
}

// ReorderFile moves a file to a new position within its parent directory
// by computing a new fractional-indexing seq_key.
func ReorderFile(rail miso.Rail, db *gorm.DB, req ReorderFileReq, user flow.User) error {
	// 1. Validate file ownership
	f, ok, err := findFile(rail, db, req.FileKey)
	if err != nil {
		return ErrFileNotFound.Wrapf(err, "failed to find file, fileKey: %v", req.FileKey)
	}
	if !ok {
		return ErrFileNotFound.New()
	}
	if f.UploaderNo != user.UserNo {
		return miso.ErrNotPermitted.New()
	}

	// Lock the directory for reorder to prevent concurrent seq_key races
	lock := newReorderDirLock(rail, req.ParentFile)
	if err := lock.Lock(); err != nil {
		return errs.Wrap(err)
	}
	defer lock.Unlock()

	// 2. Determine the new seq_key using midpoint calculation
	var aKey, bKey string // new key must satisfy: aKey < new_key < bKey

	if req.AfterKey == "" && req.BeforeKey == "" {
		return errs.NewErrf("either beforeKey or afterKey must be provided")
	}

	if req.AfterKey != "" {
		aKey = lookupNeighborSeqKey(rail, db, req.AfterKey, req.ParentFile)
	}

	if req.BeforeKey != "" {
		bKey = lookupNeighborSeqKey(rail, db, req.BeforeKey, req.ParentFile)
	}

	// If neighbors don't have seq_key, bootstrap first, then re-lookup
	if (req.AfterKey != "" && aKey == "") || (req.BeforeKey != "" && bKey == "") {
		if err := bootstrapCustomOrder(rail, db, req.ParentFile, user); err != nil {
			return errs.Wrapf(err, "failed to bootstrap custom order for reorder")
		}
		if req.AfterKey != "" {
			aKey = lookupNeighborSeqKey(rail, db, req.AfterKey, req.ParentFile)
		}
		if req.BeforeKey != "" {
			bKey = lookupNeighborSeqKey(rail, db, req.BeforeKey, req.ParentFile)
		}
	}

	var lower, upper *string
	if aKey != "" {
		lower = &aKey
	}
	if bKey != "" {
		upper = &bKey
	}
	newKey, err := generateKeyBetween(lower, upper, base62Digits, base52Digits)
	if err != nil {
		return errs.Wrap(err)
	}

	// 3. Update the file's seq_key
	err = dbquery.NewQuery(rail, db).
		Table("file_info").
		Set("seq_key", newKey).
		Set("update_by", user.Username).
		Eq("uuid", req.FileKey).
		Eq("uploader_no", user.UserNo).
		UpdateAny()
	if err != nil {
		return errs.Wrapf(err, "failed to update seq_key for fileKey: %v", req.FileKey)
	}

	// 4. Set Redis flag indicating custom order exists for this directory
	customOrderKey := "vfm:custom_order:" + req.ParentFile
	if err := redis.GetRedis().Set(rail.Context(), customOrderKey, "1", 0).Err(); err != nil {
		rail.Warnf("failed to set custom order redis flag for dir %v: %v", req.ParentFile, err)
		// non-fatal
	}

	rail.Infof("File %v reordered: new seq_key=%v (afterKey=%v, beforeKey=%v)", req.FileKey, newKey, req.AfterKey, req.BeforeKey)
	return nil
}

// ===== order-by preference cache =====

type OrderByPreferenceReq struct {
	OrderBy string `json:"orderBy"`
	DirKey  string `json:"dirKey"`
}

type GetOrderByPreferenceReq struct {
	DirKey string `form:"dirKey"`
}

type OrderByPreferenceRes struct {
	OrderBy string `json:"orderBy"`
}

// SaveOrderByPreference caches the order-by preference for a user+directory, expires in 30 days.
func SaveOrderByPreference(rail miso.Rail, userNo string, dirKey string, orderBy string) error {
	key := "vfm:orderby:" + userNo + ":" + dirKey
	return redis.Set(rail, key, orderBy, 30*24*time.Hour)
}

// GetOrderByPreference retrieves the cached order-by preference for a user+directory.
// Returns empty string if not found or expired.
func GetOrderByPreference(rail miso.Rail, userNo string, dirKey string) (string, error) {
	key := "vfm:orderby:" + userNo + ":" + dirKey
	val, err := redis.GetStr(key)
	if err != nil {
		if redis.IsNil(err) {
			return "", nil
		}
		return "", err
	}
	return val, nil
}
