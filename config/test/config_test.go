package test

import (
	"fmt"
	"gearbox/config"
	"gearbox/gearbox"
	"gearbox/status"
	"gearbox/status/is"
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
	gb := &gearbox.Gearbox{
		OsSupport: mock.NewOsSupport(t),
	}
	c := config.NewConfig(gb.GetOsSupport())
	gb.SetConfig(c)

	t.Run("Initialize", func(t *testing.T) {
		sts := c.Initialize()
		if status.IsError(sts) {
			t.Error(sts.Message())
		}
	})
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
					if is.Error(sts) {
						t.Error(sts.Message())
						return
					}
					var basedir types.AbsoluteDir
					basedir, sts = c.GetHostBasedir(p.Basedir)
					if is.Error(sts) {
						t.Error(sts.Message())
						return
					}
					path := util.ExtractRelativePath(types.AbsoluteFilepath(util.FileDir(fp)), basedir)
					if _, ok := ProjectPaths[path]; !ok {
						t.Error(fmt.Sprintf("path '%s' not found", path))
					}
				})
			})
		}
	})

}
