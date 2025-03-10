package vfm

import (
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	vault "github.com/curtisnewbie/user-vault/api"
	"gorm.io/gorm"
)

// Internal endpoint. System create file.
//
//   - misoapi-http: POST /internal/v1/file/create
//   - misoapi-desc: Internal endpoint, System create file
func ApiSysCreateFile(rail miso.Rail, db *gorm.DB, req SysCreateFileReq) (string, error) {

	user, err := vault.FindUserCommon(rail, vault.FindUserCommonReq{
		UserNo: req.UserNo,
	})
	if err != nil {
		return "", miso.ErrUnknownError.Wrapf(err, "failed to find user %v", req.UserNo)
	}

	fk, err := CreateFile(rail, db, CreateFileReq{
		Filename:         req.Filename,
		FakeFstoreFileId: req.FakeFstoreFileId,
		ParentFile:       req.ParentFile,
	}, user)

	return fk, err
}

type InternalCheckDuplicateReq struct {
	Filename      string `form:"fileName"`
	ParentFileKey string `form:"parentFileKey"`
	UserNo        string
}

// Internal endpoint, Preflight check for duplicate file uploads.
//
//   - misoapi-http: GET /internal/file/upload/duplication/preflight
//   - misoapi-desc: Internal endpoint, Preflight check for duplicate file uploads
func ApiInternalCheckDuplicate(rail miso.Rail, db *gorm.DB, req InternalCheckDuplicateReq) (bool, error) {
	pcq := PreflightCheckReq{Filename: req.Filename, ParentFileKey: req.ParentFileKey}
	return FileExists(rail, db, pcq, req.UserNo)
}

type InternalCheckFileAccessReq struct {
	FileKey string
	UserNo  string
}

// Internal endpoint, Check if user has access to the file
//
//   - misoapi-http: POST /internal/file/check-access
//   - misoapi-desc: Internal endpoint, Check if user has access to the file
func ApiInternalCheckFileAccess(rail miso.Rail, db *gorm.DB, req InternalCheckFileAccessReq) error {
	return ValidateFileAccess(rail, db, req.FileKey, req.UserNo)
}

type InternalFetchFileInfoReq struct {
	FileKey string `vaild:"notEmpty"`
}

type InternalFetchFileInfoRes struct {
	Name        string
	UploadTime  util.ETime
	SizeInBytes int64
	FileType    string
}

// Internal endpoint. Fetch file info.
//
//   - misoapi-http: POST /internal/file/fetch-info
//   - misoapi-desc: Internal endpoint. Fetch file info.
func ApiInternalFetchFileInfo(rail miso.Rail, db *gorm.DB, req InternalFetchFileInfoReq) (InternalFetchFileInfoRes, error) {
	return InternalFetchFileInfo(rail, db, req)
}

type SysMakeDirReq struct {
	ParentFile string `valid:"notEmpty"`
	UserNo     string `valid:"notEmpty"`
	Name       string `valid:"notEmpty"`
}

// Internal endpoint. System make directory.
//
//   - misoapi-http: POST /internal/v1/file/make-dir
//   - misoapi-desc: Internal endpoint, System make directory.
func ApiSysMakeDir(rail miso.Rail, db *gorm.DB, req SysMakeDirReq) (string, error) {

	user, err := vault.FindUserCommon(rail, vault.FindUserCommonReq{
		UserNo: req.UserNo,
	})
	if err != nil {
		return "", miso.ErrUnknownError.Wrapf(err, "failed to find user %v", req.UserNo)
	}

	dirKey, err := CheckDirExists(rail, db, CheckDirExistsReq{ParentFile: req.ParentFile, Name: req.Name}, user)
	if err != nil {
		return "", err
	}
	if dirKey != "" {
		return dirKey, nil
	}

	fk, err := MakeDir(rail, db, MakeDirReq{
		ParentFile: req.ParentFile,
		Name:       req.Name,
	}, user)

	return fk, err
}
