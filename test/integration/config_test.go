package integration

import (
	"fmt"
	"gearbox/config"
	"gearbox/gearbox"
	"gearbox/status"
	"gearbox/test"
	"gearbox/test/mock"
	"gearbox/types"
	"gearbox/util"
	"testing"
)

var ProjectPaths = map[types.RelativePath]bool{
	"/site1.local": true,
	"/site2.local": true,
	"/site3.test":  true,
	"/site4.test":  true,
	"/site5":       true,
	"/site6":       true,
}

func TestEmptyConfig(t *testing.T) {
	mgb := &gearbox.Gearbox{
		OsSupport: mock.NewOsSupport(t),
	}
	c := config.NewConfig(mgb.GetOsSupport())
	mgb.SetConfig(c)

	t.Run("Initialize", func(t *testing.T) {
		sts := c.Initialize()
		if status.IsError(sts) {
			t.Error(sts.Message())
		}
	})
	t.Run("Brandname", func(t *testing.T) {
		if c.GetBoxname() != gearbox.Brandname {
			t.Error(fmt.Sprintf("Want: %s, Got %s",
				gearbox.Brandname,
				c.GetBoxname(),
			))
		}
	})
	basedir := mgb.OsSupport.GetUserHomeDir()
	t.Run("ProjectMap", func(t *testing.T) {
		pm, sts := c.GetProjectMap()
		if status.IsError(sts) {
			t.Error(sts.Message())
			return
		}
		for k, p := range pm {
			t.Run(string(k), func(t *testing.T) {
				t.Run(string(k), func(t *testing.T) {
					fp, sts := p.GetFilepath()
					if status.IsError(sts) {
						t.Error(sts.Message())
						return
					}
					path := test.ParseRelativePath(
						types.AbsoluteDir(basedir),
						types.AbsoluteFilepath(util.FileDir(fp)),
					)
					if _, ok := ProjectPaths[path]; !ok {
						t.Error(fmt.Sprintf("path '%s' not found", path))
					}
				})
			})
		}
	})

}
