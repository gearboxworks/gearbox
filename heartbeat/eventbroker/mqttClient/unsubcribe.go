package mqttClient

import (
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/only"
)


////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

// Unsubscribe a service by method defined by a UUID reference.
func (me *MqttClient) UnsubscribeByUuid(u messages.MessageAddress) error {

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
			me.services[u].State.SetNewAction(states.ActionUnsubscribe)
			channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)

			me.instance.client.Unsubscribe(me.services[u].Entry.Topic.String())

			delete(me.services, u)

			eblog.Debug(me.EntityId, "unregistered service %s OK", u.String())
		}

		me.Channels.PublishCallerState(&u, &states.Status{Current: states.StateUnsubscribed})
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}

// Unsubscribe a service via a channel defined by a UUID reference.
func (me *MqttClient) UnsubscribeByChannel(caller messages.MessageAddress, u messages.MessageAddress) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		//unreg := me.EntityId.Construct(me.EntityId, states.ActionUnsubscribe, messages.MessageText(u.String()))
		unreg := caller.ConstructMessage(me.EntityId, states.ActionUnsubscribe, messages.MessageText(u.String()))
		err = me.Channels.Publish(unreg)
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "unsubscribed service by channel %s OK", u.String())
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


////////////////////////////////////////////////////////////////////////////////
// Executed from a channel.

// Non-exposed channel function that responds to an "unsubscribe" channel request.
func unsubscribeTopic(event *messages.Message, i channels.Argument) channels.Return {

	var me *MqttClient
	var err error

	for range only.Once {
		me, err = InterfaceToTypeMqttClient(i)
		if err != nil {
			break
		}

		//fmt.Printf("MESSAGE Rx:\n[%v]\n", event.Text.String())

		// Use message element as the UUID.
		err = me.UnsubscribeByUuid(event.Text.ToUuid())
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "unsubscribed service by channel %s OK", event.Text.ToUuid())
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}

