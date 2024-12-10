package config

import (
	"github.com/curtisnewbie/miso/miso"
)

const (
	PropDefaultUserRole = "user-vault.default-role"
)

func DefaultUserRole() string {
	return miso.GetPropStr(PropDefaultUserRole)
}
