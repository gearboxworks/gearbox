package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/network"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/only"
	"github.com/kardianos/service"
	"net/url"
	"os/exec"
	"path/filepath"
)


////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

// Register a service by method defined by a *CreateEntry structure.
func (me *Daemon) Register(c CreateEntry) (*Service, error) {

	var err error
	var sc Service
	var state states.Status

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = me.daemons.IsExisting(c)
		if err != nil {
			break
		}

		// Create new daemon entry.
		for range only.Once {
			sc.EntityId = messages.GenerateAddress()
			sc.IsManaged = true
			sc.channels = me.Channels
			sc.State.SetNewWantState(states.StateRegistered)

			// Create a new service entry.
			sc.Entry, err = me.createEntry(c)
			if err != nil {
				break
			}

			sc.instance.Config = sc.Entry.ToServiceType()

			sc.instance.cmd = &exec.Cmd{
				Path: sc.instance.Executable,
				Args: sc.instance.Arguments,
				Env:  sc.Entry.Env,
				Dir:  sc.instance.WorkingDirectory,
			}

			sc.instance.exit = make(chan struct{})

			sc.instance.service, err = service.New(&sc.instance, sc.instance.Config)
			if err != nil {
				break
			}

			// Attempt to install new service on O/S.
			//state, err = sc.Status()
			//switch {
			//	case state.Current == states.StateUnknown:
			//		// Drop through.
			//
			//	case state.Current == states.StateStarted:
			//		err = sc.instance.service.Stop()
			//		if err != nil {
			//			break
			//		}
			//		err = sc.instance.service.Uninstall()
			//		if err != nil {
			//			break
			//		}
			//
			//	case state.Current == states.StateStopped:
			//		err = sc.instance.service.Uninstall()
			//		if err != nil {
			//			break
			//		}
			//}

			err = sc.instance.service.Install()
			if err != nil {
				break
			}

			state, err = sc.Status()

			// Everything created, store new instance.
			me.daemons[sc.EntityId] = &sc

			eblog.Debug("Daemon %s registered service %s OK", me.EntityId.String(), sc.Entry.Url)
		}

		sc.State.SetNewState(states.StateRegistered, err)
		sc.channels.PublishCallerState(&sc.EntityId, &sc.State)
	}
	eblog.LogIfError(&me, err)

	return &sc, err
}

// Register a service via a channel defined by a *CreateEntry structure and
// returns a *Service structure if successful.
func (me *Daemon) RegisterByChannel(caller messages.MessageAddress, s CreateEntry) (*network.Service, error) {

	var err error
	var j []byte
	var sc *network.Service

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		j, err = json.Marshal(s)
		if err != nil {
			break
		}

		// reg := me.EntityId.Construct(caller, states.ActionRegister, messages.MessageText(j))
		reg := caller.ConstructMessage(me.EntityId, states.ActionRegister, messages.MessageText(j))
		err = me.Channels.Publish(reg)
		if err != nil {
			break
		}

		rs, err := me.Channels.GetCallbackReturn(reg, 100)
		if err != nil {
			break
		}

		sc, err = network.InterfaceToTypeService(rs)	// sc = rs.(*Service)
		if err != nil {
			break
		}

		eblog.Debug("Daemon %s registered service %s via channel", me.EntityId.String())
	}
	eblog.LogIfError(&me, err)

	return sc, err
}


// Register a service by method defined by a *CreateEntry structure.
func (me *Daemon) RegisterByFile(f string) (*Service, error) {

	var err error
	var sc *CreateEntry
	var s *Service

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		sc, err = ReadJsonConfig(f)
		if err != nil {
			break
		}

		s, err = me.Register(*sc)
		if err != nil {
			break
		}

		s.JsonFile = f

		eblog.Debug("Daemon %s registered service %s by file OK", me.EntityId.String(), s.Entry.Url)
	}
	eblog.LogIfError(&me, err)

	return s, err
}


