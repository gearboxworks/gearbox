package network

import (
	"gearbox/eventbroker/channels"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/osdirs"
	"gearbox/eventbroker/states"
	"gearbox/eventbroker/tasks"
	"github.com/gearboxworks/go-status/only"
	"sync"
	"time"
)

type ZeroConf struct {
	EntityId msgs.Address
	Boxname  string
	State    *states.Status
	Task     *tasks.Task
	Channels *channels.Channels

	mutex           sync.RWMutex // Mutex control for map.
	channelHandler  *channels.Subscriber
	restartAttempts int
	waitTime        time.Duration
	domain          string
	services        ServicesMap
	scannedServices ServicesMap
	OsPaths         *osdirs.BaseDirs
}
type Args ZeroConf

func (me *ZeroConf) New(args ...Args) error {

	var _args Args
	var err error

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		if _args.Channels == nil {
			err = msgs.MakeError(me.EntityId, "channel pointer is nil")
			break
		}

		if _args.OsPaths == nil {
			err = msgs.MakeError(me.EntityId, "ospaths is nil")
			break
		}

		if _args.EntityId == "" {
			_args.EntityId = entity.NetworkEntityName
		}
		_args.State = states.New(_args.EntityId, _args.EntityId, entity.SelfEntityName)

		if _args.Boxname == "" {
			_args.Boxname = entity.NetworkEntityName
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

	me.Channels.PublishState(me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

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

		me.State.SetNewAction(states.ActionStart)
		me.Channels.PublishState(me.State)

		for range only.Once {
			me.Task, err = tasks.StartTask(initZeroConf, startZeroConf, monitorZeroConf, stopZeroConf, me)
			if err != nil {
				break
			}
		}

		me.State.SetNewState(states.StateStarted, err)
		me.Channels.PublishState(me.State)
		eblog.Debug(me.EntityId, "started task handler")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

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

		me.State.SetNewAction(states.ActionStop)
		me.Channels.PublishState(me.State)

		for range only.Once {
			_ = me.StopServices()
			// Ignore error, will clean up when program exits.

			err = me.Task.Stop()
		}

		me.State.SetNewState(states.StateStopped, err)
		me.Channels.PublishState(me.State)
		eblog.Debug(me.EntityId, "stopped task handler")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return err
}

func (me *ZeroConf) StopServices() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for u := range me.services {
			if me.services[u].IsManaged {
				_ = me.UnregisterByEntityId(u)
				// Ignore error, will clean up when program exits.
			}
		}
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

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
