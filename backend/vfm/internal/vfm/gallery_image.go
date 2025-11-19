package vfm

import (
	"time"

	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/async"
	"github.com/curtisnewbie/miso/util/atom"
	"github.com/curtisnewbie/miso/util/errs"
	"github.com/curtisnewbie/miso/util/randutil"
	"gorm.io/gorm"
)

// GalleryImage.status (doesn't really matter anymore)
type ImgStatus string

const (
	TableGalleryImage = "gallery_image"

	NORMAL  ImgStatus = "NORMAL"
	DELETED ImgStatus = "DELETED"

	// 40mb is the maximum size for an image
	IMAGE_SIZE_THRESHOLD int64 = 40 * 1048576
)

type TransferGalleryImageReq struct {
	Images []CreateGalleryImageCmd `json:"images"`
}

type TransferGalleryImageInDirReq struct {
	// gallery no
	GalleryNo string `json:"galleryNo" validation:"notEmpty"`

	// file key of the directory
	FileKey string `json:"fileKey" validation:"notEmpty"`
}

// Image that belongs to a Gallery
type GalleryImage struct {
	ID         int64
	GalleryNo  string
	ImageNo    string
	Name       string
	FileKey    string
	Status     ImgStatus
	CreateTime time.Time
	CreateBy   string
	UpdateTime time.Time
	UpdateBy   string
}

func (GalleryImage) TableName() string {
	return "gallery_image"
}

type ThumbnailInfo struct {
	Name string
	Path string
}

type CreateGalleryImgEvent struct {
	Username     string `json:"username"`
	UserNo       string `json:"userNo"`
	DirFileKey   string `json:"dirFileKey"`
	DirName      string `json:"dirName"`
	ImageName    string `json:"imageName"`
	ImageFileKey string `json:"imageFileKey"`
}

type ListGalleryImagesCmd struct {
	GalleryNo   string `json:"galleryNo" validation:"notEmpty"`
	miso.Paging `json:"paging"`
}

type ListGalleryImagesResp struct {
	Images []ImageInfo `json:"images"`
	Paging miso.Paging `json:"paging"`
}

type ImageInfo struct {
	FileKey         string `json:"fileKey"`
	ThumbnailToken  string `json:"thumbnailToken"`
	FileTempToken   string `json:"fileTempToken"`
	ImageFileId     string `json:"-"`
	ThumbnailFileId string `json:"-"`
}

type CreateGalleryImageCmd struct {
	GalleryNo string `json:"galleryNo"`
	Name      string `json:"name"`
	FileKey   string `json:"fileKey"`
}

func DeleteGalleryImage(rail miso.Rail, tx *gorm.DB, fileKey string) (bool, error) {
	n, err := dbquery.NewQuery(rail, tx).Exec("delete from gallery_image where file_key = ?", fileKey)
	if err != nil {
		return false, errs.Wrapf(err, "failed to update gallery_image, uuid: %v", fileKey)
	}
	if n > 0 {
		rail.Infof("Removed file %v from all galleries", fileKey)
	}
	return n > 0, nil
}

func DeleteDirGallery(rail miso.Rail, db *gorm.DB, dirFileKey string) error {
	lock := NewGalleryDirLock(rail, dirFileKey)
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	var galleryNo string
	ok, err := dbquery.NewQuery(rail, db).
		Table("gallery").
		Eq("dir_file_key", dirFileKey).
		Select("gallery_no").
		ScanAny(&galleryNo)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	ok, err = dbquery.NewQuery(rail, db).
		Table("gallery_image").
		Eq("gallery_no", galleryNo).
		HasAny()
	if err != nil {
		return err
	}

	// just in case the user adds images to the dir gallery
	// if the dir gallery is empty, we just delete it
	// if not, user may reuse the gallery for whatever reason
	if !ok {
		err = dbquery.NewQuery(rail, db).
			Table("gallery").
			Eq("gallery_no", galleryNo).
			DeleteAny()
		if err == nil {
			rail.Infof("Deleted gallery for dir: %v", dirFileKey)
		}
	}
	return err
}

