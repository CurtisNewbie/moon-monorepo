package vfm

import (
	"fmt"

	"github.com/curtisnewbie/miso/errs"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/async"
	"github.com/curtisnewbie/miso/util/atom"
	"github.com/curtisnewbie/miso/util/randutil"
	"gorm.io/gorm"
)

// Gallery
type Gallery struct {
	Id         int64  `json:"id"`
	GalleryNo  string `json:"galleryNo"`
	UserNo     string `json:"userNo"`
	Name       string `json:"name"`
	DirFileKey string `json:"dirFileKey"`
	CreateBy   string `json:"createBy"`
	UpdateBy   string `json:"updateBy"`
	IsDel      bool   `json:"isDel"`
}

func (Gallery) TableName() string {
	return "gallery"

}

type CreateGalleryCmd struct {
	Name string `json:"name" validation:"notEmpty"`
}

type CreateGalleryForDirCmd struct {
	DirName    string
	DirFileKey string
	Username   string
	UserNo     string
}

type UpdateGalleryCmd struct {
	GalleryNo string `json:"galleryNo" validation:"notEmpty"`
	Name      string `json:"name" validation:"notEmpty"`
}

type ListGalleriesCmd struct {
	Paging miso.Paging `json:"paging"`
}

type DeleteGalleryCmd struct {
	GalleryNo string `json:"galleryNo" validation:"notEmpty"`
}

type VGalleryBrief struct {
	GalleryNo string `json:"galleryNo"`
	Name      string `json:"name"`
}

type VGallery struct {
	ID             int64     `json:"id"`
	GalleryNo      string    `json:"galleryNo"`
	UserNo         string    `json:"userNo"`
	Name           string    `json:"name"`
	CreateTime     atom.Time `json:"-"`
	UpdateTime     atom.Time `json:"-"`
	CreateBy       string    `json:"createBy"`
	UpdateBy       string    `json:"updateBy"`
	IsOwner        bool      `json:"isOwner"`
	CreateTimeStr  string    `json:"createTime"`
	UpdateTimeStr  string    `json:"updateTime"`
	DirFileKey     string    `json:"dirFileKey"`
	ThumbnailToken string    `json:"thumbnailToken"`
}

// List owned gallery briefs
func ListOwnedGalleryBriefs(rail miso.Rail, user common.User, tx *gorm.DB) ([]VGalleryBrief, error) {
	var briefs []VGalleryBrief
	err := dbquery.NewQuery(rail, tx).
		Raw(`select gallery_no, name from gallery where user_no = ? AND is_del = 0`, user.UserNo).
		ScanVal(&briefs)
	if err != nil {
		return nil, err
	}
	if briefs == nil {
		briefs = []VGalleryBrief{}
	}

	return briefs, nil
}

/* List Galleries */
func ListGalleries(rail miso.Rail, cmd ListGalleriesCmd, user common.User, db *gorm.DB) (miso.PageRes[VGallery], error) {
	return dbquery.NewPagedQuery[VGallery](db).
		WithBaseQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.Table("gallery g").
				Where("g.is_del = 0").
				Where("g.user_no = ? OR EXISTS (select * from gallery_user_access ga where ga.user_no = ? AND ga.is_del = 0 AND ga.gallery_no = g.gallery_no)", user.UserNo, user.UserNo)
		}).
		WithSelectQuery(func(q *dbquery.Query) *dbquery.Query {
			return q.Select("g.*").Order("g.update_time DESC")
		}).
		TransformAsync(func(g VGallery) async.Future[VGallery] {
			if g.UserNo == user.UserNo {
				g.IsOwner = true
			}
			if !g.IsOwner {
				g.DirFileKey = ""
			}
			g.CreateTimeStr = g.CreateTime.FormatClassic()
			g.UpdateTimeStr = g.UpdateTime.FormatClassic()

			return async.Run[VGallery](func() (VGallery, error) {
				var thumbnailFileId string
				ok, err := dbquery.NewQuery(rail).
					Table("gallery_image gi").
					Joins("LEFT JOIN file_info fi ON gi.file_key = fi.uuid").
					Eq("gi.gallery_no", g.GalleryNo).
					Eq("fi.is_logic_deleted", false).
					Select("fi.thumbnail").
					OrderAsc("gi.name").
					Limit(1).
					ScanAny(&thumbnailFileId)

				if err != nil {
					rail.Errorf("Failed to find thumbnail file_key for gallery: %v, %v", g.GalleryNo, err)
				} else if ok {
					tkn, err := GetFstoreTmpToken(rail.NextSpan(), thumbnailFileId, "thumbnail.jpg")
					if err != nil {
						rail.Errorf("Failed to generate tmp token for gallery thumbnail, galleryNo: %v, thumbnailFileId: %v, %v",
							g.GalleryNo, thumbnailFileId, err)
					} else {
						g.ThumbnailToken = tkn
					}
				}
				return g, nil
			})
		}).
		Scan(rail, cmd.Paging)
}

