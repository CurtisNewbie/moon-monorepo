package vault

import (
	"fmt"

	binlog "github.com/curtisnewbie/event-pump/binlog"
	pump "github.com/curtisnewbie/event-pump/client"
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/user-vault/api"
)

const (
	BinlogStreamUserCreated = "user-vault:binlog:user-created"
)

func SubscribeBinlogEvent() {
	pipeline := pump.Pipeline{
		Schema:     miso.GetPropStr(mysql.PropMySQLSchema),
		Table:      "user",
		EventTypes: []pump.EventType{pump.EventTypeInsert},
		Stream:     BinlogStreamUserCreated,
	}
	concurrency := 2
	binlog.SubscribeBinlogEventsOnBootstrap(pipeline, concurrency, func(rail miso.Rail, t pump.StreamEvent) error {
		username, ok := t.ColumnAfter("username")
		if !ok || username == "" {
			return nil
		}

		err := api.CreateNotifiByAccessPipeline.Send(rail, api.CreateNotifiByAccessEvent{
			Title:   fmt.Sprintf("Review user %v's registration", username),
			Message: fmt.Sprintf("Please review new user %v's registration. A role should be assigned for the new user.", username),
			ResCode: ResourceManagerUser,
		})
		if err != nil {
			rail.Errorf("failed to create notification for UserRegister, %v", err)
		}
		return nil
	})
}
