package mqttClient

import "gearbox/heartbeat/eventbroker/messages"


func (me *MqttClient) GetEntities() messages.MessageAddresses {

	var ret messages.MessageAddresses

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for s, _ := range me.services {	// Managed by Mutex
		ret = append(ret, s)
	}

	return ret
}


func (me *MqttClient) GetManagedEntities() messages.MessageAddresses {

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


func (me *MqttClient) DeleteEntity(entity messages.MessageAddress) {

	me.mutex.Lock()
	defer me.mutex.Unlock()
	delete(me.services, entity)	// Managed by Mutex

	return
}


func (me *MqttClient) EnsureDaemonNotNil(entity messages.MessageAddress) error {

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
func (me *MqttClient) IsExisting(him CreateEntry) (*Service, error) {

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


func (me *MqttClient) GetTopics() messages.SubTopics {

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


//func (me *MqttClient) _GetEntities() messages.MessageAddresses {
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
//func (me *MqttClient) _GetManagedEntities() messages.MessageAddresses {
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
//func (me *MqttClient) _GetEntityId(u messages.MessageAddress) (messages.MessageAddress, error) {
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
//func (me *MqttClient) _DeleteEntity(entity messages.MessageAddress) {
//
//	me.mutex.Lock()
//	defer me.mutex.Unlock()
//	delete(me.services, entity)	// Managed by Mutex
//
//	return
//}
//
//
//func (me *MqttClient) _GetIsManaged(u messages.MessageAddress) bool {
//
//	me.mutex.RLock()
//	defer me.mutex.RUnlock()
//	return me.services[u].IsManaged	// Managed by Mutex
//}
//
//
//func (me *MqttClient) _EnsureDaemonNotNil(entity messages.MessageAddress) error {
//
//	var err error
//
//	me.mutex.RLock()
//	defer me.mutex.RUnlock()
//
//	if _, ok := me.services[entity]; !ok {		// Managed by Mutex
//		err = me.EntityId.ProduceError("service doesn't exist")
//	} else {
//		err = me.services[entity].EnsureNotNil()	// Managed by Mutex
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

