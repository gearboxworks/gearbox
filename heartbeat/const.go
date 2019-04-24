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
)

const pidName = "[Gearbox]"

const (
	DefaultLogo = "heartbeat/img/logo.ico"
	DefaultUp = "heartbeat/img/UpArrow.ico"
	DefaultDown = "heartbeat/img/DownArrow.ico"
)