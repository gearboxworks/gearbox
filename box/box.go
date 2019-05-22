package box

import (
	"bufio"
	"fmt"
	"gearbox/global"
	"gearbox/help"
	"gearbox/only"
	"gearbox/os_support"
	"gearbox/ssh"
	"gearbox/util"
	lbssh "github.com/apcera/libretto/ssh"
	"github.com/apcera/libretto/virtualmachine/virtualbox"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	// dmvb "github.com/docker/machine/drivers/virtualbox"
	// "github.com/docker/machine/libmachine/drivers/plugin"
	"net"
	"regexp"
	"strings"
	"time"
)


type Box struct {
	Boxname    string
	VmInstance virtualbox.VM
	State      BoxState
	OvaFile    string
	VmBaseDir  string
	VmIsoFile  string

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

	if _args.VmBaseDir == "" {
		_args.VmBaseDir = string(OsSupport.GetUserConfigDir() + "/box/vm")
	}

	if _args.VmIsoFile == "" {
		_args.VmIsoFile = string(OsSupport.GetUserConfigDir() + "/box/iso/gearbox-0.5.0.iso")
	}

	_args.VmInstance = virtualbox.VM{
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

//		_, err := os.Stat(me.OvaFile)
//		if os.IsExist(err) {
//			break
//		}
//		err = vm.RestoreAssets(string(cfgdir), strings.TrimLeft(OvaFileName, string(os.PathSeparator)))
//		if err != nil {
//			sts = status.Wrap(err, &status.Args{
//				Message: fmt.Sprintf("%s: VM OVA file cannot be created as'%s'.", global.Brandname, me.OvaFile),
//			})
//			break
//		}
	}
	return sts
}


func (me *Box) WaitForVmState(displayString string) bool {

	found := false
	var waitCount int

	spinner := util.NewSpinner(util.SpinnerArgs{
		Text:    displayString,
		ExitOK:  displayString + " - OK",
		ExitNOK: displayString + " - FAILED",
	})
	spinner.Start()

	for waitCount = 0; waitCount < me.WaitRetries; waitCount++ {

		_, sts := me.GetState()
		if is.Error(sts) {
			found = false
			break
		}

		if me.State.VM.CurrentState == me.State.VM.WantState {
			found = true
			break
		}

		time.Sleep(me.WaitDelay)
		spinner.Update(fmt.Sprintf("%s [%d]", displayString, waitCount))
	}

	spinner.Stop(found)

	return found
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


func (me *Box) Start() (sts status.Status) {

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		_, sts := me.GetState()
		if is.Error(sts) {
			break
		}

		switch {
			case me.State.VM.CurrentState == VmStateUnknown:
				sts = me.State.LastSts
				break

			case me.State.VM.CurrentState == VmStateStarting:
				sts = me.State.LastSts
				break

			case me.State.VM.CurrentState == VmStateRunning:
				sts = me.State.LastSts
				break

			case me.State.VM.CurrentState == VmStateStopping:
				sts = me.State.LastSts
				break

			case me.State.VM.CurrentState == VmStatePowerOff:
				// fall-through
		}

		me.State.VM.WantState = VmStateRunning
		me.State.API.WantState = VmStateRunning
		err := me.VmInstance.Start()
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("%s VM failed to start", global.Brandname),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    ErrorState,
			})
			break
		}

		if me.NoWait != false {
			break
		}

		if me.WaitForVmState(fmt.Sprintf("%s VM: Starting", global.Brandname)) == true {
			err = me.GetApiStatus(fmt.Sprintf("%s API: Starting", global.Brandname), 30)
		}
	}

	return sts
}


func (me *Box) Stop() (sts status.Status) {

	var err error

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		_, sts := me.GetState()
		if is.Error(sts) {
			break
		}

		switch {
			case me.State.VM.CurrentState == VmStateUnknown:
				sts = me.State.LastSts
				break

			case me.State.VM.CurrentState == VmStateStarting:
				sts = me.State.LastSts
				break

			case me.State.VM.CurrentState == VmStateRunning:
				// fall-through

			case me.State.VM.CurrentState == VmStateStopping:
				sts = me.State.LastSts
				break

			case me.State.VM.CurrentState == VmStatePowerOff:
				sts = me.State.LastSts
				break
		}

		me.State.VM.WantState = VmStatePowerOff
		me.State.API.WantState = VmStatePowerOff
		err = me.VmInstance.Halt()
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("%s VM failed to stop", global.Brandname),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    ErrorState,
			})
			break
		}

		if me.NoWait != false {
			break
		}

		if me.WaitForVmState(fmt.Sprintf("%s VM: Stopping", global.Brandname)) == true {
			break
		}
	}

	if err == nil {
		sts = status.Success("%s VM stopped", global.Brandname)
	}

	return sts
}


