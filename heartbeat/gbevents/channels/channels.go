package channels

import (
	"errors"
	"fmt"
	"gearbox/box"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/only"
	oss "gearbox/os_support"
	"github.com/jinzhu/copier"
	"github.com/olebedev/emitter"
)


func (me *Channels) New(OsSupport oss.OsSupporter, args ...Args) error {

	var _args Args
	var err error

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		_args.OsSupport = OsSupport
		foo := box.Args{}
		err = copier.Copy(&foo, &_args)
		if err != nil {
			break
		}

		if _args.EntityId == "" {
			_args.EntityId = DefaultEntityId
		}

		_args.instance.emitter = emitter.Emitter{}
		_args.Subscribers = make(Subscribers)

		*me = Channels(_args)

		messages.Debug("New Channel instance (%s).", me.EntityId.String())
	}

	if err != nil {
		messages.Debug("Error: %v", err)
	}

	// Save last state.
	me.Error = err
	return err
}


func (me *Channels) Subscribe(topic messages.Topic, fn Callback, i interface{}) (SubTopics, error) {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = topic.EnsureNotNil()
		if err != nil {
			break
		}

		if fn == nil {
			err = errors.New("channel callback function is empty")
			break
		}

		if me.Subscribers == nil {
			me.Subscribers = make(Subscribers)
		}

		if _, ok := me.Subscribers[topic.Address]; !ok {
			me.Subscribers[topic.Address] = SubTopics{
				Address: topic.Address,
				Callbacks: make(Callbacks),
				Interfaces: make(Interfaces),
				instance: &me.instance,
			}
		}

		me.Subscribers[topic.Address].Callbacks[topic.SubTopic] = fn
		me.Subscribers[topic.Address].Interfaces[topic.SubTopic] = i
		foo := me.Subscribers[topic.Address]
		foo.List()
		me.Subscribers.List()

		messages.Debug("New subscriber: %s\n", messages.SprintfTopic(topic.Address, topic.SubTopic))
	}

	if err != nil {
		messages.Debug("Error: %v", err)
	}
	// Save last state.
	me.Error = err

	if _, ok := me.Subscribers[topic.Address]; ok {
		return me.Subscribers[topic.Address], err
	} else {
		return SubTopics{}, err
	}
}


func (me *Channels) UnSubscribe(topic messages.Topic) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = topic.EnsureNotNil()
		if err != nil {
			break
		}

		delete(me.Subscribers[topic.Address].Callbacks, topic.SubTopic)
		delete(me.Subscribers[topic.Address].Interfaces, topic.SubTopic)
		messages.Debug("Unsubscribed: %s\n", messages.SprintfTopic(topic.Address, topic.SubTopic))
	}

	if err != nil {
		messages.Debug("Error: %v", err)
	}

	// Save last state.
	me.Error = err
	return err
}


func (me *SubTopics) Subscribe(topic messages.SubTopic, fn Callback, i interface{}) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = topic.EnsureNotNil()
		if err != nil {
			break
		}

		if fn == nil {
			err = errors.New("channel callback function is empty")
			break
		}

		me.Callbacks[topic] = fn
		me.Interfaces[topic] = i
		me.List()

		messages.Debug("New subscriber: %s\n", messages.SprintfTopic(me.Address, topic))
	}

	if err != nil {
		messages.Debug("Error: %v", err)
	}

	return err
}


func (me *SubTopics) UnSubscribe(topic messages.SubTopic) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = topic.EnsureNotNil()
		if err != nil {
			break
		}

		delete(me.Callbacks, topic)
		delete(me.Interfaces, topic)
		messages.Debug("Unsubscribed: %s\n", messages.SprintfTopic(me.Address, topic))
	}

	if err != nil {
		messages.Debug("Error: %v", err)
	}

	return err
}


func (me *Channels) StopHandler(client messages.MessageAddress)  {

	topicStop := messages.Topic{
		Address: client,
		SubTopic: "exit",
	}
	messages.Debug("StopHandler:'%s'\n", topicStop.String())
	me.instance.emitter.Off(topicStop.String())

	return
}


func (me *SubTopics) StopHandler()  {

	topicStop := messages.Topic{
		Address: me.Address,
		SubTopic: "exit",
	}
	messages.Debug("StopHandler:'%s'\n", topicStop.String())
	me.instance.emitter.Off(topicStop.String())

	return
}


func (me *Channels) StartHandler(client messages.MessageAddress) (SubTopics, error) {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if me.Subscribers == nil {
			me.Subscribers = make(Subscribers)
		}

		if _, ok := me.Subscribers[client]; !ok {
			me.Subscribers[client] = SubTopics{
				Address: client,
				Callbacks: make(Callbacks),
				Interfaces: make(Interfaces),
				instance: &me.instance,
			}
		}

		go func() {
			err = me.handler(client)
			if err != nil {
				messages.Debug("GBevents - handler errored '%v'.", err)
			}
		}()

		messages.Debug("started channel event handler for %s", client.String())
	}

	if err != nil {
		messages.Debug("Error: %v", err)
	}
	// Save last state.
	me.Error = err

	if _, ok := me.Subscribers[client]; ok {
		return me.Subscribers[client], err
	} else {
		return SubTopics{}, err
	}
}


func (me *Channels) handler(client messages.MessageAddress) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
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

					if sub.Callbacks[topic] == nil {
						continue
					}

					if _, ok := sub.Interfaces[topic]; !ok {
						sub.Interfaces[topic] = nil
					}

					// Execute callback in thread.
					go func(c int) {
						//defer wg.Done()
						// fmt.Printf("Callback(%s)	Time:%v	Src:%s	Text:%s\n", msg.Topic, msg.Time.Convert().Unix(), msg.Src, msg.Text)
						err := sub.Callbacks[topic](&msg, sub.Interfaces[topic])
						//fmt.Printf("Done:[%d]\n", c)
						if err != nil {
							messages.Debug("Error: ", err)
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

		// Remove client from map.
		delete(me.Subscribers, client)
	}

	if err != nil {
		messages.Debug("Error: ", err)
	}

	// Save last state.
	me.Error = err
	return err
}


func (me *Channels) Publish(msg messages.Message) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = msg.Topic.EnsureNotNil()
		if err != nil {
			break
		}

		if msg.Time.IsNil() {
			msg.Time = msg.Time.Now()
		}

		if msg.Source.IsNil() {
			msg.Source = me.EntityId
		}

		if msg.Topic.Address.IsNil() {
			err = errors.New("no destination for channel message")
			break
		}

		// fmt.Printf("Publish(%s) =>\n\tmsg.CreateTopic():%v\n\tme.instance.emitter:%v\n\n", msg.Topic.String(), msg, me.instance.emitter)
		me.instance.emits = me.instance.emitter.Emit(msg.Topic.String(), msg)
		if me.instance.emits == nil {
			err = errors.New("failed to send channel message")
		}

		/*
			select {
				case <-me.emits:
					// err = errors.New("channel message sent OK")

				case <-time.After(time.Second * 10):
					err = errors.New("timeout sending channel message")
					close(me.emits)
			}
		*/
	}

	if err != nil {
		messages.Debug("Error: ", err)
	}

	// Save last state.
	me.Error = err
	return err
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
func (me *Channels) allHandler() error {

	var err error

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

		err = errors.New("GBevents - Poller(FINISHED)")
	}

	if !is.Success(sts) {
		sts.Log()
	}

	// Save last state.
	me.Sts = sts

	return sts
}
*/
