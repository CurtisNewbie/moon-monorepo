package server

import (
	"os"

	"github.com/curtisnewbie/acct/internal/flow"
	"github.com/curtisnewbie/acct/internal/web"
	"github.com/curtisnewbie/miso/miso"

	_ "github.com/curtisnewbie/acct/internal/config"
)

func init() {
	miso.PreServerBootstrap(func(rail miso.Rail) error {
		rail.Infof("acct (monorepo) version: %v", Version)
		return nil
	})
}

func BootstrapServer() {
	miso.PreServerBootstrap(PreServerBootstrap)
	miso.PostServerBootstrapped(PostServerBootstrap)
	miso.BootstrapServer(os.Args)
}

func PreServerBootstrap(rail miso.Rail) error {
	// declare http endpoints, jobs/tasks, and other components here
	web.PrepareWebServer(rail)
	flow.LoadCategoryConfs(rail)
	return nil
}

func PostServerBootstrap(rail miso.Rail) error {
	// do stuff right after server being fully bootstrapped
	return nil
}
