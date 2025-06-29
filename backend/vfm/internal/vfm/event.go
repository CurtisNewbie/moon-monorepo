package vfm

import (
	"fmt"

	ep "github.com/curtisnewbie/event-pump/client"
	fstore "github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/middleware/rabbit"
	"github.com/curtisnewbie/miso/miso"
	"gorm.io/gorm"
)

const (
	AddFileToVFolderEventBus        = "event.bus.vfm.file.vfolder.add"
	CompressImgNotifyEventBus       = "vfm.image.compressed.event"
	GenVideoThumbnailNotifyEventBus = "vfm.video.thumbnail.generate"
	UnzipResultNotifyEventBus       = "vfm.unzip.result.notify.event"
)

var (
	UnzipResultNotifyPipeline       = rabbit.NewEventPipeline[fstore.UnzipFileReplyEvent](UnzipResultNotifyEventBus)
	GenVideoThumbnailNotifyPipeline = rabbit.NewEventPipeline[fstore.GenVideoThumbnailReplyEvent](GenVideoThumbnailNotifyEventBus)
	CompressImgNotifyPipeline       = rabbit.NewEventPipeline[fstore.ImageCompressReplyEvent](CompressImgNotifyEventBus)
	AddFileToVFolderPipeline        = rabbit.NewEventPipeline[AddFileToVfolderEvent](AddFileToVFolderEventBus)

	// TODO: deprecated, remove this in v0.0.3
	CalcDirSizePipeline = rabbit.NewEventPipeline[CalcDirSizeEvt]("event.bus.vfm.dir.size.calc").Listen(1, func(rail miso.Rail, t CalcDirSizeEvt) error {
		return nil
	})
)

func PrepareEventBus(rail miso.Rail) error {
	UnzipResultNotifyPipeline.Listen(2, OnUnzipFileReplyEvent)
	GenVideoThumbnailNotifyPipeline.Listen(2, OnVidoeThumbnailGenerated)
	CompressImgNotifyPipeline.Listen(2, OnImageCompressed)
	AddFileToVFolderPipeline.Listen(2, OnAddFileToVfolderEvent)
	return nil
}

type NotifyFileDeletedEvent struct {
	FileKey string `json:"fileKey"`
}

// event-pump send binlog event when a file_info record is saved.
// vfm guesses if the file is an image by file name,
// if so, vfm sends events to hammer to compress the image as a thumbnail
func OnFileSaved(rail miso.Rail, evt ep.StreamEvent) error {
	if evt.Type != ep.EventTypeInsert {
		return nil
	}

	var uuid string
	uuidCol, ok := evt.Columns["uuid"]
	if !ok {
		rail.Errorf("Event doesn't contain uuid, %+v", evt)
		return nil
	}
	uuid = uuidCol.After

	// lock before we do anything about it
	lock := fileLock(rail, uuid)
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	f, err := findFile(rail, mysql.GetMySQL(), uuid)
	if err != nil {
		return err
	}
	if f == nil {
		rail.Infof("file is deleted, %v", uuid)
		return nil // file already deleted
	}

	// reload user's dir tree cache
	if err := userDirTreeCache.Del(rail, f.UploaderNo); err != nil {
		rail.Errorf("Failed to reload user %v file directory tree cache, %v", f.UploaderNo, err)
	} else {
		rail.Infof("Reloaded user %v file directory tree cache", f.UploaderNo)
	}

	if f.FileType != FileTypeFile {
		rail.Infof("file is dir, %v", uuid)
		return nil // a directory
	}

	if f.Thumbnail != "" {
		rail.Infof("file has thumbnail aleady, %v", uuid)
		return nil // already has a thumbnail
	}

	if isImage(f.Name) {
		evt := fstore.ImgThumbnailTriggerEvent{Identifier: f.Uuid, FileId: f.FstoreFileId, ReplyTo: CompressImgNotifyEventBus}
		if err := fstore.GenImgThumbnailPipeline.Send(rail, evt); err != nil {
			return fmt.Errorf("failed to send %#v, uuid: %v, %v", evt, f.Uuid, err)
		}
		return nil
	}

	if isVideo(f.Name) {
		evt := fstore.VidThumbnailTriggerEvent{
			Identifier: f.Uuid,
			FileId:     f.FstoreFileId,
			ReplyTo:    GenVideoThumbnailNotifyEventBus,
		}
		if err := fstore.GenVidThumbnailPipeline.Send(rail, evt); err != nil {
			return fmt.Errorf("failed to send %#v, uuid: %v, %v", evt, f.Uuid, err)
		}
		return nil
	}

	return nil
}

// hammer sends event message when the thumbnail image is compressed and saved on mini-fstore
func OnImageCompressed(rail miso.Rail, evt fstore.ImageCompressReplyEvent) error {
	rail.Infof("Receive %#v", evt)
	return OnThumbnailGenerated(rail, mysql.GetMySQL(), evt.Identifier, evt.FileId)
}