func GalleryNoOfDir(rail miso.Rail, dirFileKey string, tx *gorm.DB) (string, error) {
	var gallery Gallery
	err := dbquery.NewQuery(rail, tx).
		Raw(`SELECT g.gallery_no from gallery g WHERE g.dir_file_key = ? and g.is_del = 0 limit 1`, dirFileKey).
		ScanVal(&gallery)
	if err != nil {
		return "", err
	}

	return gallery.GalleryNo, nil
}

// Check if the name is already used by current user
func IsGalleryNameUsed(rail miso.Rail, name string, userNo string, tx *gorm.DB) (bool, error) {
	var gallery Gallery
	n, err := dbquery.NewQuery(rail, tx).
		Raw(`SELECT g.id from gallery g WHERE g.user_no = ? and g.name = ? AND g.is_del = 0`, userNo, name).
		Scan(&gallery)

	if err != nil {
		return false, err
	}

	return n > 0, nil
}

func NewGalleryDirLock(rail miso.Rail, dirFileKey string) *redis.RLock {
	return redis.NewRLockf(rail, "fantahsea:gallery:create:dir:%v", dirFileKey)
}

// Create a new Gallery for dir
func CreateGalleryForDir(rail miso.Rail, cmd CreateGalleryForDirCmd, db *gorm.DB) (string, error) {
	lock := NewGalleryDirLock(rail, cmd.DirFileKey)
	if err := lock.Lock(); err != nil {
		return "", err
	}
	defer lock.Unlock()

	galleryNo, err := GalleryNoOfDir(rail, cmd.DirFileKey, db)
	if err != nil {
		return "", err
	}

	if galleryNo == "" {
		galleryNo = randutil.GenNoL("GAL", 25)
		rail.Infof("Creating gallery (%s) for directory %s (%s)", galleryNo, cmd.DirName, cmd.DirFileKey)

		err := db.Transaction(func(tx *gorm.DB) error {
			gallery := &Gallery{
				GalleryNo:  galleryNo,
				Name:       cmd.DirName,
				DirFileKey: cmd.DirFileKey,
				UserNo:     cmd.UserNo,
				CreateBy:   cmd.Username,
				UpdateBy:   cmd.Username,
				IsDel:      false,
			}
			return dbquery.NewQuery(rail, tx).Table("gallery").Omit("CreateTime", "UpdateTime").CreateAny(gallery)
		})
		if err != nil {
			return galleryNo, err
		}
	}
	return galleryNo, nil
}

// Create a new Gallery
func CreateGallery(rail miso.Rail, cmd CreateGalleryCmd, user common.User, tx *gorm.DB) (*Gallery, error) {
	rail.Infof("Creating gallery, cmd: %#v, user: %#v", cmd, user)

	gal, er := redis.RLockRun(rail, "fantahsea:gallery:create:"+user.UserNo, func() (*Gallery, error) {

		if isUsed, err := IsGalleryNameUsed(rail, cmd.Name, user.UserNo, tx); isUsed || err != nil {
			if err != nil {
				return nil, err
			}
			return nil, errs.NewErrf("You already have a gallery with the same name, please change and try again")
		}

		galleryNo := randutil.GenNoL("GAL", 25)
		gallery := &Gallery{
			GalleryNo: galleryNo,
			Name:      cmd.Name,
			UserNo:    user.UserNo,
			CreateBy:  user.Username,
			UpdateBy:  user.Username,
			IsDel:     false,
		}
		err := dbquery.NewQuery(rail, tx).Omit("CreateTime", "UpdateTime").CreateAny(gallery)
		return gallery, err
	})

	if er != nil {
		return nil, er
	}

	return gal, nil
}

