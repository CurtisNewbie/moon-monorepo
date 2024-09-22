package vfm

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/auth"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	vault "github.com/curtisnewbie/user-vault/api"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

const (
	ManageFilesResource    = "manage-files"
	ResourceManageBookmark = "manage-bookmarks"
)

var (
	ErrUnknown     = miso.NewErrf("Unknown error, please try again")
	ErrUploadFiled = miso.NewErrf("Upload failed, please try again")
)

func RegisterHttpRoutes(rail miso.Rail) error {

	auth.ExposeResourceInfo([]auth.Resource{
		{Code: ManageFilesResource, Name: "Manage files"},
		{Code: ResourceManageBookmark, Name: "Manage Bookmarks"},
	})

	return nil
}

// misoapi-http: GET /open/api/file/upload/duplication/preflight
// misoapi-desc: Preflight check for duplicate file uploads
// misoapi-resource: ref(ManageFilesResource)
func DupPreflightCheckEp(inb *miso.Inbound, req PreflightCheckReq) (bool, error) {
	rail := inb.Rail()
	return FileExists(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

// misoapi-http: GET /open/api/file/parent
// misoapi-desc: User fetch parent file info
// misoapi-resource: ref(ManageFilesResource)
func GetParentFileEp(inb *miso.Inbound, req FetchParentFileReq) (*ParentFileInfo, error) {
	rail := inb.Rail()
	if req.FileKey == "" {
		return nil, miso.NewErrf("fileKey is required")
	}
	pf, e := FindParentFile(rail, mysql.GetMySQL(), req, common.GetUser(rail))
	if e != nil {
		return nil, e
	}
	if pf.Zero {
		return nil, nil
	}
	return &pf, nil
}

// misoapi-http: POST /open/api/file/move-to-dir
// misoapi-desc: User move file into directory
// misoapi-resource: ref(ManageFilesResource)
func MoveFileToDirEp(inb *miso.Inbound, req MoveIntoDirReq) (any, error) {
	rail := inb.Rail()
	return nil, MoveFileToDir(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

type BatchMoveIntoDirReq struct {
	Instructions []MoveIntoDirReq
}

// misoapi-http: POST /open/api/file/batch-move-to-dir
// misoapi-desc: User move files into directory
// misoapi-resource: ref(ManageFilesResource)
func BatchMoveFileToDirEp(inb *miso.Inbound, req BatchMoveIntoDirReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	db := mysql.GetMySQL()
	for _, r := range req.Instructions {
		if err := MoveFileToDir(rail, db, r, user); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// misoapi-http: POST /open/api/file/make-dir
// misoapi-desc: User make directory
// misoapi-resource: ref(ManageFilesResource)
func MakeDirEp(inb *miso.Inbound, req MakeDirReq) (string, error) {
	rail := inb.Rail()
	return MakeDir(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

// misoapi-http: GET /open/api/file/dir/list
// misoapi-desc: User list directories
// misoapi-resource: ref(ManageFilesResource)
func ListDirEp(inb *miso.Inbound) ([]ListedDir, error) {
	rail := inb.Rail()
	return ListDirs(rail, mysql.GetMySQL(), common.GetUser(rail))
}

// misoapi-http: POST /open/api/file/list
// misoapi-desc: User list files
// misoapi-resource: ref(ManageFilesResource)
func ListFilesEp(inb *miso.Inbound, req ListFileReq) (miso.PageRes[ListedFile], error) {
	rail := inb.Rail()
	return ListFiles(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

// misoapi-http: POST /open/api/file/delete
// misoapi-desc: User delete file
// misoapi-resource: ref(ManageFilesResource)
func DeleteFileEp(inb *miso.Inbound, req DeleteFileReq) (any, error) {
	rail := inb.Rail()
	return nil, DeleteFile(rail, mysql.GetMySQL(), req, common.GetUser(rail), nil)
}

// misoapi-http: POST /open/api/file/dir/truncate
// misoapi-desc: User delete truncate directory recursively
// misoapi-resource: ref(ManageFilesResource)
func TruncateDirEp(inb *miso.Inbound, req DeleteFileReq) (any, error) {
	rail := inb.Rail()
	return nil, TruncateDir(rail, mysql.GetMySQL(), req, common.GetUser(rail), true)
}

type BatchDeleteFileReq struct {
	FileKeys []string
}

// misoapi-http: POST /open/api/file/delete/batch
// misoapi-desc: User delete file in batch
// misoapi-resource: ref(ManageFilesResource)
func BatchDeleteFileEp(inb *miso.Inbound, req BatchDeleteFileReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	if len(req.FileKeys) < 31 {
		for i := range req.FileKeys {
			fk := req.FileKeys[i]
			if err := DeleteFile(rail, mysql.GetMySQL(), DeleteFileReq{fk}, user, nil); err != nil {
				rail.Errorf("failed to delete file, fileKey: %v, %v", fk, err)
				return nil, err
			}
		}
		return nil, nil
	}

	// too many file keys, delete files asynchronously
	for i := range req.FileKeys {
		fk := req.FileKeys[i]
		vfmPool.Go(func() {
			rrail := rail.NextSpan()
			if err := DeleteFile(rrail, mysql.GetMySQL(), DeleteFileReq{fk}, user, nil); err != nil {
				rrail.Errorf("failed to delete file, fileKey: %v, %v", fk, err)
			}
		})
	}
	return nil, nil
}

// misoapi-http: POST /open/api/file/create
// misoapi-desc: User create file
// misoapi-resource: ref(ManageFilesResource)
func CreateFileEp(inb *miso.Inbound, req CreateFileReq) (any, error) {
	rail := inb.Rail()
	_, err := CreateFile(rail, mysql.GetMySQL(), req, common.GetUser(rail))
	return nil, err
}

// misoapi-http: POST /open/api/file/info/update
// misoapi-desc: User update file
// misoapi-resource: ref(ManageFilesResource)
func UpdateFileEp(inb *miso.Inbound, req UpdateFileReq) (any, error) {
	rail := inb.Rail()
	return nil, UpdateFile(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

// misoapi-http: POST /open/api/file/token/generate
// misoapi-desc: User generate temporary token
// misoapi-resource: ref(ManageFilesResource)
func GenFileTknEp(inb *miso.Inbound, req GenerateTempTokenReq) (string, error) {
	rail := inb.Rail()
	return GenTempToken(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

// misoapi-http: POST /open/api/file/unpack
// misoapi-desc: User unpack zip
// misoapi-resource: ref(ManageFilesResource)
func UnpackZipEp(inb *miso.Inbound, req UnpackZipReq) (any, error) {
	rail := inb.Rail()
	err := UnpackZip(rail, mysql.GetMySQL(), common.GetUser(rail), req)
	return nil, err
}

// misoapi-http: GET /open/api/file/token/qrcode
// misoapi-desc: User generate qrcode image for temporary token
// misoapi-query-doc: token: Generated temporary file key
// misoapi-scope: PUBLIC
func GenFileTknQRCodeEp(inb *miso.Inbound) {
	w, r := inb.Unwrap()
	rail := inb.Rail()
	token := r.URL.Query().Get("token")
	if util.IsBlankStr(token) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url := miso.GetPropStr(PropVfmSiteHost) + "/fstore/file/raw?key=" + url.QueryEscape(token)
	png, err := qrcode.Encode(url, qrcode.Medium, 512)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rail.Errorf("Failed to generate qrcode image, fileKey: %v, %v", token, err)
		return
	}

	reader := bytes.NewReader(png)
	_, err = io.Copy(w, reader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rail.Errorf("Failed to tranfer qrcode image, fileKey: %v, %v", token, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// misoapi-http: GET /open/api/vfolder/brief/owned
// misoapi-desc: User list virtual folder briefs
// misoapi-resource: ref(ManageFilesResource)
func ListVFolderBriefEp(inb *miso.Inbound) ([]VFolderBrief, error) {
	rail := inb.Rail()
	return ListVFolderBrief(rail, mysql.GetMySQL(), common.GetUser(rail))
}

// misoapi-http: POST /open/api/vfolder/list
// misoapi-desc: User list virtual folders
// misoapi-resource: ref(ManageFilesResource)
func ListVFoldersEp(inb *miso.Inbound, req ListVFolderReq) (ListVFolderRes, error) {
	rail := inb.Rail()
	return ListVFolders(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

// misoapi-http: POST /open/api/vfolder/create
// misoapi-desc: User create virtual folder
// misoapi-resource: ref(ManageFilesResource)
func CreateVFolderEp(inb *miso.Inbound, req CreateVFolderReq) (string, error) {
	rail := inb.Rail()
	return CreateVFolder(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

// misoapi-http: POST /open/api/vfolder/file/add
// misoapi-desc: User add file to virtual folder
// misoapi-resource: ref(ManageFilesResource)
func VFolderAddFileEp(inb *miso.Inbound, req AddFileToVfolderReq) (any, error) {
	rail := inb.Rail()
	return nil, AddFileToVFolder(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

// misoapi-http: POST /open/api/vfolder/file/remove
// misoapi-desc: User remove file from virtual folder
// misoapi-resource: ref(ManageFilesResource)
func VFolderRemoveFileEp(inb *miso.Inbound, req RemoveFileFromVfolderReq) (any, error) {
	rail := inb.Rail()
	return nil, RemoveFileFromVFolder(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

// misoapi-http: POST /open/api/vfolder/share
// misoapi-desc: Share access to virtual folder
// misoapi-resource: ref(ManageFilesResource)
func ShareVFolderEp(inb *miso.Inbound, req ShareVfolderReq) (any, error) {
	rail := inb.Rail()
	sharedTo, e := vault.FindUser(rail, vault.FindUserReq{Username: &req.Username})
	if e != nil {
		rail.Warnf("Unable to find user, sharedTo: %s, %v", req.Username, e)
		return nil, miso.NewErrf("Failed to find user")
	}
	return nil, ShareVFolder(rail, mysql.GetMySQL(), sharedTo, req.FolderNo, common.GetUser(rail))
}

// misoapi-http: POST /open/api/vfolder/access/remove
// misoapi-desc: Remove granted access to virtual folder
// misoapi-resource: ref(ManageFilesResource)
func RemoveVFolderAccessEp(inb *miso.Inbound, req RemoveGrantedFolderAccessReq) (any, error) {
	rail := inb.Rail()
	return nil, RemoveVFolderAccess(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

// misoapi-http: POST /open/api/vfolder/granted/list
// misoapi-desc: List granted access to virtual folder
// misoapi-resource: ref(ManageFilesResource)
func ListVFolderAccessEp(inb *miso.Inbound, req ListGrantedFolderAccessReq) (ListGrantedFolderAccessRes, error) {
	rail := inb.Rail()
	return ListGrantedFolderAccess(rail, mysql.GetMySQL(), req, common.GetUser(rail))
}

// misoapi-http: POST /open/api/vfolder/remove
// misoapi-desc: Remove virtual folder
// misoapi-resource: ref(ManageFilesResource)
func RemoveVFolderEp(inb *miso.Inbound, req RemoveVFolderReq) (any, error) {
	rail := inb.Rail()
	return nil, RemoveVFolder(rail, mysql.GetMySQL(), common.GetUser(rail), req)
}

// misoapi-http: GET /open/api/gallery/brief/owned
// misoapi-desc: List owned gallery brief info
// misoapi-resource: ref(ManageFilesResource)
func ListGalleryBriefsEp(inb *miso.Inbound) ([]VGalleryBrief, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return ListOwnedGalleryBriefs(rail, user, mysql.GetMySQL())
}

// misoapi-http: POST /open/api/gallery/new
// misoapi-desc: Create new gallery
// misoapi-resource: ref(ManageFilesResource)
func CreateGalleryEp(inb *miso.Inbound, cmd CreateGalleryCmd) (*Gallery, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return CreateGallery(rail, cmd, user, mysql.GetMySQL())
}

// misoapi-http: POST /open/api/gallery/update
// misoapi-desc: Update gallery
// misoapi-resource: ref(ManageFilesResource)
func UpdateGalleryEp(inb *miso.Inbound, cmd UpdateGalleryCmd) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	e := UpdateGallery(rail, cmd, user, mysql.GetMySQL())
	return nil, e
}

// misoapi-http: POST /open/api/gallery/delete
// misoapi-desc: Delete gallery
// misoapi-resource: ref(ManageFilesResource)
func DeleteGalleryEp(inb *miso.Inbound, cmd DeleteGalleryCmd) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	e := DeleteGallery(rail, mysql.GetMySQL(), cmd, user)
	return nil, e
}

// misoapi-http: POST /open/api/gallery/list
// misoapi-desc: List galleries
// misoapi-resource: ref(ManageFilesResource)
func ListGalleriesEp(inb *miso.Inbound, cmd ListGalleriesCmd) (miso.PageRes[VGallery], error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return ListGalleries(rail, cmd, user, mysql.GetMySQL())
}

// misoapi-http: POST /open/api/gallery/access/grant
// misoapi-desc: Grant access to the galleries
// misoapi-resource: ref(ManageFilesResource)
func GranteGalleryAccessEp(inb *miso.Inbound, cmd PermitGalleryAccessCmd) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	e := GrantGalleryAccessToUser(rail, mysql.GetMySQL(), cmd, user)
	return nil, e
}

// misoapi-http: POST /open/api/gallery/access/remove
// misoapi-desc: Remove access to the galleries
// misoapi-resource: ref(ManageFilesResource)
func RemoveGalleryAccessEp(inb *miso.Inbound, cmd RemoveGalleryAccessCmd) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	e := RemoveGalleryAccess(rail, mysql.GetMySQL(), cmd, user)
	return nil, e
}

// misoapi-http: POST /open/api/gallery/access/list
// misoapi-desc: List granted access to the galleries
// misoapi-resource: ref(ManageFilesResource)
func ListGalleryAccessEp(inb *miso.Inbound, cmd ListGrantedGalleryAccessCmd) (miso.PageRes[ListedGalleryAccessRes], error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return ListedGrantedGalleryAccess(rail, mysql.GetMySQL(), cmd, user)
}

// misoapi-http: POST /open/api/gallery/images
// misoapi-desc: List images of gallery
// misoapi-resource: ref(ManageFilesResource)
func ListGalleryImagesEp(inb *miso.Inbound, cmd ListGalleryImagesCmd) (*ListGalleryImagesResp, error) {
	rail := inb.Rail()
	return ListGalleryImages(rail, mysql.GetMySQL(), cmd, common.GetUser(rail))
}

// misoapi-http: POST /open/api/gallery/image/transfer
// misoapi-desc: Host selected images on gallery
// misoapi-resource: ref(ManageFilesResource)
func TransferGalleryImageEp(inb *miso.Inbound, cmd TransferGalleryImageReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return BatchTransferAsync(rail, cmd, user, mysql.GetMySQL())
}

// misoapi-http: POST /open/api/versioned-file/list
// misoapi-desc: List versioned files
// misoapi-resource: ref(ManageFilesResource)
func ApiListVersionedFile(inb *miso.Inbound, req ApiListVerFileReq) (miso.PageRes[ApiListVerFileRes], error) {
	return ListVerFile(inb.Rail(), mysql.GetMySQL(), req, common.GetUser(inb.Rail()))
}

type ApiListVerFileHistoryReq struct {
	Paging    miso.Paging `desc:"paging params"`
	VerFileId string      `desc:"versioned file id" valid:"notEmpty"`
}

type ApiListVerFileHistoryRes struct {
	Name        string     `desc:"file name"`
	FileKey     string     `desc:"file key"`
	SizeInBytes int64      `desc:"size in bytes"`
	UploadTime  util.ETime `desc:"last upload time"`
	Thumbnail   string     `desc:"thumbnail token"`
}

// misoapi-http: POST /open/api/versioned-file/history
// misoapi-desc: List versioned file history
// misoapi-resource: ref(ManageFilesResource)
func ApiListVersionedFileHistory(inb *miso.Inbound, req ApiListVerFileHistoryReq) (miso.PageRes[ApiListVerFileHistoryRes], error) {
	return ListVerFileHistory(inb.Rail(), mysql.GetMySQL(), req, common.GetUser(inb.Rail()))
}

type ApiQryVerFileAccuSizeReq struct {
	VerFileId string `desc:"versioned file id" valid:"notEmpty"`
}

type ApiQryVerFileAccuSizeRes struct {
	SizeInBytes int64 `desc:"total size in bytes"`
}

// misoapi-http: POST /open/api/versioned-file/accumulated-size
// misoapi-desc: Query versioned file log accumulated size
// misoapi-resource: ref(ManageFilesResource)
func ApiQryVersionedFileAccuSize(inb *miso.Inbound, req ApiQryVerFileAccuSizeReq) (ApiQryVerFileAccuSizeRes, error) {
	return CalcVerFileAccuSize(inb.Rail(), mysql.GetMySQL(), req, common.GetUser(inb.Rail()))
}

// misoapi-http: POST /open/api/versioned-file/create
// misoapi-desc: Create versioned file
// misoapi-resource: ref(ManageFilesResource)
func ApiCreateVersionedFile(inb *miso.Inbound, req ApiCreateVerFileReq) (ApiCreateVerFileRes, error) {
	return CreateVerFile(inb.Rail(), mysql.GetMySQL(), req, common.GetUser(inb.Rail()))
}

// misoapi-http: POST /open/api/versioned-file/update
// misoapi-desc: Update versioned file
// misoapi-resource: ref(ManageFilesResource)
func ApiUpdateVersionedFile(inb *miso.Inbound, req ApiUpdateVerFileReq) (any, error) {
	return nil, UpdateVerFile(inb.Rail(), mysql.GetMySQL(), req, common.GetUser(inb.Rail()))
}

// misoapi-http: POST /open/api/versioned-file/delete
// misoapi-desc: Delete versioned file
// misoapi-resource: ref(ManageFilesResource)
func ApiDelVersionedFile(inb *miso.Inbound, req ApiDelVerFileReq) (any, error) {
	return nil, DelVerFile(inb.Rail(), mysql.GetMySQL(), req, common.GetUser(inb.Rail()))
}

// misoapi-http: POST /compensate/thumbnail
// misoapi-desc: Compensate thumbnail generation
func CompensateThumbnailEp(rail miso.Rail, db *gorm.DB) (any, error) {
	return nil, CompensateThumbnail(rail, db)
}

// misoapi-http: POST /compensate/dir/calculate-size
// misoapi-desc: Calculate size of all directories recursively
func ImMemBatchCalcDirSizeEp(rail miso.Rail, db *gorm.DB) (any, error) {
	return nil, ImMemBatchCalcDirSize(rail, mysql.GetMySQL())
}

type ListBookmarksReq struct {
	Name *string

	Paging      miso.Paging
	Blacklisted bool `gorm:"-" json:"-"`
}

// Upload bookmark file endpoint.
//
// misoapi-http: PUT /bookmark/file/upload
// misoapi-desc: Upload bookmark file
// misoapi-resource: ref(ResourceManageBookmark)
func UploadBookmarkFileEp(inb *miso.Inbound) (any, error) {
	rail := inb.Rail()
	_, r := inb.Unwrap()
	user := common.GetUser(rail)
	path, err := TransferTmpFile(rail, r.Body)
	if err != nil {
		return nil, err
	}
	defer os.Remove(path)

	lock := redis.NewRLock(rail, "docindexer:bookmark:"+user.UserNo)
	if err := lock.Lock(); err != nil {
		rail.Errorf("failed to lock for bookmark upload, user: %v, %v", user.Username, err)
		return nil, miso.NewErrf("Please try again later")
	}
	defer lock.Unlock()

	if err := ProcessUploadedBookmarkFile(rail, path, user); err != nil {
		rail.Errorf("ProcessUploadedBookmarkFile failed, user: %v, path: %v, %v", user.Username, path, err)
		return nil, miso.NewErrf("Failed to parse bookmark file")
	}

	return nil, nil
}

// List bookmarks endpoint.
//
// misoapi-http: POST /bookmark/list
// misoapi-desc: List bookmarks
// misoapi-resource: ref(ResourceManageBookmark)
func ListBookmarksEp(inb *miso.Inbound, req ListBookmarksReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return ListBookmarks(rail, mysql.GetMySQL(), req, user.UserNo)
}

type RemoveBookmarkReq struct {
	Id int64
}

// Remove bookmark endpoint.
//
// misoapi-http: POST /bookmark/remove
// misoapi-desc: Remove bookmark
// misoapi-resource: ref(ResourceManageBookmark)
func RemoveBookmarkEp(inb *miso.Inbound, req RemoveBookmarkReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return nil, RemoveBookmark(rail, mysql.GetMySQL(), req.Id, user.UserNo)
}

// misoapi-http: POST /bookmark/blacklist/list
// misoapi-desc: List bookmark blacklist
// misoapi-resource: ref(ResourceManageBookmark)
func ListBlacklistedBookmarksEp(inb *miso.Inbound, req ListBookmarksReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	req.Blacklisted = true
	return ListBookmarks(rail, mysql.GetMySQL(), req, user.UserNo)
}

// misoapi-http: POST /bookmark/blacklist/remove
// misoapi-desc: Remove bookmark blacklist
// misoapi-resource: ref(ResourceManageBookmark)
func RemoveBookmarkBlacklistEp(inb *miso.Inbound, req RemoveBookmarkReq) (any, error) {
	rail := inb.Rail()
	user := common.GetUser(rail)
	return nil, RemoveBookmarkBlacklist(rail, mysql.GetMySQL(), req.Id, user.UserNo)
}
