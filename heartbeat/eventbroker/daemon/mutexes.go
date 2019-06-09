package daemon

import "gearbox/heartbeat/eventbroker/messages"


func (me *Daemon) GetIsManaged(u messages.MessageAddress) bool {

	me.mutex.RLock()
	defer me.mutex.RUnlock()
	return me.daemons[u].IsManaged	// Managed by Mutex
}


func (me *Daemon) GetAllDaemonUuids() []messages.MessageAddress {

	var ret []messages.MessageAddress

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for u, _ := range me.daemons {	// Managed by Mutex
		ret = append(ret, u)
	}

	return ret
}


func (me *Daemon) GetManagedDaemonUuids() []messages.MessageAddress {

	var ret []messages.MessageAddress

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for u, _ := range me.daemons {	// Managed by Mutex
		if me.daemons[u].IsManaged {	// Managed by Mutex
			ret = append(ret, u)
		}
	}

	return ret
}


func (me *Daemon) GetEntityId(u messages.MessageAddress) messages.MessageAddress {

	me.mutex.RLock()
	defer me.mutex.RUnlock()
	return me.daemons[u].EntityId	// Managed by Mutex
}


func (me *Daemon) EnsureDaemonNotNil(u messages.MessageAddress) error {

	var err error

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	_, ok := me.daemons[u]	// Managed by Mutex
	if !ok {
		err = me.EntityId.ProduceError("service doesn't exist")
	} else {
		err = me.daemons[u].EnsureNotNil()	// Managed by Mutex
	}

	return err
}


// Mutex handling.
func (me *Daemon) DeleteDaemon(u messages.MessageAddress) {

	me.mutex.Lock()
	defer me.mutex.Unlock()
	delete(me.daemons, u)	// Managed by Mutex

	return
}

