package eventbroker

import (
	"fmt"
	"gearbox/eventbroker/channels"
	"gearbox/eventbroker/daemon"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/mqttClient"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/network"
	"gearbox/eventbroker/osdirs"
	"gearbox/eventbroker/states"
	"gearbox/global"
	"github.com/gearboxworks/go-status/only"
	"time"
)

func New(args ...Args) (*EventBroker, error) {

	var _args Args
	var err error

	me := &EventBroker{}

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		if _args.Boxname == "" {
			_args.Boxname = global.Brandname
		}

		if _args.EntityId == "" {
			_args.EntityId = DefaultEntityName
		}
		_args.State = states.New(_args.EntityId, _args.EntityId, entity.SelfEntityName)

		if _args.SubBaseDir == "" {
			_args.SubBaseDir = osdirs.DefaultBaseDir
		}
		_args.OsDirs = osdirs.New(_args.SubBaseDir)
		err = _args.OsDirs.CreateIfNotExists()
		if err != nil {
			break
		}

		_args.Services = make(Services)

		*me = EventBroker(_args)

		// 0. Logger - enable logging.
		me.Logger, err = eblog.NewLogger(me.OsDirs)
		if err != nil {
			break
		}

		// 1. Channel - provides inter-thread communications.
		err = me.Channels.New(channels.Args{
			Boxname: me.Boxname,
			OsDirs:  me.OsDirs,
		})
		if err != nil {
			break
		}
		err = me.StartChannelHandler()
		if err != nil {
			break
		}

		// 2. ZeroConf - provides discovery and management of network services.
		err = me.ZeroConf.New(network.Args{
			Boxname:  me.Boxname,
			Channels: &me.Channels,
			OsPaths:  me.OsDirs,
		})
		if err != nil {
			break
		}

		// 3. Daemon - provides control over arbitrary services.
		err = me.Daemon.New(daemon.Args{
			Boxname:  me.Boxname,
			Channels: &me.Channels,
			BaseDirs: me.OsDirs,
		})
		if err != nil {
			break
		}

		// 4. MQTT broker - provides inter-process communications.

		// 5. MQTT client - provides inter-process communications.
		err = me.MqttClient.New(mqttClient.Args{
			Boxname:  me.Boxname,
			Channels: &me.Channels,
			OsPaths:  me.OsDirs,
		})
		if err != nil {
			break
		}

		me.State.SetWant(states.StateIdle)
		if me.State.SetNewState(states.StateIdle, err) {
			eblog.Debug(me.EntityId, "init complete")
		}
	}

	me.Channels.PublishState(me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return me, err
}

func (me *EventBroker) Start() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.State.SetNewAction(states.ActionStart)
		me.Channels.PublishState(me.State)

		// 1. Channel - provides inter-thread communications.
		// Start the inter-thread service.
		// Note: These will be started dynamically as clients are registered.

		// 2. ZeroConf - start discovery and management of network services.
		err = me.ZeroConf.StartHandler()
		if err != nil {
			break
		}

		// 3. Daemon - starts the daemon handler.
		err = me.Daemon.StartHandler()
		if err != nil {
			break
		}

		// 4. MQTT broker - start the inter-process communications.
		// This will be started via the daemons process.

		// 5. MQTT client - start the inter-process communications.
		err = me.MqttClient.StartHandler()
		if err != nil {
			break
		}

		me.State.SetNewState(states.StateStarted, err)
		me.Channels.PublishState(me.State)
		eblog.Debug(me.EntityId, "event broker started OK")
	}

	me.Channels.PublishState(me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}

func (me *EventBroker) Stop() error {

	var err error

	me.State.SetNewAction(states.ActionStop)
	me.Channels.PublishState(me.State)

	err = me.MqttClient.StopHandler()
	if err != nil {
		eblog.Debug(me.EntityId, "MqttClient shutdown error %v", err)
	}

	err = me.Daemon.StopHandler()
	if err != nil {
		eblog.Debug(me.EntityId, "Daemon shutdown error %v", err)
	}

	err = me.ZeroConf.StopHandler()
	if err != nil {
		eblog.Debug(me.EntityId, "ZeroConf shutdown error %v", err)
	}

	me.State.SetNewState(states.StateStopped, err)
	me.Channels.PublishState(me.State)

	// Add in wait thingy.
	time.Sleep(time.Second * 2)

	err = me.StopChannelHandler()
	if err != nil {
		eblog.Debug(me.EntityId, "Channels shutdown error %v", err)
	}

	err = me.Channels.StopHandler()
	if err != nil {
		eblog.Debug(me.EntityId, "Channels shutdown error %v", err)
	}

	return err
}

func (me *EventBroker) Restart() error {
	fmt.Printf("(me *EventBroker) RestartService() error\n")

	return nil
}

func (me *EventBroker) Status() error {
	fmt.Printf("(me *EventBroker) ServiceStatus() error\n")

	return nil
}

func (me *EventBroker) Create() error {
	fmt.Printf("(me *EventBroker) CreateService() error\n")

	return nil
}

func (me *EventBroker) StatusOf(client msgs.Address) (states.Status, error) {

	var ret states.Status
	var err error

	msg := msgs.Message{
		Source: me.EntityId,
		Topic: msgs.Topic{
			Address:  client,
			SubTopic: "status",
		},
		Text: "",
	}

	wrapper := msgs.Message{
		Source: me.EntityId,
		Topic: msgs.Topic{
			Address:  entity.BroadcastEntityName,
			SubTopic: "get",
		},
		Text: msg.ToMessageText(),
	}

	for range only.Once {
		i, err := me.Channels.PublishAndWaitForReturn(wrapper, 400)
		if err != nil {
			break
		}

		f, err := states.InterfaceToTypeStatus(i)
		if err != nil {
			break
		}

		ret = *f
	}

	return ret, err
}

func (me *EventBroker) GetSimpleStatus() (SimpleStatus, error) {

	ret := make(SimpleStatus)
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for _, v := range me.Services {
			v.mutex.RLock()
			ret[*v.State.EntityName] = v.State.Current
			v.mutex.RUnlock()
		}
	}

	return ret, err
}