/* Update a Gallery */
func UpdateGallery(rail miso.Rail, cmd UpdateGalleryCmd, user common.User, tx *gorm.DB) error {
	galleryNo := cmd.GalleryNo

	gallery, e := FindGallery(rail, tx, galleryNo)
	if e != nil {
		return e
	}

	// only owner can update the gallery
	if user.UserNo != gallery.UserNo {
		return errs.NewErrf("You are not allowed to update this gallery")
	}

	err := dbquery.NewQuery(rail, tx).
		Table("gallery").
		Where("gallery_no = ?", galleryNo).
		SetCols(Gallery{
			Name:     cmd.Name,
			UpdateBy: user.Username,
		}).
		UpdateAny()

	if err != nil {
		rail.Warnf("Failed to update gallery, gallery_no: %v, e: %v", galleryNo, err)
		return errs.NewErrf("Failed to update gallery, please try again later")
	}

	return nil
}

/* Find Gallery's creator by gallery_no */
func FindGalleryCreator(rail miso.Rail, galleryNo string, tx *gorm.DB) (*string, error) {
	var gallery Gallery
	n, err := dbquery.NewQuery(rail, tx).Raw(`
		SELECT g.user_no from gallery g
		WHERE g.gallery_no = ?
		AND g.is_del = 0`, galleryNo).Scan(&gallery)

	if err != nil || n < 1 {
		if err != nil {
			rail.Warnf("failed to find gallery %v, %v", galleryNo, err)
			return nil, err
		}
		rail.Warnf("Could not find gallery %v", galleryNo)
		return nil, errs.NewErrf("Gallery doesn't exist")
	}
	return &gallery.UserNo, nil
}

/* Find Gallery by gallery_no */
func FindGallery(rail miso.Rail, tx *gorm.DB, galleryNo string) (*Gallery, error) {
	var gallery Gallery
	n, err := dbquery.NewQuery(rail, tx).
		Raw(`SELECT g.* from gallery g WHERE g.gallery_no = ? AND g.is_del = 0`, galleryNo).
		Scan(&gallery)

	if err != nil || n < 1 {
		if err != nil {
			return nil, fmt.Errorf("failed to find gallery, %v", err)
		}
		return nil, errs.NewErrf("Gallery doesn't exist")
	}
	return &gallery, nil
}

/* Delete a gallery */
func DeleteGallery(rail miso.Rail, tx *gorm.DB, cmd DeleteGalleryCmd, user common.User) error {
	galleryNo := cmd.GalleryNo
	if access, err := HasAccessToGallery(rail, tx, user.UserNo, galleryNo); !access || err != nil {
		if err != nil {
			return err
		}
		return errs.NewErrf("You are not allowed to delete this gallery")
	}

	return dbquery.ExecSQL(rail, tx, `UPDATE gallery g SET g.is_del = 1 WHERE gallery_no = ? AND g.is_del = 0`, galleryNo)
}

// Check if the gallery exists
func GalleryExists(rail miso.Rail, tx *gorm.DB, galleryNo string) (bool, error) {
	var gallery Gallery

	n, err := dbquery.NewQuery(rail, tx).Raw(`SELECT g.id from gallery g WHERE g.gallery_no = ? AND g.is_del = 0`, galleryNo).
		Scan(&gallery)

	if err != nil || n < 1 {
		return false, err
	}

	return true, nil
}

func OnCreateGalleryImgEvent(rail miso.Rail, evt CreateGalleryImgEvent) error {
	rail.Infof("Received CreateGalleryImgEvent %+v", evt)
	tx := mysql.GetMySQL()

	// it's meant to be used for adding image to the gallery that belongs to the directory
	if evt.DirFileKey == "" {
		return nil
	}

	// create gallery for the directory if necessary
	galleryNo, err := CreateGalleryForDir(rail, CreateGalleryForDirCmd{
		Username:   evt.Username,
		UserNo:     evt.UserNo,
		DirName:    evt.DirName,
		DirFileKey: evt.DirFileKey,
	}, tx)

	if err != nil {
		return err
	}

	// add image to the gallery
	return CreateGalleryImage(rail,
		CreateGalleryImageCmd{
			GalleryNo: galleryNo,
			Name:      evt.ImageName,
			FileKey:   evt.ImageFileKey,
		},
		evt.UserNo,
		evt.Username, tx)
}

func OnNotifyFileDeletedEvent(rail miso.Rail, evt NotifyFileDeletedEvent) error {
	rail.Infof("Received NotifyFileDeletedEvent: %+v", evt)
	ok, err := DeleteGalleryImage(rail, mysql.GetMySQL(), evt.FileKey)
	if err != nil {
		return err
	}
	if !ok {
		if err := DeleteDirGallery(rail, mysql.GetMySQL(), evt.FileKey); err != nil {
			return err
		}
	}
	return RemoveDeletedFileFromAllVFolder(rail, mysql.GetMySQL(), evt.FileKey)
}
