package channels

import (
	"fmt"
	"gearbox/eventbroker/msgs"
	"github.com/gearboxworks/go-status/only"
)

// Mutex handling.

func (me *Channels) GetEntities() msgs.Addresses {

	var ret msgs.Addresses

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for s := range me.subscribers { // Managed by Mutex
		ret = append(ret, s)
	}

	return ret
}

func (me *Channels) GetManagedEntities() msgs.Addresses {

	var ret msgs.Addresses

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

func (me *Channels) AddEntity(client msgs.Address, sc *Subscriber) error {
	var err error

	me.mutex.Lock()
	defer me.mutex.Unlock()

	if _, ok := me.subscribers[client]; !ok { // Managed by Mutex
		me.subscribers[client] = sc
	} else {
		err = msgs.MakeError(me.EntityId, "service %s already exists", client)
	}

	return err
}

func (me *Channels) DeleteEntity(client msgs.Address) error {

	var err error

	me.mutex.Lock()
	defer me.mutex.Unlock()

	for range only.Once {
		if _, ok := me.subscribers[client]; !ok { // Managed by Mutex
			err = msgs.MakeError(me.EntityId, "service doesn't exist")
			break
		}

		delete(me.subscribers, client) // Managed by Mutex
	}

	return err
}

func (me *Subscriber) GetTopic(topic msgs.SubTopic) (error, Callback, Argument, Return, ReturnType) {

	var err error
	var cb Callback
	var args Argument
	var ret Return
	var retType ReturnType

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for range only.Once {
		if _, ok := me.topics[topic]; !ok { // Managed by Mutex
			err = msgs.MakeError(me.EntityId, "channel topic doesn't exist")
			break
		}

		if me.topics[topic].Return == nil { // Managed by Mutex
			err = msgs.MakeError(me.EntityId, "channel return not defined")
			break
		}

		if me.topics[topic].ReturnType == "" { // Managed by Mutex
			err = msgs.MakeError(me.EntityId, "channel return type not defined")
			break
		}

		cb = me.topics[topic].Callback        // Managed by Mutex
		args = me.topics[topic].Argument      // Managed by Mutex
		ret = me.topics[topic].Return         // Managed by Mutex
		retType = me.topics[topic].ReturnType // Managed by Mutex
	}

	return err, cb, args, ret, retType
}

func (me *Subscriber) GetTopics() msgs.SubTopics {

	var ret msgs.SubTopics

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

func (me *Channels) GetListeners(topic msgs.Topic) ([]string, error) {

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

func (me *Channels) GetListenerTopics() (msgs.Topics, error) {

	var topics msgs.Topics
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for _, t := range me.instance.emitter.Topics() {
			topics = append(topics, msgs.StringToTopic(t))
		}
	}

	return topics, err
}

func (me *Subscriber) GetExecuted(sub msgs.SubTopic) bool {

	if me == nil {
		return false
	}

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	return me.topics[sub].Executed // Managed by Mutex
}

func (me *Subscriber) SetExecuted(sub msgs.SubTopic, v bool) {

	if me == nil {
		return
	}

	me.mutex.Lock()
	defer me.mutex.Unlock()

	me.topics[sub].Executed = v // Managed by Mutex

	return
}

func (me *Subscriber) AddTopic(topic msgs.SubTopic, callback Callback, argInterface Argument, retType ReturnType) {
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

func (me *Subscriber) DeleteTopic(client msgs.SubTopic) error {

	var err error

	me.mutex.Lock()
	defer me.mutex.Unlock()

	for range only.Once {
		_, ok := me.topics[client] // Managed by Mutex
		if !ok {
			err = msgs.MakeError(me.EntityId, "service doesn't exist")
			break
		}

		delete(me.topics, client) // Managed by Mutex
	}

	return err
}

func (me *Subscriber) GetReturns(sub msgs.SubTopic) Return {

	me.mutex.RLock()
	r := me.topics[sub].Return // Managed by Mutex
	me.mutex.RUnlock()

	return r
}

func (me *Subscriber) SetReturns(sub msgs.SubTopic, v Return) {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	me.topics[sub].Return = v // Managed by Mutex

	return
}

//func (me *Subscriber) GetEntityId(client msg.Address) msg.Address {
//
//	me.daemonsMutex.RLock()
//	defer me.daemonsMutex.RUnlock()
//	return me.daemons[client].EntityId	// Managed by Mutex
//}
//
//
//func (me *Subscriber) EnsureSubscriberNotNil(client msg.Address) error {
//
//	var err error
//
//	me.daemonsMutex.RLock()
//	defer me.daemonsMutex.RUnlock()
//
//	_, ok := me.daemons[client]	// Managed by Mutex
//	if !ok {
//		err = msgs.MakeError(me.EntityId,"service doesn't exist")
//	} else {
//		err = me.daemons[client].EnsureNotEmpty()	// Managed by Mutex
//	}
//
//	return err
//}
//
//
//// Mutex handling.
//func (me *Subscriber) DeleteSubscriber(client msg.Address) {
//
//	me.daemonsMutex.Lock()
//	defer me.daemonsMutex.Unlock()
//	delete(me.daemons, client)	// Managed by Mutex
//
//	return
//}
