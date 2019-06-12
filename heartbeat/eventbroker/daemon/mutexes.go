package daemon

import (
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/only"
	"gearbox/heartbeat/eventbroker/states"
	"time"
)


func (me *Daemon) GetEntities() messages.MessageAddresses {

	var ret messages.MessageAddresses

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for s, _ := range me.daemons {	// Managed by Mutex
		ret = append(ret, s)
	}

	return ret
}


func (me *Daemon) GetManagedEntities() messages.MessageAddresses {

	var ret messages.MessageAddresses

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for s, _ := range me.daemons {	// Managed by Mutex
		if me.daemons[s].IsManaged {	// Managed by Mutex
			ret = append(ret, s)
		}
	}

	return ret
}


func (me *Daemon) AddEntity(entity messages.MessageAddress, sc *Service) error {
	var err error

	me.mutex.Lock()
	defer me.mutex.Unlock()

	if _, ok := me.daemons[entity]; !ok { // Managed by Mutex
		me.daemons[entity] = sc
	} else {
		err = me.EntityId.ProduceError("service %s already exists", entity)
	}

	return err
}


func (me *Daemon) DeleteEntity(entity messages.MessageAddress) error {

	var err error

	me.mutex.Lock()
	defer me.mutex.Unlock()

	for range only.Once {
		if _, ok := me.daemons[entity]; !ok { // Managed by Mutex
			err = me.EntityId.ProduceError("service doesn't exist")
			break
		}

		delete(me.daemons, entity) // Managed by Mutex
	}

	return err
}


func (me *Daemon) EnsureDaemonNotNil(entity messages.MessageAddress) error {

	var err error

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	if _, ok := me.daemons[entity]; !ok {		// Managed by Mutex
		err = me.EntityId.ProduceError("service doesn't exist")
	} else {
		err = me.daemons[entity].EnsureNotNil()	// Managed by Mutex
	}

	return err
}


// Ensure we don't duplicate services.
func (me *Daemon) FindExistingConfig(him ServiceConfig) (*Service, error) {

	var err error
	var sc *Service

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for _, ce := range me.daemons {	// Managed by Mutex
		err = ce.IsExisting(him)
		if err != nil {
			sc = ce
			break
		}
	}

	return sc, err
}


// Ensure we don't duplicate services.
func (me *Daemon) IsExisting(s messages.MessageAddress) *Service {

	var sc *Service

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for _, sc = range me.daemons {	// Managed by Mutex
		if sc.EntityId == s {
			break
		}
	}

	return sc
}


func (me *Daemon) GetTopics() messages.SubTopics {

	return me.channelHandler.GetTopics()
}


func (me *Service) GetIsManaged() bool {

	me.mutex.RLock()
	defer me.mutex.RUnlock()
	return me.IsManaged	// Managed by Mutex
}


func (me *Service) GetEntityId() (messages.MessageAddress, error) {

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	err := me.EnsureNotNil()
	if err != nil {
		return "", err
	}

	return me.EntityId, err		// Managed by Mutex
}


func (me *Service) GetConfig() (ServiceConfig, error) {

	var sc ServiceConfig

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	err := me.EnsureNotNil()
	if err != nil {
		return sc, err
	}

	return sc, err		// Managed by Mutex
}


func (me *Service) GetStatus() (*states.Status, error) {

	var sc *states.Status

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	err := me.EnsureNotNil()
	if err == nil {
		sc = &me.State		// Managed by Mutex
	}

	return sc, err
}


func (me *Daemon) GetServiceFiles() map[string]time.Time {

	jc := make(map[string]time.Time)

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for _, ce := range me.daemons {	// Managed by Mutex
		jc[ce.JsonFile.Name] = ce.JsonFile.LastModTime
	}

	return jc
}


//
//
//func (me *Daemon) _GetEntityId(u messages.MessageAddress) (messages.MessageAddress, error) {
//
//	me.mutex.RLock()
//	defer me.mutex.RUnlock()
//
//	err := me.EnsureNotNil(u)
//	if err != nil {
//		return "", err
//	}
//
//	return me.daemons[u].EntityId, err	// Managed by Mutex
//}