func (me *Box) Restart() (sts status.Status) {

	for range only.Once {

		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		me.Stop()
		if me.State.VM.CurrentState != me.State.VM.WantState {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("%s VM in an unknown state: %s", global.Brandname, me.State),
				Data:    VmStateUnknown,
			})
			break
		}

		me.Start()
		if me.State.VM.CurrentState != me.State.VM.WantState {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("%s VM in an unknown state: %s", global.Brandname, me.State),
				Data:    VmStateUnknown,
			})
			break
		}
	}

	if me.State.VM.CurrentState == me.State.VM.WantState {
		sts = status.Success("%s VM restarted OK", global.Brandname)
	}

	return sts
}


func (me *Box) GetCachedState() (state BoxState, sts status.Status) {

	// This is required so that not more than one

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		state = me.State
	}

	return
}


func (me *Box) GetState() (BoxState, status.Status) {

	// Possible VM states:
	// running
	// paused
	// saved
	// poweroff

	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		sts = me.GetVmStatus()
		if is.Error(sts) {
			break
		}
		if me.State.VM.CurrentState != VmStateRunning {
			break
		}

		sts = me.GetApiStatus("", 10)
		if is.Error(sts) {
			break
		}

		vmState := me.State.GetStateMeaning()
		sts = vmState.Sts
		me.State.LastSts = sts
		//fmt.Printf("FOO:", vmState)

		//		if err == nil {
//			sts = status.Success("%s VM in a valid state: %s", global.Brandname, state)
//			sts.SetData(state)
//
//		}

		//sts = me.State.LastSts
		//fmt.Printf("STATE2: %v\n", sts)
	}

	return me.State, sts
}


func (me *Box) GetVmStatus() (status.Status) {

	// Possible VM states:
	// running
	// paused
	// saved
	// poweroff

	var sts status.Status
	// var kvm KeyValueMap

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		if me.State.VM.WantState == VmStateInit {
			me.State.VM.WantState = me.State.VM.CurrentState
		}

		_, sts = me.cmdListVm()
		// fmt.Printf("Box '%s' is in state: '%s'\n", kvm["name"], kvm["VMState"])
		if is.Error(sts) {
			me.State.VM.CurrentState = VmStateNotPresent
		}

/*
		// First check on the VM.
		// state, err := me.VmInstance.GetState()
		switch kvm["VMState"] {
			case VmStateRunning:
				me.State.VM.CurrentState = VmStateRunning

			case VmStatePowerOff:
				fallthrough
			case VmStateSaved:
				fallthrough
			case VmStatePaused:
				me.State.VM.CurrentState = VmStatePowerOff
				me.State.API.CurrentState = VmStatePowerOff
		}
*/

/*
		switch {
			case err != nil:
				// No Gearbox VM available - need to create one.
				me.State.VM.CurrentState = VmStateUnknown
				me.State.LastSts = status.Fail(&status.Args{
					Message: fmt.Sprintf("%s VM - needs to be created", global.Brandname),
					Help:    help.ContactSupportHelp(), // @TODO need better support here
					Data:    me.State.VM.CurrentState,
				})
				break

			case (me.State.VM.CurrentState == me.State.VM.WantState) && (state == VmStatePowerOff):
				// If we are not changing states and the VM is halted.
				me.State.VM.CurrentState = VmStatePowerOff
				me.State.LastSts = status.Success("%s VM - halted", global.Brandname)
				break

			case (me.State.VM.CurrentState == me.State.VM.WantState) && (state == VmStateRunning):
				// If we are not changing states and the VM is running.
				me.State.VM.CurrentState = VmStateRunning
				me.State.LastSts = status.Success("%s VM - running", global.Brandname)
				// Don't break here - need to check on the API.

			case (me.State.VM.CurrentState != me.State.VM.WantState) && (state == VmStatePowerOff):
				// If we are changing states then the VM is halting.
				me.State.VM.CurrentState = VmStateStopping
				me.State.LastSts = status.Success("%s VM - stopping", global.Brandname)
				// Don't break here - need to check on the API.

			case (me.State.VM.CurrentState != me.State.VM.WantState) && (state == VmStateRunning):
				// If we are changing states then the VM is starting.
				me.State.VM.CurrentState = VmStateStarting
				me.State.LastSts = status.Success("%s VM - starting", global.Brandname)
				// Don't break here - need to check on the API.
		}
*/

		fmt.Printf("vmState: %v\n", sts)
	}

	return sts
}


