package vfm

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"

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
	ResManageFiles    = "manage-files"
	ResManageBookmark = "manage-bookmarks"
	ResVfmMaintenance = "vfm:server:maintenance"
)

var (
	ErrUnknown     = miso.NewErrf("Unknown error, please try again")
	ErrUploadFiled = miso.NewErrf("Upload failed, please try again")
)

func RegisterHttpRoutes(rail miso.Rail) error {

	auth.ExposeResourceInfo([]auth.Resource{
		{Code: ResManageFiles, Name: "Manage files"},
		{Code: ResManageBookmark, Name: "Manage Bookmarks"},
		{Code: ResVfmMaintenance, Name: "VFM Server Maintenance"},
	})

	return nil
}

// Check duplicate file in dir.
//
//   - misoapi-http: GET /open/api/file/upload/duplication/preflight
//   - misoapi-desc: Preflight check for duplicate file uploads
//   - misoapi-resource: ref(ResManageFiles)
func ApiPreflightCheckDuplicate(rail miso.Rail, db *gorm.DB, req PreflightCheckReq, user common.User) (bool, error) {
	return FileExists(rail, db, req, user.UserNo)
}

// Fetch parent file.
//
//   - misoapi-http: GET /open/api/file/parent
//   - misoapi-desc: User fetch parent file info
//   - misoapi-resource: ref(ResManageFiles)
func ApiGetParentFile(rail miso.Rail, db *gorm.DB, req FetchParentFileReq, user common.User) (*ParentFileInfo, error) {
	if req.FileKey == "" {
		return nil, miso.NewErrf("fileKey is required")
	}
	pf, e := FindParentFile(rail, db, req, user)
	if e != nil {
		return nil, e
	}
	if pf.Zero {
		return nil, nil
	}
	return &pf, nil
}

// Move file to dir.
//
//   - misoapi-http: POST /open/api/file/move-to-dir
//   - misoapi-desc: User move file into directory
//   - misoapi-resource: ref(ResManageFiles)
func ApiMoveFileToDir(rail miso.Rail, db *gorm.DB, req MoveIntoDirReq, user common.User) (any, error) {
	return nil, MoveFileToDir(rail, db, req, user)
}

type BatchMoveIntoDirReq struct {
	Instructions []MoveIntoDirReq
}

