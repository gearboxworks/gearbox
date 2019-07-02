package vmbox

import (
	"errors"
	"fmt"
	"gearbox/box/external/virtualbox"
	"gearbox/eventbroker/channels"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/osdirs"
	"gearbox/eventbroker/states"
	"gearbox/eventbroker/tasks"
	"gearbox/global"
	"github.com/gearboxworks/go-status/only"
	"sync"
	"time"
)

type VmMap map[msgs.Address]*Vm

func (me *VmMap) Print() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for u, s := range *me {
			fmt.Printf("# Entry: %s\n", u)
			err = s.Print()
			if err != nil {
				break
			}
		}
	}

	return err
}

type VmBox struct {
	EntityId     msgs.Address
	EntityName   msgs.Address
	EntityParent msgs.Address
	Boxname      string
	State        *states.Status
	Task         *tasks.Task
	Channels     *channels.Channels
	Releases     *Releases

	mutex           sync.RWMutex // Mutex control for map.
	channelHandler  *channels.Subscriber
	restartAttempts int
	waitTime        time.Duration
	vms             VmMap
	OsPaths         *osdirs.BaseDirs
}
type Args VmBox

func New(args ...Args) (*VmBox, error) {

	var _args Args
	var err error

	me := &VmBox{}

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		//foo := Args{}
		//err = copier.Copy(&foo, &_args)
		//if err != nil {
		//	err = msgs.MakeError(_args.EntityId,"unable to copy config args")
		//	break
		//}

		if _args.Channels == nil {
			err = msgs.MakeError(_args.EntityId, "channel pointer is nil")
			break
		}

		if _args.OsPaths == nil {
			err = msgs.MakeError(_args.EntityId, "ospaths is nil")
			break
		}

		if _args.EntityId == "" {
			_args.EntityId = entity.VmBoxEntityName
		}

		if _args.EntityName == "" {
			_args.EntityName = _args.EntityId
		}

		if _args.EntityParent == "" {
			_args.EntityParent = _args.EntityId
		}

		_args.State = states.New(_args.EntityId, _args.EntityId, entity.SelfEntityName)

		if _args.Boxname == "" {
			_args.Boxname = global.Brandname
		}

		if _args.waitTime == 0 {
			_args.waitTime = DefaultVmWaitTime
		}

		_args.Releases, _ = NewReleases(_args.Channels)

		_args.vms = make(VmMap)

		*me = VmBox(_args)

		me.SetWantState(states.StateIdle)
		me.SetStateError(err)
		me.SetState(states.StateIdle)
		eblog.Debug(me.EntityId, "init complete")
	}

	me.PublishChannelState(me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return me, err
}

// Start the VmBox handler.
func (me *VmBox) Start() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.SetStateAction(states.ActionStart)
		me.PublishChannelState(me.State)

		for range only.Once {
			me.Task, err = tasks.StartTask(initVmBox, startVmBox, monitorVmBox, stopVmBox, me)
			me.SetStateError(err)
			if err != nil {
				eblog.LogIfError(err)
				break
			}
		}

		if err == nil {
			me.SetState(states.StateStarted)
		}

		me.PublishChannelState(me.State)
		eblog.Debug(me.EntityId, "started task handler")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return err
}

// Stop the VmBox handler.
func (me *VmBox) Stop() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.SetStateAction(states.ActionStop)
		me.PublishChannelState(me.State)

		err = me.StopVms()
		if err != nil {
			eblog.LogIfError(err)
		}

		err = me.Task.Stop()
		if err != nil {
			me.SetStateError(err)
		} else {
			me.SetState(states.StateStopped)
		}

		me.PublishChannelState(me.State)
		eblog.Debug(me.EntityId, "stopped task handler")

	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return err
}

func (me *VmBox) StopVms() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for u := range me.vms {
			if !me.vms[u].IsManaged {
				continue
			}
			err = me.vms[u].Stop()
			eblog.LogIfError(err)
		}
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return err
}

func (me *VmBox) EnsureNotNil() (err error) {
	if me == nil {
		err = errors.New("VmBox instance is nil")
	}
	return err
}

func (me *VmBox) SetStateError(err error) {
	me.State.Error = err
}

