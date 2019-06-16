package vmbox

import (
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/only"
	"gearbox/eventbroker/states"
	"gearbox/eventbroker/tasks"
)


////////////////////////////////////////////////////////////////////////////////
// Executed as a task.


// Non-exposed task function - M-DNS initialization.
func initVmBox(task *tasks.Task, i ...interface{}) error {

	var me *VmBox
	var err error

	for range only.Once {
		me, err = InterfaceToTypeVmBox(i[0])
		if err != nil {
			break
		}

		for range only.Once {
			me.State.SetNewAction(states.ActionInitialize)
			me.Channels.PublishState(me.State)

			_ = task.SetRetryLimit(DefaultRetries)
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

			me.State.SetNewState(states.StateInitialized, err)
			me.Channels.PublishState(me.State)
			eblog.Debug(me.EntityId, "task handler init completed OK")
		}

		eblog.LogIfNil(me, err)
		eblog.LogIfError(me.EntityId, err)
	}

	return err
}


// Non-exposed task function - M-DNS start.
func startVmBox(task *tasks.Task, i ...interface{}) error {

	var me *VmBox
	var err error

	for range only.Once {
		me, err = InterfaceToTypeVmBox(i[0])
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
		eblog.LogIfError(me.EntityId, err)
	}

	return err
}


// Non-exposed task function - M-DNS monitoring.
func monitorVmBox(task *tasks.Task, i ...interface{}) error {

	var me *VmBox
	var err error

	for range only.Once {
		me, err = InterfaceToTypeVmBox(i[0])
		if err != nil {
			break
		}


		//// First monitor my current state.
		//if me.State.GetCurrent() != states.StateStarted {
		//	err = messages.ProduceError(entity.VmBoxEntityName, "task needs restarting")
		//	break
		//}


		// Next do something else.
		for range only.Once {

			me.Channels.PublishState(me.State)
			eblog.Debug(me.EntityId, "task handler status OK")
		}

		eblog.LogIfNil(me, err)
		eblog.LogIfError(me.EntityId, err)
	}

	return err
}


// Non-exposed task function - M-DNS stop.
func stopVmBox(task *tasks.Task, i ...interface{}) error {

	var me *VmBox
	var err error

	for range only.Once {
		me, err = InterfaceToTypeVmBox(i[0])
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
		eblog.LogIfError(me.EntityId, err)
	}

	return err
}

