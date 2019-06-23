package box

import (
	"fmt"
	"gearbox/global"
	"gearbox/help"
	"github.com/gearboxworks/go-osbridge"
	"github.com/gearboxworks/go-status/only"

	//	"gearbox/os_support"
	"gearbox/ssh"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"

	"time"
)

type Box struct {
	Boxname      string
	State        BoxState
	VmBaseDir    string
	VmIsoDir     string
	VmIsoVersion string
	VmIsoFile    string
	VmIsoUrl     string
	VmIsoDlIndex int

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

	OsSupport osbridge.OsBridger
}
type Args Box

// These will be re-implemented using the EventBroker framework.

func NewBox(OsSupport osbridge.OsBridger, args ...Args) *Box {

	var _args Args
	//var sts status.Status
	box := &Box{}

	for range only.Once {

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

		if _args.VmIsoDir == "" {
			_args.VmIsoDir = string(OsSupport.GetUserConfigDir() + "/box/iso")
		}

		_args.VmIsoDlIndex = 100

		*box = Box(_args)

	}

	return box
}

func (me *Box) Initialize() (sts status.Status) {

	return sts
}

//
//
//func (me *Box) WaitForVmState(displayString string) bool {
//
//	found := false
//	var waitCount int
//
//	spinner := util.NewSpinner(util.SpinnerArgs{
//		Text:    displayString,
//		ExitOK:  displayString + " - OK",
//		ExitNOK: displayString + " - FAILED",
//	})
//	spinner.Start()
//
//	for waitCount = 0; waitCount < me.WaitRetries; waitCount++ {
//
//		_, sts := me.GetState()
//		if is.Error(sts) {
//			found = false
//			break
//		}
//
//		if me.State.VM.CurrentState == me.State.VM.WantState {
//			found = true
//			break
//		}
//
//		time.Sleep(me.WaitDelay)
//		spinner.Update(fmt.Sprintf("%s [%d]", displayString, waitCount))
//	}
//
//	spinner.Stop(found)
//
//	return found
//}
//
//
//func closeDialConnection(conn net.Conn) {
//	_ = conn.Close()
//}
//
//
//func newSpinner(displayString string) *util.Spinner {
//	return util.NewSpinner(util.SpinnerArgs{
//		Text:    displayString,
//		ExitOK:  displayString + " - OK",
//		ExitNOK: displayString + " - FAILED",
//	})
//}
//
//
//func (me *Box) heartbeatOk(b []byte, n int) (sts status.Status) {
//
//	for range only.Once {
//		apiSplit := strings.Split(string(b[:n]), ";")
//		if len(apiSplit) <= 1 {
//			break
//		}
//
//		match, _ := regexp.MatchString(me.ConsoleOkString, apiSplit[1])
//		if !match {
//			break
//		}
//
//		fmt.Printf("API:%v\n", apiSplit)
//		if len(apiSplit) < 2 {
//			sts = status.Fail(&status.Args{
//				Message: fmt.Sprintf("did not receive 'OK' from console: %s",
//					apiSplit[2],
//				),
//				Data: NotOkState,
//			})
//			break
//		}
//
//		if apiSplit[2] != "OK" {
//			sts = status.Fail(&status.Args{
//				Message: fmt.Sprintf("did not receive 'OK' from console: %s",
//					apiSplit[2],
//				),
//				Data: NotOkState,
//			})
//			break
//		}
//
//		sts = status.Success("received 'OK' from console")
//		sts.SetData(OkState)
//	}
//
//	return sts
//}

func (me *Box) Start() (sts status.Status) {

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		fmt.Printf("Not implemented.\n")
	}

	return sts
}

func (me *Box) Stop() (sts status.Status) {

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		fmt.Printf("Not implemented.\n")
	}

	return sts
}

func (me *Box) Restart() (sts status.Status) {

	for range only.Once {

		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		fmt.Printf("Not implemented.\n")
	}

	return sts
}

func (me *Box) GetState() (BoxState, status.Status) {

	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		fmt.Printf("Not implemented.\n")
	}

	return me.State, sts
}

func (me *Box) CreateBox() (BoxState, status.Status) {

	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		fmt.Printf("Not implemented.\n")
	}

	return me.State, sts
}

func (me *Box) GetVmStatus() status.Status {

	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		fmt.Printf("Not implemented.\n")
	}

	return sts
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
