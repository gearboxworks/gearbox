package virtualbox

import (
	"gearbox/box/external/hypervisor"
	"github.com/gearboxworks/go-status/only"
)

type HostOnlyNic struct {
	Name        string
	Index       int
	Ip          string
	Netmask     string
	DhcpLowerIp string
	DhcpUpperIp string
}
type HostOnlyNicArgs HostOnlyNic

func NewHostOnlyNic(args ...*HostOnlyNicArgs) *HostOnlyNic {

	if len(args) == 0 {
		args = make([]*HostOnlyNicArgs, 1)
	}

	if args[0].Ip == "" {
		args[0].Ip = DefaultHostOnlyIp
	}

	if args[0].Name == "" {
		args[0].Name = args[0].Ip
	}

	if args[0].Netmask == "" {
		args[0].Netmask = DefaultHostOnlyNetmask
	}

	if args[0].DhcpLowerIp == "" {
		args[0].DhcpLowerIp = DefaultHostOnlyDhcpLowerIp
	}

	if args[0].DhcpUpperIp == "" {
		args[0].DhcpUpperIp = DefaultHostOnlyDhcpUpperIp
	}
	hon := HostOnlyNic(*args[0])
	return &hon

}

func (me *HostOnlyNic) initialize(vm hypervisor.VirtualMachiner) (err error) {
	for range only.Once {
		var n string
		n, err = CmdFindHostOnlyNet(vm)
		if err != nil {
			break
		}
		me.Name = n
	}
	return err
}
