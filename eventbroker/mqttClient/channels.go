package mqttClient

import (
	"encoding/json"
	"fmt"
	"gearbox/eventbroker/channels"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/states"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gearboxworks/go-status/only"
)

////////////////////////////////////////////////////////////////////////////////
// Executed from a channel

// Non-exposed channel function that responds to an "stop" channel request.
func stopHandler(event *msgs.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *MqttClient

	for range only.Once {
		me, err = InterfaceToTypeMqttClient(i)
		if err != nil {
			break
		}

		err = me.StopHandler()
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "requested service stop via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return &err
}

// Non-exposed channel function that responds to an "start" channel request.
func startHandler(event *msgs.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *MqttClient

	for range only.Once {
		me, err = InterfaceToTypeMqttClient(i)
		if err != nil {
			break
		}

		err = me.StartHandler()
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "requested service start via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return &err
}

// Non-exposed channel function that responds to a "status" channel request.
func statusHandler(event *msgs.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *MqttClient
	var ret *states.Status

	for range only.Once {
		me, err = InterfaceToTypeMqttClient(i)
		if err != nil {
			break
		}

		if event.Text.String() == "" {
			// Get status of Daemon by default
			ret = me.State.GetStatus()
		} else {
			// Get status of specific sub
			sc := me.IsExisting(msgs.Address(event.Text))
			if sc != nil {
				ret, err = sc.GetStatus()
			}
		}

		eblog.Debug(me.EntityId, "statusHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return ret
}

// Non-exposed channel function that responds to a "register" channel request.
func subscribeTopic(event *msgs.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var me *MqttClient
	var ret *Service
	var err error

	for range only.Once {
		me, err = InterfaceToTypeMqttClient(i)
		if err != nil {
			break
		}

		//fmt.Printf("Rx: %v\n", event)

		var ce ServiceConfig
		err = json.Unmarshal(event.Text.ByteArray(), &ce)

		ret, err = me.Subscribe(ce)
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "subscribed by channel %s OK", ret.EntityId.String())
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return ret
}

// Non-exposed channel function that responds to an "unsubscribe" channel request.
func unsubscribeTopic(event *msgs.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var me *MqttClient
	var err error

	for range only.Once {
		me, err = InterfaceToTypeMqttClient(i)
		if err != nil {
			break
		}

		//fmt.Printf("MESSAGE Rx:\n[%v]\n", event.Text.String())

		// Use message element as the UUID.
		err = me.UnsubscribeByUuid(event.Text.ToAddress())
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "unsubscribed service by channel %s OK", event.Text.ToAddress())
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return &err
}

// Non-exposed channel function that responds to a "get" channel request.
func getHandler(event *msgs.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *MqttClient
	var ret msgs.SubTopics

	for range only.Once {
		me, err = InterfaceToTypeMqttClient(i)
		if err != nil {
			break
		}

		switch event.Text.String() {
		case "topics":
			ret = me.channelHandler.GetTopics()
		case "topics/subs":
			ret = me.channelHandler.GetTopics()
		}

		fmt.Printf("topics: %v\n", ret)

		eblog.Debug(me.EntityId, "topicsHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return &ret
}

func foo2(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
}

func defaultCallback(client mqtt.Client, msg mqtt.Message) {

	fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
}
