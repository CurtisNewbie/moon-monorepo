package config

import (
	"github.com/curtisnewbie/miso/miso"
)

// misoconfig-section: User Configuration
const (

	// misoconfig-prop: Default role no for new user |
	PropDefaultUserRole = "user-vault.default-role"

	// misoconfig-prop: Coalescing window for unread notification count changes. A burst of changes for one user within this window is folded into a single announcement. | 1s
	PropNotificationCountCoalesceWindow = "user-vault.notification-count-coalesce-window"
)

// misoconfig-default-start
func init() {
	miso.SetDefProp(PropNotificationCountCoalesceWindow, "1s")
}

// misoconfig-default-end

func DefaultUserRole() string {
	return miso.GetPropStr(PropDefaultUserRole)
}
