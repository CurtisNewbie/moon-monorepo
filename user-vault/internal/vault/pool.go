package vault

import "github.com/curtisnewbie/miso/util"

var (
	monitorPool = util.NewAsyncPool(500, 10)
	commonPool  = util.NewAsyncPool(500, 10)
)
