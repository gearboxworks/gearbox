package gbevents

import (
	"fmt"
	"gearbox/app/logger"
	"gearbox/box"
	"gearbox/global"
	"gearbox/heartbeat/gbevents/channels"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/heartbeat/gbevents/network"
	"gearbox/help"
	"gearbox/only"
	oss "gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/jinzhu/copier"
	"os"
	"path/filepath"
	"time"
)


func New(OsSupport oss.OsSupporter, args ...Args) (*EventBroker, status.Status) {

	var _args Args
	var sts status.Status
	se := &EventBroker{}

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		if _args.Boxname == "" {
			_args.Boxname = global.Brandname
		}

		if _args.EntityId == "" {
			_args.EntityId = "gbevents" // uuid.New()
		}

		_args.OsSupport = OsSupport
		foo := box.Args{}
		err := copier.Copy(&foo, &_args)
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("unable to configure event handler").
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		// Channels allow inter-thread communications.
		err = _args.Channels.New(OsSupport, channels.Args{})
		if err != nil {
			sts = status.Fail().
				SetMessage("failed to init channels handler").
				SetAdditional("", ).
				SetData(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		// 2. ZeroConf
		err = _args.ZeroConf.New(OsSupport, network.Args{Channels: &_args.Channels})
		if err != nil {
			sts = status.Fail().
				SetMessage("failed to init zeroconf handler").
				SetAdditional("", ).
				SetData(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		//_, _, _ = _args.FindMqttBroker()
		//
		//// 3. MQTT allows inter-process communications.
		//err = _args.MqttClient.New(OsSupport, mqttClient.Args{})
		//if err != nil {
		//	sts = status.Fail().
		//		SetMessage("failed to init MQTT client handler").
		//		SetAdditional("", ).
		//		SetData(err).
		//		SetHelp(status.AllHelp, help.ContactSupportHelp())
		//	break
		//}

		// _args.ZeroConf.Browse("_workstation._tcp")
		// daemon.SimpleWaitLoop("ZeroConf", 2000, time.Second * 5)

		//sts = _args.MqttBroker.New(OsSupport, )
		//if is.Error(sts) {
		//	break
		//}

		_args.PidFile = filepath.FromSlash(fmt.Sprintf("%s/%s", _args.OsSupport.GetAdminRootDir(), defaultPidFile))

		*se = EventBroker(_args)

		sts = status.Success("created new event broker")
	}
	// status.Log(sts)

	return se, sts
}


func (me *EventBroker) RegisterService(topic string, args ...ServiceData) {
	fmt.Printf("RegisterService\n")

	// .

	return
}


func (me *EventBroker) Create() status.Status {
	fmt.Printf("(me *EventBroker) CreateService() status.Status\n")

	return nil
}


func (me *EventBroker) Start() status.Status {

	var err error
	var sts status.Status

	for range only.Once {
		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

/*
		me.Identifier = "thisisme"
		// Start the inter-thread channels service.
		fmt.Printf("# DEBUG1a\n")
		me1, _ := me.Channels.StartHandler(messages.MessageAddress(me.Identifier))
		//me.Channels.Subscribe(messages.StringsToTopic(me.Identifier, "testme"), testme)
		//me.Channels.Subscribe(messages.StringsToTopic(me.Identifier, "exit"), checkExit)
		me1.Subscribe(messages.SubTopic("testme"), testme)
		me1.Subscribe(messages.SubTopic("exit"), checkExit)
		me1.Subscribe(messages.SubTopic("harry"), testme)

		Identifier2 := "hello"
		// Start the inter-thread channels service.
		fmt.Printf("# DEBUG1b\n")
		me2, _ := me.Channels.StartHandler(messages.MessageAddress(Identifier2))
		// me.Channels.Subscribe(messages.StringsToTopic(Identifier2, "testme"), testme)
		me2.Subscribe(messages.SubTopic("testme"), testme)


		fmt.Printf("# DEBUG2\n")
		time.Sleep(time.Second * 2)
		me.Channels.Publish(messages.Message{Topic: messages.StringsToTopic(Identifier2, "question"), Text: "What about now?"}).Log()
		me.Channels.Publish(messages.Message{Topic: messages.StringsToTopic(me.Identifier, "statement"), Text: "not now"}).Log()

		time.Sleep(time.Second * 1)
		me.Channels.Publish(messages.Message{Topic: messages.StringsToTopic(Identifier2, "statement"), Text: "Come on!"}).Log()
		me.Channels.Publish(messages.Message{Topic: messages.StringsToTopic(me.Identifier, "statement"), Text: "wait for it"}).Log()

		me.Channels.Publish(messages.Message{Topic: messages.StringsToTopic(Identifier2, "question"), Text: "Really?"}).Log()
		me.Channels.Publish(messages.Message{Topic: messages.StringsToTopic(me.Identifier, "statement"), Text: "not yet"}).Log()

		me.Channels.Publish(messages.Message{Topic: messages.StringsToTopic(Identifier2, "statement"), Text: "About time."}).Log()
		me.Channels.Publish(messages.Message{Topic: messages.StringsToTopic(me.Identifier, "statement"), Text: "almost there"}).Log()

		time.Sleep(time.Second * 1)

		fmt.Printf("# DEBUG3\n")
		//me.Channels.Publish(messages.Message{Topic: messages.StringsToTopic(me.Identifier, "exit"), Text: "now"}).Log()
		// me.Channels.StopHandler(messages.MessageAddress(me.Identifier))
		me1.StopHandler()
		time.Sleep(time.Second * 4)
		// me.Channels.StopHandler(messages.MessageAddress(me.Identifier))
		me2.StopHandler()

		os.Exit(0)
*/


		// Start the inter-thread service.
		// Note: These will be started dynamically as clients are registered.

		// Start the inter-process service.
		//go func() {
		//	err = me.MqttBroker.StartHandler()
		//	if err != nil {
		//		sts = status.Fail().
		//			SetMessage("failed to start MQTT handler").
		//			SetAdditional("", ).
		//			SetData(err).
		//			SetHelp(status.AllHelp, help.ContactSupportHelp())
		//		break
		//	}
		//}()

		// Start the zeroconf service client.
		//go func() {
		err = me.ZeroConf.StartHandler()
		if err != nil {
			sts = status.Fail().
				SetMessage("failed to start zeroconf handler").
				SetAdditional("", ).
				SetData(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		// Start the inter-process service.
		err = me.MqttClient.StartHandler()
		if err != nil {
			sts = status.Fail().
				SetMessage("failed to start MQTT handler").
				SetAdditional("", ).
				SetData(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}
		fmt.Printf("Waiting...\n")

		time.Sleep(time.Second * 6)

		//me.CreateEntity("HELO")
		_, _, _ = me.FindMqttBroker()

		time.Sleep(time.Second * 20)

		err = me.ZeroConf.StopHandler()
		time.Sleep(time.Second * 3)
		fmt.Printf("Sleeping...\n")
		os.Exit(0)

		//}()

		sts = status.Success("started event broker")
	}

	return sts
}


func (me *EventBroker) FindMqttBroker() (string, int, error) {

	var err error
	var mqttService *network.Service

	fmt.Printf("\n\n################################################################################\n")
	_ = me.ZeroConf.PrintServices()

	mqtt := network.CreateEntry{
		Name:   network.Name("_gearbox-mqtt"),
		Type:   "_mqtt._tcp",
		Domain: "local",
	}

	mqttService, err = me.ZeroConf.FindService(mqtt)
	if err != nil {
		fmt.Printf("Error(me.ZeroConf.FindService): %v\n", err)
	} else {
		fmt.Printf("Found: %v\n", mqttService)
	}

	return mqttService.Entry.HostName, mqttService.Entry.Port, err
}


func (me *EventBroker) zcByChannel(s network.CreateEntry) (*network.Service, error) {

	var err error
	var sc *network.Service

	fmt.Printf("Register service by channel...\n")
	sc, err = me.ZeroConf.RegisterByChannel(messages.MessageAddress(me.EntityId), s)

	return sc, err
}


func (me *EventBroker) zcByMethod(s network.CreateEntry) (*network.Service, error) {

	var err error
	var sc *network.Service

	fmt.Printf("Register service by method...\n")
	sc, err = me.ZeroConf.Register(s)

	return sc, err
}


func (me *EventBroker) CreateEntity(serviceName string) {

	var err error

	fmt.Printf("\n\n################################################################################\n")
	_ = me.ZeroConf.PrintServices()

	s1 := network.CreateEntry{
		Name: network.Name(serviceName + "1"),
		Type: "_gearbox._tcp",
		Domain: "local",
		Port: 0,
	}

	s2 := network.CreateEntry{
		Name: network.Name(serviceName + "2"),
		Type: "_gearbox._tcp",
		Domain: "local",
		Port: 0,
	}

	s3 := network.CreateEntry{
		Name: network.Name(serviceName + "3"),
		Type: "_gearbox._tcp",
		Domain: "local",
		Port: 0,
	}

	var s1ref *network.Service
	s1ref, err = me.zcByChannel(s1)
	fmt.Printf("Response(me.zcByChannel): %v\n", err)
	// s1ref, err = me.ZeroConf.GetReference(s1)
	fmt.Printf("Response(me.ZeroConf.GetReference): %v\n%v\n", err, s1ref)
	_ = s1ref.Print()

	var s2ref *network.Service
	s2ref, err = me.zcByMethod(s2)
	fmt.Printf("Response(me.zcByMethod): %v\n", err)
	_ = s2ref.Print()

	var s3ref *network.Service
	s3ref, err = me.zcByMethod(s3)
	fmt.Printf("Response(me.zcByMethod): %v\n", err)
	_ = s3ref.Print()


	time.Sleep(time.Second * 7)


	fmt.Printf("Listeners...\n")
	_ = me.ZeroConf.PrintServices()

	//time.Sleep(time.Minute * 600)
	//me.FindMqtt()

	time.Sleep(time.Second * 700)

	err = me.ZeroConf.UnregisterByChannel(s1ref.EntityId)
	fmt.Printf("Response(me.ZeroConf.UnregisterByChannel): %v\n", err)

	err = me.ZeroConf.UnregisterByUuid(s1ref.EntityId)
	fmt.Printf("Response(me.ZeroConf.UnregisterByUuid): %v\n", err)

	err = s3ref.Unregister()
	fmt.Printf("Response(s3ref.Unregister): %v\n", err)


	//fmt.Printf("Start channel...\n")
	//channelService, _ := me.Channels.StartHandler(messages.MessageAddress(serviceName))
	//err = channelService.Subscribe(messages.SubTopic("start"), manageService, s1ref)
	//if err != nil {
	//	return
	//}
	//err = channelService.Subscribe(messages.SubTopic("stop"), manageService, s1ref)
	//if err != nil {
	//	return
	//}
	//err = channelService.Subscribe(messages.SubTopic("status"), manageService, s1ref)
	//if err != nil {
	//	return
	//}
	//fmt.Printf("List channel...\n")
	//channelService.List()

	time.Sleep(time.Second * 1)

	//fmt.Printf("Stopping channel...\n")
	//_ = channelService.StopHandler()
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
		logger.Debug("Error: %v", err)
	}

	return err
}


func (me *EventBroker) Stop() status.Status {
	fmt.Printf("(me *EventBroker) StopService() status.Status\n")

	return nil
}


func (me *EventBroker) Restart() status.Status {
	fmt.Printf("(me *EventBroker) RestartService() status.Status\n")

	return nil
}


func (me *EventBroker) Status() status.Status {
	fmt.Printf("(me *EventBroker) ServiceStatus() status.Status\n")

	return nil
}
