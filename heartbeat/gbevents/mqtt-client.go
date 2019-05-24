package gbevents

import (
	"fmt"
	"gearbox/global"
	"gearbox/help"
	"gearbox/only"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"net/url"
	"os"
	"time"
)


func (me *ServiceEvents) startMqttClient() status.Status {

	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		fmt.Printf("GBevents - MQTT client(STARTED)\n")
		uri, err := url.Parse(os.Getenv("CLOUDMQTT_URL"))
		if err != nil {
			sts = status.Fail(&status.Args{
				Message: "GBevents - MQTT client error",
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}
		topic := uri.Path[1:len(uri.Path)]
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

		fmt.Printf("GBevents - MQTT client(FINISHED)\n")
		sts = status.Success("%s GBevents - MQTT client exited.", global.Brandname)
	}

	return sts
}


func (me *ServiceEvents) connect(clientId string, uri *url.URL) (mqtt.Client, status.Status) {

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
			sts = status.Fail(&status.Args{
				Message: "GBevents - MQTT client error",
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		sts = status.Success("%s GBevents - MQTT client connected OK.", global.Brandname)
	}

	return client, sts
}


func (me *ServiceEvents) createClientOptions(clientId string, uri *url.URL) (*mqtt.ClientOptions, status.Status) {

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


func (me *ServiceEvents) subscribe(uri *url.URL, topic string) status.Status {

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

