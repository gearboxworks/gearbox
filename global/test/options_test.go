package test

import (
	"gearbox/gearbox"
	gopt "gearbox/global"
	"gearbox/test/includes"
	"gearbox/test/mock"
	"github.com/gearboxworks/go-status/is"
	"testing"
)

const T = true
const F = false

var GlobalOptionsTable = []*gopt.Options{
	{NoCache: T, IsDebug: T},
	{NoCache: F, IsDebug: F},
	{NoCache: F, IsDebug: T},
	{NoCache: T, IsDebug: F},
}

func TestGlobalOptions(t *testing.T) {
	for _, glopt := range GlobalOptionsTable {
		t.Run(glopt.Debug(), func(t *testing.T) {
			testGlobalOption(t, glopt)
		})
	}
}

func testGlobalOption(t *testing.T, glopt *gopt.Options) {
	gb := gearbox.NewGearbox(&gearbox.Args{
		OsSupport:     mock.NewOsSupport(t),
		GlobalOptions: glopt,
	})
	gb.SetConfig(includes.NewTestConfig(gb))

	sts := gb.Initialize()
	if is.Error(sts) {
		t.Error(sts.Message())
	}
}
