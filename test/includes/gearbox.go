package includes

import (
	"gearbox/api"
	"gearbox/gearbox"
	"gearbox/global"
	"github.com/gearboxworks/go-status"
	"testing"
)

func GearboxConstructAndInitialize(t *testing.T) (gearbox.Gearboxer, status.Status) {
	//gb := mock.NewGearbox(&mock.GearboxArgs{
	//	OsBridge: mock.NewOsBridge(t),
	//	GlobalOptions: &gearbox.GlobalOptions{
	//		NoCache: true,
	//		IsDebug: false,
	//	},
	//})
	gb := gearbox.NewGearbox(&gearbox.Args{
		OsBridge: mock.NewOsBridge(t),
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
