package daemon

import (
	"encoding/json"
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
	"time"
)


////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

// Register a service by method defined by a *CreateEntry structure.
func (me *Daemon) Register(c CreateEntry) (*Service, error) {

	var err error
	var sc Service

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = me.IsExisting(c)
		if err != nil {
			fmt.Printf("PIP! %v\n", time.Now().Unix())
			break
		}

		// Create new daemon entry.
		for range only.Once {
			sc.State.SetNewAction(states.ActionRegister)
			sc.EntityId = messages.GenerateAddress()
			sc.IsManaged = true
			sc.channels = me.Channels
			channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)

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

			// Make sure it's not already present.
			sc.State, err = sc.Status()

			//me.Channels.PublishCallerState(&u, &state)
			//s := sc.State.GetCurrent()

			// Already started. Stop it.
			if sc.State.Current == states.StateStarted {
				err = sc.instance.service.Uninstall()
				if err != nil {
					break
				}
			}

			// Already registered. Remove it.
			if sc.State.Current == states.StateStopped {
				err = sc.instance.service.Uninstall()
				if err != nil {
					break
				}
			}

			// Attempt to install new service on O/S.
			err = sc.instance.service.Install()
			if err != nil {
				break
			}

			sc.State, err = sc.Status()
			if err != nil {
				break
			}

			me.mutex.Lock()
			me.daemons[sc.EntityId] = &sc	// Managed by Mutex
			me.mutex.Unlock()

			eblog.Debug(me.EntityId, "registered service %s OK", sc.Entry.Url)
		}

		sc.State.SetNewState(states.StateRegistered, err)
		me.Channels.PublishCallerState(&sc.EntityId, &sc.State)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

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

		eblog.Debug(me.EntityId, "registered service by channel %s OK", sc.EntityId.String())
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

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

		eblog.Debug(me.EntityId, "registered service by file %s OK", f)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

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

		eblog.Debug(me.EntityId, "registered service by channel %s OK", sc.EntityId.String())
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return sc
}

