package mqttClient

import (
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/states"
	"github.com/gearboxworks/go-status/only"
)

func (me *MqttClient) GetEntities() msgs.Addresses {

	var ret msgs.Addresses

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for s := range me.services { // Managed by Mutex
		ret = append(ret, s)
	}

	return ret
}

func (me *MqttClient) GetManagedEntities() msgs.Addresses {

	var ret msgs.Addresses

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for s := range me.services { // Managed by Mutex
		if me.services[s].IsManaged { // Managed by Mutex
			ret = append(ret, s)
		}
	}

	return ret
}

func (me *MqttClient) AddEntity(client msgs.Address, sc *Service) error {
	var err error

	me.mutex.Lock()
	defer me.mutex.Unlock()

	if _, ok := me.services[client]; !ok { // Managed by Mutex
		me.services[client] = sc
	} else {
		err = msgs.MakeError(me.EntityId, "service %s already exists", client)
	}

	return err
}

func (me *MqttClient) DeleteEntity(client msgs.Address) error {

	var err error

	me.mutex.Lock()
	defer me.mutex.Unlock()

	for range only.Once {
		if _, ok := me.services[client]; !ok { // Managed by Mutex
			err = msgs.MakeError(me.EntityId, "service doesn't exist")
			break
		}

		delete(me.services, client) // Managed by Mutex
	}

	return err
}

func (me *MqttClient) EnsureDaemonNotNil(client msgs.Address) error {

	var err error

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	if _, ok := me.services[client]; !ok { // Managed by Mutex
		err = msgs.MakeError(me.EntityId, "service doesn't exist")
	} else {
		err = me.services[client].EnsureNotNil() // Managed by Mutex
	}

	return err
}

// Ensure we don't duplicate services.
func (me *MqttClient) FindExistingConfig(him ServiceConfig) (*Service, error) {

	var err error
	var sc *Service

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for _, ce := range me.services { // Managed by Mutex
		err = ce.IsExisting(him)
		if err != nil {
			sc = ce
			break
		}
	}

	return sc, err
}

// Ensure we don't duplicate services.
func (me *MqttClient) IsExisting(s msgs.Address) *Service {

	var sc *Service

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for _, sc = range me.services { // Managed by Mutex
		if sc.EntityId == s {
			break
		}
	}

	return sc
}

func (me *MqttClient) GetTopics() msgs.SubTopics {

	return me.channelHandler.GetTopics()
}

func (me *Service) GetIsManaged() bool {

	me.mutex.RLock()
	defer me.mutex.RUnlock()
	return me.IsManaged // Managed by Mutex
}

func (me *Service) GetEntityId() (msgs.Address, error) {

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	err := me.EnsureNotNil()
	if err != nil {
		return "", err
	}

	return me.EntityId, err // Managed by Mutex
}

func (me *Service) GetConfig() (ServiceConfig, error) {

	var sc ServiceConfig

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	err := me.EnsureNotNil()
	if err != nil {
		return sc, err
	}

	return sc, err // Managed by Mutex
}

func (me *Service) GetStatus() (*states.Status, error) {

	var sc *states.Status

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	err := me.EnsureNotNil()
	if err == nil {
		sc = me.State // Managed by Mutex
	}

	return sc, err
}

//func (me *MqttClient) _GetEntities() msg.Addresses {
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
//func (me *MqttClient) _GetManagedEntities() msg.Addresses {
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
//func (me *MqttClient) _GetEntityId(u msg.Address) (msg.Address, error) {
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
//func (me *MqttClient) _DeleteEntity(client msg.Address) {
//
//	me.mutex.Lock()
//	defer me.mutex.Unlock()
//	delete(me.services, client)	// Managed by Mutex
//
//	return
//}
//
//
//func (me *MqttClient) _GetIsManaged(u msg.Address) bool {
//
//	me.mutex.RLock()
//	defer me.mutex.RUnlock()
//	return me.services[u].IsManaged	// Managed by Mutex
//}
//
//
//func (me *MqttClient) _EnsureDaemonNotNil(client msg.Address) error {
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
//func (me *Service) _GetTopics() (msg.SubTopics, error) {
//
//	var ret msg.SubTopics
//	var err error
//
//	err = me.EnsureNotEmpty()
//	if err != nil {
//		return ret, err
//	}
//
//	ret = me.channelHandler.GetTopics()
//
//	return ret, err
//}
