package vmbox

import (
	"bufio"
	"fmt"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/states"
	"github.com/gearboxworks/go-status/only"
	"net"
	"regexp"
	"strings"
	"time"
)



func (me *VmBox) New(c ServiceConfig) (*Vm, error) {

	var err error

	sc := Vm{}

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}


		sc.Entry = &c


		err = me.OsPaths.EnsureNotNil()
		if err != nil {
			break
		}
		sc.osPaths = me.OsPaths
		sc.channels = me.Channels


		// Check config, set defaults if needed.
		err = sc.VerifyConfig()
		if err != nil {
			break
		}


		// Load up config file, if it exists. Supersceding previous config.
		err = sc.ConfigExists()
		if err == nil {
			err = sc.ReadConfig()
			if err != nil {
				break
			}
		}


		// Fetch ISO releases.
		err = me.Releases.UpdateReleases()
		if err != nil {
			break
		}

		var rel *Release
		rel, err = me.Releases.SelectRelease(ReleaseSelector{SpecificVersion: sc.Entry.Version})
		if err != nil {
			break
		}

		// Fetch ISO.
		if !me.Releases.Selected.IsDownloading {
			err = me.Releases.Selected.GetIso()
		}

		// Update specific version fetched.
		sc.Entry.Version = string(rel.Version)


		sc.EntityId = *messages.GenerateAddress()
		sc.EntityName = sc.Entry.Name
		sc.EntityParent = &me.EntityId
		sc.State = states.New(&sc.EntityId, &sc.EntityName, me.EntityId)
		sc.State.SetNewAction(states.ActionStop)
		a := messages.MessageAddress(entity.ApiEntityName)
		sc.ApiState = states.New(&sc.EntityId, &a, me.EntityId)
		sc.ApiState.SetNewAction(states.ActionStop)
		sc.IsManaged = true
		sc.osRelease = rel


		sc.channels.PublishState(sc.State)


		var state states.State
		state, err = sc.vbCreate()
		if err != nil {
			fmt.Printf("Gearbox: Error creating VM %s with '%v'\n", sc.EntityName, err)
			break
		}

		err = sc.WriteConfig()
		if err != nil {
			break
		}

		switch state {
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
				err = me.AddEntity(sc.EntityId, &sc)
				if err != nil {
					break
				}
				sc.State.SetNewAction(states.ActionStart)

			case states.StateUnregistered:
				eblog.Debug(me.EntityId, "VM not created")


			case states.StateUnknown:
				eblog.Debug(me.EntityId, "%v", err)
		}


		sc.State.SetNewState(state, err)
		sc.channels.PublishState(sc.State)
		eblog.Debug(me.EntityId, "registered VM OK")
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


		// Check for ISO image first.
		var i int
		i, err = me.osRelease.IsIsoFilePresent()
		if i != IsoFileDownloaded {
			break
		}


		me.ChangeRequested = true
		defer func() { me.ChangeRequested = false }()


		me.State.SetNewAction(states.ActionStart)
		me.channels.PublishState(me.State)

		ok, err = me.vbStart(true)
		if !ok {
			fmt.Printf("Gearbox: failed to start VM\n")
			break
		}

		// Publish new state.
		me.State.SetNewState(states.StateStarted, err)
		me.channels.PublishState(me.State)


		// Now wait for API.
		me.ApiState.SetNewAction(states.ActionStart)
		me.channels.PublishState(me.ApiState)

		var state states.State
		state, err = me.waitForApiState(DefaultBootWaitTime, true)
		me.ApiState.SetNewState(state, err)
		me.channels.PublishState(me.ApiState)

		eblog.Debug(me.EntityId, "started VM OK")
		fmt.Printf("Gearbox: started VM OK\n")
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

		me.ChangeRequested = true
		defer func() { me.ChangeRequested = false }()

		me.State.SetNewAction(states.ActionStop)
		me.channels.PublishState(me.State)
		me.ApiState.SetNewAction(states.ActionStop)
		me.channels.PublishState(me.ApiState)

		ok, err = me.vbStop(false, true)
		if !ok {
			ok, err = me.vbStop(true, true)
			break
		}


		me.State.SetNewState(states.StateStopped, err)
		me.channels.PublishState(me.State)

		me.ApiState.SetNewState(states.StateStopped, err)
		me.channels.PublishState(me.ApiState)

		eblog.Debug(me.EntityId, "stopped VM OK")
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