func OnVidoeThumbnailGenerated(rail miso.Rail, evt fstore.GenVideoThumbnailReplyEvent) error {
	rail.Infof("Receive %#v", evt)
	return OnThumbnailGenerated(rail, mysql.GetMySQL(), evt.Identifier, evt.FileId)
}

func OnThumbnailGenerated(rail miso.Rail, tx *gorm.DB, identifier string, fileId string) error {
	fileKey := identifier
	lock := fileLock(rail, fileKey)
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	f, e := findFile(rail, tx, fileKey)
	if e != nil {
		rail.Errorf("Unable to find file, uuid: %v, %v", fileKey, e)
		return nil
	}
	if f == nil {
		rail.Errorf("File not found, uuid: %v", fileKey)
		return nil
	}

	if f.Thumbnail == fileId {
		rail.Infof("Thumbnail for %v remain the same, %v", fileKey, fileId)
		return nil
	} else if f.Thumbnail != "" {
		if err := fstore.DeleteFile(rail, f.Thumbnail); err != nil {
			rail.Errorf("Delete previous thumbnail failed, %v, %v", f.Thumbnail, err)
			return err
		}
		rail.Infof("Deleted previous thumbnail for %v, %v", fileKey, fileId)
	}

	err := tx.Exec("UPDATE file_info SET thumbnail = ? WHERE uuid = ?", fileId, fileKey).
		Error
	if err == nil {
		rail.Infof("Updated file's thumbnail to %v, fileKey: %v", fileId, fileKey)
	}
	return err
}

// event-pump send binlog event when a file_info's thumbnail is updated.
// vfm receives the event and check if the file has a thumbnail,
// if so, sends events to fantahsea to create a gallery image,
// adding current image to the gallery for its directory
func OnThumbnailUpdated(rail miso.Rail, evt ep.StreamEvent) error {
	if evt.Type != ep.EventTypeUpdate {
		return nil
	}

	uuid, ok := evt.ColumnAfter("uuid")
	if !ok {
		rail.Errorf("Event doesn't contain uuid column, %+v", evt)
		return nil
	}

	thumbnail, ok := evt.ColumnAfter("thumbnail")
	if !ok || thumbnail == "" {
		return nil
	}

	// lock before we do anything about it
	lock := fileLock(rail, uuid)
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	f, err := findFile(rail, mysql.GetMySQL(), uuid)
	if err != nil {
		return err
	}
	if f == nil || f.FileType != FileTypeFile {
		return nil
	}

	if f.Thumbnail == "" || f.ParentFile == "" {
		return nil
	}
	if !isImage(f.Name) {
		return nil
	}

	pf, err := findFile(rail, mysql.GetMySQL(), f.ParentFile)
	if err != nil {
		return err
	}
	if pf == nil {
		rail.Infof("parent file not found, %v", f.ParentFile)
		return nil
	}

	user, err := CachedFindUser(rail, f.UploaderNo)
	if err != nil {
		return err
	}

	cfi := CreateGalleryImgEvent{
		Username:     user.Username,
		UserNo:       user.UserNo,
		DirFileKey:   pf.Uuid,
		DirName:      pf.Name,
		ImageName:    f.Name,
		ImageFileKey: f.Uuid,
	}
	return OnCreateGalleryImgEvent(rail, cfi)
}

// event-pump send binlog event when a file_info is deleted (is_logic_deleted changed)
// vfm notifies fantahsea about the delete
func OnFileDeleted(rail miso.Rail, evt ep.StreamEvent) error {
	if evt.Type != ep.EventTypeUpdate {
		return nil
	}

	uuid, ok := evt.ColumnAfter("uuid")
	if !ok {
		rail.Errorf("Event doesn't contain uuid column, %+v", evt)
		return nil
	}

	isLogicDeleted, ok := evt.ColumnAfter("is_logic_deleted")
	if !ok {
		rail.Errorf("Event doesn't contain is_logic_deleted, %+v", evt)
		return nil
	}

	if isLogicDeleted != "1" { // FILE_LDEL_Y
		return nil
	}

	rail.Infof("File logically deleted, %v", uuid)

	if e := OnNotifyFileDeletedEvent(rail, NotifyFileDeletedEvent{FileKey: uuid}); e != nil {
		return fmt.Errorf("failed to send NotifyFileDeletedEvent, uuid: %v, %v", uuid, e)
	}
	return nil
}

type AddFileToVfolderEvent struct {
	Username string
	UserNo   string
	FolderNo string
	FileKeys []string
}

func OnAddFileToVfolderEvent(rail miso.Rail, evt AddFileToVfolderEvent) error {
	return HandleAddFileToVFolderEvent(rail, mysql.GetMySQL(), evt)
}

