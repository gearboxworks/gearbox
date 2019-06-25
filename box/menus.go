package box

import (
	"errors"
	"fmt"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/states"
	"gearbox/box/external/vmbox"
	"gearbox/only"
	"github.com/getlantern/systray"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)
// Disabled because of GOOS=windows & GOOS=linux
// 	"github.com/sqweek/dialog"


const (
	menuVmAdmin  = "admin"
	menuVmCreate = "create"
	menuVmUpdate = "update"
	menuVmStart  = "start"
	menuVmStop   = "stop"
	menuVmSsh    = "ssh"
)


func (me *Box) CreateMenus() {

	systray.SetIcon(me.getIcon(DefaultLogo))
	systray.SetTitle("")

	me.menu = make(Menus)

	me.menu["help"] = &Menu{
		MenuItem: systray.AddMenuItem("About Gearbox", "Contact Gearbox help for"+me.Boxname),
		PrefixToolTip: "",
		PrefixMenu: "",
		CurrentIcon: "",
	}

	me.menu["version"] = &Menu{
		MenuItem: systray.AddMenuItem("Gearbox (v0.5.9)", "Running v0.5.0"),
		PrefixToolTip: "",
		PrefixMenu: "",
		CurrentIcon: "",
	}


	systray.AddSeparator()


	me.menu[entity.VmEntityName] = &Menu{
		MenuItem: systray.AddMenuItem("Gearbox OS: Idle", "Current state of Gearbox VM"),
		PrefixToolTip: "",
		PrefixMenu: "Gearbox OS: ",
		CurrentIcon: DefaultLogo,
	}
	me.menu[entity.VmEntityName].MenuItem.SetIcon(me.getIcon(me.menu[entity.VmEntityName].CurrentIcon))

	me.menu[entity.ApiEntityName] = &Menu{
		MenuItem: systray.AddMenuItem("Gearbox API: Idle", "Current state of Gearbox API"),
		PrefixToolTip: "",
		PrefixMenu: "Gearbox API: ",
		CurrentIcon: DefaultLogo,
	}
	me.menu[entity.ApiEntityName].MenuItem.SetIcon(me.getIcon(me.menu[entity.ApiEntityName].CurrentIcon))

	me.menu[entity.UnfsdEntityName] = &Menu{
		MenuItem: systray.AddMenuItem("Gearbox FS: Idle", "Current state of Gearbox NFS service"),
		PrefixToolTip: "",
		PrefixMenu: "Gearbox FS: ",
		CurrentIcon: DefaultLogo,
	}
	me.menu[entity.UnfsdEntityName].MenuItem.SetIcon(me.getIcon(me.menu[entity.UnfsdEntityName].CurrentIcon))


	systray.AddSeparator()


	me.menu[menuVmAdmin] = &Menu{
		MenuItem: systray.AddMenuItem("Admin", "Open Gearbox admin interface"),
		PrefixToolTip: "",
		PrefixMenu: "Admin",
		CurrentIcon: "",
	}

	me.menu[menuVmCreate] = &Menu{
		MenuItem: systray.AddMenuItem("Create Box", "Create a Gearbox OS instance"),
		PrefixToolTip: "",
		PrefixMenu: "Create Box",
		CurrentIcon: "",
	}

	me.menu[menuVmUpdate] = &Menu{
		MenuItem: systray.AddMenuItem("Update Box", "Check for Gearbox OS updates"),
		PrefixToolTip: "",
		PrefixMenu: "Update Box",
		CurrentIcon: "",
	}

	me.menu[menuVmStart] = &Menu{
		MenuItem: systray.AddMenuItem("Start Box", "Start Gearbox OS instance"),
		PrefixToolTip: "",
		PrefixMenu: "Start Box",
		CurrentIcon: "",
	}

	me.menu[menuVmStop] = &Menu{
		MenuItem: systray.AddMenuItem("Stop Box", "Stop Gearbox OS instance"),
		PrefixToolTip: "",
		PrefixMenu: "Stop Box",
		CurrentIcon: "",
	}

	me.menu[menuVmSsh] = &Menu{
		MenuItem: systray.AddMenuItem("SSH", "Connect to Gearbox OS via SSH"),
		PrefixToolTip: "",
		PrefixMenu: "SSH",
		CurrentIcon: "",
	}


	systray.AddSeparator()


	pid := os.Getpid()
	//me.menu["restart"] = &Menu{
	//	MenuItem: systray.AddMenuItem("Restart Box", fmt.Sprintf("Restart this app [pid:%v]", pid)),
	//	PrefixToolTip: "",
	//	PrefixMenu: "",
	//	CurrentIcon: "",
	//}

	me.menu["quit"] = &Menu{
		MenuItem: systray.AddMenuItem("Quit", fmt.Sprintf("Terminate this app [pid:%v]", pid)),
		PrefixToolTip: "",
		PrefixMenu: "",
		CurrentIcon: "",
	}

}


