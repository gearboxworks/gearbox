package channels

import (
	"errors"
	"gearbox/app/logger"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/only"
)


func (me *Channels) Subscribe(topic messages.Topic, callback Callback, argInterface Argument) (Subscriber, error) {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = topic.EnsureNotNil()
		if err != nil {
			break
		}

		if callback == nil {
			err = errors.New("channel callback function is empty")
			break
		}

		if me.subscribers == nil {
			me.subscribers = make(Subscribers)
		}

		if _, ok := me.subscribers[topic.Address]; !ok {
			*me.subscribers[topic.Address] = Subscriber{
				Address: topic.Address,
				Callbacks: make(Callbacks),
				Arguments: make(Arguments),
				Returns: make(Returns),
				//References: make(References),
				instance: &me.instance,
			}
		}

		me.subscribers[topic.Address].Callbacks[topic.SubTopic] = callback
		me.subscribers[topic.Address].Arguments[topic.SubTopic] = argInterface
		me.subscribers[topic.Address].Returns[topic.SubTopic] = nil

		// me.subscribers[topic.Address].List()
		// me.subscribers.List()

		logger.Debug("New subscriber: %s", messages.SprintfTopic(topic.Address, topic.SubTopic))
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}
	// Save last state.
	me.Error = err

	if _, ok := me.subscribers[topic.Address]; ok {
		return *me.subscribers[topic.Address], err
	} else {
		return Subscriber{}, err
	}
}


func (me *Subscriber) Subscribe(topic messages.SubTopic, callback Callback, argInterface Argument) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = topic.EnsureNotNil()
		if err != nil {
			break
		}

		if callback == nil {
			err = errors.New("channel callback function is empty")
			break
		}

		me.Callbacks[topic] = callback
		me.Arguments[topic] = argInterface
		me.Returns[topic] = nil
		// me.List()

		logger.Debug("New subscriber: %s", messages.SprintfTopic(me.Address, topic))
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return err
}
