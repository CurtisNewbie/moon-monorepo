package fstore

import (
	"github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/mini-fstore/internal/config"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/errs"
	"gorm.io/gorm"
)

var (
	ErrInvalidAuth = errs.NewErrf("Invalid authorization").WithCode(api.InvalidAuthorization)
)

type BackupFileInf struct {
	Id     int64  `json:"id"`
	FileId string `json:"fileId"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Size   int64  `json:"size"`
	Md5    string `json:"md5"`
}

type ListBackupFileReq struct {
	Limit    int64 `json:"limit"`
	IdOffset int   `json:"idOffset"`
}

type ListBackupFileResp struct {
	Files []BackupFileInf `json:"files"`
}

func ListBackupFiles(rail miso.Rail, tx *gorm.DB, req ListBackupFileReq) (ListBackupFileResp, error) {
	var files []BackupFileInf
	_, err := dbquery.NewQuery(rail, tx).
		Table("file").
		Select("id, file_id, name, status, size, md5").
		Where("id > ?", req.IdOffset).
		Order("id ASC").
		Limit(int(req.Limit)).
		Scan(&files)
	if err != nil {
		return ListBackupFileResp{}, ErrUnknownError.WithInternalMsg("Failed to list back up files, req %+v, %v", req, err)
	}
	if files == nil {
		files = []BackupFileInf{}
	}
	return ListBackupFileResp{Files: files}, nil
}

func CheckBackupAuth(rail miso.Rail, auth string) error {
	rail.Debugf("Checking backup auth, auth: %v", auth)
	if auth == "" {
		return ErrInvalidAuth.WithInternalMsg("auth is empty")
	}
	secret := miso.GetPropStr(config.PropBackupAuthSecret)
	if secret != auth {
		return ErrInvalidAuth.WithInternalMsg("secret != auth")
	}
	return nil
}
