package virtualbox

import (
	"fmt"
	"gearbox/box/external/hypervisor"
	"gearbox/ensure"
	"gearbox/eventbroker/eblog"
	"github.com/gearboxworks/go-status/only"
)


func findFirstNic(vm hypervisor.VirtualMachiner) error {

	var err error

	for range only.Once {
		err = ensure.NotNil(vm)
		if err != nil {
			break
		}

		var cr *CommandResult
		cr = ManageVm(vm, "list", "bridgedifs", "-s")
		if cr.Error != nil {
			logger.Error("unable to list bridged interfaces for VM '%s' %s",
				vm.GetId(),
				cr.String(),
			)
			break
		}

		var nic KeyValueMap
		dr, ok := cr.Decode(Stdout, ':')
		if ok == true {
			var nics KeyValuesMap
			nics, ok = dr.decodeBridgeIfs()
			if ok == false {
				err = fmt.Errorf("no NICs found for VM '%s'", vm.GetName())
				break
			}

			for _, nic = range nics {
				if nic["FirstNic"] == "Yes" {
					break
				}
			}
		}

		if nic == nil {
			err = fmt.Errorf("no NICs found for VM '%s'", vm.GetName())
			break
		}

		logger.Debug("using NIC '%s' for VM '%s'", nic, vm.GetName())
	}

	eblog.LogIfNil(vm, err)
	eblog.LogIfError(err)

	return err
}
