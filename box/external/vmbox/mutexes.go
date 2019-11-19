package vmbox

import (
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/states"
	"github.com/gearboxworks/go-status/only"
)


func (me *VmBox) AddEntity(client messages.MessageAddress, sc *Vm) error {
	var err error

	me.mutex.Lock()
	defer me.mutex.Unlock()

	if _, ok := me.vms[client]; !ok { // Managed by Mutex
		me.vms[client] = sc
	} else {
		err = me.EntityId.ProduceError("VM %s already exists", client)
	}

	return err
}


func (me *VmBox) DeleteEntity(client messages.MessageAddress) error {

	var err error

	me.mutex.Lock()
	defer me.mutex.Unlock()

	for range only.Once {
		if _, ok := me.vms[client]; !ok { // Managed by Mutex
			err = me.EntityId.ProduceError("VM doesn't exist")
			break
		}

		delete(me.vms, client) // Managed by Mutex
	}

	return err
}


func (me *Vm) GetIsManaged() bool {

	me.mutex.RLock()
	defer me.mutex.RUnlock()
	return me.IsManaged	// Managed by Mutex
}


func (me *Vm) GetEntityId() (messages.MessageAddress, error) {

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	err := me.EnsureNotNil()
	if err != nil {
		return "", err
	}

	return me.EntityId, err		// Managed by Mutex
}


func (me *Vm) GetConfig() (ServiceConfig, error) {

	var sc ServiceConfig

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	err := me.EnsureNotNil()
	if err != nil {
		return sc, err
	}

	return sc, err		// Managed by Mutex
}


func (me *Vm) GetStatus() (*states.Status, error) {

	var sc *states.Status

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	err := me.EnsureNotNil()
	if err == nil {
		sc = me.State		// Managed by Mutex
	}

	return sc, err
}


//func (me *VmBox) GetEntities() messages.MessageAddresses {
//
//	var ret messages.MessageAddresses
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
//func (me *VmBox) GetManagedEntities() messages.MessageAddresses {
//
//	var ret messages.MessageAddresses
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
//func (me *VmBox) EnsureDaemonNotNil(client messages.MessageAddress) error {
//
//	var err error
//
//	me.mutex.RLock()
//	defer me.mutex.RUnlock()
//
//	if _, ok := me.services[client]; !ok {		// Managed by Mutex
//		err = me.EntityId.ProduceError("service doesn't exist")
//	} else {
//		err = me.services[client].EnsureNotNil()	// Managed by Mutex
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
//func (me *VmBox) IsExisting(s messages.MessageAddress) *Service {
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
//func (me *VmBox) GetTopics() messages.SubTopics {
//
//	return me.channelHandler.GetTopics()
//}
//
//
//func (me *VmBox) _GetEntities() messages.MessageAddresses {
//
//	var ret messages.MessageAddresses
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
//func (me *VmBox) _GetManagedEntities() messages.MessageAddresses {
//
//	var ret messages.MessageAddresses
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
//func (me *VmBox) _GetEntityId(u messages.MessageAddress) (messages.MessageAddress, error) {
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
//func (me *VmBox) _DeleteEntity(client messages.MessageAddress) {
//
//	me.mutex.Lock()
//	defer me.mutex.Unlock()
//	delete(me.services, client)	// Managed by Mutex
//
//	return
//}
//
//
//func (me *VmBox) _GetIsManaged(u messages.MessageAddress) bool {
//
//	me.mutex.RLock()
//	defer me.mutex.RUnlock()
//	return me.services[u].IsManaged	// Managed by Mutex
//}
//
//
//func (me *VmBox) _EnsureDaemonNotNil(client messages.MessageAddress) error {
//
//	var err error
//
//	me.mutex.RLock()
//	defer me.mutex.RUnlock()
//
//	if _, ok := me.services[client]; !ok {		// Managed by Mutex
//		err = me.EntityId.ProduceError("service doesn't exist")
//	} else {
//		err = me.services[client].EnsureNotNil()	// Managed by Mutex
//	}
//
//	return err
//}
//
//
//func (me *Service) _GetTopics() (messages.SubTopics, error) {
//
//	var ret messages.SubTopics
//	var err error
//
//	err = me.EnsureNotNil()
//	if err != nil {
//		return ret, err
//	}
//
//	ret = me.channelHandler.GetTopics()
//
//	return ret, err
//}

