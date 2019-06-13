package eventbroker

import (
	"errors"
	"fmt"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/network"
	"gearbox/heartbeat/eventbroker/only"
	"gearbox/heartbeat/eventbroker/states"
	"net/url"
	"reflect"
	"sort"
)


func (me *EventBroker) EnsureNotNil() error {

	var err error

	if me == nil {
		err = errors.New("unexpected software error")
	}
	return err
}


func (me *EventBroker) RegisterService(topic string, args ...string) {
	fmt.Printf("RegisterService\n")

	// .

	return
}


var ServiceMqtt = network.ServiceConfig{
	Name:   "_gearbox-mqtt",
	Type:   "_mqtt._tcp",
	Domain: "local",
}


func InterfaceToTypeEventBroker(i interface{}) (*EventBroker, error) {

	var err error
	var me *EventBroker

	for range only.Once {
		if i == nil {
			err = errors.New("interface is nil, should be" + InterfaceTypeEventBroker)
			break
		}

		checkType := reflect.ValueOf(i)
		if checkType.Type().String() != InterfaceTypeEventBroker {
			err = errors.New("interface type not " + InterfaceTypeEventBroker)
			break
		}

		me = i.(*EventBroker)

		err = me.EnsureNotNil()
		if err != nil {
			break
		}
	}

	return me, err
}


func (me *EventBroker) FindMqttBroker() (*url.URL, error) {

	var err error
	var mqttService *network.Service
	var u *url.URL

	for range only.Once {

		//fmt.Printf("\n\n################################################################################\n")
		//err = me.ZeroConf.PrintServices()

		mqttService, err = me.ZeroConf.FindService(ServiceMqtt)
		if err != nil {
			fmt.Printf("Error(me.ZeroConf.FindService): %v\n", err)
			break
		}

		u, err = url.Parse(fmt.Sprintf("tcp://%s:%d", mqttService.Entry.HostName, mqttService.Entry.Port))
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "Found MQTT broker %s", u.String())
	}

	return u, err
}


func (me *Services) AddState(state states.Status) error {

	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("status.Status is nil")
			break
		}

		if state.EntityId == nil {
			err = errors.New("status.EntityId is nil")
			break
		}

		me.mutex.Lock()
		defer me.mutex.Unlock()

		client := state.EntityId

		if _, ok := me.States[*client]; !ok {
			me.States[*client] = &states.Status{}		// states.New(client, client, entity.BroadcastEntityName)
		}

		me.States[*client] = &state
		// mutex:      &sync.RWMutex{}

		//me.Logs = append(me.Logs, Log{
		//	State: state,
		//	When: time.Now(),
		//})
	}

	return err
}


func (me *Services) PrintStates() error {

	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("status.Status is nil")
			break
		}

		var keys []string
		for k := range me.States {
			keys = append(keys, k.String())
		}
		sort.Strings(keys)

		for i, e := range keys {
			s := me.States[messages.MessageAddress(e)]
			fmt.Printf("%d\tCurrent:%s\tName:%s\tEntityId:%s\tLast:%s\tAction:%s\tParent:%s\n",
				i,
				s.Current.String(),
				s.EntityName.String(),
				s.EntityId.String(),
				s.Last.String(),
				s.Action.String(),
				s.ParentId.String(),
			)
			//fmt.Printf("%d  EntityId:%s  Name:%s  Parent:%s  Action:%s  Want:%s  Current:%s  Last:%s  LastWhen:%v  Error:%v\n",
			//	i,
			//	s.EntityId.String(),
			//	s.EntityName.String(),
			//	s.ParentId.String(),
			//	s.Action.String(),
			//	s.Want.String(),
			//	s.Current.String(),
			//	s.Last.String(),
			//	s.LastWhen.Unix(),
			//	s.Error,
			//)
		}
	}

	return err
}


