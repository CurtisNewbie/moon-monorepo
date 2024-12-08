package vault

import "github.com/curtisnewbie/miso/util"

var (
	monitorPool = util.NewIOAsyncPool()
	commonPool  = util.NewIOAsyncPool()
)
