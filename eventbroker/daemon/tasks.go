package daemon

import (
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/states"
	"gearbox/eventbroker/tasks"
	"github.com/gearboxworks/go-status/only"
)

////////////////////////////////////////////////////////////////////////////////
// Executed as a task.

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
			me.State.SetNewAction(states.ActionStart)
			me.Channels.PublishState(me.State)

			_ = task.SetRetryLimit(defaultRetries)
			_ = task.SetRetryDelay(DefaultRetryDelay)

			me.channelHandler, err = me.Channels.StartClientHandler(me.EntityId)
			if err != nil {
				break
			}

			err = me.channelHandler.Subscribe(states.ActionStop, stopHandler, me, states.InterfaceTypeError)
			if err != nil {
				break
			}
			err = me.channelHandler.Subscribe(states.ActionStart, startHandler, me, states.InterfaceTypeError)
			if err != nil {
				break
			}
			err = me.channelHandler.Subscribe(states.ActionStatus, statusHandler, me, states.InterfaceTypeStatus)
			if err != nil {
				break
			}

			err = me.channelHandler.Subscribe(states.ActionRegister, registerService, me, InterfaceTypeService)
			if err != nil {
				break
			}
			err = me.channelHandler.Subscribe(states.ActionUnregister, unregisterService, me, states.InterfaceTypeError)
			if err != nil {
				break
			}
			err = me.channelHandler.Subscribe(msgs.SubTopic("get"), getHandler, me, msgs.InterfaceTypeSubTopics)
			if err != nil {
				break
			}
			err = me.channelHandler.Subscribe(msgs.SubTopic("load"), loadConfigHandler, me, states.InterfaceTypeError)
			if err != nil {
				break
			}

			me.State.SetNewState(states.StateStarted, err)
			me.Channels.PublishState(me.State)
			eblog.Debug(me.EntityId, "task handler init completed OK")
		}

		eblog.LogIfNil(me, err)
		eblog.LogIfError(err)
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
			me.Channels.PublishState(me.State)

			// Function to restart services should they die.
			// Will only be executed if monitor task fails.

			me.State.SetNewState(states.StateStarted, err)
			me.Channels.PublishState(me.State)
			eblog.Debug(me.EntityId, "task handler init completed OK")
		}

		eblog.LogIfNil(me, err)
		eblog.LogIfError(err)
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

		// First monitor my current state.
		if me.State.GetCurrent() != states.StateStarted {
			err = msgs.MakeError(me.EntityId, "task needs restarting")
			break
		}

		// Next do something else.
		for range only.Once {
			_ = me.LoadServiceFiles()

			for _, u := range me.GetManagedEntities() {
				var state states.Status

				state, err = me.daemons[u].Status(PublishState) // Managed by Mutex
				if err != nil {
					continue
				}
				switch state.Current {
				case states.StateUnregistered:
					err = me.daemons[u].instance.service.Install() // Mutex not required
					if err != nil {
						continue
					}

				case states.StateUnknown:
					fallthrough
				case states.StateStopped:
					err = me.daemons[u].Start() // Managed by Mutex
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
		eblog.LogIfError(err)
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
			me.Channels.PublishState(me.State)

			err = me.channelHandler.StopHandler()
			if err != nil {
				break
			}

			me.State.SetNewState(states.StateStopped, err)
			me.Channels.PublishState(me.State)
			eblog.Debug(me.EntityId, "task handler stopped OK")
		}

		eblog.LogIfNil(me, err)
		eblog.LogIfError(err)
	}

	return err
}
