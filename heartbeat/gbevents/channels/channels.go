package channels

import (
	"fmt"
	"gearbox/app/logger"
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

		_args.osSupport = OsSupport
		foo := box.Args{}
		err = copier.Copy(&foo, &_args)
		if err != nil {
			break
		}

		if _args.entityId == "" {
			_args.entityId = DefaultEntityId
		}

		_args.instance.emitter = emitter.Emitter{}
		_args.subscribers = make(Subscribers)

		*me = Channels(_args)

		logger.Debug("New Channel instance (%s).", me.entityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	// Save last state.
	me.Error = err
	return err
}


func (me *Channels) StopHandler(client messages.MessageAddress)  {

	topicStop := messages.Topic{
		Address: client,
		SubTopic: DefaultExitString,
	}
	logger.Debug("StopHandler:'%s'", topicStop.String())
	me.instance.emitter.Off(topicStop.String())

	return
}


func (me *Subscriber) StopHandler() error {

	topicStop := messages.Topic{
		Address: me.Address,
		SubTopic: DefaultExitString,
	}
	logger.Debug("StopHandler:'%s'", topicStop.String())
	me.instance.emitter.Off(topicStop.String())

	return nil
}


func (me *Channels) StartHandler(client messages.MessageAddress) (*Subscriber, error) {

	var err error

	fmt.Print("HELLO\n")
	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if me.subscribers == nil {
			me.subscribers = make(Subscribers)
		}

		if _, ok := me.subscribers[client]; !ok {
			sh := Subscriber{
				Address: client,
				Callbacks: make(Callbacks),
				Arguments: make(Arguments),
				Returns: make(Returns),
				Executed: make(Executed),
				instance: &me.instance,
			}
			me.subscribers[client] = &sh
		}

		go func() {
			err = me.rxHandler(client)
			if err != nil {
				logger.Debug("GBevents - handler errored '%v'.", err)
			}
		}()

		logger.Debug("started channel event handler for %s", client.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}
	// Save last state.
	me.Error = err

	if _, ok := me.subscribers[client]; ok {
		return me.subscribers[client], err
	} else {
		var empty Subscriber
		return &empty, err
	}
}


func (me *Channels) rxHandler(client messages.MessageAddress) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		//wgChannel := make(chan int)
		//var wg sync.WaitGroup
		child := 0

		logger.Debug("channels handler started '%s'.", client.String())
		topicGlob := messages.CreateTopicGlob(client).String()
		topicExit := messages.CreateTopic(client, DefaultExitString).String()

		for me.instance.events = range me.instance.emitter.On(topicGlob) {
			if me.instance.events.Args == nil {
				logger.Debug("channels handler saw zero args")
				continue
			}

			// Only one message ever sent.
			msg := me.instance.events.Args[0].(messages.Message)

			logger.Debug("Event(%s) Time:%d Src:%s Text:%s", msg.Topic.String(), msg.Time.Convert().Unix(), msg.Source.String(), msg.Text.String())
			if me.instance.events.OriginalTopic == topicExit { //} && (msg.Text.String() == me.entityId.String()) {
				logger.Debug("EXIT TIME: %s => %s", me.instance.events.OriginalTopic, topicGlob)
				me.instance.emitter.Off(topicGlob)
			}

			// Always replace topic with the correct one. Never trust calling entity.
			msg.Topic = messages.StringToTopic(me.instance.events.OriginalTopic)

			// Split topic from the /address/topic format
			topicAddress := msg.Topic.Address
			topic := msg.Topic.SubTopic

			if sub, ok := me.subscribers[topicAddress]; ok {

				// Now check topics the subscriber is subscribed to, else continue to next.
				if _, ok := sub.Callbacks[topic]; !ok {
					// No callback defined, ignore.
					continue
				}

				if sub.Callbacks[topic] == nil {
					continue
				}

				if _, ok := sub.Arguments[topic]; !ok {
					sub.Arguments[topic] = nil
				}

				//logger.Debug("LOOP:[%d]", child)
				// Execute callback in thread.
				go func(c int) {
					//defer wg.Done()
					// logger.Debug("Callback(%s)	Time:%v	Src:%s	Text:%s", msg.Topic, msg.Time.Convert().Unix(), msg.Src, msg.Text)
					sub.Executed[topic] = false

					if _, ok := sub.Returns[topic]; ok {
						sub.Returns[topic] = sub.Callbacks[topic](&msg, sub.Arguments[topic])
						//logger.Debug("# Return1: %v", sub.Returns[topic])
					} else {
						_ = sub.Callbacks[topic](&msg, sub.Arguments[topic])
						//logger.Debug("# Return2: %v", sub.Returns[topic])
					}
					sub.Executed[topic] = true

					//wgChannel <- c
				}(child)
				//wg.Add(1)
				child++
			}
		}

		//logger.Debug("WAIT")
		//debug.PrintStack()
		//wg.Wait()

		logger.Debug("channels handler stopped '%s'.", client.String())

		// Remove client from map.
		delete(me.subscribers, client)
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	// Save last state.
	me.Error = err
	return err
}


