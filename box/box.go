package box

import (
	"fmt"
	"gearbox/box/external/unfsd"
	"gearbox/box/external/vmbox"
	"gearbox/eventbroker"
	"gearbox/eventbroker/daemon"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/osdirs"
	"gearbox/eventbroker/states"
	"gearbox/global"
	"github.com/gearboxworks/go-osbridge"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/gearboxworks/go-status/only"
//	"github.com/gearboxworks/go-systray"
	"github.com/getlantern/systray"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Box struct {
	EntityId    msgs.Address
	EntityName  msgs.Address
	Boxname     string
	Version     Version
	NfsExports  *unfsd.Unfsd
	State       *states.Status
	menu        Menus
	EventBroker *eventbroker.EventBroker
	VmBox       *vmbox.VmBox

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

	baseDir  osdirs.Dir
	pidFile  string
	osBridge osbridge.OsBridger
	osPaths  *osdirs.BaseDirs
}
type Args Box

func (me *Args) SetOsBridge(b osbridge.OsBridger) {
	me.osBridge = b
}

func New(args ...*Args) (b *Box, sts status.Status) {

	var _args Args
	b = &Box{}

	for range only.Once {

		if len(args) > 0 {
			_args = *args[0]
		}

		if _args.EntityId == "" {
			_args.EntityId = msgs.MakeAddress()
		}

		if _args.EntityName == "" {
			_args.EntityName = DefaultEntityName
		}

		if _args.Boxname == "" {
			_args.Boxname = global.Brandname
		}

		if _args.Version == "" {
			_args.Version = LatestVersion
		}

		_args.osPaths = osdirs.New()
		_args.baseDir = _args.osPaths.AppendToUserConfigDir(DefaultBaseDir)
		_args.pidFile = osdirs.AddFilef(_args.baseDir, DefaultPidFile)

		_args.State = states.New(
			_args.EntityId,
			_args.EntityName,
			entity.SelfEntityName,
		)

		*b = Box(_args)
	}

	return b, sts
}

func (me *Box) RunAsDaemon() (sts status.Status) {

	var err error

	for range only.Once {
		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		if daemon.IsParentInit() {
			//if !daemon.IsParentInit() {
			fmt.Printf("Gearbox: Sub-command not available for user.\n")
			sts = status.Fail().
				SetMessage("daemon mode cannot be run directly by user.").
				SetAllHelp("@TODO — Explain how daemon can be run")
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

		me.NfsExports, err = unfsd.New(unfsd.Args{
			Channels: &me.EventBroker.Channels,
			BaseDirs: me.EventBroker.OsDirs,
			Boxname:  me.Boxname,
		})
		if err != nil {
			break
		}

		me.VmBox, err = vmbox.New(vmbox.Args{
			Channels: &me.EventBroker.Channels,
			OsPaths:  me.EventBroker.OsDirs,
			Boxname:  me.Boxname,
		})
		if err != nil {
			break
		}

		err = me.VmBox.Start()
		if err != nil {
			sts = status.Wrap(err).SetMessage("VM manager was not able to start")
			break
		}

		// Setup systray menus.
		fmt.Printf("Gearbox: Starting systray.\n")
		systray.Run(me.onReady, me.onExit)
	}

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

		fmt.Printf("Gearbox: The alpha release runs Gearbox in the foreground. Please keep this shell open.\n\n")
		err = me.RunAsDaemon()
		if err != nil {
			fmt.Printf("Gearbox: Error %v\n", err)
			break
		}

	}

	if err != nil {
		sts = status.Fail().
			SetMessage("Box terminated with error").
			SetData(err)
	}

	return sts
}

func (me *Box) StopBox() (sts status.Status) {

	sts = me.EnsureNotNil()
	if is.Error(sts) {
		sts.Log()
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

	}

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

func (me *Box) DestroyBox() (sts status.Status) {

	for range only.Once {

		var err error

		me.VmBox, err = vmbox.New(vmbox.Args{
			Channels: &me.EventBroker.Channels,
			OsPaths:  me.EventBroker.OsDirs,
			Boxname:  me.Boxname,
		})
		if err != nil {
			break
		}

		err = me.VmBox.Stop()
		if err != nil {
			sts = status.Wrap(err).SetMessage("VM manager was not able to start")
			break
		}

	}

	return sts
}
