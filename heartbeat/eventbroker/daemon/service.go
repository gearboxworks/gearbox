package daemon

import (
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/only"
	"github.com/kardianos/service"
)


func (me *Service) Start() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = me.instance.service.Start()
		if err != nil {
			break
		}

		// Notify channels
		err = me.channels.PublishCallerState(me.EntityId, states.StateStarted)
		if err != nil {
			break
		}

		//// Now register the service with zeroconf.
		//zc := network.CreateEntry{
		//	Name: network.Name("gearbox_" + me.Entry.Name),
		//	Type: network.Type(fmt.Sprintf("_%s._tcp", me.Entry.MdnsType)),
		//	Domain: "local",
		//	Port: network.Port(me.port),
		//}
		//me.mdns, err = me.RegisterByChannel(me.EntityId, zc)
		//if err != nil {
		//	break
		//}

		me.IsManaged = true

		eblog.Debug("Daemon stopped %s", me.Entry.Name)
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
	}

	return nil
}


func (me *Service) Stop() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = me.instance.service.Stop()
		if err != nil {
			break
		}

		// Notify channels
		err = me.channels.PublishCallerState(me.EntityId, states.StateStarted)
		if err != nil {
			break
		}

		//// Now unregister the service with zeroconf.
		//zc := network.CreateEntry{
		//	Name: network.Name("gearbox_" + me.Entry.Name),
		//	Type: network.Type(fmt.Sprintf("_%s._tcp", me.Entry.MdnsType)),
		//	Domain: "local",
		//	Port: network.Port(me.port),
		//}
		//me.mdns, err = me.RegisterByChannel(me.EntityId, zc)
		//if err != nil {
		//	break
		//}

		me.IsManaged = false

		eblog.Debug("Daemon stopped %s", me.Entry.Name)
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
	}

	return err
}


func (me *Service) Status() (states.Status, error) {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		var serviceState service.Status

		serviceState, err = me.instance.service.Status()
		// We want to translate the states defined from the service package.

		// ErrNameFieldRequired is returned when Config.Name is empty.
		// service.ErrNameFieldRequired
		// ErrNoServiceSystemDetected is returned when no system was detected.
		// service.ErrNoServiceSystemDetected
		// ErrNotInstalled is returned when the service is not installed
		// service.ErrNotInstalled

		var newState states.State
		switch {
			case err == service.ErrNameFieldRequired:
				newState = states.StateError

			case err == service.ErrNoServiceSystemDetected:
				newState = states.StateError

			case err == service.ErrNotInstalled:
				newState = states.StateUnregistered

			case serviceState == service.StatusUnknown:
				newState = states.StateUnknown

			case serviceState == service.StatusStopped:
				newState = states.StateStopped

			case serviceState == service.StatusRunning:
				newState = states.StateStarted

			default:
				newState = states.StateError
		}
		me.State.SetNewState(newState, err)

		if me.State.Last != me.State.Current {
			eblog.Debug("Daemon %s status current:%s last:%s", me.EntityId.String(), me.State.Current.String(), me.State.Last.String())
			channels.PublishCallerState(me.channels, &me.EntityId, &me.State)
		}
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
	}

	return me.State, err
}


func translateState(i service.Status, err error) (states.State, error) {

	var s states.State
	// We're also passing through err unchanged.

	// ErrNameFieldRequired is returned when Config.Name is empty.
	// service.ErrNameFieldRequired
	// ErrNoServiceSystemDetected is returned when no system was detected.
	// service.ErrNoServiceSystemDetected
	// ErrNotInstalled is returned when the service is not installed
	// service.ErrNotInstalled

	switch {
		case err == service.ErrNameFieldRequired:
			s = states.StateError

		case err == service.ErrNoServiceSystemDetected:
			s = states.StateError

		case err == service.ErrNotInstalled:
			s = states.StateUnregistered

		case i == service.StatusUnknown:
			s = states.StateUnknown

		case i == service.StatusStopped:
			s = states.StateStopped

		case i == service.StatusRunning:
			s = states.StateStarted

		default:
			s = states.StateError
	}

	return s, err
}

