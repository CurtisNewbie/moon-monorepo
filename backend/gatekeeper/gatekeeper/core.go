package gatekeeper

import (
	"net/http"
	"strings"

	"github.com/curtisnewbie/miso/errs"
	"github.com/curtisnewbie/miso/flow"
	"github.com/curtisnewbie/miso/middleware/jwt"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util/hash"
	"github.com/curtisnewbie/miso/util/strutil"
	uvault "github.com/curtisnewbie/user-vault/api"
	"github.com/spf13/cast"
)

var (
	errPathNotFound = errs.NewErrf("Path not found")

	timerHisto    = miso.NewPromHisto("gatekeeper_all_request_duration")
	timerExclPath = hash.NewSet[string]()

	whitelistPatterns []string
)

const (
	AttrAuthInfo = "gk.auth.info"

	HeaderAuthorization = "Authorization"
	CookieAuthorization = "Gatekeeper_Authorization"
)

type ServicePath struct {
	ServiceName string
	Path        string
}

func Bootstrap(args []string) {
	miso.PreServerBootstrap(prepareServer)
	miso.BootstrapServer(args)
}

func prepareServer(rail miso.Rail) error {

	miso.Infof("gatekeeper (monorepo) version: %v", Version)

	// disable trace propagation, we are the entry point
	miso.SetProp(miso.PropServerPropagateInboundTrace, false)

	// whitelisted path patterns
	whitelistPatterns = miso.GetPropStrSlice(PropWhitelistPathPatterns)

	// create proxy
	proxy := miso.NewHttpProxy("/", miso.NewDynProxyTargetResolver())
	proxy.AddReqTimeLogFilter(func(path string) bool {
		return timerExclPath.Has(path)
	})
	proxy.AddFilter(IpFilter)

	// healthcheck filter
	proxy.AddHealthcheckFilter()

	// metrics filter
	proxy.AddMetricsFilter(timerHisto, func(path string) bool {
		return timerExclPath.Has(path)
	})

	// pprof filter for gatekeeper itself
	if err := proxy.AddDebugFilter(true); err != nil {
		return err
	}
	rail.Infof("Enabled pprof/trace api for gatekeeper")

	proxy.AddFilter(AuthFilter)
	proxy.AddFilter(AccessFilter)
	proxy.AddFilter(TraceFilter)

	// paths that are not measured by prometheus timer
	timerExclPath.AddAll(miso.GetPropStrSlice(PropTimerExclPath))
	timerExclPath.Add(miso.GetPropStr(miso.PropMetricsRoute))
	rail.Infof("Timer excluded paths: %v", timerExclPath)
	return nil
}

func AuthFilter(pc *miso.ProxyContext, next func()) {

	rail := pc.Rail
	_, r := pc.Inb.Unwrap()

	// header
	authorization := r.Header.Get(HeaderAuthorization)

	// fallback to cookie
	if authorization == "" {
		ck, err := r.Cookie(CookieAuthorization)
		if err == nil && ck != nil {
			authorization = ck.Value
		}
	}
	rail.Debugf("Authorization: %v", authorization)

	// no token available
	if authorization == "" {
		next()
		return
	}

	// parse bearer
	if s, ok := miso.ParseBearer(authorization); ok {
		authorization = s
	}

	// decode jwt token, extract claims and build a user struct as attr
	tkn, err := jwt.JwtDecode(authorization)
	rail.Debugf("DecodeToken, tkn: %v, err: %v", tkn, err)

	// token invalid, but the public endpoints are still accessible, so we don't stop here
	if err != nil || !tkn.Valid {
		rail.Infof("Token invalid, %v", err)
		next()
		return
	}

	// extract the user info from it
	claims := tkn.Claims
	var user flow.User

	if v, ok := claims["username"]; ok {
		user.Username = cast.ToString(v)
	}
	if v, ok := claims["userno"]; ok {
		user.UserNo = cast.ToString(v)
	}
	if v, ok := claims["roleno"]; ok {
		user.RoleNo = cast.ToString(v)
	}
	pc.SetAttr(AttrAuthInfo, user)
	rail.Debugf("user: %#v", user)
	rail.Debugf("set user to proxyContext: %v", pc)

	next()
}

func AccessFilter(pc *miso.ProxyContext, next func()) {

	w, r := pc.Inb.Unwrap()
	rail := pc.Rail

	var roleNo string
	var u flow.User = flow.NilUser()

	if v, ok := pc.GetAttr(AttrAuthInfo); ok && v != nil {
		u = v.(flow.User)
		roleNo = u.RoleNo
	}

	inWhitelist := false
	for _, pat := range whitelistPatterns {
		if ok := strutil.MatchPath(pat, r.URL.Path); ok {
			inWhitelist = true
			break
		}
	}

	var cr uvault.CheckResAccessResp
	if inWhitelist {
		cr = uvault.CheckResAccessResp{Valid: true}
	} else {
		var err error
		cr, err = ValidateResourceAccess(*rail, uvault.CheckResAccessReq{
			Url:    r.URL.Path,
			Method: r.Method,
			RoleNo: roleNo,
		})

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			rail.Warnf("Request forbidden, err: %v", err)
			return
		}
	}

	if !cr.Valid {
		rail.Warnf("Request forbidden (resource access not authorized), url: %v, user: %+v", r.URL.Path, u)

		// authenticated, but doesn't have enough authority to access the endpoint
		if !u.IsNil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// token invalid or expired
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	next()
}

func TraceFilter(pc *miso.ProxyContext, next func()) {

	v, ok := pc.GetAttr(AttrAuthInfo)
	if ok && v != nil {
		u := v.(flow.User)
		*pc.Rail = flow.StoreUser(*pc.Rail, u)
		pc.Rail.Debugf("Setup trace for user info, rail: %+v", pc.Rail)
	}

	next()
}

func IpFilter(pc *miso.ProxyContext, next func()) {
	// IP is provided by nginx, but just in case the nginx is missing, we identify the remote IP ourselves
	_, r := pc.Inb.Unwrap()

	if miso.GetPropBool(PropOverwriteRemoteIp) || r.Header.Get("x-forwarded-for") == "" {
		v := r.RemoteAddr
		if i := strings.LastIndexByte(v, ':'); i > -1 {
			v = v[0:i]
		}
		r.Header.Set("x-forwarded-for", v)
		pc.Rail.Debugf("Overwrote remote IP: %v", v)
	}
	next()
}
