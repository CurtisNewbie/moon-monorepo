package vfm

import "github.com/curtisnewbie/miso/util"

var (
	vfmPool = util.NewAsyncPool(500, 20)
)
