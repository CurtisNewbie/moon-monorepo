package vfm

import (
	"github.com/curtisnewbie/event-pump/binlog"
	"github.com/curtisnewbie/event-pump/client"
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/miso"
)

func SubscribeBinlogChanges(rail miso.Rail) error {

	binlog.SubscribeBinlogEventsOnBootstrapV2(binlog.SubscribeBinlogOption{
		Pipeline: client.Pipeline{
			Schema:     miso.GetPropStr(mysql.PropMySQLSchema),
			Table:      "file_info",
			EventTypes: []client.EventType{client.EventTypeInsert},
			Stream:     "event.bus.vfm.file.saved",
		},
		Concurrency:   2,
		ContinueOnErr: true,
		Listener:      OnFileSaved,
	})

	binlog.SubscribeBinlogEventsOnBootstrapV2(binlog.SubscribeBinlogOption{
		Pipeline: client.Pipeline{
			Schema:     miso.GetPropStr(mysql.PropMySQLSchema),
			Table:      "file_info",
			EventTypes: []client.EventType{client.EventTypeUpdate},
			Stream:     "event.bus.vfm.file.thumbnail.updated",
			Condition: client.Condition{
				ColumnChanged: []string{"thumbnail"},
			},
		},
		Concurrency:   2,
		ContinueOnErr: true,
		Listener:      OnThumbnailUpdated,
	})

	binlog.SubscribeBinlogEventsOnBootstrapV2(binlog.SubscribeBinlogOption{
		Pipeline: client.Pipeline{
			Schema:     miso.GetPropStr(mysql.PropMySQLSchema),
			Table:      "file_info",
			EventTypes: []client.EventType{client.EventTypeUpdate},
			Stream:     "event.bus.vfm.file.logic.deleted",
			Condition: client.Condition{
				ColumnChanged: []string{"is_logic_deleted"},
			},
		},
		Concurrency:   2,
		ContinueOnErr: true,
		Listener:      OnFileDeleted,
	})

	binlog.SubscribeBinlogEventsOnBootstrapV2(binlog.SubscribeBinlogOption{
		Pipeline: client.Pipeline{
			Schema:     miso.GetPropStr(mysql.PropMySQLSchema),
			Table:      "file_info",
			EventTypes: []client.EventType{client.EventTypeUpdate},
			Stream:     "event.bus.vfm.file.moved",
			Condition: client.Condition{
				ColumnChanged: []string{"parent_file"},
			},
		},
		Concurrency:   2,
		ContinueOnErr: true,
		Listener:      OnFileMoved,
	})

	return nil
}
