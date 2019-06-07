package channels

import (
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/only"
)


func (me *Channels) Subscribe(topic messages.MessageTopic, callback Callback, argInterface Argument) (*Subscriber, error) {

	var err error
	var sub Subscriber

	for range only.Once {
		err = EnsureNotNil(me)
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
					Callbacks: make(Callbacks),
					Arguments: make(Arguments),
					Returns:   make(Returns),
					//References: make(References),
					parentInstance: &me.instance,
				}
				*me.subscribers[topic.Address] = sub
			}

			me.subscribers[topic.Address].Callbacks[topic.SubTopic] = callback
			me.subscribers[topic.Address].Arguments[topic.SubTopic] = argInterface
			me.subscribers[topic.Address].Returns[topic.SubTopic] = nil

			// me.subscribers[topic.Address].List()
			// me.subscribers.List()

			// Save last state.
			me.subscribers[topic.Address].State.Error = err
			eblog.Debug("channel new subscriber: %s", messages.SprintfTopic(topic.Address, topic.SubTopic))
		}
	}
	eblog.LogIfError(&me, err)

	return &sub, err
}


func (me *Subscriber) Subscribe(topic messages.SubTopic, callback Callback, argInterface Argument) error {

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

			me.Callbacks[topic] = callback
			me.Arguments[topic] = argInterface
			me.Returns[topic] = nil
			// me.List()

			// Save last state.
			me.State.Error = err
			eblog.Debug("channel new subscriber: %s", messages.SprintfTopic(me.EntityId, topic))
		}
	}
	eblog.LogIfError(&me, err)

	return err
}
