package vault

import (
	"testing"

	"github.com/curtisnewbie/miso/miso"
)

func TestLoadResourceMonitorConf(t *testing.T) {
	if err := miso.LoadConfigFromFile("conf.yml", miso.EmptyRail()); err != nil {
		t.Fatal(err)
	}
	c := LoadMonitoredServices()
	t.Logf("loaded: %+v", c)
}
