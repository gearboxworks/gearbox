package channels

import (
	"fmt"
	"gearbox/eventbroker/messages"
	"github.com/gearboxworks/go-status/only"
)

// Mutex handling.

func (me *Channels) GetEntities() messages.MessageAddresses {

	var ret messages.MessageAddresses

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for s := range me.subscribers { // Managed by Mutex
		ret = append(ret, s)
	}

	return ret
}

func (me *Channels) GetManagedEntities() messages.MessageAddresses {

	var ret messages.MessageAddresses

	if me == nil {
		return ret
	}

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for s := range me.subscribers { // Managed by Mutex
		if me.subscribers[s].IsManaged { // Managed by Mutex
			ret = append(ret, s)
		}
	}

	return ret
}

func (me *Channels) AddEntity(client messages.MessageAddress, sc *Subscriber) error {
	var err error

	me.mutex.Lock()
	defer me.mutex.Unlock()

	if _, ok := me.subscribers[client]; !ok { // Managed by Mutex
		me.subscribers[client] = sc
	} else {
		err = me.EntityId.ProduceError("service %s already exists", client)
	}

	return err
}

func (me *Channels) DeleteEntity(client messages.MessageAddress) error {

	var err error

	me.mutex.Lock()
	defer me.mutex.Unlock()

	for range only.Once {
		if _, ok := me.subscribers[client]; !ok { // Managed by Mutex
			err = me.EntityId.ProduceError("service doesn't exist")
			break
		}

		delete(me.subscribers, client) // Managed by Mutex
	}

	return err
}

func (me *Subscriber) GetTopic(topic messages.SubTopic) (error, Callback, Argument, Return, ReturnType) {

	var err error
	var cb Callback
	var args Argument
	var ret Return
	var retType ReturnType

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for range only.Once {
		if _, ok := me.topics[topic]; !ok { // Managed by Mutex
			err = me.EntityId.ProduceError("channel topic doesn't exist")
			break
		}

		if me.topics[topic].Return == nil { // Managed by Mutex
			err = me.EntityId.ProduceError("channel return not defined")
			break
		}

		if me.topics[topic].ReturnType == "" { // Managed by Mutex
			err = me.EntityId.ProduceError("channel return type not defined")
			break
		}

		cb = me.topics[topic].Callback        // Managed by Mutex
		args = me.topics[topic].Argument      // Managed by Mutex
		ret = me.topics[topic].Return         // Managed by Mutex
		retType = me.topics[topic].ReturnType // Managed by Mutex
	}

	return err, cb, args, ret, retType
}

func (me *Subscriber) GetTopics() messages.SubTopics {

	var ret messages.SubTopics

	if me == nil {
		return ret
	}

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for t := range me.topics {
		ret = append(ret, t)
	}

	return ret
}

func (me *Channels) GetListeners(topic messages.MessageTopic) ([]string, error) {

	var ret []string
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		foo := me.instance.emitter.Listeners(topic.String())[0]

		for f := range foo {
			ret = append(ret, f.Topic)
			fmt.Printf("[%s] - '%s' '%s' '%s'\n", f.Topic, f.OriginalTopic, f.Args, f.Flags)
		}
	}

	return ret, err
}

func (me *Channels) GetListenerTopics() (messages.Topics, error) {

	var topics messages.Topics
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for _, t := range me.instance.emitter.Topics() {
			topics = append(topics, messages.StringToTopic(t))
		}
	}

	return topics, err
}

func (me *Subscriber) GetExecuted(sub messages.SubTopic) bool {

	if me == nil {
		return false
	}

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	return me.topics[sub].Executed // Managed by Mutex
}

func (me *Subscriber) SetExecuted(sub messages.SubTopic, v bool) {

	if me == nil {
		return
	}

	me.mutex.Lock()
	defer me.mutex.Unlock()

	me.topics[sub].Executed = v // Managed by Mutex

	return
}

func (me *Subscriber) AddTopic(topic messages.SubTopic, callback Callback, argInterface Argument, retType ReturnType) {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	if _, ok := me.topics[topic]; !ok {
		var ret *Return
		me.topics[topic] = &Reference{
			Callback:   callback,     // Managed by Mutex
			Argument:   argInterface, // Managed by Mutex
			Return:     ret,          // Managed by Mutex
			ReturnType: retType,      // Managed by Mutex
			Executed:   false,        // Managed by Mutex
		}
	}

	return
}

func (me *Subscriber) DeleteTopic(client messages.SubTopic) error {

	var err error

	me.mutex.Lock()
	defer me.mutex.Unlock()

	for range only.Once {
		_, ok := me.topics[client] // Managed by Mutex
		if !ok {
			err = me.EntityId.ProduceError("service doesn't exist")
			break
		}

		delete(me.topics, client) // Managed by Mutex
	}

	return err
}

func (me *Subscriber) GetReturns(sub messages.SubTopic) Return {

	me.mutex.RLock()
	r := me.topics[sub].Return // Managed by Mutex
	me.mutex.RUnlock()

	return r
}

func (me *Subscriber) SetReturns(sub messages.SubTopic, v Return) {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	me.topics[sub].Return = v // Managed by Mutex

	return
}

//func (me *Subscriber) GetEntityId(client messages.MessageAddress) messages.MessageAddress {
//
//	me.daemonsMutex.RLock()
//	defer me.daemonsMutex.RUnlock()
//	return me.daemons[client].EntityId	// Managed by Mutex
//}
//
//
//func (me *Subscriber) EnsureSubscriberNotNil(client messages.MessageAddress) error {
//
//	var err error
//
//	me.daemonsMutex.RLock()
//	defer me.daemonsMutex.RUnlock()
//
//	_, ok := me.daemons[client]	// Managed by Mutex
//	if !ok {
//		err = me.EntityId.ProduceError("service doesn't exist")
//	} else {
//		err = me.daemons[client].EnsureNotNil()	// Managed by Mutex
//	}
//
//	return err
//}
//
//
//// Mutex handling.
//func (me *Subscriber) DeleteSubscriber(client messages.MessageAddress) {
//
//	me.daemonsMutex.Lock()
//	defer me.daemonsMutex.Unlock()
//	delete(me.daemons, client)	// Managed by Mutex
//
//	return
//}
