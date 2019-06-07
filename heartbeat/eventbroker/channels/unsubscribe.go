package channels

import (
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/only"
)

func (me *Channels) UnSubscribe(topic messages.MessageTopic) error {

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

		eblog.Debug("Unsubscribed: %s", messages.SprintfTopic(topic.Address, topic.SubTopic))
	}
	// Save last state.
	me.State.Error = err
	eblog.LogIfError(&me, err)

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

		eblog.Debug("Unsubscribed: %s", messages.SprintfTopic(me.EntityId, subtopic))
	}
	// Save last state.
	me.State.Error = err
	eblog.LogIfError(&me, err)

	return err
}
