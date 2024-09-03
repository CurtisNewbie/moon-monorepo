package main

import (
	"os"

	"github.com/curtisnewbie/vfm/internal/vfm"
)

func main() {
	vfm.BootstrapServer(os.Args)
}
