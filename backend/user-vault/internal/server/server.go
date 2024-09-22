package server

import (
	"github.com/curtisnewbie/miso/middleware/logbot"
	"github.com/curtisnewbie/miso/middleware/user-vault/auth"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/user-vault/internal/postbox"
	"github.com/curtisnewbie/user-vault/internal/vault"
)

func BootstrapServer(args []string) {
	common.LoadBuiltinPropagationKeys()
	logbot.EnableLogbotErrLogReport()

	miso.PreServerBootstrap(vault.SubscribeBinlogEvent)
	miso.PreServerBootstrap(func(rail miso.Rail) error {
		vault.RegisterInternalPathResourcesOnBootstrapped([]auth.Resource{
			{Code: vault.ResourceManageResources, Name: "Manage Resources Access"},
			{Code: vault.ResourceManagerUser, Name: "Admin Manage Users"},
			{Code: vault.ResourceBasicUser, Name: "Basic User Operation"},
			{Code: postbox.ResourceQueryNotification, Name: "Query Notifications"},
			{Code: postbox.ResourceCreateNotification, Name: "Create Notifications"},
		})
		return nil
	})

	miso.PreServerBootstrap(printVersion)
	miso.PreServerBootstrap(vault.ScheduleTasks)
	miso.PreServerBootstrap(postbox.RegisterRoutes)
	miso.PreServerBootstrap(postbox.InitPipeline)
	miso.PostServerBootstrapped(vault.CreateMonitoredServiceWatches)
	miso.BootstrapServer(args)
}

func printVersion(rail miso.Rail) error {
	rail.Infof("user-vault version: %v", Version)
	return nil
}
