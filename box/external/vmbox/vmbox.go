package vmbox

import (
	"errors"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/states"
	"gearbox/eventbroker/tasks"
	"gearbox/global"
	"github.com/gearboxworks/go-status/only"
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
		//	err = msgs.MakeError(_args.EntityId,"unable to copy config args")
		//	break
		//}

		if _args.Channels == nil {
			err = msgs.MakeError(_args.EntityId, "channel pointer is nil")
			break
		}

		if _args.OsPaths == nil {
			err = msgs.MakeError(_args.EntityId, "ospaths is nil")
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

		_args.State = states.New(_args.EntityId, _args.EntityId, entity.SelfEntityName)

		if _args.Boxname == "" {
			_args.Boxname = global.Brandname
		}

		if _args.waitTime == 0 {
			_args.waitTime = DefaultVmWaitTime
		}

		_args.Releases, _ = NewReleases(_args.Channels)

		_args.vms = make(VmMap)

		*me = VmBox(_args)

		me.SetWantState(states.StateIdle)
		me.SetStateError(err)
		me.SetState(states.StateIdle)
		eblog.Debug(me.EntityId, "init complete")
	}

	me.PublishChannelState(me.State)
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

		me.SetStateAction(states.ActionStart)
		me.PublishChannelState(me.State)

		for range only.Once {
			me.Task, err = tasks.StartTask(initVmBox, startVmBox, monitorVmBox, stopVmBox, me)
			me.SetStateError(err)
			if err != nil {
				eblog.LogIfError(msgs.Address("unable to start tasks"), err)
				break
			}
		}

		if err == nil {
			me.SetState(states.StateStarted)
		}

		me.PublishChannelState(me.State)
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

		me.SetStateAction(states.ActionStop)
		me.PublishChannelState(me.State)

		err = me.StopVms()
		if err != nil {
			eblog.LogIfError(msgs.Address("VM failed to stop"), err)
		}

		err = me.Task.Stop()
		if err != nil {
			me.SetStateError(err)
		} else {
			me.SetState(states.StateStopped)
		}

		me.PublishChannelState(me.State)
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

		for u := range me.vms {
			if !me.vms[u].IsManaged {
				continue
			}
			err = me.vms[u].Stop()
			eblog.LogIfError(msgs.Address("VM failed to stop"), err)
		}
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}

func (me *VmBox) EnsureNotNil() (err error) {
	if me == nil {
		err = errors.New("VmBox instance is nil")
	}
	return err
}

func (me *VmBox) SetStateError(err error) {
	me.State.Error = err
}

func (me *VmBox) SetState(s states.State) {
	me.State.SetState(s)
}

func (me *VmBox) SetStateAction(a states.Action) {
	me.State.SetNewAction(a)
}

func (me *VmBox) PublishChannelState(state *states.Status) {
	me.Channels.PublishState(me.State)
}

func (me *VmBox) SetWantState(s states.State) {
	me.State.SetWant(s)
}
