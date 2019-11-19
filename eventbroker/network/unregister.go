package network

import (
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/states"
	"github.com/gearboxworks/go-status/only"
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

		me.services[client].State.SetNewAction(states.ActionStop)	// Was states.ActionUnregister
		me.services[client].channels.PublishState(me.State)

		me.services[client].instance.Shutdown()

		me.services[client].State.SetNewState(states.StateStopped, err)	// Was states.StateUnregistered
		me.services[client].channels.PublishState(me.services[client].State)

		err = me.DeleteEntity(client)
		if err != nil {
			break
		}

		//me.Channels.PublishSpecificState(&client, states.State(states.StateUnregistered))
		eblog.Debug(me.EntityId, "unregistered service %s OK", client.String())
	}

	me.Channels.PublishState(me.State)
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

	me.Channels.PublishState(me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}

