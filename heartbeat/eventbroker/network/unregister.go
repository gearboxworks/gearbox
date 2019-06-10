package network

import (
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/only"
)


////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

// Unregister a service by method defined by a UUID reference.
func (me *ZeroConf) UnregisterByEntityId(u messages.MessageAddress) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if _, ok := me.services[u]; !ok {
			err = me.EntityId.ProduceError("no service defined")
			break
		}

		err = me.services[u].EnsureNotNil()
		if err != nil {
			break
		}

		for range only.Once {
			me.services[u].State.SetNewAction(states.ActionUnregister)
			channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)

			me.services[u].instance.Shutdown()
			delete(me.services, u)

			eblog.Debug(me.EntityId, "unregistered service %s OK", u.String())
		}

		me.Channels.PublishCallerState(&u, &states.Status{Current: states.StateUnregistered})
	}

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

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}

