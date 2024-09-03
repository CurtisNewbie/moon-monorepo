package vfm

import (
	"fmt"

	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"gorm.io/gorm"
)

func CheckVerFile(rail miso.Rail, db *gorm.DB, fileKey string, userNo string) (*FileInfo, error) {
	f, err := findFile(rail, db, fileKey)
	if err != nil {
		return f, miso.NewErrf("File not found", "fileKey: %v, %v", fileKey, err)
	}
	if f == nil {
		return f, miso.NewErrf("File not found", "fileKey: %v", fileKey)
	}
	if f.UploaderNo != userNo {
		return f, miso.NewErrf("Not permitted")
	}
	if f.FileType != FileTypeFile {
		return f, miso.NewErrf("Illegal File Type")
	}
	if f.IsLogicDeleted == LDelY {
		return f, miso.NewErrf("File already deleted")
	}
	return f, nil
}

type ApiCreateVerFileReq struct {
	Filename         string `json:"filename" valid:"notEmpty"`
	FakeFstoreFileId string `json:"fstoreFileId" valid:"notEmpty"`
}

type ApiCreateVerFileRes struct {
	VerFileId string `desc:"Versioned File Id"`
}

func CreateVerFile(rail miso.Rail, db *gorm.DB, req ApiCreateVerFileReq, user common.User) (ApiCreateVerFileRes, error) {
	var res ApiCreateVerFileRes

	fk, err := CreateFile(rail, db, CreateFileReq{Filename: req.Filename, FakeFstoreFileId: req.FakeFstoreFileId, Hidden: true}, user)
	if err != nil {
		return res, fmt.Errorf("failed to CreateFile, %v, %#v", err, req)
	}

	verFileId := util.GenIdP("verf_")
	rail.Infof("file_info record created, fileKey: %s, req: %#v", fk, req)

	f, err := findFile(rail, db, fk)
	if err != nil {
		return res, fmt.Errorf("failed to find file_info record, %s, %v", fk, err)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		err := tx.Exec(`
			INSERT INTO versioned_file (ver_file_id,file_key,name,size_in_bytes,uploader_no,uploader_name,upload_time,created_by)
			VALUES (?,?,?,?,?,?,?,?)
		`, verFileId, f.Uuid, f.Name, f.SizeInBytes, f.UploaderNo, f.UploaderName, util.Now(), user.Username).Error
		if err != nil {
			return fmt.Errorf("failed to insert versioned_file, req: #%v, %w", req, err)
		}
		if err := SaveVerFileLog(rail, tx,
			SaveVerFileLogReq{VerFileId: verFileId, FileKey: fk, Username: user.Username}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return res, err
	}
	rail.Infof("Versioned file %v created for %v", verFileId, f.Uuid)
	return ApiCreateVerFileRes{VerFileId: verFileId}, nil
}

type ApiUpdateVerFileReq struct {
	VerFileId        string `valid:"notEmpty" desc:"versioned file id"`
	Filename         string `json:"filename" valid:"notEmpty"`
	FakeFstoreFileId string `json:"fstoreFileId" valid:"notEmpty"`
}

type UpdateVerFileInf struct {
	FileKey    string
	UploaderNo string
	Deleted    bool
}

func NewVerFileLock(rail miso.Rail, verFileId string) *redis.RLock {
	return redis.NewRLockf(rail, "rlock:version_file:%v", verFileId)
}

func UpdateVerFile(rail miso.Rail, db *gorm.DB, req ApiUpdateVerFileReq, user common.User) error {
	lock := NewVerFileLock(rail, req.VerFileId)
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	fk, err := CreateFile(rail, db, CreateFileReq{Filename: req.Filename, FakeFstoreFileId: req.FakeFstoreFileId, Hidden: true}, user)
	if err != nil {
		return fmt.Errorf("failed to CreateFile, %v, req: %#v", err, req)
	}
	rail.Infof("file_info record created, fileKey: %s, req: %#v", fk, req)

	f, err := findFile(rail, db, fk)
	if err != nil {
		return fmt.Errorf("failed to find file_info record, %s, %v", fk, err)
	}

	var uvf UpdateVerFileInf
	tx := db.Raw(`
		SELECT file_key,uploader_no,deleted
		FROM versioned_file
		WHERE ver_file_id = ?
	`, req.VerFileId).Scan(&uvf)
	if tx.Error != nil {
		return fmt.Errorf("failed to query versioned_file, req: %#v, %v", req, tx.Error)
	}
	if tx.RowsAffected < 1 {
		return miso.NewErrf("File not found", "ver_file_id not found, %v", req.VerFileId)
	}
	if uvf.UploaderNo != user.UserNo {
		return miso.NewErrf("Not permitted")
	}
	if uvf.Deleted {
		return miso.NewErrf("File already deleted")
	}
	if uvf.FileKey == fk {
		rail.Infof("Versioned File %v already using file_key: %v", req.VerFileId, fk)
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		svlr := SaveVerFileLogReq{VerFileId: req.VerFileId, FileKey: fk, Username: user.Username}
		if err := SaveVerFileLog(rail, tx, svlr); err != nil {
			return err
		}

		err = tx.Exec(`
			UPDATE versioned_file
			SET file_key = ?,
				name = ?,
				size_in_bytes = ?,
				upload_time = ?,
				updated_by = ?
			WHERE ver_file_id = ?
		`, f.Uuid, f.Name, f.SizeInBytes, util.Now(), user.Username, req.VerFileId).Error
		if err != nil {
			return fmt.Errorf("failed to update versioned_file, req: #%v, %w", req, err)
		}

		rail.Infof("Versioned file %v updated using %v", req.VerFileId, f.Uuid)

		return nil
	})

}