func (me *Box) UpdateMenus() {

	s, err := me.EventBroker.GetSimpleStatus()
	if err != nil {
		return
	}

	me.SetStateMenu(entity.VmEntityName, s[entity.VmEntityName])
	me.SetStateMenu(entity.ApiEntityName, s[entity.ApiEntityName])
	me.SetStateMenu(entity.UnfsdEntityName, s[entity.UnfsdEntityName])
	me.SetControlMenu(entity.VmEntityName, s[entity.VmEntityName])

	//for k, v := range s {
	//	me.SetStateMenu(k, v)
	//}
	//
	//control := messages.MessageAddresses{"admin", "create", "update", "start", "stop", "ssh"}
	//for _, v := range control {
	//	if me.menu.Exists(v) {
	//		me.SetControlMenu(v, states.ActionIdle)
	//	}
	//
	//}
	//
	//me.SetMenu("api", "")
	//me.SetMenu("unfsd", "")

}


func (me *Box) SetStateMenu(m messages.MessageAddress, state states.State) {
	// This can clearly be refactored a LOT.

	if _, ok := me.menu[m]; !ok {
		return
	}

	if me.menu[m].MenuItem == nil {
		return
	}

	mi := me.menu[m]
	switch state {
		case states.StateUnknown:
			mi.MenuItem.SetIcon(me.getIcon(IconError))

		case states.StateStopping:
			mi.MenuItem.SetIcon(me.getIcon(IconStopping))

		case states.StateStarting:
			mi.MenuItem.SetIcon(me.getIcon(IconStarting))

		case states.StateStarted:
			mi.MenuItem.SetIcon(me.getIcon(IconUp))

		case states.StateStopped:
			mi.MenuItem.SetIcon(me.getIcon(IconDown))

		default:
			mi.MenuItem.SetIcon(me.getIcon(IconWarning))
	}
	mi.MenuItem.SetTitle(mi.PrefixMenu + state.String())
	mi.MenuItem.SetTooltip(mi.PrefixToolTip + state.String())

	return
}


func (me *Box) SetControlMenu(m messages.MessageAddress, state states.State) {
	// This can clearly be refactored a LOT.

	if _, ok := me.menu[m]; !ok {
		return
	}

	if me.menu[m].MenuItem == nil {
		return
	}

	if m != entity.VmEntityName {
		return
	}

	// admin
	// create
	// update
	// start
	// stop
	// ssh

	mi := me.menu[m]
	switch state {
		case states.StateIdle:
			_ = me.menu[menuVmAdmin].Enable()
			_ = me.menu[menuVmCreate].Disable()
			_ = me.menu[menuVmUpdate].Enable()
			_ = me.menu[menuVmStart].Disable()
			_ = me.menu[menuVmStop].Disable()
			_ = me.menu[menuVmSsh].Disable()

		case states.StateUnknown:
			_ = me.menu[menuVmAdmin].Enable()
			_ = me.menu[menuVmCreate].Enable()
			_ = me.menu[menuVmUpdate].Enable()
			_ = me.menu[menuVmStart].Disable()
			_ = me.menu[menuVmStop].Disable()
			_ = me.menu[menuVmSsh].Disable()

		case states.StateStopping:
			_ = me.menu[menuVmAdmin].Enable()
			_ = me.menu[menuVmCreate].Disable()
			_ = me.menu[menuVmUpdate].Disable()
			_ = me.menu[menuVmStart].Disable()
			_ = me.menu[menuVmStop].Disable()
			_ = me.menu[menuVmSsh].Disable()

		case states.StateStarting:
			_ = me.menu[menuVmAdmin].Enable()
			_ = me.menu[menuVmCreate].Disable()
			_ = me.menu[menuVmUpdate].Disable()
			_ = me.menu[menuVmStart].Disable()
			_ = me.menu[menuVmStop].Disable()
			_ = me.menu[menuVmSsh].Disable()

		case states.StateStarted:
			_ = me.menu[menuVmAdmin].Enable()
			_ = me.menu[menuVmCreate].Disable()
			_ = me.menu[menuVmUpdate].Disable()
			_ = me.menu[menuVmStart].Disable()
			_ = me.menu[menuVmStop].Enable()
			_ = me.menu[menuVmSsh].Enable()

		case states.StateStopped:
			_ = me.menu[menuVmAdmin].Enable()
			_ = me.menu[menuVmCreate].Disable()
			_ = me.menu[menuVmUpdate].Disable()
			_ = me.menu[menuVmStart].Enable()
			_ = me.menu[menuVmStop].Disable()
			_ = me.menu[menuVmSsh].Disable()

		case states.StateUnregistered:
			_ = me.menu[menuVmAdmin].Enable()
			_ = me.menu[menuVmCreate].Enable()
			_ = me.menu[menuVmUpdate].Enable()
			_ = me.menu[menuVmStart].Disable()
			_ = me.menu[menuVmStop].Disable()
			_ = me.menu[menuVmSsh].Disable()

		case states.StateUpdating:
			_ = me.menu[menuVmAdmin].Enable()
			_ = me.menu[menuVmCreate].Disable()
			_ = me.menu[menuVmUpdate].Enable()
			_ = me.menu[menuVmStart].Disable()
			_ = me.menu[menuVmStop].Disable()
			_ = me.menu[menuVmSsh].Disable()

		default:
			_ = me.menu[menuVmAdmin].Enable()
			_ = me.menu[menuVmCreate].Enable()
			_ = me.menu[menuVmUpdate].Enable()
			_ = me.menu[menuVmStart].Disable()
			_ = me.menu[menuVmStop].Disable()
			_ = me.menu[menuVmSsh].Disable()
	}
	mi.MenuItem.SetTitle(mi.PrefixMenu + state.String())
	mi.MenuItem.SetTooltip(mi.PrefixToolTip + state.String())

	return
}


