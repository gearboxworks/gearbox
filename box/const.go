package box

import (
	"fmt"
	"time"
)

const (
	Brandname   = "Gearbox"
	Basedir     = "/home/gearbox/projects"
	OvaFileName = "/box/vm/Gearbox.ova"
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
	UnknownState: fmt.Sprintf("%s VM & Heartbeat in an unknown state", Brandname),
	HaltedState:  fmt.Sprintf("%s VM & Heartbeat halted", Brandname),
	RunningState: fmt.Sprintf("%s VM running, Heartbeat halted", Brandname),
	StartedState: fmt.Sprintf("%s VM running, Heartbeat halted", Brandname),
	NotOkState:   fmt.Sprintf("%s VM running, Heartbeat halted", Brandname),
	OkState:      fmt.Sprintf("%s VM running, Heartbeat running", Brandname),
}

const (
	DefaultWaitDelay       = time.Second
	DefaultWaitRetries     = 90
	DefaultConsoleHost     = "127.0.0.1"
	DefaultConsolePort     = "2023"
	DefaultConsoleOkString = "Gearbox Heartbeat"
	DefaultShowConsole     = false
	DefaultConsoleReadWait = time.Second * 5
)
