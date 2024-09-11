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

// ------------------------------- entity start

// Gallery
type Gallery struct {
	Id         int64
	GalleryNo  string
	UserNo     string
	Name       string
	DirFileKey string
	CreateTime util.ETime
	CreateBy   string
	UpdateTime util.ETime
	UpdateBy   string
	IsDel      bool
}

func (Gallery) TableName() string {
	return "gallery"

}

// ------------------------------- entity end

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
	ID            int64      `json:"id"`
	GalleryNo     string     `json:"galleryNo"`
	UserNo        string     `json:"userNo"`
	Name          string     `json:"name"`
	CreateTime    util.ETime `json:"-"`
	UpdateTime    util.ETime `json:"-"`
	CreateBy      string     `json:"createBy"`
	UpdateBy      string     `json:"updateBy"`
	IsOwner       bool       `json:"isOwner"`
	CreateTimeStr string     `json:"createTime"`
	UpdateTimeStr string     `json:"updateTime"`
}

// List owned gallery briefs
func ListOwnedGalleryBriefs(rail miso.Rail, user common.User, tx *gorm.DB) ([]VGalleryBrief, error) {
	var briefs []VGalleryBrief
	t := tx.Raw(`select gallery_no, name from gallery where user_no = ? AND is_del = 0`, user.UserNo).
		Scan(&briefs)

	if e := t.Error; e != nil {
		return nil, e
	}
	if briefs == nil {
		briefs = []VGalleryBrief{}
	}

	return briefs, nil
}

/* List Galleries */
func ListGalleries(rail miso.Rail, cmd ListGalleriesCmd, user common.User, db *gorm.DB) (miso.PageRes[VGallery], error) {
	return mysql.NewPageQuery[VGallery]().
		WithPage(cmd.Paging).
		WithBaseQuery(func(tx *gorm.DB) *gorm.DB {
			return tx.Table("gallery g").
				Joins("LEFT JOIN gallery_user_access ga ON (g.gallery_no = ga.gallery_no)").
				Where("g.is_del = 0").
				Where("g.user_no = ? OR (ga.user_no = ? AND ga.is_del = 0)", user.UserNo, user.UserNo)
		}).
		WithSelectQuery(func(tx *gorm.DB) *gorm.DB {
			tx = tx.Select("g.*").Order("g.update_time DESC")
			return tx
		}).
		ForEach(func(g VGallery) VGallery {
			if g.UserNo == user.UserNo {
				g.IsOwner = true
			}
			g.CreateTimeStr = g.CreateTime.FormatClassic()
			g.UpdateTimeStr = g.UpdateTime.FormatClassic()
			return g
		}).
		Exec(rail, db)
}

func GalleryNoOfDir(dirFileKey string, tx *gorm.DB) (string, error) {
	var gallery Gallery
	tx = tx.Raw(`SELECT g.gallery_no from gallery g WHERE g.dir_file_key = ? and g.is_del = 0 limit 1`, dirFileKey).
		Scan(&gallery)

	if e := tx.Error; e != nil {
		return "", tx.Error
	}

	return gallery.GalleryNo, nil
}

// Check if the name is already used by current user
func IsGalleryNameUsed(name string, userNo string, tx *gorm.DB) (bool, error) {
	var gallery Gallery
	t := tx.Raw(`SELECT g.id from gallery g WHERE g.user_no = ? and g.name = ? AND g.is_del = 0`, userNo, name).
		Scan(&gallery)

	if e := t.Error; e != nil {
		return false, t.Error
	}

	return t.RowsAffected > 0, nil
}

// Create a new Gallery for dir
func CreateGalleryForDir(rail miso.Rail, cmd CreateGalleryForDirCmd, tx *gorm.DB) (string, error) {

	return redis.RLockRun(rail, "fantahsea:gallery:create:"+cmd.UserNo,
		func() (string, error) {
			galleryNo, err := GalleryNoOfDir(cmd.DirFileKey, tx)
			if err != nil {
				return "", err
			}

			if galleryNo == "" {
				galleryNo = util.GenNoL("GAL", 25)
				rail.Infof("Creating gallery (%s) for directory %s (%s)", galleryNo, cmd.DirName, cmd.DirFileKey)

				err := tx.Transaction(func(tx *gorm.DB) error {
					gallery := &Gallery{
						GalleryNo:  galleryNo,
						Name:       cmd.DirName,
						DirFileKey: cmd.DirFileKey,
						UserNo:     cmd.UserNo,
						CreateBy:   cmd.Username,
						UpdateBy:   cmd.Username,
						IsDel:      false,
					}
					result := tx.Omit("CreateTime", "UpdateTime").Create(gallery)
					return result.Error
				})
				if err != nil {
					return galleryNo, err
				}
			}
			return galleryNo, nil
		})
}

