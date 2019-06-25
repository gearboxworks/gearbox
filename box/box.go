package box

import (
	"fmt"
	"gearbox/box/external/unfsd"
	"gearbox/eventbroker"
	"gearbox/eventbroker/daemon"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/ospaths"
	"gearbox/eventbroker/states"
	"gearbox/global"
	"gearbox/box/external/vmbox"
	"gearbox/only"
	"github.com/gearboxworks/go-osbridge"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/getlantern/systray"
	"os"
	"os/signal"
	"syscall"
)


func New(OsBridge osbridge.OsBridger, args ...Args) (*Box, status.Status) {

	var _args Args
	var sts status.Status
	hb := &Box{}

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		//foo := box.Args{}
		//err := copier.Copy(&foo, &_args)
		//if err != nil {
		//	sts = status.Wrap(err).SetMessage("unable to copy Box config")
		//	break
		//}

		if _args.EntityId == "" {
			_args.EntityId = *messages.GenerateAddress()
		}

		if _args.EntityName == "" {
			_args.EntityName = DefaultEntityName
		}

		if _args.Boxname == "" {
			_args.Boxname = global.Brandname
		}

		if _args.Version == "" {
			_args.Version = "latest"
		}

		_args.osBridge = OsBridge
		_args.osPaths = ospaths.New("")
		_args.baseDir = _args.osPaths.UserConfigDir.AddToPath(DefaultBaseDir)
		_args.pidFile = _args.baseDir.AddFileToPath(DefaultPidFile).String()

		_args.State = states.New(&_args.EntityId, &_args.EntityName, entity.SelfEntityName)


		*hb = Box(_args)
	}

	return hb, sts
}


func (me *Box) BoxDaemon() (sts status.Status) {

	var err error

	for range only.Once {
		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		if daemon.IsParentInit() {
		//if !daemon.IsParentInit() {
			fmt.Printf("Gearbox: Sub-command not available for user.\n")
			sts = status.Fail().SetMessage("daemon mode cannot be run by user specifically")
			break
		}
		fmt.Printf("Gearbox: Starting Box daemon.\n")

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
			fmt.Printf("Gearbox exiting.\n")

			os.Exit(0)
		}()


		me.EventBroker, err = eventbroker.New(eventbroker.Args{Boxname: me.Boxname})
		if err != nil {
			break
		}

		_, _, err = me.EventBroker.AttachCallback(entity.UnfsdEntityName, myCallback, me)
		if err != nil {
			eblog.Debug(me.EntityId, "failed to attach callback")
			break
		}

		_, _, err = me.EventBroker.AttachCallback(entity.MqttBrokerEntityName, myCallback, me)
		if err != nil {
			eblog.Debug(me.EntityId, "failed to attach callback")
			break
		}

		_, _, err = me.EventBroker.AttachCallback(entity.VmEntityName, myCallback, me)
		if err != nil {
			eblog.Debug(me.EntityId, "failed to attach callback")
			break
		}

		_, _, err = me.EventBroker.AttachCallback(entity.ApiEntityName, myCallback, me)
		if err != nil {
			eblog.Debug(me.EntityId, "failed to attach callback")
			break
		}

		_, _, err = me.EventBroker.AttachCallback(menuVmUpdate, myCallback, me)
		if err != nil {
			eblog.Debug(me.EntityId, "failed to attach callback")
			break
		}


		err = me.EventBroker.Start()
		if err != nil {
			sts = status.Wrap(err).SetMessage("EventBroker was not able to start")
			break
		}

		me.NfsExports, err = unfsd.New(unfsd.Args{Channels: &me.EventBroker.Channels, OsPaths: me.EventBroker.OsPaths, Boxname: me.Boxname})
		if err != nil {
			break
		}

		fmt.Printf("Dropping in.\n")
		me.VmBox, err = vmbox.New(vmbox.Args{Channels: &me.EventBroker.Channels, OsPaths: me.EventBroker.OsPaths, Boxname: me.Boxname})
		if err != nil {
			break
		}

		err = me.VmBox.Start()
		if err != nil {
			sts = status.Wrap(err).SetMessage("VM manager was not able to start")
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
		//me.NfsInstance, sts = unfsd.NewUnfsd(me.OsBridge)

		// Should never exit, unless we get a signal to do so.
	}

	//eblog.LogIfNil(me, err)
	//eblog.LogIfError(me.EntityId, err)

	if err != nil {
		sts = status.Fail().
			SetMessage("Box terminated with error").
			SetData(err)
	}

	return sts
}


func (me *Box) StartBox() (sts status.Status) {

	var err error

	for range only.Once {

		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

//		if me.DaemonInstance.IsRunning() {
//			fmt.Printf("%s Box - Restarting service.\n", global.Brandname)
//			sts = me.DaemonInstance.Unload()
//			if is.Error(sts) {
//				break
//			}
//		}
//
//		if me.DaemonInstance.IsLoaded() {
//			fmt.Printf("%s Box - Restarting service.\n", global.Brandname)
//			sts = me.DaemonInstance.Unload()
//			if is.Error(sts) {
//				break
//			}
//		}

		//sts = me.DaemonInstance.Load()
		fmt.Printf("Gearbox: The alpha release runs Gearbox in the foreground. Please keep this shell open.\n\n")
		err = me.BoxDaemon()
		if err != nil {
			fmt.Printf("Gearbox: Error %v\n", err)
			break
		}

	}

	//eblog.LogIfNil(me, err)
	//eblog.LogIfError(me.EntityId, err)

	if err != nil {
		sts = status.Fail().
			SetMessage("Box terminated with error").
			SetData(err)
	}

	return sts
}


func (me *Box) StopBox() (sts status.Status) {

	var err error

	for range only.Once {

		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		//sts = me.DaemonInstance.Unload()
		if err != nil {
			break
		}
		//fmt.Printf("%s\n", sts.Message())
		// fmt.Printf("%s Box - Started service.\n", global.Brandname)

	}

	//eblog.LogIfNil(me, err)
	//eblog.LogIfError(me.EntityId, err)

	if err != nil {
		sts = status.Fail().
			SetMessage("Box terminated with error").
			SetData(err)
	}

	return sts
}


func (me *Box) RestartBox() (sts status.Status) {

	for range only.Once {

		sts = me.StopBox()
		if is.Error(sts) {
			break
		}

		sts = me.StartBox()
		if is.Error(sts) {
			break
		}
	}

	return sts
}


func (me *Box) GetState() (sts status.Status) {

	var err error

	for range only.Once {

		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		if me == nil {
			sts = status.Fail().
				SetMessage("Box terminated with error").
				SetData(err)
			break
		}

		//sts = me.DaemonInstance.GetState()
		if err != nil {
			break
		}
	}

	//eblog.LogIfNil(me, err)
	//eblog.LogIfError(me.EntityId, err)

	if err != nil {
		sts = status.Fail().
			SetMessage("Box terminated with error").
			SetData(err)
	}

	return sts
}


func (me *Box) CreateBox() (sts status.Status) {

	for range only.Once {

	fmt.Printf("Not implemented.\n")

	}

	return sts
}