// Move multiple files to dir.
//
//   - misoapi-http: POST /open/api/file/batch-move-to-dir
//   - misoapi-desc: User move files into directory
//   - misoapi-resource: ref(ResManageFiles)
func ApiBatchMoveFileToDir(rail miso.Rail, db *gorm.DB, req BatchMoveIntoDirReq, user common.User) (any, error) {
	for _, r := range req.Instructions {
		if err := MoveFileToDir(rail, db, r, user); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// Make dir.
//
//   - misoapi-http: POST /open/api/file/make-dir
//   - misoapi-desc: User make directory
//   - misoapi-resource: ref(ResManageFiles)
func ApiMakeDir(rail miso.Rail, db *gorm.DB, req MakeDirReq, user common.User) (string, error) {
	return MakeDir(rail, db, req, user)
}

// List dirs.
//
//   - misoapi-http: GET /open/api/file/dir/list
//   - misoapi-desc: User list directories
//   - misoapi-resource: ref(ResManageFiles)
func ApiListDir(rail miso.Rail, db *gorm.DB, user common.User) ([]ListedDir, error) {
	return ListDirs(rail, db, user)
}

// List files.
//
//   - misoapi-http: POST /open/api/file/list
//   - misoapi-desc: User list files
//   - misoapi-resource: ref(ResManageFiles)
func ApiListFiles(rail miso.Rail, db *gorm.DB, req ListFileReq, user common.User) (miso.PageRes[ListedFile], error) {
	return ListFiles(rail, db, req, user)
}

// Delete files.
//
//   - misoapi-http: POST /open/api/file/delete
//   - misoapi-desc: User delete file
//   - misoapi-resource: ref(ResManageFiles)
func ApiDeleteFiles(rail miso.Rail, db *gorm.DB, req DeleteFileReq, user common.User) (any, error) {
	return nil, DeleteFile(rail, db, req, user, nil)
}

// Truncate dir.
//
//   - misoapi-http: POST /open/api/file/dir/truncate
//   - misoapi-desc: User delete truncate directory recursively
//   - misoapi-resource: ref(ResManageFiles)
func ApiTruncateDir(rail miso.Rail, db *gorm.DB, req DeleteFileReq, user common.User) (any, error) {
	return nil, TruncateDir(rail, db, req, user, true)
}

type FetchDirTreeReq struct {
	FileKey string
}

type DirBottomUpTreeNode struct {
	FileKey string
	Name    string
	Child   *DirBottomUpTreeNode
}

// Fetch dir trees, bottom up.
//
//   - misoapi-http: POST /open/api/file/dir/bottom-up-tree
//   - misoapi-desc: Fetch directory tree bottom up.
//   - misoapi-resource: ref(ResManageFiles)
func ApiFetchDirBottomUpTree(inb *miso.Inbound, db *gorm.DB, req FetchDirTreeReq, user common.User) (*DirBottomUpTreeNode, error) {
	return FetchDirTreeBottomUp(inb.Rail(), db, req, user)
}

type DirTopDownTreeNode struct {
	FileKey string
	Name    string
	Child   []*DirTopDownTreeNode
}

// Fetch dir trees, top down.
//
//   - misoapi-http: GET /open/api/file/dir/top-down-tree
//   - misoapi-desc: Fetch directory tree top down.
//   - misoapi-resource: ref(ResManageFiles)
func ApiFetchDirTopDownTree(inb *miso.Inbound, db *gorm.DB, user common.User) (*DirTopDownTreeNode, error) {
	return FetchDirTreeTopDown(inb.Rail(), db, user)
}

type BatchDeleteFileReq struct {
	FileKeys []string
}

// Delete multiple files.
//
//   - misoapi-http: POST /open/api/file/delete/batch
//   - misoapi-desc: User delete file in batch
//   - misoapi-resource: ref(ResManageFiles)
func ApiBatchDeleteFile(rail miso.Rail, db *gorm.DB, req BatchDeleteFileReq, user common.User) (any, error) {
	if len(req.FileKeys) < 31 {
		for i := range req.FileKeys {
			fk := req.FileKeys[i]
			if err := DeleteFile(rail, db, DeleteFileReq{fk}, user, nil); err != nil {
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
			if err := DeleteFile(rrail, db, DeleteFileReq{fk}, user, nil); err != nil {
				rrail.Errorf("failed to delete file, fileKey: %v, %v", fk, err)
			}
		})
	}
	return nil, nil
}

// User Create file.
//
//   - misoapi-http: POST /open/api/file/create
//   - misoapi-desc: User create file
//   - misoapi-resource: ref(ResManageFiles)
func ApiCreateFile(rail miso.Rail, db *gorm.DB, req CreateFileReq, user common.User) (any, error) {
	_, err := CreateFile(rail, db, req, user)
	return nil, err
}

// Update file.
//
//   - misoapi-http: POST /open/api/file/info/update
//   - misoapi-desc: User update file
//   - misoapi-resource: ref(ResManageFiles)
func ApiUpdateFile(rail miso.Rail, db *gorm.DB, req UpdateFileReq, user common.User) (any, error) {
	return nil, UpdateFile(rail, db, req, user)
}

// Generate file token.
//
//   - misoapi-http: POST /open/api/file/token/generate
//   - misoapi-desc: User generate temporary token
//   - misoapi-resource: ref(ResManageFiles)
func ApiGenFileTkn(rail miso.Rail, db *gorm.DB, req GenerateTempTokenReq, user common.User) (string, error) {
	return GenTempToken(rail, db, req, user)
}

// Unzip file in new dir.
//
//   - misoapi-http: POST /open/api/file/unpack
//   - misoapi-desc: User unpack zip
//   - misoapi-resource: ref(ResManageFiles)
func ApiUnpackZip(rail miso.Rail, db *gorm.DB, req UnpackZipReq, user common.User) (any, error) {
	err := UnpackZip(rail, db, user, req)
	return nil, err
}

// Generate file token in QRCode.
//
//   - misoapi-http: GET /open/api/file/token/qrcode
//   - misoapi-desc: User generate qrcode image for temporary token
//   - misoapi-query-doc: token: Generated temporary file key
//   - misoapi-scope: PUBLIC
func ApiGenFileTknQRCode(inb *miso.Inbound) {
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

// List virtual folder brief infos.
//
//   - misoapi-http: GET /open/api/vfolder/brief/owned
//   - misoapi-desc: User list virtual folder briefs
//   - misoapi-resource: ref(ResManageFiles)
func ApiListVFolderBrief(rail miso.Rail, db *gorm.DB, user common.User) ([]VFolderBrief, error) {
	return ListVFolderBrief(rail, db, user)
}

// List virtual folders.
//
//   - misoapi-http: POST /open/api/vfolder/list
//   - misoapi-desc: User list virtual folders
//   - misoapi-resource: ref(ResManageFiles)
func ApiListVFolders(rail miso.Rail, db *gorm.DB, req ListVFolderReq, user common.User) (ListVFolderRes, error) {
	return ListVFolders(rail, db, req, user)
}

// Create virtual folder.
//
//   - misoapi-http: POST /open/api/vfolder/create
//   - misoapi-desc: User create virtual folder
//   - misoapi-resource: ref(ResManageFiles)
func ApiCreateVFolder(rail miso.Rail, db *gorm.DB, req CreateVFolderReq, user common.User) (string, error) {
	return CreateVFolder(rail, db, req, user)
}

// Add file to virtual folder.
//
//   - misoapi-http: POST /open/api/vfolder/file/add
//   - misoapi-desc: User add file to virtual folder
//   - misoapi-resource: ref(ResManageFiles)
func ApiVFolderAddFile(rail miso.Rail, db *gorm.DB, req AddFileToVfolderReq, user common.User) (any, error) {
	return nil, AddFileToVFolder(rail, db, req, user)
}

// Remove file from virtual folder.
//
//   - misoapi-http: POST /open/api/vfolder/file/remove
//   - misoapi-desc: User remove file from virtual folder
//   - misoapi-resource: ref(ResManageFiles)
func ApiVFolderRemoveFile(rail miso.Rail, db *gorm.DB, req RemoveFileFromVfolderReq, user common.User) (any, error) {
	return nil, RemoveFileFromVFolder(rail, db, req, user)
}

// Share virtual folder.
//
//   - misoapi-http: POST /open/api/vfolder/share
//   - misoapi-desc: Share access to virtual folder
//   - misoapi-resource: ref(ResManageFiles)
func ApiShareVFolder(rail miso.Rail, db *gorm.DB, req ShareVfolderReq, user common.User) (any, error) {
	sharedTo, e := vault.FindUser(rail, vault.FindUserReq{Username: &req.Username})
	if e != nil {
		rail.Warnf("Unable to find user, sharedTo: %s, %v", req.Username, e)
		return nil, miso.NewErrf("Failed to find user")
	}
	return nil, ShareVFolder(rail, db, sharedTo, req.FolderNo, user)
}

// Remove access to virtual folder.
//
//   - misoapi-http: POST /open/api/vfolder/access/remove
//   - misoapi-desc: Remove granted access to virtual folder
//   - misoapi-resource: ref(ResManageFiles)
func ApiRemoveVFolderAccess(rail miso.Rail, db *gorm.DB, req RemoveGrantedFolderAccessReq, user common.User) (any, error) {
	return nil, RemoveVFolderAccess(rail, db, req, user)
}

// List accesses to virtual folder.
//
//   - misoapi-http: POST /open/api/vfolder/granted/list
//   - misoapi-desc: List granted access to virtual folder
//   - misoapi-resource: ref(ResManageFiles)
func ApiListVFolderAccess(rail miso.Rail, db *gorm.DB, req ListGrantedFolderAccessReq, user common.User) (ListGrantedFolderAccessRes, error) {
	return ListGrantedFolderAccess(rail, db, req, user)
}

// Remove virtual folder.
//
//   - misoapi-http: POST /open/api/vfolder/remove
//   - misoapi-desc: Remove virtual folder
//   - misoapi-resource: ref(ResManageFiles)
func ApiRemoveVFolder(rail miso.Rail, db *gorm.DB, req RemoveVFolderReq, user common.User) (any, error) {
	return nil, RemoveVFolder(rail, db, user, req)
}

// List gallery brief infos.
//
//   - misoapi-http: GET /open/api/gallery/brief/owned
//   - misoapi-desc: List owned gallery brief info
//   - misoapi-resource: ref(ResManageFiles)
func ApiListGalleryBriefs(rail miso.Rail, db *gorm.DB, user common.User) ([]VGalleryBrief, error) {
	return ListOwnedGalleryBriefs(rail, user, db)
}

// Create gallery.
//
//   - misoapi-http: POST /open/api/gallery/new
//   - misoapi-desc: Create new gallery
//   - misoapi-resource: ref(ResManageFiles)
func ApiCreateGallery(rail miso.Rail, cmd CreateGalleryCmd, db *gorm.DB, user common.User) (*Gallery, error) {
	return CreateGallery(rail, cmd, user, db)
}

// Update gallery.
//
//   - misoapi-http: POST /open/api/gallery/update
//   - misoapi-desc: Update gallery
//   - misoapi-resource: ref(ResManageFiles)
func ApiUpdateGallery(rail miso.Rail, cmd UpdateGalleryCmd, db *gorm.DB, user common.User) (any, error) {
	return nil, UpdateGallery(rail, cmd, user, db)
}

// Delete gallery.
//
//   - misoapi-http: POST /open/api/gallery/delete
//   - misoapi-desc: Delete gallery
//   - misoapi-resource: ref(ResManageFiles)
func ApiDeleteGallery(rail miso.Rail, db *gorm.DB, cmd DeleteGalleryCmd, user common.User) (any, error) {
	return nil, DeleteGallery(rail, db, cmd, user)
}

// List galleries.
//
//   - misoapi-http: POST /open/api/gallery/list
//   - misoapi-desc: List galleries
//   - misoapi-resource: ref(ResManageFiles)
func ApiListGalleries(rail miso.Rail, db *gorm.DB, cmd ListGalleriesCmd, user common.User) (miso.PageRes[VGallery], error) {
	return ListGalleries(rail, cmd, user, db)
}

// Grante access to gallery.
//
//   - misoapi-http: POST /open/api/gallery/access/grant
//   - misoapi-desc: Grant access to the galleries
//   - misoapi-resource: ref(ResManageFiles)
func ApiGranteGalleryAccess(rail miso.Rail, db *gorm.DB, cmd PermitGalleryAccessCmd, user common.User) (any, error) {
	return nil, GrantGalleryAccessToUser(rail, db, cmd, user)
}

// Remove access to gallery.
//
//   - misoapi-http: POST /open/api/gallery/access/remove
//   - misoapi-desc: Remove access to the galleries
//   - misoapi-resource: ref(ResManageFiles)
func ApiRemoveGalleryAccess(rail miso.Rail, db *gorm.DB, cmd RemoveGalleryAccessCmd, user common.User) (any, error) {
	return nil, RemoveGalleryAccess(rail, db, cmd, user)
}

// List accesses to gallery.
//
//   - misoapi-http: POST /open/api/gallery/access/list
//   - misoapi-desc: List granted access to the galleries
//   - misoapi-resource: ref(ResManageFiles)
func ApiListGalleryAccess(rail miso.Rail, db *gorm.DB, cmd ListGrantedGalleryAccessCmd, user common.User) (miso.PageRes[ListedGalleryAccessRes], error) {
	return ListedGrantedGalleryAccess(rail, db, cmd, user)
}

// List images in gallery.
//
//   - misoapi-http: POST /open/api/gallery/images
//   - misoapi-desc: List images of gallery
//   - misoapi-resource: ref(ResManageFiles)
func ApiListGalleryImages(rail miso.Rail, db *gorm.DB, cmd ListGalleryImagesCmd, user common.User) (*ListGalleryImagesResp, error) {
	return ListGalleryImages(rail, db, cmd, user)
}

// Add image to gallery.
//
//   - misoapi-http: POST /open/api/gallery/image/transfer
//   - misoapi-desc: Host selected images on gallery
//   - misoapi-resource: ref(ResManageFiles)
func ApiTransferGalleryImage(rail miso.Rail, db *gorm.DB, cmd TransferGalleryImageReq, user common.User) (any, error) {
	return BatchTransferAsync(rail, cmd, user, db)
}

// List versioned files.
//
//   - misoapi-http: POST /open/api/versioned-file/list
//   - misoapi-desc: List versioned files
//   - misoapi-resource: ref(ResManageFiles)
func ApiListVersionedFile(rail miso.Rail, db *gorm.DB, req ApiListVerFileReq, user common.User) (miso.PageRes[ApiListVerFileRes], error) {
	return ListVerFile(rail, db, req, user)
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

// List history of versioned file.
//
//   - misoapi-http: POST /open/api/versioned-file/history
//   - misoapi-desc: List versioned file history
//   - misoapi-resource: ref(ResManageFiles)
func ApiListVersionedFileHistory(rail miso.Rail, db *gorm.DB, req ApiListVerFileHistoryReq, user common.User) (miso.PageRes[ApiListVerFileHistoryRes], error) {
	return ListVerFileHistory(rail, db, req, user)
}

type ApiQryVerFileAccuSizeReq struct {
	VerFileId string `desc:"versioned file id" valid:"notEmpty"`
}

type ApiQryVerFileAccuSizeRes struct {
	SizeInBytes int64 `desc:"total size in bytes"`
}

// Fetch versioned file accumulated size.
//
//   - misoapi-http: POST /open/api/versioned-file/accumulated-size
//   - misoapi-desc: Query versioned file log accumulated size
//   - misoapi-resource: ref(ResManageFiles)
func ApiQryVersionedFileAccuSize(rail miso.Rail, db *gorm.DB, req ApiQryVerFileAccuSizeReq, user common.User) (ApiQryVerFileAccuSizeRes, error) {
	return CalcVerFileAccuSize(rail, db, req, user)
}

// Create versioned file.
//
//   - misoapi-http: POST /open/api/versioned-file/create
//   - misoapi-desc: Create versioned file
//   - misoapi-resource: ref(ResManageFiles)
func ApiCreateVersionedFile(rail miso.Rail, db *gorm.DB, req ApiCreateVerFileReq, user common.User) (ApiCreateVerFileRes, error) {
	return CreateVerFile(rail, db, req, user)
}

// Update versioned file.
//
//   - misoapi-http: POST /open/api/versioned-file/update
//   - misoapi-desc: Update versioned file
//   - misoapi-resource: ref(ResManageFiles)
func ApiUpdateVersionedFile(rail miso.Rail, db *gorm.DB, req ApiUpdateVerFileReq, user common.User) (any, error) {
	return nil, UpdateVerFile(rail, db, req, user)
}

// Delete versioned file.
//
//   - misoapi-http: POST /open/api/versioned-file/delete
//   - misoapi-desc: Delete versioned file
//   - misoapi-resource: ref(ResManageFiles)
func ApiDelVersionedFile(rail miso.Rail, db *gorm.DB, req ApiDelVerFileReq, user common.User) (any, error) {
	return nil, DelVerFile(rail, db, req, user)
}

// Compensate thumbnail generation
//
//   - misoapi-http: POST /compensate/thumbnail
//   - misoapi-desc: Compensate thumbnail generation
//   - misoapi-resource: ref(ResVfmMaintenance)
func ApiCompensateThumbnail(rail miso.Rail, db *gorm.DB) (any, error) {
	return nil, CompensateThumbnail(rail, db)
}

// Regenerate video thumbnails
//
//   - misoapi-http: POST /compensate/regenerate-video-thumbnails
//   - misoapi-desc: Regenerate video thumbnails
//   - misoapi-resource: ref(ResVfmMaintenance)
func ApiRegenerateVideoThumbnail(rail miso.Rail, db *gorm.DB) error {
	return RegenerateVideoThumbnails(rail, db)
}

type ListBookmarksReq struct {
	Name *string

	Paging      miso.Paging
	Blacklisted bool `gorm:"-" json:"-"`
}

// Upload bookmark file endpoint.
//
//   - misoapi-http: PUT /bookmark/file/upload
//   - misoapi-desc: Upload bookmark file
//   - misoapi-resource: ref(ResManageBookmark)
func ApiUploadBookmarkFile(inb *miso.Inbound, user common.User) (any, error) {
	rail := inb.Rail()
	_, r := inb.Unwrap()
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
//   - misoapi-http: POST /bookmark/list
//   - misoapi-desc: List bookmarks
//   - misoapi-resource: ref(ResManageBookmark)
func ApiListBookmarks(rail miso.Rail, db *gorm.DB, req ListBookmarksReq, user common.User) (any, error) {
	return ListBookmarks(rail, db, req, user.UserNo)
}

type RemoveBookmarkReq struct {
	Id int64
}

// Remove bookmark endpoint.
//
//   - misoapi-http: POST /bookmark/remove
//   - misoapi-desc: Remove bookmark
//   - misoapi-resource: ref(ResManageBookmark)
func ApiRemoveBookmark(rail miso.Rail, db *gorm.DB, req RemoveBookmarkReq, user common.User) (any, error) {
	return nil, RemoveBookmark(rail, db, req.Id, user.UserNo)
}

// List bookmark blacklists.
//
//   - misoapi-http: POST /bookmark/blacklist/list
//   - misoapi-desc: List bookmark blacklist
//   - misoapi-resource: ref(ResManageBookmark)
func ApiListBlacklistedBookmarks(rail miso.Rail, db *gorm.DB, req ListBookmarksReq, user common.User) (any, error) {
	req.Blacklisted = true
	return ListBookmarks(rail, db, req, user.UserNo)
}

// Remove bookmark blacklist.
//
//   - misoapi-http: POST /bookmark/blacklist/remove
//   - misoapi-desc: Remove bookmark blacklist
//   - misoapi-resource: ref(ResManageBookmark)
func ApiRemoveBookmarkBlacklist(rail miso.Rail, db *gorm.DB, req RemoveBookmarkReq, user common.User) (any, error) {
	return nil, RemoveBookmarkBlacklist(rail, db, req.Id, user.UserNo)
}

// List user browse history.
//
//   - misoapi-http: GET /history/list-browse-history
//   - misoapi-desc: List user browse history
//   - misoapi-resource: ref(ResManageFiles)
func ApiListBrowseHistory(rail miso.Rail, db *gorm.DB, user common.User) ([]ListBrowseRecordRes, error) {
	return ListBrowseHistory(rail, db, user)
}

// Record user browse history, only files that are directly owned by the user is recorded.
//
//   - misoapi-http: POST /history/record-browse-history
//   - misoapi-desc: Record user browse history, only files that are directly owned by the user is recorded
//   - misoapi-resource: ref(ResManageFiles)
func ApiRecordBrowseHistory(rail miso.Rail, db *gorm.DB, user common.User, req RecordBrowseHistoryReq) error {
	return RecordBrowseHistory(rail, db, user, req)
}

// Check server maintenance status.
//
//   - misoapi-http: GET /maintenance/status
//   - misoapi-desc: Check server maintenance status
//   - misoapi-resource: ref(ResVfmMaintenance)
func ApiFetchMaintenanceStatus() (MaintenanceStatus, error) {
	return CheckMaintenanceStatus()
}
