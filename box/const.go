package box

import (
	"fmt"
	"gearbox/global"
	"github.com/gearboxworks/go-status"
	"time"
)

const (
	Basedir     = "/home/gearbox/projects"
	OvaFileName = "/box/vm/Gearbox.ova"
)

const (
	VmStateInit			= "init"
	VmStateNotPresent	= "not present"
	VmStateUnknown		= "unknown"
	VmStatePowerOff 	= "poweroff"	// Valid VM state return from listvm
	VmStatePaused 		= "paused"		// Valid VM state return from listvm
	VmStateSaved 		= "saved"		// Valid VM state return from listvm
	VmStateRunning  	= "running"		// Valid VM state return from listvm
	VmStateStarting		= "starting"
	VmStateStarting2	= "starting up"
	VmStateStopping		= "stopping"
	VmStateStopping2	= "shutting down"
	VmStateDontCare		= ""
	VmStateIgnore		= "ignore"

	ErrorState   = "error"
	VmStateHalted  = "halted"
	OkState      = "ok"
	NotOkState   = "nok"
)

const (
	StateInit		= 0
	StateNotPresent	= 1
	StateUnknown	= 2
	StateDown		= 3
	StateStarting	= 4
	StateUp			= 5
	StateStopping	= 6
	StateDontCare	= 7
)

const (
	DefaultLogo = "heartbeat/img/IconLogo.ico"
	DefaultUp = "heartbeat/img/UpArrow.ico"
	DefaultDown = "heartbeat/img/DownArrow.ico"

	IconLogo = "heartbeat/img/IconLogo.ico"
	IconError = "heartbeat/img/IconError.ico"
	IconWarning = "heartbeat/img/IconWarning.ico"
	IconUp = "heartbeat/img/IconUp.ico"
	IconDown = "heartbeat/img/IconDown.ico"
	IconStarting = "heartbeat/img/IconStarting.ico"
	IconStopping = "heartbeat/img/IconStopping.ico"
)


type BoxEntity struct {
	Name         string
	CurrentState string
	WantState    string
}

type BoxState struct {
	VM      BoxEntity
	API     BoxEntity
	LastSts status.Status
}

type VmState struct {
	Name			string
	VmState			string
	VmWantState		string
	VmIconState		string
	VmTitleState	string
	VmHintState		string

	ApiState		string
	ApiWantState	string
	ApiIconState	string
	ApiTitleState	string
	ApiHintState	string

	Title			string
	Hint			string
	Sts	    		status.Status
}


