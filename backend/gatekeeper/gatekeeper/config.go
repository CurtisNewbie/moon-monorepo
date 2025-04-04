package gatekeeper

import "github.com/curtisnewbie/miso/miso"

// misoconfig-section: Gatekeeper Configuration
const (
	// misoconfig-prop: slice of paths that are not measured by prometheus timer
	PropTimerExclPath = "gatekeeper.timer.path.excl"

	// misoconfig-prop: slice of path patterns that do not require authorization and authentication
	PropWhitelistPathPatterns = "gatekeeper.whitelist.path.patterns"

	// misoconfig-prop: always overwrite remote ip address in `x-forwarded-for` header (by default, there should be a nginx sitting right before the gatekeeper as a reverse proxy, this the default value for this setting is false) | false
	PropOverwriteRemoteIp = "gatekeeper.overwrite-remote-ip"

	// misoconfig-prop: bearer authentication token for pprof endpoints (not just gatekeeper, but also servers behind it), mandatory for production mode; if bearer token is missing, all requests to `*/debug/pprof/*` are rejected |
	PropProxyPprofBearer = "gatekeeper.proxy.pprof.bearer"
)

// misoconfig-default-start
func init() {
	miso.SetDefProp(PropOverwriteRemoteIp, false)
}

// misoconfig-default-end
