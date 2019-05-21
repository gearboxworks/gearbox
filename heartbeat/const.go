package heartbeat

import (
	"time"
)


const (
	UnknownState = "unknown"
)

const (
	DefaultWaitDelay       = time.Second
	DefaultWaitRetries     = 90
	DefaultConsoleHost     = "127.0.0.1"
	DefaultConsolePort     = "2023"
	DefaultConsoleOkString = "Gearbox Heartbeat"
	DefaultShowConsole     = false
	DefaultConsoleReadWait = time.Second * 5
	DefaultPidFile         = "heartbeat/heartbeat.pid"
)

const pidName = "[Gearbox]"

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