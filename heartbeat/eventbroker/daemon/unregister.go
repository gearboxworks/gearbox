package daemon

import (
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/only"
)


////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

// Unregister a service by method defined by a UUID reference.
func (me *Daemon) UnregisterByUuid(u messages.MessageAddress) error {

	var err error
	var state states.Status

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if _, ok := me.daemons[u]; !ok {
			err = me.EntityId.ProduceError("no service defined")
			break
		}

		err = me.daemons[u].EnsureNotNil()
		if err != nil {
			break
		}

		for range only.Once {
			me.daemons[u].State.SetNewAction(states.ActionUnregister)
			channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)

			state, err = me.daemons[u].Status()
			s := state.GetCurrent()
			switch s {
				case states.StateUnknown:
					//

				case states.StateStarted:
					err = me.daemons[u].instance.service.Stop()
					if err != nil {
						break
					}

				case states.StateStopped:
					//
			}

			err = me.daemons[u].instance.service.Uninstall()
			if err != nil {
				break
			}

			delete(me.daemons, u)

			eblog.Debug(me.EntityId, "unregistered service %s OK", u.String())
		}

		me.Channels.PublishSpecificCallerState(&u, states.StateUnregistered)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


// Unregister a service via a channel defined by a UUID reference.
func (me *Daemon) UnregisterByChannel(caller messages.MessageAddress, u messages.MessageAddress) error {

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


////////////////////////////////////////////////////////////////////////////////
// Executed from a channel.

// Non-exposed channel function that responds to an "unregister" channel request.
func unregisterService(event *messages.Message, i channels.Argument) channels.Return {

	var me *Daemon
	var err error

	for range only.Once {
		me, err = InterfaceToTypeDaemon(i)
		if err != nil {
			break
		}

		//fmt.Printf("MESSAGE Rx:\n[%v]\n", event.Text.String())

		// Use message element as the UUID.
		err = me.UnregisterByUuid(event.Text.ToUuid())
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "unregistered service by channel %s OK", event.Text.ToUuid())
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}

