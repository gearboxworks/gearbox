package mqttClient

import (
	"errors"
	"fmt"
	"gearbox/app/logger"
	"gearbox/box"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/heartbeat/gbevents/tasks"
	"gearbox/only"
	oss "gearbox/os_support"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fhmq/hmq/lib/topics"
	"github.com/jinzhu/copier"
	"net/url"
)


func (me *MqttClient) New(OsSupport oss.OsSupporter, args ...Args) error {

	var _args Args
	var err error

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		_args.osSupport = OsSupport
		foo := box.Args{}
		err = copier.Copy(&foo, &_args)
		if err != nil {
			err = errors.New("unable to copy MQTT client config")
			break
		}

		if _args.EntityId == "" {
			_args.EntityId = DefaultEntityId
		}

		if _args.Server == nil {
			_args.Server, err = url.Parse(DefaultServer)
		}

		if _args.waitTime == 0 {
			_args.waitTime = defaultWaitTime
		}

		_args.instance.options = mqtt.NewClientOptions()
		if _args.instance.options == nil {
			err = errors.New("unable to create MQTT client options")
			break
		}
		_args.instance.options.SetClientID(_args.EntityId.String())

		_args.services = make(ServicesMap)

		*me = MqttClient(_args)

		err = me.ConnectToServer(DefaultServer)
		if err != nil {
			_ = me.Channels.PublishCallerState(me.EntityId, messages.MessageStateUnconfigured)

			// If this fails, it'll be handled within the task.
			// So, reset error to nil to avoid parent thinking everything has gone South.
			err = nil
			break
		}

		logger.Debug("MQTTclient init (%s).", me.EntityId.String())
		_ = me.Channels.PublishCallerState(me.EntityId, messages.MessageStateIdle)
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	if me.EnsureNotNil() == nil {
		// Save last state.
		me.Error = err
	}

	return err
}


// Start the MQTT handler.
func (me *MqttClient) StartHandler() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.Task, err = tasks.StartTask(initMqttClient, startMqttClient, monitorMqttClient, stopMqttClient, me)
		if err != nil {
			break
		}

		logger.Debug("started MQTT client handler for %s", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	if me.EnsureNotNil() == nil {
		// Save last state.
		me.Error = err
	}

	return err
}


// Stop the MQTT handler.
func (me *MqttClient) StopHandler() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for u, _ := range me.services {
			logger.Debug("Unsubscribe MQTT client entry %s.", u.String())
			err = me.services[u].Unsubscribe()
			if err != nil {
				break
			}
		}

		err = me.Task.Stop()
		if err != nil {
			break
		}

		logger.Debug("stopped MQTT client handler for %s", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	if me.EnsureNotNil() == nil {
		// Save last state.
		me.Error = err
	}

	return err
}


func (me *MqttClient) ConnectToServer(u string) error {

	var err error

	for range only.Once {

		if u == "" {
			_ = me.Channels.PublishCallerState(me.EntityId, messages.MessageStateUnconfigured)
			err = errors.New(string(messages.MessageStateUnconfigured))
			break
		}

		fmt.Printf("Server: %v\n", me.Server)
		me.Server, err = url.Parse(u)
		if err != nil {
			err = errors.New("unable to parse MQTT client config")
			break
		}

		me.instance.options.AddBroker(me.Server.String())
		if me.Server.User != nil {
			me.instance.options.SetUsername(me.Server.User.Username())
			p, ok := me.Server.User.Password()
			if !ok {
				_ = me.Channels.PublishCallerState(me.EntityId, messages.MessageStateUnconfigured)
				err = errors.New(string(messages.MessageStateUnconfigured))
				break
			}
			me.instance.options.SetPassword(p)
		}

		fmt.Printf("Options: %v\n", me.instance.options)

		myWill := me.EntityId.CreateTopic(messages.SubTopicState).String()
		me.instance.options.SetWill(myWill, "stopped", topics.QosFailure, false)
		fmt.Printf("Will: %v\n", myWill)

		me.instance.client = mqtt.NewClient(me.instance.options)
		if me.instance.client == nil {
			err = errors.New("unable to create MQTT client config")
			break
		}

		topic := me.EntityId.ConstructTopic(me.EntityId, messages.SubTopicGlob)

		// uri.Path[1:len(uri.Path)]
		err = topic.EnsureNotNil()
		if err != nil {
			break
		}

		me.instance.token = me.instance.client.Connect()
		if me.instance.token == nil {
			err = errors.New(fmt.Sprintf("MQTT client %s unable to connect to %s", me.EntityId.String(), me.Server.String()))
			break
		}

		for !me.instance.token.WaitTimeout(me.waitTime) {
		}

		err := me.instance.token.Error()
		if err != nil {
			break
		}
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	if me.EnsureNotNil() == nil {
		// Save last state.
		me.Error = err
	}

	return err
}

