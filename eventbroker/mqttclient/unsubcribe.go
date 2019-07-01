package mqttClient

import (
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/states"
	"github.com/gearboxworks/go-status/only"
)

////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

// Unsubscribe a service by method defined by a UUID reference.
func (me *MqttClient) UnsubscribeByUuid(client msgs.Address) error {

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

		me.services[client].State.SetNewAction(states.ActionStop) // Was states.ActionUnsubscribe
		me.services[client].channels.PublishState(me.State)

		me.instance.client.Unsubscribe(me.services[client].Entry.Topic.String())

		me.services[client].State.SetNewState(states.StateStopped, err) // Was states.StateUnsubscribed
		me.services[client].channels.PublishState(me.services[client].State)

		err = me.DeleteEntity(client)
		if err != nil {
			break
		}

		//me.Channels.PublishSpecificState(&client, states.State(states.StateUnsubscribed))
		eblog.Debug(me.EntityId, "unregistered service %s OK", client.String())
	}

	me.Channels.PublishState(me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return err
}

// Unsubscribe a service via a channel defined by a UUID reference.
func (me *MqttClient) UnsubscribeByChannel(caller msgs.Address, u msgs.Address) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		//unreg := me.EntityId.Construct(me.EntityId, states.ActionUnsubscribe, msg.Text(u.String()))
		unreg := caller.MakeMessage(me.EntityId, states.ActionUnsubscribe, msgs.Text(u.String()))
		err = me.Channels.Publish(unreg)
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "unsubscribed service by channel %s OK", u.String())
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return err
}
