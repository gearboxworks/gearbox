package includes

import (
	"gearbox/config"
	"gearbox/gearbox"
	"gearbox/status"
)

var _ config.Configer = (*TestConfig)(nil)

type TestConfig struct {
	config.Configer
}

func NewTestConfig(gb gearbox.Gearboxer) config.Configer {
	return &TestConfig{
		Configer: config.NewConfig(gb.GetOsSupport()),
	}
}

func (me *TestConfig) Initialize() (sts status.Status) {
	sts = me.Load()
	if !status.IsError(sts) {
		sts = me.WriteFile()
	}
	return sts
}

func (me *TestConfig) WriteFile() (sts status.Status) {
	return status.Success("Test file written (wink, wink)")
}
