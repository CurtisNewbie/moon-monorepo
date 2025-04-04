package config

import (
	"github.com/curtisnewbie/miso/miso"
)

// misoconfig-section: User Configuration
const (

	// misoconfig-prop: Default role no for new user |
	PropDefaultUserRole = "user-vault.default-role"
)

func DefaultUserRole() string {
	return miso.GetPropStr(PropDefaultUserRole)
}
