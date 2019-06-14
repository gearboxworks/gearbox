package heartbeat

import (
	"errors"
	"fmt"
	"gearbox/box"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/only"
	"github.com/gearboxworks/go-status/is"
	"github.com/getlantern/systray"
	"github.com/sqweek/dialog"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)


func (me *Heartbeat) CreateMenus() {

	systray.SetIcon(me.getIcon(DefaultLogo))
	systray.SetTitle("")

	me.menu = make(Menus)

	me.menu["help"] = &Menu{
		MenuItem: systray.AddMenuItem("About Gearbox", "Contact Gearbox help for"+me.BoxInstance.Boxname),
		PrefixToolTip: "",
		PrefixMenu: "",
		CurrentIcon: "",
	}

	me.menu["version"] = &Menu{
		MenuItem: systray.AddMenuItem(fmt.Sprintf("Gearbox (v%s)", me.BoxInstance.VmIsoVersion), "Running v"+me.BoxInstance.VmIsoVersion),
		PrefixToolTip: "",
		PrefixMenu: "",
		CurrentIcon: "",
	}


	systray.AddSeparator()


	me.menu["vm"] = &Menu{
		MenuItem: systray.AddMenuItem("Gearbox OS: Idle", "Current state of Gearbox VM"),
		PrefixToolTip: "",
		PrefixMenu: "Gearbox OS: ",
		CurrentIcon: DefaultLogo,
	}
	me.menu["vm"].MenuItem.SetIcon(me.getIcon(DefaultLogo))

	me.menu["api"] = &Menu{
		MenuItem: systray.AddMenuItem("Gearbox API: Idle", "Current state of Gearbox API"),
		PrefixToolTip: "",
		PrefixMenu: "Gearbox API: ",
		CurrentIcon: DefaultLogo,
	}
	me.menu["api"].MenuItem.SetIcon(me.getIcon(me.menu["api"].CurrentIcon))

	me.menu["unfsd"] = &Menu{
		MenuItem: systray.AddMenuItem("Gearbox FS: Idle", "Current state of Gearbox NFS service"),
		PrefixToolTip: "",
		PrefixMenu: "Gearbox FS: ",
		CurrentIcon: DefaultLogo,
	}
	me.menu["unfsd"].MenuItem.SetIcon(me.getIcon(me.menu["unfsd"].CurrentIcon))


	systray.AddSeparator()


	me.menu["admin"] = &Menu{
		MenuItem: systray.AddMenuItem("Admin", "Open Gearbox admin interface"),
		PrefixToolTip: "",
		PrefixMenu: "",
		CurrentIcon: "",
	}

	me.menu["create"] = &Menu{
		MenuItem: systray.AddMenuItem("Create Box", "Create a Gearbox OS instance"),
		PrefixToolTip: "",
		PrefixMenu: "",
		CurrentIcon: "",
	}

	me.menu["update"] = &Menu{
		MenuItem: systray.AddMenuItem("Update Box", "Check for Gearbox OS updates"),
		PrefixToolTip: "",
		PrefixMenu: "",
		CurrentIcon: "",
	}

	me.menu["start"] = &Menu{
		MenuItem: systray.AddMenuItem("Start Box", "Start Gearbox OS instance"),
		PrefixToolTip: "",
		PrefixMenu: "",
		CurrentIcon: "",
	}

	me.menu["stop"] = &Menu{
		MenuItem: systray.AddMenuItem("Stop Box", "Stop Gearbox OS instance"),
		PrefixToolTip: "",
		PrefixMenu: "",
		CurrentIcon: "",
	}

	me.menu["ssh"] = &Menu{
		MenuItem: systray.AddMenuItem("SSH", "Connect to Gearbox OS via SSH"),
		PrefixToolTip: "",
		PrefixMenu: "",
		CurrentIcon: "",
	}


	systray.AddSeparator()


	pid := os.Getpid()
	me.menu["restart"] = &Menu{
		MenuItem: systray.AddMenuItem("Restart Heartbeat", fmt.Sprintf("Restart this app [pid:%v]", pid)),
		PrefixToolTip: "",
		PrefixMenu: "",
		CurrentIcon: "",
	}

	me.menu["quit"] = &Menu{
		MenuItem: systray.AddMenuItem("Quit", fmt.Sprintf("Terminate this app [pid:%v]", pid)),
		PrefixToolTip: "",
		PrefixMenu: "",
		CurrentIcon: "",
	}

}


