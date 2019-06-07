package daemon

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
func initDaemon(task *tasks.Task, i ...interface{}) error {

	var me *Daemon
	var err error

	for range only.Once {
		me, err = InterfaceToTypeDaemon(i[0])
		if err != nil {
			break
		}
		_ = me.Channels.PublishCallerState(me.EntityId, states.StateStarting)

		_ = task.SetRetryLimit(defaultRetries)
		_ = task.SetRetryDelay(DefaultRetryDelay)

		me.ChannelHandler, err = me.Channels.StartClientHandler(me.EntityId)
		if err != nil {
			break
		}
		err = me.ChannelHandler.Subscribe(messages.SubTopic("register"), registerService, me)
		if err != nil {
			break
		}
		err = me.ChannelHandler.Subscribe(messages.SubTopic("unregister"), unregisterService, me)
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

		eblog.Debug("Daemon %s initialized OK", me.EntityId.String())
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
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

		// Already started as part of initDaemon().

		_ = me.Channels.PublishCallerState(me.EntityId, states.StateStarted)
		eblog.Debug("Daemon %s started OK", me.EntityId.String())
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
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

		for u, _ := range me.daemons {
			if me.daemons[u].IsManaged {
				var state states.Status

				state, err = me.daemons[u].Status()
				if (state.Current == states.StateUnknown) || (state.Current == states.StateStopped) {
					err = me.daemons[u].Start()
					if err != nil {
						continue
					}

				} else if state.Current == states.StateStarted {
					//
				}
			}
		}

		eblog.Debug("Daemon %s status", me.EntityId.String())
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
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
		_ = me.Channels.PublishCallerState(me.EntityId, states.StateStopping)

		err = me.ChannelHandler.StopHandler()
		if err != nil {
			break
		}

		_ = me.Channels.PublishCallerState(me.EntityId, states.StateStopped)
		eblog.Debug("Daemon %s stopped", me.EntityId.String())
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
		_ = me.Channels.PublishCallerState(me.EntityId, states.StateError)
	}

	return err
}
