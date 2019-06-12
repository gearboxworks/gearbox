package daemon

import (
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/entity"
	"gearbox/heartbeat/eventbroker/network"
	"gearbox/heartbeat/eventbroker/only"
	"gearbox/heartbeat/eventbroker/states"
	"github.com/kardianos/service"
	"os"
	"path/filepath"
	"strings"
)


func (srv *Service) Start() error {

	var err error

	for range only.Once {
		err = srv.EnsureNotNil()
		if err != nil {
			break
		}

		srv.State.SetNewAction(states.ActionStart)
		srv.channels.PublishState(&srv.EntityId, &srv.State)

		err = srv.instance.service.Start()
		if err != nil {
			break
		}

		var state states.State
		state, err = srv.decodeServiceState()
		if err != nil {
			break
		}
		// Didn't start.
		if state != states.StateStarted {
			break
		}

		err = srv.RegisterMDNS()
		if err != nil {
			err = srv.Stop()
			if err != nil {
				break
			}

			break
		}

		srv.State.SetNewState(states.StateStarted, err)
		srv.channels.PublishState(&srv.EntityId, &srv.State)
		eblog.Debug(srv.EntityId, "service started OK")
	}

	eblog.LogIfNil(srv, err)
	eblog.LogIfError(srv.EntityId, err)

	return nil
}


func (srv *Service) Stop() error {

	var err error

	for range only.Once {
		err = srv.EnsureNotNil()
		if err != nil {
			break
		}

		srv.State.SetNewAction(states.ActionStop)
		srv.channels.PublishState(&srv.EntityId, &srv.State)

		err = srv.instance.service.Stop()
		if err != nil {
			break
		}

		var state states.State
		state, err = srv.decodeServiceState()
		if err != nil {
			break
		}
		// Didn't stop.
		if state != states.StateStopped {
			break
		}

		err = srv.UnregisterMDNS()
		if err != nil {
			break
		}

		srv.State.SetNewState(states.StateStopped, err)
		srv.channels.PublishState(&srv.EntityId, &srv.State)
		eblog.Debug(srv.EntityId, "service stopped OK")
	}

	eblog.LogIfNil(srv, err)
	eblog.LogIfError(srv.EntityId, err)

	return err
}


func (srv *Service) Status(publish bool) (states.Status, error) {

	var err error

	for range only.Once {
		err = srv.EnsureNotNil()
		if err != nil {
			break
		}

		var state states.State
		state, err = srv.decodeServiceState()

		srv.State.SetNewState(state, err)

		if srv.State.HasChangedState() {
			eblog.Debug(srv.EntityId, "status current:%s last:%s", srv.State.GetCurrent().String(), srv.State.GetLast().String())

			if publish {
				srv.channels.PublishState(&srv.EntityId, &srv.State)
			}
		}
	}

	eblog.LogIfNil(srv, err)
	eblog.LogIfError(srv.EntityId, err)

	return srv.State, err
}


func (srv *Service) decodeServiceState() (states.State, error) {

	var err error
	var state states.State

	for range only.Once {
		err = srv.EnsureNotNil()
		if err != nil {
			break
		}

		var serviceState service.Status
		serviceState, err = srv.instance.service.Status()
		// We want to translate the states defined from the service package.

		// ErrNameFieldRequired is returned when Config.Name is empty.
		// service.ErrNameFieldRequired
		// ErrNoServiceSystemDetected is returned when no system was detected.
		// service.ErrNoServiceSystemDetected
		// ErrNotInstalled is returned when the service is not installed
		// service.ErrNotInstalled

		switch {
			case err == service.ErrNameFieldRequired:
				state = states.StateError

			case err == service.ErrNoServiceSystemDetected:
				state = states.StateError

			case err == service.ErrNotInstalled:
				state = states.StateUnregistered

			case serviceState == service.StatusUnknown:
				state = states.StateUnknown

			case serviceState == service.StatusStopped:
				state = states.StateStopped

			case serviceState == service.StatusRunning:
				state = states.StateStarted

			default:
				state = states.StateError
		}
	}

	return state, err
}


func (srv *Service) RegisterMDNS() error {

	var err error

	for range only.Once {
		err = srv.EnsureNotNil()
		if err != nil {
			break
		}

		// Register with ZeroConf
		var i interface{}
		// msg := network.ConstructMdnsRegisterMessage(srv.EntityId, entity.NetworkEntityName, *srv.MdnsEntry)
		msg := network.ConstructMdnsMessage(srv.EntityId, entity.NetworkEntityName, *srv.MdnsEntry, states.ActionRegister)
		i, err = srv.channels.PublishAndWaitForReturn(msg, 400)
		if err != nil {
			break
		}

		// Might have to store the ZeroConf *sc in the *Service entity.
		_, err = network.InterfaceToTypeService(i)
		if err != nil {
			eblog.Debug(srv.EntityId, "couldn't register MDNS service")
			break
		}

		eblog.Debug(srv.EntityId, "registered MDNS service OK")
	}

	eblog.LogIfNil(srv, err)
	eblog.LogIfError(srv.EntityId, err)

	return nil
}


func (srv *Service) UnregisterMDNS() error {

	var err error

	for range only.Once {
		err = srv.EnsureNotNil()
		if err != nil {
			break
		}

		// Register with ZeroConf
		var i interface{}
		//var sc *network.Service
		// msg := network.ConstructMdnsRegisterMessage(srv.EntityId, entity.NetworkEntityName, *srv.MdnsEntry)
		msg := network.ConstructMdnsMessage(srv.EntityId, entity.NetworkEntityName, *srv.MdnsEntry, states.ActionUnregister)
		i, err = srv.channels.PublishAndWaitForReturn(msg, 400)
		if err != nil {
			break
		}

		_, err = states.InterfaceToTypeError(i)
		if err != nil {
			eblog.Debug(srv.EntityId, "couldn't unregister MDNS service")
			break
		}

		eblog.Debug(srv.EntityId, "unregistered MDNS service OK")
	}

	eblog.LogIfNil(srv, err)
	eblog.LogIfError(srv.EntityId, err)

	return nil
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


func (me *Daemon) FindServiceFiles() ([]string, error) {

	var err error
	var files []string

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		checkIn := me.OsPaths.EventBrokerEtcDir.AddToPath(DefaultJsonDir)
		err = checkIn.DirExists()
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "Finding service files within %s", checkIn)
		err = filepath.Walk(checkIn.String(), func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, ".json") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			break
		}
	}

	return files, err
}