func (me *VmBox) SetState(s states.State) {
	me.State.SetState(s)
}

func (me *VmBox) SetStateAction(a states.Action) {
	me.State.SetNewAction(a)
}

func (me *VmBox) PublishChannelState(state *states.Status) {
	me.Channels.PublishState(me.State)
}

func (me *VmBox) SetWantState(s states.State) {
	me.State.SetWant(s)
}

func (me *VmBox) MakeVm(c ServiceConfig) (*Vm, error) {

	var err error

	vm := Vm{}

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		vm.Entry = &c

		err = me.OsPaths.EnsureNotNil()
		if err != nil {
			break
		}
		vm.osDirs = me.OsPaths
		vm.channels = me.Channels

		// Check config, set defaults if needed.
		err = vm.VerifyConfig()
		if err != nil {
			break
		}

		// Load up config file, if it exists. Supersceding previous config.
		err = vm.ConfigExists()
		if err == nil {
			err = vm.ReadConfig()
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
		rel, err = me.Releases.SelectRelease(ReleaseSelector{SpecificVersion: vm.Entry.Version})
		if err != nil {
			break
		}

		// Fetch ISO.
		if !me.Releases.Selected.IsDownloading {
			err = me.Releases.Selected.GetIso()
		}

		// Update specific version fetched.
		vm.Entry.Version = string(rel.Version)

		vm.EntityId = msgs.MakeAddress()
		vm.EntityName = vm.Entry.Name
		vm.EntityParent = me.EntityId
		vm.State = states.New(vm.EntityId, vm.EntityName, me.EntityId)
		vm.State.SetNewAction(states.ActionStop)
		a := msgs.Address(entity.ApiEntityName)
		vm.ApiState = states.New(vm.EntityId, a, me.EntityId)
		vm.ApiState.SetNewAction(states.ActionStop)
		vm.IsManaged = true
		vm.osRelease = rel

		vm.channels.PublishState(vm.State)

		var state states.State
		state, err = virtualbox.CreateVm(&vm)
		if err != nil {
			fmt.Printf("Gearbox: Error creating VM %s with '%v'\n", vm.EntityName, err)
			break
		}

		err = vm.WriteConfig()
		if err != nil {
			break
		}

		switch state {
		case states.StateError:
			eblog.Debug(me.EntityId, "%v", err)

		case states.StateStopped:
			err = me.AddEntity(vm.EntityId, &vm)
			if err != nil {
				break
			}
			vm.State.SetNewAction(states.ActionStop)
			eblog.Debug(me.EntityId, "VM registered OK")

		case states.StateStarted:
			// VM already created but started.
			err = me.AddEntity(vm.EntityId, &vm)
			if err != nil {
				break
			}
			vm.State.SetNewAction(states.ActionStart)

		case states.StateUnregistered:
			eblog.Debug(me.EntityId, "VM not created")

		case states.StateUnknown:
			eblog.Debug(me.EntityId, "%v", err)
		}

		vm.State.SetNewState(state, err)
		vm.channels.PublishState(vm.State)
		eblog.Debug(me.EntityId, "registered VM OK")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return &vm, err
}

// Ensure we don't duplicate services.
func (me *VmBox) IsExisting(client msgs.Address) *Vm {

	var ret *Vm

	for _, v := range me.vms {
		if v.EntityName == client {
			ret = v
			break
		}

		if v.EntityId == client {
			ret = v
			break
		}
	}
	return ret
}

func (me *VmMap) EnsureNotNil() (err error) {
	if me == nil {
		err = errors.New("VmBox VmMap instance is nil")
	}
	return err
}

///////////[ Mutexes ]///////////
func (me *VmBox) AddEntity(client msgs.Address, sc *Vm) error {
	var err error

	me.mutex.Lock()
	defer me.mutex.Unlock()

	if _, ok := me.vms[client]; !ok { // Managed by Mutex
		me.vms[client] = sc
	} else {
		err = msgs.MakeError(me.EntityId, "VM %s already exists", client)
	}

	return err
}

func (me *VmBox) DeleteEntity(client msgs.Address) error {

	var err error

	me.mutex.Lock()
	defer me.mutex.Unlock()

	for range only.Once {
		if _, ok := me.vms[client]; !ok { // Managed by Mutex
			err = msgs.MakeError(me.EntityId, "VM doesn't exist")
			break
		}

		delete(me.vms, client) // Managed by Mutex
	}

	return err
}

//
//// Ensure we don't duplicate services.
//func (me *VmBox) IsExisting(him ServiceConfig) error {
//
//	var err error
//
//	for _, ce := range me.vms {
//		err = ce.IsExisting(him)
//		if err != nil {
//			break
//		}
//	}
//
//	return err
//}

////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

//// Unsubscribe a service by method defined by a UUID reference.
//func (me *VmBox) UnsubscribeByUuid(client msg.MessageAddress) error {
//
//	var err error
//
//	for range only.Once {
//		err = me.EnsureNotEmpty()
//		if err != nil {
//			break
//		}
//
//		err = me.services[client].EnsureNotEmpty()
//		if err != nil {
//			break
//		}
//
//		me.services[client].State.SetNewAction(states.ActionStop)	// Was states.ActionUnsubscribe
//		me.services[client].channels.PublishState(me.State)
//
//		// Do something
//
//		me.services[client].State.SetNewState(states.StateStopped, err)	// Was states.StateUnsubscribed
//		me.services[client].channels.PublishState(me.services[client].State)
//
//		err = me.DeleteEntity(client)
//		if err != nil {
//			break
//		}
//
//		//me.Channels.PublishSpecificState(&client, states.State(states.StateUnsubscribed))
//		eblog.Debug(me.EntityId, "unregistered service %s OK", client.String())
//	}
//
//	me.Channels.PublishState(me.State)
//	eblog.LogIfNil(me, err)
//	eblog.LogIfError(me.EntityId, err)
//
//	return err
//}
//
//// Unsubscribe a service via a channel defined by a UUID reference.
//func (me *VmBox) UnsubscribeByChannel(caller msg.MessageAddress, u msg.MessageAddress) error {
//
//	var err error
//
//	for range only.Once {
//		err = me.EnsureNotEmpty()
//		if err != nil {
//			break
//		}
//
//		//unreg := me.EntityId.Construct(me.EntityId, states.ActionUnsubscribe, msg.MessageText(u.String()))
//		unreg := caller.MakeMessage(me.EntityId, states.ActionUnsubscribe, msg.MessageText(u.String()))
//		err = me.Channels.Publish(unreg)
//		if err != nil {
//			break
//		}
//
//		eblog.Debug(me.EntityId, "unsubscribed service by channel %s OK", u.String())
//	}
//
//	eblog.LogIfNil(me, err)
//	eblog.LogIfError(me.EntityId, err)
//
//	return err
//}
//
//// Register a service by method defined by a *NewTopic structure.
//func (me *VmBox) Subscribe(ce ServiceConfig) (*Service, error) {
//
//	var err error
//	var sc Service
//
//	for range only.Once {
//		err = me.EnsureNotEmpty()
//		if err != nil {
//			break
//		}
//
//		err = me.services.IsExisting(ce)
//		if err != nil {
//			break
//		}
//
//		// Create new client entry.
//		sc.EntityId = *msg.GenerateAddress()
//		sc.EntityName = msg.MessageAddress(ce.Name)
//		sc.EntityParent = &me.EntityId
//		sc.State = states.New(&sc.EntityId, &sc.EntityName, me.EntityId)
//		sc.State.SetNewAction(states.ActionSubscribe)
//		sc.IsManaged = true
//		sc.channels = me.Channels
//		sc.channels.PublishState(sc.State)
//
//		err = me.AddEntity(sc.EntityId, &sc)
//		if err != nil {
//			break
//		}
//
//		sc.State.SetNewState(states.StateSubscribed, err)
//		sc.channels.PublishState(sc.State)
//		eblog.Debug(me.EntityId, "subscribed %s OK", sc.EntityId.String())
//	}
//
//	me.Channels.PublishState(me.State)
//	eblog.LogIfNil(me, err)
//	eblog.LogIfError(me.EntityId, err)
//
//	return &sc, err
//}
//func (me *VmBox) GetEntities() msg.Addresses {
//
//	var ret msg.Addresses
//
//	me.mutex.RLock()
//	defer me.mutex.RUnlock()
//
//	for s, _ := range me.services {	// Managed by Mutex
//		ret = append(ret, s)
//	}
//
//	return ret
//}
//
//
//func (me *VmBox) GetManagedEntities() msg.Addresses {
//
//	var ret msg.Addresses
//
//	me.mutex.RLock()
//	defer me.mutex.RUnlock()
//
//	for s, _ := range me.services {	// Managed by Mutex
//		if me.services[s].IsManaged {	// Managed by Mutex
//			ret = append(ret, s)
//		}
//	}
//
//	return ret
//}
//
//
//func (me *VmBox) EnsureDaemonNotNil(client msg.Address) error {
//
//	var err error
//
//	me.mutex.RLock()
//	defer me.mutex.RUnlock()
//
//	if _, ok := me.services[client]; !ok {		// Managed by Mutex
//		err = msgs.MakeError(me.EntityId,"service doesn't exist")
//	} else {
//		err = me.services[client].EnsureNotEmpty()	// Managed by Mutex
//	}
//
//	return err
//}
//
//
//// Ensure we don't duplicate services.
//func (me *VmBox) FindExistingConfig(him ServiceConfig) (*Service, error) {
//
//	var err error
//	var sc *Service
//
//	me.mutex.RLock()
//	defer me.mutex.RUnlock()
//
//	for _, ce := range me.services {	// Managed by Mutex
//		err = ce.IsExisting(him)
//		if err != nil {
//			sc = ce
//			break
//		}
//	}
//
//	return sc, err
//}
//
//
//// Ensure we don't duplicate services.
//func (me *VmBox) IsExisting(s msg.Address) *Service {
//
//	var sc *Service
//
//	me.mutex.RLock()
//	defer me.mutex.RUnlock()
//
//	for _, sc = range me.services {	// Managed by Mutex
//		if sc.EntityId == s {
//			break
//		}
//	}
//
//	return sc
//}
//
//
//func (me *VmBox) GetTopics() msg.SubTopics {
//
//	return me.channelHandler.GetTopics()
//}
//
//
//func (me *VmBox) _GetEntities() msg.Addresses {
//
//	var ret msg.Addresses
//
//	me.mutex.RLock()
//	defer me.mutex.RUnlock()
//
//	for u, _ := range me.services {	// Managed by Mutex
//		ret = append(ret, u)
//	}
//
//	return ret
//}
//
//
//func (me *VmBox) _GetManagedEntities() msg.Addresses {
//
//	var ret msg.Addresses
//
//	me.mutex.RLock()
//	defer me.mutex.RUnlock()
//
//	for u, _ := range me.services {	// Managed by Mutex
//		if me.services[u].IsManaged {	// Managed by Mutex
//			ret = append(ret, u)
//		}
//	}
//
//	return ret
//}
//
//
//func (me *VmBox) _GetEntityId(u msg.Address) (msg.Address, error) {
//
//	me.mutex.RLock()
//	defer me.mutex.RUnlock()
//
//	err := me.EnsureDaemonNotNil(u)
//	if err != nil {
//		return "", err
//	}
//
//	return me.services[u].EntityId, err	// Managed by Mutex
//}
//
//
//func (me *VmBox) _DeleteEntity(client msg.Address) {
//
//	me.mutex.Lock()
//	defer me.mutex.Unlock()
//	delete(me.services, client)	// Managed by Mutex
//
//	return
//}
//
//
//func (me *VmBox) _GetIsManaged(u msg.Address) bool {
//
//	me.mutex.RLock()
//	defer me.mutex.RUnlock()
//	return me.services[u].IsManaged	// Managed by Mutex
//}
//
//
//func (me *VmBox) _EnsureDaemonNotNil(client msg.Address) error {
//
//	var err error
//
//	me.mutex.RLock()
//	defer me.mutex.RUnlock()
//
//	if _, ok := me.services[client]; !ok {		// Managed by Mutex
//		err = msgs.MakeError(me.EntityId,"service doesn't exist")
//	} else {
//		err = me.services[client].EnsureNotEmpty()	// Managed by Mutex
//	}
//
//	return err
//}
//
