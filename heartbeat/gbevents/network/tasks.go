package network

import (
	"fmt"
	"gearbox/app/logger"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/heartbeat/gbevents/tasks"
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
		_ = me.Channels.PublishCallerState(me.EntityId, messages.MessageStateStarting)

		_ = task.SetRetryLimit(defaultRetries)
		_ = task.SetRetryDelay(DefaultRetryDelay)

		me.ChannelHandler, err = me.Channels.StartHandler(me.EntityId)
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

		logger.Debug("ZeroConf %s initialized OK", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
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

		logger.Debug("ZeroConf %s started OK", me.EntityId.String())
		_ = me.Channels.PublishCallerState(me.EntityId, messages.MessageStateUp)
	}

	if err != nil {
		logger.Debug("Error: %v", err)
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
			fmt.Printf("Browse(%v, %v) => %v\n", s, me.domain, found)
			if err != nil {
				_ = me.Channels.PublishCallerState(me.EntityId, messages.MessageStateError)
				break
			}

			for k, v := range found {
				me.scannedServices[k] = v
			}
		}

		// Update services
		err = me.updateRegisteredServices()
		if err != nil {
			_ = me.Channels.PublishCallerState(me.EntityId, messages.MessageStateError)
			break
		}

		logger.Debug("ZeroConf %s status", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
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
		_ = me.Channels.PublishCallerState(me.EntityId, messages.MessageStateStopping)

		err = me.ChannelHandler.StopHandler()
		if err != nil {
			_ = me.Channels.PublishCallerState(me.EntityId, messages.MessageStateError)
			break
		}

		logger.Debug("ZeroConf %s stopped", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return err
}
