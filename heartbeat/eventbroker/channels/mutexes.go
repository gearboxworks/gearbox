package channels

import (
	"gearbox/heartbeat/eventbroker/messages"
)


//func (me *Subscriber) GetEntityId(u messages.MessageAddress) messages.MessageAddress {
//
//	me.daemonsMutex.RLock()
//	defer me.daemonsMutex.RUnlock()
//	return me.daemons[u].EntityId	// Managed by Mutex
//}
//
//
//func (me *Subscriber) EnsureSubscriberNotNil(u messages.MessageAddress) error {
//
//	var err error
//
//	me.daemonsMutex.RLock()
//	defer me.daemonsMutex.RUnlock()
//
//	_, ok := me.daemons[u]	// Managed by Mutex
//	if !ok {
//		err = me.EntityId.ProduceError("service doesn't exist")
//	} else {
//		err = me.daemons[u].EnsureNotNil()	// Managed by Mutex
//	}
//
//	return err
//}
//
//
//// Mutex handling.
//func (me *Subscriber) DeleteSubscriber(u messages.MessageAddress) {
//
//	me.daemonsMutex.Lock()
//	defer me.daemonsMutex.Unlock()
//	delete(me.daemons, u)	// Managed by Mutex
//
//	return
//}



// Mutex handling.
func (me *Subscriber) GetExecuted(sub messages.SubTopic) bool {
	me.mutex.RLock()
	defer me.mutex.RUnlock()

	return me.topics[sub].Executed	// Managed by Mutex
}

func (me *Subscriber) SetExecuted(sub messages.SubTopic, v bool) {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	me.topics[sub].Executed = v		// Managed by Mutex

	return
}


// Mutex handling.
func (me *Subscriber) AddTopic(topic messages.SubTopic, callback Callback, argInterface Argument) {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	if _, ok := me.topics[topic]; !ok {
		var ret *Return
		me.topics[topic] = &Reference{
			Callback: callback,		// Managed by Mutex
			Argument: argInterface,	// Managed by Mutex
			Return:   ret,			// Managed by Mutex
			Executed: false,		// Managed by Mutex
		}
	}

	return
}

// Mutex handling.
func (me *Subscriber) DeleteTopic(sub messages.SubTopic) {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	delete(me.topics, sub)

	return
}


// Mutex handling.
func (me *Subscriber) GetTopic(topic messages.SubTopic) (error, Callback, Argument, Return) {
	me.mutex.RLock()
	defer me.mutex.RUnlock()

	if _, ok := me.topics[topic]; !ok {	// Managed by Mutex
		return me.EntityId.ProduceError("channel topic doesn't exist"), nil, nil, nil
	}

	if me.topics[topic].Return == nil {	// Managed by Mutex
		return me.EntityId.ProduceError("channel return not defined"), nil, nil, nil
	}

	return nil, me.topics[topic].Callback, me.topics[topic].Argument, me.topics[topic].Return	// Managed by Mutex
}


// Mutex handling.
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

