package vfm

import "github.com/curtisnewbie/miso/miso"

// misoconfig-section: General Configuration
const (
	// misoconfig-prop: Externally accessible host |
	PropVfmSiteHost = "vfm.site.host"
)

// misoconfig-section: VFM Configuration
const (
	// misoconfig-prop: Temporary file path for bootmarks files | "/tmp/vfm"
	PropTempPath = "vfm.temp-path"
)

// misoconfig-default-start
func init() {
	miso.SetDefProp(PropTempPath, "/tmp/vfm")
}

// misoconfig-default-end
