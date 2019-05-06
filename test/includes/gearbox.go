package includes

import (
	"gearbox/api"
	"gearbox/gearbox"
	"gearbox/global"
	"gearbox/test/mock"
	"github.com/gearboxworks/go-status"
	"testing"
)

func GearboxConstructAndInitialize(t *testing.T) (gearbox.Gearboxer, status.Status) {
	//gb := mock.NewGearbox(&mock.GearboxArgs{
	//	OsSupport: mock.NewOsSupport(t),
	//	GlobalOptions: &gearbox.GlobalOptions{
	//		NoCache: true,
	//		IsDebug: false,
	//	},
	//})
	gb := gearbox.NewGearbox(&gearbox.Args{
		OsSupport: mock.NewOsSupport(t),
		GlobalOptions: &global.Options{
			NoCache: true,
			IsDebug: false,
		},
	})
	gb.SetConfig(NewTestConfig(gb))
	gb.SetApi(api.NewApi(gb))
	sts := gb.Initialize()
	return gb, sts
}
