package includes

import (
	"gearbox"
	"gearbox/stat"
)

var _ gearbox.Config = (*TestConfig)(nil)

type TestConfig struct {
	gearbox.Config
}

func NewTestConfig(gb gearbox.Gearbox) gearbox.Config {
	return &TestConfig{
		Config: gearbox.NewConfiguration(gb),
	}
}

func (me *TestConfig) Initialize() (status stat.Status) {
	status = me.Load()
	if !status.IsError() {
		status = me.WriteFile()
	}
	return status
}

func (me *TestConfig) WriteFile() (status stat.Status) {
	return stat.NewOkStatus("Test file written (wink, wink)")
}
