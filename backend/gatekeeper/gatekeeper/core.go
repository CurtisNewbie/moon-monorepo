package gatekeeper

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"strings"
	"time"

	"github.com/curtisnewbie/miso/middleware/jwt"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	uvault "github.com/curtisnewbie/user-vault/api"
	"github.com/spf13/cast"
)

var (
	errPathNotFound = miso.NewErrf("Path not found")

	timerHisto    = miso.NewPromHisto("gatekeeper_all_request_duration")
	timerExclPath = util.NewSet[string]()

	whitelistPatterns []string
)

const (
	AttrAuthInfo           = "gk.auth.info"
	AttrPprofAuthenticated = "gk.pprof.auth.pass"

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
	common.LoadBuiltinPropagationKeys()
	miso.SetProp(miso.PropServerPropagateInboundTrace, false)

	// whitelisted path patterns
	whitelistPatterns = miso.GetPropStrSlice(PropWhitelistPathPatterns)

	// create proxy
	proxy := miso.NewHttpProxy("/", ResolveServiceTarget)
	proxy.AddFilter(ReqTimeLogFilter)
	proxy.AddFilter(IpFilter)

	// healthcheck filter
	healthcheckPath := miso.GetPropStr(miso.PropHealthCheckUrl)
	if !util.IsBlankStr(healthcheckPath) {
		miso.PerfLogExclPath(healthcheckPath)
		proxy.AddFilter(HealthcheckFilter)
	}

	// metrics filter
	metricsEndpoint := miso.GetPropStr(miso.PropMetricsRoute)
	if !util.IsBlankStr(metricsEndpoint) {
		miso.PerfLogExclPath(metricsEndpoint)
		if miso.GetPropBool(miso.PropMetricsEnabled) {
			proxy.AddFilter(MetricsFilter)
		}
	}

	// pprof filter for gatekeeper itself
	if !miso.IsProdMode() || miso.GetPropBool(miso.PropServerPprofEnabled) {
		bearer := ""
		if miso.IsProdMode() {
			bearer = miso.GetPropStr(miso.PropServerPprofAuthBearer)
			if bearer == "" {
				return miso.NewErrf("Configuration '%v' for pprof authentication is missing, but pprof authentication is enabled", miso.PropServerPprofAuthBearer)
			}
		}
		miso.PerfLogExclPath("/debug/pprof")
		miso.PerfLogExclPath("/debug/pprof/cmdline")
		miso.PerfLogExclPath("/debug/pprof/profile")
		miso.PerfLogExclPath("/debug/pprof/symbol")
		miso.PerfLogExclPath("/debug/pprof/trace")
		proxy.AddFilter(PProfFilter(bearer))
		rail.Infof("Enabled pprof api for gatekeeper")
	}

	proxy.AddFilter(ProxyPprofAuthFilter)
	proxy.AddFilter(AuthFilter)
	proxy.AddFilter(AccessFilter)
	proxy.AddFilter(TraceFilter)

	// paths that are not measured by prometheus timer
	timerExclPath.AddAll(miso.GetPropStrSlice(PropTimerExclPath))
	timerExclPath.Add(miso.GetPropStr(miso.PropMetricsRoute))
	rail.Infof("Timer excluded paths: %v", timerExclPath)
	return nil
}

func parseServicePath(url string) (ServicePath, error) {
	rurl := []rune(url)[1:] // remove leading '/'

	// root path, invalid request
	if len(rurl) < 1 {
		return ServicePath{}, errPathNotFound
	}

	start := 0
	for i := range rurl {
		if rurl[i] == '/' && i > 0 {
			start = i
			break
		}
	}

	if start < 1 {
		return ServicePath{}, errPathNotFound
	}

	return ServicePath{
		ServiceName: string(rurl[0:start]),
		Path:        string(rurl[start:]),
	}, nil
}

