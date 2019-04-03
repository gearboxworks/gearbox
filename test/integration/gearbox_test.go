package integration

import (
	"gearbox"
	"gearbox/hostapi"
	"gearbox/status"
	"gearbox/test/includes"
	"gearbox/test/mock"
	"testing"
)

func TestGearboxConstructAndInitialize(t *testing.T) {
	_, sts := GearboxConstructAndInitialize(t)
	if status.IsError(sts) {
		t.Error(sts.Message())
	}
}

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
		GlobalOptions: &gopt.GlobalOptions{
			NoCache: true,
			IsDebug: false,
		},
	})
	gb.SetConfig(includes.NewTestConfig(gb))
	gb.SetHostApi(hostapi.NewHostApi(gb))
	sts := gb.Initialize()
	return gb, sts
}
