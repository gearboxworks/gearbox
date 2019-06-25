package daemon

import (
	"encoding/json"
	"fmt"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/network"
	"gearbox/eventbroker/only"
	"gearbox/eventbroker/ospaths"
	"gearbox/eventbroker/states"
	"github.com/kardianos/service"
	"os"
	"os/exec"
	"time"
)


////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

// Register a service by method defined by a *CreateEntry structure.
func (me *Daemon) Register(c ServiceConfig) (*Service, error) {

	var err error
	sc := &Service{}

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		var check *Service
		check, err = me.FindExistingConfig(c)
		if err != nil {
			if check == nil {
				break
			}

			// It's not an error, but rather we've already got a service registered.
			err = nil

			if check.IsRegistered() {
				break
			}

			fmt.Printf("PIP! %v\n", time.Now().Unix())
			break
		}

		// Check platform.
		if c.SkipPlatform() {
			// err = me.EntityId.ProduceError("service shouldn't run on this host")
			sc = nil
			break
		}

		// Create new daemon entry.
		for range only.Once {
			sc.EntityId = *messages.GenerateAddress()
			if c.EntityName != "" {
				sc.EntityName = messages.MessageAddress(c.EntityName)
			} else {
				if c.Name != "" {
					sc.EntityName = messages.MessageAddress(c.Name)
				} else {
					sc.EntityName = sc.EntityId
				}
			}
			sc.EntityParent = &me.EntityId
			sc.State = states.New(&sc.EntityId, &sc.EntityName, me.EntityId)
			sc.State.SetNewAction(states.ActionRegister)
			sc.IsManaged = true
			sc.channels = me.Channels
			sc.channels.PublishState(sc.State)

			sc.Entry, err = me.createEntry(c)
			if err != nil {
				break
			}

			sc.MdnsEntry, err = sc.CreateMdnsEntry()

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

			// Make sure it's not already present.
			var state states.State
			state, err = sc.decodeServiceState()
			// Already registered or started. Stop it.
			if state == states.StateStopped {
				err = sc.instance.service.Uninstall()
				if err != nil {
					break
				}
			} else if state == states.StateStarted {
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

			_, err = sc.decodeServiceState()
			if err != nil {
				break
			}

			err = me.AddEntity(sc.EntityId, sc)
			if err != nil {
				break
			}

			sc.State.SetNewState(states.StateRegistered, err)
			sc.channels.PublishState(sc.State)
			eblog.Debug(me.EntityId, "registered service %s OK", sc.Entry.UrlPtr.String())
		}
	}

	me.Channels.PublishState(me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return sc, err
}


// Register a service via a channel defined by a *CreateEntry structure and
// returns a *Service structure if successful.
func (me *Daemon) RegisterByChannel(caller messages.MessageAddress, s ServiceConfig) (*network.Service, error) {

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

	me.Channels.PublishState(me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return sc, err
}


// Register a service by method defined by a *CreateEntry structure.
func (me *Daemon) RegisterByFile(f string) (*Service, error) {

	var err error
	var sc *ServiceConfig
	var s *Service
	var ok bool

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		ok = me.IsFileRegistered(f)
		if ok {
			ok, err = me.HasFileChanged(f)
			if !ok {
				break
			}
		}

		//err = me.RemoveFileIfExist(f)
		//if err != nil {
		//	break
		//}

		sc, err = ReadJsonConfig(f)
		if err != nil {
			break
		}


		s, err = me.Register(*sc)
		if err != nil {
			break
		}
		if s == nil {
			break
		}


		info, err := os.Stat(f)
		if err != nil {
			break
		}
		s.JsonFile.Name = f
		s.JsonFile.LastModTime = info.ModTime()

		eblog.Debug(me.EntityId, "registered service by file %s OK", f)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return s, err
}


func (me *Daemon) LoadServiceFiles() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		var files []string
		files, err = me.FindServiceFiles()
		if err != nil {
			break
		}

		for _, file := range files {
			var sc *Service
			sc, err = me.RegisterByFile(file)
			if sc == nil {
				//eblog.Debug(me.EntityId, "Unloaded service file %s", file)
				continue
			}
			if err != nil {
				eblog.Debug(me.EntityId, "Loading service file %s failed with '%v'\n", file, err)
				continue
			}
			eblog.Debug(me.EntityId, "Loaded service file %s", file)

			fmt.Printf("Starting service: %s\n", file)
			err = sc.Start()
			if err != nil {
				fmt.Printf("Loading file: %s - FAILED\n", file)
			}
		}
	}

	return err
}


