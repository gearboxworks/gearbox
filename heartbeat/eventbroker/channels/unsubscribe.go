package channels

import (
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/only"
)

func (me *Channels) UnSubscribe(client messages.MessageTopic) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = client.EnsureNotNil()
		if err != nil {
			break
		}

		err = me.EnsureSubscriberNotNil(client.Address)
		if err != nil {
			break
		}

		err, _, _, _, _ = me.subscribers[client.Address].GetTopic(client.SubTopic)
		if err != nil {
			break
		}

		me.subscribers[client.Address].DeleteTopic(client.SubTopic)	// Managed by Mutex

		eblog.Debug(me.EntityId, "Unsubscribed: %s", messages.SprintfTopic(client.Address, client.SubTopic))
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Subscriber) UnSubscribe(topic messages.SubTopic) error {

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

		me.DeleteTopic(topic)	// Managed by Mutex

		eblog.Debug(me.EntityId, "Unsubscribed: %s", messages.SprintfTopic(me.EntityId, topic))
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}
