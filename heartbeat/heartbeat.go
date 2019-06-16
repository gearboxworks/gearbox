package heartbeat

import (
	"fmt"
	"gearbox/app/logger"
	"gearbox/box"
	"gearbox/eventbroker"
	"gearbox/eventbroker/daemon"
	"gearbox/eventbroker/entity"
	"gearbox/heartbeat/external/vmbox"
	"gearbox/help"
	"gearbox/only"
	"gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/getlantern/systray"
	"github.com/jinzhu/copier"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)


func New(OsSupport oss.OsSupporter, args ...Args) (*Heartbeat, status.Status) {

	var _args Args
	var sts status.Status
	hb := &Heartbeat{}

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		_args.OsSupport = OsSupport
		foo := box.Args{}
		err := copier.Copy(&foo, &_args)
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("unable to copy Heartbeat config").
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		// Start a new VM Box instance.
		//_args.BoxInstance = box.NewBox(OsSupport, foo)
		//
		//if _args.Boxname == "" {
		//	_args.BoxInstance.Boxname = global.Brandname
		//} else {
		//	_args.BoxInstance.Boxname = _args.Boxname
		//}
		//_args.BoxInstance.Boxname = "harry" // DEBUG
		//
		//if _args.WaitDelay == 0 {
		//	_args.BoxInstance.WaitDelay = DefaultWaitDelay
		//} else {
		//	_args.BoxInstance.WaitDelay = _args.WaitDelay
		//}
		//
		//if _args.WaitRetries == 0 {
		//	_args.BoxInstance.WaitRetries = DefaultWaitRetries
		//} else {
		//	_args.BoxInstance.WaitRetries = _args.WaitRetries
		//}
		//
		//if _args.ConsoleHost == "" {
		//	_args.BoxInstance.ConsoleHost = DefaultConsoleHost
		//} else {
		//	_args.BoxInstance.ConsoleHost = _args.ConsoleHost
		//}
		//
		//if _args.ConsolePort == "" {
		//	_args.BoxInstance.ConsolePort = DefaultConsolePort
		//} else {
		//	_args.BoxInstance.ConsolePort = _args.ConsolePort
		//}
		//
		//if _args.ConsoleOkString == "" {
		//	_args.BoxInstance.ConsoleOkString = DefaultConsoleOkString
		//} else {
		//	_args.BoxInstance.ConsoleOkString = _args.ConsoleOkString
		//}
		//
		//if _args.ConsoleReadWait == 0 {
		//	_args.BoxInstance.ConsoleReadWait = DefaultConsoleReadWait
		//} else {
		//	_args.BoxInstance.ConsoleReadWait = _args.ConsoleReadWait
		//}
		//
		//if _args.SshUsername == "" {
		//	_args.BoxInstance.SshUsername = ssh.DefaultUsername
		//} else {
		//	_args.BoxInstance.SshUsername = _args.SshUsername
		//}
		//
		//if _args.SshPassword == "" {
		//	_args.BoxInstance.SshPassword = ssh.DefaultPassword
		//} else {
		//	_args.BoxInstance.SshPassword = _args.SshPassword
		//}
		//
		//if _args.SshPublicKey == "" {
		//	_args.BoxInstance.SshPublicKey = ssh.DefaultKeyFile
		//} else {
		//	_args.BoxInstance.SshPublicKey = _args.SshPublicKey
		//}

		_args.baseDir = filepath.FromSlash(fmt.Sprintf("%s/dist/heartbeat", _args.OsSupport.GetUserConfigDir()))
		_args.PidFile = filepath.FromSlash(fmt.Sprintf("%s/%s", _args.baseDir, DefaultPidFile))


		//execPath, _ := os.Executable()
		//execCwd := string(_args.OsSupport.GetAdminRootDir()) + "/heartbeat" // os.Getwd()
		//// Start a new Daemon instance.
		//_args.DaemonInstance = daemon.NewDaemon(_args.OsSupport, daemon.Args{
		//	Boxname: _args.Boxname,
		//	ServiceData: daemon.PlistData{
		//		Label:       "com.gearbox.heartbeat",
		//		Program:     execPath,
		//		ProgramArgs: []string{"heartbeat", "daemon"},
		//		Path:        execCwd,
		//		PidFile:     _args.PidFile,
		//		KeepAlive:   true,
		//		RunAtLoad:   true,
		//	},
		//})

		*hb = Heartbeat(_args)
	}

	return hb, sts
}


