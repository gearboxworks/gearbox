package heartbeat

import (
	"fmt"
	"gearbox/box"
	"gearbox/global"
	"gearbox/heartbeat/daemon"
	"gearbox/heartbeat/monitor"
	"gearbox/help"
	"gearbox/only"
	"gearbox/os_support"
	"gearbox/ssh"
	"github.com/apcera/libretto/virtualmachine/virtualbox"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/sqweek/dialog"
	"io/ioutil"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	//	"github.com/apcera/libretto/virtualmachine/virtualbox"
	//	lbssh "github.com/apcera/libretto/ssh"
	"github.com/getlantern/systray"
	"github.com/jinzhu/copier"
	"os"
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
	vmStatusEntry  *systray.MenuItem
	apiStatusEntry  *systray.MenuItem
	unfsdStatusEntry  *systray.MenuItem

	startEntry  *systray.MenuItem
	stopEntry   *systray.MenuItem
	adminEntry  *systray.MenuItem
	sshEntry    *systray.MenuItem
	quitEntry   *systray.MenuItem
	restartEntry *systray.MenuItem
	helpEntry   *systray.MenuItem
	createEntry *systray.MenuItem
}


type GearboxVM struct {
	virtualbox.VM
}


func New(OsSupport oss.OsSupporter, args ...Args) (*Heartbeat, status.Status) {

	var _args Args
	var sts status.Status

	if len(args) > 0 {
		_args = args[0]
	}

	_args.OsSupport = OsSupport
	foo := box.Args{}
	copier.Copy(&foo, &_args)

	// Start a new VM Box instance.
	_args.BoxInstance = box.NewBox(OsSupport, foo)

	if _args.Boxname == "" {
		_args.BoxInstance.Boxname = global.Brandname
	} else {
		_args.BoxInstance.Boxname = _args.Boxname
	}
	//_args.BoxInstance.Boxname = "harry"		// DEBUG

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

	execPath, _ := os.Executable()
	execCwd := string(_args.OsSupport.GetAdminRootDir()) + "/heartbeat"	// os.Getwd()

	_args.PidFile = filepath.FromSlash(fmt.Sprintf("%s/%s", _args.OsSupport.GetAdminRootDir(), DefaultPidFile))

	// Start a new Daemon instance.
	_args.DaemonInstance = daemon.NewDaemon(_args.OsSupport, daemon.Args{
		Boxname: _args.Boxname,
		ServiceData: daemon.PlistData {
			Label: "com.gearbox.heartbeat",
			Program:   execPath,
			ProgramArgs: []string{"heartbeat", "daemon"},
			Path:      execCwd,
			PidFile: _args.PidFile,
			KeepAlive: true,
			RunAtLoad: true,
		},
	})

	hb := &Heartbeat{}
	*hb = Heartbeat(_args)

	return hb, sts
}


func (me *Heartbeat) HeartbeatDaemon() (sts status.Status) {

	for range only.Once {

		sts = me.GetState()
		if is.Error(sts) {
			break
		}

		//if !daemon.IsParentInit() {
		if daemon.IsParentInit() {
			fmt.Printf("Gearbox: Sub-command not available for user.\n")
			sts = status.Fail(&status.Args{
				Message: "Sub-command not available for user",
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    UnknownState,
			})
			break
		}
		fmt.Printf("Gearbox: Starting Heartbeat daemon.\n")

		// Doesn't seem to work properly - need a workaround of some sort.
		//		if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		//			fmt.Printf("Sub-command not available for user.\n")
		//			//break
		//		}

		// Handle exit signals.
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigs
			os.Exit(0)
		}()

		// Create a new VM Box instance.
		fmt.Printf("Gearbox: Creating unfsd instance.\n")
		me.NfsInstance, sts = monitor.NewUnfsd(me.OsSupport)

		fmt.Printf("Gearbox: Starting systray.\n")
		systray.Run(me.onReady, me.onExit)
		// Should never exit, unless we get a signal to do so.
	}

	return sts
}


