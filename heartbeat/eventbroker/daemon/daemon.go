package daemon

import (
	"fmt"
	"gearbox/box"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/heartbeat/eventbroker/tasks"
	"gearbox/only"
	oss "gearbox/os_support"
	"github.com/jinzhu/copier"
	"time"
)


func (me *Daemon) New(OsSupport oss.OsSupporter, args ...Args) error {

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

		//if _args.waitTime == 0 {
		//	_args.waitTime = defaultWaitTime
		//}
		//if _args.restartAttempts == 0 {
		//	_args.restartAttempts = defaultRetries
		//}

		_args.daemons = make(ServicesMap)

		*me = Daemon(_args)


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
func (me *Daemon) StartHandler() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.State.SetNewState(states.StateStarting, err)
		channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)

		for range only.Once {
			me.Task, err = tasks.StartTask(initDaemon, startDaemon, monitorDaemon, stopDaemon, me)
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


// Stop the daemon handler.
func (me *Daemon) StopHandler() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.State.SetNewState(states.StateStopping, err)
		channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)

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


func (me *Daemon) StopServices() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for u, _ := range me.daemons {
			if me.daemons[u].IsManaged {
				_ = me.daemons[u].Stop()
				// Ignore error, will clean up when program exits.
			}
		}
	}

	channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Daemon) TestMe() error {

	var err error

	fmt.Printf("DEBUG STARTED\n")

	var s *Service
	s, err = me.RegisterByFile("/Users/mick/.gearbox/admin/dist/eventbroker/unfsd/unfsd.json")

	time.Sleep(time.Second * 10)
	if err == nil {
		var state states.Status

		state, err = s.Status()
		fmt.Printf("Status: %v\n", state)
		if err != nil {
			fmt.Printf("Woops!\n")
		}
		//
		//err = s.instance.service.Start()
		//if err != nil {
		//	fmt.Printf("Woops!\n")
		//}
		//
		//state, err = s.Status()
		//fmt.Printf("Status: %v\n", state)
		//if err != nil {
		//	fmt.Printf("Woops!\n")
		//}
		//
		//err = s.Stop()
		//if err != nil {
		//	fmt.Printf("Woops!\n")
		//}
		//
		//err = s.instance.service.Uninstall()
		//if err != nil {
		//	fmt.Printf("Woops!\n")
		//}
	}

	fmt.Printf("DEBUG SLEEPING")
	time.Sleep(time.Hour * 4200)

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}

