package postbox

import "github.com/curtisnewbie/miso/util"

var (
	PostboxPool *util.AsyncPool
)

func init() {
	PostboxPool = util.NewAsyncPool(500, 10)
}
