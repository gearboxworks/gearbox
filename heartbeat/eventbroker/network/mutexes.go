package network

import "gearbox/heartbeat/eventbroker/messages"


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


func (me *ZeroConf) DeleteEntity(entity messages.MessageAddress) {

	me.mutex.Lock()
	defer me.mutex.Unlock()
	delete(me.services, entity)	// Managed by Mutex

	return
}


func (me *ZeroConf) EnsureDaemonNotNil(entity messages.MessageAddress) error {

	var err error

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	if _, ok := me.services[entity]; !ok {		// Managed by Mutex
		err = me.EntityId.ProduceError("service doesn't exist")
	} else {
		err = me.services[entity].EnsureNotNil()	// Managed by Mutex
	}

	return err
}


// Ensure we don't duplicate services.
func (me *ZeroConf) IsExisting(him CreateEntry) (*Service, error) {

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

