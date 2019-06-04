package network

import (
	"errors"
	"gearbox/app/logger"
	"gearbox/box"
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

		logger.Debug(" ZeroConf init (%s).", me.EntityId.String())
		_ = me.Channels.PublishCallerState(me.EntityId, messages.MessageStateIdle)
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