type CalcDirSizeEvt struct {
	FileKey string
}

func OnUnzipFileReplyEvent(rail miso.Rail, evt fstore.UnzipFileReplyEvent) error {
	rail.Infof("received UnzipFileReplyEvent: %+v", evt)
	return HandleZipUnpackResult(rail, mysql.GetMySQL(), evt)
}

// file is moved to another directory.
func OnFileMoved(rail miso.Rail, evt ep.StreamEvent) error {

	fileKey, ok := evt.ColumnAfter("uuid")
	if !ok {
		rail.Errorf("Event doesn't contain uuid column, %+v", evt)
		return nil
	}

	rail.Infof("File %v moved", fileKey)

	// update dir parent cache
	if err := dirParentCache.Del(rail, fileKey); err != nil {
		rail.Errorf("Failed to update file %v parent dir cache, %v", fileKey, err)
	} else {
		rail.Infof("Updated file %v parent dir cache", fileKey)
	}

	parentFile, ok := evt.Columns["parent_file"]
	if !ok {
		rail.Errorf("Event doesn't contain parent_file column, %+v", evt)
		return nil
	}
	rail.Infof("Filed %v is moved from %v to %v", fileKey, parentFile.Before, parentFile.After)

	db := mysql.GetMySQL()
	f, err := findFile(rail, db, fileKey)
	if err != nil {
		return err
	}

	// reload user's dir tree cache
	err = userDirTreeCache.Del(rail, f.UploaderNo)
	if err != nil {
		rail.Errorf("Failed to reload user %v file directory tree cache, %v", f.UploaderNo, err)
	} else {
		rail.Infof("Reloaded user %v file directory tree cache", f.UploaderNo)
	}

	if parentFile.Before != "" {
		// remove from previous directory's gallery
		err := RemoveGalleryImage(rail, db, parentFile.Before, fileKey)
		if err != nil {
			rail.Errorf("RemoveGalleryImage failed, fileKey: %v, dirFileKey: %v", fileKey, parentFile.Before)
		} else {
			rail.Infof("Removed image from gallery, fileKey: %v, dirFileKey: %v", fileKey, parentFile.Before)
		}
	}

	if parentFile.After != "" {
		// lock before we do anything about it
		lock := fileLock(rail, fileKey)
		if err := lock.Lock(); err != nil {
			return err
		}
		defer lock.Unlock()

		if f == nil || f.FileType != FileTypeFile ||
			f.Thumbnail == "" || !isImage(f.Name) {
			return nil
		}

		pf, err := findFile(rail, db, parentFile.After)
		if err != nil {
			return err
		}
		if pf == nil {
			rail.Infof("parent file not found, %v", f.ParentFile)
			return nil
		}

		user, err := CachedFindUser(rail, f.UploaderNo)
		if err != nil {
			return err
		}

		cfi := CreateGalleryImgEvent{
			Username:     user.Username,
			UserNo:       user.UserNo,
			DirFileKey:   pf.Uuid,
			DirName:      pf.Name,
			ImageName:    f.Name,
			ImageFileKey: f.Uuid,
		}
		return OnCreateGalleryImgEvent(rail, cfi)
	}
	return nil
}

func OnDirNameUpdated(rail miso.Rail, evt ep.StreamEvent) error {

	fileKey, ok := evt.ColumnAfter("uuid")
	if !ok {
		rail.Errorf("Event doesn't contain uuid column, %+v", evt)
		return nil
	}

	db := mysql.GetMySQL()
	f, err := findFile(rail, db, fileKey)
	if err != nil {
		return err
	}

	// reload user's dir tree cache
	if err := userDirTreeCache.Del(rail, f.UploaderNo); err != nil {
		rail.Errorf("Failed to reload user %v file directory tree cache, %v", f.UploaderNo, err)
	} else {
		rail.Infof("Reloaded user %v file directory tree cache", f.UploaderNo)
	}

	// reload dir name cache
	if err := dirNameCache.Del(rail, f.Uuid); err != nil {
		rail.Errorf("Failed to reloaded directory name cache, %v, %v", f.Uuid, err)
	} else {
		rail.Infof("Reloaded directory name cache, %v", f.Uuid)
	}

	if f.FileType != FileTypeDir {
		return nil
	}

	rail.Infof("Directory name changed, updating directory's gallery name, fileKey: %v", fileKey)

	galleryNo, err := GalleryNoOfDir(fileKey, db)
	if err != nil || galleryNo == "" {
		return err
	}

	_, err = dbquery.NewQueryRail(rail, db).
		Where("gallery_no = ?", galleryNo).
		SetCols(Gallery{
			Name:     f.Name,
			UpdateBy: "vfm",
		}).
		Update()
	return err
}
