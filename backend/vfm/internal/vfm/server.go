package vfm

import (
	"embed"
	"os"

	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
)

var (
	SchemaFs embed.FS
)

func PrepareServer() {
	common.LoadBuiltinPropagationKeys()
	miso.PreServerBootstrap(SubscribeBinlogChanges)
	miso.PreServerBootstrap(PrintVersion)
	miso.PreServerBootstrap(PrepareEventBus)
	miso.PreServerBootstrap(RegisterHttpRoutes)
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
