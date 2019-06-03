package network

import (
	"errors"
	"fmt"
	"gearbox/app/logger"
	"gearbox/box"
	"gearbox/heartbeat/gbevents/channels"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/heartbeat/gbevents/tasks"
	"gearbox/only"
	oss "gearbox/os_support"
	"github.com/jinzhu/copier"
)


// Instantiate a new network structure.
func (me *ZeroConf) New(OsSupport oss.OsSupporter, args ...Args) error {

	var _args Args
	var err error

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		_args.osSupport = OsSupport
		foo := box.Args{}
		err = copier.Copy(&foo, &_args)
		if err != nil {
			err = errors.New("unable to copy MQTT client config")
			break
		}

		if _args.EntityId == "" {
			_args.EntityId = defaultEntityId
		}

		if _args.domain == "" {
			_args.domain = defaultDomain
		}

		if _args.waitTime == 0 {
			_args.waitTime = defaultWaitTime
		}

		if _args.restartAttempts == 0 {
			_args.restartAttempts = defaultRetries
		}

		_args.services = make(ServicesMap)

		*me = ZeroConf(_args)
		logger.Debug("GBevents - ZeroConf init (%s).", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	// Save last state.
	me.Error = err
	return err
}

// Start the M-DNS network handler.
func (me *ZeroConf) StartHandler() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.Task, err = tasks.StartTask(initZeroConf, startZeroConf, monitorZeroConf, stopZeroConf, me)
		if err != nil {
			break
		}

		logger.Debug("started zeroconf handler for %s", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	// Save last state.
	me.Error = err
	return err
}

// Stop the M-DNS network handler.
func (me *ZeroConf) StopHandler() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for u, _ := range me.services {
			logger.Debug("Unregister zeroconf entry %s.", u.String())
			err = me.services[u].Unregister()
			if err != nil {
				break
			}
		}

		err = me.Task.Stop()
		if err != nil {
			break
		}

		logger.Debug("stopped zeroconf handler for %s", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	// Save last state.
	me.Error = err
	return err
}

// Print all services registered under M-DNS that I manage.
func (me *ZeroConf) PrintServices() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		_ = me.services.Print()
	}

	return err
}


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

		err = nil
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

		err = nil
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

		logger.Debug("ZeroConf %s status", me.EntityId.String())

		err = nil
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

		err = me.ChannelHandler.StopHandler()
		if err != nil {
			break
		}

		logger.Debug("ZeroConf %s stopped", me.EntityId.String())

		err = nil
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return err
}


////////////////////////////////////////////////////////////////////////////////
// Executed from a channel

// Non-exposed channel function that responds to a "status" channel request.
// Produces the status of the M-DNS handler via a channel.
func statusHandler(event *messages.Message, i channels.Argument) channels.Return {

	var err error
	var me *ZeroConf

	for range only.Once {
		me, err = InterfaceToTypeZeroConf(i)
		if err != nil {
			break
		}

		logger.Debug("ZeroConf %s handler status OK", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return err
}

// Non-exposed channel function that responds to an "stop" channel request.
// Causes the M-DNS handler task to stop via a channel.
func stopHandler(event *messages.Message, i channels.Argument) channels.Return {

	var err error
	var me *ZeroConf

	for range only.Once {
		me, err = InterfaceToTypeZeroConf(i)
		if err != nil {
			break
		}

		err = me.StopHandler()
		if err != nil {
			break
		}

		logger.Debug("ZeroConf %s handler stopped OK", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return err
}


// Non-exposed channel function that responds to an "start" channel request.
// Causes the M-DNS handler task to start via a channel.
func startHandler(event *messages.Message, i channels.Argument) channels.Return {

	var err error
	var me *ZeroConf

	for range only.Once {
		me, err = InterfaceToTypeZeroConf(i)
		if err != nil {
			break
		}

		err = me.StartHandler()
		if err != nil {
			break
		}

		logger.Debug("ZeroConf %s handler started OK", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return err
}
