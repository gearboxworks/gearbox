package mqttClient

import (
	"gearbox/app/logger"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/heartbeat/gbevents/tasks"
	"gearbox/only"
)


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
		_ = me.Channels.PublishCallerState(me.EntityId, messages.MessageStateStarting)

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
		_ = me.Channels.PublishCallerState(me.EntityId, messages.MessageStateUp)
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
		_ = me.Channels.PublishCallerState(me.EntityId, messages.MessageStateStopping)

		err = me.ChannelHandler.StopHandler()
		if err != nil {
			_ = me.Channels.PublishCallerState(me.EntityId, messages.MessageStateError)
			break
		}

		logger.Debug("MqttClient %s stopped", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return err
}