func HealthcheckFilter(pc *miso.ProxyContext, next func()) {
	healthcheckPath := miso.GetPropStr(miso.PropHealthCheckUrl)

	// check if it's a healthcheck endpoint (for consul), we don't really return anything, so it's fine to expose it
	if pc.ProxyPath == healthcheckPath {
		w, _ := pc.Inb.Unwrap()
		if miso.IsHealthcheckPass(*pc.Rail) {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		return
	}

	next()
}

func MetricsFilter(pc *miso.ProxyContext, next func()) {

	metricsEndpoint := miso.GetPropStr(miso.PropMetricsRoute)

	w, r := pc.Inb.Unwrap()
	if r.URL.Path == metricsEndpoint {
		miso.PrometheusHandler().ServeHTTP(w, r)
		return
	}

	if timerExclPath.Has(r.URL.Path) {
		next()
		return
	}

	timer := miso.NewHistTimer(timerHisto)
	defer timer.ObserveDuration()

	next()
}

func PProfFilter(bearer string) func(pc *miso.ProxyContext, next func()) {
	return func(pc *miso.ProxyContext, next func()) {
		w, r := pc.Inb.Unwrap()

		p := r.URL.Path
		if v, ok := strings.CutPrefix(r.URL.Path, "/gatekeeper"); ok {
			p = v
		}
		if strings.HasPrefix(p, "/debug/pprof") {
			if bearer != "" {
				token, ok := miso.ParseBearer(r.Header.Get("Authorization"))
				if !ok || token != bearer {
					miso.Debugf("Bearer authorization failed, missing bearer token or token mismatch, %v %v", r.Method, r.RequestURI)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}

			if p == "/debug/pprof/cmdline" {
				pprof.Cmdline(w, r)
				return
			} else if p == "/debug/pprof/profile" {
				pprof.Profile(w, r)
				return
			} else if p == "/debug/pprof/symbol" {
				pprof.Symbol(w, r)
				return
			} else if p == "/debug/pprof/trace" {
				pprof.Trace(w, r)
				return
			} else {
				if name, found := strings.CutPrefix(p, "/debug/pprof/"); found && name != "" {
					pprof.Handler(name).ServeHTTP(w, r)
					return
				}

				pprof.Index(w, r)
				return
			}
		}

		next()
	}
}

type GatewayError struct {
	StatusCode int
}

func (g GatewayError) Status() int {
	return g.StatusCode
}

func (g GatewayError) Error() string {
	return fmt.Sprintf("gateway error, %v", g.StatusCode)
}

func ResolveServiceTarget(rail miso.Rail, proxyPath string) (string, error) {
	if proxyPath == miso.GetPropStr(miso.PropHealthCheckUrl) {
		return proxyPath, nil
	}
	if proxyPath == miso.GetPropStr(miso.PropMetricsRoute) {
		return proxyPath, nil
	}
	if strings.HasPrefix(proxyPath, "/debug/pprof") {
		return proxyPath, nil
	}

	// parse the request path, extract service name, and the relative url for the backend server
	var sp ServicePath
	var err error
	if sp, err = parseServicePath(proxyPath); err != nil {
		rail.Warnf("Invalid request, %v", err)
		return "", GatewayError{StatusCode: 404}
	}
	rail.Debugf("parsed service path: %#v", sp)
	target, err := miso.GetServiceRegistry().ResolveUrl(miso.EmptyRail(), sp.ServiceName, sp.Path)
	if err != nil {
		rail.Warnf("ServiceRegistry ResolveUrl failed, %v", err)
		return "", GatewayError{StatusCode: 404}
	}
	return target, nil
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
	var user common.User

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

	if strings.Contains(r.URL.Path, "/debug/pprof") {
		if v, ok := pc.GetAttr(AttrPprofAuthenticated); ok && v.(bool) {
			next()
			return
		}
	}

	var roleNo string
	var u common.User = common.NilUser()

	if v, ok := pc.GetAttr(AttrAuthInfo); ok && v != nil {
		u = v.(common.User)
		roleNo = u.RoleNo
	}

	inWhitelist := false
	for _, pat := range whitelistPatterns {
		if ok := util.MatchPath(pat, r.URL.Path); ok {
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
		u := v.(common.User)
		*pc.Rail = common.StoreUser(*pc.Rail, u)
		pc.Rail.Debugf("Setup trace for user info, rail: %+v", pc.Rail)
	}

	next()
}

func ReqTimeLogFilter(pc *miso.ProxyContext, next func()) {
	_, r := pc.Inb.Unwrap()

	if timerExclPath.Has(r.URL.Path) {
		next()
		return
	}

	start := time.Now()
	next()
	pc.Rail.Infof("%-6v %-60v [%s]", r.Method, r.RequestURI, time.Since(start))
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

func ProxyPprofAuthFilter(pc *miso.ProxyContext, next func()) {
	w, r := pc.Inb.Unwrap()
	if strings.Contains(r.URL.Path, "/debug/pprof") {
		bearer := miso.GetPropStr(PropProxyPprofBearer)

		if bearer == "" && miso.IsProdMode() { // production must enable pprof authentication
			pc.Rail.Infof("Attempt to request '%v', authentication for pprof is mandatory in production, rejected", r.RequestURI)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if bearer != "" {
			authorization := r.Header.Get("Authorization")
			provided, ok := miso.ParseBearer(authorization)
			if !ok || bearer != provided {
				pc.Rail.Infof("Attempt to request '%v', but bearer authentication failed, rejected", r.RequestURI)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
	}
	pc.SetAttr(AttrPprofAuthenticated, true)
	next()
}
