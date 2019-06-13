package network

import (
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/only"
	"gearbox/heartbeat/eventbroker/states"
)


func (me *ZeroConf) GetEntities() messages.MessageAddresses {

	var ret messages.MessageAddresses

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for s, _ := range me.services {	// Managed by Mutex
		ret = append(ret, s)
	}

	return ret
}


func (me *ZeroConf) GetManagedEntities() messages.MessageAddresses {

	var ret messages.MessageAddresses

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for s, _ := range me.services {	// Managed by Mutex
		if me.services[s].IsManaged {	// Managed by Mutex
			ret = append(ret, s)
		}
	}

	return ret
}


func (me *ZeroConf) AddEntity(client messages.MessageAddress, sc *Service) error {
	var err error

	me.mutex.Lock()
	defer me.mutex.Unlock()

	if _, ok := me.services[client]; !ok { // Managed by Mutex
		me.services[client] = sc
	} else {
		err = me.EntityId.ProduceError("service %s already exists", client)
	}

	return err
}


func (me *ZeroConf) DeleteEntity(client messages.MessageAddress) error {

	var err error

	me.mutex.Lock()
	defer me.mutex.Unlock()

	for range only.Once {
		if _, ok := me.services[client]; !ok { // Managed by Mutex
			err = me.EntityId.ProduceError("service doesn't exist")
			break
		}

		delete(me.services, client) // Managed by Mutex
	}

	return err
}


func (me *ZeroConf) EnsureDaemonNotNil(client messages.MessageAddress) error {

	var err error

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	if _, ok := me.services[client]; !ok {		// Managed by Mutex
		err = me.EntityId.ProduceError("service doesn't exist")
	} else {
		err = me.services[client].EnsureNotNil()	// Managed by Mutex
	}

	return err
}


// Ensure we don't duplicate services.
func (me *ZeroConf) FindExistingConfig(him ServiceConfig) (*Service, error) {

	var err error
	var sc *Service

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for _, ce := range me.services {	// Managed by Mutex
		err = ce.IsExisting(him)
		if err != nil {
			sc = ce
			break
		}
	}

	return sc, err
}


// Ensure we don't duplicate services.
func (me *ZeroConf) IsExisting(s messages.MessageAddress) *Service {

	var sc *Service

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for _, sc = range me.services {	// Managed by Mutex
		if sc.EntityId == s {
			break
		}
	}

	return sc
}


func (me *ZeroConf) GetTopics() messages.SubTopics {

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
		sc = me.State		// Managed by Mutex
	}

	return sc, err
}

