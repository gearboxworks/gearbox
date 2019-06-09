package channels

import (
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/only"
	"sync"
)


func (me *Channels) Subscribe(topic messages.MessageTopic, callback Callback, argInterface Argument) (*Subscriber, error) {

	var err error
	var sub Subscriber

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for range only.Once {
			err = topic.EnsureNotNil()
			if err != nil {
				break
			}

			if callback == nil {
				err = me.EntityId.ProduceError("callback function is empty")
				break
			}

			if me.subscribers == nil {
				me.subscribers = make(Subscribers)
			}

			if _, ok := me.subscribers[topic.Address]; !ok {
				sub = Subscriber{
					EntityId:  topic.Address,
					topics: make(References),
					mutex: sync.RWMutex{},
					//Callbacks: make(Callbacks),
					//Arguments: make(Arguments),
					//Returns:   make(Returns),
					//Executed:   make(Executed),
					//
					//mutexArguments: sync.RWMutex{},	// Mutex control for map.
					//mutexReturns: sync.RWMutex{},	// Mutex control for map.
					//mutexExecuted: sync.RWMutex{},	// Mutex control for map.

					parentInstance: &me.instance,
				}
				*me.subscribers[topic.Address] = sub
			}

			//me.subscribers[topic.Address].Callbacks[topic.SubTopic] = callback
			//// MUTEX me.subscribers[topic.Address].Arguments[topic.SubTopic] = argInterface
			//me.subscribers[topic.Address].SetArguments(topic.SubTopic, argInterface)
			//// MUTEX me.subscribers[topic.Address].Returns[topic.SubTopic] = nil
			//me.subscribers[topic.Address].SetReturns(topic.SubTopic, nil)

			me.subscribers[topic.Address].AddTopic(topic.SubTopic, callback, argInterface)

			// me.subscribers[topic.Address].List()
			// me.subscribers.List()

			me.subscribers[topic.Address].State.SetNewState(states.StateSubscribed, err)
			eblog.Debug(me.EntityId, "channel new subscriber: %s", messages.SprintfTopic(topic.Address, topic.SubTopic))
		}
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return &sub, err
}


func (me *Subscriber) Subscribe(topic messages.SubTopic, callback Callback, argInterface Argument) error {//, ret Return) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for range only.Once {
			err = topic.EnsureNotNil()
			if err != nil {
				break
			}

			if callback == nil {
				err = me.EntityId.ProduceError("callback function is empty")
				break
			}

			//me.Callbacks[topic] = callback
			//// MUTEX me.Arguments[topic] = argInterface
			//me.SetArguments(topic, argInterface)
			//// MUTEX me.Returns[topic] = nil
			//me.SetReturns(topic, nil)

			me.AddTopic(topic, callback, argInterface)

			// me.List()

			me.State.SetNewState(states.StateSubscribed, err)
			eblog.Debug(me.EntityId, "channel new subscriber: %s", messages.SprintfTopic(me.EntityId, topic))
		}
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}
