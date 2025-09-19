package vfm

import (
	"context"
	"errors"
	"fmt"
	"time"

	fstore "github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/redis"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/errs"
	"gorm.io/gorm"
)

var (
	serverMaintainanceKey = "vfm:maintenance"

	serverMaintainanceTicker *miso.TickRunner = nil
)

type FileProcInf struct {
	Id           int
	Name         string
	Uuid         string
	FstoreFileId string
}

func CompensateThumbnail(rail miso.Rail, db *gorm.DB) error {
	ok, err := EnterMaintenance(rail)
	if err != nil {
		return err
	}
	if !ok {
		return errs.NewErrf("Server is already in maintenance")
	}
	defer LeaveMaintenance(rail)

	rail.Info("CompensateThumbnail start")
	defer miso.TimeOp(rail, time.Now(), "CompensateThumbnail")

	limit := 500
	minId := 0

	for {
		var files []FileProcInf
		n, err := dbquery.NewQueryRail(rail, db).
			Raw(`SELECT id, name, uuid, fstore_file_id
			FROM file_info
			WHERE id > ?
			AND file_type = 'FILE'
			AND is_logic_deleted = 0
			AND thumbnail = ''
			ORDER BY id ASC
			LIMIT ?`, minId, limit).
			Scan(&files)
		if err != nil {
			return err
		}
		if n < 1 {
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

		if len(files) < limit {
			return nil // the end
		}
	}
}

func RegenerateVideoThumbnails(rail miso.Rail, db *gorm.DB) error {
	ok, err := EnterMaintenance(rail)
	if err != nil {
		return err
	}
	if !ok {
		return errs.NewErrf("Server is already in maintenance")
	}
	defer LeaveMaintenance(rail)

	defer miso.TimeOp(rail, time.Now(), "RegenerateVideoThumbnails")

	limit := 500
	minId := 0
	var maxId int
	scanned, err := dbquery.NewQueryRail(rail, db).
		Table("file_info").
		Select("max(id)").
		Eq("file_type", "FILE").
		Eq("is_logic_deleted", 0).Scan(&maxId)

	if err != nil {
		rail.Errorf("find max id failed, %v", err)
		return err
	}
	if scanned < 1 {
		maxId = -1
	}
	rail.Infof("file_info max id: %v", maxId)

	for {
		var files []FileProcInf
		n, err := dbquery.NewQueryRail(rail, db).
			Raw(`SELECT id, name, uuid, fstore_file_id
			FROM file_info
			WHERE id > ?
			AND file_type = 'file'
			AND is_logic_deleted = 0
			AND id <= ?
			ORDER BY id ASC
			LIMIT ?`, minId, maxId, limit).
			Scan(&files)
		if err != nil {
			return err
		}
		if n < 1 || len(files) < 1 {
			return nil // the end
		}

		for _, f := range files {
			if isVideo(f.Name) {
				evt := fstore.VidThumbnailTriggerEvent{
					Identifier: f.Uuid,
					FileId:     f.FstoreFileId,
					ReplyTo:    GenVideoThumbnailNotifyEventBus,
				}
				if e := fstore.GenVidThumbnailPipeline.Send(rail, evt); e != nil {
					return fmt.Errorf("failed to send %#v, uuid: %v, %v", evt, f.Uuid, e)
				}
				rail.Infof("GenVidThumbnailPipeline sent, %#v", evt)
			}
		}

		minId = files[len(files)-1].Id
		rail.Infof("RegenerateVideoThumbnails, minId: %v", minId)

		if len(files) < limit {
			return nil // the end
		}
	}
}

func LeaveMaintenance(rail miso.Rail) error {
	serverMaintainanceTicker.Stop()
	c := redis.GetRedis().Del(rail.Context(), serverMaintainanceKey)
	if c.Err() != nil {
		if redis.IsNil(c.Err()) {
			return nil
		}
		rail.Errorf("Failed to delete redis server maintainance flag, %v", c.Err())
		return c.Err()
	}
	return nil
}

func EnterMaintenance(rail miso.Rail) (bool, error) {
	c := redis.GetRedis().SetNX(rail.Context(), serverMaintainanceKey, 1, time.Second*30)
	if c.Err() != nil {
		return false, c.Err()
	}
	if !c.Val() {
		return false, nil
	}

	serverMaintainanceTicker = miso.NewTickRuner(time.Second*5, func() {
		rail := rail.NextSpan()
		c := redis.GetRedis().SetXX(rail.Context(), serverMaintainanceKey, 1, time.Second*30)
		if c.Err() != nil {
			if !errors.Is(c.Err(), redis.Nil) {
				rail.Errorf("failed to maintain redis server maintenance flag, %v", c.Err())
			}
			return
		}
		rail.Info("Refreshed redis server maintenance flag")
	})
	serverMaintainanceTicker.Start()
	return true, nil
}

type MaintenanceStatus struct {
	UnderMaintenance bool
}

func CheckMaintenanceStatus() (MaintenanceStatus, error) {
	cmd := redis.GetRedis().Exists(context.Background(), serverMaintainanceKey)
	if cmd.Err() != nil {
		return MaintenanceStatus{}, cmd.Err()
	}
	return MaintenanceStatus{
		UnderMaintenance: cmd.Val() > 0,
	}, nil
}
