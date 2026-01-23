package binlog

import (
	"fmt"

	"github.com/curtisnewbie/event-pump/binlog"
	"github.com/curtisnewbie/event-pump/client"
	"github.com/curtisnewbie/miso/middleware/mysql"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/user-vault/api"
	"github.com/curtisnewbie/user-vault/internal/postbox"
	"github.com/curtisnewbie/user-vault/internal/vault"

	pump "github.com/curtisnewbie/event-pump/client"
)

const (
	BinlogStreamUserCreated       = "user-vault:binlog:user-created"
	BinlogStreamReloadAccessCache = "user-vault:binlog:reload-access-cache"
)

func SubscribeBinlogEvents(rail miso.Rail) error {
	binlog.SubscribeBinlogEventsOnBootstrapV2(binlog.SubscribeBinlogOption{
		Pipeline: client.Pipeline{
			Schema:     miso.GetPropStr(mysql.PropMySQLSchema),
			Table:      "notification",
			EventTypes: []client.EventType{client.EventTypeInsert, client.EventTypeUpdate},
			Stream:     "event.bus.postbox.notification.count.changed",
		},
		Concurrency:        2,
		ContinueOnErr:      true,
		Listener:           postbox.EvictNotifCountCache,
		ListenerLogPayload: false,
	})

	binlog.SubscribeBinlogEventsOnBootstrapV2(
		binlog.SubscribeBinlogOption{
			ContinueOnErr: true,
			Pipeline: pump.Pipeline{
				Schema:     miso.GetPropStr(mysql.PropMySQLSchema),
				Table:      "user",
				EventTypes: []pump.EventType{pump.EventTypeInsert},
				Stream:     BinlogStreamUserCreated,
			},
			Concurrency: 2,
			Listener: func(rail miso.Rail, t pump.StreamEvent) error {
				username, ok := t.ColumnAfter("username")
				if !ok || username == "" {
					return nil
				}

				err := api.CreateNotifiByAccessPipeline.Send(rail, api.CreateNotifiByAccessEvent{
					Title:   fmt.Sprintf("Review user %v's registration", username),
					Message: fmt.Sprintf("Please review new user %v's registration. A role should be assigned for the new user.", username),
					ResCode: "manage-users",
				})
				if err != nil {
					rail.Errorf("failed to create notification for UserRegister, %v", err)
				}
				return nil
			},
			ListenerLogPayload: true,
		},
	)

	binlog.SubscribeBinlogEventsOnBootstrapV3(
		binlog.SubscribeBinlogOptionV3{
			ContinueOnErr: true,
			MergedPipeline: pump.MergedPipeline{
				Stream: BinlogStreamReloadAccessCache,
				Pipelines: []pump.MPipeline{
					{
						Table: "path",
					},
					{
						Table: "role_resource",
					},
					{
						Table: "path_resource",
					},
				},
			},
			Concurrency: 1,
			Listener: func(rail miso.Rail, t pump.StreamEvent) error {
				if err := vault.BatchLoadRoleAccessCache(rail); err != nil {
					rail.Errorf("Failed to BatchLoadRoleAccessCache, %v", err)
				}
				if err := vault.LoadPublicAccessCache(rail); err != nil {
					rail.Errorf("Failed to LoadPublicAccessCache, %v", err)
				}
				return nil
			},
			ListenerLogPayload: true,
		},
	)

	return nil
}
