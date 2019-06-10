package eventbroker

import (
	"fmt"
	"gearbox/box"
	"gearbox/global"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/daemon"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/mqttClient"
	"gearbox/heartbeat/eventbroker/network"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/only"
	oss "gearbox/os_support"
	"github.com/jinzhu/copier"
	"path/filepath"
	"time"
)


func New(OsSupport oss.OsSupporter, args ...Args) (*EventBroker, error) {

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

		_args.osSupport = OsSupport
		foo := box.Args{}
		err = copier.Copy(&foo, &_args)
		if err != nil {
			err = me.EntityId.ProduceError("unable to copy config args")
			break
		}

		_args.PidFile = filepath.FromSlash(fmt.Sprintf("%s/%s", _args.osSupport.GetAdminRootDir(), defaultPidFile))

		_args.Entities = make(Entities)


		// 1. Channel - provides inter-thread communications.
		err = _args.Channels.New(OsSupport, channels.Args{})
		if err != nil {
			break
		}


		// 2. ZeroConf - provides discovery and management of network services.
		err = _args.ZeroConf.New(OsSupport, network.Args{Channels: &_args.Channels})
		if err != nil {
			break
		}


		// 3. Daemon - provides control over arbitrary services.
		err = _args.Daemon.New(OsSupport, daemon.Args{Channels: &_args.Channels})
		if err != nil {
			break
		}

		// 4. MQTT broker - provides inter-process communications.


		// 5. MQTT client - provides inter-process communications.
		err = _args.MqttClient.New(OsSupport, mqttClient.Args{Channels: &_args.Channels})
		if err != nil {
			break
		}

		*me = EventBroker(_args)


		me.State.SetWant(states.StateIdle)
		if me.State.SetNewState(states.StateIdle, err) {
			eblog.Debug(me.EntityId, "init complete")
		}
	}

	channels.PublishCallerState(&me.Channels, &me.EntityId, &me.State)
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


		// 1. Channel - provides inter-thread communications.
		// Start the inter-thread service.
		// Note: These will be started dynamically as clients are registered.
		me.channelHandler, err = me.Channels.StartClientHandler(messages.BroadcastAddress)
		if err != nil {
			break
		}
		_, err := me.Channels.Subscribe(messages.MessageTopic{
			Address: messages.BroadcastAddress,
			SubTopic: "status",
		}, statusHandler, me, states.InterfaceTypeStatus)
		if err != nil {
			break
		}


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
		//_ = me.Daemon.LoadFiles()


		// 5. MQTT client - start the inter-process communications.
		err = me.MqttClient.StartHandler()
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "event broker started OK")
	}

	channels.PublishCallerState(&me.Channels, &me.EntityId, &me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *EventBroker) Stop() error {
	var err error

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


func (me *EventBroker) TempLoop() error {

	var err error

	fmt.Printf("(me *EventBroker) TempLoop()\n")

	//msg := messages.Message{
	//	Source: me.EntityId,
	//	Topic: messages.MessageTopic{
	//		Address: "eventbrokerdaemon",
	//		SubTopic: "get",
	//	},
	//	Text: "status",
	//}
	msg := messages.Message{
		Source: me.EntityId,
		Topic: messages.MessageTopic{
			Address: "eventbroker-daemon",
			SubTopic: "status",
		},
		Text: "now",
	}
	//time.Sleep(time.Second * 8)
	//err = me.Daemon.LoadFiles()
	//fmt.Printf("me.Daemon.LoadFiles(): %v\n", err)
	//me.CreateEntity("BEEP")
	me.Foo()
	me.SimpleLoop()

	time.Sleep(time.Hour * 8000)

	err = me.Daemon.LoadFiles()
	fmt.Printf("me.Daemon.LoadFiles(): %v\n", err)

	err = me.Daemon.LoadFiles()
	fmt.Printf("me.Daemon.LoadFiles(): %v\n", err)

	//time.Sleep(time.Hour * 8000)

	go func() {
		index := 0
		for {
			fmt.Printf("####################################################################\nPING\n")
			//me.Daemon.ChannelHandler.List()
			//me.Daemon.Channels.ListSubscribers()

			//fmt.Printf("Error1: %v\n", me.Daemon.State.Error)
			//me.Daemon.State.SetError(errors.New(fmt.Sprintf("Loop #%d (%s)", index, me.Daemon.Fluff)))
			//fmt.Printf("Error2: %v\n", me.Daemon.State.Error)

			fmt.Printf("\n\n%d gbevents before: %v\n", time.Now().Unix(), me.Daemon.State.GetError())
			i, _ := me.Channels.PublishAndWaitForReturn(msg, 400)
			//foo := reflect.ValueOf(i)
			//fmt.Printf("ERROR: %v\t\tRESPONSE: %v\n", err, i)
			//fmt.Printf("Reflect %s, %s\n", foo.Type(), foo.String())
			f, err := states.InterfaceToTypeStatus(i)
			if err == nil {
				fmt.Printf("%d gbevents after: %v\n", time.Now().Unix(), f.GetError())
				// fmt.Printf("%d gbevents after: %v (%v)\n", time.Now().Unix(), f.GetError(), me.Daemon.Fluff)
			} else {
				fmt.Printf("%d gbevents after: is nil!\n", time.Now().Unix())
			}

			index++
			//me.Daemon.State.Error = nil
			time.Sleep(time.Second * 20)
		}
	}()

	//time.Sleep(time.Second * 6)
	//
	////me.CreateEntity("HELO")
	//var u *url.URL
	//u, err = me.FindMqttBroker()
	//if err == nil {
	//	err = me.MqttClient.ConnectToServer(u.String())
	//}
	//if err != nil {
	//	eblog.Debug(me.EntityId, "Aaaaargh! => %v", err)
	//}
	//
	//_ = me.MqttClient.GlobSubscribe(me.MqttClient.EntityId)
	//
	//time.Sleep(time.Second * 200)
	//
	//err = me.Stop()
	//
	//fmt.Printf("Exiting...\n")

	return err
}


func (me *EventBroker) Foo() {

	var st messages.SubTopics
	var t messages.Topics
	var ma messages.MessageAddresses

	msg := messages.Message{
		Source: me.EntityId,
		Topic: messages.MessageTopic{
			Address: "eventbroker-daemon",
			SubTopic: "get",
		},
		Text: "topics",
	}
	fmt.Printf("\n\n%d gbevents before: %v\n", time.Now().Unix(), me.Daemon.State.GetError())
	i, _ := me.Channels.PublishAndWaitForReturn(msg, 400)
	f, err := messages.InterfaceToTypeSubTopics(i)
	if err == nil {
		fmt.Printf("%d gbevents after: %v\n", time.Now().Unix(), f)
	} else {
		fmt.Printf("%d gbevents after: is nil!\n", time.Now().Unix())
	}

	//time.Sleep(time.Hour * 8000)

	fmt.Printf("\n### Channels\n")
	fmt.Printf("me.Channels.GetManagedEntities\n")
	ma = me.Channels.GetManagedEntities()
	for _, f := range ma {
		fmt.Printf("me.Channels.GetManagedEntities\t=> %s\n", f.String())
	}
	fmt.Printf("me.Channels.GetEntities\n")
	ma = me.Channels.GetEntities()
	for _, f := range ma {
		fmt.Printf("me.Channels.GetEntities\t=> %s\n", f.String())
	}
	fmt.Printf("me.Channels.GetListenerTopics\n")
	t = me.Channels.GetListenerTopics()
	for _, f := range t {
		fmt.Printf("me.Channels.GetListenerTopics\t=> %s\n", f.String())
	}

	//fmt.Printf("\n### eventbroker\n")
	//fmt.Printf("me.GetEntities\n")
	//ma = me.GetEntities()
	//for _, f := range ma {
	//	fmt.Printf("me.GetEntities => %s\n", f.String())
	//}
	//
	//fmt.Printf("me.GetManagedEntities\n")
	//ma = me.Daemon.GetManagedEntities()
	//for _, f := range ma {
	//	fmt.Printf("me.GetManagedEntities => %s\n", f.String())
	//}
	//
	//fmt.Printf("me.GetTopics\n")
	//st = me.Daemon.GetTopics()
	//for _, f := range st {
	//	fmt.Printf("me.GetTopics => %s\n", f.String())
	//}

	fmt.Printf("\n### ZeroConf\n")
	fmt.Printf("me.ZeroConf.GetEntities\n")
	ma = me.ZeroConf.GetEntities()
	for _, f := range ma {
		fmt.Printf("me.ZeroConf.GetEntities\t=> %s\n", f.String())
	}

	fmt.Printf("me.ZeroConf.GetManagedEntities\n")
	ma = me.ZeroConf.GetManagedEntities()
	for _, f := range ma {
		fmt.Printf("me.ZeroConf.GetManagedEntities\t=> %s\n", f.String())
	}

	fmt.Printf("me.ZeroConf.GetTopics\n")
	st = me.ZeroConf.GetTopics()
	for _, f := range st {
		fmt.Printf("me.ZeroConf.GetTopics\t=> %s\n", f.String())
	}


	fmt.Printf("\n### Daemon\n")
	fmt.Printf("me.Daemon.GetEntities\n")
	ma = me.Daemon.GetEntities()
	for _, f := range ma {
		fmt.Printf("me.Daemon.GetEntities\t=> %s\n", f.String())
	}

	fmt.Printf("me.Daemon.GetManagedEntities\n")
	ma = me.Daemon.GetManagedEntities()
	for _, f := range ma {
		fmt.Printf("me.Daemon.GetManagedEntities\t=> %s\n", f.String())
	}

	fmt.Printf("me.Daemon.GetTopics\n")
	st = me.Daemon.GetTopics()
	for _, f := range st {
		fmt.Printf("me.Daemon.GetTopics\t=> %s\n", f.String())
	}


	fmt.Printf("\n### MqttClient\n")
	fmt.Printf("me.MqttClient.GetEntities\n")
	ma = me.MqttClient.GetEntities()
	for _, f := range ma {
		fmt.Printf("me.MqttClient.GetEntities\t=> %s\n", f.String())
	}

	fmt.Printf("me.MqttClient.GetManagedEntities\n")
	ma = me.MqttClient.GetManagedEntities()
	for _, f := range ma {
		fmt.Printf("me.MqttClient.GetManagedEntities\t=> %s\n", f.String())
	}

	fmt.Printf("me.MqttClient.GetTopics\n")
	st = me.MqttClient.GetTopics()
	for _, f := range st {
		fmt.Printf("me.MqttClient.GetTopics\t=> %s\n", f.String())
	}
}


func (me *EventBroker) SimpleLoop() {

	var err error
	fmt.Printf("SimpleLoop()\n")

	for range only.Once {
		//fmt.Printf("daemon == %v\n", d)
		//fmt.Printf("daemon topics == %v\n", d.GetTopics())
		time.Sleep(time.Hour * 100)

		msg := messages.Message{
			Source: me.EntityId,
			Topic: messages.MessageTopic{
				Address: "eventbroker-daemon",
				SubTopic: "status",
			},
			Text: "hello",
		}

		for i := 0; i < 10; i++ {

			err = me.Daemon.LoadFiles()
			fmt.Printf("me.Daemon.LoadFiles(): %v\n", err)
			time.Sleep(time.Second * 10)

			fmt.Printf("\n\n%d gbevents before: %v\n", time.Now().Unix(), me.Daemon.State.GetError())
			i, _ := me.Channels.PublishAndWaitForReturn(msg, 400)
			f, err := messages.InterfaceToTypeSubTopics(i)
			if err == nil {
				fmt.Printf("%d gbevents after: %v\n", time.Now().Unix(), f)
			} else {
				fmt.Printf("%d gbevents after: is nil!\n", time.Now().Unix())
			}


			err = me.Daemon.UnLoadFiles()
			fmt.Printf("me.Daemon.UnLoadFiles(): %v\n", err)

			time.Sleep(time.Second * 20)
		}

		time.Sleep(time.Hour * 60)
	}

}


//
//
//// Non-exposed channel function that responds to a "register" channel request.
//func registerService(event *messages.Message, i channels.Argument, r channels.ReturnType) channels.Return {
//
//	var me *Daemon
//	var sc *Service
//	var err error
//
//	for range only.Once {
//		me, err = InterfaceToTypeDaemon(i)
//		if err != nil {
//			break
//		}
//
//		//fmt.Printf("Rx: %v\n", event)
//
//		ce := ServiceConfig{}
//		err = json.Unmarshal(event.Text.ByteArray(), &ce)
//		if err != nil {
//			break
//		}
//
//		sc, err = me.Register(ce)
//		if err != nil {
//			break
//		}
//
//		eblog.Debug(me.EntityId, "registered service by channel %s OK", sc.EntityId.String())
//	}
//
//	eblog.LogIfNil(me, err)
//	eblog.LogIfError(me.EntityId, err)
//
//	return sc
//}
//
//
//// Non-exposed channel function that responds to an "unregister" channel request.
//func unregisterService(event *messages.Message, i channels.Argument, r channels.ReturnType) channels.Return {
//
//	var me *Daemon
//	var err error
//
//	for range only.Once {
//		me, err = InterfaceToTypeDaemon(i)
//		if err != nil {
//			break
//		}
//
//		//fmt.Printf("MESSAGE Rx:\n[%v]\n", event.Text.String())
//
//		// Use message element as the UUID.
//		err = me.UnregisterByEntityId(event.Text.ToMessageAddress())
//		if err != nil {
//			break
//		}
//
//		eblog.Debug(me.EntityId, "unregistered service by channel %s OK", event.Text.ToMessageAddress())
//	}
//
//	eblog.LogIfNil(me, err)
//	eblog.LogIfError(me.EntityId, err)
//
//	return err
//}


// Non-exposed channel function that responds to a "status" channel request.
func statusHandler(event *messages.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *EventBroker
	var ret *states.Status

	for range only.Once {
		me, err = InterfaceToTypeEventBroker(i)
		if err != nil {
			break
		}

		fmt.Printf("Event(%s) Time:%d Src:%s Text:%s\n", event.Topic.String(), event.Time.Convert().Unix(), event.Source.String(), event.Text.String())

		if _, ok := me.Entities[event.Source]; !ok {
			sub := Entity{
				Src:   event.Topic.Address,
				State: &states.Status{},
				StateString: states.State(event.Text),
			}
			me.Entities[event.Source] = &sub
		}

		fmt.Printf("Entities\n")
		for n, e := range me.Entities {
			fmt.Printf("%s	%s %s %v\n",
				n,
				e.Src,
				e.StateString,
				e.State,
				)
		}

		switch event.Text {
			case states.StateStarted:
				// .

			case states.StateStopped:
				// .
		}

		//msg := messages.Message{
		//	Source: me.EntityId,
		//	Topic: messages.MessageTopic{
		//		Address: event.Source,
		//		SubTopic: "status",
		//	},
		//	Text: "",
		//}
		//fmt.Printf("\n\n%d gbevents before: %v\n", time.Now().Unix(), me.Daemon.State.GetError())
		//i, _ := me.Channels.PublishAndWaitForReturn(msg, 400)
		//f, err := states.InterfaceToTypeStatus(i)
		//if err == nil {
		//	fmt.Printf("%d gbevents after: %v\n", time.Now().Unix(), f)
		//} else {
		//	fmt.Printf("%d gbevents after: is nil!\n", time.Now().Unix())
		//}

		eblog.Debug(me.EntityId, "statusHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return ret
}