/*
func (me *Channels) txHandler(client messages.MessageAddress) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		//wgChannel := make(chan int)
		//var wg sync.WaitGroup
		child := 0


		for me.instance.events = range me.instance.emitter.On(topicGlob) {
			if me.instance.events.Args == nil {
				logger.Debug("channels handler saw zero args")
				continue
			}

			// Only one message ever sent.
			msg := me.instance.events.Args[0].(messages.Message)

			logger.Debug("Event(%s)	Time:%v	Src:%s	Text:%s", msg.Topic.String(), msg.Time.Convert().Unix(), msg.Source, msg.Text)
			if me.instance.events.OriginalTopic == topicExit { //} && (msg.Text.String() == me.entityId.String()) {
				logger.Debug("EXIT TIME: %s => %s", me.instance.events.OriginalTopic, topicGlob)
				me.instance.emitter.Off(topicGlob)
			}

			// Always replace topic with the correct one. Never trust calling entity.
			msg.Topic = messages.StringToTopic(me.instance.events.OriginalTopic)

			// Split topic from the /address/topic format
			topicAddress := msg.Topic.Address
			topic := msg.Topic.SubTopic

			if sub, ok := me.subscribers[topicAddress]; ok {

				// Now check topics the subscriber is subscribed to, else continue to next.
				if _, ok := sub.Callbacks[topic]; ok {
					//logger.Debug("LOOP:[%d]", child)

					if sub.Callbacks[topic] == nil {
						continue
					}

					if _, ok := sub.Interfaces[topic]; !ok {
						sub.Interfaces[topic] = nil
					}

					// Execute callback in thread.
					go func(c int) {
						//defer wg.Done()
						// logger.Debug("Callback(%s)	Time:%v	Src:%s	Text:%s", msg.Topic, msg.Time.Convert().Unix(), msg.Src, msg.Text)
						sub.Errors[topic] = sub.Callbacks[topic](&msg, sub.Interfaces[topic])
						//logger.Debug("Done:[%d]", c)
						if sub.Errors[topic] != nil {
							logger.Debug("Error: ", err)
						}
						//wgChannel <- c
					}(child)
					//wg.Add(1)
					child++
				}
			}
		}

		//logger.Debug("WAIT")
		//debug.PrintStack()
		//wg.Wait()

		logger.Debug("channels handler stopped '%s'.", client.String())

		// Remove client from map.
		delete(me.subscribers, client)
	}

	if err != nil {
		logger.Debug("Error: ", err)
	}

	// Save last state.
	me.Error = err
	return err
}
*/


func (me *Channels) Listeners(topic messages.Topic)  {
	fmt.Printf("Listeners\n")

	foo := me.instance.emitter.Listeners(topic.String())[0]

	for f := range foo {
		fmt.Printf("[%s] - '%s' '%s' '%s'\n", f.Topic, f.OriginalTopic, f.Args, f.Flags)
	}

	return
}


func (me *Channels) Topics() (topics messages.Topics) {
	logger.Debug("Topics")

	for _, t := range me.instance.emitter.Topics() {
		topics = append(topics, messages.StringToTopic(t))
	}

	return
}


func (me *Channels) GetId() messages.MessageAddress {

	if me == nil {
		return messages.MessageAddress("")
	}

	return me.entityId
}


func (me *Channels) GetSubscribers() *Subscribers {

	empty := Subscribers{}

	if me == nil {
		return &empty
	}

	return &me.subscribers
}


func (me *Channels) ListSubscribers() {

	if me == nil {
		return
	}

	me.subscribers.List()
}


func (me *Channels) off(topic messages.Topic, channels ...<-chan emitter.Event)  {
	logger.Debug("Off")

	me.instance.emitter.Off(topic.String(), channels...)

	return
}


func (me *Channels) on(topic messages.Topic, middleware ...func(emitter *emitter.Event)) <-chan emitter.Event {
	logger.Debug("On")

	// me.instance.events = <-me.instance.emitter.On(topic.String(), middleware...)
	// me.group.Add(me.instance.emitter.On(topic.String()))

	return me.instance.emitter.On(topic.String(), middleware...)
}


func (me *Channels) once(topic messages.Topic, middleware ...func(emitter *emitter.Event)) <-chan emitter.Event {
	logger.Debug("Once")

	// me.instance.events = <-me.instance.emitter.Once(topic.String(), middleware...)
	// me.instance.events.String(1)

	return me.instance.emitter.Once(topic.String(), middleware...)
}


func (me *Channels) use(pattern string, middleware ...func(emitter *emitter.Event))  {
	logger.Debug("Use")

	me.instance.emitter.Use(pattern, middleware...)

	return
}

