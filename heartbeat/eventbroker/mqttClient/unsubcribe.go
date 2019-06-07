package mqttClient

import (
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/only"
	"github.com/google/uuid"
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

		me.instance.client.Unsubscribe(me.services[u].Entry.Topic.String())

		eblog.Debug("MqttClient %s unsubscribe via UUID (%s).", me.EntityId.String(), u.String())
		delete(me.services, u)
	}
	eblog.LogIfError(&me, err)

	return err
}

// Unsubscribe a service via a channel defined by a UUID reference.
func (me *MqttClient) UnsubscribeByChannel(caller messages.MessageAddress, u uuid.UUID) error {

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

		eblog.Debug("MqttClient %s unsubscribe via channel (%s).", me.EntityId.String(), u.String())
	}
	eblog.LogIfError(&me, err)

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

		eblog.Debug("MqttClient %s unsubscribeed service OK", me.EntityId.String())
	}
	eblog.LogIfError(&me, err)

	return err
}

