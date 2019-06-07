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


		me.State.SetNewWantState(states.StateIdle)
		if me.State.SetNewState(states.StateIdle, err) {
			eblog.Debug(me.EntityId, "init complete")
		}
	}

	channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
	eblog.LogIfError(me, err)

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

		eblog.Debug("started zeroconf handler for %s", me.EntityId.String())
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
	}

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
			eblog.Debug("Unregister zeroconf entry %s.", u.String())
			err = me.services[u].Unregister()
			if err != nil {
				break
			}
		}

		err = me.Task.Stop()
		if err != nil {
			break
		}

		eblog.Debug("stopped zeroconf handler for %s", me.EntityId.String())
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
	}

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