// Create a new Gallery
func CreateGallery(rail miso.Rail, cmd CreateGalleryCmd, user common.User, tx *gorm.DB) (*Gallery, error) {
	rail.Infof("Creating gallery, cmd: %#v, user: %#v", cmd, user)

	gal, er := redis.RLockRun(rail, "fantahsea:gallery:create:"+user.UserNo, func() (*Gallery, error) {

		if isUsed, err := IsGalleryNameUsed(cmd.Name, user.UserNo, tx); isUsed || err != nil {
			if err != nil {
				return nil, err
			}
			return nil, miso.NewErrf("You already have a gallery with the same name, please change and try again")
		}

		galleryNo := util.GenNoL("GAL", 25)
		gallery := &Gallery{
			GalleryNo: galleryNo,
			Name:      cmd.Name,
			UserNo:    user.UserNo,
			CreateBy:  user.Username,
			UpdateBy:  user.Username,
			IsDel:     false,
		}
		result := tx.Omit("CreateTime", "UpdateTime").Create(gallery)
		return gallery, result.Error
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
		return miso.NewErrf("You are not allowed to update this gallery")
	}

	t := tx.Where("gallery_no = ?", galleryNo).
		Updates(Gallery{
			Name:     cmd.Name,
			UpdateBy: user.Username,
		})

	if e := t.Error; e != nil {
		rail.Warnf("Failed to update gallery, gallery_no: %v, e: %v", galleryNo, t.Error)
		return miso.NewErrf("Failed to update gallery, please try again later")
	}

	return nil
}

/* Find Gallery's creator by gallery_no */
func FindGalleryCreator(rail miso.Rail, galleryNo string, tx *gorm.DB) (*string, error) {
	var gallery Gallery
	t := tx.Raw(`
		SELECT g.user_no from gallery g
		WHERE g.gallery_no = ?
		AND g.is_del = 0`, galleryNo).Scan(&gallery)

	if e := t.Error; e != nil || t.RowsAffected < 1 {
		if e != nil {
			rail.Warnf("failed to find gallery %v, %v", galleryNo, t.Error)
			return nil, t.Error
		}
		rail.Warnf("Could not find gallery %v", galleryNo)
		return nil, miso.NewErrf("Gallery doesn't exist")
	}
	return &gallery.UserNo, nil
}

/* Find Gallery by gallery_no */
func FindGallery(rail miso.Rail, tx *gorm.DB, galleryNo string) (*Gallery, error) {
	var gallery Gallery
	t := tx.Raw(`SELECT g.* from gallery g WHERE g.gallery_no = ? AND g.is_del = 0`, galleryNo).
		Scan(&gallery)

	if e := t.Error; e != nil || t.RowsAffected < 1 {
		if e != nil {
			return nil, fmt.Errorf("failed to find gallery, %v", t.Error)
		}
		return nil, miso.NewErrf("Gallery doesn't exist")
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
		return miso.NewErrf("You are not allowed to delete this gallery")
	}

	t := tx.Exec(`UPDATE gallery g SET g.is_del = 1 WHERE gallery_no = ? AND g.is_del = 0`, galleryNo)
	if e := t.Error; e != nil {
		return t.Error
	}

	return nil
}

// Check if the gallery exists
func GalleryExists(rail miso.Rail, tx *gorm.DB, galleryNo string) (bool, error) {
	var gallery Gallery
	tx = tx.Raw(`SELECT g.id from gallery g WHERE g.gallery_no = ? AND g.is_del = 0`, galleryNo).
		Scan(&gallery)

	if e := tx.Error; e != nil || tx.RowsAffected < 1 {
		if e != nil {
			return false, tx.Error
		}
		return false, nil
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
	return DeleteGalleryImage(rail, mysql.GetMySQL(), evt.FileKey)
}
