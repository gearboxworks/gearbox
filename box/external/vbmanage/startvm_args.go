package vbmanage

import (
	"fmt"
	"github.com/gearboxworks/go-status/only"
)

type StartVmArgs struct {
	Name string
	Type string
	Headless bool
}

func (me *StartVmArgs) Strings() Strings {
	return []string{
		me.Name,
		"--type",
		me.Type,
	}
}

func (me *StartVmArgs) Validate() (err error) {
	for range only.Once {
		if me.Name == "" {
			err = fmt.Errorf("VB name empty for '%s' command", StartVmCmd)
			break
		}
		if me.Type == "" {
			err = fmt.Errorf("VB --type empty for '%s' command", StartVmCmd)
			break
		}
	}
	return err
}

