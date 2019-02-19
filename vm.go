package gearbox

import (
	"fmt"
	"github.com/apcera/libretto/virtualmachine/virtualbox"
	"time"
)


type Vm struct {
	VmName        string
	Instance      virtualbox.VM
	Status        string
	WaitDelay     time.Duration
	WaitRetries   int
	WaitIndicator waitIndicator
}
type VmArgs Vm

type waitIndicator func(int)


const VmUnknown = "unknown"
const VmStopped = "stopped"
const VmHalted = "halted"
const VmRunning = "running"
const VmStarted = "started"
const VmStarting = "starting"
const VmStopping = "stopping"



// //////////////////////////////////////////////////////////////////////////////
// Gearbox related
func (me *Gearbox) StartVM(nowait bool) error {

	vm := NewVm(*me, VmArgs{
		VmName: me.Config.VmName,
	})
	err := vm.StartVm(nowait)

	return err
}


func (me *Gearbox) StopVM(nowait bool) error {

	vm := NewVm(*me, VmArgs{
		VmName: me.Config.VmName,
	})
	err := vm.StopVm(nowait)

	return err
}


func (me *Gearbox) RestartVM(nowait bool) error {

	vm := NewVm(*me, VmArgs{
		VmName: me.Config.VmName,
	})

	err := vm.RestartVm(nowait)

	return err
}


func (me *Gearbox) StatusVM() (string, error) {

	vm := NewVm(*me, VmArgs{
		VmName: me.Config.VmName,
	})
	state, err := vm.StatusVm()

	fmt.Printf("%s VM is in state: %s\n", me.Config.VmName, state)

	return state, err
}


// //////////////////////////////////////////////////////////////////////////////
// Low-level related
func NewVm(gb Gearbox, args ...VmArgs) *Vm {
	var _args VmArgs
	if len(args)>0 {
		_args = args[0]
	}

	if _args.VmName == "" {
		_args.VmName = "Gearbox"
	}

	if _args.WaitDelay == 0 {
		_args.WaitDelay = time.Second
	}

	if _args.WaitRetries == 0 {
		_args.WaitRetries = 30
	}

	// Should be refactored elsewhere.
	if _args.WaitIndicator == nil {
		_args.WaitIndicator = UserWaitIndicator
	}

	_args.Instance = virtualbox.VM{Name: "Gearbox", }

	vm := &Vm{}
	*vm = Vm(_args)

	// Query VB to see if it exists.
	// If not return nil.

	return vm
}


func Exists(gb Gearbox, args ...VmArgs) *Vm {

	return nil
}


func (me *Vm) WaitForState(waitForState string) error {

	for i := 0; i < me.WaitRetries; i++ {
		state, err := me.StatusVm()
		if err != nil {
			return err
		}
		if state == waitForState {
			me.WaitIndicator(WIStopOK)
			break
		}

		me.WaitIndicator(WISpin)
		time.Sleep(me.WaitDelay)

		if i == 1 {
			// Generally only just indicate to the user any waiting if we are indeed waiting.
			me.WaitIndicator(WIStart)
		}
	}

	return nil
}


func (me *Vm) StartVm(nowait bool) error {

	if me == nil {
		// Throw software error.
		return nil
	}

	state, err := me.StatusVm()
	if err != nil {
		return err
	}
	if state == VmRunning || state == VmStarted || state == VmStarting {
		return nil
	}

	err = me.Instance.Start()
	if err != nil {
		return err
	}
	if nowait == false {
		err := me.WaitForState(VmRunning)
		if err != nil {
			return err
		}
	}

	return nil
}


func (me *Vm) StopVm(nowait bool) error {

	if me == nil {
		// Throw software error.
		return nil
	}

	state, err := me.StatusVm()
	if err != nil {
		return err
	}
	if state == VmHalted || state == VmStopped || state == VmStopping {
		return nil
	}

	err = me.Instance.Halt()
	if err != nil {
		return err
	}
	if nowait == false {
		err := me.WaitForState(VmHalted)
		if err != nil {
			return err
		}
	}

	return nil
}


func (me *Vm) RestartVm(nowait bool) error {

	var stateBefore string
	var stateAfter string
	var err error

	if me == nil {
		// Throw software error.
		return nil
	}

	stateBefore, err = me.StatusVm()
	if err != nil {
		return err
	}

	switch {
		case stateBefore == VmRunning:
			fallthrough
		case stateBefore == VmStarted:
			err := me.StopVm(false)
			if err != nil {
				return err
			}
			err = me.StartVm(nowait)
			if err != nil {
				return err
			}

		case stateBefore == VmStopped:
			fallthrough
		case stateBefore == VmHalted:
			err := me.StartVm(nowait)
			if err != nil {
				return err
			}

		case stateBefore == VmStarting:

		case stateBefore == VmStopping:

		case stateBefore == VmUnknown:

	}

	stateAfter, err = me.StatusVm()
	if err != nil {
		return err
	}

	if stateAfter != VmRunning {
		// Throw an error.
	}

	return err
}


func (me *Vm) StatusVm() (string, error) {

	if me == nil {
		// Throw software error.
		return "", nil
	}

	state, err := me.Instance.GetState()
	if err != nil {
		return state, err
	}

	return state, nil
}


// Wait Indicators.
const WINew = 0
const WIStart = 1
const WISpin = 2
const WIStopOK = 3
const WIStopNOK = 4

func DefaultWaitIndicator(position int) {

	// Do nothing.
	switch position {
		case WINew:

		case WIStart:

		case WISpin:

		case WIStopOK:

		case WIStopNOK:
	}

}


func UserWaitIndicator(position int) {

	switch position {
		case WINew:

		case WIStart:
			fmt.Printf("\nWaiting")

		case WISpin:
			fmt.Printf(".")

		case WIStopOK:
			fmt.Printf("OK\n")

		case WIStopNOK:
			fmt.Printf("Not OK!\n")
	}

}