func (me *Heartbeat) UpdateMenus() {

	s, err := me.EventBroker.GetSimpleStatus()
	if err != nil {
		return
	}

	for k, v := range s {
		me.SetMenu(k, v)
	}

	control := messages.MessageAddresses{"admin", "create", "update", "start", "stop", "ssh"}
	for _, v := range control {
		if me.menu.Exists(v) {
			me.SetMenu(v, states.ActionIdle)
		}

	}


	// admin
	// create
	// update
	// start
	// stop
	// ssh


	//me.SetMenu("api", "")
	//me.SetMenu("unfsd", "")

}


func (me Menus) EnsureNotNil() error {

	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("oops")
			break
		}
	}

	return err
}


func (me Menus) Exists(item messages.MessageAddress) bool {

	var ret bool

	if _, ok := me[item]; ok {
		ret = true
	}

	return ret
}


func (me *Menu) Disable() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.MenuItem.Disable()
	}

	return err
}


func (me *Menu) Enable() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.MenuItem.Enable()
	}

	return err
}


func (me *Menu) EnsureNotNil() error {

	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("oops")
			break
		}

		if me.MenuItem == nil {
			err = errors.New("oops")
			break
		}
	}

	return err
}


//func (me *Heartbeat) SetMenuState() (returnValue string) {
//	// This can clearly be refactored a LOT.
//
//	if me.BoxInstance.VmIsoDlIndex < 100 {
//		me.menu.versionEntry.SetTitle(fmt.Sprintf("Gearbox (v%s) updating ...", me.BoxInstance.VmIsoVersion))
//		me.menu.versionEntry.SetTooltip("Updating v" + me.BoxInstance.VmIsoVersion)
//
//		me.menu.updateEntry.SetTitle(fmt.Sprintf("Updating Box (%d%%)", me.BoxInstance.VmIsoDlIndex))
//		me.menu.updateEntry.SetTooltip(fmt.Sprintf("Downloading v%s - %d%% complete", me.BoxInstance.VmIsoVersion, me.BoxInstance.VmIsoDlIndex))
//		me.menu.updateEntry.Disable()
//
//	} else {
//		me.menu.versionEntry.SetTitle(fmt.Sprintf("Gearbox (v%s)", me.BoxInstance.VmIsoVersion))
//		me.menu.versionEntry.SetTooltip("Running v" + me.BoxInstance.VmIsoVersion)
//
//		me.menu.updateEntry.SetTitle("Update Box")
//		me.menu.updateEntry.SetTooltip("Check for Gearbox OS updates")
//		me.menu.updateEntry.Enable()
//	}
//	// 		sts = me.IsIsoFilePresent()
//	//		if is.Error(sts) {
//	//			break
//	//		}
//
//	vmState := me.State.Box.GetStateMeaning()
//	me.menu.apiStatusEntry.SetTooltip(vmState.ApiHintState)
//	me.menu.apiStatusEntry.SetIcon(me.getIcon(vmState.ApiIconState))
//	me.menu.apiStatusEntry.SetTitle(vmState.ApiTitleState)
//
//	me.menu.vmStatusEntry.SetTooltip(vmState.VmHintState)
//	me.menu.vmStatusEntry.SetIcon(me.getIcon(vmState.VmIconState))
//	me.menu.vmStatusEntry.SetTitle(vmState.VmTitleState)
//
//	//if me.State.Box.VM.LastSts != nil {
//	//	me.menu.vmStatusEntry.SetTooltip(me.State.Box.VM.LastSts.Message())
//	//}
//	//
//	//switch {
//	//	case me.State.Box.VM.CurrentState == box.VmStateNotPresent:
//	//		me.menu.vmStatusEntry.SetIcon(me.getIcon(IconLogo))
//	//		me.menu.vmStatusEntry.SetTitle("Box: VM not created")
//	//
//	//	case me.State.Box.VM.CurrentState == box.VmStateUnknown:
//	//		me.menu.vmStatusEntry.SetIcon(me.getIcon(IconError))
//	//		me.menu.vmStatusEntry.SetTitle("Box: unknown error")
//	//
//	//	case (me.State.Box.VM.CurrentState == box.VmStateRunning) && (me.State.Box.VM.WantState == box.VmStatePowerOff):
//	//		me.menu.vmStatusEntry.SetIcon(me.getIcon(IconStopping))
//	//		me.menu.vmStatusEntry.SetTitle("Box: stopping")
//	//
//	//	case (me.State.Box.VM.CurrentState == box.VmStatePowerOff) && (me.State.Box.VM.WantState == box.VmStateRunning):
//	//		me.menu.vmStatusEntry.SetIcon(me.getIcon(IconStarting))
//	//		me.menu.vmStatusEntry.SetTitle("Box: starting")
//	//
//	//	case (me.State.Box.VM.CurrentState == box.VmStateRunning) && (me.State.Box.VM.WantState == box.VmStateRunning):
//	//		me.menu.vmStatusEntry.SetIcon(me.getIcon(IconUp))
//	//		me.menu.vmStatusEntry.SetTitle("Box: running")
//	//
//	//	case (me.State.Box.VM.CurrentState == box.VmStatePowerOff) && (me.State.Box.VM.WantState == box.VmStatePowerOff):
//	//		me.menu.vmStatusEntry.SetIcon(me.getIcon(IconDown))
//	//		me.menu.vmStatusEntry.SetTitle("Box: halted")
//	//
//	//	default:
//	//		me.menu.vmStatusEntry.SetIcon(me.getIcon(IconWarning))
//	//		me.menu.vmStatusEntry.SetTitle("Box: unknown")
//	//}
//	//
//	//if me.State.Box.API.LastSts != nil {
//	//	me.menu.apiStatusEntry.SetTooltip(me.State.Box.API.LastSts.Message())
//	//}
//	//
//	//switch {
//	//	case me.State.Box.API.CurrentState == box.VmStateUnknown:
//	//		me.menu.apiStatusEntry.SetIcon(me.getIcon(IconDown))
//	//		me.menu.apiStatusEntry.SetTitle("API: halted")
//	//
//	//	case (me.State.Box.API.CurrentState == box.VmStateRunning && me.State.Box.API.WantState == box.VmStatePowerOff) || (me.State.Box.API.CurrentState == box.VmStateStopping):
//	//		me.menu.apiStatusEntry.SetIcon(me.getIcon(IconStopping))
//	//		me.menu.apiStatusEntry.SetTitle("API: stopping")
//	//
//	//	case (me.State.Box.API.CurrentState == box.VmStatePowerOff && me.State.Box.API.WantState == box.VmStateRunning) || (me.State.Box.API.CurrentState == box.VmStateStarting):
//	//		me.menu.apiStatusEntry.SetIcon(me.getIcon(IconStarting))
//	//		me.menu.apiStatusEntry.SetTitle("API: starting")
//	//
//	//	case (me.State.Box.API.CurrentState == box.VmStateRunning) && (me.State.Box.API.WantState == box.VmStateRunning):
//	//		me.menu.apiStatusEntry.SetIcon(me.getIcon(IconUp))
//	//		me.menu.apiStatusEntry.SetTitle("API: running")
//	//
//	//	case (me.State.Box.API.CurrentState == box.VmStatePowerOff) && (me.State.Box.API.WantState == box.VmStatePowerOff):
//	//		me.menu.apiStatusEntry.SetIcon(me.getIcon(IconDown))
//	//		me.menu.apiStatusEntry.SetTitle("API: halted")
//	//
//	//	default:
//	//		me.menu.apiStatusEntry.SetIcon(me.getIcon(IconWarning))
//	//		me.menu.apiStatusEntry.SetTitle("API: unknown")
//	//}
//	////fmt.Printf("me.State.Unfsd=%v\n", me.State.Unfsd)
//	//if me.State.Unfsd.LastSts != nil {
//	//	me.menu.unfsdStatusEntry.SetTooltip(me.State.Unfsd.LastSts.Message())
//	//}
//	//switch {
//	//	case me.State.Unfsd.CurrentState == external.StateUnknown:
//	//		me.menu.unfsdStatusEntry.SetIcon(me.getIcon(IconError))
//	//		me.menu.unfsdStatusEntry.SetTitle("FS: unknown error")
//	//
//	//	case (me.State.Unfsd.CurrentState == external.StateRunning) && (me.State.Unfsd.WantState == external.StatePowerOff):
//	//		me.menu.unfsdStatusEntry.SetIcon(me.getIcon(IconStopping))
//	//		me.menu.unfsdStatusEntry.SetTitle("FS: stopping")
//	//
//	//	case (me.State.Unfsd.CurrentState == external.StatePowerOff) && (me.State.Unfsd.WantState == external.StateRunning):
//	//		me.menu.unfsdStatusEntry.SetIcon(me.getIcon(IconStarting))
//	//		me.menu.unfsdStatusEntry.SetTitle("FS: starting")
//	//
//	//	case (me.State.Unfsd.CurrentState == external.StateRunning) && (me.State.Unfsd.WantState == external.StateRunning):
//	//		me.menu.unfsdStatusEntry.SetIcon(me.getIcon(IconUp))
//	//		me.menu.unfsdStatusEntry.SetTitle("FS: running")
//	//
//	//	case (me.State.Unfsd.CurrentState == external.StatePowerOff) && (me.State.Unfsd.WantState == external.StatePowerOff):
//	//		me.menu.unfsdStatusEntry.SetIcon(me.getIcon(IconDown))
//	//		me.menu.unfsdStatusEntry.SetTitle("FS: halted")
//	//
//	//	default:
//	//		me.menu.unfsdStatusEntry.SetIcon(me.getIcon(IconWarning))
//	//		me.menu.unfsdStatusEntry.SetTitle("FS: unknown")
//	//}
//
//	switch vmState.Name {
//	case box.VmStateNotPresent:
//		fmt.Printf("STATE: NOT PRESENT\n")
//		systray.SetIcon(me.getIcon(IconWarning))
//		systray.SetTooltip("Gearbox VM needs to be created.")
//
//		returnValue = box.VmStateNotPresent
//		me.menu.stopEntry.Hide()
//		me.menu.startEntry.Hide()
//		me.menu.sshEntry.Hide()
//		me.menu.createEntry.Show()
//
//	case box.VmStateUnknown:
//		fmt.Printf("STATE: UNKNOWN\n")
//		systray.SetIcon(me.getIcon(IconWarning))
//		systray.SetTooltip("Gearbox is in an unknown state.")
//
//		returnValue = box.VmStateUnknown
//		me.menu.stopEntry.Hide()
//		me.menu.startEntry.Hide()
//		me.menu.sshEntry.Hide()
//		me.menu.createEntry.Show()
//
//	case box.VmStatePaused:
//		fallthrough
//	case box.VmStateSaved:
//		fallthrough
//	case box.VmStatePowerOff:
//		// fmt.Printf("STATE: HALTED\n")
//		systray.SetIcon(me.getIcon(IconDown))
//		systray.SetTooltip("Gearbox is halted.")
//
//		returnValue = box.VmStatePowerOff
//		me.menu.stopEntry.Hide()
//		me.menu.startEntry.Show()
//		me.menu.sshEntry.Hide()
//		me.menu.createEntry.Hide()
//
//	case box.VmStateRunning:
//		// fmt.Printf("STATE: RUNNING\n")
//		systray.SetIcon(me.getIcon(IconUp))
//		systray.SetTooltip("Gearbox is running.")
//
//		returnValue = box.VmStateRunning
//		me.menu.stopEntry.Show()
//		me.menu.startEntry.Hide()
//		me.menu.sshEntry.Show()
//		me.menu.createEntry.Hide()
//
//	case box.VmStateStarting:
//		fmt.Printf("STATE: STARTING\n")
//		systray.SetIcon(me.getIcon(IconStarting))
//		systray.SetTooltip("Gearbox starting up.")
//
//		returnValue = box.VmStateStarting
//		me.menu.stopEntry.Hide()
//		me.menu.startEntry.Hide()
//		me.menu.sshEntry.Hide()
//		me.menu.createEntry.Hide()
//
//	case box.VmStateStopping:
//		fmt.Printf("STATE: STOPPING\n")
//		systray.SetIcon(me.getIcon(IconStopping))
//		systray.SetTooltip("Gearbox is stopping.")
//
//		returnValue = box.VmStateStopping
//		me.menu.stopEntry.Hide()
//		me.menu.startEntry.Hide()
//		me.menu.sshEntry.Hide()
//		me.menu.createEntry.Hide()
//
//	}
//
//	return
//}


