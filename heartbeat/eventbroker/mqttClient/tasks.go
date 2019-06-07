package mqttClient

import (
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/heartbeat/eventbroker/tasks"
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
		_ = me.Channels.PublishCallerState(me.EntityId, states.StateStarting)

		_ = task.SetRetryLimit(DefaultRetries)
		_ = task.SetRetryDelay(DefaultRetryDelay)

		me.ChannelHandler, err = me.Channels.StartClientHandler(me.EntityId)
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

		eblog.Debug("MqttClient %s initialized OK", me.EntityId.String())
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
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

		eblog.Debug("MqttClient %s started OK", me.EntityId.String())
		_ = me.Channels.PublishCallerState(me.EntityId, states.StateStarted)
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
	}

	return err
}


// Non-exposed task function - M-DNS monitoring.
func monitorMqttClient(task *tasks.Task, i ...interface{}) error {

	var me *MqttClient
	var err error
	var ok bool

	for range only.Once {
		me, err = InterfaceToTypeMqttClient(i[0])
		if err != nil {
			break
		}

		err = me.ConnectToServer(DefaultServer)
		if err != nil {
			me.Channels.PublishSpecificCallerState(&me.EntityId, states.StateUnregistered)

			// If this fails, it'll be handled within the task.
			// So, reset error to nil to avoid parent thinking everything has gone South.
			err = nil
			break
		}


		ok, _ = me.IsConnected()
		if !ok {
			// We're not connected. Try and connect to whatever is defined in me.Server
			err = me.ConnectToServer("")
			if err != nil {
				_ = me.Channels.PublishCallerState(me.EntityId, states.StateError)
			}
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

		eblog.Debug("MqttClient %s status", me.EntityId.String())
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
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
		_ = me.Channels.PublishCallerState(me.EntityId, states.StateStopping)

		err = me.ChannelHandler.StopHandler()
		if err != nil {
			_ = me.Channels.PublishCallerState(me.EntityId, states.StateError)
			break
		}

		eblog.Debug("MqttClient %s stopped", me.EntityId.String())
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
	}

	return err
}

