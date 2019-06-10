package daemon

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

var globalCheck = ""

// Non-exposed task function - M-DNS initialization.
func initDaemon(task *tasks.Task, i ...interface{}) error {

	var me *Daemon
	var err error

	for range only.Once {
		me, err = InterfaceToTypeDaemon(i[0])
		if err != nil {
			break
		}

		for range only.Once {
			me.State.SetNewAction(states.ActionInitialize)
			channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)

			_ = task.SetRetryLimit(defaultRetries)
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
			err = me.channelHandler.Subscribe(messages.SubTopic("load"), loadConfigHandler, me, InterfaceTypeError)
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
func startDaemon(task *tasks.Task, i ...interface{}) error {

	var me *Daemon
	var err error

	for range only.Once {
		me, err = InterfaceToTypeDaemon(i[0])
		if err != nil {
			break
		}

		for range only.Once {
			me.State.SetNewAction(states.ActionStart)
			channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)

			// Read in any new files and load them up.
			err = me.LoadFiles()

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
func monitorDaemon(task *tasks.Task, i ...interface{}) error {

	var me *Daemon
	var err error

	for range only.Once {
		me, err = InterfaceToTypeDaemon(i[0])
		if err != nil {
			break
		}

		for range only.Once {
			_ = me.LoadFiles()

			for _, u := range me.GetManagedEntities() {
				var state states.Status

				state, err = me.daemons[u].Status()	// Managed by Mutex
				if err != nil {
					continue
				}
				switch state.Current {
					case states.StateUnregistered:
						err = me.daemons[u].instance.service.Install()	// Mutex not required
						if err != nil {
							continue
						}

					case states.StateUnknown:
						fallthrough
					case states.StateStopped:
						err = me.daemons[u].Start()	// Managed by Mutex
						if err != nil {
							continue
						}

					case states.StateStarted:
				}
				//if (state.Current == states.StateUnknown) || (state.Current == states.StateStopped) {
				//	err = me.daemons[u].Start()
				//	if err != nil {
				//		continue
				//	}
				//
				//} else if state.Current == states.StateStarted {
				//	//
				//}
			}

			eblog.Debug(me.EntityId, "task handler status OK")
		}

		eblog.LogIfNil(me, err)
		eblog.LogIfError(me.EntityId, err)
	}

	return err
}


// Non-exposed task function - M-DNS stop.
func stopDaemon(task *tasks.Task, i ...interface{}) error {

	var me *Daemon
	var err error

	for range only.Once {
		me, err = InterfaceToTypeDaemon(i[0])
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

