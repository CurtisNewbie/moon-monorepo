package server

import (
	"github.com/curtisnewbie/miso/middleware/logbot"
	"github.com/curtisnewbie/miso/middleware/user-vault/auth"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/user-vault/internal/postbox"
	"github.com/curtisnewbie/user-vault/internal/vault"
	"github.com/curtisnewbie/user-vault/internal/web"
)

func BootstrapServer(args []string) {
	common.LoadBuiltinPropagationKeys()
	logbot.EnableLogbotErrLogReport()

	miso.PreServerBootstrap(vault.SubscribeBinlogEvent)
	miso.PreServerBootstrap(postbox.PrepareLongPollHandler)
	miso.PreServerBootstrap(func(rail miso.Rail) error {
		vault.RegisterInternalPathResourcesOnBootstrapped([]auth.Resource{
			{Code: web.ResourceManageResources, Name: "Manage Resources Access"},
			{Code: web.ResourceManagerUser, Name: "Admin Manage Users"},
			{Code: web.ResourceBasicUser, Name: "Basic User Operation"},
			{Code: web.ResourceQueryNotification, Name: "Query Notifications"},
			{Code: web.ResourceCreateNotification, Name: "Create Notifications"},
		})
		return nil
	})

	miso.PreServerBootstrap(printVersion)
	miso.PreServerBootstrap(vault.ScheduleTasks)
	miso.PreServerBootstrap(web.RegisterRoutes)
	miso.PreServerBootstrap(postbox.InitPipeline)
	miso.PreServerBootstrap(postbox.SubscribeBinlogChanges)
	miso.PostServerBootstrapped(vault.CreateMonitoredServiceWatches)
	miso.BootstrapServer(args)
}

func printVersion(rail miso.Rail) error {
	rail.Infof("user-vault (monorepo) version: %v", Version)
	return nil
}