func (me *Box) onReady() {

	me.CreateMenus()
	me.UpdateMenus()

	go func() {
		for {
			select {
				case <- me.menu["help"].MenuItem.ClickedCh:
					fmt.Printf("Menu: Help.\n")
					me.openAbout()

				case <- me.menu["version"].MenuItem.ClickedCh:
					fmt.Printf("Menu: Version\n")


				case <- me.menu[entity.VmEntityName].MenuItem.ClickedCh:
					// Ignore.
				case <- me.menu[entity.ApiEntityName].MenuItem.ClickedCh:
					// Ignore.
				case <- me.menu[entity.UnfsdEntityName].MenuItem.ClickedCh:
					// Ignore.


				case <- me.menu[menuVmStart].MenuItem.ClickedCh:
					fmt.Printf("Menu: Start VM.\n")
					msg := vmbox.ConstructVmMessage(entity.VmBoxEntityName, entity.VmEntityName, states.ActionStart)
					_ = me.EventBroker.Channels.Publish(msg)

				case <- me.menu[menuVmStop].MenuItem.ClickedCh:
					fmt.Printf("Menu: Stop VM.\n")
					msg := vmbox.ConstructVmMessage(entity.VmBoxEntityName, entity.VmEntityName, states.ActionStop)
					_ = me.EventBroker.Channels.Publish(msg)

				case <- me.menu[menuVmAdmin].MenuItem.ClickedCh:
					fmt.Printf("Menu: Admin console.\n")
					me.openAdmin()

				case <- me.menu[menuVmSsh].MenuItem.ClickedCh:
					fmt.Printf("Menu: SSH\n")
					me.openTerminal()

				case <- me.menu[menuVmCreate].MenuItem.ClickedCh:
					fmt.Printf("Menu: Create VM.\n")
					msg := vmbox.ConstructVmMessage(entity.VmBoxEntityName, entity.VmEntityName, states.ActionRegister)
					_ = me.EventBroker.Channels.Publish(msg)

				case <- me.menu[menuVmUpdate].MenuItem.ClickedCh:
					fmt.Printf("Menu: Update VM ISO.\n")
					msg := vmbox.ConstructVmMessage(entity.VmBoxEntityName, entity.VmEntityName, states.ActionUpdate)
					_ = me.EventBroker.Channels.Publish(msg)


				//case <- me.menu["restart"].MenuItem.ClickedCh:
				//	fmt.Printf("Menu: Restart VM.\n")
				//	if me.confirmDialog("Restart Gearbox", "This will restart Gearbox Box, but keep services running.\nAre you sure?") {
				//		fmt.Printf("Shutting down!")
				//		systray.Quit()
				//	}

				case <- me.menu["quit"].MenuItem.ClickedCh:
					fmt.Printf("Menu: Quit.\n")
					fmt.Printf("Gearbox: Shutting down. (May take up to 2 minutes.)\n")
					if me.confirmDialog("Shutdown Gearbox", "This will shutdown Gearbox and all Gearbox related services.\nAre you sure?") {
						_ = me.VmBox.Stop()
						_ = me.EventBroker.Stop()
						_ = me.StopBox()

						systray.Quit()
					}
			}
		}
	}()

}


func (me *Box) fileDialog(t string, m string) bool {
	//dialog.Message("%s", "Please select a file").Title("Hello world!").Info()
	//file, err := dialog.File().Title("Save As").Filter("All Files", "*").Save()
	//fmt.Println(file)
	//fmt.Println("Error:", err)
	//dialog.Message("You chose file: %s", file).Title("Goodbye world!").Error()

	return true
}


func (me *Box) confirmDialog(t string, m string) bool {

	// Disabled because of GOOS=windows & GOOS=linux
	//ok := dialog.Message("%s", m).Title(t).YesNo()
	ok := true

	return ok
}


func (me *Box) openAdmin() error {

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


func (me *Box) openTerminal() error {

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


func (me *Box) openAbout() error {

	cmd := exec.Command("open", "https://gearbox.works/")
	err := cmd.Run()

	return err
}


func (me *Box) onExit() {
	// Cleaning stuff here.
}


func getClockTime(tz string) string {
	t := time.Now()
	utc, _ := time.LoadLocation(tz)

	return t.In(utc).Format("15:04:05")
}


func (me *Box) getIcon(s string) []byte {

	if s == "" {
		return nil
	}

	fp := me.baseDir.AddFileToPath(s).String()

	b, err := ioutil.ReadFile(fp)
	if err != nil {
		fmt.Print(err)
	}

	return b
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

