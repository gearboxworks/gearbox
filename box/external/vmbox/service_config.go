package vmbox

import (
	"bytes"
	"errors"
	"gearbox/box/external/virtualbox"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/osdirs"
	"time"
)

type ServiceConfig struct {
	Name        msgs.Address
	Version     string
	Console     Console
	Ssh         Ssh
	IconFile    osdirs.File
	VmDir       osdirs.Dir
	HostOnlyNic *virtualbox.HostOnlyNic

	retryMax   int
	retryDelay time.Duration
	cmdStdout  bytes.Buffer
	cmdStderr  bytes.Buffer
	vmInfo     virtualbox.KeyValueMap
	vmNics     virtualbox.KeyValuesMap
}

func (me *ServiceConfig) EnsureNotNil() error {
	var err error

	switch {
	case me == nil:
		err = errors.New("VmBox Service instance is nil")
	}

	return err
}
