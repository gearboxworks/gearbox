package box

import (
	"bufio"
	"bytes"
	"fmt"
	"gearbox/box/vm"
	"gearbox/global"
	"gearbox/help"
	"gearbox/only"
	"gearbox/os_support"
	"gearbox/ssh"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/util"
	lbssh "github.com/apcera/libretto/ssh"
	"github.com/apcera/libretto/virtualmachine"
	"github.com/apcera/libretto/virtualmachine/virtualbox"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type State string

type Box struct {
	Boxname  string
	Instance virtualbox.VM
	State    State
	OvaFile  string

	// SSH related - Need to fix this. It's used within CreateBox()
	SshUsername  string
	SshPassword  string
	SshPublicKey string

	// State polling delays.
	NoWait      bool
	WaitDelay   time.Duration
	WaitRetries int

	// Console related.
	ConsoleHost     string
	ConsolePort     string
	ConsoleOkString string
	ConsoleReadWait time.Duration
	ShowConsole     bool

	OsSupport oss.OsSupporter
}
type Args Box

func NewBox(OsSupport oss.OsSupporter, args ...Args) *Box {
	var _args Args
	if len(args) > 0 {
		_args = args[0]
	}

	_args.OsSupport = OsSupport

	if _args.Boxname == "" {
		_args.Boxname = global.Brandname
	}

	if _args.WaitDelay == 0 {
		_args.WaitDelay = DefaultWaitDelay
	}

	if _args.WaitRetries == 0 {
		_args.WaitRetries = DefaultWaitRetries
	}

	if _args.ConsoleHost == "" {
		_args.ConsoleHost = DefaultConsoleHost
	}

	if _args.ConsolePort == "" {
		_args.ConsolePort = DefaultConsolePort
	}

	if _args.ConsoleOkString == "" {
		_args.ConsoleOkString = DefaultConsoleOkString
	}

	if _args.ConsoleReadWait == 0 {
		_args.ConsoleReadWait = DefaultConsoleReadWait
	}

	if _args.SshUsername == "" {
		_args.SshUsername = ssh.DefaultUsername
	}

	if _args.SshPassword == "" {
		_args.SshPassword = ssh.DefaultPassword
	}

	if _args.SshPublicKey == "" {
		_args.SshPublicKey = ssh.DefaultKeyFile
	}

	_args.Instance = virtualbox.VM{
		Name: _args.Boxname,
		Src:  _args.OvaFile,
		Credentials: lbssh.Credentials{
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

func (me *Box) Initialize() (sts status.Status) {
	for range only.Once {
		if me.OvaFile != "" {
			break
		}

		cfgdir := me.OsSupport.GetUserConfigDir()

		me.OvaFile = fmt.Sprintf("%s/%s", cfgdir, OvaFileName)

		// The OvaFile is created from an export from within VirtualBox.
		// VBoxManage export Parent -o Parent.ova --options manifest
		// This was the best way to create a base template, avoiding too much code bloat.
		// And allows multiple VM frameworks to be used with libretto.
		// It doesn't include the ISO image yet as it is too large.
		// Once the ISO image size has been reduced, we can do this:
		// VBoxManage export Parent -o Parent.ova --options iso,manifest

		_, err := os.Stat(me.OvaFile)
		if os.IsExist(err) {
			break
		}
		err = vm.RestoreAssets(string(cfgdir), strings.TrimLeft(OvaFileName, string(os.PathSeparator)))
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("%s: VM OVA file cannot be created as'%s'.", global.Brandname, me.OvaFile),
			})
			break
		}
	}
	return sts
}

func (me *Box) WaitForState(waitForState State, displayString string) (sts status.Status) {

	for range only.Once {
		var waitCount int

		spinner := util.NewSpinner(util.SpinnerArgs{
			Text:    displayString,
			ExitOK:  displayString + " - OK",
			ExitNOK: displayString + " - FAILED",
		})
		spinner.Start()

		for waitCount = 0; waitCount < me.WaitRetries; waitCount++ {
			sts = me.GetState()
			if is.Error(sts) {
				spinner.Stop(false)
				break
			}
			if sts.GetData().(State) == waitForState {
				spinner.Stop(true)
				break
			}
			time.Sleep(me.WaitDelay)
			spinner.Update(fmt.Sprintf("%s [%d]", displayString, waitCount))
		}

	}
	if is.Error(sts) {
		sts = status.Wrap(sts, &status.Args{
			Message: fmt.Sprintf("%s VM failed to stop", global.Brandname),
			Help:    help.ContactSupportHelp(), // @TODO need better support here
			Data:    ErrorState,
		})
	}
	return sts
}

func closeDialConnection(conn net.Conn) {
	_ = conn.Close()
}

func newSpinner(displayString string) *util.Spinner {
	return util.NewSpinner(util.SpinnerArgs{
		Text:    displayString,
		ExitOK:  displayString + " - OK",
		ExitNOK: displayString + " - FAILED",
	})
}

func (me *Box) heartbeatOk(b []byte, n int) (sts status.Status) {
	for range only.Once {
		apiSplit := strings.Split(string(b[:n]), ";")
		if len(apiSplit) <= 1 {
			break
		}
		match, _ := regexp.MatchString(me.ConsoleOkString, apiSplit[1])
		if !match {
			break
		}
		if apiSplit[2] != "OK" {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("did not receive 'OK' from console: %s",
					apiSplit[2],
				),
				Data: NotOkState,
			})
			break
		}
		sts = status.Success("received 'OK' from console")
		sts.SetData(OkState)
	}
	return sts
}