// Create a service by method defined by a *CreateEntry structure.
func (me *Daemon) createEntry(c ServiceConfig) (*ServiceConfig, error) {

	var err error
	var sc *ServiceConfig

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		sc = &c


		// Basic sanity checks.
		if sc.MdnsType == "" {
			err = me.EntityId.ProduceError("service MdnsType not defined")
			break
		}

		if sc.Config.Name == "" {
			err = me.EntityId.ProduceError("service Name not defined")
			break
		}

		if sc.Config.DisplayName == "" {
			err = me.EntityId.ProduceError("service DisplayName not defined")
			break
		}

		if sc.Config.Description == "" {
			err = me.EntityId.ProduceError("service Description not defined")
			break
		}

		if sc.Config.Executable == "" {
			err = me.EntityId.ProduceError("service Executable not defined")
			break
		}


		sc.autoHost = "0.0.0.0"
		sc.autoPort = "0"
		// Check URL.
		//sc.Url = me.ParsePaths(*sc, sc.Url)
		sc.UrlPtr, err = network.ParseUrl(sc.Url)
		sc.Url = sc.UrlPtr.String()
		sc.autoHost = sc.UrlPtr.Hostname()
		sc.autoPort = sc.UrlPtr.Port()

		//sc.Host = sc.Url.Hostname()
		//sc.Port = sc.Url.Port()
		//if len(sc.Dependencies) == 0 {
		//	sc.Dependencies = []string{}
		//}


		// Parse paths.
		dirs := ospaths.NewPath()
		sc.Config.ChRoot = me.ParsePaths(*sc, sc.Config.ChRoot)
		dirs = dirs.AppendDir(sc.Config.ChRoot)


		sc.Config.Executable = me.ParsePaths(*sc, sc.Config.Executable)
		dirs = dirs.AppendFile(sc.Config.Executable)
		_, err = ospaths.FileExists(sc.Config.Executable)
		if err != nil {
			break
		}
		_, err = ospaths.FileSetExecutePerms(sc.Config.Executable)
		if err != nil {
			break
		}


		if sc.Config.WorkingDirectory == "" {
			sc.Config.WorkingDirectory = me.OsPaths.EventBrokerWorkingDir.String()
		} else {
			sc.Config.WorkingDirectory = me.ParsePaths(*sc, sc.Config.WorkingDirectory)
		}
		dirs = dirs.AppendDir(sc.Config.WorkingDirectory)


		if sc.Stdout == "" {
			sc.Stdout = me.OsPaths.EventBrokerLogDir.AddFileToPath("%s-error.log", sc.Name).String()
		} else {
			sc.Stdout = me.ParsePaths(*sc, sc.Stdout)
		}
		dirs = dirs.AppendFile(sc.Stdout)


		if sc.Stderr == "" {
			sc.Stderr = me.OsPaths.EventBrokerLogDir.AddFileToPath("%s.log", sc.Name).String()
		} else {
			sc.Stderr = me.ParsePaths(*sc, sc.Stderr)
		}
		dirs = dirs.AppendFile(sc.Stderr)


		err = dirs.CreateIfNotExists()
		if err != nil {
			break
		}


		// Parse envs and arguments.
		for k, v := range sc.Env {
			sc.Env[k] = me.ParsePaths(*sc, v)
		}

		for k, v := range sc.Arguments {
			sc.Arguments[k] = me.ParsePaths(*sc, v)
		}
	}

	return sc, err
}

