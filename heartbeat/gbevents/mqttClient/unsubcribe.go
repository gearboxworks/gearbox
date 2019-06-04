package mqttClient

import (
	"errors"
	"gearbox/app/logger"
	"gearbox/heartbeat/gbevents/channels"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/only"
	"github.com/google/uuid"
)


////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

// Unsubscribe a service by *Service method reference.
func (me *Service) Unsubscribe() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		logger.Debug("MqttClient unsubscribe service (%s).", me.EntityId)
		me.instance = nil
	}

	return err
}

// Unsubscribe a service by method defined by a UUID reference.
func (me *MqttClient) UnsubscribeByUuid(u uuid.UUID) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if _, ok := me.services[u]; !ok {
			err = errors.New("no service defined")
			break
		}

		err = me.services[u].EnsureNotNil()
		if err != nil {
			break
		}

		logger.Debug("MqttClient %s unsubscribe via UUID (%s).", me.EntityId.String(), u.String())
		delete(me.services, u)
	}

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

		//unreg := me.EntityId.Construct(me.EntityId, messages.SubTopicUnsubscribe, messages.MessageText(u.String()))
		unreg := caller.ConstructMessage(me.EntityId, messages.SubTopicSubscribe, messages.MessageText(u.String()))
		err = me.Channels.Publish(unreg)
		if err != nil {
			break
		}

		logger.Debug("MqttClient %s unsubscribe via channel (%s).", me.EntityId.String(), u.String())
	}

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

		// Use message element as the UUID.
		var u uuid.UUID
		//fmt.Printf("MESSAGE Rx:\n[%v]\n", event.Text.String())
		u, err = uuid.Parse(event.Text.String())
		if err != nil {
			logger.Debug("MqttClient %s unsubscribeed service OK", me.EntityId.String())
			break
		}

		err = me.UnsubscribeByUuid(u)
		if err != nil {
			break
		}

		logger.Debug("MqttClient %s unsubscribeed service OK", me.EntityId.String())
		err = nil
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return err
}

