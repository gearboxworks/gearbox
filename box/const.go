package box

import (
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/states"
	//	oss "gearbox/os_support"
//	"github.com/gearboxworks/go-systray"
	"github.com/getlantern/systray"
	"time"
)

// This really needs to be refactored!
//type State struct {
//	Code    int
//	Overall string
//	Box     box.BoxState
//	Unfsd   unfsd.UnfsdState
//}

type Version = string

const LatestVersion Version = "latest"

type Menus map[msgs.Address]*Menu
type Menu struct {
	MenuItem      *systray.MenuItem
	PrefixToolTip string
	PrefixMenu    string
	CurrentIcon   string
	State         states.State
}
type MenuItem systray.MenuItem

const (
	UnknownState = "unknown"
)

const (
	DefaultEntityName = "Heartbeat"

	DefaultWaitDelay       = time.Second
	DefaultWaitRetries     = 90
	DefaultConsoleHost     = "127.0.0.1"
	DefaultConsolePort     = "2023"
	DefaultConsoleOkString = "Gearbox Heartbeat"
	DefaultShowConsole     = false
	DefaultConsoleReadWait = time.Second * 5
	DefaultPidFile         = "heartbeat.pid"
	DefaultBaseDir         = "app/dist/heartbeat"
)

const pidName = "[Gearbox]"

const (
	DefaultLogo = "img/IconLogo.ico"
	DefaultUp   = "img/UpArrow.ico"
	DefaultDown = "img/DownArrow.ico"

	IconLogo     = "img/IconLogo.ico"
	IconError    = "img/IconError.ico"
	IconWarning  = "img/IconWarning.ico"
	IconUp       = "img/IconUp.ico"
	IconDown     = "img/IconDown.ico"
	IconStarting = "img/IconStarting.ico"
	IconStopping = "img/IconStopping.ico"
)
