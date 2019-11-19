package mqttBroker

import (
	"gearbox/box"
	"gearbox/eventbroker/channels"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/states"
	"gearbox/eventbroker/tasks"
	"github.com/gearboxworks/go-status/only"
	//	oss "gearbox/os_support"
	"github.com/jinzhu/copier"
)


func (me *mqttBroker) New(OsBridge osbridge.OsBridger, args ...Args) error {

	var _args Args
	var err error

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		_args.osSupport = OsBridge
		foo := box.Args{}
		err = copier.Copy(&foo, &_args)
		if err != nil {
			err = me.EntityId.ProduceError("unable to copy config args")
			break
		}

		if _args.EntityId == "" {
			_args.EntityId = DefaultEntityId
		}

		//if _args.Servers == nil {
		//	_args.Servers, err = url.Parse(DefaultServer)
		//}

		if _args.waitTime == 0 {
			_args.waitTime = defaultWaitTime
		}

		*me = mqttBroker(_args)


		eblog.Debug(me.EntityId, "init complete")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


// Start the MQTT handler.
func (me *mqttBroker) StartHandler() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.State.SetNewState(states.StateStarting, err)
		channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)

		for range only.Once {
			me.Task, err = tasks.StartTask(initMqttClient, startMqttClient, monitorMqttClient, stopMqttClient, me)
			me.State.SetError(err)
			if err != nil {
				break
			}
		}

		me.State.SetNewState(states.StateStarted, err)
		eblog.Debug(me.EntityId, "started task handler")
	}

	channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


// Stop the MQTT handler.
func (me *mqttBroker) StopHandler() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.State.SetNewState(states.StateStopping, err)
		channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)

		for range only.Once {
			_ = me.StopServices()
			// Ignore error, will clean up when program exits.

			err = me.Task.Stop()
		}

		me.State.SetNewState(states.StateStopped, err)
		eblog.Debug(me.EntityId, "stopped task handler")
	}

	channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *mqttBroker) StopServices() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for u := range me.services {
			if me.services[u].IsManaged {
				_ = me.UnsubscribeByUuid(u)
				// Ignore error, will clean up when program exits.
			}
		}
	}

	channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}

