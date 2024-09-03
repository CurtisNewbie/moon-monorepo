package vfm

import (
	"embed"
	"os"

	"github.com/curtisnewbie/miso/middleware/logbot"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/vfm/internal/schema"
)

var (
	SchemaFs embed.FS
)

func PrepareServer() {
	common.LoadBuiltinPropagationKeys()
	schema.EnableSchemaMigrate()
	logbot.EnableLogbotErrLogReport()
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
	rail.Infof("vfm version: %v", Version)
	return nil
}