// We have to have some way to block access to other concurrent processes/threads
// So, we're simply establishing a boolean that indicates this fact.
// var alreadyRunning = false
// @TODO - OK, so that's not working out.
func (me *Box) GetApiStatus(displayString string, waitFor time.Duration) status.Status {

	for range only.Once {
		sts := EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		if me.State.API.WantState == VmStateInit {
			me.State.API.WantState = me.State.API.CurrentState
		}

		spinner := newSpinner(displayString)
		displaySpinner := !me.ShowConsole && displayString != ""

		if displaySpinner {
			// We want to display just a spinner instead of console output.
			spinner.Start()
		}

		// Connect to this console
		conn, err := net.Dial("tcp", me.ConsoleHost+":"+me.ConsolePort)
		if err != nil {
			me.State.API.CurrentState = VmStatePowerOff
			me.State.LastSts = status.Fail(&status.Args{
				Message: fmt.Sprintf("%s API - timeout", global.Brandname),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    me.State.API.CurrentState,
			})
			break
		}
		// defer closeDialConnection(conn)
		defer conn.Close()

		// Set default state before we begin.
		me.State.API.CurrentState = VmStateUnknown
		me.State.LastSts = status.Fail(&status.Args{
			Message: fmt.Sprintf("%s API - no data", global.Brandname),
			Help:    help.ContactSupportHelp(), // @TODO need better support here
			Data:    me.State.API.CurrentState,
		})

		exitWhen := time.Now().Add(time.Second * waitFor)
		readBuffer := make([]byte, 512)
		for waitCount := 0; time.Now().Unix() < exitWhen.Unix(); waitCount++ {
			err = conn.SetDeadline(time.Now().Add(me.ConsoleReadWait))
			if err != nil {
				me.State.API.CurrentState = VmStateUnknown
				me.State.LastSts = status.Fail(&status.Args{
					Message: fmt.Sprintf("%s API - deadline", global.Brandname),
					Help:    help.ContactSupportHelp(), // @TODO need better support here
					Data:    me.State.API.CurrentState,
				})
				break
			}

			bytesRead, err := bufio.NewReader(conn).Read(readBuffer)
			// bytesRead, err := conn.Read(readBuffer)
			// readBuffer, err := bufio.NewReader(conn).ReadString('\n')
			// bytesRead := len(readBuffer)
			if err != nil {
				me.State.API.CurrentState = VmStateUnknown
				me.State.LastSts = status.Fail(&status.Args{
					Message: fmt.Sprintf("%s API - no data", global.Brandname),
					Help:    help.ContactSupportHelp(), // @TODO need better support here
					Data:    me.State.API.CurrentState,
				})
				break
			}

			if bytesRead > 0 {
				if me.ShowConsole {
					fmt.Printf("%s", string(readBuffer[:bytesRead]))
				}

				sts = me.heartbeatOk(readBuffer, bytesRead)
				if sts != nil {
					me.State.API.CurrentState = VmStateRunning
					me.State.LastSts = status.Success("%s API - running", global.Brandname)
					break

				} else {
					me.State.API.CurrentState = VmStateStarting
					me.State.LastSts = status.Success("%s API - starting", global.Brandname)
					// Do not break.
				}
			}

			time.Sleep(me.WaitDelay)
			if displaySpinner {
				spinner.Update(fmt.Sprintf("%s [%d]", displayString, waitCount))
			}
		}

		if me.ShowConsole {
			fmt.Printf("\n\n# Exiting Console.\n")
		}

		if displaySpinner {
			spinner.Stop(false)
		}
	}

	fmt.Printf("vmState: %v\n", me.State.LastSts)

	return me.State.LastSts
}


func EnsureNotNil(bx *Box) (sts status.Status) {
	if bx == nil {
		sts = status.Fail(&status.Args{
			Message: "unexpected error",
			Help:    help.ContactSupportHelp(), // @TODO need better support here
			Data:    VmStateUnknown,
		})
	}

	return sts
}
