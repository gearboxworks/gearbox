package box

import (
	"fmt"
	"gearbox/global"
	"time"
)

const (
	Basedir     = "/home/gearbox/projects"
	OvaFileName = "/box/vm/Gearbox.ova"
)

const (
	ErrorState   = "error"
	UnknownState = "unknown"
	VmStateHalted  = "halted"
	VmStateRunning = "running"
	OkState      = "ok"
	NotOkState   = "nok"
)

const (
	StateInit		= 0
	StateUnknown	= 1
	StateDontCare	= 2
	StateDown		= 3
	StateStarting	= 4
	StateUp			= 5
	StateStopping	= 6
)


type States struct {
	VmState		int
	ApiState	int
	Title		string
	Hint		string
}

var StateMeaning = [...]States {
	States {
		VmState:	StateUnknown,
		ApiState:	StateUnknown,
		Title:	fmt.Sprintf("%s: UNKNOWN", global.Brandname),
		Hint:	fmt.Sprintf("%s is in an unknown state.", global.Brandname),
	},
}


/*
TODO
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
