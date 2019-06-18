package includes

import (
	"gearbox/config"
	"gearbox/gearbox"
	"github.com/gearboxworks/go-status"
)

var _ config.Configer = (*TestConfig)(nil)

type TestConfig struct {
	config.Configer
}

func NewTestConfig(gb gearbox.Gearboxer) config.Configer {
	return &TestConfig{
		Configer: config.NewConfig(gb.GetOsBridge()),
	}
}

func (me *TestConfig) Initialize() (sts status.Status) {
	sts = me.Load()
	if !status.IsError(sts) {
		sts = me.WriteFile()
	}
	return sts
}
