package mqttClient

import (
	"fmt"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/only"
	"gearbox/heartbeat/eventbroker/states"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fhmq/hmq/lib/topics"
	"net/url"
)


func (me *MqttClient) IsConnected() (bool, error) {

	var err error
	var ok bool

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if me.instance.client != nil {
			if me.instance.client.IsConnected() {
				ok = true
			}
		}
	}

	return ok, err
}


func (me *MqttClient) ConnectToServer(u string) error {

	var err error
	var ok bool

	for range only.Once {

		// Default to whatever is defined in me.Server.
		if u == "" {
			u = me.Server.String()
		}

		ok, _ = me.IsConnected()
		if ok {
			if me.Server.String() == u {
				eblog.Debug(me.EntityId, "MqttClient already connected to broker %s", me.Server.String())
				break
			}

			eblog.Debug(me.EntityId, "MqttClient disconnecting from broker %s", me.Server.String())
			me.instance.client.Disconnect(500)
		}

		me.State.SetNewAction(states.ActionStart)
		me.Channels.PublishState(me.State)

		me.Server, err = url.Parse(u)
		if err != nil {
			err = me.EntityId.ProduceError("unable to parse config")
			break
		}

		me.instance.options.Servers = []*url.URL{me.Server}
		//me.instance.options.AddBroker(me.Server.String())
		if me.Server.User != nil {
			me.instance.options.SetUsername(me.Server.User.Username())
			p, ok := me.Server.User.Password()
			if !ok {
				// err = errors.New(string(states.StateUnregistered))
				err = me.EntityId.ProduceError(string(states.StateError))
				break
			}
			me.instance.options.SetPassword(p)
		}

		myWill := me.EntityId.CreateTopic(states.ActionStatus).String()
		me.instance.options.SetWill(myWill, "stopped", topics.QosFailure, false)
		eblog.Debug(me.EntityId, "MqttClient setting LWaT as '%s' on %s", myWill, me.Server.String())

		me.instance.client = mqtt.NewClient(me.instance.options)
		if me.instance.client == nil {
			err = me.EntityId.ProduceError("unable to create client")
			break
		}

		topic := me.EntityId.ConstructTopic(me.EntityId, states.ActionGlob)

		// uri.Path[1:len(uri.Path)]
		err = topic.EnsureNotNil()
		if err != nil {
			break
		}

		me.instance.token = me.instance.client.Connect()
		if me.instance.token == nil {
			err = me.EntityId.ProduceError("unable to connect to %s", me.EntityId.String(), me.Server.String())
			break
		}

		for !me.instance.token.WaitTimeout(me.waitTime) {
		}

		err := me.instance.token.Error()
		if err != nil {
			break
		}

		me.State.SetNewState(states.StateStarted, err)
		me.Channels.PublishState(me.State)
		eblog.Debug(me.EntityId, "connected to broker %s", me.Server)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *MqttClient) GlobSubscribe(client messages.MessageAddress) error {

	var err error
	var ok bool

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		ok, err = me.IsConnected()
		if (!ok) || (err != nil) {
			break
		}

		topic := "/" + client.String() + "/#"

		fmt.Printf("Subcribe to: %s\n", topic)
		// me.instance.client.Subscribe(topic, topics.QosAtLeastOnce, MessageHandler)
		me.instance.client.Subscribe(topic, 0, MessageHandler)
		//me.instance.client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		//	fmt.Printf(">> * [%s] %s\n", msg.Topic(), string(msg.Payload()))
		//})
	}

	return err
}


func MessageHandler(client mqtt.Client, message mqtt.Message) {

	fmt.Printf("MessageHandler =>\n")
	fmt.Printf(">> * [%s] %s\n", message.Topic(), string(message.Payload()))

	//fmt.Printf("mqtt.Client => %v\n", client)
	//fmt.Printf("MessageHandler => %v\n", message)
}

