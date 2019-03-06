package gearbox

import (
	"bufio"
	"fmt"
	"gearbox/box/vm"
	"gearbox/util"
	"github.com/apcera/libretto/ssh"
	"github.com/apcera/libretto/virtualmachine/virtualbox"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

type Vm struct {
	VmName   string
	Instance virtualbox.VM
	Status   string
	OvaFile   string

	// SSH related - Need to fix this. It's used within CreateVm()
	SshUsername  string
	SshPassword  string
	SshPublicKey string

	// Status polling delays.
	NoWait      bool
	WaitDelay   time.Duration
	WaitRetries int

	// Console related.
	ConsoleHost     string
	ConsolePort     string
	ConsoleOkString string
	ConsoleReadWait time.Duration
	ShowConsole     bool
}
type VmArgs Vm

const VmDefaultName = "Gearbox"
const VmDefaultOvaFileName = "box/vm/Gearbox.ova"
const VmDefaultWaitDelay = time.Second
const VmDefaultWaitRetries = 30
const VmDefaultConsoleHost = "127.0.0.1"
const VmDefaultConsolePort = "2023"
const VmDefaultConsoleOkString = "Gearbox Heartbeat"
const VmDefaultShowConsole = false
const VmDefaultConsoleReadWait = time.Second * 5

const VmError = "error"
const VmUnknown = "unknown"
const VmAbsent = "absent"
const VmHalted = "halted"
const VmRunning = "running"
const VmStarted = "started"
const VmGearBoxOK = "ok"
const VmGearBoxNOK = "nok"

// //////////////////////////////////////////////////////////////////////////////
// Gearbox related
func (me *Gearbox) StartVM(vmArgs VmArgs) error {

	vm := NewVm(*me, vmArgs)

	err := vm.StartVm()

	return err
}


func (me *Gearbox) StopVM(vmArgs VmArgs) error {

	vm := NewVm(*me, vmArgs)

	err := vm.StopVm()

	return err
}


func (me *Gearbox) RestartVM(vmArgs VmArgs) error {

	vm := NewVm(*me, vmArgs)

	err := vm.RestartVm()

	return err
}


func (me *Gearbox) StatusVM(vmArgs VmArgs) (string, error) {

	vm := NewVm(*me, vmArgs)

	err := vm.StatusVm()
	if err != nil {
		return vm.Status, err
	}

	err = vm.StatusApi()
	if err != nil {
		return vm.Status, err
	}

	switch vm.Status {
		case VmUnknown:
			fmt.Printf("\r👎 %s: VM & API in an unknown state.\n", me.Config.VmName)

		case VmHalted:
			fmt.Printf("\r👎 %s: VM halted. API halted.\n", me.Config.VmName)

		case VmRunning:
			fmt.Printf("\r👎 %s: VM running. API halted.\n", me.Config.VmName)

		case VmStarted:
			fmt.Printf("\r👎 %s: VM running. API halted.\n", me.Config.VmName)

		case VmGearBoxOK:
			fmt.Printf("\r👍 %s: VM running. API running.\n", me.Config.VmName)

		case VmGearBoxNOK:
			fmt.Printf("\r👎 %s: VM running. API halted.\n", me.Config.VmName)
	}

	return vm.Status, err
}


func (me *Gearbox) CreateVM(vmArgs VmArgs) (string, error) {

	vm := NewVm(*me, vmArgs)

	err := vm.CreateVm()
	if err != nil {
		return vm.Status, err
	}

	err = vm.StatusApi()
	if err != nil {
		return vm.Status, err
	}

	fmt.Printf("\r👎 %s: VM & API in an unknown state.\n", me.Config.VmName)

	return vm.Status, err
}


// //////////////////////////////////////////////////////////////////////////////
// Low-level related
func NewVm(gb Gearbox, args ...VmArgs) *Vm {
	var _args VmArgs
	if len(args) > 0 {
		_args = args[0]
	}

	if _args.VmName == "" {
		_args.VmName = VmDefaultName
	}

	if _args.OvaFile == "" {

		// The '/' will become a problem on Windows
		_args.OvaFile = gb.HostConnector.GetUserConfigDir() + "/" + VmDefaultOvaFileName
		// The OvaFile is created from an export from within VirtualBox.
		// VBoxManage export Gearbox -o Gearbox.ova --options manifest
		// This was the best way to create a base template, avoiding too much code bloat.
		// And allows multiple VM frameworks to be used with libretto.
		// It doesn't include the ISO image yet as it is too large.
		// Once the ISO image size has been reduced, we can do this:
		// VBoxManage export Gearbox -o Gearbox.ova --options iso,manifest
		if _, err := os.Stat(_args.OvaFile); os.IsNotExist(err) {
			err := vm.RestoreAssets(gb.HostConnector.GetUserConfigDir(), VmDefaultOvaFileName)
			if err != nil {
				fmt.Printf("\r👎 %s: VM OVA file cannot be created in %s.\n", _args.VmName, _args.OvaFile)
			}
		}
	}

	if _args.WaitDelay == 0 {
		_args.WaitDelay = VmDefaultWaitDelay
	}

	if _args.WaitRetries == 0 {
		_args.WaitRetries = VmDefaultWaitRetries
	}

	if _args.ConsoleHost == "" {
		_args.ConsoleHost = VmDefaultConsoleHost
	}

	if _args.ConsolePort == "" {
		_args.ConsolePort = VmDefaultConsolePort
	}

	if _args.ConsoleOkString == "" {
		_args.ConsoleOkString = VmDefaultConsoleOkString
	}

	if _args.ConsoleReadWait == 0 {
		_args.ConsoleReadWait = VmDefaultConsoleReadWait
	}

	if _args.SshUsername == "" {
		_args.SshUsername = SshDefaultUsername
	}

	if _args.SshPassword == "" {
		_args.SshPassword = SshDefaultPassword
	}

	if _args.SshPublicKey == "" {
		_args.SshPublicKey = SshDefaultKeyFile
	}

	_args.Instance = virtualbox.VM{
		Name: _args.VmName,
		Src: _args.OvaFile,
		Credentials: ssh.Credentials{
			// Need a way of obtaining this.
			SSHUser: _args.SshUsername,
			SSHPassword: _args.SshPassword,
			SSHPrivateKey: _args.SshPublicKey,
		},
	}

	vm := &Vm{}
	*vm = Vm(_args)

	// Query VB to see if it exists.
	// If not return nil.

	return vm
}

func (me *Vm) WaitForVmState(waitForState string, displayString string) error {

	var waitCount int

	spinner := util.NewSpinner(util.SpinnerArgs{
		Text:    displayString,
		ExitOK:  displayString + " - OK",
		ExitNOK: displayString + " - FAILED",
	})
	spinner.Start()

	for waitCount = 0; waitCount < me.WaitRetries; waitCount++ {
		err := me.StatusVm()
		if err != nil {
			spinner.Stop(false)
			return err
		}
		if me.Status == waitForState {
			spinner.Stop(true)
			break
		}

		time.Sleep(me.WaitDelay)
		spinner.Update(fmt.Sprintf("%s [%d]", displayString, waitCount))
	}

	return nil
}

func (me *Vm) WaitForConsole(displayString string, waitFor time.Duration) error {

	if me == nil {
		// Throw software error.
		return nil
	}

	err := me.StatusVm()
	if err != nil {
		return err
	}

	// TRUE - show the spinner on console.
	displaySpinner := me.ShowConsole == false && displayString != ""

	if me.Status == VmRunning {
		spinner := util.NewSpinner(util.SpinnerArgs{
			Text:    displayString,
			ExitOK:  displayString + " - OK",
			ExitNOK: displayString + " - FAILED",
		})

		if displaySpinner == true {
			// We want to display just a spinner instead of console output.
			spinner.Start()
		}

		// Connect to this console
		conn, err := net.Dial("tcp", me.ConsoleHost+":"+me.ConsolePort)
		defer conn.Close()
		if err != nil {
			return err
		}

		exitWhen := time.Now().Add(time.Second * waitFor)

		readBuffer := make([]byte, 512)
		for waitCount := 0; time.Now().Unix() < exitWhen.Unix(); waitCount++ {
			err = conn.SetDeadline(time.Now().Add(me.ConsoleReadWait))
			if err != nil {
				return err
			}

			bytesRead, err := bufio.NewReader(conn).Read(readBuffer)
			// bytesRead, err := conn.Read(readBuffer)
			// readBuffer, err := bufio.NewReader(conn).ReadString('\n')
			// bytesRead := len(readBuffer)
			if err != nil {
				me.Status = VmGearBoxNOK
				if displaySpinner == true {
					spinner.Stop(false)
				}
				break
			}

			if bytesRead > 0 {
				if me.ShowConsole == true {
					fmt.Printf("%s", string(readBuffer[:bytesRead]))
				}

				apiSplit := strings.Split(string(readBuffer[:bytesRead]), ";")
				if len(apiSplit) > 1 {
					match, _ := regexp.MatchString(me.ConsoleOkString, apiSplit[1])
					if match == true {
						if apiSplit[2] == "OK" {
							me.Status = VmGearBoxOK
						} else {
							me.Status = VmGearBoxNOK
						}
						if displaySpinner == true {
							spinner.Stop(true)
						}
						break
					}
				}
			}

			time.Sleep(me.WaitDelay)
			if displaySpinner == true {
				spinner.Update(fmt.Sprintf("%s [%d]", displayString, waitCount))
			}
		}

		if me.ShowConsole == true {
			fmt.Printf("\n\n# Exiting Console.\n")
		}
	}

	return nil
}

func (me *Vm) StartVm() error {

	if me == nil {
		// Throw software error.
		return nil
	}

	err := me.StatusVm()
	if err != nil {
		return err
	}
	if me.Status == VmRunning || me.Status == VmStarted {
		return nil
	}

	err = me.Instance.Start()
	if err != nil {
		return err
	}
	if me.NoWait == false {
		err := me.WaitForVmState(VmRunning, me.VmName+" VM: Starting")
		if err != nil {
			return err
		}

		// Check final state of the system from the top down.
		err = me.WaitForConsole(me.VmName+" : Starting", 30)
		if err != nil {
			return err
		}
	}

	return nil
}

func (me *Vm) StopVm() error {

	if me == nil {
		// Throw software error.
		return nil
	}

	err := me.StatusVm()
	if err != nil {
		return err
	}
	if me.Status == VmHalted {
		return nil
	}

	err = me.Instance.Halt()
	if err != nil {
		return err
	}

	if me.NoWait == false {
		err := me.WaitForVmState(VmHalted, me.VmName+" VM: Stopping")
		if err != nil {
			return err
		}

		// Check final state of the system from the top down.
		err = me.WaitForConsole(me.VmName+" : Stopping", 30)
		if err != nil {
			return err
		}
	}

	return nil
}

func (me *Vm) RestartVm() error {

	var err error

	if me == nil {
		// Throw software error.
		return nil
	}

	err = me.StatusVm()
	if err != nil {
		return err
	}

	switch me.Status {
	case VmGearBoxOK:
		fallthrough
	case VmGearBoxNOK:
		fallthrough
	case VmRunning:
		fallthrough
	case VmStarted:
		err := me.StopVm()
		if err != nil {
			return err
		}
		err = me.StartVm()
		if err != nil {
			return err
		}

	case VmHalted:
		err := me.StartVm()
		if err != nil {
			return err
		}

	case VmUnknown:

	}

	if me.Status != VmRunning {
		// Throw an error.
	}

	return err
}

func (me *Vm) StatusVm() error {

	if me == nil {
		// Throw software error.
		me.Status = VmUnknown
		return nil
	}

	state, err := me.Instance.GetState()
	if err != nil {
		me.Status = VmError
	} else {
		me.Status = state
	}

	return err
}

func (me *Vm) StatusApi() error {

	if me == nil {
		// Throw software error.
		return nil
	}

	if me.Status == VmRunning {
		err := me.WaitForConsole("", 10)
		if err != nil {
			return err
		}
	}

	return nil
}

func (me *Vm) CreateVm() error {

	if me == nil {
		// Throw software error.
		me.Status = VmUnknown
		return nil
	}

	// Check if the VM is already there.
	state, err := me.Instance.GetState()
	if err != nil {
		// Doesn't exist - great!
		if _, err := os.Stat(me.OvaFile); os.IsNotExist(err) {
			fmt.Printf("\r👎 %s: VM OVA file does not exist in %s.\n", me.VmName, me.OvaFile)
			return err
		}

		err = me.Instance.Provision()
		if err != nil {
			me.Status = VmError
		}
	} else {
		// Already created!
		fmt.Printf("\r👎 %s: Cannot create. VM already exists and is in state %s.\n", me.VmName, state)
	}

	return err
}


func (me *Vm) ValidateVm() error {

	if me == nil {
		// Throw software error.
		me.Status = VmUnknown
		return nil
	}
/*
	state, err := me.Instance.Provision()
	if err != nil {
		me.Status = VmError
	} else {
		me.Status = state
	}
	return err
*/
	return nil
}


