package box

import (
	"fmt"
	"gearbox/global"
	"time"
)

const (
	Basedir     = "/home/gearbox/projects"
	OvaFileName = "/box/vm/Parent.ova"
)

const (
	ErrorState   = "error"
	UnknownState = "unknown"
	AbsentState  = "absent"
	HaltedState  = "halted"
	RunningState = "running"
	StartedState = "started"
	OkState      = "ok"
	NotOkState   = "nok"
)

var StateMeaning = map[State]string{
	UnknownState: fmt.Sprintf("%s VM & Heartbeat in an unknown state", global.Brandname),
	HaltedState:  fmt.Sprintf("%s VM & Heartbeat halted", global.Brandname),
	RunningState: fmt.Sprintf("%s VM running, Heartbeat halted", global.Brandname),
	StartedState: fmt.Sprintf("%s VM running, Heartbeat halted", global.Brandname),
	NotOkState:   fmt.Sprintf("%s VM running, Heartbeat halted", global.Brandname),
	OkState:      fmt.Sprintf("%s VM running, Heartbeat running", global.Brandname),
}

const (
	DefaultWaitDelay       = time.Second
	DefaultWaitRetries     = 90
	DefaultConsoleHost     = "127.0.0.1"
	DefaultConsolePort     = "2023"
	DefaultConsoleOkString = "Parent Heartbeat"
	DefaultShowConsole     = false
	DefaultConsoleReadWait = time.Second * 5
)
