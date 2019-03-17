package integration

import (
	"fmt"
	"gearbox"
	"gearbox/test"
	"gearbox/test/mock"
	"path/filepath"
	"testing"
)

var ProjectPaths = map[string]bool{
	"/site1.local": true,
	"/site2.local": true,
	"/site3.test":  true,
	"/site4.test":  true,
	"/site5":       true,
	"/site6":       true,
}

func TestEmptyConfig(t *testing.T) {
	mgb := &mock.GearboxObj{
		HostConnector: mock.NewHostConnector(t),
	}
	c := gearbox.NewConfiguration(mgb)
	mgb.SetConfig(c)

	t.Run("Initialize", func(t *testing.T) {
		status := c.Initialize()
		if status.IsError() {
			t.Error(status.Message)
		}
	})
	t.Run("Boxname", func(t *testing.T) {
		if c.GetBoxname() != gearbox.Boxname {
			t.Error(fmt.Sprintf("Want: %s, Got %s",
				gearbox.Boxname,
				c.GetBoxname(),
			))
		}
	})
	basedir := mgb.HostConnector.GetUserHomeDir()
	t.Run("ProjectMap", func(t *testing.T) {
		for k, p := range c.GetProjectMap() {
			t.Run(k, func(t *testing.T) {
				fp, status := p.GetFilepath()
				if status.IsError() {
					t.Error(status.Message)
				}
				path := test.ParseRelativePath(basedir, filepath.Dir(fp))
				if _, ok := ProjectPaths[path]; !ok {
					t.Error(fmt.Sprintf("path '%s' not found", path))
				}
			})
		}
	})

}
