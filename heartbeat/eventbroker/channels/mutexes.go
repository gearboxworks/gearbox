package channels

import (
	"fmt"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/only"
)

// Mutex handling.

func (me *Channels) GetEntities() messages.MessageAddresses {

	var ret messages.MessageAddresses

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for s, _ := range me.subscribers {	// Managed by Mutex
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

	for s, _ := range me.subscribers {	// Managed by Mutex
		if me.subscribers[s].IsManaged {	// Managed by Mutex
			ret = append(ret, s)
		}
	}

	return ret
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

		cb = me.topics[topic].Callback			// Managed by Mutex
		args = me.topics[topic].Argument		// Managed by Mutex
		ret = me.topics[topic].Return			// Managed by Mutex
		retType = me.topics[topic].ReturnType	// Managed by Mutex
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

	for t, _ := range me.topics {
		ret = append(ret, t)
	}

	return ret
}


func (me *Channels) GetListeners(topic messages.MessageTopic) []string {

	var ret []string

	if me == nil {
		return ret
	}

	foo := me.instance.emitter.Listeners(topic.String())[0]

	for f := range foo {
		ret = append(ret, f.Topic)
		fmt.Printf("[%s] - '%s' '%s' '%s'\n", f.Topic, f.OriginalTopic, f.Args, f.Flags)
	}

	return ret
}


func (me *Channels) GetListenerTopics() messages.Topics {

	var topics messages.Topics

	if me == nil {
		return topics
	}

	for _, t := range me.instance.emitter.Topics() {
		topics = append(topics, messages.StringToTopic(t))
	}

	return topics
}


func (me *Subscriber) GetExecuted(sub messages.SubTopic) bool {

	if me == nil {
		return false
	}

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	return me.topics[sub].Executed	// Managed by Mutex
}

func (me *Subscriber) SetExecuted(sub messages.SubTopic, v bool) {

	if me == nil {
		return
	}

	me.mutex.Lock()
	defer me.mutex.Unlock()

	me.topics[sub].Executed = v		// Managed by Mutex

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

func (me *Subscriber) DeleteTopic(sub messages.SubTopic) {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	delete(me.topics, sub)

	return
}


func (me *Subscriber) GetReturns(sub messages.SubTopic) Return {

	me.mutex.RLock()
	r := me.topics[sub].Return	// Managed by Mutex
	me.mutex.RUnlock()

	return r
}

func (me *Subscriber) SetReturns(sub messages.SubTopic, v Return) {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	me.topics[sub].Return = v	// Managed by Mutex

	return
}


//func (me *Subscriber) GetEntityId(entity messages.MessageAddress) messages.MessageAddress {
//
//	me.daemonsMutex.RLock()
//	defer me.daemonsMutex.RUnlock()
//	return me.daemons[entity].EntityId	// Managed by Mutex
//}
//
//
//func (me *Subscriber) EnsureSubscriberNotNil(entity messages.MessageAddress) error {
//
//	var err error
//
//	me.daemonsMutex.RLock()
//	defer me.daemonsMutex.RUnlock()
//
//	_, ok := me.daemons[entity]	// Managed by Mutex
//	if !ok {
//		err = me.EntityId.ProduceError("service doesn't exist")
//	} else {
//		err = me.daemons[entity].EnsureNotNil()	// Managed by Mutex
//	}
//
//	return err
//}
//
//
//// Mutex handling.
//func (me *Subscriber) DeleteSubscriber(entity messages.MessageAddress) {
//
//	me.daemonsMutex.Lock()
//	defer me.daemonsMutex.Unlock()
//	delete(me.daemons, entity)	// Managed by Mutex
//
//	return
//}