// Create a gallery image record
func CreateGalleryImage(rail miso.Rail, cmd CreateGalleryImageCmd, userNo string, username string, db *gorm.DB) error {
	creator, err := FindGalleryCreator(rail, cmd.GalleryNo, db)
	if err != nil {
		return err
	}

	if *creator != userNo {
		return errs.NewErrf("You are not allowed to upload image to this gallery")
	}

	lock := NewGalleryFileLock(rail, cmd.GalleryNo, cmd.FileKey)
	if err := lock.Lock(); err != nil {
		return errs.Wrapf(err, "failed to obtain gallery image lock, gallery:%v, fileKey: %v", cmd.GalleryNo, cmd.FileKey)
	}
	defer lock.Unlock()

	if isCreated, e := isImgCreatedAlready(rail, db, cmd.GalleryNo, cmd.FileKey); isCreated || e != nil {
		if e != nil {
			return e
		}
		rail.Infof("Image '%s' added already", cmd.Name)
		return nil
	}

	imageNo := randutil.GenNoL("IMG", 25)
	return db.Transaction(func(tx *gorm.DB) error {
		if err := dbquery.NewQuery(rail, tx).
			ExecAny(`insert into gallery_image (gallery_no, image_no, name, file_key, create_by) values (?, ?, ?, ?, ?)`,
				cmd.GalleryNo, imageNo, cmd.Name, cmd.FileKey, username); err != nil {
			return err
		}

		return dbquery.NewQuery(rail, tx).
			ExecAny(`UPDATE gallery SET update_time = ? WHERE gallery_no = ?`, atom.Now(), cmd.GalleryNo)
	})
}

type FstoreTmpToken struct {
	FileId  string
	TempKey string
}

func GenFstoreTknBatch(rail miso.Rail, futures *async.AwaitFutures[FstoreTmpToken], fileId string, name string) {
	futures.SubmitAsync(func() (FstoreTmpToken, error) {
		tkn, err := GetFstoreTmpToken(rail.NextSpan(), fileId, name)
		if err != nil {
			return FstoreTmpToken{FileId: fileId}, err
		}
		return FstoreTmpToken{
			FileId:  fileId,
			TempKey: tkn,
		}, nil
	})
}

func GenFstoreTknAsync(rail miso.Rail, fileId string, name string) async.Future[FstoreTmpToken] {
	return async.Submit[FstoreTmpToken](vfmPool,
		func() (FstoreTmpToken, error) {
			tkn, err := GetFstoreTmpToken(rail.NextSpan(), fileId, name)
			if err != nil {
				return FstoreTmpToken{}, err
			}
			return FstoreTmpToken{
				FileId:  fileId,
				TempKey: tkn,
			}, nil
		})
}

