package mqttClient

import (
	"encoding/json"
	"errors"
	"fmt"
	"gearbox/app/logger"
	"gearbox/heartbeat/gbevents/channels"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/only"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)


////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

// Register a service by method defined by a *CreateTopic structure.
func (me *MqttClient) Subscribe(msg Topic) (*Service, error) {

	var err error
	var u uuid.UUID
	var sc Service

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		cb := msgCallback{Topic: msg, Function: foo2}
		sc.instance = me.instance.client.Subscribe(cb.Topic.String(), 0, cb.Function)
		if sc.instance == nil {
			err = errors.New("mqttClient unable to subscribe")
			break
		}

		u, err = uuid.NewUUID()
		if err != nil {
			return &sc, err
		}
		fmt.Printf("HARRY: %s\n", u.String())

		sc.EntityId = u
		sc.IsManaged = true

		me.services[u] = &sc

		logger.Debug("MqttClient %s registered service %s OK", me.EntityId.String(), u.String())
	}

	return &sc, err
}


func foo2(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
}


func callbackFunc(client mqtt.Client, msg mqtt.Message) {

	fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
}


// Subscribe a service via a channel defined by a *CreateTopic structure and
// returns a *Service structure if successful.
func (me *MqttClient) SubscribeByChannel(caller messages.MessageAddress, s Topic) (*Service, error) {

	var err error
	var sc *Service

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if s == "" {
			err = errors.New("mqttClient unable to subscribe")
			break
		}

		reg := caller.ConstructMessage(me.EntityId, messages.SubTopicSubscribe, messages.MessageText(s))
		err = me.Channels.Publish(reg)
		if err != nil {
			break
		}

		rs, err := me.Channels.GetCallbackReturn(reg, 100)
		if err != nil {
			break
		}

		sc, err = InterfaceToTypeService(rs)	// sc = rs.(*Service)
		if err != nil {
			break
		}

		logger.Debug("MqttClient %s registered service %s via channel", me.EntityId.String(), sc.EntityId.String())
	}

	return sc, err
}


////////////////////////////////////////////////////////////////////////////////
// Executed from a channel.

// Non-exposed channel function that responds to a "register" channel request.
func subscribeTopic(event *messages.Message, i channels.Argument) channels.Return {

	var me *MqttClient
	var sc *Service
	var err error

	for range only.Once {
		me, err = InterfaceToTypeMqttClient(i)
		if err != nil {
			break
		}

		//fmt.Printf("Rx: %v\n", event)

		ce := Topic(event.Text)
		err = json.Unmarshal(event.Text.ByteArray(), &ce)

		sc, err = me.Subscribe(ce)
		if err != nil {
			break
		}
		// logger.Debug("Service: %v", sc)

		err = nil
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return sc
}

