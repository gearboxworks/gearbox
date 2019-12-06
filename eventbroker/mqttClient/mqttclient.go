package mqttClient

import (
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/states"
	"gearbox/eventbroker/tasks"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gearboxworks/go-status/only"
	"net/url"
)

func (me *MqttClient) New(args ...Args) error {

	var _args Args
	var err error

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		if _args.Channels == nil {
			err = msgs.MakeError(me.EntityId, "channel pointer is nil")
			break
		}

		if _args.OsPaths == nil {
			err = msgs.MakeError(me.EntityId, "ospaths is nil")
			break
		}

		if _args.EntityId == "" {
			_args.EntityId = entity.MqttClientEntityName
		}
		_args.State = states.New(_args.EntityId, _args.EntityId, entity.SelfEntityName)

		if _args.Boxname == "" {
			_args.Boxname = entity.MqttClientEntityName
		}

		if _args.Server == nil {
			_args.Server, err = url.Parse(DefaultServer)
		}

		if _args.waitTime == 0 {
			_args.waitTime = defaultWaitTime
		}

		_args.instance.options = mqtt.NewClientOptions()
		if _args.instance.options == nil {
			err = msgs.MakeError(me.EntityId, "unable to create options")
			break
		}
		_args.instance.options.SetClientID(_args.EntityId.String())

		_args.services = make(ServicesMap)

		*me = MqttClient(_args)

		me.State.SetWant(states.StateIdle)
		me.State.SetNewState(states.StateIdle, err)
		eblog.Debug(me.EntityId, "init complete")
	}

	me.Channels.PublishState(me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return err
}

// Start the MQTT handler.
func (me *MqttClient) StartHandler() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.State.SetNewAction(states.ActionStart)
		me.Channels.PublishState(me.State)

		for range only.Once {
			me.Task, err = tasks.StartTask(initMqttClient, startMqttClient, monitorMqttClient, stopMqttClient, me)
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
	eblog.LogIfError(err)

	return err
}

// Stop the MQTT handler.
func (me *MqttClient) StopHandler() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.State.SetNewAction(states.ActionStop)
		me.Channels.PublishState(me.State)

		for range only.Once {
			_ = me.StopServices()
			// Ignore error, will clean up when program exits.

			err = me.Task.Stop()
		}

		me.State.SetNewState(states.StateStopped, err)
		me.Channels.PublishState(me.State)
		eblog.Debug(me.EntityId, "stopped task handler")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return err
}

func (me *MqttClient) StopServices() error {

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

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return err
}
