package vfm

import (
	"os"

	"github.com/curtisnewbie/miso/miso"
)

func PrepareServer() {

	miso.PreServerBootstrap(SubscribeBinlogChanges)
	miso.PreServerBootstrap(PrintVersion)
	miso.PreServerBootstrap(PrepareEventBus)
	miso.PreServerBootstrap(PrepareWebServer)
	miso.PreServerBootstrap(MakeTempDirs)
}

func BootstrapServer(args []string) {
	PrepareServer()
	miso.BootstrapServer(os.Args)
}

func PrintVersion(rail miso.Rail) error {
	rail.Infof("vfm (monorepo) version: %v", Version)
	return nil
}