// List gallery images
func ListGalleryImages(rail miso.Rail, tx *gorm.DB, cmd ListGalleryImagesCmd, user common.User) (*ListGalleryImagesResp, error) {
	if hasAccess, err := HasAccessToGallery(rail, tx, user.UserNo, cmd.GalleryNo); err != nil || !hasAccess {
		if err != nil {
			return nil, errs.Wrapf(err, "check HasAccessToGallery failed")
		}
		return nil, errs.ErrNotPermitted.New()
	}

	var galleryImages []GalleryImage
	_, err := dbquery.NewQuery(rail, tx).
		Table("gallery_image").
		Select("image_no, file_key").
		Eq("gallery_no", cmd.GalleryNo).
		Order("name ASC").
		Offset(cmd.Paging.GetOffset()).
		Limit(cmd.Paging.GetLimit()).
		Scan(&galleryImages)
	if err != nil {
		return nil, errs.Wrapf(err, "select gallery_image failed")
	}
	if galleryImages == nil {
		galleryImages = []GalleryImage{}
	}

	// count total asynchronoulsy (normally, when the SELECT is successful, the COUNT doesn't really fail)
	countFuture := async.Submit(vfmPool, func() (int, error) {
		var total int
		err := dbquery.NewQuery(rail, tx).
			Raw(`SELECT COUNT(*) FROM gallery_image WHERE gallery_no = ?`, cmd.GalleryNo).
			ScanVal(&total)
		if err == nil {
			return total, nil
		}
		return total, errs.Wrapf(err, "failed to count gallery_image")
	})

	// generate temp tokens for the actual files and the thumbnail, these are served by mini-fstore
	images := []ImageInfo{}
	if len(galleryImages) > 0 {
		awaitFutures := async.NewAwaitFutures[FstoreTmpToken](vfmPool)
		for _, img := range galleryImages {
			fi, ok, e := findFile(rail, tx, img.FileKey)
			if e != nil || !ok {
				rail.Errorf("findFile failed, fileKey: %v, %v", img.FileKey, e)
				continue
			}

			// original
			GenFstoreTknBatch(rail, awaitFutures, fi.FstoreFileId, fi.Name)

			// thumbnail
			thumbnailFileId := fi.Thumbnail
			if thumbnailFileId == "" {
				thumbnailFileId = fi.FstoreFileId
			} else {
				GenFstoreTknBatch(rail, awaitFutures, thumbnailFileId, fi.Name)
			}
			images = append(images, ImageInfo{ImageFileId: fi.FstoreFileId, ThumbnailFileId: thumbnailFileId, FileKey: fi.Uuid})
		}

		genTknFutures := awaitFutures.Await()
		tokens := make([]FstoreTmpToken, 0, len(genTknFutures))
		for i := range genTknFutures {
			res, err := genTknFutures[i].Get()
			if err != nil {
				rail.Errorf("Failed to get mini-fstore temp token for fstore_file_id: %v, %v", res.FileId, err)
				continue
			}
			tokens = append(tokens, res)
		}

		idTknMap := map[string]string{}
		for _, t := range tokens {
			idTknMap[t.FileId] = t.TempKey
		}
		for i, im := range images {
			im.ThumbnailToken = idTknMap[im.ThumbnailFileId]
			im.FileTempToken = idTknMap[im.ImageFileId]
			images[i] = im
		}
	}

	total, errCnt := countFuture.Get()
	if errCnt != nil {
		return nil, errCnt
	}

	return &ListGalleryImagesResp{Images: images, Paging: miso.RespPage(cmd.Paging, total)}, nil
}

func BatchTransferAsync(rail miso.Rail, cmd TransferGalleryImageReq, user common.User, tx *gorm.DB) error {
	if len(cmd.Images) < 1 {
		return nil
	}

	// validate the keys first
	for _, img := range cmd.Images {
		if isValid, e := ValidateFileOwner(rail, tx, ValidateFileOwnerReq{
			FileKey: img.FileKey,
			UserNo:  user.UserNo,
		}); e != nil || !isValid {
			if e != nil {
				return e
			}
			return errs.NewErrf("Only file's owner can make it a gallery image ('%s')", img.Name)
		}
	}

	// start transferring
	go func(rail miso.Rail, images []CreateGalleryImageCmd) {
		for _, cmd := range images {
			fi, ok, er := findFile(rail, tx, cmd.FileKey)
			if er != nil || !ok {
				rail.Errorf("Failed to fetch file info while transferring selected images, fi's fileKey: %s, error: %v", cmd.FileKey, er)
				continue
			}

			if fi.FileType == FileTypeFile { // a file
				if fi.FstoreFileId == "" {
					continue // doesn't have fstore fileId, cannot be transferred
				}

				if GuessIsImage(rail, fi) {
					nc := CreateGalleryImageCmd{GalleryNo: cmd.GalleryNo, Name: fi.Name, FileKey: fi.Uuid}
					if err := CreateGalleryImage(rail, nc, user.UserNo, user.Username, tx); err != nil {
						rail.Errorf("Failed to create gallery image, fi's fileKey: %s, error: %v", cmd.FileKey, err)
						continue
					}
				}
			} else { // a directory
				treq := TransferGalleryImageInDirReq{
					GalleryNo: cmd.GalleryNo,
					FileKey:   cmd.FileKey,
				}
				if err := TransferImagesInDir(rail, treq, user, tx); err != nil {
					rail.Errorf("Failed to transfer images in directory, fi's fileKey: %s, error: %v", cmd.FileKey, err)
					continue
				}
			}
		}
	}(rail.NewCtx(), cmd.Images)

	return nil
}

