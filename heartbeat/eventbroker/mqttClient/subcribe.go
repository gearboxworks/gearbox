package mqttClient

import (
	"encoding/json"
	"fmt"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/only"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)


////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

// Register a service by method defined by a *CreateTopic structure.
func (me *MqttClient) Subscribe(ce CreateEntry) (*Service, error) {

	var err error
	var sc Service

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = me.services.IsExisting(ce)
		if err != nil {
			break
		}

		for range only.Once {
			sc.State.SetNewWantState(states.StateSubscribed)

			sc.instance = me.instance.client.Subscribe(ce.Topic.String(), ce.Qos, ce.callback)
			if sc.instance == nil {
				err = me.EntityId.ProduceError("unable to subscribe")
				break
			}

			sc.EntityId = messages.GenerateAddress()
			sc.IsManaged = true
			sc.channels = me.channels

			me.services[sc.EntityId] = &sc

			me.services[sc.EntityId].State.SetNewState(states.StateSubscribed)
			eblog.Debug("MqttClient %s registered service %s OK", me.EntityId.String(), sc.EntityId.String())
		}

		// Save last state.
		sc.State.Error = err
		channels.PublishCallerState(me.services[sc.EntityId].channels, &me.services[sc.EntityId].EntityId, &me.services[sc.EntityId].State)
	}
	// Save last state.
	me.State.Error = err
	eblog.LogIfError(&me, err)

	return &sc, err
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
			err = me.EntityId.ProduceError("unable to subscribe")
			break
		}

		reg := caller.ConstructMessage(me.EntityId, states.ActionSubscribe, messages.MessageText(s))
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

		eblog.Debug("MqttClient %s registered service %s via channel", me.EntityId.String(), sc.EntityId.String())
	}
	eblog.LogIfError(&me, err)

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

		eblog.Debug("MqttClient %s registered service %s via channel", me.EntityId.String(), sc.EntityId.String())
	}
	eblog.LogIfError(&me, err)

	return sc
}



func foo2(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
}


func callbackFunc(client mqtt.Client, msg mqtt.Message) {

	fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
}