// Create a service by method defined by a *CreateEntry structure.
func (me *Daemon) createEntry(c CreateEntry) (*CreateEntry, error) {

	var err error
	var sc *CreateEntry
	var u *url.URL

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		sc = &c

		switch {
			case sc.MdnsType == "":
				err = me.EntityId.ProduceError("service MdnsType not defined")
				break

			case sc.Name == "":
				err = me.EntityId.ProduceError("service Name not defined")
				break

			case sc.DisplayName == "":
				err = me.EntityId.ProduceError("service DisplayName not defined")
				break

			case sc.Description == "":
				err = me.EntityId.ProduceError("service Description not defined")
				break

			case sc.Executable == "":
				err = me.EntityId.ProduceError("service Executable not defined")
				break
		}

		// u, err = url.Parse(fmt.Sprintf("tcp://%s:%d", mqttService.Entry.HostName, mqttService.Entry.Port))
		u, err = url.Parse(sc.Url)
		if err != nil {
			break
		}

		if u.Port() == "0" {
			var p network.Port
			p, err = network.GetFreePort()
			if err != nil {
				break
			}
			u.Host = u.Host + ":" + p.String()
		}
		sc.Url = u.String()
		sc.Host = network.Host(u.Hostname())
		sc.Port = network.StringToPort(u.Port())

		//if sc.UserName == "" {
		//	sc.UserName = ""
		//}
		//
		//if len(sc.Dependencies) == 0 {
		//	sc.Dependencies = []string{}
		//}

		sc.ChRoot = me.ParsePaths(*sc, sc.ChRoot)
		err = me.CreateDirPaths(sc.ChRoot)
		if err != nil {
			break
		}

		sc.Executable = me.ParsePaths(*sc, sc.Executable)
		err = me.CreateDirPaths(sc.Executable)
		if err != nil {
			break
		}

		if sc.WorkingDirectory == "" {
			sc.WorkingDirectory = filepath.FromSlash(fmt.Sprintf("%s/%s", me.osSupport.GetAdminRootDir(), defaultBaseDir))
		} else {
			sc.WorkingDirectory = me.ParsePaths(*sc, sc.WorkingDirectory)
		}
		err = me.CreateDirPaths(sc.WorkingDirectory)
		if err != nil {
			break
		}

		if sc.Stdout == "" {
			sc.Stdout = filepath.FromSlash(fmt.Sprintf("%s/%s/%s.log",
				me.osSupport.GetAdminRootDir(),
				defaultLogBaseDir,
				sc.Name))
		} else {
			sc.Stdout = me.ParsePaths(*sc, sc.Stdout)
		}
		err = me.CreateDirPaths(sc.Stdout)
		if err != nil {
			break
		}

		if sc.Stderr == "" {
			sc.Stderr = filepath.FromSlash(fmt.Sprintf("%s/%s/%s-error.log",
				me.osSupport.GetAdminRootDir(),
				defaultLogBaseDir,
				sc.Name))
		} else {
			sc.Stderr = me.ParsePaths(*sc, sc.Stderr)
		}
		err = me.CreateDirPaths(sc.Stderr)
		if err != nil {
			break
		}

		for k, v := range sc.Env {
			sc.Env[k] = me.ParsePaths(*sc, v)
			err = me.CreateDirPaths(sc.Env[k])
			if err != nil {
				break
			}
		}

		for k, v := range sc.Arguments {
			sc.Arguments[k] = me.ParsePaths(*sc, v)
			err = me.CreateDirPaths(sc.Arguments[k])
			if err != nil {
				break
			}
		}
	}

	return sc, err
}


////////////////////////////////////////////////////////////////////////////////
// Executed from a channel.

// Non-exposed channel function that responds to a "register" channel request.
func registerService(event *messages.Message, i channels.Argument) channels.Return {

	var me *Daemon
	var sc *Service
	var err error

	for range only.Once {
		me, err = InterfaceToTypeDaemon(i)
		if err != nil {
			break
		}

		//fmt.Printf("Rx: %v\n", event)

		ce := CreateEntry{}
		err = json.Unmarshal(event.Text.ByteArray(), &ce)
		if err != nil {
			break
		}

		sc, err = me.Register(ce)
		if err != nil {
			break
		}
	}
	eblog.LogIfError(&me, err)

	return sc
}


