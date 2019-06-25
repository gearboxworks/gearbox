package vmbox

import (
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/only"
	"gearbox/eventbroker/states"
	"gearbox/eventbroker/tasks"
	"gearbox/global"
)


func New(args ...Args) (*VmBox, error) {

	var _args Args
	var err error

	me := &VmBox{}

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		//foo := Args{}
		//err = copier.Copy(&foo, &_args)
		//if err != nil {
		//	err = _args.EntityId.ProduceError("unable to copy config args")
		//	break
		//}

		if _args.Channels == nil {
			err = _args.EntityId.ProduceError("channel pointer is nil")
			break
		}

		if _args.OsPaths == nil {
			err = _args.EntityId.ProduceError("ospaths is nil")
			break
		}


		if _args.EntityId == "" {
			_args.EntityId = entity.VmBoxEntityName
		}

		if _args.EntityName == "" {
			_args.EntityName = _args.EntityId
		}

		if _args.EntityParent == "" {
			_args.EntityParent = _args.EntityId
		}

		_args.State = states.New(&_args.EntityId, &_args.EntityId, entity.SelfEntityName)

		if _args.Boxname == "" {
			_args.Boxname = global.Brandname
		}

		if _args.waitTime == 0 {
			_args.waitTime = DefaultVmWaitTime
		}

		_args.Releases, _ = NewReleases(_args.Channels)

		_args.vms = make(VmMap)

		*me = VmBox(_args)


		me.State.SetWant(states.StateIdle)
		me.State.SetNewState(states.StateIdle, err)
		eblog.Debug(me.EntityId, "init complete")
	}

	me.Channels.PublishState(me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return me, err
}


// Start the VmBox handler.
func (me *VmBox) Start() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.State.SetNewAction(states.ActionStart)
		me.Channels.PublishState(me.State)

		for range only.Once {
			me.Task, err = tasks.StartTask(initVmBox, startVmBox, monitorVmBox, stopVmBox, me)
			me.State.SetError(err)
			if err != nil {
				break
			}
		}

		me.State.SetNewState(states.StateStarted, err)
		me.Channels.PublishState(me.State)
		eblog.Debug(me.EntityId, "started task handler")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


// Stop the VmBox handler.
func (me *VmBox) Stop() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.State.SetNewAction(states.ActionStop)
		me.Channels.PublishState(me.State)

		for range only.Once {
			_ = me.StopVms()
			// Ignore error, will clean up when program exits.

			err = me.Task.Stop()
		}

		me.State.SetNewState(states.StateStopped, err)
		me.Channels.PublishState(me.State)
		eblog.Debug(me.EntityId, "stopped task handler")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *VmBox) StopVms() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for u, _ := range me.vms {
			if me.vms[u].IsManaged {
				_ = me.vms[u].Stop()
				// Ignore error, will clean up when program exits.
			}
		}
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}