var DecodeVmState = map[string]VmState{
	VmStateInit: {
		Name:			VmStateInit,
		VmState:		VmStateInit,
		VmWantState:	VmStateDontCare,
		ApiState:		VmStateDontCare,
		ApiWantState:	VmStateDontCare,

		VmTitleState:	fmt.Sprintf("%s VM: idle", global.Brandname),
		VmHintState:	fmt.Sprintf("%s VM: idle", global.Brandname),
		VmIconState:	IconLogo,

		ApiTitleState:	fmt.Sprintf("%s API: halted", global.Brandname),
		ApiHintState:	fmt.Sprintf("%s API: halted", global.Brandname),
		ApiIconState:	IconLogo,

		Title:			fmt.Sprintf("%s is idle", global.Brandname),
		Hint:			fmt.Sprintf("%s is idle.", global.Brandname),
	},

	VmStateNotPresent: {
		Name:			VmStateNotPresent,
		VmState:		VmStateNotPresent,
		VmWantState:	VmStateDontCare,
		ApiState:		VmStateDontCare,
		ApiWantState:	VmStateDontCare,

		VmTitleState:	fmt.Sprintf("%s VM: is not present", global.Brandname),
		VmHintState:	fmt.Sprintf("%s VM: is not present", global.Brandname),
		VmIconState:	IconWarning,

		ApiTitleState:	fmt.Sprintf("%s API: halted", global.Brandname),
		ApiHintState:	fmt.Sprintf("%s API: halted", global.Brandname),
		ApiIconState:	IconLogo,

		Title:			fmt.Sprintf("%s VM creation", global.Brandname),
		Hint:			fmt.Sprintf("%s VM needs to be created.", global.Brandname),
	},

	VmStateUnknown: {
		Name:			VmStateUnknown,
		VmState:		VmStateUnknown,
		VmWantState:	VmStateDontCare,
		ApiState:		VmStateDontCare,
		ApiWantState:	VmStateDontCare,

		VmTitleState:	fmt.Sprintf("%s VM: unknown", global.Brandname),
		VmHintState:	fmt.Sprintf("%s VM: unknown", global.Brandname),
		VmIconState:	IconWarning,

		ApiTitleState:	fmt.Sprintf("%s API: halted", global.Brandname),
		ApiHintState:	fmt.Sprintf("%s API: halted", global.Brandname),
		ApiIconState:	IconLogo,

		Title:			fmt.Sprintf("%s VM is unknown", global.Brandname),
		Hint:			fmt.Sprintf("%s VM is in an unknown state.", global.Brandname),
	},

	VmStatePowerOff: {
		Name:			VmStatePowerOff,
		VmState:		VmStatePowerOff,
		VmWantState:	VmStateDontCare,
		ApiState:		VmStateDontCare,
		ApiWantState:	VmStateDontCare,

		VmTitleState:	fmt.Sprintf("%s VM: halted", global.Brandname),
		VmHintState:	fmt.Sprintf("%s VM: halted", global.Brandname),
		VmIconState:	IconDown,

		ApiTitleState:	fmt.Sprintf("%s API: halted", global.Brandname),
		ApiHintState:	fmt.Sprintf("%s API: halted", global.Brandname),
		ApiIconState:	IconLogo,

		Title:			fmt.Sprintf("%s VM halted", global.Brandname),
		Hint:			fmt.Sprintf("%s VM is in a halted state.", global.Brandname),
	},

	VmStatePaused: {
		Name:			VmStatePaused,
		VmState:		VmStatePaused,
		VmWantState:	VmStateDontCare,
		ApiState:		VmStateDontCare,
		ApiWantState:	VmStateDontCare,

		VmTitleState:	fmt.Sprintf("%s VM: paused", global.Brandname),
		VmHintState:	fmt.Sprintf("%s VM: paused", global.Brandname),
		VmIconState:	IconDown,

		ApiTitleState:	fmt.Sprintf("%s API: halted", global.Brandname),
		ApiHintState:	fmt.Sprintf("%s API: halted", global.Brandname),
		ApiIconState:	IconLogo,

		Title:			fmt.Sprintf("%s VM paused", global.Brandname),
		Hint:			fmt.Sprintf("%s VM is in a paused state.", global.Brandname),
	},

	VmStateSaved: {
		Name:			VmStateSaved,
		VmState:		VmStateSaved,
		VmWantState:	VmStateDontCare,
		ApiState:		VmStateDontCare,
		ApiWantState:	VmStateDontCare,

		VmTitleState:	fmt.Sprintf("%s VM: saved", global.Brandname),
		VmHintState:	fmt.Sprintf("%s VM: saved", global.Brandname),
		VmIconState:	IconDown,

		ApiTitleState:	fmt.Sprintf("%s API: halted", global.Brandname),
		ApiHintState:	fmt.Sprintf("%s API: halted", global.Brandname),
		ApiIconState:	IconLogo,

		Title:			fmt.Sprintf("%s VM saved", global.Brandname),
		Hint:			fmt.Sprintf("%s VM is in a saved state.", global.Brandname),
	},

	VmStateRunning: {
		Name:			VmStateRunning,
		VmState:		VmStateRunning,
		VmWantState:	VmStateDontCare,
		ApiState:		VmStateRunning,
		ApiWantState:	VmStateDontCare,

		VmTitleState:	fmt.Sprintf("%s VM: running", global.Brandname),
		VmHintState:	fmt.Sprintf("%s VM: running", global.Brandname),
		VmIconState:	IconUp,

		ApiTitleState:	fmt.Sprintf("%s API: running", global.Brandname),
		ApiHintState:	fmt.Sprintf("%s API: running", global.Brandname),
		ApiIconState:	IconUp,

		Title:			fmt.Sprintf("%s running", global.Brandname),
		Hint:			fmt.Sprintf("%s is running.", global.Brandname),
	},

	VmStateStarting: {
		Name:			VmStateStarting,
		VmState:		VmStateRunning,
		VmWantState:	VmStateDontCare,
		ApiState:		VmStateStarting,
		ApiWantState:	VmStateDontCare,

		VmTitleState:	fmt.Sprintf("%s VM: running", global.Brandname),
		VmHintState:	fmt.Sprintf("%s VM: running", global.Brandname),
		VmIconState:	IconUp,

		ApiTitleState:	fmt.Sprintf("%s API: starting", global.Brandname),
		ApiHintState:	fmt.Sprintf("%s API: starting", global.Brandname),
		ApiIconState:	IconStarting,

		Title:			fmt.Sprintf("%s API starting", global.Brandname),
		Hint:			fmt.Sprintf("%s API is starting.", global.Brandname),
	},

	VmStateStopping: {
		Name:			VmStateStopping,
		VmState:		VmStateRunning,
		VmWantState:	VmStateDontCare,
		ApiState:		VmStateStopping,
		ApiWantState:	VmStateDontCare,

		VmTitleState:	fmt.Sprintf("%s VM: stopping", global.Brandname),
		VmHintState:	fmt.Sprintf("%s VM: stopping", global.Brandname),
		VmIconState:	IconUp,

		ApiTitleState:	fmt.Sprintf("%s API: stopping", global.Brandname),
		ApiHintState:	fmt.Sprintf("%s API: stopping", global.Brandname),
		ApiIconState:	IconStopping,

		Title:			fmt.Sprintf("%s VM stopping", global.Brandname),
		Hint:			fmt.Sprintf("%s VM is stopping.", global.Brandname),
	},

	VmStateStarting2: {
		Name:			VmStateStarting,
		VmState:		VmStatePowerOff,
		VmWantState:	VmStateRunning,
		ApiState:		VmStateDontCare,
		ApiWantState:	VmStateDontCare,

		VmTitleState:	fmt.Sprintf("%s VM: running", global.Brandname),
		VmHintState:	fmt.Sprintf("%s VM: running", global.Brandname),
		VmIconState:	IconUp,

		ApiTitleState:	fmt.Sprintf("%s API: starting", global.Brandname),
		ApiHintState:	fmt.Sprintf("%s API: starting", global.Brandname),
		ApiIconState:	IconStarting,

		Title:			fmt.Sprintf("%s VM starting", global.Brandname),
		Hint:			fmt.Sprintf("%s VM is starting.", global.Brandname),
	},

	VmStateStopping2: {
		Name:			VmStateStopping,
		VmState:		VmStateRunning,
		VmWantState:	VmStatePowerOff,
		ApiState:		VmStateDontCare,
		ApiWantState:	VmStateDontCare,

		VmTitleState:	fmt.Sprintf("%s VM: stopping", global.Brandname),
		VmHintState:	fmt.Sprintf("%s VM: stopping", global.Brandname),
		VmIconState:	IconUp,

		ApiTitleState:	fmt.Sprintf("%s API: stopping", global.Brandname),
		ApiHintState:	fmt.Sprintf("%s API: stopping", global.Brandname),
		ApiIconState:	IconStopping,

		Title:			fmt.Sprintf("%s VM stopping", global.Brandname),
		Hint:			fmt.Sprintf("%s VM is stopping.", global.Brandname),
	},
}


