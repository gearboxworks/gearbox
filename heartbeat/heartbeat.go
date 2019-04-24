package heartbeat

import (
	"fmt"
	"gearbox/box"
	"gearbox/global"
	"gearbox/heartbeat/daemon"
	"gearbox/help"
	"gearbox/only"
	"gearbox/os_support"
	"gearbox/ssh"
	"gearbox/status"
	"gearbox/status/is"
	"github.com/apcera/libretto/virtualmachine/virtualbox"
	"os/exec"
	"os/signal"
	"syscall"
	//	"github.com/apcera/libretto/virtualmachine/virtualbox"
	//	lbssh "github.com/apcera/libretto/ssh"
	"github.com/getlantern/systray"
	"github.com/jinzhu/copier"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type State string

type Heartbeat struct {
	Boxname  string
	BoxInstance *box.Box
	DaemonInstance *daemon.Daemon
	State    State
	OvaFile  string

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
	statusEntry  *systray.MenuItem
	startEntry  *systray.MenuItem
	stopEntry   *systray.MenuItem
	adminEntry  *systray.MenuItem
	sshEntry    *systray.MenuItem
	quitEntry   *systray.MenuItem
	unloadEntry *systray.MenuItem
	helpEntry   *systray.MenuItem
	createEntry *systray.MenuItem
}


type GearboxVM struct {
	virtualbox.VM
}

var intentDelay = false		// Used to change delay times when we have just performed an action.

func NewHeartbeat(OsSupport oss.OsSupporter, args ...Args) *Heartbeat {
	var _args Args
	if len(args) > 0 {
		_args = args[0]
	}

	_args.OsSupport = OsSupport
	foo := box.Args{}
	copier.Copy(&foo, &_args)
	_args.BoxInstance = box.NewBox(OsSupport, foo)

	if _args.Boxname == "" {
		_args.BoxInstance.Boxname = global.Brandname
	} else {
		_args.BoxInstance.Boxname = _args.Boxname
	}

	if _args.WaitDelay == 0 {
		_args.BoxInstance.WaitDelay = DefaultWaitDelay
	} else {
		_args.BoxInstance.WaitDelay = _args.WaitDelay
	}

	if _args.WaitRetries == 0 {
		_args.BoxInstance.WaitRetries = DefaultWaitRetries
	} else {
		_args.BoxInstance.WaitRetries = _args.WaitRetries
	}

	if _args.ConsoleHost == "" {
		_args.BoxInstance.ConsoleHost = DefaultConsoleHost
	} else {
		_args.BoxInstance.ConsoleHost = _args.ConsoleHost
	}

	if _args.ConsolePort == "" {
		_args.BoxInstance.ConsolePort = DefaultConsolePort
	} else {
		_args.BoxInstance.ConsolePort = _args.ConsolePort
	}

	if _args.ConsoleOkString == "" {
		_args.BoxInstance.ConsoleOkString = DefaultConsoleOkString
	} else {
		_args.BoxInstance.ConsoleOkString = _args.ConsoleOkString
	}

	if _args.ConsoleReadWait == 0 {
		_args.BoxInstance.ConsoleReadWait = DefaultConsoleReadWait
	} else {
		_args.BoxInstance.ConsoleReadWait = _args.ConsoleReadWait
	}

	if _args.SshUsername == "" {
		_args.BoxInstance.SshUsername = ssh.DefaultUsername
	} else {
		_args.BoxInstance.SshUsername = _args.SshUsername
	}

	if _args.SshPassword == "" {
		_args.BoxInstance.SshPassword = ssh.DefaultPassword
	} else {
		_args.BoxInstance.SshPassword = _args.SshPassword
	}

	if _args.SshPublicKey == "" {
		_args.BoxInstance.SshPublicKey = ssh.DefaultKeyFile
	} else {
		_args.BoxInstance.SshPublicKey = _args.SshPublicKey
	}

	_args.DaemonInstance = daemon.NewDaemon(_args.OsSupport, daemon.Args{
		Boxname: _args.Boxname,
	})

	hb := &Heartbeat{}
	*hb = Heartbeat(_args)

	// Query VB to see if it exists.
	// If not return nil.

	return hb
}


func (me *Heartbeat) onReady() {

	var menu menuStruct
	fmt.Printf("Gearbox: Heartbeat restarted.\n")

	systray.SetIcon(me.getIcon(DefaultLogo))
	systray.SetTitle("Gearbox:")

	menu.statusEntry = systray.AddMenuItem("Starting", "Current state of Gearbox")

	systray.AddSeparator()
	menu.helpEntry = systray.AddMenuItem("About Gearbox", "Contact Gearbox help")


	systray.AddSeparator()
	menu.adminEntry = systray.AddMenuItem("Admin", "Open Gearbox admin interface")
	menu.createEntry = systray.AddMenuItem("Create Box", "Create a Gearbox OS instance")
	menu.startEntry = systray.AddMenuItem("Start Box", "Start Gearbox OS instance")
	menu.stopEntry = systray.AddMenuItem("Stop Box", "Stop Gearbox OS instance")

	menu.sshEntry = systray.AddMenuItem("SSH", "Connect to Gearbox OS via SSH")
	//menu.consoleEntry = systray.AddMenuItem("Console", "Show the Gearbox OS console")

	systray.AddSeparator()
	pid := os.Getpid()
	menu.quitEntry = systray.AddMenuItem("Restart App", fmt.Sprintf("Restart this app [pid:%v]", pid))
	menu.unloadEntry = systray.AddMenuItem("Terminate App",fmt.Sprintf("Terminate this app [pid:%v]", pid))

	go func() {
		for {
			if intentDelay {
				// User has requested a change, check on cached results faster.
				// results will be updated by concurrent functions.
				//fmt.Printf("CACHE POLL\n")
				_, state := me.BoxInstance.GetCachedState()
				me.SetState(menu, state)
				time.Sleep(time.Second)

			} else {
				// Normal polling.
				//fmt.Printf("NORMAL POLL\n")
				_, state := me.BoxInstance.GetState()
				me.SetState(menu, state)
				time.Sleep(5 * time.Second)
			}
		}
	}()

	go func() {
		for {
			select {
				case <- menu.startEntry.ClickedCh:
					fmt.Printf("Menu: Start\n")
					intentDelay = true
					me.BoxInstance.StartBox()
					intentDelay = false

				case <- menu.stopEntry.ClickedCh:
					fmt.Printf("Menu: Stop\n")
					intentDelay = true
					me.BoxInstance.StopBox()
					intentDelay = false

				case <- menu.adminEntry.ClickedCh:
					fmt.Printf("Menu: Admin\n")
					me.openAdmin()

				//case <- menu.consoleEntry.ClickedCh:
				//	fmt.Printf("Menu: Console\n")

				case <- menu.sshEntry.ClickedCh:
					fmt.Printf("Menu: SSH\n")
					me.openTerminal()

				case <- menu.helpEntry.ClickedCh:
					fmt.Printf("Menu: Help\n")
					me.openAbout()

				case <- menu.createEntry.ClickedCh:
					fmt.Printf("Menu: Create\n")

				case <- menu.unloadEntry.ClickedCh:
					me.StopHeartbeat()
					systray.Quit()
					return

				case <- menu.quitEntry.ClickedCh:
					systray.Quit()
					return
			}
		}
	}()
}


func (me *Heartbeat) openAdmin() error {

	execPath, err := os.Executable()
	if err == nil {
		fmt.Printf("Menu: Admin - %s\n", execPath)
	}

	execCwd, err := os.Getwd()
	if err == nil {
		fmt.Printf("Menu: Admin - %s\n", execCwd)
	}

	cmd := exec.Command(execPath,"admin")
	err = cmd.Run()

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	return err
}


func (me *Heartbeat) openTerminal() error {

	execPath, err := os.Executable()
	if err == nil {
		fmt.Printf("Menu: Admin - %s\n", execPath)
	}

	execCwd, err := os.Getwd()
	if err == nil {
		fmt.Printf("Menu: Admin - %s\n", execCwd)
	}

	cmd := exec.Command("open", "-b", "com.apple.terminal", execPath, "--args", "ssh")
	err = cmd.Run()

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	return err
}


func (me *Heartbeat) openAbout() error {

	cmd := exec.Command("open", "https://gearbox.works/")
	err := cmd.Run()

	return err
}


func (me *Heartbeat) onExit() {
	// Cleaning stuff here.
}


func getClockTime(tz string) string {
	t := time.Now()
	utc, _ := time.LoadLocation(tz)

	return t.In(utc).Format("15:04:05")
}


func (me *Heartbeat) getIcon(s string) []byte {

	fp := filepath.FromSlash(fmt.Sprintf("%s/%s", me.OsSupport.GetAdminRootDir(), s))
	if fp == "" {
		return nil
	}

	b, err := ioutil.ReadFile(fp)
	if err != nil {
		fmt.Print(err)
	}

	return b
}


func (me *Heartbeat) SetState(menu menuStruct, state box.BoxState) (returnValue int) {

	switch {
		case (state.VM.CurrentState == box.StateUnknown):
			fmt.Printf("STATE: UNKNOWN\n")
			systray.SetTitle("Gearbox: UNKNOWN")
			systray.SetTooltip("Gearbox is in an unknown state.")
			menu.statusEntry.SetIcon(me.getIcon(DefaultLogo))
			menu.statusEntry.SetTitle("State: unknown")

			returnValue = box.StateUnknown
			menu.statusEntry.Enable()
			menu.stopEntry.Disable()
			menu.startEntry.Disable()
			menu.sshEntry.Disable()
			menu.createEntry.Enable()


		case (state.VM.CurrentState == box.StateDown):
			// fmt.Printf("STATE: HALTED\n")
			systray.SetTitle("Gearbox: HALTED")
			systray.SetTooltip("Gearbox is halted.")
			menu.statusEntry.SetIcon(me.getIcon(DefaultDown))
			menu.statusEntry.SetTitle("State: halted")

			returnValue = box.StateDown
			menu.statusEntry.Disable()
			menu.stopEntry.Disable()
			menu.startEntry.Enable()
			menu.sshEntry.Disable()
			menu.createEntry.Disable()


		case (state.VM.CurrentState == box.StateUp) && (state.API.CurrentState == box.StateUp):
			// fmt.Printf("STATE: RUNNING\n")
			systray.SetTitle("Gearbox: RUNNING")
			systray.SetTooltip("Gearbox is running.")
			menu.statusEntry.SetIcon(me.getIcon(DefaultUp))
			menu.statusEntry.SetTitle("State: running")

			returnValue = box.StateUp
			menu.statusEntry.Disable()
			menu.stopEntry.Enable()
			menu.startEntry.Disable()
			menu.sshEntry.Enable()
			menu.createEntry.Disable()


		case (state.VM.WantState == box.StateUp) && (state.VM.CurrentState != box.StateUp):
			fallthrough
		case (state.API.WantState == box.StateUp) && (state.API.CurrentState != box.StateUp):
			fallthrough
		case (state.VM.CurrentState == box.StateStarting):
			fmt.Printf("STATE: STARTING\n")
			systray.SetTitle("Gearbox: STARTING")
			systray.SetTooltip("Gearbox starting up.")
			menu.statusEntry.SetIcon(me.getIcon(DefaultUp))
			menu.statusEntry.SetTitle("State: starting")

			returnValue = box.StateStarting
			menu.statusEntry.Disable()
			menu.stopEntry.Disable()
			menu.startEntry.Disable()
			menu.sshEntry.Disable()
			menu.createEntry.Disable()


		case (state.VM.WantState == box.StateDown) && (state.VM.CurrentState != box.StateDown):
			fallthrough
		case (state.API.WantState == box.StateDown) && (state.API.CurrentState != box.StateDown):
			fallthrough
		case (state.VM.CurrentState == box.StateStopping):
			fmt.Printf("STATE: STOPPING\n")
			systray.SetTitle("Gearbox: STOPPING")
			systray.SetTooltip("Gearbox is stopping.")
			menu.statusEntry.SetIcon(me.getIcon(DefaultDown))
			menu.statusEntry.SetTitle("State: stopping")

			returnValue = box.StateStopping
			menu.statusEntry.Disable()
			menu.stopEntry.Disable()
			menu.startEntry.Disable()
			menu.sshEntry.Disable()
			menu.createEntry.Disable()


		default:
			fmt.Printf("STATE: UNKNOWN DEFAULT\n")
			systray.SetTitle("Gearbox: UNKNOWN")
			systray.SetTooltip("Gearbox is in an unknown state.")
			menu.statusEntry.SetIcon(me.getIcon(DefaultLogo))
			menu.statusEntry.SetTitle("State: unknown")

			returnValue = box.StateUnknown
			menu.statusEntry.Enable()
			menu.stopEntry.Disable()
			menu.startEntry.Disable()
			menu.sshEntry.Disable()
			menu.createEntry.Enable()
	}

	return
}


func (me *Heartbeat) Initialize() (sts status.Status) {

	//var boxArgs box.Args

	for range only.Once {

		cfgdir := me.OsSupport.GetUserConfigDir()

		// sts := gearbox.BoxInstance.StartBox(boxArgs)

		me.OvaFile = fmt.Sprintf("%s/%s", cfgdir, "foo")

		// The OvaFile is created from an export from within VirtualBox.
		// VBoxManage export Parent -o Parent.ova --options manifest
		// This was the best way to create a base template, avoiding too much code bloat.
		// And allows multiple VM frameworks to be used with libretto.
		// It doesn't include the ISO image yet as it is too large.
		// Once the ISO image size has been reduced, we can do this:
		// VBoxManage export Parent -o Parent.ova --options iso,manifest

/*
		_, err := os.Stat(me.OvaFile)
		if os.IsExist(err) {
			break
		}
		err = vm.RestoreAssets(string(cfgdir), strings.TrimLeft(OvaFileName, string(os.PathSeparator)))
		if err != nil {
			sts = status.Wrap(err, &status.Args{
				Message: fmt.Sprintf("%s: VM OVA file cannot be created as'%s'.", global.Brandname, me.OvaFile),
			})
			break
		}
*/
	}

	return sts
}


func (me *Heartbeat) HeartbeatDaemon() (sts status.Status) {

	for range only.Once {

		sts = me.GetState()
		if is.Error(sts) {
			break
		}

		if !me.DaemonInstance.IsParentInit() {
			fmt.Printf("Sub-command not available for user.\n")
			sts = status.Fail(&status.Args{
				Message: "Sub-command not available for user",
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    UnknownState,
			})
			break
		}

// Doesn't seem to work properly.
//		if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
//			fmt.Printf("Sub-command not available for user.\n")
//			//break
//		}

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigs
			os.Exit(0)
		}()

		systray.Run(me.onReady, me.onExit)
	}

	return sts
}


