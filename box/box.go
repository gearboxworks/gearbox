package box

import (
	"fmt"
	"gearbox/app/logger"
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
<<<<<<< HEAD
	"github.com/gearboxworks/go-osbridge"

	//	"gearbox/os_support"
=======
	"gearbox/ssh"
	"gearbox/util"
	"github.com/gearboxworks/go-osbridge"
>>>>>>> master
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/getlantern/systray"
	"os"
	"os/signal"
	"syscall"
)

<<<<<<< HEAD

func New(OsBridge osbridge.OsBridger, args ...Args) (*Box, status.Status) {
=======
type Box struct {
	Boxname      string
	State        State
	VmBaseDir    string
	VmIsoDir     string
	VmIsoVersion string
	VmIsoFile    string
	VmIsoUrl     string
	VmIsoInfo    Release
	VmIsoDlIndex int

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

	OsBridge osbridge.OsBridger
}
type Args Box

func NewBox(OsBridge osbridge.OsBridger, args ...Args) *Box {
>>>>>>> master

	var _args Args
	var sts status.Status
	hb := &Box{}

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

<<<<<<< HEAD
		//foo := box.Args{}
		//err := copier.Copy(&foo, &_args)
		//if err != nil {
		//	sts = status.Wrap(err).SetMessage("unable to copy Box config")
		//	break
		//}
=======
		_args.OsBridge = OsBridge

		if _args.Boxname == "" {
			_args.Boxname = global.Brandname
		}

		if _args.WaitDelay == 0 {
			_args.WaitDelay = DefaultWaitDelay
		}
>>>>>>> master

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

<<<<<<< HEAD
	return hb, sts
}

=======
		if _args.VmBaseDir == "" {
			_args.VmBaseDir = string(OsBridge.GetUserConfigDir() + "/box/vm")
		}

		if _args.VmIsoDir == "" {
			_args.VmIsoDir = string(OsBridge.GetUserConfigDir() + "/box/iso")
		}
>>>>>>> master

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

<<<<<<< HEAD
		// Doesn't seem to work properly - need a workaround of some sort.
		//		if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		//			fmt.Printf("Sub-command not available for user.\n")
		//			//break
		//		}
=======
	return box
}

