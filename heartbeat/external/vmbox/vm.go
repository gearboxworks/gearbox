package vmbox

import (
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/states"
	"gearbox/only"
	"gearbox/os_support"

	"time"
)


type Box struct {
	Boxname         string
	//State           BoxState
	VmBaseDir       string
	VmIsoDir        string
	VmIsoVersion    string
	VmIsoFile       string
	VmIsoUrl 		string
	VmIsoInfo	    Release
	VmIsoDlIndex	int

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


func (me *VmBox) New(c ServiceConfig) (*Vm, error) {

	var err error

	sc := Vm{}

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		//if _args.WaitDelay == 0 {
		//	_args.WaitDelay = DefaultWaitDelay
		//}
		//
		//if _args.WaitRetries == 0 {
		//	_args.WaitRetries = DefaultWaitRetries
		//}
		//
		//if _args.ConsoleHost == "" {
		//	_args.ConsoleHost = DefaultConsoleHost
		//}
		//
		//if _args.ConsolePort == "" {
		//	_args.ConsolePort = DefaultConsolePort
		//}
		//
		//if _args.ConsoleOkString == "" {
		//	_args.ConsoleOkString = DefaultConsoleOkString
		//}
		//
		//if _args.ConsoleReadWait == 0 {
		//	_args.ConsoleReadWait = DefaultConsoleReadWait
		//}
		//
		//if _args.SshUsername == "" {
		//	_args.SshUsername = ssh.DefaultUsername
		//}
		//
		//if _args.SshPassword == "" {
		//	_args.SshPassword = ssh.DefaultPassword
		//}
		//
		//if _args.SshPublicKey == "" {
		//	_args.SshPublicKey = ssh.DefaultKeyFile
		//}
		//
		//if _args.VmBaseDir == "" {
		//	_args.VmBaseDir = string(OsSupport.GetUserConfigDir() + "/box/vm")
		//}
		//
		//if _args.VmIsoDir == "" {
		//	_args.VmIsoDir = string(OsSupport.GetUserConfigDir() + "/box/iso")
		//}
		//
		//_args.VmIsoDlIndex = 100


		if c.Name == "" {
			err = me.EntityName.ProduceError("VM doesn't have a name")
			break
		}

		if c.Version == "" {
			c.Version = "latest"
		}

		if c.ConsolePort == "" {
			c.ConsolePort = "2232"
		}

		if c.SshPort == "" {
			c.SshPort = "2222"
		}


		var rel *Release
		rel, err = me.Releases.SelectRelease(ReleaseSelector{SpecificVersion: "latest"})
		if err != nil {
			break
		}

		err = me.Releases.ShowReleases()
		if err != nil {
			break
		}


		sc.EntityId = *messages.GenerateAddress()
		sc.EntityName = c.Name
		sc.EntityParent = &me.EntityId
		sc.State = states.New(&sc.EntityId, &sc.EntityName, me.EntityId)
		sc.State.SetNewAction(states.ActionRegister)
		sc.IsManaged = true
		sc.Entry = &ServiceConfig{
			Name: sc.EntityName,
			baseDir: me.OsPaths.UserConfigDir.AddToPath("vm"),
			osRelease: rel,
			ConsolePort: c.ConsolePort,
			SshPort: c.SshPort,
			retryMax: DefaultRetries,
			retryDelay: DefaultVmWaitTime,
		}
		sc.channels = me.Channels
		sc.channels.PublishState(sc.State)


		var state states.Status
		state, err = sc.Status()
		if err != nil {
			break
		}

		switch state.Current {
			case states.StateError:
				eblog.Debug(me.EntityId, "%v", err)


			case states.StateStopped:
				err = me.AddEntity(sc.EntityId, &sc)
				if err != nil {
					break
				}
				sc.State.SetNewAction(states.ActionStop)
				eblog.Debug(me.EntityId, "VM registered OK")

			case states.StateStarted:
				// VM already created but started.
				sc.State.SetNewAction(states.ActionStart)

			case states.StateUnregistered:
				eblog.Debug(me.EntityId, "VM not created")


			case states.StateUnknown:
				eblog.Debug(me.EntityId, "%v", err)
		}


		sc.State.SetNewState(state.Current, err)
		sc.channels.PublishState(me.State)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return &sc, err
}


func (me *Vm) Start() error {

	var err error
	var ok bool

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.State.SetNewAction(states.ActionStart)
		me.channels.PublishState(me.State)

		err = me.vbStart()
		if err != nil {
			break
		}

		ok, err = me.vbWaitForState(states.StateStarted)
		if err != nil {
			break
		}
		if !ok {
			break
		}

		me.State.SetNewState(states.StateStarted, err)
		me.channels.PublishState(me.State)
		eblog.Debug(me.EntityId, "started VM")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Vm) Stop() error {

	var err error
	var ok bool

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.State.SetNewAction(states.ActionStop)
		me.channels.PublishState(me.State)

		err = me.vbStart()
		if err != nil {
			break
		}

		ok, err = me.vbWaitForState(states.StateStopped)
		if err != nil {
			break
		}
		if !ok {
			break
		}

		me.State.SetNewState(states.StateStopped, err)
		me.channels.PublishState(me.State)
		eblog.Debug(me.EntityId, "stopped VM")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Vm) Restart() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = me.Stop()
		if err != nil {
			break
		}

		err = me.Start()
		if err != nil {
			break
		}
	}

	return err
}


