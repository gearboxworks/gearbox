package box

import (
	"gearbox/box/external/unfsd"
	"gearbox/box/external/vmbox"
	"gearbox/eventbroker"
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/ospaths"
	"gearbox/eventbroker/states"
	"github.com/gearboxworks/go-osbridge"

	//	oss "gearbox/os_support"
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

type Box struct {
	EntityId       messages.MessageAddress
	EntityName     messages.MessageAddress
	Boxname        string
	Version        string
	NfsInstance    *unfsd.Unfsd
	State          *states.Status
	menu           Menus
	EventBroker    *eventbroker.EventBroker
	VmBox          *vmbox.VmBox

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

	baseDir         *ospaths.Dir
	pidFile         string
	osBridge 		osbridge.OsBridger
	osPaths         *ospaths.BasePaths
}
type Args Box


type Menus map[messages.MessageAddress]*Menu
type Menu struct {
	MenuItem         *systray.MenuItem
	PrefixToolTip    string
	PrefixMenu       string
	CurrentIcon      string
	State            states.State
}
type MenuItem             systray.MenuItem

const (
	UnknownState = "unknown"
)

const (
	DefaultEntityName      = "Heartbeat"

	DefaultWaitDelay       = time.Second
	DefaultWaitRetries     = 90
	DefaultConsoleHost     = "127.0.0.1"
	DefaultConsolePort     = "2023"
	DefaultConsoleOkString = "Gearbox Heartbeat"
	DefaultShowConsole     = false
	DefaultConsoleReadWait = time.Second * 5
	DefaultPidFile         = "heartbeat.pid"
	DefaultBaseDir         = "dist/heartbeat"
)


const pidName = "[Gearbox]"


const (
	DefaultLogo = "img/IconLogo.ico"
	DefaultUp = "img/UpArrow.ico"
	DefaultDown = "img/DownArrow.ico"

	IconLogo = "img/IconLogo.ico"
	IconError = "img/IconError.ico"
	IconWarning = "img/IconWarning.ico"
	IconUp = "img/IconUp.ico"
	IconDown = "img/IconDown.ico"
	IconStarting = "img/IconStarting.ico"
	IconStopping = "img/IconStopping.ico"
)

