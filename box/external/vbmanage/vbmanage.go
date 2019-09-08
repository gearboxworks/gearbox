package vbmanage

import (
	"fmt"
	"gearbox/box/external/hypervisor"
	"github.com/gearboxworks/go-status/only"
	"os/exec"
	"strings"
)

type VbManage struct{
	Vm   hypervisor.VirtualMachiner
}

func NewVbManage(vm hypervisor.VirtualMachiner) *VbManage {
	return &VbManage{Vm:vm}
}


func (me *VbManage) StartVm(args *StartVmArgs) (cr *CommandResult) {
	cr = NewCommandResult()
	cr.Error = args.Validate()
	if cr.Error == nil {
		cr = me.ManageVm(StartVmCmd, args.Strings())
	}
	return cr
}

// ManageVm runs a VBoxManage command.
func (me *VbManage) ManageVm(name VbCmd, args Strings) (cr *CommandResult) {

	var err error
	cr = NewCommandResult()
	command := strings.Join(args, " ")
	for range only.Once {
		vm := me.Vm
		var path string
		path, err = exec.LookPath(ExecutableName)
		if err != nil {
			path = VBoxManagePath
		}

		logger.Debug("EXEC[%v]: %v '%v'", vm.GetId(), path, command)

		cmd := exec.Command(path, args...)
		cmd.Stdout = cr.GetStdout()
		cmd.Stderr = cr.GetStderr()
		err = cmd.Run()
		if err == nil {
			break
		}

		cr.ExitCode = strings.TrimPrefix(err.Error(), "exit status ")

		switch cr.ExitCode {
		case exitCodeMissingVm:
			err = fmt.Errorf("VirtualBox[%s] command error '%v'", vm.GetName(), cmd.Stderr)

		default:
			err = fmt.Errorf("VirtualBox[%s] failed to run command '%v'",
				vm.GetName(),
				err.Error(),
			)
		}

	}

	if err != nil {
		cr.Error = err
		logger.Error("VBoxManage command '%s' failed [%s]: %v",
			command,
			cr.String(),
			err.Error(),
		)
	}

	return cr

}
