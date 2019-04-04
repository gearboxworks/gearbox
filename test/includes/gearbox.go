package includes

import (
	"gearbox/gearbox"
	"gearbox/global"
	"gearbox/hostapi"
	"gearbox/status"
	"gearbox/test/mock"
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
	gb.SetHostApi(hostapi.NewHostApi(gb))
	sts := gb.Initialize()
	return gb, sts
}
