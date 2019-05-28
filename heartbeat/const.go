package heartbeat

import (
	"gearbox/box"
	"gearbox/heartbeat/daemon"
	"gearbox/heartbeat/gbevents"
	"gearbox/heartbeat/monitor"
	oss "gearbox/os_support"
	"github.com/getlantern/systray"
	"time"
)


// This really needs to be refactored!
type State struct {
	Code    int
	Overall string
	Box     box.BoxState
	Unfsd   monitor.UnfsdState
}

type Heartbeat struct {
	Boxname        string
	EventBroker		   *gbevents.EventBroker
	BoxInstance    *box.Box
	DaemonInstance *daemon.Daemon
	NfsInstance    *monitor.Unfsd
	State          State
	OvaFile        string
	PidFile        string

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


type menuStruct struct {
	vmStatusEntry    *systray.MenuItem
	apiStatusEntry   *systray.MenuItem
	unfsdStatusEntry *systray.MenuItem
	startEntry       *systray.MenuItem
	stopEntry        *systray.MenuItem
	adminEntry       *systray.MenuItem
	sshEntry         *systray.MenuItem
	quitEntry        *systray.MenuItem
	restartEntry     *systray.MenuItem

	helpEntry        *systray.MenuItem
	versionEntry     *systray.MenuItem
	updateEntry      *systray.MenuItem
	createEntry      *systray.MenuItem
}


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

