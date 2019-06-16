package heartbeat

import (
	"gearbox/box"
	"gearbox/eventbroker"
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/states"
	"gearbox/heartbeat/external/unfsd"
	"gearbox/heartbeat/external/vmbox"
	oss "gearbox/os_support"
	"github.com/getlantern/systray"
	"time"
)


// This really needs to be refactored!
type State struct {
	Code    int
	Overall string
	Box     box.BoxState
	Unfsd   unfsd.UnfsdState
}

type Heartbeat struct {
	Boxname        string
	EventBroker    *eventbroker.EventBroker
	VmBox          *vmbox.VmBox
	Version        string
	//osRelease      *vmbox.Release
	//BoxInstance    *box.Box
	//DaemonInstance *daemon.Daemon
	NfsInstance    *unfsd.Unfsd
	State          State
	PidFile        string
	menu           Menus
	baseDir        string

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

	OsSupport oss.OsSupporter
}
type Args Heartbeat


type Menus map[messages.MessageAddress]*Menu
type Menu struct {
	MenuItem         *systray.MenuItem
	PrefixToolTip    string
	PrefixMenu       string
	CurrentIcon      string
	State            states.State

	//vmStatusEntry    *systray.MenuItem
	//apiStatusEntry   *systray.MenuItem
	//unfsdStatusEntry *systray.MenuItem
	//
	//startEntry       *systray.MenuItem
	//stopEntry        *systray.MenuItem
	//adminEntry       *systray.MenuItem
	//updateEntry      *systray.MenuItem
	//createEntry      *systray.MenuItem
	//sshEntry         *systray.MenuItem
	//quitEntry        *systray.MenuItem
	//restartEntry     *systray.MenuItem
	//
	//helpEntry        *systray.MenuItem
	//versionEntry     *systray.MenuItem
}
type MenuItem             systray.MenuItem

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
	DefaultPidFile         = "heartbeat.pid"
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