func (me *Box) WaitForConsole(displayString string, waitFor time.Duration) (sts status.Status) {

	spinner := newSpinner(displayString)
	displaySpinner := !me.ShowConsole && displayString != ""

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}
		sts = me.GetState()
		if is.Error(sts) {
			break
		}

		state, sts := sts.GetString()
		if is.Error(sts) {
			break
		}
		if state != RunningState {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("unable to wait for console with %s VM in '%s' state",
					global.Brandname,
					state,
				),
			})
			break
		}

		if displaySpinner {
			// We want to display just a spinner instead of console output.
			spinner.Start()
		}

		// Connect to this console
		conn, err := net.Dial("tcp", me.ConsoleHost+":"+me.ConsolePort)
		defer closeDialConnection(conn)
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: "unable to connect to console",
			})
			break
		}

		exitWhen := time.Now().Add(time.Second * waitFor)

		readBuffer := make([]byte, 512)
		for waitCount := 0; time.Now().Unix() < exitWhen.Unix(); waitCount++ {
			err = conn.SetDeadline(time.Now().Add(me.ConsoleReadWait))
			if err != nil {
				sts = status.Wrap(err, &status.Args{
					Message: "unable to set deadline while waiting for console connection",
				})
				break
			}

			bytesRead, err := bufio.NewReader(conn).Read(readBuffer)
			// bytesRead, err := conn.Read(readBuffer)
			// readBuffer, err := bufio.NewReader(conn).ReadString('\n')
			// bytesRead := len(readBuffer)
			if err != nil {
				sts = status.Wrap(err, &status.Args{
					Message: "unable to set read from connection while waiting for console",
					Data:    NotOkState,
				})
				break
			}

			if bytesRead > 0 {
				if me.ShowConsole {
					fmt.Printf("%s", string(readBuffer[:bytesRead]))
				}
				sts = me.heartbeatOk(readBuffer, bytesRead)
				if sts != nil {
					break
				}
			}
			time.Sleep(me.WaitDelay)
			if displaySpinner {
				spinner.Update(fmt.Sprintf("%s [%d]", displayString, waitCount))
			}
		}

	}

	if me.ShowConsole {
		fmt.Printf("\n\n# Exiting Console.\n")
	}

	if displaySpinner {
		spinner.Stop(false)
	}

	return sts
}

func (me *Box) StartBox() (sts status.Status) {
	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		sts := me.GetState()
		if is.Error(sts) {
			break
		}

		if me.State == StartedState {
			sts = status.Success("%s VM is starting", global.Brandname)
		}

		if me.State == RunningState {
			sts = status.Success("%s VM is running", global.Brandname)
		}

		err := me.Instance.Start()
		if err == nil {
			if me.NoWait {
				break
			}
			err = me.WaitForState(RunningState, fmt.Sprintf("%s VM: Starting", global.Brandname))
			if err == nil {
				err = me.WaitForConsole(fmt.Sprintf("%s VM: Starting", global.Brandname), 30)
			}
		}

		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("%s VM failed to start", global.Brandname),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    ErrorState,
			})
			break
		}

	}

	return sts
}

