package vfm

import (
	"github.com/curtisnewbie/miso/miso"
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
