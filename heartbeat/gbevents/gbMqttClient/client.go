package gbMqttClient

import (
	"fmt"
	"gearbox/box"
	"gearbox/heartbeat/daemon"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/help"
	"gearbox/only"
	oss "gearbox/os_support"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/jinzhu/copier"
	"net/url"
	"time"
)


func (me *Client) New(OsSupport oss.OsSupporter, args ...Args) status.Status {

	var _args Args
	var sts status.Status

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		_args.OsSupport = OsSupport
		foo := box.Args{}
		err := copier.Copy(&foo, &_args)
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("unable to copy MQTT client config").
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		if _args.EntityId == "" {
			_args.EntityId = "gearbox-mqtt-client"
		}

		status.Success("GBevents - MQTT client(INIT)").Log()
		if _args.Server == nil {

		}
		fmt.Printf("Server: %v\n", _args.Server)
		me.Server, err = url.Parse(_args.Server.String())
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("unable to parse MQTT client config").
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		_args.Config = mqtt.NewClientOptions()
		_args.Config.AddBroker(me.Server.String())
		_args.Config.SetUsername(me.Server.User.Username())
		password, _ := me.Server.User.Password()
		_args.Config.SetPassword(password)
		_args.Config.SetClientID(_args.EntityId.String())

//		topic := messages.CreateTopic(_args.EntityId.String())
//		_args.Config.SetWill(topic, "Last will and testament", topics.QosFailure, true)

		_args.client = mqtt.NewClient(_args.Config)

		*me = Client(_args)
	}

	return sts
}


func (me *Client) Start() status.Status {

	var sts status.Status
	status.Success("GBevents - MQTT client(STARTED)").Log()

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		topic := messages.Topic{
			Address: me.EntityId,
			SubTopic: "*",
		}
		// uri.Path[1:len(uri.Path)]
		sts = topic.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		sts = me.connect()
		if is.Error(sts) {
			break
		}

		sts = me.subscribe(topic, nil)
		if is.Error(sts) {
			break
		}

		s := daemon.WaitForSignal()

		status.Success("GBevents - MQTT broker(STOPPED)").Log()
		sts = status.Success("MQTT client exited with signal %v.", s)
	}
	me.Sts = sts
	status.Log(sts)

	return sts
}


func (me *Client) connect() (status.Status) {

	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		me.token = me.client.Connect()
		if me.token == nil {
			sts = status.Fail().
				SetMessage("unexpected software error").
				SetAdditional("", ).
				SetData("").
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		for !me.token.WaitTimeout(3 * time.Second) {
		}

		if err := me.token.Error(); err != nil {
			sts = status.Wrap(err).
				SetMessage("unable to connect to MQTT broker %s", me.Server.String()).
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		sts = status.Success("MQTT client connected to %s OK.", me.Server.String())
	}

	return sts
}


func (me *Client) subscribe(topic messages.Topic, foo mqtt.MessageHandler) status.Status {

	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		cb := msgCallback{Topic: messages.Topic(topic), Function: foo2}
		me.client.Subscribe(cb.Topic.String(), 0, cb.Function)

		sts = status.Success("MQTT client subscribed OK")
	}

	return sts
}


func foo2(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
}


func callbackFunc(client mqtt.Client, msg mqtt.Message) {

	fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
}


