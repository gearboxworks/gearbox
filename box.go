package gearbox

import (
	"bufio"
	"bytes"
	"fmt"
	"gearbox/box/vm"
	"gearbox/util"
	"github.com/apcera/libretto/ssh"
	lvm "github.com/apcera/libretto/virtualmachine"
	"github.com/apcera/libretto/virtualmachine/virtualbox"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type Box struct {
	BoxName  string
	Instance virtualbox.VM
	Status   string
	OvaFile  string

	// SSH related - Need to fix this. It's used within CreateBox()
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
type BoxArgs Box

const BoxDefaultName = "Gearbox"
const BoxDefaultOvaFileName = "box/vm/Gearbox.ova"
const BoxDefaultWaitDelay = time.Second
const BoxDefaultWaitRetries = 90
const BoxDefaultConsoleHost = "127.0.0.1"
const BoxDefaultConsolePort = "2023"
const BoxDefaultConsoleOkString = "GearBox Heartbeat"
const BoxDefaultShowConsole = false
const BoxDefaultConsoleReadWait = time.Second * 5

const BoxError = "error"
const BoxUnknown = "unknown"
const BoxAbsent = "absent"
const BoxHalted = "halted"
const BoxRunning = "running"
const BoxStarted = "started"
const BoxGearBoxOK = "ok"
const BoxGearBoxNOK = "nok"

// //////////////////////////////////////////////////////////////////////////////
// Gearbox related
func (me *Gearbox) StartBox(boxArgs BoxArgs) error {

	box := NewBox(*me, boxArgs)

	err := box.StartBox()

	return err
}

func (me *Gearbox) StopBox(boxArgs BoxArgs) error {

	box := NewBox(*me, boxArgs)

	err := box.StopBox()

	return err
}

func (me *Gearbox) RestartBox(boxArgs BoxArgs) error {

	box := NewBox(*me, boxArgs)

	err := box.RestartBox()

	return err
}

func (me *Gearbox) GetBoxStatus(boxArgs BoxArgs) (string, error) {

	box := NewBox(*me, boxArgs)

	err := box.GetBoxStatus()
	if err != nil {
		return box.Status, err
	}

	err = box.GetApiStatus()
	if err != nil {
		return box.Status, err
	}

	switch box.Status {
	case BoxUnknown:
		fmt.Printf("\r👎 %s: Box status: VM & Heartbeat in an unknown state.\n", me.Config.BoxName)

	case BoxHalted:
		fmt.Printf("\r👎 %s: Box status VM & Heartbeat halted.\n", me.Config.BoxName)

	case BoxRunning:
		fmt.Printf("\r👎 %s: Box status: VM running, Heartbeat halted.\n", me.Config.BoxName)

	case BoxStarted:
		fmt.Printf("\r👎 %s: Box status: VM running, Heartbeat halted.\n", me.Config.BoxName)

	case BoxGearBoxNOK:
		fmt.Printf("\r👎 %s: Box status: VM running, Heartbeat halted.\n", me.Config.BoxName)

	case BoxGearBoxOK:
		fmt.Printf("\r👍 %s: Box status: VM running, Heartbeat running.\n", me.Config.BoxName)

	}

	return box.Status, err
}

func (me *Gearbox) CreateBox(boxArgs BoxArgs) (string, error) {

	box := NewBox(*me, boxArgs)

	err := box.CreateBox()
	if err != nil {
		return box.Status, err
	}

	err = box.GetApiStatus()
	if err != nil {
		return box.Status, err
	}

	fmt.Printf("\r👎 %s: Box status: VM & Heartbeat in an unknown state.\n", me.Config.BoxName)

	return box.Status, err
}

// //////////////////////////////////////////////////////////////////////////////
// Low-level related
func NewBox(gb Gearbox, args ...BoxArgs) *Box {
	var _args BoxArgs
	if len(args) > 0 {
		_args = args[0]
	}

	if _args.BoxName == "" {
		_args.BoxName = BoxDefaultName
	}

	if _args.OvaFile == "" {

		// The '/' will become a problem on Windows
		_args.OvaFile = gb.HostConnector.GetUserConfigDir() + "/" + BoxDefaultOvaFileName
		// The OvaFile is created from an export from within VirtualBox.
		// VBoxManage export Gearbox -o Gearbox.ova --options manifest
		// This was the best way to create a base template, avoiding too much code bloat.
		// And allows multiple VM frameworks to be used with libretto.
		// It doesn't include the ISO image yet as it is too large.
		// Once the ISO image size has been reduced, we can do this:
		// VBoxManage export Gearbox -o Gearbox.ova --options iso,manifest
		if _, err := os.Stat(_args.OvaFile); os.IsNotExist(err) {
			err := vm.RestoreAssets(gb.HostConnector.GetUserConfigDir(), BoxDefaultOvaFileName)
			if err != nil {
				fmt.Printf("\r👎 %s: VM OVA file cannot be created in %s.\n", _args.BoxName, _args.OvaFile)
			}
		}
	}

	if _args.WaitDelay == 0 {
		_args.WaitDelay = BoxDefaultWaitDelay
	}

	if _args.WaitRetries == 0 {
		_args.WaitRetries = BoxDefaultWaitRetries
	}

	if _args.ConsoleHost == "" {
		_args.ConsoleHost = BoxDefaultConsoleHost
	}

	if _args.ConsolePort == "" {
		_args.ConsolePort = BoxDefaultConsolePort
	}

	if _args.ConsoleOkString == "" {
		_args.ConsoleOkString = BoxDefaultConsoleOkString
	}

	if _args.ConsoleReadWait == 0 {
		_args.ConsoleReadWait = BoxDefaultConsoleReadWait
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
		Name: _args.BoxName,
		Src:  _args.OvaFile,
		Credentials: ssh.Credentials{
			// Need a way of obtaining this.
			SSHUser:       _args.SshUsername,
			SSHPassword:   _args.SshPassword,
			SSHPrivateKey: _args.SshPublicKey,
		},
	}

	box := &Box{}
	*box = Box(_args)

	// Query VB to see if it exists.
	// If not return nil.

	return box
}

func (me *Box) WaitForBoxState(waitForState string, displayString string) error {

	var waitCount int

	spinner := util.NewSpinner(util.SpinnerArgs{
		Text:    displayString,
		ExitOK:  displayString + " - OK",
		ExitNOK: displayString + " - FAILED",
	})
	spinner.Start()

	for waitCount = 0; waitCount < me.WaitRetries; waitCount++ {
		err := me.GetBoxStatus()
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

func (me *Box) WaitForConsole(displayString string, waitFor time.Duration) error {

	if me == nil {
		// Throw software error.
		return nil
	}

	err := me.GetBoxStatus()
	if err != nil {
		return err
	}

	// TRUE - show the spinner on console.
	displaySpinner := me.ShowConsole == false && displayString != ""

	if me.Status == BoxRunning {
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
				me.Status = BoxGearBoxNOK
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
							me.Status = BoxGearBoxOK
						} else {
							me.Status = BoxGearBoxNOK
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

func (me *Box) StartBox() error {

	if me == nil {
		// Throw software error.
		return nil
	}

	err := me.GetBoxStatus()
	if err != nil {
		return err
	}
	if me.Status == BoxRunning || me.Status == BoxStarted {
		return nil
	}

	err = me.Instance.Start()
	if err != nil {
		return err
	}
	if me.NoWait == false {
		err := me.WaitForBoxState(BoxRunning, me.BoxName+" Box (VM): Starting")
		if err != nil {
			return err
		}

		// Check final state of the system from the top down.
		err = me.WaitForConsole(me.BoxName+" : Starting", 30)
		if err != nil {
			return err
		}
	}

	return nil
}

func (me *Box) StopBox() error {

	if me == nil {
		// Throw software error.
		return nil
	}

	err := me.GetBoxStatus()
	if err != nil {
		return err
	}
	if me.Status == BoxHalted {
		return nil
	}

	err = me.Instance.Halt()
	if err != nil {
		return err
	}

	if me.NoWait == false {
		err := me.WaitForBoxState(BoxHalted, me.BoxName+" Box (VM): Stopping")
		if err != nil {
			return err
		}

		// Check final state of the system from the top down.
		err = me.WaitForConsole(me.BoxName+" : Stopping", 30)
		if err != nil {
			return err
		}
	}

	return nil
}

var runner virtualbox.Runner

// This is here because it's not implemented in libretto.
func (me *Box) ReplacementBoxHalt() error {

	_, err := me.RunCombinedError("controlvm", me.BoxName, "acpipowerbutton")
	if err != nil {
		return lvm.WrapErrors(lvm.ErrStoppingVM, err)
	}
	return nil
}

// Run runs a VBoxManage command.
func (me *Box) Run(args ...string) (string, string, error) {
	var vboxManagePath string
	// If vBoxManage is not found in the system path, fall back to the
	// hard coded path.
	if path, err := exec.LookPath("VBoxManage"); err == nil {
		vboxManagePath = path
	} else {
		vboxManagePath = virtualbox.VBOXMANAGE
	}
	cmd := exec.Command(vboxManagePath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// RunCombinedError runs a VBoxManage command.  The output is stdout and the the
// combined err/stderr from the command.
func (me *Box) RunCombinedError(args ...string) (string, error) {
	wout, werr, err := me.Run(args...)
	if err != nil {
		if werr != "" {
			return wout, fmt.Errorf("%s: %s", err, werr)
		}
		return wout, err
	}

	return wout, nil
}

func (me *Box) RestartBox() error {

	var err error

	if me == nil {
		// Throw software error.
		return nil
	}

	err = me.GetBoxStatus()
	if err != nil {
		return err
	}

	switch me.Status {
	case BoxGearBoxOK:
		fallthrough
	case BoxGearBoxNOK:
		fallthrough
	case BoxRunning:
		fallthrough
	case BoxStarted:
		err := me.StopBox()
		if err != nil {
			return err
		}
		err = me.StartBox()
		if err != nil {
			return err
		}

	case BoxHalted:
		err := me.StartBox()
		if err != nil {
			return err
		}

	case BoxUnknown:

	}

	if me.Status != BoxRunning {
		// Throw an error.
	}

	return err
}

func (me *Box) GetBoxStatus() error {

	if me == nil {
		// Throw software error.
		me.Status = BoxUnknown
		return nil
	}

	state, err := me.Instance.GetState()
	if err != nil {
		me.Status = BoxError
	} else {
		me.Status = state
	}

	return err
}

func (me *Box) GetApiStatus() error {

	if me == nil {
		// Throw software error.
		return nil
	}

	if me.Status == BoxRunning {
		err := me.WaitForConsole("", 10)
		if err != nil {
			return err
		}
	}

	return nil
}

func (me *Box) CreateBox() error {

	if me == nil {
		// Throw software error.
		me.Status = BoxUnknown
		return nil
	}

	// Check if the VM is already there.
	state, err := me.Instance.GetState()
	if err != nil {
		// Doesn't exist - great!
		if _, err := os.Stat(me.OvaFile); os.IsNotExist(err) {
			fmt.Printf("\r👎 %s: VM OVA file does not exist in %s.\n", me.BoxName, me.OvaFile)
			return err
		}

		err = me.Instance.Provision()
		if err != nil {
			me.Status = BoxError
		}
	} else {
		// Already created!
		fmt.Printf("\r👎 %s: Cannot create. VM already exists and is in state %s.\n", me.BoxName, state)
	}

	return err
}

func (me *Box) ValidateBox() error {

	if me == nil {
		// Throw software error.
		me.Status = BoxUnknown
		return nil
	}
	/*
		state, err := me.Instance.Provision()
		if err != nil {
			me.Status = BoxError
		} else {
			me.Status = state
		}
		return err
	*/
	return nil
}