//// Create a service by method defined by a *CreateEntry structure.
//func (me *Daemon) createService(c CreateEntry) (*Service, error) {
//
//	var err error
//	var sc Service
//	var u *url.URL
//
//	for range only.Once {
//		err = me.EnsureNotNil()
//		if err != nil {
//			break
//		}
//
//		sc.Entry = &c
//
//		switch {
//			case sc.MdnsType == "":
//			err = me.EntityId.ProduceError("service MdnsType not defined")
//			break
//
//			case sc.Name == "":
//			err = me.EntityId.ProduceError("service Name not defined")
//			break
//
//			case sc.DisplayName == "":
//			err = me.EntityId.ProduceError("service DisplayName not defined")
//			break
//
//			case sc.Description == "":
//			err = me.EntityId.ProduceError("service Description not defined")
//			break
//
//			case sc.Executable == "":
//			err = me.EntityId.ProduceError("service Executable not defined")
//			break
//		}
//
//		//if sc.Entry.UserName == "" {
//		//	sc.Entry.UserName = ""
//		//}
//		//
//		//if len(sc.Entry.Dependencies) == 0 {
//		//	sc.Entry.Dependencies = []string{}
//		//}
//
//		sc.Entry.ChRoot = me.ParsePaths(sc, sc.Entry.ChRoot)
//		err = me.CreateDirPaths(sc.Entry.ChRoot)
//		if err != nil {
//			break
//		}
//
//		sc.Entry.Executable = me.ParsePaths(sc, sc.Entry.Executable)
//		err = me.CreateDirPaths(sc.Entry.Executable)
//		if err != nil {
//			break
//		}
//
//		if sc.Entry.WorkingDirectory == "" {
//			sc.Entry.WorkingDirectory = filepath.FromSlash(fmt.Sprintf("%s/%s", me.osSupport.GetAdminRootDir(), defaultBaseDir))
//		} else {
//			sc.Entry.WorkingDirectory = me.ParsePaths(sc, sc.Entry.WorkingDirectory)
//		}
//		err = me.CreateDirPaths(sc.Entry.WorkingDirectory)
//		if err != nil {
//			break
//		}
//
//		if sc.Entry.Stdout == "" {
//			sc.Entry.Stdout = filepath.FromSlash(fmt.Sprintf("%s/%s/%s.log",
//				me.osSupport.GetAdminRootDir(),
//				defaultLogBaseDir,
//				sc.Entry.Name))
//		} else {
//			sc.Entry.Stdout = me.ParsePaths(sc, sc.Entry.Stdout)
//		}
//		err = me.CreateDirPaths(sc.Entry.Stdout)
//		if err != nil {
//			break
//		}
//
//		if sc.Entry.Stderr == "" {
//			sc.Entry.Stderr = filepath.FromSlash(fmt.Sprintf("%s/%s/%s-error.log",
//				me.osSupport.GetAdminRootDir(),
//				defaultLogBaseDir,
//				sc.Entry.Name))
//		} else {
//			sc.Entry.Stderr = me.ParsePaths(sc, sc.Entry.Stderr)
//		}
//		err = me.CreateDirPaths(sc.Entry.Stderr)
//		if err != nil {
//			break
//		}
//
//		// u, err = url.Parse(fmt.Sprintf("tcp://%s:%d", mqttService.Entry.HostName, mqttService.Entry.Port))
//		u, err = url.Parse(sc.Entry.Url)
//		if err != nil {
//			break
//		}
//
//		if u.Port() == "0" {
//			var p network.Port
//			p, err = network.GetFreePort()
//			if err != nil {
//				break
//			}
//			u.Host = u.Host + ":" + p.String()
//		}
//		sc.Entry.Url = u.String()
//		sc.host = network.Host(u.Hostname())
//		sc.port = network.StringToPort(u.Port())
//
//		for k, v := range sc.Entry.Env {
//			sc.Entry.Env[k] = me.ParsePaths(sc, v)
//			err = me.CreateDirPaths(sc.Entry.Env[k])
//			if err != nil {
//				break
//			}
//		}
//
//		for k, v := range sc.Entry.Arguments {
//			sc.Entry.Arguments[k] = me.ParsePaths(sc, v)
//			err = me.CreateDirPaths(sc.Entry.Arguments[k])
//			if err != nil {
//				break
//			}
//		}
//
//		sc.instance.Config = sc.Entry.Convert()
//
//		sc.instance.cmd = &exec.Cmd{
//			Path: sc.instance.Executable,
//			Args: sc.instance.Arguments,
//			Env: sc.Entry.Env,
//			Dir: sc.instance.WorkingDirectory,
//		}
//
//		sc.instance.exit = make(chan struct{})
//
//		sc.instance.service, err = service.New(&sc.instance, sc.instance.Config)
//		if err != nil {
//			break
//		}
//
//		sc.EntityId = messages.GenerateAddress()
//		sc.IsManaged = true
//		sc.channels = me.Channels
//		sc.State.SetNewState(states.StateRegistered)
//
//		me.daemons[sc.EntityId] = &sc
//
//		eblog.Debug("Daemon %s registered service %s OK", me.EntityId.String(), u.String())
//	}
//
//	return &sc, err
//}