func (me *Box) Initialize() (sts status.Status) {
>>>>>>> master


<<<<<<< HEAD
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
=======
func (me *Box) WaitForVmState(displayString string) bool {
>>>>>>> master


		me.EventBroker, err = eventbroker.New(eventbroker.Args{Boxname: me.Boxname})
		if err != nil {
			break
		}

		_, _, err = me.EventBroker.AttachCallback(entity.UnfsdEntityName, myCallback, me)
		if err != nil {
			eblog.Debug(me.EntityId, "failed to attach callback")
			break
		}

<<<<<<< HEAD
		_, _, err = me.EventBroker.AttachCallback(entity.MqttBrokerEntityName, myCallback, me)
		if err != nil {
			eblog.Debug(me.EntityId, "failed to attach callback")
=======
		time.Sleep(me.WaitDelay)
		spinner.Update(fmt.Sprintf("%s [%d]", displayString, waitCount))
	}

	spinner.Stop(found)

	return found
}

func closeDialConnection(conn net.Conn) {
	_ = conn.Close()
}

func newSpinner(displayString string) *util.Spinner {
	return util.NewSpinner(util.SpinnerArgs{
		Text:    displayString,
		ExitOK:  displayString + " - OK",
		ExitNOK: displayString + " - FAILED",
	})
}

func (me *Box) heartbeatOk(b []byte, n int) (sts status.Status) {

	for range only.Once {
		apiSplit := strings.Split(string(b[:n]), ";")
		if len(apiSplit) <= 1 {
>>>>>>> master
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

<<<<<<< HEAD

		err = me.EventBroker.Start()
		if err != nil {
			sts = status.Wrap(err).SetMessage("EventBroker was not able to start")
			break
		}

=======
		sts = status.Success("received 'OK' from console")
		_ = sts.SetData(OkState)
	}

	return sts
}

func (me *Box) Start() (sts status.Status) {
>>>>>>> master

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

<<<<<<< HEAD
		//me.osRelease = me.VmBox.Releases.Selected
		//me.menu["version"].MenuItem.SetTitle(fmt.Sprintf("Gearbox (v%s)", me.osRelease.Version))
		//me.menu["version"].MenuItem.SetTooltip(fmt.Sprintf("Running v%s", me.osRelease.Version))


		// Setup systray menus.
		fmt.Printf("Gearbox: Starting systray.\n")
		systray.Run(me.onReady, me.onExit)


		//time.Sleep(time.Second * 10)
		//state, _ := me.EventBroker.GetSimpleStatus()
		//fmt.Printf("STATUS:\n%s", state.String())
=======
		/*
			switch {
				case me.State.VM.CurrentState == VmStateUnknown:
					sts = me.State.LastSts
					break

				case me.State.VM.CurrentState == VmStateStarting:
					sts = me.State.LastSts
					break

				case me.State.VM.CurrentState == VmStateRunning:
					sts = me.State.LastSts
					break

				case me.State.VM.CurrentState == VmStateStopping:
					sts = me.State.LastSts
					break

				case me.State.VM.CurrentState == VmStatePowerOff:
					// fall-through
			}
		*/
>>>>>>> master

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

<<<<<<< HEAD

func (me *Box) StartBox() (sts status.Status) {
=======
func (me *Box) Stop() (sts status.Status) {
>>>>>>> master

	for range only.Once {

		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

<<<<<<< HEAD
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
=======
		/*
			switch {
				case me.State.VM.CurrentState == VmStateUnknown:
					sts = me.State.LastSts
					break

				case me.State.VM.CurrentState == VmStateStarting:
					sts = me.State.LastSts
					break

				case me.State.VM.CurrentState == VmStateRunning:
					// fall-through

				case me.State.VM.CurrentState == VmStateStopping:
					sts = me.State.LastSts
					break

				case me.State.VM.CurrentState == VmStatePowerOff:
					sts = me.State.LastSts
					break
			}
		*/

		_, sts = me.StopBox()
		if is.Error(sts) {
			break
		}
>>>>>>> master

		//sts = me.DaemonInstance.Load()
		fmt.Printf("For now, we're running in the forground.\n")
		err = me.BoxDaemon()
		if err != nil {
			break
		}
		fmt.Printf("%s\n", sts.Message())

	}

<<<<<<< HEAD
	//eblog.LogIfNil(me, err)
	//eblog.LogIfError(me.EntityId, err)

	if err != nil {
		sts = status.Fail().
			SetMessage("Box terminated with error").
			SetData(err)
=======
	if !is.Error(sts) {
		sts = status.Success("%s VM stopped", global.Brandname)
>>>>>>> master
	}

	return sts
}

<<<<<<< HEAD

func (me *Box) StopBox() (sts status.Status) {

	var err error
=======
func (me *Box) Restart() (sts status.Status) {
>>>>>>> master

	for range only.Once {

		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

<<<<<<< HEAD
		//sts = me.DaemonInstance.Unload()
		if err != nil {
=======
		sts = me.Stop()
		if is.Error(sts) {
			break
		}
		if me.State.VM.CurrentState != me.State.VM.WantState {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("%s VM in an unknown state: %s", global.Brandname, me.State),
				Data:    VmStateUnknown,
			})
>>>>>>> master
			break
		}
		//fmt.Printf("%s\n", sts.Message())
		// fmt.Printf("%s Box - Started service.\n", global.Brandname)

<<<<<<< HEAD
=======
		sts = me.Start()
		if is.Error(sts) {
			break
		}
		if me.State.VM.CurrentState != me.State.VM.WantState {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("%s VM in an unknown state: %s", global.Brandname, me.State),
				Data:    VmStateUnknown,
			})
			break
		}
>>>>>>> master
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

<<<<<<< HEAD

func (me *Box) RestartBox() (sts status.Status) {
=======
func (me *Box) GetCachedState() (state State, sts status.Status) {

	// This is required so that not more than one process bashes VB at the same time.
	// This causes no end of issues.
>>>>>>> master

	for range only.Once {

<<<<<<< HEAD
		sts = me.StopBox()
=======
	return
}

func (me *Box) GetState() (State, status.Status) {

	// Possible VM states:
	// running
	// paused
	// saved
	// poweroff

	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
>>>>>>> master
		if is.Error(sts) {
			break
		}

		sts = me.StartBox()
		if is.Error(sts) {
			break
		}
<<<<<<< HEAD
=======
		if me.State.VM.CurrentState != VmStateRunning {
			me.State.API.CurrentState = VmStatePowerOff
			break
		}
		me.State.API.WantState = me.State.VM.WantState

		sts = me.GetApiStatus("", 10)

		if me.State.VM.WantState == VmStateInit {
			me.State.VM.WantState = me.State.VM.CurrentState
		}
		//me.State.GetStateMeaning()
		//me.State.LastSts = sts
		//fmt.Printf("FOO:", vmState)

		//		if err == nil {
		//			sts = status.Success("%s VM in a valid state: %s", global.Brandname, state)
		//			sts.SetData(state)
		//
		//		}
>>>>>>> master

	}

	return sts
}

<<<<<<< HEAD

func (me *Box) GetState() (sts status.Status) {
=======
func (me *Box) GetVmStatus() status.Status {
>>>>>>> master

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

<<<<<<< HEAD
	//eblog.LogIfNil(me, err)
	//eblog.LogIfError(me.EntityId, err)
=======
		/*
			// First check on the VM.
			// state, err := me.VmInstance.GetState()
			switch kvm["VMState"] {
				case VmStateRunning:
					me.State.VM.CurrentState = VmStateRunning

				case VmStatePowerOff:
					fallthrough
				case VmStateSaved:
					fallthrough
				case VmStatePaused:
					me.State.VM.CurrentState = VmStatePowerOff
					me.State.API.CurrentState = VmStatePowerOff
			}
		*/

		/*
			switch {
				case err != nil:
					// No Gearbox VM available - need to create one.
					me.State.VM.CurrentState = VmStateUnknown
					me.State.LastSts = status.Fail(&status.Args{
						Message: fmt.Sprintf("%s VM - needs to be created", global.Brandname),
						Help:    help.ContactSupportHelp(), // @TODO need better support here
						Data:    me.State.VM.CurrentState,
					})
					break

				case (me.State.VM.CurrentState == me.State.VM.WantState) && (state == VmStatePowerOff):
					// If we are not changing states and the VM is halted.
					me.State.VM.CurrentState = VmStatePowerOff
					me.State.LastSts = status.Success("%s VM - halted", global.Brandname)
					break

				case (me.State.VM.CurrentState == me.State.VM.WantState) && (state == VmStateRunning):
					// If we are not changing states and the VM is running.
					me.State.VM.CurrentState = VmStateRunning
					me.State.LastSts = status.Success("%s VM - running", global.Brandname)
					// Don't break here - need to check on the API.

				case (me.State.VM.CurrentState != me.State.VM.WantState) && (state == VmStatePowerOff):
					// If we are changing states then the VM is halting.
					me.State.VM.CurrentState = VmStateStopping
					me.State.LastSts = status.Success("%s VM - stopping", global.Brandname)
					// Don't break here - need to check on the API.

				case (me.State.VM.CurrentState != me.State.VM.WantState) && (state == VmStateRunning):
					// If we are changing states then the VM is starting.
					me.State.VM.CurrentState = VmStateStarting
					me.State.LastSts = status.Success("%s VM - starting", global.Brandname)
					// Don't break here - need to check on the API.
			}
		*/
>>>>>>> master

	if err != nil {
		sts = status.Fail().
			SetMessage("Box terminated with error").
			SetData(err)
	}

	return sts
}

<<<<<<< HEAD

func (me *Box) CreateBox() (sts status.Status) {
=======
// We have to have some way to block access to other concurrent processes/threads
// So, we're simply establishing a boolean that indicates this fact.
// var alreadyRunning = false
// @TODO - OK, so that's not working out.
func (me *Box) GetApiStatus(displayString string, waitFor time.Duration) (sts status.Status) {
>>>>>>> master

	for range only.Once {

	fmt.Printf("Not implemented.\n")

	}

	return sts
}

<<<<<<< HEAD
=======
func EnsureNotNil(bx *Box) (sts status.Status) {
	if bx == nil {
		sts = status.Fail(&status.Args{
			Message: "unexpected error",
			Help:    help.ContactSupportHelp(), // @TODO need better support here
			Data:    VmStateUnknown,
		})
	}

	return sts
}
>>>>>>> master