type ApiListVerFileReq struct {
	Paging miso.Paging `desc:"paging params"`
	Name   *string     `desc:"file name"`
}

type ApiListVerFileRes struct {
	VerFileId   string     `desc:"versioned file id"`
	Name        string     `desc:"file name"`
	FileKey     string     `desc:"file key"`
	SizeInBytes int64      `desc:"size in bytes"`
	UploadTime  util.ETime `desc:"last upload time"`
	CreateTime  util.ETime `desc:"create time of the versioned file record"`
	UpdateTime  util.ETime `desc:"Update time of the versioned file record"`
	Thumbnail   string     `desc:"thumbnail token"`
}

func ListVerFile(rail miso.Rail, db *gorm.DB, req ApiListVerFileReq, user common.User) (miso.PageRes[ApiListVerFileRes], error) {
	return mysql.NewPageQuery[ApiListVerFileRes]().
		WithPage(req.Paging).
		WithBaseQuery(func(tx *gorm.DB) *gorm.DB {
			tx = tx.Table(`versioned_file f`).
				Joins("LEFT JOIN file_info fi on f.file_key = fi.uuid").
				Where(`f.uploader_no = ?`, user.UserNo).
				Where(`f.deleted = 0`)

			if req.Name != nil && *req.Name != "" {
				tx = tx.Where("match(f.name) against (? IN NATURAL LANGUAGE MODE)", *req.Name)
			}

			return tx
		}).
		WithSelectQuery(func(tx *gorm.DB) *gorm.DB {
			return tx.Select(`f.ver_file_id,f.name,f.file_key,f.size_in_bytes,f.upload_time,f.create_time,f.update_time,fi.thumbnail`)
		}).
		ForEach(func(t ApiListVerFileRes) ApiListVerFileRes {
			if t.Thumbnail != "" {
				tkn, err := GetFstoreTmpToken(rail, t.Thumbnail, "")
				if err != nil {
					rail.Errorf("failed to generate file token for thumbnail: %v, %v", t.Thumbnail, err)
					t.Thumbnail = ""
				} else {
					t.Thumbnail = tkn
				}
			}
			return t
		}).Exec(rail, db)
}

type SaveVerFileLogReq struct {
	VerFileId string
	FileKey   string
	Username  string
}

func SaveVerFileLog(rail miso.Rail, db *gorm.DB, req SaveVerFileLogReq) error {
	err := db.Exec(`INSERT INTO versioned_file_log (ver_file_id, file_key, created_by) VALUES (?,?,?)`,
		req.VerFileId, req.FileKey, req.Username).Error
	if err != nil {
		return fmt.Errorf("failed to save versioned_file_log, %#v, %v", req, err)
	}
	return err
}

type ApiDelVerFileReq struct {
	VerFileId string `desc:"Versioned File Id" valid:"notEmpty"`
}