// Transfer images in dir
func TransferImagesInDir(rail miso.Rail, cmd TransferGalleryImageInDirReq, user common.User, tx *gorm.DB) error {
	fi, ok, e := findFile(rail, tx, cmd.FileKey)
	if e != nil {
		return e
	}
	if !ok {
		return ErrFileNotFound.New()
	}

	// only the owner of the directory can do this, by default directory is only visible to the uploader
	if fi.UploaderNo != user.UserNo {
		return errs.ErrNotPermitted.New()
	}

	if fi.FileType != FileTypeDir {
		return errs.NewErrf("This is not a directory")
	}

	if fi.IsLogicDeleted == DelY || fi.IsPhysicDeleted == DelY {
		return errs.NewErrf("Directory is already deleted")
	}
	dirFileKey := cmd.FileKey
	galleryNo := cmd.GalleryNo
	start := time.Now()

	page := 1
	for {
		// dirFileKey, 100, page
		res, err := ListFilesInDir(rail, tx, ListFilesInDirReq{
			FileKey: dirFileKey,
			Limit:   100,
			Page:    page,
		})
		if err != nil {
			rail.Errorf("Failed to list files in dir, dir's fileKey: %s, error: %v", dirFileKey, err)
			break
		}
		if len(res) < 1 {
			break
		}

		// starts fetching file one by one
		for i := 0; i < len(res); i++ {
			fk := res[i]
			fi, ok, er := findFile(rail, tx, fk)
			if er != nil || !ok {
				rail.Errorf("Failed to fetch file info while looping files in dir, fi's fileKey: %s, error: %v", fk, er)
				continue
			}

			if GuessIsImage(rail, fi) {
				cmd := CreateGalleryImageCmd{GalleryNo: galleryNo, Name: fi.Name, FileKey: fi.Uuid}
				if err := CreateGalleryImage(rail, cmd, user.UserNo, user.Username, tx); err != nil {
					rail.Errorf("Failed to create gallery image, fi's fileKey: %s, error: %v", fk, err)
				}
			}
		}

		page += 1
	}

	rail.Infof("Finished TransferImagesInDir, dir's fileKey: %s, took: %s", dirFileKey, time.Since(start))
	return nil
}

// Guess whether a file is an image
func GuessIsImage(rail miso.Rail, f FileInfo) bool {
	if f.SizeInBytes > IMAGE_SIZE_THRESHOLD {
		return false
	}
	if f.FileType != FileTypeFile {
		return false
	}
	if f.Thumbnail == "" {
		rail.Infof("File doesn't have thumbnail, fileKey: %v", f.Uuid)
		return false
	}
	return isImage(f.Name)
}

// check whether the gallery image is created already
func isImgCreatedAlready(rail miso.Rail, tx *gorm.DB, galleryNo string, fileKey string) (created bool, er error) {
	ok, err := dbquery.NewQuery(rail, tx).
		Table("gallery_image").
		Eq("gallery_no", galleryNo).
		Eq("file_key", fileKey).
		HasAny()
	return ok, err
}

func NewGalleryFileLock(rail miso.Rail, galleryNo string, fileKey string) *redis.RLock {
	return redis.NewRLockf(rail, "gallery:image:%v:%v", galleryNo, fileKey)
}

func RemoveGalleryImage(rail miso.Rail, db *gorm.DB, dirFileKey string, imageFileKey string) error {
	galleryNo, err := GalleryNoOfDir(rail, dirFileKey, db)
	if err != nil {
		return err
	}
	rail.Infof("Found gallery_no of dir %v: %v", dirFileKey, galleryNo)

	lock := NewGalleryFileLock(rail, galleryNo, imageFileKey)
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	return dbquery.NewQuery(rail, db).
		Table(TableGalleryImage).
		Eq("gallery_no", galleryNo).
		Eq("file_key", imageFileKey).
		DeleteAny()
}
