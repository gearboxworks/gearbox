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

		// MUTEX
		me.subscribers[topic.Address].DeleteSubTopic(topic.SubTopic)

		eblog.Debug(me.EntityId, "Unsubscribed: %s", messages.SprintfTopic(topic.Address, topic.SubTopic))
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

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

		// MUTEX
		me.DeleteSubTopic(subtopic)

		eblog.Debug(me.EntityId, "Unsubscribed: %s", messages.SprintfTopic(me.EntityId, subtopic))
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}
