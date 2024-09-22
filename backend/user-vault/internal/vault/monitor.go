package vault

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
	"github.com/spf13/cast"
)

const (
	DefaultMonitorPath = "/auth/resource"
)

var (
	monitorServiceTickser = []*miso.TickRunner{}
)

type MonitorConf struct {
	Monitor []MonitoredService
}

type MonitoredService struct {
	Service string
	Path    string
	All     bool
}

func LoadMonitoredServices() []MonitoredService {
	var c MonitorConf
	miso.UnmarshalFromProp(&c)
	for i, m := range c.Monitor {
		if m.Path == "" {
			m.Path = DefaultMonitorPath
			c.Monitor[i] = m
		}
	}
	return c.Monitor
}

type QueryResourcePathRes struct {
	Resources []CreateResReq
	Paths     []CreatePathReq
}

func QueryResourcePath(rail miso.Rail, server miso.Server, service string, path string) (QueryResourcePathRes, error) {
	var resp miso.GnResp[QueryResourcePathRes]
	err := miso.NewTClient(rail, server.BuildUrl(path)).
		Require2xx().
		Get().
		Json(&resp)
	if err != nil {
		return QueryResourcePathRes{}, fmt.Errorf("failed to query resource path from monitored service, server: %+v, service: %v, %w",
			server, service, err)
	}
	return resp.Res()
}

func CreateMonitoredServiceWatches(rail miso.Rail) error {
	services := LoadMonitoredServices()
	for i := range services {
		s := services[i]
		if err := CreateMonitoredServiceWatch(rail, s); err != nil {
			return err
		}
	}
	return nil
}

func QueryResourcePathAsync(rail miso.Rail, server miso.Server, m MonitoredService) {
	monitorPool.Go(func() {
		rail.Debugf("Polling resource path for service: %v from %#v", m.Service, server)
		res, err := QueryResourcePath(rail, server, m.Service, m.Path)
		if err != nil {
			rail.Errorf("monitor service %v failed, %v", m.Service, err)
		} else {
			rail.Debugf("service %v (%v:%v), returned resouces/paths: %+v", m.Service, server.Address, server.Port, res)
			user := common.NilUser() // just to satisfy the method, it's always a zero value
			for _, r := range res.Resources {
				if err := CreateResourceIfNotExist(rail, r, user); err != nil {
					rail.Errorf("failed to create resource, req: %+v, %v", r, err)
				}
			}
			for _, r := range res.Paths {
				if err := CreatePath(rail, r, user); err != nil {
					rail.Errorf("failed to create path, req: %+v, %v", r, err)
				}
			}
		}
	})
}

func TriggerResourcePathCollection(rail miso.Rail, m MonitoredService) {
	servers := miso.GetServerList().ListServers(rail, m.Service)
	if len(servers) < 1 {
		return
	}

	if m.All {
		for i := range servers {
			server := servers[i]
			QueryResourcePathAsync(rail, server, m)
		}
	} else {
		sort.Slice(servers, func(i, j int) bool {
			var irt int = math.MaxInt
			var jrt int = math.MaxInt
			if irts, ok := servers[i].Meta[miso.ConsulMetaRegisterTime]; ok {
				irt = cast.ToInt(irts)
			}
			if jrts, ok := servers[j].Meta[miso.ConsulMetaRegisterTime]; ok {
				jrt = cast.ToInt(jrts)
			}
			if irt > jrt {
				return true
			}
			return false
		})
		server := servers[0] // the most recently created service
		QueryResourcePathAsync(rail, server, m)
	}
}

func CreateMonitoredServiceWatch(rail miso.Rail, m MonitoredService) error {
	if err := miso.SubscribeServerChanges(rail, m.Service, func() {
		TriggerResourcePathCollection(miso.EmptyRail(), m)
	}); err != nil {
		return fmt.Errorf("failed to subscribe server chagnes, service: %v, %v", m.Service, err)
	}

	tr := miso.NewTickRuner(time.Minute*5, func() {
		TriggerResourcePathCollection(miso.EmptyRail(), m)
	})
	monitorServiceTickser = append(monitorServiceTickser, tr)
	tr.Start()
	return nil
}
