package eventbroker

import (
	"fmt"
	"gearbox/box"
	"gearbox/global"
	"gearbox/eventbroker/channels"
	"gearbox/eventbroker/daemon"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/mqttClient"
	"gearbox/eventbroker/network"
	"gearbox/eventbroker/only"
	"gearbox/eventbroker/ospaths"
	"gearbox/eventbroker/states"
	"github.com/jinzhu/copier"
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

		foo := box.Args{}
		err = copier.Copy(&foo, &_args)
		if err != nil {
			err = me.EntityId.ProduceError("unable to copy config args")
			break
		}

		if _args.Boxname == "" {
			_args.Boxname = global.Brandname
		}

		if _args.EntityId == "" {
			_args.EntityId = DefaultEntityName
		}
		_args.State = states.New(&_args.EntityId, &_args.EntityId, entity.SelfEntityName)

		if _args.SubBaseDir == "" {
			_args.SubBaseDir = ospaths.DefaultBaseDir
		}
		_args.OsPaths = ospaths.New(_args.SubBaseDir)
		err = _args.OsPaths.CreateIfNotExists()
		if err != nil {
			break
		}

		//_args.Services.States = make(States)
		//_args.Services.Callbacks = make(Callbacks)
		//_args.Services.CallbackLocks = make(CallbackLocks)
		//_args.Services.Logs = make(Logs, 0)
		_args.Services = make(Services)


		*me = EventBroker(_args)


		// 1. Channel - provides inter-thread communications.
		err = me.Channels.New(channels.Args{Boxname: me.Boxname, OsPaths: me.OsPaths})
		if err != nil {
			break
		}
		err = me.StartChannelHandler()
		if err != nil {
			break
		}


		// 2. ZeroConf - provides discovery and management of network services.
		err = me.ZeroConf.New(network.Args{Boxname: me.Boxname, Channels: &me.Channels, OsPaths: me.OsPaths})
		if err != nil {
			break
		}


		// 3. Daemon - provides control over arbitrary services.
		err = me.Daemon.New(daemon.Args{Boxname: me.Boxname, Channels: &me.Channels, OsPaths: me.OsPaths})
		if err != nil {
			break
		}


		// 4. MQTT broker - provides inter-process communications.


		// 5. MQTT client - provides inter-process communications.
		err = me.MqttClient.New(mqttClient.Args{Boxname: me.Boxname, Channels: &me.Channels, OsPaths: me.OsPaths})
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


func (me *EventBroker) SimpleLoop() {

	//var state states.Status
	var err error

	fmt.Printf("SimpleLoop()\n")
	//services := ServiceData{
	//	Now: make(Entities),
	//	Logs: []Log{},
	//}

	var Ret1 string

	var sc1 *Service
	sc1, _, err = me.AttachCallback(messages.MessageAddress(entity.UnfsdEntityName), myCallback, &Ret1)
	if err != nil {
		fmt.Printf("Ooops\n")
	}
	_ = sc1.PrintState()

	var sc2 *Service
	sc2, _, err = me.AttachCallback(messages.MessageAddress(entity.MqttBrokerEntityName), myCallback, &Ret1)
	if err != nil {
		fmt.Printf("Ooops\n")
	}
	_ = sc2.PrintState()

	for i := 0; i < 1000; i++ {
		fmt.Printf("\n############\n")
		//for _, e := range entity.PartialEntities {
		//	services.Now[e].State, err = me.StatusOf(e)
		//	if err != nil {
		//		fmt.Printf("%s is at error %v\n", e.String(), err)
		//	} else {
		//		services.Logs = append(services.Logs, Log{
		//			When: time.Now(),
		//			State: services.Now[e].State,
		//		})
		//		//fmt.Printf("%s is at state '%s'\n", e.String(), state.Current.String())
		//	}
		//}

		me.Services.PrintStates()
		//fmt.Printf("States: N:'%s'\tD:'%s'\tM:'%s'\n",
		//	me.Services.States[entity.NetworkEntityName].Current.String(),
		//	me.Services.States[entity.DaemonEntityName].Current.String(),
		//	me.Services.States[entity.MqttClientEntityName].Current.String())

		time.Sleep(time.Second * 20)
	}

	//for i := 0; i < 1000; i++ {
	//	fmt.Printf("\n############\n")
	//	for _, e := range entity.PartialEntities {
	//		services.Now[e].State, err = me.StatusOf(e)
	//		if err != nil {
	//			fmt.Printf("%s is at error %v\n", e.String(), err)
	//		} else {
	//			services.Logs = append(services.Logs, Log{
	//				When: time.Now(),
	//				Status: services.Now[e].State,
	//			})
	//			//fmt.Printf("%s is at state '%s'\n", e.String(), state.Current.String())
	//		}
	//	}
	//
	//	fmt.Printf("States: N:'%s'\tD:'%s'\tM:'%s'\n",
	//		services.Now[entity.NetworkEntityName].State.Current.String(),
	//		services.Now[entity.DaemonEntityName].State.Current.String(),
	//		services.Now[entity.MqttClientEntityName].State.Current.String())
	//
	//	time.Sleep(time.Second * 3)
	//}

	time.Sleep(time.Hour * 60)
}

var RetFunc string

func myCallback(args interface{}, state states.Status) error {

	var err error

	ret := args.(string)

	fmt.Printf("CB state: %s (%v)\n", state.ShortString(), ret)
	//fmt.Printf("HELLO state: %s\n", PrintServiceState(&state))
	//fmt.Printf("EntityId:%s  Name:%s  Parent:%s  Action:%s  Want:%s  Current:%s  Last:%s  LastWhen:%v  Error:%v\n",
	//	state.EntityId.String(),
	//	state.EntityName.String(),
	//	state.ParentId.String(),
	//	state.Action.String(),
	//	state.Want.String(),
	//	state.Current.String(),
	//	state.Last.String(),
	//	state.LastWhen.Unix(),
	//	state.Error,
	//)


	return err
}


func (me *EventBroker) StatusOf(client messages.MessageAddress) (states.Status, error) {

	var ret states.Status
	var err error

	msg := messages.Message{
		Source: me.EntityId,
		Topic: messages.MessageTopic{
			Address: client,
			SubTopic: "status",
		},
		Text: "",
	}

	wrapper := messages.Message{
		Source: me.EntityId,
		Topic: messages.MessageTopic{
			Address: entity.BroadcastEntityName,
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


type SimpleState map[messages.MessageAddress]states.State
func (me *EventBroker) GetSimpleStatus() (SimpleState, error) {

	ret := make(SimpleState)
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


func (me SimpleState) String() string {

	var ret string

	for range only.Once {
		for k, v := range me {
			ret += fmt.Sprintf("%s %s\n", k, v)
		}
	}

	return ret
}


//func (me *EventBroker) TempLoop() error {
//
//	var err error
//
//	fmt.Printf("(me *EventBroker) TempLoop()\n")
//
//	//msg := messages.Message{
//	//	Source: me.EntityId,
//	//	Topic: messages.MessageTopic{
//	//		Address: "eventbrokerdaemon",
//	//		SubTopic: "get",
//	//	},
//	//	Text: "status",
//	//}
//	msg := messages.Message{
//		Source: me.EntityId,
//		Topic: messages.MessageTopic{
//			Address: "eventbroker-daemon",
//			SubTopic: "status",
//		},
//		Text: "now",
//	}
//	//time.Sleep(time.Second * 8)
//	//err = me.Daemon.LoadFiles()
//	//fmt.Printf("me.Daemon.LoadFiles(): %v\n", err)
//	//me.CreateEntity("BEEP")
//	fmt.Printf("####################################################################\nPING\n")
//	//err = me.Daemon.LoadFiles()
//	fmt.Printf("me.Daemon.LoadFiles(): %v\n", err)
//	fmt.Printf("####################################################################\nPING\n")
//
//
//	//me.Foo()
//	me.SimpleLoop()
//
//	time.Sleep(time.Hour * 8000)
//
//	go func() {
//		index := 0
//		for {
//			fmt.Printf("####################################################################\nPING\n")
//			//me.Daemon.ChannelHandler.List()
//			//me.Daemon.Channels.ListSubscribers()
//
//			//fmt.Printf("Error1: %v\n", me.Daemon.State.Error)
//			//me.Daemon.State.SetError(errors.New(fmt.Sprintf("Loop #%d (%s)", index, me.Daemon.Fluff)))
//			//fmt.Printf("Error2: %v\n", me.Daemon.State.Error)
//
//			fmt.Printf("\n\n%d gbevents before: %v\n", time.Now().Unix(), me.Daemon.State.GetError())
//			i, _ := me.Channels.PublishAndWaitForReturn(msg, 400)
//			//foo := reflect.ValueOf(i)
//			//fmt.Printf("ERROR: %v\t\tRESPONSE: %v\n", err, i)
//			//fmt.Printf("Reflect %s, %s\n", foo.Type(), foo.String())
//			f, err := states.InterfaceToTypeStatus(i)
//			if err == nil {
//				fmt.Printf("%d gbevents after: %v\n", time.Now().Unix(), f.GetError())
//				// fmt.Printf("%d gbevents after: %v (%v)\n", time.Now().Unix(), f.GetError(), me.Daemon.Fluff)
//			} else {
//				fmt.Printf("%d gbevents after: is nil!\n", time.Now().Unix())
//			}
//
//			index++
//			//me.Daemon.State.Error = nil
//			time.Sleep(time.Second * 20)
//		}
//	}()
//
//	//time.Sleep(time.Second * 6)
//	//
//	////me.CreateEntity("HELO")
//	//var u *url.URL
//	//u, err = me.FindMqttBroker()
//	//if err == nil {
//	//	err = me.MqttClient.ConnectToServer(u.String())
//	//}
//	//if err != nil {
//	//	eblog.Debug(me.EntityId, "Aaaaargh! => %v", err)
//	//}
//	//
//	//_ = me.MqttClient.GlobSubscribe(me.MqttClient.EntityId)
//	//
//	//time.Sleep(time.Second * 200)
//	//
//	//err = me.Stop()
//	//
//	//fmt.Printf("Exiting...\n")
//
//	return err
//}
//
//
//func (me *EventBroker) Foo() {
//
//	var st messages.SubTopics
//	var t messages.Topics
//	var ma messages.MessageAddresses
//
//	msg := messages.Message{
//		Source: me.EntityId,
//		Topic: messages.MessageTopic{
//			Address: "eventbroker-daemon",
//			SubTopic: "get",
//		},
//		Text: "topics",
//	}
//	fmt.Printf("\n\n%d gbevents before: %v\n", time.Now().Unix(), me.Daemon.State.GetError())
//	i, _ := me.Channels.PublishAndWaitForReturn(msg, 400)
//	f, err := messages.InterfaceToTypeSubTopics(i)
//	if err == nil {
//		fmt.Printf("%d gbevents after: %v\n", time.Now().Unix(), f)
//	} else {
//		fmt.Printf("%d gbevents after: is nil!\n", time.Now().Unix())
//	}
//
//	//time.Sleep(time.Hour * 8000)
//
//	fmt.Printf("\n### Channels\n")
//	fmt.Printf("me.Channels.GetManagedEntities\n")
//	ma = me.Channels.GetManagedEntities()
//	for _, f := range ma {
//		fmt.Printf("me.Channels.GetManagedEntities\t=> %s\n", f.String())
//	}
//	fmt.Printf("me.Channels.GetEntities\n")
//	ma = me.Channels.GetEntities()
//	for _, f := range ma {
//		fmt.Printf("me.Channels.GetEntities\t=> %s\n", f.String())
//	}
//	fmt.Printf("me.Channels.GetListenerTopics\n")
//	t, err = me.Channels.GetListenerTopics()
//	if err != nil {
//		fmt.Printf("me.Channels.GetListenerTopics\tERR:%v\n", err)
//	}
//	for _, f := range t {
//		fmt.Printf("me.Channels.GetListenerTopics\t=> %s\n", f.String())
//	}
//
//	//fmt.Printf("\n### eventbroker\n")
//	//fmt.Printf("me.GetEntities\n")
//	//ma = me.GetEntities()
//	//for _, f := range ma {
//	//	fmt.Printf("me.GetEntities => %s\n", f.String())
//	//}
//	//
//	//fmt.Printf("me.GetManagedEntities\n")
//	//ma = me.Daemon.GetManagedEntities()
//	//for _, f := range ma {
//	//	fmt.Printf("me.GetManagedEntities => %s\n", f.String())
//	//}
//	//
//	//fmt.Printf("me.GetTopics\n")
//	//st = me.Daemon.GetTopics()
//	//for _, f := range st {
//	//	fmt.Printf("me.GetTopics => %s\n", f.String())
//	//}
//
//	fmt.Printf("\n### ZeroConf\n")
//	fmt.Printf("me.ZeroConf.GetEntities\n")
//	ma = me.ZeroConf.GetEntities()
//	for _, f := range ma {
//		fmt.Printf("me.ZeroConf.GetEntities\t=> %s\n", f.String())
//	}
//
//	fmt.Printf("me.ZeroConf.GetManagedEntities\n")
//	ma = me.ZeroConf.GetManagedEntities()
//	for _, f := range ma {
//		fmt.Printf("me.ZeroConf.GetManagedEntities\t=> %s\n", f.String())
//	}
//
//	fmt.Printf("me.ZeroConf.GetTopics\n")
//	st = me.ZeroConf.GetTopics()
//	for _, f := range st {
//		fmt.Printf("me.ZeroConf.GetTopics\t=> %s\n", f.String())
//	}
//
//
//	fmt.Printf("\n### Daemon\n")
//	fmt.Printf("me.Daemon.GetEntities\n")
//	ma = me.Daemon.GetEntities()
//	for _, f := range ma {
//		fmt.Printf("me.Daemon.GetEntities\t=> %s\n", f.String())
//	}
//
//	fmt.Printf("me.Daemon.GetManagedEntities\n")
//	ma = me.Daemon.GetManagedEntities()
//	for _, f := range ma {
//		fmt.Printf("me.Daemon.GetManagedEntities\t=> %s\n", f.String())
//	}
//
//	fmt.Printf("me.Daemon.GetTopics\n")
//	st = me.Daemon.GetTopics()
//	for _, f := range st {
//		fmt.Printf("me.Daemon.GetTopics\t=> %s\n", f.String())
//	}
//
//
//	fmt.Printf("\n### MqttClient\n")
//	fmt.Printf("me.MqttClient.GetEntities\n")
//	ma = me.MqttClient.GetEntities()
//	for _, f := range ma {
//		fmt.Printf("me.MqttClient.GetEntities\t=> %s\n", f.String())
//	}
//
//	fmt.Printf("me.MqttClient.GetManagedEntities\n")
//	ma = me.MqttClient.GetManagedEntities()
//	for _, f := range ma {
//		fmt.Printf("me.MqttClient.GetManagedEntities\t=> %s\n", f.String())
//	}
//
//	fmt.Printf("me.MqttClient.GetTopics\n")
//	st = me.MqttClient.GetTopics()
//	for _, f := range st {
//		fmt.Printf("me.MqttClient.GetTopics\t=> %s\n", f.String())
//	}
//}


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