func (me *Heartbeat) onReady() {

	//var intentDelay = false
	// Used to change delay times when the user has just performed an action.
	//sts := me.NfsInstance.Daemon.Load()
	//if is.Error(sts) {
	//	fmt.Printf("%s\n", sts.Message())
	//	return
	//}
	// Concurrent process: Provide status updates on systray.
	// Ideally, this should also send messages on message bus for actions to be taken. EG: Retry startup, disk full, etc.
	// Even further, these should be brokem out into methods to avoid having to hard code specific entities to monitor.
	//go func() {
	//	var state State
	//	var sts status.Status
	//
	//	for {
	//		if intentDelay {
	//			// User has requested a change, check on cached results faster.
	//			// results will be updated by concurrent functions.
	//			//fmt.Printf("CACHE POLL\n")
	//
	//			// Check state of VM.
	//			me.State.Box, sts = me.BoxInstance.GetCachedState()
	//			if is.Error(sts) {
	//				// .
	//			}
	//
	//			// Check state of UNFSD.
	//			//me.State.Unfsd, sts = me.NfsInstance.GetState()
	//			if is.Error(sts) || is.Error(state.Unfsd.LastSts) {
	//				// .
	//			}
	//
	//			me.SetMenuState()
	//			time.Sleep(time.Second)
	//
	//		} else {
	//			// Normal polling.
	//			//fmt.Printf("NORMAL POLL\n")
	//
	//			// Check state of VM.
	//			me.State.Box, sts = me.BoxInstance.GetState()
	//			//fmt.Printf("STATE:\n%v\n%v\n", me.State.Box, sts)
	//			if is.Error(sts) {
	//				// .
	//			}
	//
	//			// Check state of UNFSD.
	//			//me.State.Unfsd, sts = me.NfsInstance.GetState()
	//			if is.Error(sts) || is.Error(state.Unfsd.LastSts) {
	//				// .
	//			}
	//
	//			if me.BoxInstance.VmIsoDlIndex == 100 {
	//				sts = me.BoxInstance.IsIsoFilePresent()
	//				if !is.Success(sts) {
	//					fmt.Printf("Get ready agent: %v\n", sts)
	//					fmt.Printf("Downloading...\n")
	//					me.BoxInstance.VmIsoDlIndex = 0
	//					intentDelay = true
	//					go me.BoxInstance.GetIso()
	//					intentDelay = false
	//					// var b struct{}
	//					// me.menu.updateEntry.ClickedCh <- b
	//				}
	//			}
	//
	//			me.SetMenuState()
	//			time.Sleep(10 * time.Second)
	//		}
	//	}
	//}()

	me.CreateMenus()
	me.UpdateMenus()

	go func() {
		for {
		select {
			case <- me.menu["help"].MenuItem.ClickedCh:
				fmt.Printf("Menu: Help\n")
				me.openAbout()

			case <- me.menu["version"].MenuItem.ClickedCh:
				fmt.Printf("Menu: Version\n")

			case <- me.menu["vm"].MenuItem.ClickedCh:
				// Ignore.
			case <- me.menu["api"].MenuItem.ClickedCh:
				// Ignore.
			case <- me.menu["unfsd"].MenuItem.ClickedCh:
				// Ignore.

			case <- me.menu["start"].MenuItem.ClickedCh:
				fmt.Printf("Menu: Start\n")
				//intentDelay = true
				me.BoxInstance.Start()
				//intentDelay = false

			case <- me.menu["stop"].MenuItem.ClickedCh:
				fmt.Printf("Menu: Stop\n")
				//intentDelay = true
				me.BoxInstance.Stop()
				//intentDelay = false

			case <- me.menu["admin"].MenuItem.ClickedCh:
				fmt.Printf("Menu: Admin\n")
				me.openAdmin()

			case <- me.menu["ssh"].MenuItem.ClickedCh:
				fmt.Printf("Menu: SSH\n")
				me.openTerminal()

			case <- me.menu["create"].MenuItem.ClickedCh:
				fmt.Printf("Menu: Create\n")
				//intentDelay = true
				if me.BoxInstance.State.VM.CurrentState == box.VmStateNotPresent {
					sts := me.BoxInstance.CreateBox()
					if is.Error(sts) {
						dialog.Message("Error! Creating Gearbox OS VM: %s", me.Boxname).Title("GearBox OS Creation").Error()
					} else {
						dialog.Message("Success! Gearbox OS VM created: %s", me.Boxname).Title("GearBox OS Creation").Info()
					}
				}
				//intentDelay = false

			case <- me.menu["update"].MenuItem.ClickedCh:
				fmt.Printf("Menu: Update\n")
				if me.BoxInstance.VmIsoDlIndex == 100 {
					me.BoxInstance.VmIsoDlIndex = 0
					//intentDelay = true
					go me.BoxInstance.GetIso()
					//intentDelay = false
				}

			case <- me.menu["restart"].MenuItem.ClickedCh:
				fmt.Printf("Menu: Restart\n")
				if me.confirmDialog("Restart Gearbox", "This will restart Gearbox Heartbeat, but keep services running.\nAre you sure?") {
					fmt.Printf("HEY!")
					systray.Quit()
				}

			case <- me.menu["quit"].MenuItem.ClickedCh:
				fmt.Printf("Menu: Quit\n")
				if me.confirmDialog("Shutdown Gearbox", "This will shutdown Gearbox and all Gearbox related services.\nAre you sure?") {
					//intentDelay = true
					me.BoxInstance.Stop()
					//me.NfsInstance.Stop()
					//intentDelay = false

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

	fp := filepath.FromSlash(me.baseDir + "/" + s)
	if fp == "" {
		return nil
	}

	b, err := ioutil.ReadFile(fp)
	if err != nil {
		fmt.Print(err)
	}

	return b
}
