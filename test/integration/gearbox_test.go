package integration

import (
	"gearbox"
	"gearbox/stat"
	"gearbox/test/includes"
	"gearbox/test/mock"
	"testing"
)

func TestGearboxConstructAndInitialize(t *testing.T) {
	_, status := GearboxConstructAndInitialize(t)
	if status.IsError() {
		t.Error(status.Message)
	}
}

func GearboxConstructAndInitialize(t *testing.T) (gearbox.Gearbox, stat.Status) {
	gb := gearbox.NewApp(&gearbox.Args{
		HostConnector: mock.NewHostConnector(t),
		GlobalOptions: &gearbox.GlobalOptions{
			NoCache: true,
			IsDebug: false,
		},
	})
	gb.SetConfig(includes.NewTestConfig(gb))
	status := gb.Initialize()
	return gb, status
}
