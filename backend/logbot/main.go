package main

import (
	"os"

	"github.com/curtisnewbie/log-bot/logbot"
	"github.com/curtisnewbie/miso/miso"
)

func main() {
	miso.PreServerBootstrap(logbot.BeforeServerBootstrap)
	miso.PostServerBootstrap(logbot.AfterServerBootstrapped)
	miso.BootstrapServer(os.Args)
}
