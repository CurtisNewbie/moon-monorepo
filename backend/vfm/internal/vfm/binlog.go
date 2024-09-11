package vfm

import (
	"github.com/curtisnewbie/event-pump/client"
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/miso"
)

func SubscribeBinlogChanges(rail miso.Rail) error {
	pipelines := []client.Pipeline{
		{
			Schema:     miso.GetPropStr(mysql.PropMySQLSchema),
			Table:      "file_info",
			EventTypes: []client.EventType{client.EventTypeInsert},
			Stream:     FileSavedEventBus,
		},
		{
			Schema:     miso.GetPropStr(mysql.PropMySQLSchema),
			Table:      "file_info",
			EventTypes: []client.EventType{client.EventTypeUpdate},
			Stream:     ThumbnailUpdatedEventBus,
			Condition: client.Condition{
				ColumnChanged: []string{"thumbnail"},
			},
		},
		{
			Schema:     miso.GetPropStr(mysql.PropMySQLSchema),
			Table:      "file_info",
			EventTypes: []client.EventType{client.EventTypeUpdate},
			Stream:     FileLDeletedEventBus,
			Condition: client.Condition{
				ColumnChanged: []string{"is_logic_deleted"},
			},
		},
	}
	for _, p := range pipelines {
		err := client.CreatePipeline(rail, p)
		if err != nil {
			return err
		}
	}
	return nil
}
