package vmbox

import (
	"time"
)

const JsonFilePattern = "%s.json"

const (
	DefaultVmWaitTime       = time.Millisecond * 1000
	DefaultRetries          = 12
	DefaultRetryDelay       = time.Second * 10
	DefaultConsoleReadWait  = time.Second * 5
	DefaultConsoleWaitDelay = time.Second
	DefaultBootWaitTime     = time.Second * 60
	DefaultRunWaitTime      = time.Second * 5
	DefaultSshHost          = "gearbox.local"
	DefaultSshPort          = "22"

	DefaultConsoleOkString = "Gearbox Heartbeat"
	Basedir                = "/home/gearbox/projects"
	IconLogoPng            = "app/dist/heartbeat/img/IconLogo.png"
)

const (
	Package                    = "vmbox"
	InterfaceTypeVmBox         = "*" + Package + ".VmBox"
	InterfaceTypeService       = "*" + Package + ".Service"
	InterfaceTypeServiceConfig = "*" + Package + ".ServiceConfig"
)
