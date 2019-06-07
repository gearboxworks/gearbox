package daemon

import (
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/only"
	"github.com/google/uuid"
)


////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

// Unregister a service by method defined by a UUID reference.
func (me *Daemon) Unregister(u messages.MessageAddress) error {

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

		state, err = me.daemons[u].Status()
		switch state.Current {
			case states.StateUnknown:
		}

		if state.Current == states.StateUnknown {
			// Reset err, because it's not an error.
			err = nil
			break

		} else if state.Current == states.StateStarted {
			err = me.daemons[u].instance.service.Stop()
			if err != nil {
				break
			}

		} else if state.Current == states.StateStopped {
			//
		}

		err = me.daemons[u].instance.service.Uninstall()
		if err != nil {
			break
		}

		delete(me.daemons, u)

		eblog.Debug("Daemon %s unregister via UUID (%s).", me.EntityId.String(), u.String())
	}
	eblog.LogIfError(&me, err)

	return err
}

// Unregister a service via a channel defined by a UUID reference.
func (me *Daemon) UnregisterByChannel(caller messages.MessageAddress, u uuid.UUID) error {

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

		eblog.Debug("Daemon %s unregister via channel (%s).", me.EntityId.String(), u.String())
	}
	eblog.LogIfError(&me, err)

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
		err = me.Unregister(event.Text.ToUuid())
		if err != nil {
			break
		}

		eblog.Debug("Daemon %s unregistered service OK", me.EntityId.String())
	}
	eblog.LogIfError(&me, err)

	return err
}