func (me *Heartbeat) onReady() {

	var menu menuStruct
	var intentDelay = false		// Used to change delay times when the user has just performed an action.

	fmt.Printf("Gearbox: Heartbeat started.\n")

	systray.SetIcon(me.getIcon(DefaultLogo))
	systray.SetTitle("")

	menu.helpEntry = systray.AddMenuItem("About Gearbox", "Contact Gearbox help for" + me.BoxInstance.Boxname)

	systray.AddSeparator()
	menu.vmStatusEntry = systray.AddMenuItem("Box: Idle", "Current state of Gearbox VM")
	menu.vmStatusEntry.SetIcon(me.getIcon(DefaultLogo))
	menu.apiStatusEntry = systray.AddMenuItem("API: Idle", "Current state of Gearbox API")
	menu.apiStatusEntry.SetIcon(me.getIcon(DefaultLogo))
	menu.unfsdStatusEntry = systray.AddMenuItem("FS: Idle", "Current state of Gearbox NFS service")
	menu.unfsdStatusEntry.SetIcon(me.getIcon(DefaultLogo))

	systray.AddSeparator()
	menu.adminEntry = systray.AddMenuItem("Admin", "Open Gearbox admin interface")
	menu.createEntry = systray.AddMenuItem("Create Box", "Create a Gearbox OS instance")
	menu.startEntry = systray.AddMenuItem("Start Box", "Start Gearbox OS instance")
	menu.stopEntry = systray.AddMenuItem("Stop Box", "Stop Gearbox OS instance")

	menu.sshEntry = systray.AddMenuItem("SSH", "Connect to Gearbox OS via SSH")
	//menu.consoleEntry = systray.AddMenuItem("Console", "Show the Gearbox OS console")

	systray.AddSeparator()
	pid := os.Getpid()
	menu.restartEntry = systray.AddMenuItem("Restart Heartbeat", fmt.Sprintf("Restart this app [pid:%v]", pid))
	menu.quitEntry = systray.AddMenuItem("Quit",fmt.Sprintf("Terminate this app [pid:%v]", pid))

	sts := me.NfsInstance.Daemon.Load()
	if is.Error(sts) {
		fmt.Printf("%s\n", sts.Message())
		return
	}

	// Concurrent process: Provide status updates on systray.
	// Ideally, this should also send messages on message bus for actions to be taken. EG: Retry startup, disk full, etc.
	// Even further, these should be brokem out into methods to avoid having to hard code specific entities to monitor.
	go func() {
		var state State
		var sts status.Status

		for {
			if intentDelay {
				// User has requested a change, check on cached results faster.
				// results will be updated by concurrent functions.
				//fmt.Printf("CACHE POLL\n")

				// Check state of VM.
				me.State.Box, sts = me.BoxInstance.GetCachedState()
				if is.Error(sts) {
					// .
				}

				// Check state of UNFSD.
				me.State.Unfsd, sts = me.NfsInstance.GetState()
				if is.Error(sts) || is.Error(state.Unfsd.LastSts) {
					// .
				}

				me.SetMenuState(menu)
				time.Sleep(time.Second)

			} else {
				// Normal polling.
				//fmt.Printf("NORMAL POLL\n")

				// Check state of VM.
				me.State.Box, sts = me.BoxInstance.GetState()
				//fmt.Printf("STATE:\n%v\n%v\n", me.State.Box, sts)
				if is.Error(sts) {
					// .
				}

				// Check state of UNFSD.
				me.State.Unfsd, sts = me.NfsInstance.GetState()
				if is.Error(sts) || is.Error(state.Unfsd.LastSts) {
					// .
				}

				me.SetMenuState(menu)
				time.Sleep(5 * time.Second)
			}
		}
	}()

	// Concurrent process: Handle user clicky clicks on menu.
	go func() {
		for {
			select {
				case <- menu.startEntry.ClickedCh:
					fmt.Printf("Menu: Start\n")
					intentDelay = true
					me.BoxInstance.Start()
					intentDelay = false

				case <- menu.stopEntry.ClickedCh:
					fmt.Printf("Menu: Stop\n")
					intentDelay = true
					me.BoxInstance.Stop()
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
					if me.BoxInstance.State.VM.CurrentState == box.VmStateNotPresent {
						sts := me.BoxInstance.CreateBox()
						if is.Error(sts) {
							// Throw error
						}
					}

				case <- menu.restartEntry.ClickedCh:
					fmt.Printf("Menu: Restart\n")
					if me.confirmDialog("Restart Gearbox", "This will restart Gearbox Heartbeat, but keep services running.\nAre you sure?") {
						fmt.Printf("HEY!")
						systray.Quit()
					}

				case <- menu.quitEntry.ClickedCh:
					fmt.Printf("Menu: Quit\n")
					if me.confirmDialog("Shutdown Gearbox", "This will shutdown Gearbox and all Gearbox related services.\nAre you sure?") {
						intentDelay = true
						me.BoxInstance.Stop()
						me.NfsInstance.Stop()
						intentDelay = false

						me.StopHeartbeat()

						systray.Quit()
					}
			}
		}
	}()
}


func (me *Heartbeat) fileDialog(t string, m string) bool {
	dialog.Message("%s", "Please select a file").Title("Hello world!").Info()
	file, err := dialog.File().Title("Save As").Filter("All Files", "*").Save()
	fmt.Println(file)
	fmt.Println("Error:", err)
	dialog.Message("You chose file: %s", file).Title("Goodbye world!").Error()

	return true
}


