package daemon

import (
	"fmt"
	"gearbox/heartbeat/eventbroker/messages"
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


func (me *Daemon) DeleteEntity(entity messages.MessageAddress) {

	me.mutex.Lock()
	defer me.mutex.Unlock()
	delete(me.daemons, entity)	// Managed by Mutex

	return
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
func (me *Daemon) IsExisting(him ServiceConfig) (*Service, error) {

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


func (me *Daemon) GetServiceFiles() map[string]time.Time {

	jc := make(map[string]time.Time)

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for _, ce := range me.daemons {	// Managed by Mutex
		jc[ce.JsonFile.Name] = ce.JsonFile.LastModTime
	}

	return jc
}


func (me *Daemon) Foo() {

	var st messages.SubTopics
	var t messages.Topics
	var ma messages.MessageAddresses

	fmt.Printf("\nme.GetTopics\n")
	st = me.GetTopics()
	for _, f := range st {
		fmt.Printf("me.GetTopics => %s\n", f.String())
	}

	fmt.Printf("\nme.GetManagedEntities\n")
	ma = me.GetManagedEntities()
	for _, f := range ma {
		fmt.Printf("me.GetManagedEntities => %s\n", f.String())
	}

	fmt.Printf("\nme.GetEntities\n")
	ma = me.GetEntities()
	for _, f := range ma {
		fmt.Printf("me.GetEntities => %s\n", f.String())
	}

	fmt.Printf("\nme.Channels.GetManagedEntities\n")
	ma = me.Channels.GetManagedEntities()
	for _, f := range ma {
		fmt.Printf("me.Channels.GetManagedEntities => %s\n", f.String())
	}

	fmt.Printf("\nme.Channels.GetEntities\n")
	ma = me.Channels.GetEntities()
	for _, f := range ma {
		fmt.Printf("me.Channels.GetEntities => %s\n", f.String())
	}

	fmt.Printf("\nme.Channels.GetListenerTopics\n")
	t = me.Channels.GetListenerTopics()
	for _, f := range t {
		fmt.Printf("me.Channels.GetListenerTopics => %s\n", f.String())
	}

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

