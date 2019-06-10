package channels

import (
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/only"
	"sync"
)


func (me *Channels) Subscribe(topic messages.MessageTopic, callback Callback, argInterface Argument, retType ReturnType) (*Subscriber, error) {

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

			if retType == "" {
				err = me.EntityId.ProduceError("return type is empty")
				break
			}

			if me.subscribers == nil {
				me.subscribers = make(Subscribers)
			}

			if _, ok := me.subscribers[topic.Address]; !ok {
				sub = Subscriber{
					EntityId:  topic.Address,
					State: states.Status{},
					IsManaged: true,

					topics: make(References),
					mutex: sync.RWMutex{},
					parentInstance: &me.instance,
				}
				*me.subscribers[topic.Address] = sub
			}
			me.subscribers[topic.Address].AddTopic(topic.SubTopic, callback, argInterface, retType)

			me.subscribers[topic.Address].State.SetNewState(states.StateSubscribed, err)
			eblog.Debug(me.EntityId, "channel new subscriber: %s", messages.SprintfTopic(topic.Address, topic.SubTopic))
		}
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return &sub, err
}


func (me *Subscriber) Subscribe(topic messages.SubTopic, callback Callback, argInterface Argument, retType ReturnType) error {

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

			if retType == "" {
				err = me.EntityId.ProduceError("return type is empty")
				break
			}

			me.AddTopic(topic, callback, argInterface, retType)

			me.State.SetNewState(states.StateSubscribed, err)
			eblog.Debug(me.EntityId, "channel new subscriber: %s", messages.SprintfTopic(me.EntityId, topic))
		}
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}

