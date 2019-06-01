package gbevents

import (
	"fmt"
	"gearbox/box"
	"gearbox/global"
	"gearbox/heartbeat/gbevents/channels"
	"gearbox/heartbeat/gbevents/network"
	"gearbox/help"
	"gearbox/only"
	oss "gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/jinzhu/copier"
	"path/filepath"
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

		// ZeroConf
		err = _args.ZeroConf.New(OsSupport, network.Args{Channels: &_args.Channels})
		if err != nil {
			sts = status.Fail().
				SetMessage("failed to init zeroconf handler").
				SetAdditional("", ).
				SetData(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

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

		//}()

		sts = status.Success("started event broker")
	}

	return sts
}


//func checkExit(msg *messages.Message) status.Status {
//
//	var sts status.Status
//	messages.Debug("MSG:%s", msg.Topic)
//
//	if (msg.Topic.SubTopic == "exit") && (msg.Text == "now") {
//		messages.Debug("Hey! It works! Awesome.")
//	}
//
//	return sts
//}
//
//
//func testme(msg *messages.Message) status.Status {
//
//	var sts status.Status
//
//	messages.Debug("testme() '%s' == '%s'", msg.Topic, msg.Text)
//	// messages.Debug(">>>>>> testme(%s)	Time:%v	Src:%s	Text:%s\n", msg.Topic, msg.Time.Convert().Unix(), msg.Src, msg.Text)
//
//	return sts
//}


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
