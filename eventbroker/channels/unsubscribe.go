package channels

import (
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/states"
	"github.com/gearboxworks/go-status/only"
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

		me.subscribers[client.Address].State.SetNewAction(states.ActionUnsubscribe)

		err = me.subscribers[client.Address].DeleteTopic(client.SubTopic) // Managed by Mutex
		if err != nil {
			break
		}

		me.subscribers[client.Address].State.SetNewState(states.StateUnsubscribed, err)
		eblog.Debug(me.EntityId, "channel unsubscriber: %s", messages.SprintfTopic(client.Address, client.SubTopic))
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}

func (me *Subscriber) UnSubscribe(client messages.SubTopic) error {

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

		me.State.SetNewAction(states.ActionUnsubscribe)

		err = me.DeleteTopic(client) // Managed by Mutex
		if err != nil {
			break
		}

		me.State.SetNewState(states.StateUnsubscribed, err)
		eblog.Debug(me.EntityId, "channel unsubscriber: %s/%s", me.EntityId.String(), client.String())
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}
