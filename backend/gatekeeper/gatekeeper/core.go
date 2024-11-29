package gatekeeper

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/curtisnewbie/miso/middleware/jwt"
	"github.com/curtisnewbie/miso/middleware/logbot"
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/curtisnewbie/miso/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cast"
)

var (
	errPathNotFound = miso.NewErrf("Path not found")

	timerHistoVec     *prometheus.HistogramVec = miso.NewPromHistoVec("gatekeeper_request_duration", []string{"url"})
	timerExclPath                              = util.NewSet[string]()
	histoVecTimerPool                          = sync.Pool{
		New: func() any {
			return miso.NewVecTimer(timerHistoVec)
		},
	}

	whitelistPatterns []string
)

const (
	AttrAuthInfo = "gk.auth.info"

	PropTimerExclPath         = "gatekeeper.timer.path.excl"
	PropWhitelistPathPatterns = "gatekeeper.whitelist.path.patterns"
)

type ServicePath struct {
	ServiceName string
	Path        string
}

func Bootstrap(args []string) {
	logbot.EnableLogbotErrLogReport()
	miso.PreServerBootstrap(prepareServer)
	miso.BootstrapServer(args)
}

func prepareServer(rail miso.Rail) error {

	miso.Infof("gatekeeper version: %v", Version)

	// disable trace propagation, we are the entry point
	common.LoadBuiltinPropagationKeys()
	miso.SetProp(miso.PropServerPropagateInboundTrace, false)

	// whitelisted path patterns
	whitelistPatterns = miso.GetPropStrSlice(PropWhitelistPathPatterns)

	// create proxy
	proxy := miso.NewHttpProxy("/", ResolveServiceTarget)

	if !miso.IsProdMode() {
		proxy.AddFilter(ReqTimeLogFilter)
	}

	// healthcheck filter
	healthcheckPath := miso.GetPropStr(miso.PropConsulHealthcheckUrl)
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

	// pprof filter
	if !miso.IsProdMode() && miso.GetPropBool(miso.PropServerPprofEnabled) {
		miso.PerfLogExclPath("/debug/pprof")
		miso.PerfLogExclPath("/debug/pprof/cmdline")
		miso.PerfLogExclPath("/debug/pprof/profile")
		miso.PerfLogExclPath("/debug/pprof/symbol")
		miso.PerfLogExclPath("/debug/pprof/trace")
		proxy.AddFilter(PProfFilter)
	}

	// gatekeeper filter
	proxy.AddFilter(AuthFilter)
	proxy.AddFilter(AccessFilter)
	proxy.AddFilter(TraceFilter)

	// paths that are not measured by prometheus timer
	timerExclPath.AddAll(miso.GetPropStrSlice(PropTimerExclPath))
	timerExclPath.Add(miso.GetPropStr(miso.PropMetricsRoute))
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
	healthcheckPath := miso.GetPropStr(miso.PropConsulHealthcheckUrl)

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

	timer := histoVecTimerPool.Get().(*miso.VecTimer)
	timer.Reset()
	defer func() {
		timer.ObserveDuration(pc.ProxyPath)
		histoVecTimerPool.Put(timer)
	}()

	next()
}

func PProfFilter(pc *miso.ProxyContext, next func()) {

	w, r := pc.Inb.Unwrap()
	if r.URL.Path == "/debug/pprof/cmdline" {
		pprof.Cmdline(w, r)
		return
	} else if r.URL.Path == "/debug/pprof/profile" {
		pprof.Profile(w, r)
		return
	} else if r.URL.Path == "/debug/pprof/symbol" {
		pprof.Symbol(w, r)
		return
	} else if r.URL.Path == "/debug/pprof/trace" {
		pprof.Trace(w, r)
		return
	} else if strings.HasPrefix(r.URL.Path, "/debug/pprof") {
		pprof.Index(w, r)
		return
	}

	next()
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
	if proxyPath == miso.GetPropStr(miso.PropConsulHealthcheckUrl) {
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
	authorization := r.Header.Get("Authorization")
	rail.Debugf("Authorization: %v", authorization)

	// no token available
	if authorization == "" {
		next()
		return
	}

	// decode jwt token, extract claims and build a user struct as attr
	tkn, err := jwt.JwtDecode(authorization)
	rail.Debugf("DecodeToken, tkn: %v, err: %v", tkn, err)

	// token invalid, but the public endpoints are still accessible, so we don't stop here
	if err != nil || !tkn.Valid {
		rail.Debugf("Token invalid, %v", err)
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

	var roleNo string
	var u common.User = common.NilUser()

	if v, ok := pc.GetAttr(AttrAuthInfo); ok && v != nil {
		u = v.(common.User)
		roleNo = u.RoleNo
	}

	inWhitelist := false
	for _, pat := range whitelistPatterns {
		if ok, _ := path.Match(pat, r.URL.Path); ok {
			inWhitelist = true
			break
		}
	}

	var cr CheckResAccessResp
	if inWhitelist {
		cr = CheckResAccessResp{true}
	} else {
		var err error
		cr, err = ValidateResourceAccess(*rail, CheckResAccessReq{
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
	start := time.Now()
	_, r := pc.Inb.Unwrap()
	next()
	pc.Rail.Infof("%-6v %-60v [%s]", r.Method, r.RequestURI, time.Since(start))
}
