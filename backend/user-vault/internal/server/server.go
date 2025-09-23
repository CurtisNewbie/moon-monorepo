package server

import (
	"github.com/curtisnewbie/miso/middleware/dbquery"
	"github.com/curtisnewbie/miso/middleware/user-vault/auth"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/user-vault/internal/postbox"
	"github.com/curtisnewbie/user-vault/internal/vault"
	"github.com/curtisnewbie/user-vault/internal/web"
)

func BootstrapServer(args []string) {
	common.LoadBuiltinPropagationKeys()

	miso.PreServerBootstrap(
		func(rail miso.Rail) error {
			dbquery.PrepareCreateModelHook()
			dbquery.PrepareUpdateModelHook()
			return nil
		},
		vault.SubscribeBinlogEvent,
		postbox.PrepareLongPollHandler,
		func(rail miso.Rail) error {
			vault.RegisterInternalPathResourcesOnBootstrapped([]auth.Resource{
				{Code: web.ResourceManageResources, Name: "Manage Resources Access"},
				{Code: web.ResourceManagerUser, Name: "Admin Manage Users"},
				{Code: web.ResourceBasicUser, Name: "Basic User Operation"},
				{Code: web.ResourceQueryNotification, Name: "Query Notifications"},
				{Code: web.ResourceCreateNotification, Name: "Create Notifications"},
			})
			return nil
		},
		printVersion,
		vault.ScheduleTasks,
		web.RegisterRoutes,
		postbox.InitPipeline,
		postbox.SubscribeBinlogChanges,
	)

	miso.PostServerBootstrap(vault.CreateMonitoredServiceWatches)
	miso.BootstrapServer(args)
}

func printVersion(rail miso.Rail) error {
	rail.Infof("user-vault (monorepo) version: %v", Version)
	return nil
}
