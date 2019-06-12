package daemon

import (
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/only"
	"gearbox/heartbeat/eventbroker/states"
)


////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

// Unregister a service by method defined by a UUID reference.
func (me *Daemon) UnregisterByEntityId(client messages.MessageAddress) error {

	var err error
	var state states.Status

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = client.EnsureNotNil()
		if err != nil {
			break
		}

		me.daemons[client].State.SetNewAction(states.ActionUnregister)	// Managed by Mutex
		me.daemons[client].channels.PublishCallerState(&me.daemons[client].State)

		state, err = me.daemons[client].Status(DontPublishState)	// Mutex not required
		if err != nil {
			continue
		}
		switch state.Current {
			case states.StateUnknown:
				//

			case states.StateStarted:
				err = me.daemons[client].instance.service.Stop()	// Mutex not required
				if err != nil {
					break
				}

			case states.StateStopped:
				//
		}

		err = me.daemons[client].instance.service.Uninstall()	// Mutex not required
		if err != nil {
			break
		}

		err = me.DeleteEntity(client)
		if err != nil {
			break
		}

		me.Channels.PublishSpecificState(&client, states.State(states.StateUnsubscribed))
		eblog.Debug(me.EntityId, "unregistered service %s OK", client.String())
	}

	me.Channels.PublishState(&me.EntityId, &me.State)
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

	me.Channels.PublishState(&me.EntityId, &me.State)
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


func (me *Daemon) UnloadServiceFiles() error {

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
			sc, err = me.UnregisterByFile(file)
			if sc == nil {
				//eblog.Debug(me.EntityId, "Unloaded service file %s", file)
				continue
			}
			if err != nil {
				eblog.Debug(me.EntityId, "Unloading service file %s failed with '%v'\n", file, err)
				continue
			}
			eblog.Debug(me.EntityId, "Unloaded service file %s", file)
		}
	}

	return err
}

