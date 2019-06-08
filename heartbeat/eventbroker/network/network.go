package network

import (
	"gearbox/box"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/heartbeat/eventbroker/tasks"
	"gearbox/only"
	oss "gearbox/os_support"
	"github.com/jinzhu/copier"
)


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
			err = me.EntityId.ProduceError("unable to copy config args")
			break
		}

		if _args.Channels == nil {
			err = me.EntityId.ProduceError("channel pointer is nil")
			break
		}

		if _args.EntityId == "" {
			_args.EntityId = DefaultEntityId
		}

		if _args.domain == "" {
			_args.domain = DefaultDomain
		}

		if _args.waitTime == 0 {
			_args.waitTime = DefaultWaitTime
		}

		if _args.restartAttempts == 0 {
			_args.restartAttempts = DefaultRetries
		}

		_args.services = make(ServicesMap)

		*me = ZeroConf(_args)


		me.State.SetWant(states.StateIdle)
		me.State.SetNewState(states.StateIdle, err)
		eblog.Debug(me.EntityId, "init complete")
	}

	channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

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

		me.State.SetNewState(states.StateStarting, err)
		channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
		me.State.SetWant(states.StateStarted)

		for range only.Once {
			me.Task, err = tasks.StartTask(initZeroConf, startZeroConf, monitorZeroConf, stopZeroConf, me)
			if err != nil {
				break
			}
		}

		me.State.SetNewState(states.StateStarted, err)
		eblog.Debug(me.EntityId, "started task handler")
	}

	channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

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

		me.State.SetNewState(states.StateStopping, err)
		channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
		me.State.SetWant(states.StateStopped)

		for range only.Once {
			_ = me.StopServices()
			// Ignore error, will clean up when program exits.

			err = me.Task.Stop()
		}

		me.State.SetNewState(states.StateStopped, err)
		eblog.Debug(me.EntityId, "stopped task handler")
	}

	channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *ZeroConf) StopServices() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for u, _ := range me.services {
			if me.services[u].IsManaged {
				_ = me.UnregisterByUuid(u)
				// Ignore error, will clean up when program exits.
			}
		}
	}

	channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

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