func (me *Vm) UpdateRealStatus() error {

	var err error
	//var vm states.Status
	//var api states.Status

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}


		var v states.State
		v, err = me.vbStatus()
		me.State.SetNewState(v, err)
		me.channels.PublishState(me.State)


		var a states.State
		a, err = me.waitForApiState(DefaultRunWaitTime, false)
		me.ApiState.SetNewState(a, err)
		me.channels.PublishState(me.ApiState)


		eblog.Debug(me.EntityId, "VM is in state %s, API is in state %s", v, a)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Vm) waitForApiState(waitFor time.Duration, showConsole bool) (states.State, error) {

	var err error
	var state states.State

	state = states.StateIdle

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}


		if me.Entry.Console.Host == "" {
			err = me.EntityName.ProduceError("no VM console host defined")
			state = states.StateError
			break
		}

		if me.Entry.Console.Port == "" {
			err = me.EntityName.ProduceError("no VM console port defined")
			state = states.StateError
			break
		}


		// Ensure we only have one at a time.
		me.Entry.Console.mutex.Lock()
		defer me.Entry.Console.mutex.Unlock()


		// Connect to this console
		conn, err := net.Dial("tcp", me.Entry.Console.Host + ":" + me.Entry.Console.Port)
		if err != nil {
			err = me.EntityName.ProduceError("VM can't connect to console")
			state = states.StateIdle
			break
		}
		// defer closeDialConnection(conn)
		defer func() {
			if conn != nil {
				_ = conn.Close()
			}
		}()

		eblog.Debug(me.EntityId, "waiting for VM console")

		exitWhen := time.Now().Add(waitFor)
		readBuffer := make([]byte, 512)
		for waitCount := 0; time.Now().Unix() < exitWhen.Unix(); waitCount++ {
			err = conn.SetDeadline(time.Now().Add(me.Entry.Console.ReadWait))
			if err != nil {
				err = me.EntityName.ProduceError("VM console deadline reached")
				state = states.StateStopped	// states.StateUnknown
				break
			}

			bytesRead, err := bufio.NewReader(conn).Read(readBuffer)
			// bytesRead, err := conn.Read(readBuffer)
			// readBuffer, err := bufio.NewReader(conn).ReadString('\n')
			// bytesRead := len(readBuffer)
			if err != nil {
				err = me.EntityName.ProduceError("no VM console data")
				state = states.StateStopped	// states.StateUnknown
				break
			}

			if bytesRead > 0 {
				if showConsole {
					fmt.Printf("%s", string(readBuffer[:bytesRead]))
				}

				err = me.heartbeatOk(readBuffer, bytesRead)
				if err != nil {
					state = states.StateStarted
					break

				//} else {
				//	if me.State.API.WantState == VmStatePowerOff {
				//		me.State.API.CurrentState = VmStateStopping
				//		sts = status.Success("%s API - stopping", global.Brandname)
				//	} else if me.State.API.WantState == VmStateRunning {
				//		me.State.API.CurrentState = VmStateStarting
				//		sts = status.Success("%s API - starting", global.Brandname)
				//	}
				//	// Do not break.
				}
			}

			time.Sleep(me.Entry.Console.WaitDelay)
		}


		//me.State.SetNewState(state, err)
		//me.channels.PublishState(me.State)
		eblog.Debug(me.EntityId, "VM console started OK")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return state, err
}


func (me *Vm) heartbeatOk(b []byte, n int) error {

	var err error

	for range only.Once {
		apiSplit := strings.Split(string(b[:n]), ";")
		if len(apiSplit) <= 1 {
			break
		}

		match, _ := regexp.MatchString(me.Entry.Console.OkString, apiSplit[1])
		if !match {
			break
		}

		// Expecting "1560783374 Gearbox Heartbeat OK"
		//fmt.Printf("API[%d]:%v\n", len(apiSplit), apiSplit)
		switch {
			case len(apiSplit) < 4:
				err = me.EntityName.ProduceError("did not see OK from console - '%s'", string(b[:n]))

			case apiSplit[3] != "OK":
				err = me.EntityName.ProduceError("did not see OK from console - '%s'", string(b[:n]))
		}
	}

	return err
}

