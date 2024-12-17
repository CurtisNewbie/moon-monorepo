package gatekeeper

import (
	"github.com/curtisnewbie/miso/miso"
	uvault "github.com/curtisnewbie/user-vault/api"
)

// Check whether this role has access to the url
func ValidateResourceAccess(c miso.Rail, req uvault.CheckResAccessReq) (uvault.CheckResAccessResp, error) {
	return uvault.SendCheckResAccessReq(c, req)
}