//func (me *Services) DeleteState(client messages.SubTopic) error {
//
//	var err error
//
//	me.mutex.Lock()
//	defer me.mutex.Unlock()
//
//	for range only.Once {
//		_, ok := me.topics[client] // Managed by Mutex
//		if !ok {
//			err = me.EntityId.ProduceError("service doesn't exist")
//			break
//		}
//
//		delete(me.topics, client) // Managed by Mutex
//	}
//
//	return err
//}
//
//
//func (me *EventBroker) zcByChannel(s network.ServiceConfig) (*network.Service, error) {
//
//	var err error
//	var sc *network.Service
//
//	fmt.Printf("Register service by channel...\n")
//	sc, err = me.ZeroConf.RegisterByChannel(messages.MessageAddress(me.EntityId), s)
//
//	return sc, err
//}
//
//
//func (me *EventBroker) zcByMethod(s network.ServiceConfig) (*network.Service, error) {
//
//	var err error
//	var sc *network.Service
//
//	fmt.Printf("Register service by method...\n")
//	sc, err = me.ZeroConf.Register(s)
//
//	return sc, err
//}
//
//
//func (me *EventBroker) CreateEntity(serviceName string) {
//
//	var err error
//
//	fmt.Printf("\n\n################################################################################\n")
//	_ = me.ZeroConf.PrintServices()
//
//	s1 := network.ServiceConfig{
//		Name: network.Name(serviceName + "1"),
//		Type: "_gearbox._tcp",
//		Domain: "local",
//		Port: network.SelectRandomPort,
//	}
//
//	s2 := network.ServiceConfig{
//		Name: network.Name(serviceName + "2"),
//		Type: "_gearbox._tcp",
//		Domain: "local",
//		Port: network.SelectRandomPort,
//	}
//
//	s3 := network.ServiceConfig{
//		Name: network.Name(serviceName + "3"),
//		Type: "_gearbox._tcp",
//		Domain: "local",
//		Port: network.SelectRandomPort,
//	}
//
//	var s1ref *network.Service
//	s1ref, err = me.zcByChannel(s1)
//	fmt.Printf("Response(me.zcByChannel): %v\n", err)
//	// s1ref, err = me.ZeroConf.GetReference(s1)
//	fmt.Printf("Response(me.ZeroConf.GetReference): %v\n%v\n", err, s1ref)
//	_ = s1ref.Print()
//
//	var s2ref *network.Service
//	s2ref, err = me.zcByMethod(s2)
//	fmt.Printf("Response(me.zcByMethod): %v\n", err)
//	_ = s2ref.Print()
//
//	var s3ref *network.Service
//	s3ref, err = me.zcByMethod(s3)
//	fmt.Printf("Response(me.zcByMethod): %v\n", err)
//	_ = s3ref.Print()
//
//
//	time.Sleep(time.Second * 7)
//
//
//	fmt.Printf("Listeners...\n")
//	_ = me.ZeroConf.PrintServices()
//
//	//time.Sleep(time.Minute * 600)
//	//me.FindMqtt()
//
//	time.Sleep(time.Second * 700)
//
//	err = me.ZeroConf.UnregisterByChannel(me.EntityId, s1ref.EntityId)
//	fmt.Printf("Response(me.ZeroConf.UnregisterByChannel): %v\n", err)
//
//	err = me.ZeroConf.UnregisterByEntityId(s1ref.EntityId)
//	fmt.Printf("Response(me.ZeroConf.UnregisterByEntityId): %v\n", err)
//
//	//err = s3ref.Unregister()
//	//fmt.Printf("Response(s3ref.Unregister): %v\n", err)
//
//
//	//fmt.Printf("Start channel...\n")
//	//channelService, _ := me.Channels.StartHandler(messages.MessageAddress(serviceName))
//	//err = channelService.Subscribe(messages.SubTopic("start"), manageService, s1ref)
//	//if err != nil {
//	//	return
//	//}
//	//err = channelService.Subscribe(messages.SubTopic("stop"), manageService, s1ref)
//	//if err != nil {
//	//	return
//	//}
//	//err = channelService.Subscribe(messages.SubTopic("status"), manageService, s1ref)
//	//if err != nil {
//	//	return
//	//}
//	//fmt.Printf("List channel...\n")
//	//channelService.List()
//
//	time.Sleep(time.Second * 1)
//
//	//fmt.Printf("Stopping channel...\n")
//	//_ = channelService.StopHandler()
//}

