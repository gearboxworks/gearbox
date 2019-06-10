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

		// 2. ZeroConf - start discovery and management of network services.
		//err = me.ZeroConf.StartHandler()
		//if err != nil {
		//	break
		//}

		// 3. Daemon - starts the daemon handler.
		err = me.Daemon.StartHandler()
		if err != nil {
			break
		}

		// 4. MQTT broker - start the inter-process communications.
		// This will be started via the daemons process.

		// 5. MQTT client - start the inter-process communications.
		//err = me.MqttClient.StartHandler()
		//if err != nil {
		//	break
		//}

		eblog.Debug(me.EntityId, "event broker started OK")
	}

	channels.PublishCallerState(&me.Channels, &me.EntityId, &me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
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
	time.Sleep(time.Second * 8)

	err = me.Daemon.LoadFiles()
	fmt.Printf("me.Daemon.LoadFiles(): %v\n", err)

	//me.CreateEntity("BEEP")

	me.Daemon.Foo()

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


// Executed from a channel
func manageService(event *messages.Message, i interface{}) error {

	var err error
	var sc *network.Service

	for range only.Once {
		sc, err = network.InterfaceToTypeService(i)	// sc = rs.(*Service)
		if err != nil {
			break
		}
		fmt.Printf("DEBUG: %v\n", sc.EntityId)

		switch {
			case event.Topic.SubTopic.String() == "start":
				fmt.Printf("Start a physical process.\n")

			case event.Topic.SubTopic.String() == "stop":
				fmt.Printf("Stop a physical process.\n")

			case event.Topic.SubTopic.String() == "status":
				fmt.Printf("Determine status of service.\n")
		}


		err = nil
	}

	if err != nil {
		//eblog.Debug(me.EntityId, "Error: %v", err)
	}

	return err
}
