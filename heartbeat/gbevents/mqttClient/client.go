package mqttClient

import (
	"errors"
	"fmt"
	"gearbox/app/logger"
	"gearbox/box"
	"gearbox/heartbeat/gbevents/channels"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/heartbeat/gbevents/tasks"
	"gearbox/only"
	oss "gearbox/os_support"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gearboxworks/go-status"
	"github.com/jinzhu/copier"
	"net/url"
	"time"
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
			err = errors.New("no MQTT server url defined")
			break
		}

		*me = MqttClient(_args)
		status.Success("GBevents - MQTT client init (%s).", me.EntityId.String()).Log()
	}

	// We're going to throw away the config for now.
	// This will be handled within the task.
	_ = me.ConfigureHandler()

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	// Save last state.
	me.Error = err

	return err
}


func (me *MqttClient) ConfigureHandler() error {

	var err error

	for range only.Once {

		fmt.Printf("Server: %v\n", me.Server)
		me.Server, err = url.Parse(me.Server.String())
		if err != nil {
			err = errors.New("unable to parse MQTT client config")
			break
		}

		me.Config = mqtt.NewClientOptions()
		me.Config.AddBroker(me.Server.String())
		me.Config.SetUsername(me.Server.User.Username())
		password, _ := me.Server.User.Password()
		me.Config.SetPassword(password)
		me.Config.SetClientID(me.EntityId.String())

		//		topic := messages.CreateTopic(me.EntityId.String())
		//		me.Config.SetWill(topic, "Last will and testament", topics.QosFailure, true)

		me.client = mqtt.NewClient(me.Config)
		if me.client == nil {
			err = errors.New("unable to create MQTT client config")
			break
		}
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

		logger.Debug("started zeroconf handler for %s", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	// Save last state.
	me.Error = err
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
			logger.Debug("Unsubscribe zeroconf entry %s.", u.String())
			err = me.services[u].Unsubscribe()
			if err != nil {
				break
			}
		}

		err = me.Task.Stop()
		if err != nil {
			break
		}

		logger.Debug("stopped zeroconf handler for %s", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	// Save last state.
	me.Error = err
	return err
}



////////////////////////////////////////////////////////////////////////////////
// Executed as a task.

// Non-exposed task function - M-DNS initialization.
func initMqttClient(task *tasks.Task, i ...interface{}) error {

	var me *MqttClient
	var err error

	for range only.Once {
		me, err = InterfaceToTypeMqttClient(i[0])
		if err != nil {
			break
		}

		_ = task.SetRetryLimit(DefaultRetries)
		_ = task.SetRetryDelay(DefaultRetryDelay)

		me.ChannelHandler, err = me.Channels.StartHandler(me.EntityId)
		if err != nil {
			break
		}
		err = me.ChannelHandler.Subscribe(messages.SubTopic("subscribe"), subscribeTopic, me)
		if err != nil {
			break
		}
		err = me.ChannelHandler.Subscribe(messages.SubTopic("unsubscribe"), unsubscribeTopic, me)
		if err != nil {
			break
		}

		err = me.ChannelHandler.Subscribe(messages.SubTopic("status"), statusHandler, me)
		if err != nil {
			break
		}
		err = me.ChannelHandler.Subscribe(messages.SubTopic("stop"), stopHandler, me)
		if err != nil {
			break
		}
		err = me.ChannelHandler.Subscribe(messages.SubTopic("start"), startHandler, me)
		if err != nil {
			break
		}

		logger.Debug("MqttClient %s initialized OK", me.EntityId.String())

		err = nil
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return err
}

// Non-exposed task function - M-DNS start.
func startMqttClient(task *tasks.Task, i ...interface{}) error {

	var me *MqttClient
	var err error

	for range only.Once {
		me, err = InterfaceToTypeMqttClient(i[0])
		if err != nil {
			break
		}

		// Already started as part of initMqttClient().

		logger.Debug("MqttClient %s started OK", me.EntityId.String())

		err = nil
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return err
}

// Non-exposed task function - M-DNS monitoring.
func monitorMqttClient(task *tasks.Task, i ...interface{}) error {

	var me *MqttClient
	var err error

	for range only.Once {
		me, err = InterfaceToTypeMqttClient(i[0])
		if err != nil {
			break
		}

		//me.scannedServices = make(ServicesMap)
		//var found ServicesMap
		//for _, s := range browseList {
		//	found, err = me.Browse(s, me.domain)
		//	fmt.Printf("Browse(%v, %v) => %v\n", s, me.domain, found)
		//	if err != nil {
		//		break
		//	}
		//
		//	for k, v := range found {
		//		me.scannedServices[k] = v
		//	}
		//}
		//
		//// Update services
		//err = me.updateRegisteredServices()
		//if err != nil {
		//	break
		//}

		logger.Debug("MqttClient %s status", me.EntityId.String())

		err = nil
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return err
}

// Non-exposed task function - M-DNS stop.
func stopMqttClient(task *tasks.Task, i ...interface{}) error {

	var me *MqttClient
	var err error

	for range only.Once {
		me, err = InterfaceToTypeMqttClient(i[0])
		if err != nil {
			break
		}

		err = me.ChannelHandler.StopHandler()
		if err != nil {
			break
		}

		logger.Debug("MqttClient %s stopped", me.EntityId.String())

		err = nil
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return err
}
////////////////////////////////////////////////////////////////////////////////
// Executed from a channel

// Non-exposed channel function that responds to a "status" channel request.
// Produces the status of the M-DNS handler via a channel.
func statusHandler(event *messages.Message, i channels.Argument) channels.Return {

	var err error
	var me *MqttClient

	for range only.Once {
		me, err = InterfaceToTypeMqttClient(i)
		if err != nil {
			break
		}

		logger.Debug("MqttClient %s handler status OK", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return err
}

// Non-exposed channel function that responds to an "stop" channel request.
// Causes the M-DNS handler task to stop via a channel.
func stopHandler(event *messages.Message, i channels.Argument) channels.Return {

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

		logger.Debug("MqttClient %s handler stopped OK", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return err
}


// Non-exposed channel function that responds to an "start" channel request.
// Causes the M-DNS handler task to start via a channel.
func startHandler(event *messages.Message, i channels.Argument) channels.Return {

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

		logger.Debug("MqttClient %s handler started OK", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return err
}



func (me *MqttClient) initMqttClient() error {

	var err error

	status.Success("GBevents - MQTT client started (%s).", me.EntityId.String()).Log()

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		topic := messages.Topic{
			Address: me.EntityId,
			SubTopic: "*",
		}

		// uri.Path[1:len(uri.Path)]
		err = topic.EnsureNotNil()
		if err != nil {
			break
		}

		err = me.connect()
		if err != nil {
			break
		}

		status.Success("GBevents - MQTT client stopped (%s).", me.EntityId.String()).Log()
	}

	return err
}


func (me *MqttClient) connect() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.token = me.client.Connect()
		if me.token == nil {
			err = errors.New("unable to obtain an MQTT client token")
			break
		}

		for !me.token.WaitTimeout(3 * time.Second) {
		}

		err := me.token.Error()
		if err != nil {
			break
		}
	}

	return err
}

