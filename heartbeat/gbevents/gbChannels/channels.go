package gbChannels

import (
	"fmt"
	"gearbox/box"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/help"
	"gearbox/only"
	oss "gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/jinzhu/copier"
	"github.com/olebedev/emitter"
)


func (me *Channels) New(OsSupport oss.OsSupporter, args ...Args) status.Status {

	var _args Args
	var sts status.Status

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		_args.OsSupport = OsSupport
		foo := box.Args{}
		err := copier.Copy(&foo, &_args)
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("unable to copy Go channels config").
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		if _args.EntityId == "" {
			_args.EntityId = DefaultEntityId
		}

		_args.instance.emitter = emitter.Emitter{}
		_args.Subscribers = make(Subscribers)

		*me = Channels(_args)

		messages.Debug("GBevents - channel (%s).", me.EntityId.String())
		sts = status.Success("MQTT started OK on ", me.EntityId.String())
	}

	if !is.Success(sts) {
		sts.Log()
	}

	// Save last state.
	me.Sts = sts

	return sts
}


func (me *Channels) Subscribe(topic messages.Topic, fn Callback) (SubTopics, status.Status) {

	var sts status.Status

	for range only.Once {
		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		sts = topic.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		if fn == nil {
			sts = status.Warn("").
				SetMessage("callback function is empty").
				SetAdditional("", ).
				SetData("").
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		if me.Subscribers == nil {
			me.Subscribers = make(Subscribers)
		}

		if _, ok := me.Subscribers[topic.Address]; !ok {
			me.Subscribers[topic.Address] = SubTopics{
				Address: topic.Address,
				Callbacks: make(Callbacks),
				instance: &me.instance,
			}
		}

		me.Subscribers[topic.Address].Callbacks[topic.SubTopic] = fn
		foo := me.Subscribers[topic.Address]
		foo.List()
		me.Subscribers.List()

		sts = status.Success("New subscriber: %s\n", messages.SprintfTopic(topic.Address, topic.SubTopic))
	}

	if !is.Success(sts) {
		sts.Log()
	}

	// Save last state.
	me.Sts = sts

	if _, ok := me.Subscribers[topic.Address]; ok {
		return me.Subscribers[topic.Address], sts
	} else {
		return SubTopics{}, sts
	}
}


func (me *Channels) UnSubscribe(topic messages.Topic) status.Status {

	var sts status.Status

	for range only.Once {
		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		sts = topic.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		delete(me.Subscribers[topic.Address].Callbacks, topic.SubTopic)
		sts = status.Success("").
			SetMessage("Unsubscribed: %s\n", messages.SprintfTopic(topic.Address, topic.SubTopic)).
			SetAdditional("emitter:%v", me.instance.emitter).
			SetData(topic).
			SetHelp(status.AllHelp, help.ContactSupportHelp())
	}

	if !is.Success(sts) {
		sts.Log()
	}

	// Save last state.
	me.Sts = sts

	return sts
}


func (me *SubTopics) Subscribe(topic messages.SubTopic, fn Callback) status.Status {

	var sts status.Status

	for range only.Once {
		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		sts = topic.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		if fn == nil {
			sts = status.Warn("").
				SetMessage("callback function is empty").
				SetAdditional("", ).
				SetData("").
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		me.Callbacks[topic] = fn
		me.List()

		sts = status.Success("").
			SetMessage("New subscriber: %s\n", messages.SprintfTopic(me.Address, topic)).
			SetAdditional("emitter:%v", me.instance.emitter).
			SetData(topic).
			SetHelp(status.AllHelp, help.ContactSupportHelp())
	}

	if !is.Success(sts) {
		sts.Log()
	}

	return sts
}


func (me *SubTopics) UnSubscribe(topic messages.SubTopic) status.Status {

	var sts status.Status

	for range only.Once {
		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		sts = topic.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		delete(me.Callbacks, topic)

		sts = status.Success("").
			SetMessage("Unsubscribed: %s\n", messages.SprintfTopic(me.Address, topic)).
			SetAdditional("emitter:%v", me.instance.emitter).
			SetData(topic).
			SetHelp(status.AllHelp, help.ContactSupportHelp())
	}

	if !is.Success(sts) {
		sts.Log()
	}

	return sts
}


func (me *Channels) StopHandler(client messages.MessageAddress)  {

	topicStop := messages.Topic{
		Address: client,
		SubTopic: "exit",
	}
	fmt.Printf("StopHandler:'%s'\n", topicStop.String())
	me.instance.emitter.Off(topicStop.String())

	return
}


func (me *SubTopics) StopHandler()  {

	topicStop := messages.Topic{
		Address: me.Address,
		SubTopic: "exit",
	}
	fmt.Printf("StopHandler:'%s'\n", topicStop.String())
	me.instance.emitter.Off(topicStop.String())

	return
}


func (me *Channels) StartHandler(client messages.MessageAddress) (SubTopics, status.Status) {

	var sts status.Status

	for range only.Once {
		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		if me.Subscribers == nil {
			me.Subscribers = make(Subscribers)
		}

		if _, ok := me.Subscribers[client]; !ok {
			me.Subscribers[client] = SubTopics{
				Address: client,
				Callbacks: make(Callbacks),
				instance: &me.instance,
			}
		}

		go func() {
			sts = me.handler(client)
			sts.Log()
		}()

		sts = status.Success("").
			SetMessage("started event handler for %s", client.String()).
			SetAdditional("emitter:%v", me.instance.emitter).
			SetData(me).
			SetHelp(status.AllHelp, help.ContactSupportHelp())
	}

	if !is.Success(sts) {
		sts.Log()
	}

	// Save last state.
	me.Sts = sts

	if _, ok := me.Subscribers[client]; ok {
		return me.Subscribers[client], sts
	} else {
		return SubTopics{}, sts
	}
}


func (me *Channels) handler(client messages.MessageAddress) status.Status {

	var sts status.Status

	for range only.Once {
		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		//wgChannel := make(chan int)
		//var wg sync.WaitGroup
		child := 0

		messages.Debug("GBevents - Poller started '%s'.", client.String())
		topicGlob := messages.CreateTopicGlob(client).String()
		topicExit := messages.CreateTopic(client, "exit").String()

		fmt.Printf("topicCheck:'%s'\n", topicGlob)
		for me.instance.events = range me.instance.emitter.On(topicGlob) {
			if me.instance.events.Args == nil {
				fmt.Printf("ARGS:ZERO\n")
				continue
			}

			// Only one message ever sent.
			msg := me.instance.events.Args[0].(messages.Message)

			fmt.Printf("Event(%s)	Time:%v	Src:%s	Text:%s\n", msg.Topic, msg.Time.Convert().Unix(), msg.Source, msg.Text)
			if me.instance.events.OriginalTopic == topicExit { //} && (msg.Text.String() == me.EntityId.String()) {
				fmt.Printf("EXIT TIME: %s => %s\n", me.instance.events.OriginalTopic, topicGlob)
				me.instance.emitter.Off(topicGlob)
			}

			// Always replace topic with the correct one. Never trust calling entity.
			msg.Topic = messages.StringToTopic(me.instance.events.OriginalTopic)

			// Split topic from the /address/topic format
			topicAddress := msg.Topic.Address
			topic := msg.Topic.SubTopic
			//fmt.Printf(">>>>>> '%v' ==  d:%v => t:%v\n", msg.Topic, topicAddress, topic)

			if sub, ok := me.Subscribers[topicAddress]; ok {

				// Now check topics the subscriber is subscribed to, else continue to next.
				if _, ok := sub.Callbacks[topic]; ok {
					//fmt.Printf("LOOP:[%d]\n", child)

					callback := sub.Callbacks[topic]
					if callback == nil {
						continue
					}
					//fmt.Printf(">>>>>> callback:%v\n", callback)

					// Execute callback in thread.
					go func(c int) {
						//defer wg.Done()
						// fmt.Printf("Callback(%s)	Time:%v	Src:%s	Text:%s\n", msg.Topic, msg.Time.Convert().Unix(), msg.Src, msg.Text)
						sts := callback(&msg)
						//fmt.Printf("Done:[%d]\n", c)
						if is.Error(sts) {
							sts.Log()
						}
						//wgChannel <- c
					}(child)
					//wg.Add(1)
					child++
				}
			}
		}

		//fmt.Printf("WAIT\n")
		//debug.PrintStack()
		//wg.Wait()

		messages.Debug("GBevents - Poller stopped '%s'.", client.String())
		sts = status.Success("").
			SetMessage("GBevents - Poller stopped '%s'.", client.String()).
			SetAdditional("emitter:%v", me.instance.emitter).
			SetData("").
			SetHelp(status.AllHelp, help.ContactSupportHelp())

		// Remove client from map.
		delete(me.Subscribers, client)
	}

	if !is.Success(sts) {
		sts.Log()
	}

	// Save last state.
	me.Sts = sts

	return sts
}


func (me *Channels) Publish(msg messages.Message) status.Status {

	var sts status.Status

	for range only.Once {
		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		sts = msg.Topic.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		if msg.Time.IsNil() {
			msg.Time = msg.Time.Now()
		}

		if msg.Source.IsNil() {
			msg.Source = me.EntityId
		}

		if msg.Topic.Address.IsNil() {
			sts = status.Fail().
				SetMessage("no destination for channel message").
				SetAdditional("", ).
				SetData(msg).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		if msg.Text == "" {
			sts = status.Warn("").
				SetMessage("not sending empty channel message").
				SetAdditional("", ).
				SetData(msg).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		// fmt.Printf("Publish(%s) =>\n\tmsg.CreateTopic():%v\n\tme.instance.emitter:%v\n\n", msg.Topic.String(), msg, me.instance.emitter)
		me.instance.emits = me.instance.emitter.Emit(msg.Topic.String(), msg)
		if me.instance.emits == nil {
			sts = status.Fail().
				SetMessage("failed to send channel message").
				SetAdditional("", ).
				SetData(msg).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		/*
			select {
				case <-me.emits:
					sts = status.Success("channel message sent OK")

				case <-time.After(time.Second * 10):
					sts = status.Fail().
						SetMessage("timeout sending channel message").
						SetAdditional("", ).
						SetData(msg).
						SetHelp(status.AllHelp, help.ContactSupportHelp())
					close(me.emits)
			}
		*/

		sts = status.Success("").
			SetMessage("channel message sent OK").
			SetAdditional("emitter:%v", me.instance.emitter).
			SetData(msg).
			SetHelp(status.AllHelp, help.ContactSupportHelp())
	}

	if !is.Success(sts) {
		sts.Log()
	}

	// Save last state.
	me.Sts = sts

	return sts
}


func (me *Channels) Listeners(topic messages.Topic)  {
	fmt.Printf("Listeners\n")

	foo := me.instance.emitter.Listeners(topic.String())[0]

	fmt.Printf("%v", foo)
	for f := range foo {
		fmt.Printf("[%s] - '%s' '%s' '%s'\n", f.Topic, f.OriginalTopic, f.Args, f.Flags)
	}

	return
}


func (me *Channels) off(topic messages.Topic, channels ...<-chan emitter.Event)  {
	fmt.Printf("Off\n")

	me.instance.emitter.Off(topic.String(), channels...)

	return
}


func (me *Channels) on(topic messages.Topic, middleware ...func(emitter *emitter.Event)) <-chan emitter.Event {
	fmt.Printf("On\n")

	// me.instance.events = <-me.instance.emitter.On(topic.String(), middleware...)
	// me.group.Add(me.instance.emitter.On(topic.String()))

	return me.instance.emitter.On(topic.String(), middleware...)
}


func (me *Channels) once(topic messages.Topic, middleware ...func(emitter *emitter.Event)) <-chan emitter.Event {
	fmt.Printf("Once\n")

	// me.instance.events = <-me.instance.emitter.Once(topic.String(), middleware...)
	// me.instance.events.String(1)

	return me.instance.emitter.Once(topic.String(), middleware...)
}


func (me *Channels) Topics() (topics messages.Topics) {
	fmt.Printf("Topics\n")

	for _, t := range me.instance.emitter.Topics() {
		topics = append(topics, messages.StringToTopic(t))
	}

	return
}


func (me *Channels) Use(pattern string, middleware ...func(emitter *emitter.Event))  {
	fmt.Printf("Use\n")

	me.instance.emitter.Use(pattern, middleware...)

	return
}


////////////////////////////////////////////////////////////////////////////////
/*
func (me *Channels) allHandler() status.Status {

	var sts status.Status

	for range only.Once {
		sts = me.EnsureNotNil()
		if is.Error(sts) {
			break
		}

		wgChannel := make(chan int)
		var wg sync.WaitGroup
		child := 0

		messages.Debug("GBevents - Poller(STARTED)")
		for me.instance.events = range me.instance.emitter.On("*") {
			if me.instance.events.Args == nil {
				fmt.Printf("ARGS:ZERO\n")
				continue
			}

			// Only one message ever sent.
			msg := me.instance.events.Args[0].(messages.Message)

			fmt.Printf("Event(%s)	Time:%v	Src:%s	Text:%s\n", msg.Topic, msg.Time.Convert().Unix(), msg.Source, msg.Text)
			if me.instance.events.OriginalTopic == "exit" {
				fmt.Printf("EXIT TIME: %s\n", me.instance.events.OriginalTopic)
				me.instance.emitter.Off("exit")
			}

			// Split topic from the /address/topic format
			//topicAddress, topic := msg.Topic.Split()
			//fmt.Printf(">>>>>> '%v' ==  d:%v => t:%v\n", msg.Topic, topicAddress, topic)

			// Range through all subscribers
			for _, subscriber := range me.Subscribers {
				// If there's no subscriber, continue to the next.
				if subscriber == IsEmptySubScriber {
					continue
				}
				fmt.Printf(">>>>>> subscriber:%v\n", subscriber)

				// If empty topic, continue to next.
				if me.instance.events.OriginalTopic == "" {
					continue
				}

				// Always replace topic with the correct one. Never trust calling entity.
				msg.Topic = messages.StringToTopic(me.instance.events.OriginalTopic)
				fmt.Printf(">>>>>> msg.Topic:%v\n", msg.Topic)

				// Split topic from the /address/topic format
				topicAddress := msg.Topic.Address
				topic := msg.Topic.SubTopic
				fmt.Printf(">>>>>> d:%v => t:%v\n", topicAddress, topic)

				// Now check topics the subscriber is subscribed to, else continue to next.
				callback := subscriber[topic]
				if callback == nil {
					continue
				}
				fmt.Printf(">>>>>> callback:%v\n", callback)

				wg.Add(1)
				child++
				// Execute callback in thread.
				go func() {
					defer wg.Done()
					// fmt.Printf("Callback(%s)	Time:%v	Src:%s	Text:%s\n", msg.Topic, msg.Time.Convert().Unix(), msg.Src, msg.Text)
					sts := callback(&msg)
					if is.Error(sts) {
						sts.Log()
					}
					wgChannel <- child
				}()
				fmt.Printf(">>>>>> F4\n")
			}
		}

		fmt.Printf("WAIT\n")
		//debug.PrintStack()
		wg.Wait()

		sts = status.Success("GBevents - Poller(FINISHED)")
	}

	if !is.Success(sts) {
		sts.Log()
	}

	// Save last state.
	me.Sts = sts

	return sts
}
*/

/*
func (me *Channels) Test() {
		argh := ServiceEvents{
			Name: "mqtt",
			State: "running",
			Action: ServiceAction{
				State: "running",
				CallBack: hello,
			},
		}

		argh.CreateService()
		foo := ServiceEvents{}
		foo.State.String()


	fmt.Printf("DEBUG - 1\n")
	e1 := &emitter.Emitter{}
	e2 := &emitter.Emitter{}
	e3 := &emitter.Emitter{}


	go func() {
		fmt.Printf("poller1()\n")
		poller1(e1, e2, e3)
	}()

	go func() {
		fmt.Printf("poller2()\n")
		pollerA(e1, "/channel/testme")
	}()

	go func() {
		fmt.Printf("poller3()\n")
		pollerB(e2)
	}()


	go func() {
		time.Sleep(time.Millisecond * 1500)

		e1.Emit("change", messages.Message{Src: "1", Text: "A"})
		time.Sleep(time.Millisecond * 100)

		e1.Emit("nope", messages.Message{Src: "1", Text: "B"})
		time.Sleep(time.Millisecond * 100)

		e1.Emit("change", messages.Message{Src: "1", Text: "C"})
		time.Sleep(time.Millisecond * 100)

		e1.Emit("change", messages.Message{Src: "1", Text: "D"})

		time.Sleep(time.Second * 10)
		e1.Off("*") // unsubscribe any listeners
	}()


	fmt.Printf("DEBUG - 2\n")
	go func() {
		time.Sleep(time.Second * 2)

		e1.Emit("change", messages.Message{Src: "2", Text: "A"})
		time.Sleep(time.Millisecond * 100)

		e1.Emit("*", messages.Message{Text: "B", Src: "2"})
		time.Sleep(time.Millisecond * 100)

		e1.Emit("change", messages.Message{Src: "2", Text: "C"})

		time.Sleep(time.Second * 10)
		e2.Off("*") // unsubscribe any listeners
	}()


	fmt.Printf("DEBUG - 3\n")
	go func() {
		time.Sleep(time.Millisecond * 1100)

		e3.Emit("change", messages.Message{Src: "3", Text: "A"})
		time.Sleep(time.Millisecond * 100)

		e3.Emit("*", messages.Message{Src: "3", Text: "B"})
		time.Sleep(time.Millisecond * 100)

		e3.Emit("change", messages.Message{Src: "3", Text: "C"})

		time.Sleep(time.Second * 10)
		e3.Off("*") // unsubscribe any listeners
	}()


	// listener channel was closed
	fmt.Printf("DEBUG - 4\n")

	time.Sleep(time.Hour)

	fmt.Printf("DEBUG - ENTRY\n")
	config, err := broker.ConfigureConfig(os.Args[1:])
	if err != nil {
		log.Fatal("configure broker config error: ", err)
	}

	b, err := broker.NewBroker(config)
	if err != nil {
		log.Fatal("New Broker error: ", err)
	}
	b.Start()
	fmt.Printf("DEBUG - START\n")

	s := daemon.WaitForSignal()
	log.Println("signal received, broker closed.", s)

}
*/

func poller1(e1 *emitter.Emitter, e2 *emitter.Emitter, e3 *emitter.Emitter) {
	g := &emitter.Group{Cap: 1}
	g.Add(e1.On("*"), e2.On("*"), e3.On("*"))
	for event := range g.On() {
		if event.Args == nil {
			continue
		}

		foo := event.Args[0].(messages.Message)
		fmt.Printf("%s -> %s (%d): Event1(%s)\n", foo.Source, foo.Text, foo.Time.Convert().Unix(), event.OriginalTopic)
	}
	fmt.Printf("Event1(FINISHED)\n")
}

func pollerA(e *emitter.Emitter, topic string) {
	for event := range e.On(topic) {
		if event.Args == nil {
			continue
		}

		foo := event.Args[0].(messages.Message)
		fmt.Printf("%s -> %s (%d): Event2(%s)\n", foo.Source, foo.Text, foo.Time.Convert().Unix(), event.OriginalTopic)
	}
	fmt.Printf("Event2(FINISHED)\n")
}

func pollerB(e *emitter.Emitter, topic string) {
	fmt.Printf("Event3(STARTED)\n")
	for event := range e.On(topic) {
		if event.Args == nil {
			continue
		}

		foo := event.Args[0].(messages.Message)
		fmt.Printf("%s -> %s (%d): Event3(%s)\n", foo.Source, foo.Text, foo.Time.Convert().Unix(), event.OriginalTopic)
		//		go func() {
		//e.Emit("change", "...")
		//		}()
	}
	fmt.Printf("Event3(FINISHED)\n")
}