func (me *Heartbeat) HeartbeatDaemon() (sts status.Status) {

	var err error

	for range only.Once {

		sts = me.GetState()
		if is.Error(sts) {
			break
		}

		if daemon.IsParentInit() {
		//if !daemon.IsParentInit() {
			fmt.Printf("Gearbox: Sub-command not available for user.\n")
			sts = status.Fail().
				SetMessage("daemon mode cannot be run by user specifically").
				SetAdditional("use 'gearbox heartbeat start' to start daemon", ).
				SetData("").
				SetHelp(status.AllHelp, help.ContactSupportHelp())
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
			_ = me.VmBox.Stop()
			_ = me.EventBroker.Stop()
			logger.Debug("Goodbye!")

			os.Exit(0)
		}()


		//ebbd := filepath.FromSlash(fmt.Sprintf("%s/dist", _args.OsSupport.GetUserHomeDir()))
		//_args.EventBroker, err = eventbroker.New(eventbroker.Args{Boxname: _args.Boxname, SubBaseDir: "dist"})
		me.EventBroker, err = eventbroker.New(eventbroker.Args{Boxname: me.Boxname})
		if err != nil {
			break
		}

		_, _, err = me.EventBroker.AttachCallback(entity.UnfsdEntityName, myCallback, me)
		if err != nil {
			fmt.Printf("Ooops\n")
		}

		_, _, err = me.EventBroker.AttachCallback(entity.MqttBrokerEntityName, myCallback, me)
		if err != nil {
			fmt.Printf("Ooops\n")
		}

		_, _, err = me.EventBroker.AttachCallback(entity.VmEntityName, myCallback, me)
		if err != nil {
			fmt.Printf("Ooops\n")
		}

		_, _, err = me.EventBroker.AttachCallback(menuVmUpdate, myCallback, me)
		if err != nil {
			fmt.Printf("Ooops\n")
		}


		err = me.EventBroker.Start()
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("EventBroker was not able to start").
				SetAdditional("").
				SetData("").
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}


		fmt.Printf("Dropping in.\n")
		me.VmBox, err = vmbox.New(vmbox.Args{Channels: &me.EventBroker.Channels, OsPaths: me.EventBroker.OsPaths, Boxname: me.Boxname})
		if err != nil {
			break
		}

		err = me.VmBox.Start()
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("VM manager was not able to start").
				SetAdditional("").
				SetData("").
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		//me.osRelease = me.VmBox.Releases.Selected
		//me.menu["version"].MenuItem.SetTitle(fmt.Sprintf("Gearbox (v%s)", me.osRelease.Version))
		//me.menu["version"].MenuItem.SetTooltip(fmt.Sprintf("Running v%s", me.osRelease.Version))


		// Setup systray menus.
		fmt.Printf("Gearbox: Starting systray.\n")
		systray.Run(me.onReady, me.onExit)


		//time.Sleep(time.Second * 10)
		//state, _ := me.EventBroker.GetSimpleStatus()
		//fmt.Printf("STATUS:\n%s", state.String())

		//me.EventBroker.SimpleLoop()

		//fmt.Printf("Breaking out.\n")
		//time.Sleep(time.Second * 2)
		//_ = me.EventBroker.Stop()

		// Create a new VM Box instance.
		//fmt.Printf("Gearbox: Creating unfsd instance.\n")
		//me.NfsInstance, sts = unfsd.NewUnfsd(me.OsSupport)

		// Should never exit, unless we get a signal to do so.
	}
	status.Log(sts)

	return sts
}


func (me *Heartbeat) StartHeartbeat() (sts status.Status) {

	for range only.Once {

		sts = me.EnsureNotNil()
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
//
//		if me.DaemonInstance.IsLoaded() {
//			fmt.Printf("%s Heartbeat - Restarting service.\n", global.Brandname)
//			sts = me.DaemonInstance.Unload()
//			if is.Error(sts) {
//				break
//			}
//		}

		//sts = me.DaemonInstance.Load()
		if is.Error(sts) {
			break
		}
		fmt.Printf("%s\n", sts.Message())

	}

	return sts
}


func (me *Heartbeat) StopHeartbeat() (sts status.Status) {

	for range only.Once {

		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		//sts = me.DaemonInstance.Unload()
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

		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		if me == nil {
			sts = status.Fail().
				SetMessage("unexpected software error").
				SetAdditional("", ).
				SetData("").
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		//sts = me.DaemonInstance.GetState()
		if is.Error(sts) {
			break
		}
		if sts != nil {
			fmt.Printf(sts.Message())
		}
	}

	return sts
}

