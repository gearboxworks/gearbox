package gearbox

import (
	"fmt"
	"github.com/apcera/libretto/virtualmachine/virtualbox"
)


type Vm struct {
	VmName	string
	Instance virtualbox.VM
	Status	string
}
type VmArgs Vm


func NewVm(gb Gearbox, args ...VmArgs) *Vm {
	var _args VmArgs
	if len(args)>0 {
		_args = args[0]
	}

	if _args.VmName == ""{
		_args.VmName = "Gearbox"
	}

	vm := &Vm{}
	*vm = Vm(_args)

	// Query VB to see if it exists.
	// If not return nil.

	return vm
}


func Exists(gb Gearbox, args ...VmArgs) *Vm {

	return nil
}


func (me *Vm) StartVm() error {

	//vm.VM = virtualbox.VM{Name: "Gearbox", }

	fmt.Printf("Yup\n")

	return nil
}


func (me *Vm) StopVm() error {

	return nil
}


func (me *Vm) StatusVm() error {

	return nil
}


func (me *Vm) RestartVm() error {

	return nil
}
