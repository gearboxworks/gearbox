package channels

import (
	"gearbox/app/logger"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/only"
)

func (me *Channels) UnSubscribe(topic messages.Topic) error {

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

		delete(me.subscribers[topic.Address].Callbacks, topic.SubTopic)
		delete(me.subscribers[topic.Address].Arguments, topic.SubTopic)
		delete(me.subscribers[topic.Address].Returns, topic.SubTopic)
		logger.Debug("Unsubscribed: %s", messages.SprintfTopic(topic.Address, topic.SubTopic))
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	// Save last state.
	me.Error = err
	return err
}


func (me *Subscriber) UnSubscribe(subtopic messages.SubTopic) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = subtopic.EnsureNotNil()
		if err != nil {
			break
		}

		delete(me.Callbacks, subtopic)
		delete(me.Arguments, subtopic)
		delete(me.Returns, subtopic)
		logger.Debug("Unsubscribed: %s", messages.SprintfTopic(me.Address, subtopic))
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return err
}
