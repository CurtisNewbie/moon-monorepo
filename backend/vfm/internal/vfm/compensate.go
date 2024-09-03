package vfm

import (
	"fmt"
	"time"

	fstore "github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/miso/miso"
	"gorm.io/gorm"
)

func CompensateThumbnail(rail miso.Rail, tx *gorm.DB) error {
	rail.Info("CompensateThumbnail start")
	defer miso.TimeOp(rail, time.Now(), "CompensateThumbnail")

	type FileProcInf struct {
		Id           int
		Name         string
		Uuid         string
		FstoreFileId string
	}

	limit := 500
	minId := 0

	for {
		var files []FileProcInf
		t := tx.
			Raw(`SELECT id, name, uuid, fstore_file_id
			FROM file_info
			WHERE id > ?
			AND file_type = 'file'
			AND is_logic_deleted = 0
			AND thumbnail = ''
			ORDER BY id ASC
			LIMIT ?`, minId, limit).
			Scan(&files)
		if t.Error != nil {
			return t.Error
		}
		if t.RowsAffected < 1 || len(files) < 1 {
			return nil // the end
		}

		for _, f := range files {
			if isImage(f.Name) {
				event := fstore.ImgThumbnailTriggerEvent{Identifier: f.Uuid, FileId: f.FstoreFileId, ReplyTo: CompressImgNotifyEventBus}
				if e := fstore.GenImgThumbnailPipeline.Send(rail, event); e != nil {
					rail.Errorf("Failed to send CompressImageEvent, minId: %v, uuid: %v, %v", minId, f.Uuid, e)
					return e
				}
				continue
			}

			if isVideo(f.Name) {
				evt := fstore.VidThumbnailTriggerEvent{
					Identifier: f.Uuid,
					FileId:     f.FstoreFileId,
					ReplyTo:    GenVideoThumbnailNotifyEventBus,
				}
				if e := fstore.GenVidThumbnailPipeline.Send(rail, evt); e != nil {
					return fmt.Errorf("failed to send %#v, uuid: %v, %v", evt, f.Uuid, e)
				}
				continue
			}
		}

		minId = files[len(files)-1].Id
		rail.Infof("CompensateThumbnail, minId: %v", minId)
	}
}
