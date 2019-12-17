package daemon

import (
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/states"
	"github.com/gearboxworks/go-status/only"
	"time"
)

func (me *Daemon) GetEntities() msgs.Addresses {

	var ret msgs.Addresses

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for s := range me.daemons { // Managed by Mutex
		ret = append(ret, s)
	}

	return ret
}

func (me *Daemon) GetManagedEntities() msgs.Addresses {

	var ret msgs.Addresses

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for s := range me.daemons { // Managed by Mutex
		if me.daemons[s].IsManaged { // Managed by Mutex
			ret = append(ret, s)
		}
	}

	return ret
}

func (me *Daemon) AddEntity(client msgs.Address, sc *Service) error {
	var err error

	me.mutex.Lock()
	defer me.mutex.Unlock()

	if _, ok := me.daemons[client]; !ok { // Managed by Mutex
		me.daemons[client] = sc
	} else {
		err = msgs.MakeError(me.EntityId, "service %s already exists", client)
	}

	return err
}

func (me *Daemon) DeleteEntity(client msgs.Address) error {

	var err error

	me.mutex.Lock()
	defer me.mutex.Unlock()

	for range only.Once {
		if _, ok := me.daemons[client]; !ok { // Managed by Mutex
			err = msgs.MakeError(me.EntityId, "service doesn't exist")
			break
		}

		delete(me.daemons, client) // Managed by Mutex
	}

	return err
}

func (me *Daemon) EnsureDaemonNotNil(client msgs.Address) error {

	var err error

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	if _, ok := me.daemons[client]; !ok { // Managed by Mutex
		err = msgs.MakeError(me.EntityId, "service doesn't exist")
	} else {
		err = me.daemons[client].EnsureNotNil() // Managed by Mutex
	}

	return err
}

// Ensure we don't duplicate services.
func (me *Daemon) FindExistingConfig(him ServiceConfig) (*Service, error) {

	var err error
	var sc *Service

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for _, ce := range me.daemons { // Managed by Mutex
		err = ce.IsExisting(him)
		if err != nil {
			sc = ce
			break
		}
	}

	return sc, err
}

// Ensure we don't duplicate services.
func (me *Daemon) IsExisting(s msgs.Address) *Service {

	var sc *Service

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for _, sc = range me.daemons { // Managed by Mutex
		if sc.EntityId == s {
			break
		}
	}

	return sc
}

func (me *Daemon) GetTopics() msgs.SubTopics {

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

func (me *Daemon) GetServiceFiles() map[string]time.Time {

	jc := make(map[string]time.Time)

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for _, ce := range me.daemons { // Managed by Mutex
		jc[ce.JsonFile.Name] = ce.JsonFile.LastModTime
	}

	return jc
}

//
//
//func (me *Daemon) _GetEntityId(u msg.Address) (msg.Address, error) {
//
//	me.mutex.RLock()
//	defer me.mutex.RUnlock()
//
//	err := me.EnsureNotEmpty(u)
//	if err != nil {
//		return "", err
//	}
//
//	return me.daemons[u].EntityId, err	// Managed by Mutex
//}