func DelVerFile(rail miso.Rail, db *gorm.DB, req ApiDelVerFileReq, user common.User) error {
	lock := NewVerFileLock(rail, req.VerFileId)
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	var uvf UpdateVerFileInf
	tx := db.Raw(`
		SELECT file_key,uploader_no,deleted
		FROM versioned_file
		WHERE ver_file_id = ?
	`, req.VerFileId).Scan(&uvf)
	if tx.Error != nil {
		return fmt.Errorf("failed to query versioned_file, req: %#v, %v", req, tx.Error)
	}
	if tx.RowsAffected < 1 {
		return miso.NewErrf("File not found", "ver_file_id not found, %v", req.VerFileId)
	}
	if uvf.UploaderNo != user.UserNo {
		return miso.NewErrf("Not permitted")
	}
	if uvf.Deleted {
		return miso.NewErrf("File already deleted")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		err := tx.Exec(`UPDATE versioned_file SET deleted = 1, updated_by = ?, delete_time = ? WHERE ver_file_id = ?`,
			user.Username, util.Now(), req.VerFileId).Error
		if err != nil {
			return fmt.Errorf("failed to mark versioend_file deleted, %v, %w", req.VerFileId, err)
		}

		var fks []string
		if err := tx.Raw(`SELECT vf.file_key FROM versioned_file_log vf WHERE vf.ver_file_id = ?`, req.VerFileId).
			Scan(&fks).Error; err != nil {
			return fmt.Errorf("failed to query versioend_file_log, %v, %w", req.VerFileId, err)
		}

		for _, fk := range fks {
			if err := DeleteFile(rail, tx, DeleteFileReq{Uuid: fk}, user, nil); err != nil {
				return fmt.Errorf("failed to delete file in versioend_file_log, %v, %v, %w", req.VerFileId, fk, err)
			}
		}
		return nil
	})
}

func ListVerFileHistory(rail miso.Rail, db *gorm.DB, req ApiListVerFileHistoryReq, user common.User) (miso.PageRes[ApiListVerFileHistoryRes], error) {
	if err := checkVerFileAccess(rail, db, user.UserNo, req.VerFileId); err != nil {
		return miso.PageRes[ApiListVerFileHistoryRes]{}, err
	}

	return mysql.NewPageQuery[ApiListVerFileHistoryRes]().
		WithPage(req.Paging).
		WithBaseQuery(func(tx *gorm.DB) *gorm.DB {
			tx = tx.Table(`versioned_file_log f`).
				Joins("LEFT JOIN file_info fi on f.file_key = fi.uuid").
				Where(`f.ver_file_id = ?`, req.VerFileId)
			return tx
		}).
		WithSelectQuery(func(tx *gorm.DB) *gorm.DB {
			return tx.Select(`f.file_key,fi.name,fi.size_in_bytes,fi.upload_time,fi.thumbnail`).Order("f.id DESC")
		}).
		ForEach(func(t ApiListVerFileHistoryRes) ApiListVerFileHistoryRes {
			if t.Thumbnail != "" {
				tkn, err := GetFstoreTmpToken(rail, t.Thumbnail, "")
				if err != nil {
					rail.Errorf("failed to generate file token for thumbnail: %v, %v", t.Thumbnail, err)
					t.Thumbnail = ""
				} else {
					t.Thumbnail = tkn
				}
			}
			return t
		}).Exec(rail, db)
}

func CalcVerFileAccuSize(rail miso.Rail, db *gorm.DB, req ApiQryVerFileAccuSizeReq, user common.User) (ApiQryVerFileAccuSizeRes, error) {
	if err := checkVerFileAccess(rail, db, user.UserNo, req.VerFileId); err != nil {
		return ApiQryVerFileAccuSizeRes{}, err
	}

	var total int64
	err := db.Raw(`
	SELECT sum(fi.size_in_bytes) FROM versioned_file_log f
	LEFT JOIN file_info fi ON f.file_key = fi.uuid
	WHERE f.ver_file_id = ?
	`, req.VerFileId).Scan(&total).Error
	if err != nil {
		return ApiQryVerFileAccuSizeRes{}, err
	}
	return ApiQryVerFileAccuSizeRes{SizeInBytes: total}, nil
}

func checkVerFileAccess(rail miso.Rail, db *gorm.DB, userNo string, verFileId string) error {
	var id int
	t := db.Raw(`SELECT id FROM versioned_file WHERE uploader_no = ? and ver_file_id = ? and deleted = 0 LIMIT 1`,
		userNo, verFileId).Scan(&id)
	if t.Error != nil {
		return t.Error
	}
	if t.RowsAffected < 1 {
		return miso.NewErrf("Versioned file not found")
	}
	return nil
}