func (me *Heartbeat) StartHeartbeat() (sts status.Status) {

	for range only.Once {

		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

//		if me.DaemonInstance.IsRunning() {
//			fmt.Printf("%s Heartbeat - Restarting service.\n", global.Brandname)
//			sts = me.DaemonInstance.Unload()
//			if is.Error(sts) {
//				break
//			}
//		}

//		if me.DaemonInstance.IsLoaded() {
//			fmt.Printf("%s Heartbeat - Restarting service.\n", global.Brandname)
//			sts = me.DaemonInstance.Unload()
//			if is.Error(sts) {
//				break
//			}
//		}

		sts = me.DaemonInstance.Load()
		if is.Error(sts) {
			break
		}
		fmt.Printf("%s\n", sts.Message())

	}

	return sts
}


func (me *Heartbeat) StopHeartbeat() (sts status.Status) {

	for range only.Once {

		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		sts = me.DaemonInstance.Unload()
		if is.Error(sts) {
			break
		}
		fmt.Printf("%s\n", sts.Message())
		// fmt.Printf("%s Heartbeat - Started service.\n", global.Brandname)

	}

	return sts
}


func (me *Heartbeat) RestartHeartbeat() (sts status.Status) {

	for range only.Once {

		sts = me.StopHeartbeat()
		if is.Error(sts) {
			break
		}

		sts = me.StartHeartbeat()
		if is.Error(sts) {
			break
		}

	}

	return sts
}


func (me *Heartbeat) GetState() (sts status.Status) {

	for range only.Once {

		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		if me == nil {
			sts = status.Fail(&status.Args{
				Message: "unexpected failure",
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    UnknownState,
			})
			break
		}

		sts = me.DaemonInstance.GetState()
		if is.Error(sts) {
			break
		}
		if sts != nil {
			fmt.Printf(sts.Message())
		}

		sts, _ = me.BoxInstance.GetState()
		if is.Error(sts) {
			break
		}
		if sts != nil {
			fmt.Printf(sts.Message())
		}

	}

	return sts
}


func EnsureNotNil(bx *Heartbeat) (sts status.Status) {
	if bx == nil {
		sts = status.Fail(&status.Args{
			Message: "unexpected error",
			Help:    help.ContactSupportHelp(), // @TODO need better support here
			Data:    UnknownState,
		})
	}
	return sts
}
