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
func (me *Daemon) StartHandler() error {

	var err error

	fmt.Printf("DEBUG STARTED\n")

	var s *Service
	s, err = me.RegisterByFile("/Users/mick/.gearbox/admin/dist/eventbroker/unfsd/unfsd.json")

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.Task, err = tasks.StartTask(initDaemon, startDaemon, monitorDaemon, stopDaemon, me)
		if err != nil {
			break
		}

		eblog.Debug("started zeroconf handler for %s", me.EntityId.String())
	}

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

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
	}

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

		_ = me.StopServices()

		err = me.Task.Stop()
		if err != nil {
			break
		}

		eblog.Debug("Daemon service handler stopped")
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
	}

	return err
}


// Print all services registered under daemon that I manage.
func (me *Daemon) PrintServices() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		_ = me.daemons.Print()
	}

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
				err = me.daemons[u].Stop()
				if err != nil {
					eblog.Debug("Daemon service %s could not be stopped", me.daemons[u].Entry.Name)
				} else {
					eblog.Debug("Daemon service %s stopped", me.daemons[u].Entry.Name)
				}
			}
		}

		eblog.Debug("Daemon services stopped")
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
	}

	return err
}