func (me *Vm) Status() (states.Status, error) {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		var state states.State
		state, err = me.vbStatus()
		if err != nil {
			break
		}

		me.State.SetNewState(state, err)
		me.channels.PublishState(me.State)
		eblog.Debug(me.EntityId, "started VM")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return *me.State.GetStatus(), err
}


//func (me *Vm) WaitForVmState(displayString string) bool {
//
//	found := false
//	var waitCount int
//
//	spinner := util.NewSpinner(util.SpinnerArgs{
//		Text:    displayString,
//		ExitOK:  displayString + " - OK",
//		ExitNOK: displayString + " - FAILED",
//	})
//	spinner.Start()
//
//	for waitCount = 0; waitCount < me.WaitRetries; waitCount++ {
//
//		_, sts := me.GetState()
//		if is.Error(sts) {
//			found = false
//			break
//		}
//
//		if me.State.VM.CurrentState == me.State.VM.WantState {
//			found = true
//			break
//		}
//
//		time.Sleep(me.WaitDelay)
//		spinner.Update(fmt.Sprintf("%s [%d]", displayString, waitCount))
//	}
//
//	spinner.Stop(found)
//
//	return found
//}


//func closeDialConnection(conn net.Conn) {
//	_ = conn.Close()
//}
//
//
//func newSpinner(displayString string) *util.Spinner {
//	return util.NewSpinner(util.SpinnerArgs{
//		Text:    displayString,
//		ExitOK:  displayString + " - OK",
//		ExitNOK: displayString + " - FAILED",
//	})
//}
//
//
//func (me *Vm) heartbeatOk(b []byte, n int) (sts status.Status) {
//
//	var err error
//
//	for range only.Once {
//		err = me.EnsureNotNil()
//		if err != nil {
//			break
//		}
//
//		apiSplit := strings.Split(string(b[:n]), ";")
//		if len(apiSplit) <= 1 {
//			break
//		}
//
//		match, _ := regexp.MatchString(me.ConsoleOkString, apiSplit[1])
//		if !match {
//			break
//		}
//
//		fmt.Printf("API:%v\n", apiSplit)
//		if len(apiSplit) < 2 {
//			sts = status.Fail(&status.Args{
//				Message: fmt.Sprintf("did not receive 'OK' from console: %s",
//					apiSplit[2],
//				),
//				Data: NotOkState,
//			})
//			break
//		}
//
//		if apiSplit[2] != "OK" {
//			sts = status.Fail(&status.Args{
//				Message: fmt.Sprintf("did not receive 'OK' from console: %s",
//					apiSplit[2],
//				),
//				Data: NotOkState,
//			})
//			break
//		}
//
//		sts = status.Success("received 'OK' from console")
//		sts.SetData(OkState)
//	}
//
//	return sts
//}
//
//
//func (me *Vm) GetCachedState() (state BoxState, sts status.Status) {
//
//	var err error
//
//	// This is required so that not more than one process bashes VB at the same time.
//	// This causes no end of issues.
//
//	for range only.Once {
//		err = me.EnsureNotNil()
//		if err != nil {
//			break
//		}
//
//		state = me.State
//	}
//
//	return
//}
//
//
//func (me *Vm) GetState() (BoxState, status.Status) {
//
//	var err error
//
//	// Possible VM states:
//	// running
//	// paused
//	// saved
//	// poweroff
//
//	var sts status.Status
//
//	for range only.Once {
//		err = me.EnsureNotNil()
//		if err != nil {
//			break
//		}
//
//		sts = me.GetVmStatus()
//		if is.Error(sts) {
//			break
//		}
//		if me.State.VM.CurrentState != VmStateRunning {
//			me.State.API.CurrentState = VmStatePowerOff
//			break
//		}
//		me.State.API.WantState = me.State.VM.WantState
//
//		sts = me.GetApiStatus("", 10)
//
//		if me.State.VM.WantState == VmStateInit {
//			me.State.VM.WantState = me.State.VM.CurrentState
//		}
//		//me.State.GetStateMeaning()
//		//me.State.LastSts = sts
//		//fmt.Printf("FOO:", vmState)
//
//		//		if err == nil {
////			sts = status.Success("%s VM in a valid state: %s", global.Brandname, state)
////			sts.SetData(state)
////
////		}
//
//		//sts = me.State.LastSts
//		//fmt.Printf("STATE2: %v\n", sts)
//	}
//
//	return me.State, sts
//}
//
//
//func (me *Vm) GetVmStatus() error {
//
//	var err error
//
//	// Possible VM states:
//	// running
//	// paused
//	// saved
//	// poweroff
//
//	var sts status.Status
//	// var kvm KeyValueMap
//
//	for range only.Once {
//		err = me.EnsureNotNil()
//		if err != nil {
//			break
//		}
//
//		_, sts = me.cmdListVm()
//		// fmt.Printf("Box '%s' is in state: '%s'\n", kvm["name"], kvm["VMState"])
//		if is.Error(sts) {
//			me.State.VM.CurrentState = VmStateNotPresent
//		}
//
//		if me.State.VM.WantState == VmStateInit {
//			me.State.VM.WantState = me.State.VM.CurrentState
//		}
//
///*
//		// First check on the VM.
//		// state, err := me.VmInstance.GetState()
//		switch kvm["VMState"] {
//			case VmStateRunning:
//				me.State.VM.CurrentState = VmStateRunning
//
//			case VmStatePowerOff:
//				fallthrough
//			case VmStateSaved:
//				fallthrough
//			case VmStatePaused:
//				me.State.VM.CurrentState = VmStatePowerOff
//				me.State.API.CurrentState = VmStatePowerOff
//		}
//*/
//
///*
//		switch {
//			case err != nil:
//				// No Gearbox VM available - need to create one.
//				me.State.VM.CurrentState = VmStateUnknown
//				me.State.LastSts = status.Fail(&status.Args{
//					Message: fmt.Sprintf("%s VM - needs to be created", global.Brandname),
//					Help:    help.ContactSupportHelp(), // @TODO need better support here
//					Data:    me.State.VM.CurrentState,
//				})
//				break
//
//			case (me.State.VM.CurrentState == me.State.VM.WantState) && (state == VmStatePowerOff):
//				// If we are not changing states and the VM is halted.
//				me.State.VM.CurrentState = VmStatePowerOff
//				me.State.LastSts = status.Success("%s VM - halted", global.Brandname)
//				break
//
//			case (me.State.VM.CurrentState == me.State.VM.WantState) && (state == VmStateRunning):
//				// If we are not changing states and the VM is running.
//				me.State.VM.CurrentState = VmStateRunning
//				me.State.LastSts = status.Success("%s VM - running", global.Brandname)
//				// Don't break here - need to check on the API.
//
//			case (me.State.VM.CurrentState != me.State.VM.WantState) && (state == VmStatePowerOff):
//				// If we are changing states then the VM is halting.
//				me.State.VM.CurrentState = VmStateStopping
//				me.State.LastSts = status.Success("%s VM - stopping", global.Brandname)
//				// Don't break here - need to check on the API.
//
//			case (me.State.VM.CurrentState != me.State.VM.WantState) && (state == VmStateRunning):
//				// If we are changing states then the VM is starting.
//				me.State.VM.CurrentState = VmStateStarting
//				me.State.LastSts = status.Success("%s VM - starting", global.Brandname)
//				// Don't break here - need to check on the API.
//		}
//*/
//
//		// fmt.Printf("vmState: %v\n", sts)
//	}
//
//	return err
//}
//
//
//// We have to have some way to block access to other concurrent processes/threads
//// So, we're simply establishing a boolean that indicates this fact.
//// var alreadyRunning = false
//// @TODO - OK, so that's not working out.
//func (me *Vm) GetApiStatus(displayString string, waitFor time.Duration) (sts status.Status) {
//
//	for range only.Once {
//		sts = EnsureNotNil(me)
//		if is.Error(sts) {
//			break
//		}
//
//		spinner := newSpinner(displayString)
//		displaySpinner := !me.ShowConsole && displayString != ""
//
//		if displaySpinner {
//			// We want to display just a spinner instead of console output.
//			spinner.Start()
//		}
//
//		// Connect to this console
//		conn, err := net.Dial("tcp", me.ConsoleHost+":"+me.ConsolePort)
//		if err != nil {
//			me.State.API.CurrentState = VmStatePowerOff
//			sts = status.Fail(&status.Args{
//				Message: fmt.Sprintf("%s API - timeout", global.Brandname),
//				Help:    help.ContactSupportHelp(), // @TODO need better support here
//				Data:    me.State.API.CurrentState,
//			})
//			break
//		}
//		// defer closeDialConnection(conn)
//		defer conn.Close()
//
//		// Set default state before we begin.
//		me.State.API.CurrentState = VmStateUnknown
//		sts = status.Fail(&status.Args{
//			Message: fmt.Sprintf("%s API - no data", global.Brandname),
//			Help:    help.ContactSupportHelp(), // @TODO need better support here
//			Data:    me.State.API.CurrentState,
//		})
//
//		exitWhen := time.Now().Add(time.Second * waitFor)
//		readBuffer := make([]byte, 512)
//		for waitCount := 0; time.Now().Unix() < exitWhen.Unix(); waitCount++ {
//			err = conn.SetDeadline(time.Now().Add(me.ConsoleReadWait))
//			if err != nil {
//				me.State.API.CurrentState = VmStateUnknown
//				sts = status.Fail(&status.Args{
//					Message: fmt.Sprintf("%s API - deadline", global.Brandname),
//					Help:    help.ContactSupportHelp(), // @TODO need better support here
//					Data:    me.State.API.CurrentState,
//				})
//				break
//			}
//
//			bytesRead, err := bufio.NewReader(conn).Read(readBuffer)
//			// bytesRead, err := conn.Read(readBuffer)
//			// readBuffer, err := bufio.NewReader(conn).ReadString('\n')
//			// bytesRead := len(readBuffer)
//			if err != nil {
//				me.State.API.CurrentState = VmStateUnknown
//				sts = status.Fail(&status.Args{
//					Message: fmt.Sprintf("%s API - no data", global.Brandname),
//					Help:    help.ContactSupportHelp(), // @TODO need better support here
//					Data:    me.State.API.CurrentState,
//				})
//				break
//			}
//
//			if bytesRead > 0 {
//				if me.ShowConsole {
//					fmt.Printf("%s", string(readBuffer[:bytesRead]))
//				}
//
//				sts = me.heartbeatOk(readBuffer, bytesRead)
//				if sts != nil {
//					me.State.API.CurrentState = VmStateRunning
//					sts = status.Success("%s API - running", global.Brandname)
//					break
//
//				} else {
//					if me.State.API.WantState == VmStatePowerOff {
//						me.State.API.CurrentState = VmStateStopping
//						sts = status.Success("%s API - stopping", global.Brandname)
//					} else if me.State.API.WantState == VmStateRunning {
//						me.State.API.CurrentState = VmStateStarting
//						sts = status.Success("%s API - starting", global.Brandname)
//					}
//					// Do not break.
//				}
//			}
//
//			time.Sleep(me.WaitDelay)
//			if displaySpinner {
//				spinner.Update(fmt.Sprintf("%s [%d]", displayString, waitCount))
//			}
//		}
//
//		if me.ShowConsole {
//			fmt.Printf("\n\n# Exiting Console.\n")
//		}
//
//		if displaySpinner {
//			spinner.Stop(false)
//		}
//	}
//
//	// fmt.Printf("apiState: %v\n", sts)
//
//	return sts
//}
//
//
//func oldEnsureNotNil(bx *Vm) (sts status.Status) {
//	if bx == nil {
//		sts = status.Fail(&status.Args{
//			Message: "unexpected error",
//			Help:    help.ContactSupportHelp(), // @TODO need better support here
//			Data:    VmStateUnknown,
//		})
//	}
//
//	return sts
//}