func (me *Heartbeat) confirmDialog(t string, m string) bool {

	ok := dialog.Message("%s", m).Title(t).YesNo()

	return ok
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


func (me *Heartbeat) SetMenuState(menu menuStruct) (returnValue string) {
	// This can clearly be refactored a LOT.

	vmState := me.State.Box.GetStateMeaning()
	menu.apiStatusEntry.SetTooltip(vmState.ApiHintState)
	menu.apiStatusEntry.SetIcon(me.getIcon(vmState.ApiIconState))
	menu.apiStatusEntry.SetTitle(vmState.ApiTitleState)

	menu.vmStatusEntry.SetTooltip(vmState.VmHintState)
	menu.vmStatusEntry.SetIcon(me.getIcon(vmState.VmIconState))
	menu.vmStatusEntry.SetTitle(vmState.VmTitleState)

/*
	if me.State.Box.VM.LastSts != nil {
		menu.vmStatusEntry.SetTooltip(me.State.Box.VM.LastSts.Message())
	}

	switch {
		case me.State.Box.VM.CurrentState == box.VmStateNotPresent:
			menu.vmStatusEntry.SetIcon(me.getIcon(IconLogo))
			menu.vmStatusEntry.SetTitle("Box: VM not created")

		case me.State.Box.VM.CurrentState == box.VmStateUnknown:
			menu.vmStatusEntry.SetIcon(me.getIcon(IconError))
			menu.vmStatusEntry.SetTitle("Box: unknown error")

		case (me.State.Box.VM.CurrentState == box.VmStateRunning) && (me.State.Box.VM.WantState == box.VmStatePowerOff):
			menu.vmStatusEntry.SetIcon(me.getIcon(IconStopping))
			menu.vmStatusEntry.SetTitle("Box: stopping")

		case (me.State.Box.VM.CurrentState == box.VmStatePowerOff) && (me.State.Box.VM.WantState == box.VmStateRunning):
			menu.vmStatusEntry.SetIcon(me.getIcon(IconStarting))
			menu.vmStatusEntry.SetTitle("Box: starting")

		case (me.State.Box.VM.CurrentState == box.VmStateRunning) && (me.State.Box.VM.WantState == box.VmStateRunning):
			menu.vmStatusEntry.SetIcon(me.getIcon(IconUp))
			menu.vmStatusEntry.SetTitle("Box: running")

		case (me.State.Box.VM.CurrentState == box.VmStatePowerOff) && (me.State.Box.VM.WantState == box.VmStatePowerOff):
			menu.vmStatusEntry.SetIcon(me.getIcon(IconDown))
			menu.vmStatusEntry.SetTitle("Box: halted")

		default:
			menu.vmStatusEntry.SetIcon(me.getIcon(IconWarning))
			menu.vmStatusEntry.SetTitle("Box: unknown")
	}
*/

/*
	if me.State.Box.API.LastSts != nil {
		menu.apiStatusEntry.SetTooltip(me.State.Box.API.LastSts.Message())
	}

	switch {
		case me.State.Box.API.CurrentState == box.VmStateUnknown:
			menu.apiStatusEntry.SetIcon(me.getIcon(IconDown))
			menu.apiStatusEntry.SetTitle("API: halted")

		case (me.State.Box.API.CurrentState == box.VmStateRunning && me.State.Box.API.WantState == box.VmStatePowerOff) || (me.State.Box.API.CurrentState == box.VmStateStopping):
			menu.apiStatusEntry.SetIcon(me.getIcon(IconStopping))
			menu.apiStatusEntry.SetTitle("API: stopping")

		case (me.State.Box.API.CurrentState == box.VmStatePowerOff && me.State.Box.API.WantState == box.VmStateRunning) || (me.State.Box.API.CurrentState == box.VmStateStarting):
			menu.apiStatusEntry.SetIcon(me.getIcon(IconStarting))
			menu.apiStatusEntry.SetTitle("API: starting")

		case (me.State.Box.API.CurrentState == box.VmStateRunning) && (me.State.Box.API.WantState == box.VmStateRunning):
			menu.apiStatusEntry.SetIcon(me.getIcon(IconUp))
			menu.apiStatusEntry.SetTitle("API: running")

		case (me.State.Box.API.CurrentState == box.VmStatePowerOff) && (me.State.Box.API.WantState == box.VmStatePowerOff):
			menu.apiStatusEntry.SetIcon(me.getIcon(IconDown))
			menu.apiStatusEntry.SetTitle("API: halted")

		default:
			menu.apiStatusEntry.SetIcon(me.getIcon(IconWarning))
			menu.apiStatusEntry.SetTitle("API: unknown")
	}
*/

	fmt.Printf("me.State.Unfsd=%v\n", me.State.Unfsd)
	if me.State.Unfsd.LastSts != nil {
		menu.unfsdStatusEntry.SetTooltip(me.State.Unfsd.LastSts.Message())
	}
	switch {
		case me.State.Unfsd.CurrentState == monitor.StateUnknown:
			menu.unfsdStatusEntry.SetIcon(me.getIcon(IconError))
			menu.unfsdStatusEntry.SetTitle("FS: unknown error")

		case (me.State.Unfsd.CurrentState == monitor.StateRunning) && (me.State.Unfsd.WantState == monitor.StatePowerOff):
			menu.unfsdStatusEntry.SetIcon(me.getIcon(IconStopping))
			menu.unfsdStatusEntry.SetTitle("FS: stopping")

		case (me.State.Unfsd.CurrentState == monitor.StatePowerOff) && (me.State.Unfsd.WantState == monitor.StateRunning):
			menu.unfsdStatusEntry.SetIcon(me.getIcon(IconStarting))
			menu.unfsdStatusEntry.SetTitle("FS: starting")

		case (me.State.Unfsd.CurrentState == monitor.StateRunning) && (me.State.Unfsd.WantState == monitor.StateRunning):
			menu.unfsdStatusEntry.SetIcon(me.getIcon(IconUp))
			menu.unfsdStatusEntry.SetTitle("FS: running")

		case (me.State.Unfsd.CurrentState == monitor.StatePowerOff) && (me.State.Unfsd.WantState == monitor.StatePowerOff):
			menu.unfsdStatusEntry.SetIcon(me.getIcon(IconDown))
			menu.unfsdStatusEntry.SetTitle("FS: halted")

		default:
			menu.unfsdStatusEntry.SetIcon(me.getIcon(IconWarning))
			menu.unfsdStatusEntry.SetTitle("FS: unknown")
	}

	switch vmState.Name {
		case box.VmStateNotPresent:
			fmt.Printf("STATE: NOT PRESENT\n")
			systray.SetIcon(me.getIcon(IconWarning))
			systray.SetTooltip("Gearbox VM needs to be created.")

			returnValue = box.VmStateNotPresent
			menu.stopEntry.Disable()
			menu.startEntry.Disable()
			menu.sshEntry.Disable()
			menu.createEntry.Enable()

		case box.VmStateUnknown:
			fmt.Printf("STATE: UNKNOWN\n")
			systray.SetIcon(me.getIcon(IconWarning))
			systray.SetTooltip("Gearbox is in an unknown state.")

			returnValue = box.VmStateUnknown
			menu.stopEntry.Disable()
			menu.startEntry.Disable()
			menu.sshEntry.Disable()
			menu.createEntry.Enable()

		case box.VmStatePaused:
			fallthrough
		case box.VmStateSaved:
			fallthrough
		case box.VmStatePowerOff:
			// fmt.Printf("STATE: HALTED\n")
			systray.SetIcon(me.getIcon(IconDown))
			systray.SetTooltip("Gearbox is halted.")

			returnValue = box.VmStatePowerOff
			menu.stopEntry.Disable()
			menu.startEntry.Enable()
			menu.sshEntry.Disable()
			menu.createEntry.Disable()

		case box.VmStateRunning:
			// fmt.Printf("STATE: RUNNING\n")
			systray.SetIcon(me.getIcon(IconUp))
			systray.SetTooltip("Gearbox is running.")

			returnValue = box.VmStateRunning
			menu.stopEntry.Enable()
			menu.startEntry.Disable()
			menu.sshEntry.Enable()
			menu.createEntry.Disable()

		case box.VmStateStarting:
			fmt.Printf("STATE: STARTING\n")
			systray.SetIcon(me.getIcon(IconStarting))
			systray.SetTooltip("Gearbox starting up.")

			returnValue = box.VmStateStarting
			menu.stopEntry.Disable()
			menu.startEntry.Disable()
			menu.sshEntry.Disable()
			menu.createEntry.Disable()

		case box.VmStateStopping:
			fmt.Printf("STATE: STOPPING\n")
			systray.SetIcon(me.getIcon(IconStopping))
			systray.SetTooltip("Gearbox is stopping.")

			returnValue = box.VmStateStopping
			menu.stopEntry.Disable()
			menu.startEntry.Disable()
			menu.sshEntry.Disable()
			menu.createEntry.Disable()

	}

	return
}


func (me *Heartbeat) Initialize() (sts status.Status) {

	//var boxArgs box.Args

	for range only.Once {

		cfgdir := me.OsSupport.GetUserConfigDir()

		// sts := gearbox.BoxInstance.Start(boxArgs)

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
