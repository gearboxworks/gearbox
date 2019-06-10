package daemon

import (
	"fmt"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/only"
	"os"
	"path/filepath"
	"strings"
	"time"
)


////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

// Unregister a service by method defined by a UUID reference.
func (me *Daemon) UnregisterByEntityId(u messages.MessageAddress) error {

	var err error
	var state states.Status

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = me.EnsureDaemonNotNil(u)
		if err != nil {
			break
		}

		for range only.Once {
			me.daemons[u].State.SetNewAction(states.ActionUnregister)	// Managed by Mutex
			channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)

			state, err = me.daemons[u].Status()	// Mutex not required
			if err != nil {
				continue
			}
			switch state.Current {
				case states.StateUnknown:
					//

				case states.StateStarted:
					err = me.daemons[u].instance.service.Stop()	// Mutex not required
					if err != nil {
						break
					}

				case states.StateStopped:
					//
			}

			err = me.daemons[u].instance.service.Uninstall()	// Mutex not required
			if err != nil {
				break
			}

			me.DeleteEntity(u)

			me.State.SetNewState(states.StateUnregistered, err)
			eblog.Debug(me.EntityId, "unregistered service %s OK", u.String())
		}
	}

	channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


// Unregister a service via a channel defined by a UUID reference.
func (me *Daemon) UnregisterByChannel(caller messages.MessageAddress, u messages.MessageAddress) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		//unreg := me.EntityId.Construct(me.EntityId, states.ActionUnregister, messages.MessageText(u.String()))
		unreg := caller.ConstructMessage(me.EntityId, states.ActionUnregister, messages.MessageText(u.String()))
		err = me.Channels.Publish(unreg)
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "unregistered service by channel %s OK", u.String())
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


// Unregister a service by method defined by a *CreateEntry structure.
func (me *Daemon) UnregisterByFile(f string) (*Service, error) {

	var err error
	var sc *ServiceConfig
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

		var check *Service
		check, err = me.FindExistingConfig(*sc)
		if check == nil {
			break
		}

		//if check.IsRegistered() {
		//	break
		//}

		err = me.UnregisterByEntityId(check.EntityId)
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "unregistered service by file %s OK", f)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return s, err
}


func (me *Daemon) UnLoadFiles() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for range only.Once {
			checkIn := string(me.osSupport.GetAdminRootDir() + "/" + DefaultJsonFiles)
			fmt.Printf("%d Unloading files... from %s\n", time.Now().Unix(), checkIn)

			var files []string
			err = filepath.Walk(checkIn, func(path string, info os.FileInfo, err error) error {
				files = append(files, path)
				return nil
			})
			if err != nil {
				break
			}

			for _, file := range files {
				if strings.HasSuffix(file, ".json") {
					var sc *Service
					fmt.Printf("Unloading file: %s\n", file)
					sc, err = me.UnregisterByFile(file)
					if sc == nil {
						eblog.Debug(me.EntityId, "Unloading file: %s - already unloaded\n", file)
						continue
					}
					if err != nil {
						eblog.Debug(me.EntityId, "Unloading file: %s - FAILED: %v\n", file, err)
						continue
					}
				}
			}
		}
	}

	return err
}

