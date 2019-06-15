package mqttBroker

import (
	"gearbox/box"
	"gearbox/eventbroker/channels"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/states"
	"gearbox/eventbroker/only"
	oss "gearbox/os_support"
	"github.com/jinzhu/copier"
)


func (me *MqttBroker) New(OsSupport oss.OsSupporter, args ...Args) error {

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
		_args.State = states.New(&_args.EntityId, &_args.EntityId, entity.SelfEntityName)

		//if _args.Servers == nil {
		//	_args.Servers, err = url.Parse(DefaultServer)
		//}

		//if _args.waitTime == 0 {
		//	_args.waitTime = defaultWaitTime
		//}

		*me = MqttBroker(_args)


		me.State.SetWant(states.StateIdle)
		me.State.SetNewState(states.StateIdle, err)
		eblog.Debug(me.EntityId, "init complete")
	}

	channels.PublishState(me.Channels, &me.EntityId, &me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


// Start the MQTT handler.
func (me *MqttBroker) StartHandler() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.State.SetNewState(states.StateStarting, err)
		channels.PublishState(me.Channels, &me.EntityId, &me.State)

		// Not using tasks.
		//for range only.Once {
		//	me.Task, err = tasks.StartTask(initMqttBroker, startMqttBroker, monitorMqttBroker, stopMqttBroker, me)
		//	me.State.SetError(err)
		//	if err != nil {
		//		break
		//	}
		//}

		//s, err = me.RegisterByFile("/Users/mick/.gearbox/admin/dist/eventbroker/unfsd/unfsd.json")


		me.State.SetNewState(states.StateStarted, err)
		eblog.Debug(me.EntityId, "started task handler")
	}

	channels.PublishState(me.Channels, &me.EntityId, &me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


// Stop the MQTT handler.
func (me *MqttBroker) StopHandler() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.State.SetNewState(states.StateStopping, err)
		channels.PublishState(me.Channels, &me.EntityId, &me.State)

		for range only.Once {
			err = me.Task.Stop()
		}

		me.State.SetNewState(states.StateStopped, err)
		eblog.Debug(me.EntityId, "stopped task handler")
	}

	channels.PublishState(me.Channels, &me.EntityId, &me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


//func (me *MqttBroker) StopServices() error {
//
//	var err error
//
//	for range only.Once {
//		err = me.EnsureNotNil()
//		if err != nil {
//			break
//		}
//
//		for u, _ := range me.services {
//			if me.services[u].IsManaged {
//				_ = me.UnsubscribeByUuid(u)
//				// Ignore error, will clean up when program exits.
//			}
//		}
//	}
//
//	channels.PublishState(me.Channels, &me.EntityId, &me.State)
//	eblog.LogIfNil(me, err)
//	eblog.LogIfError(me.EntityId, err)
//
//	return err
//}
//
