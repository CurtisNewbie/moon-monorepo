package postbox

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
			Table:      "notification",
			EventTypes: []client.EventType{client.EventTypeInsert, client.EventTypeUpdate},
			Stream:     "event.bus.postbox.notification.count.changed",
		},
		Concurrency:        2,
		ContinueOnErr:      true,
		Listener:           evictNotifCountCache,
		ListenerLogPayload: true,
	})
	return nil
}