func (me *Box) StopBox() (sts status.Status) {

	var err error
	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		sts := me.GetState()
		if is.Error(sts) {
			break
		}

		if me.State == HaltedState {
			break
		}

		err = me.Instance.Halt()
		if err != nil {
			break
		}

		if me.NoWait != false {
			break
		}

		sts = me.WaitForState(HaltedState, fmt.Sprintf("%s VM stopping", global.Brandname))
		if is.Error(sts) {
			break
		}

		err = me.WaitForConsole(fmt.Sprintf("%s VM stopping", global.Brandname), 30)

	}
	if err == nil {
		sts = status.Success("%s VM stopped", global.Brandname)
	} else {
		sts = status.Wrap(err, &status.Args{
			Message: fmt.Sprintf("%s VM failed to stop", global.Brandname),
			Help:    help.ContactSupportHelp(), // @TODO need better support here
			Data:    ErrorState,
		})
	}
	return sts
}

//var runner virtualbox.Runner

// This is here because it's not implemented in libretto.
func (me *Box) ReplacementBoxHalt() error {

	_, err := me.RunCombinedError("controlvm", global.Brandname, "acpipowerbutton")
	if err != nil {
		return virtualmachine.WrapErrors(virtualmachine.ErrStoppingVM, err)
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

func (me *Box) RestartBox() (sts status.Status) {

	for range only.Once {

		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		sts = me.GetState()
		if is.Error(sts) {
			break
		}

		switch me.State {
		case OkState:
			fallthrough
		case NotOkState:
			fallthrough
		case RunningState:
			fallthrough
		case StartedState:
			sts = me.StopBox()
			if is.Error(sts) {
				break
			}
			sts = me.StartBox()
			if is.Error(sts) {
				break
			}

		case HaltedState:
			sts = me.StartBox()
			if is.Error(sts) {
				break
			}

		case UnknownState:
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("%s VM in an unknown state: %s", global.Brandname, me.State),
				Data:    UnknownState,
			})
		}

	}
	if me.State != RunningState {
		sts = status.Wrap(sts, &status.Args{
			Message: fmt.Sprintf("%s VM failed to restart", global.Brandname),
		})
	}
	return sts
}

func (me *Box) GetState() (sts status.Status) {
	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}
		state, err := me.Instance.GetState()
		if err != nil {
			sts := status.Success("%s VM in a valid state: %s", global.Brandname, state)
			sts.SetData(state)
		} else {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("%s VM in an invalid state", global.Brandname),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    state,
			})
			break
		}
		if state == RunningState {
			err := me.WaitForConsole("", 10)
			if err != nil {
				sts = status.Wrap(err, &status.Args{
					Message: fmt.Sprintf("%s VM's API failed to respond", global.Brandname),
					Help:    help.ContactSupportHelp(), // @TODO need better support here
					Data:    state,
				})
			}
		}
	}
	return sts
}

func (me *Box) CreateBox() (sts status.Status) {

	for range only.Once {
		if me == nil {
			sts = status.Fail(&status.Args{
				Message: "unexpected failure",
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    UnknownState,
			})
			break
		}

		state, err := me.Instance.GetState()
		if err == nil {
			sts = status.Success("%s VM already exists and is in state %s.\n", global.Brandname, state)
			break
		}

		if _, err := os.Stat(me.OvaFile); os.IsNotExist(err) {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("%s VM OVA file '%s' does not exist", global.Brandname, me.OvaFile),
				Data:    UnknownState,
			})
			break
		}

		err = me.Instance.Provision()
		if err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("failed to provision %s VM", global.Brandname),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    UnknownState,
			})
			break
		}
	}

	return sts
}

func EnsureNotNil(bx *Box) (sts status.Status) {
	if bx == nil {
		sts = status.Fail(&status.Args{
			Message: "unexpected error",
			Help:    help.ContactSupportHelp(), // @TODO need better support here
			Data:    UnknownState,
		})
	}
	return sts
}

func GetStateMeaning(state State) string {
	m, _ := StateMeaning[state]
	return m
}
