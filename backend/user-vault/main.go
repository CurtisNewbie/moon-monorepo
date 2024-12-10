package main

import (
	"os"

	"github.com/curtisnewbie/user-vault/internal/server"
)

func main() {
	server.BootstrapServer(os.Args)
}
