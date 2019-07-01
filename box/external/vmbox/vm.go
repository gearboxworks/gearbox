package vmbox

import (
	"bufio"
	"errors"
	"fmt"
	"gearbox/box/external/virtualbox"
	"gearbox/eventbroker/channels"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/osdirs"
	"gearbox/eventbroker/states"
	"github.com/gearboxworks/go-status/only"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"
)

var _ virtualbox.VirtualMachiner = (*Vm)(nil)

type Vm struct {
	// @TODO Three properties that begin with same prefix is a code smell.
	//		 Indicates that this should be a delegated to a different type.
	EntityId        msgs.Address
	EntityName      msgs.Address
	EntityParent    msgs.Address
	State           *states.Status
	ApiState        *states.Status
	ChangeRequested bool
	IsManaged       bool
	Entry           *ServiceConfig // @TODO Why is a ServiceConfig named "Entry?"

	mutex          sync.RWMutex // Mutex control for this struct.
	channels       *channels.Channels
	channelHandler *channels.Subscriber
	osRelease      *Release
	osDirs         *osdirs.BaseDirs
}

func (me *Vm) GetRetryMax() int {
	return me.Entry.retryMax
}

func (me *Vm) GetRetryDelay() time.Duration {
	return me.Entry.retryDelay
}

func (me *Vm) GetReleaser() virtualbox.Releaser {
	return me.osRelease
}

func (me *Vm) GetSsh() virtualbox.SecureSheller {
	ssh := me.Entry.Ssh
	return &ssh
}

func (me *Vm) GetNic() *virtualbox.HostOnlyNic {
	return me.Entry.HostOnlyNic
}

func (me *Vm) SetNic(nic *virtualbox.HostOnlyNic) {
	me.Entry.HostOnlyNic = nic
}

func (me *Vm) GetConsole() virtualbox.Consoler {
	c := me.Entry.Console
	return &c
}

func (me *Vm) GetNics() virtualbox.KeyValuesMap {
	return me.Entry.vmNics
}

func (me *Vm) SetNics(kvm virtualbox.KeyValuesMap) {
	me.Entry.vmNics = kvm
}

func (me *Vm) GetIconFile() string {
	return me.Entry.IconFile
}

func (me *Vm) GetId() string {
	return me.EntityId.String()
}

func (me *Vm) GetName() string {
	return me.EntityName.String()
}

func (me *Vm) GetVmDir() string {
	return me.Entry.VmDir
}

func (me *Vm) SetInfo(kvm virtualbox.KeyValueMap) {
	me.Entry.vmInfo = kvm
}

func (me *Vm) GetInfo() virtualbox.KeyValueMap {
	return me.Entry.vmInfo
}

func (me *Vm) GetInfoValue(k string) string {
	v, _ := me.Entry.vmInfo[k]
	return v
}

/////////////////////////

func (me *Vm) GetIsManaged() bool {

	me.mutex.RLock()
	defer me.mutex.RUnlock()
	return me.IsManaged // Managed by Mutex
}

func (me *Vm) GetEntityId() (msgs.Address, error) {

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	err := me.EnsureNotNil()
	if err != nil {
		return "", err
	}

	return me.EntityId, err // Managed by Mutex
}

func (me *Vm) GetConfig() (ServiceConfig, error) {

	var sc ServiceConfig

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	err := me.EnsureNotNil()
	if err != nil {
		return sc, err
	}

	return sc, err // Managed by Mutex
}

func (me *Vm) GetStatus() (*states.Status, error) {

	var sc *states.Status

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	err := me.EnsureNotNil()
	if err == nil {
		sc = me.State // Managed by Mutex
	}

	return sc, err
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

		ok, err = virtualbox.StartVm(me, true)
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
	eblog.LogIfError(err)

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

		ok, err = virtualbox.StopVm(me, false, true)
		if !ok {
			ok, err = virtualbox.StopVm(me, true, true)
			break
		}

		me.State.SetNewState(states.StateStopped, err)
		me.channels.PublishState(me.State)

		me.ApiState.SetNewState(states.StateStopped, err)
		me.channels.PublishState(me.ApiState)

		eblog.Debug(me.EntityId, "stopped VM OK")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

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
		v, err = virtualbox.VmStatus(me)
		me.State.SetNewState(v, err)
		me.channels.PublishState(me.State)

		var a states.State
		a, err = me.waitForApiState(DefaultRunWaitTime, false)
		me.ApiState.SetNewState(a, err)
		me.channels.PublishState(me.ApiState)

		eblog.Debug(me.EntityId, "VM is in state %s, API is in state %s", v, a)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return err
}

func (me *Vm) Print() error {

	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("software error")
			break
		}

		if me.IsManaged == true {
			fmt.Printf("# Entry(deleted): %v", me.EntityId)
		} else {
			fmt.Printf("# Entry: %v", me.EntityId)
		}
		//err = me.Entry.Print()
		//if err != nil {
		//	break
		//}
	}

	return err
}

func (me *Vm) EnsureNotNil() error {
	var err error

	switch {
	case me == nil:
		err = errors.New("VmBox Service instance is nil")
	}

	return err
}

/////////[ VirtualBox specific ]//////////

/////////[ Commands ]//////////

//////////////[ private funcs ]//////////////
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
			err = msgs.MakeError(me.EntityName, "no VM console host defined")
			state = states.StateError
			break
		}

		if me.Entry.Console.Port == "" {
			err = msgs.MakeError(me.EntityName, "no VM console port defined")
			state = states.StateError
			break
		}

		// Ensure we only have one at a time.
		me.Entry.Console.mutex.Lock()
		defer me.Entry.Console.mutex.Unlock()

		// Connect to this console
		conn, err := net.Dial("tcp", me.Entry.Console.Host+":"+me.Entry.Console.Port)
		if err != nil {
			err = msgs.MakeError(me.EntityName, "VM can't connect to console")
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
				err = msgs.MakeError(me.EntityName, "VM console deadline reached")
				state = states.StateStopped // states.StateUnknown
				break
			}

			bytesRead, err := bufio.NewReader(conn).Read(readBuffer)
			// bytesRead, err := conn.Read(readBuffer)
			// readBuffer, err := bufio.NewReader(conn).ReadString('\n')
			// bytesRead := len(readBuffer)
			if err != nil {
				err = msgs.MakeError(me.EntityName, "no VM console data")
				state = states.StateStopped // states.StateUnknown
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
	eblog.LogIfError(err)

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
			err = msgs.MakeError(me.EntityName, "did not see OK from console - '%s'", string(b[:n]))

		case apiSplit[3] != "OK":
			err = msgs.MakeError(me.EntityName, "did not see OK from console - '%s'", string(b[:n]))
		}
	}

	return err
}
