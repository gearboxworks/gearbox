package unit

import (
	"gearbox"
	"gearbox/test/includes"
	"gearbox/test/mock"
	"testing"
)

var GlobalOptionsTable = []*gearbox.GlobalOptions{
	{NoCache: T, IsDebug: T},
	{NoCache: F, IsDebug: F},
	{NoCache: F, IsDebug: T},
	{NoCache: T, IsDebug: F},
}

func TestHostApiGlobalOptions(t *testing.T) {
	for _, glopt := range GlobalOptionsTable {
		t.Run(glopt.Debug(), func(t *testing.T) {
			testGlobalOption(t, glopt)
		})
	}
}

func testGlobalOption(t *testing.T, glopt *gearbox.GlobalOptions) {
	gb := gearbox.NewApp(&gearbox.Args{
		HostConnector: mock.NewHostConnector(t),
		GlobalOptions: glopt,
	})
	gb.SetConfig(includes.NewTestConfig(gb))
	status := gb.Initialize()
	if status.IsError() {
		t.Error(status.Message)
	}
}
