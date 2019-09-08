package vmbox

import (
	"gearbox/box/external/hypervisor"
)

var _ hypervisor.SecureSheller = (*Ssh)(nil)

type Ssh struct {
	Host string
	Port string
}

func (me *Ssh) GetHost() string {
	return me.Host
}

func (me *Ssh) GetPort() string {
	return me.Port
}
