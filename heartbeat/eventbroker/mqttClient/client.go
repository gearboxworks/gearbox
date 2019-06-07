package mqttClient

import (
	"gearbox/box"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/heartbeat/eventbroker/tasks"
	"gearbox/only"
	oss "gearbox/os_support"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jinzhu/copier"
)


func (me *MqttClient) New(OsSupport oss.OsSupporter, args ...Args) error {

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

		//if _args.Servers == nil {
		//	_args.Servers, err = url.Parse(DefaultServer)
		//}

		if _args.waitTime == 0 {
			_args.waitTime = defaultWaitTime
		}

		_args.instance.options = mqtt.NewClientOptions()
		if _args.instance.options == nil {
			err = me.EntityId.ProduceError("unable to create options")
			break
		}
		_args.instance.options.SetClientID(_args.EntityId.String())

		_args.services = make(ServicesMap)

		*me = MqttClient(_args)


		me.State.SetNewWantState(states.StateIdle)
		if me.State.SetNewState(states.StateIdle, err) {
			eblog.Debug(me.EntityId, "init complete")
		}
	}

	channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
	eblog.LogIfError(me, err)

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

		me.Task, me.State.Error = tasks.StartTask(initMqttClient, startMqttClient, monitorMqttClient, stopMqttClient, me)
		if me.State.Error != nil {
			break
		}

		eblog.Debug(me.EntityId, "started task handler")
	}

	channels.PublishCallerState(me.channels, &me.EntityId, &me.State)
	eblog.LogIfError(me, err)

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

		for u, _ := range me.services {
			eblog.Debug(me.EntityId, "unsubscribing %s", u.String())
			err = me.UnsubscribeByUuid(u)
			if err != nil {
				eblog.Debug(me.EntityId, "ERROR failed to unsubscribe %s", u.String())
				break
			}
		}

		err = me.Task.Stop()
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "stopped task handler")
	}

	channels.PublishCallerState(me.channels, &me.EntityId, &me.State)
	eblog.LogIfError(me, err)

	return err
}

