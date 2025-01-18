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
	BinlogStreamUserCreated    = "user-vault:binlog:user-created"
	BinlogStreamPathUpdated    = "user-vault:binlog:path-updated"
	BinlogStreamRoleResUpdated = "user-vault:binlog:role-resource-updated"
	BinlogStreamPathResUpdated = "user-vault:binlog:path-resource-updated"
)

func SubscribeBinlogEvent(rail miso.Rail) error {
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
					ResCode: ResourceManagerUser,
				})
				if err != nil {
					rail.Errorf("failed to create notification for UserRegister, %v", err)
				}
				return nil
			},
		},
	)

	binlog.SubscribeBinlogEventsOnBootstrapV2(
		binlog.SubscribeBinlogOption{
			ContinueOnErr: true,
			Pipeline: pump.Pipeline{
				Schema: miso.GetPropStr(mysql.PropMySQLSchema),
				Table:  "path",
				Stream: BinlogStreamPathUpdated,
			},
			Concurrency: 1,
			Listener: func(rail miso.Rail, t pump.StreamEvent) error {
				rail.Infof("Subscribed to path changes")
				if err := BatchLoadRoleAccessCache(rail); err != nil {
					rail.Errorf("Failed to BatchLoadRoleAccessCache, %v", err)
				}
				if err := LoadPublicAccessCache(rail); err != nil {
					rail.Errorf("Failed to LoadPublicAccessCache, %v", err)
				}
				return nil
			},
		},
	)

	binlog.SubscribeBinlogEventsOnBootstrapV2(
		binlog.SubscribeBinlogOption{
			ContinueOnErr: true,
			Pipeline: pump.Pipeline{
				Schema: miso.GetPropStr(mysql.PropMySQLSchema),
				Table:  "role_resource",
				Stream: BinlogStreamRoleResUpdated,
			},
			Concurrency: 1,
			Listener: func(rail miso.Rail, t pump.StreamEvent) error {
				rail.Infof("Subscribed to role_resource changes")
				if err := BatchLoadRoleAccessCache(rail); err != nil {
					rail.Errorf("Failed to BatchLoadRoleAccessCache, %v", err)
				}
				return nil
			},
		},
	)

	binlog.SubscribeBinlogEventsOnBootstrapV2(
		binlog.SubscribeBinlogOption{
			ContinueOnErr: true,
			Pipeline: pump.Pipeline{
				Schema: miso.GetPropStr(mysql.PropMySQLSchema),
				Table:  "path_resource",
				Stream: BinlogStreamPathResUpdated,
			},
			Concurrency: 1,
			Listener: func(rail miso.Rail, t pump.StreamEvent) error {
				rail.Infof("Subscribed to path_resource changes")
				if err := BatchLoadRoleAccessCache(rail); err != nil {
					rail.Errorf("Failed to BatchLoadRoleAccessCache, %v", err)
				}
				return nil
			},
		},
	)

	return nil
}
