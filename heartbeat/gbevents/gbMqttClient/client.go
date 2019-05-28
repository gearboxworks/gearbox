package gbMqttClient

import (
	"fmt"
	"gearbox/box"
	"gearbox/global"
	"gearbox/heartbeat/daemon"
	"gearbox/help"
	"gearbox/only"
	oss "gearbox/os_support"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/jinzhu/copier"
	"net/url"
	"os"
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

		*me = Client(_args)
	}

	return sts
}


func (me *Client) Start() status.Status {

	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		status.Success("GBevents - MQTT client(STARTED)").Log()
		uri, err := url.Parse(os.Getenv("CLOUDMQTT_URL"))
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("unable to parse MQTT client config").
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		topic := "" // uri.Path[1:len(uri.Path)]
		if topic == "" {
			topic = "test"
		}

		sts = me.subscribe(uri, topic)
		if is.Error(sts) {
			break
		}

/*
		client := connect("pub", uri)
		timer := time.NewTicker(1 * time.Second)
		for t := range timer.C {
			client.Publish(topic, 0, false, t.String())
		}
*/

		s := daemon.WaitForSignal()

		status.Success("GBevents - MQTT broker(STOPPED)").Log()
		sts = status.Success("MQTT client exited with signal %v.", s)
	}
	me.Sts = sts
	status.Log(sts)

	return sts
}


func (me *Client) connect(clientId string, uri *url.URL) (mqtt.Client, status.Status) {

	var sts status.Status
	var client mqtt.Client
	var opts *mqtt.ClientOptions

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		opts, sts = me.createClientOptions(clientId, uri)
		if is.Error(sts) {
			break
		}

		client = mqtt.NewClient(opts)
		token := client.Connect()
		for !token.WaitTimeout(3 * time.Second) {
		}

		if err := token.Error(); err != nil {
			sts = status.Wrap(err).
				SetMessage("unable to connect to MQTT broker").
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		sts = status.Success("%s GBevents - MQTT client connected OK.", global.Brandname)
	}

	return client, sts
}


func (me *Client) createClientOptions(clientId string, uri *url.URL) (*mqtt.ClientOptions, status.Status) {

	var sts status.Status
	var opts *mqtt.ClientOptions

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		opts = mqtt.NewClientOptions()
		opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
		opts.SetUsername(uri.User.Username())
		password, _ := uri.User.Password()
		opts.SetPassword(password)
		opts.SetClientID(clientId)

		sts = status.Success("%s GBevents - MQTT client connected OK.", global.Brandname)
	}

	return opts, sts
}


func (me *Client) subscribe(uri *url.URL, topic string) status.Status {

	var sts status.Status
	var client mqtt.Client

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		client, sts = me.connect("sub", uri)
		if is.Error(sts) {
			break
		}

		client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
			fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
		})

		sts = status.Success("%s GBevents - MQTT client connected OK.", global.Brandname)
	}

	return sts
}

