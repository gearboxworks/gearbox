package vmbox

import (
	"fmt"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/messages"
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
			err = me.channelHandler.Subscribe(states.ActionUpdate, updateHandler, me, states.InterfaceTypeError)
			if err != nil {
				break
			}
			err = me.channelHandler.Subscribe(states.ActionRegister, createHandler, me, states.InterfaceTypeStatus)
			if err != nil {
				break
			}


			// Hard coded for one Vm named "gearbox" for now...
			sc := ServiceConfig{
				Name: messages.MessageAddress(me.Boxname),
				Version: "latest",
			}
			myVM, err := me.New(sc)
			if err == nil {
				fmt.Printf("VM: %v\n", myVM.State.GetStatus())
			}

			fmt.Printf("F: %v\n", me.vms)


			//var state states.Status
			//state, err = myVM.Status()
			//fmt.Printf("Status: %s\n", state.String())


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
		//
		//var state int
		////state, err = me.Releases.Selected.IsIsoFilePresent()
		//state = IsoFileIsDownloading
		//switch state {
		//	case IsoFileNeedsToDownload:
		//		// Leave it to the user.
		//
		//	case IsoFileIsDownloading:
		//		// Send periodic messages to Heartbeat.
		//		//msg := messages.Message{
		//		//	Source: me.EntityId,
		//		//	Topic: messages.MessageTopic{
		//		//		Address:  entity.BroadcastEntityName,
		//		//		SubTopic: states.ActionUpdate,
		//		//	},
		//		//	Text: messages.MessageText(fmt.Sprintf("%s:%d", me.Releases.Selected.File.String(), me.Releases.Selected.DlIndex)),
		//		//}
		//		//_ = me.Channels.Publish(msg)
		//
		//	case IsoFileDownloaded:
		//		// Nothing to do.
		//}
		//if state != IsoFileDownloaded {
		//	//
		//}
		//
		//f := messages.MessageAddress("update")
		//msg := f.ConstructMessage(entity.BroadcastEntityName, states.ActionStatus, "90%")
		////fmt.Printf("EXPECTING: %s\n", msg.String())
		//_ = me.Channels.Publish(msg)


		// Next do something else.
		for range only.Once {

			for _, v := range me.vms {
				state, err := v.Status()
				if err != nil {
					continue
				}
				v.channels.PublishState(&state)
			}

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

