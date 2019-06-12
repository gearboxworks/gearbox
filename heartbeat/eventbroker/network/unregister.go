package network

import (
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/only"
	"gearbox/heartbeat/eventbroker/states"
)


////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

// Unregister a service by method defined by a UUID reference.
func (me *ZeroConf) UnregisterByEntityId(client messages.MessageAddress) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = me.services[client].EnsureNotNil()
		if err != nil {
			break
		}

		me.services[client].State.SetNewAction(states.ActionUnregister)
		me.services[client].channels.PublishCallerState(&me.State)

		me.services[client].instance.Shutdown()

		err = me.DeleteEntity(client)
		if err != nil {
			break
		}

		me.Channels.PublishSpecificState(&client, states.State(states.StateUnregistered))
		eblog.Debug(me.EntityId, "unregistered service %s OK", client.String())
	}

	me.Channels.PublishState(&me.EntityId, &me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}

// Unregister a service via a channel defined by a UUID reference.
func (me *ZeroConf) UnregisterByChannel(caller messages.MessageAddress, u messages.MessageAddress) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		//unreg := me.EntityId.Construct(me.EntityId, states.ActionUnregister, messages.MessageText(u.String()))
		unreg := caller.ConstructMessage(me.EntityId, states.ActionUnregister, messages.MessageText(u.String()))
		err = me.Channels.Publish(unreg)
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "unregistered service by channel %s OK", u.String())
	}

	me.Channels.PublishState(&me.EntityId, &me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}

