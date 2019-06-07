package network

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
func initZeroConf(task *tasks.Task, i ...interface{}) error {

	var me *ZeroConf
	var err error

	for range only.Once {
		me, err = InterfaceToTypeZeroConf(i[0])
		if err != nil {
			break
		}
		me.Channels.PublishCallerState(me.EntityId, states.StateStarting)

		_ = task.SetRetryLimit(DefaultRetries)
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
		err = me.ChannelHandler.Subscribe(messages.SubTopic("scan"), scanServices, me)
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

		eblog.Debug("ZeroConf %s initialized OK", me.EntityId.String())
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
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

		// Already started as part of initZeroConf().

		eblog.Debug("ZeroConf %s started OK", me.EntityId.String())
		me.Channels.PublishCallerState(me.EntityId, states.StateStarted)
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
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

		me.scannedServices = make(ServicesMap)
		var found ServicesMap
		for _, s := range browseList {
			found, err = me.Browse(s, me.domain)
			eblog.Debug("ZeroConf browse results(%v, %v)\n", s, me.domain)
			if err != nil {
				me.Channels.PublishCallerState(me.EntityId, states.StateError)
				break
			}

			for k, v := range found {
				me.scannedServices[k] = v
			}
		}

		// Update services
		err = me.updateRegisteredServices()
		if err != nil {
			me.Channels.PublishCallerState(me.EntityId, states.StateError)
			break
		}

		eblog.Debug("ZeroConf %s status", me.EntityId.String())
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
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
		me.Channels.PublishCallerState(me.EntityId, states.StateStopping)

		err = me.ChannelHandler.StopHandler()
		if err != nil {
			me.Channels.PublishCallerState(me.EntityId, states.StateError)
			break
		}

		eblog.Debug("ZeroConf %s stopped", me.EntityId.String())
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
	}

	return err
}
