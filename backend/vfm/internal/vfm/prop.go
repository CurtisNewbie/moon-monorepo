package vfm

import (
	fstore "github.com/curtisnewbie/mini-fstore/api"
	"github.com/curtisnewbie/miso/miso"
)

// misoconfig-section: General Configuration
const (
	// misoconfig-prop: Externally accessible host |
	PropVfmSiteHost = "vfm.site.host"
)

// misoconfig-section: VFM Configuration
const (
	// misoconfig-prop: Temporary file path for bootmarks files | "/tmp/vfm"
	PropTempPath = "vfm.temp-path"

	// misoconfig-prop: Maximum depth of an archive's folder tree when unpacking, counted in path segments from the target directory | 32
	PropUnpackMaxDepth = "vfm.unpack.max-depth"
)

// misoconfig-default-start
func init() {
	miso.SetDefProp(PropTempPath, "/tmp/vfm")
	miso.SetDefProp(PropUnpackMaxDepth, fstore.DefaultUnzipMaxDepth)
}

// misoconfig-default-end
