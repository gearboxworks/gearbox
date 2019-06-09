package mqttClient

import (
	"gearbox/heartbeat/eventbroker/channels"
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

		for range only.Once {
			me.State.SetNewAction(states.ActionInitialize)
			channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)

			_ = task.SetRetryLimit(DefaultRetries)
			_ = task.SetRetryDelay(DefaultRetryDelay)

			me.channelHandler, err = me.Channels.StartClientHandler(me.EntityId)
			if err != nil {
				break
			}
			err = me.channelHandler.Subscribe(messages.SubTopic("subscribe"), subscribeTopic, me)
			if err != nil {
				break
			}
			err = me.channelHandler.Subscribe(messages.SubTopic("unsubscribe"), unsubscribeTopic, me)
			if err != nil {
				break
			}

			err = me.channelHandler.Subscribe(messages.SubTopic("status"), statusHandler, me)
			if err != nil {
				break
			}
			err = me.channelHandler.Subscribe(messages.SubTopic("stop"), stopHandler, me)
			if err != nil {
				break
			}
			err = me.channelHandler.Subscribe(messages.SubTopic("start"), startHandler, me)
			if err != nil {
				break
			}

			me.State.SetNewState(states.StateInitialized, err)
			eblog.Debug(me.EntityId, "task handler init completed OK")
		}

		channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
		eblog.LogIfNil(me, err)
		eblog.LogIfError(me.EntityId, err)
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

		//for range only.Once {
		//	me.State.SetNewAction(states.ActionStart)
		//	channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
		//
		//	// Already started as part of initDaemon().
		//
		//	me.State.SetNewState(states.StateStarted, err)
		//	eblog.Debug(me.EntityId, "task handler init completed OK")
		//}

		channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
		eblog.LogIfNil(me, err)
		eblog.LogIfError(me.EntityId, err)
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

		for range only.Once {
			ok, _ = me.IsConnected()
			if !ok {
				// We're not connected. Try and connect to whatever is defined in me.Server

				err = me.ConnectToServer("")
				if err != nil {
					// If this fails, try again in the next loop.
					break
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

			eblog.Debug(me.EntityId, "task handler status OK")
		}

		channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
		eblog.LogIfNil(me, err)
		eblog.LogIfError(me.EntityId, err)
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

		for range only.Once {
			me.State.SetNewAction(states.ActionStop)
			channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)

			err = me.channelHandler.StopHandler()
			if err != nil {
				break
			}

			me.State.SetNewState(states.StateStopped, err)
			eblog.Debug(me.EntityId, "task handler stopped OK")
		}

		channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
		eblog.LogIfNil(me, err)
		eblog.LogIfError(me.EntityId, err)
	}

	return err
}

