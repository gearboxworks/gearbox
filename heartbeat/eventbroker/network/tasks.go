package network

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
func initZeroConf(task *tasks.Task, i ...interface{}) error {

	var me *ZeroConf
	var err error

	for range only.Once {
		me, err = InterfaceToTypeZeroConf(i[0])
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

			err = me.channelHandler.Subscribe(messages.SubTopic("stop"), stopHandler, me, InterfaceTypeError)
			if err != nil {
				break
			}
			err = me.channelHandler.Subscribe(messages.SubTopic("start"), startHandler, me, InterfaceTypeError)
			if err != nil {
				break
			}
			err = me.channelHandler.Subscribe(messages.SubTopic("status"), statusHandler, me, states.InterfaceTypeStatus)
			if err != nil {
				break
			}

			err = me.channelHandler.Subscribe(messages.SubTopic("register"), registerService, me, InterfaceTypeService)
			if err != nil {
				break
			}
			err = me.channelHandler.Subscribe(messages.SubTopic("unregister"), unregisterService, me, InterfaceTypeError)
			if err != nil {
				break
			}
			err = me.channelHandler.Subscribe(messages.SubTopic("get"), getHandler, me, messages.InterfaceTypeSubTopics)
			if err != nil {
				break
			}
			err = me.channelHandler.Subscribe(messages.SubTopic("scan"), scanServices, me, InterfaceTypeError)
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
func startZeroConf(task *tasks.Task, i ...interface{}) error {

	var me *ZeroConf
	var err error

	for range only.Once {
		me, err = InterfaceToTypeZeroConf(i[0])
		if err != nil {
			break
		}

		for range only.Once {
			me.State.SetNewAction(states.ActionStart)
			channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)

			// Already started as part of initDaemon().

			me.State.SetNewState(states.StateStarted, err)
			eblog.Debug(me.EntityId, "task handler init completed OK")
		}

		channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
		eblog.LogIfNil(me, err)
		eblog.LogIfError(me.EntityId, err)
	}

	return err
}


// Non-exposed task function - M-DNS monitoring.
func monitorZeroConf(task *tasks.Task, i ...interface{}) error {

	var me *ZeroConf
	var err error

	for range only.Once {
		me, err = InterfaceToTypeZeroConf(i[0])
		if err != nil {
			break
		}

		for range only.Once {
			me.scannedServices = make(ServicesMap)

			var found ServicesMap

			for _, s := range browseList {
				found, err = me.Browse(s, me.domain)
				eblog.Debug(me.EntityId, "task handler status completed OK (%v, %v)", s, me.domain)
				if err != nil {
					break
				}

				for k, v := range found {
					me.scannedServices[k] = v
				}
			}

			// Update services
			err = me.updateRegisteredServices()
			if err != nil {
				break
			}

			eblog.Debug(me.EntityId, "task handler scanning completed OK")
		}

		channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
		eblog.LogIfNil(me, err)
		eblog.LogIfError(me.EntityId, err)
	}

	return err
}


// Non-exposed task function - M-DNS stop.
func stopZeroConf(task *tasks.Task, i ...interface{}) error {

	var me *ZeroConf
	var err error

	for range only.Once {
		me, err = InterfaceToTypeZeroConf(i[0])
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