// FSM handling.
func (state *BoxState) GetStateMeaning() (VmState) {

	var s VmState
	found := false

	for _, s = range DecodeVmState {
		switch {
			case (s.VmState == VmStateDontCare) && (s.ApiState == state.API.CurrentState) && (s.VmWantState == VmStateDontCare) && (s.ApiWantState == VmStateDontCare):
				// Handle don't care states first.
				fmt.Printf("1: '%v' : '%v'	'%v' : '%v'\n", s.VmState, state.VM.CurrentState, s.ApiState, state.API.CurrentState)
				found = true

			case (s.VmState == state.VM.CurrentState) && (s.ApiState == VmStateDontCare) && (s.VmWantState == VmStateDontCare) && (s.ApiWantState == VmStateDontCare):
				// Handle don't care states first.
				fmt.Printf("2: '%v' : '%v'	'%v' : '%v'\n", s.VmState, state.VM.CurrentState, s.ApiState, state.API.CurrentState)
				found = true

			case (s.VmState == state.VM.CurrentState) && (s.ApiState == state.API.CurrentState) && (s.VmWantState == VmStateDontCare) && (s.ApiWantState == VmStateDontCare):
				// .
				fmt.Printf("3: '%v' : '%v'	'%v' : '%v'\n", s.VmState, state.VM.CurrentState, s.ApiState, state.API.CurrentState)
				found = true

			case (s.VmState == VmStateDontCare) && (s.ApiState == VmStateDontCare) && (s.VmWantState == state.VM.WantState) && (s.VmState == state.VM.CurrentState):
				// Handle want states last.
				fmt.Printf("4: '%v' : '%v'	'%v' : '%v'\n", s.VmState, state.VM.CurrentState, s.ApiState, state.API.CurrentState)
				found = true
		}

		if found == true {
			break
		}
	}
	fmt.Printf("(%v) '%v' : '%v'\n", s.Name, s.VmState, s.ApiState)

	return s
}


/*
@TODO
VM states can be one of these possibilities:
VM		API		Description
UNKNOWN	X		No VM created.
DOWN	DOWN	VM shutdown and of course no API running.
DOWN	UP		INVALID STATE.
UP		DOWN	VM running, but API not yet running or not functioning.
UP		UP		Everything working OK.

VM		API		States
UNKNOWN	X		AbsentState
DOWN	DOWN	HaltedState | StoppedState | StoppingState
DOWN	UP		UnknownState
UP		DOWN	RunningState | StartedState | StartingState
UP		UP		RunningApiState


*/

/*
var StateMeaning = map[string]string{
	UnknownState: fmt.Sprintf("%s VM & Heartbeat in an unknown state", global.Brandname),
	HaltedState:  fmt.Sprintf("%s VM & Heartbeat halted", global.Brandname),
	RunningState: fmt.Sprintf("%s VM running, Heartbeat halted", global.Brandname),
	StartedState: fmt.Sprintf("%s VM running, Heartbeat halted", global.Brandname),
	NotOkState:   fmt.Sprintf("%s VM running, Heartbeat halted", global.Brandname),
	OkState:      fmt.Sprintf("%s VM running, Heartbeat running", global.Brandname),
}
*/

const (
	DefaultWaitDelay       = time.Second
	DefaultWaitRetries     = 90
	DefaultConsoleHost     = "127.0.0.1"
	DefaultConsolePort     = "2023"
	DefaultConsoleOkString = "Gearbox Heartbeat"
	DefaultShowConsole     = false
	DefaultConsoleReadWait = time.Second * 5
)
