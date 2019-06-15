package channels

import (
	"gearbox/box"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/states"
	"gearbox/eventbroker/only"
	"github.com/jinzhu/copier"
	"github.com/olebedev/emitter"
	"sync"
)


func (me *Channels) New(args ...Args) error {

	var _args Args
	var err error

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

		if _args.OsPaths == nil {
			err = me.EntityId.ProduceError("ospaths is nil")
			break
		}


		if _args.EntityId == "" {
			_args.EntityId = entity.ChannelEntityName
		}
		_args.State = states.New(&_args.EntityId, &_args.EntityId, entity.SelfEntityName)

		if _args.Boxname == "" {
			_args.Boxname = entity.ChannelEntityName
		}

		_args.instance.emitter = &emitter.Emitter{}
		_args.subscribers = make(Subscribers)

		*me = Channels(_args)


		me.State.SetWant(states.StateIdle)
		me.State.SetNewState(states.StateIdle, err)
		eblog.Debug(me.EntityId, "init complete")
	}

	me.PublishState(me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Channels) StartHandler() error {

	// Just a stub function.
	var err error

	//for range only.Once {
	//	err = me.EnsureNotNil()
	//	if err != nil {
	//		break
	//	}
	//
	//me.State.SetNewState(states.StateStarting, err)
	//PublishCallerState(me.Channels, &me.EntityId, &me.State)
	//
	//	for range only.Once {
	//		me.Task, err = tasks.StartTask(initDaemon, startDaemon, monitorDaemon, stopDaemon, me)
	//		if err != nil {
	//			break
	//		}
	//	}
	//
	//	if me.State.SetNewState(states.StateStarted, err) {
	//		eblog.Debug(me.EntityId, "started task handler")
	//	}
	//}
	//
	//channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
	//eblog.LogIfError(me, err)

	return err
}


func (me *Channels) StopHandler() error {

	var err error

	for n, s := range me.subscribers {
		err = s.StopHandler()
		if err != nil {
			eblog.Debug(me.EntityId, "channel '%s' stopped OK", n)
		} else {
			eblog.Debug(me.EntityId, "channel '%s' didn't stop", n)
		}
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Channels) StopClientHandler(client messages.MessageAddress)  {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		topicStop := client.CreateTopic(states.ActionStop)
		//topicStop := messages.MessageTopic{
		//	Address:  client,
		//	SubTopic: messages.SubTopicStop,
		//}

		eblog.Debug(me.EntityId, "StopHandler:'%s'", topicStop.String())
		me.instance.emitter.Off(topicStop.String())
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return
}


func (me *Subscriber) StopHandler() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		topicStop := me.EntityId.CreateTopic(states.ActionStop)

		eblog.Debug(me.EntityId, "StopHandler:'%s'", topicStop.String())
		me.parentInstance.emitter.Off(topicStop.String())
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return nil
}


func (me *Channels) StartClientHandler(client messages.MessageAddress) (*Subscriber, error) {

	var err error
	var sub Subscriber

	for range only.Once {
		err = EnsureNotNil(me)
		if err != nil {
			break
		}

		err = client.EnsureNotNil()
		if err != nil {
			break
		}

		if me.subscribers == nil {
			me.subscribers = make(Subscribers)
		}

		sub = Subscriber{
			EntityId:  client,
			EntityName:  client,
			EntityParent: &me.EntityId,
			State: states.New(&client, &client, me.EntityId),
			IsManaged: true,

			topics: make(References),
			mutex: sync.RWMutex{},
			parentInstance: &me.instance,
		}
		err = me.AddEntity(client, &sub)
		if err != nil {
			break
		}

		go func() {
			err = me.rxHandler(client)
			if err != nil {
				eblog.Debug(me.EntityId, "GBevents - handler errored '%v'.", err)
			}
		}()

		eblog.Debug(me.EntityId, "started channel event handler for %s", client.String())
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return &sub, err
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

		eblog.Debug(me.EntityId, "channels handler started '%s'.", client.String())
		topicGlob := client.CreateTopicGlob().String()
		topicExit := client.CreateTopic(states.ActionStop).String()

		for me.instance.events = range me.instance.emitter.On(topicGlob) {
			if me.instance.events.Args == nil {
				eblog.Debug(me.EntityId, "channels handler saw zero args")
				continue
			}

			// Only one message ever sent.
			msg := me.instance.events.Args[0].(messages.Message)

			eblog.Debug(me.EntityId, "Event(%s) Time:%d Src:%s Text:%s", msg.Topic.String(), msg.Time.Convert().Unix(), msg.Source.String(), msg.Text.String())
			if me.instance.events.OriginalTopic == topicExit { //} && (msg.Text.String() == me.EntityId.String()) {
				eblog.Debug(me.EntityId, "EXIT TIME: %s => %s", me.instance.events.OriginalTopic, topicGlob)
				me.instance.emitter.Off(topicGlob)
			}

			// Always replace topic with the correct one. Never trust calling entity.
			msg.Topic = messages.StringToTopic(me.instance.events.OriginalTopic)

			// Split topic from the /address/topic format
			client := msg.Topic.Address
			topic := msg.Topic.SubTopic

			if sub, ok := me.subscribers[client]; ok {

				// Now check topics the subscriber is subscribed to, else continue to next.
				err, callback, args, ret, retType := me.subscribers[client].GetTopic(topic)
				if err != nil {
					continue
				}

				// Execute callback in thread.
				go func(c int) {
					//defer wg.Done()
					// eblog.Debug(me.EntityId, "Callback(%s)	Time:%v	Src:%s	Text:%s", msg.Topic, msg.Time.Convert().Unix(), msg.Src, msg.Text)
					//fmt.Printf("AAARGH: %v/%v - %s\n", client, topic, msg.String())
					sub.SetExecuted(topic, false)
					if ret != nil {
						r := callback(&msg, args, retType)
						sub.SetReturns(topic, r)
					} else {
						 _ = callback(&msg, args, retType)
					}
					sub.SetExecuted(topic, true)

					//wgChannel <- c
				}(child)
				//wg.Add(1)
				child++
			}
		}

		//eblog.Debug(me.EntityId, "WAIT")
		//debug.PrintStack()
		//wg.Wait()

		eblog.Debug(me.EntityId, "channels handler stopped '%s'.", client.String())

		// Remove client from map.
		delete(me.subscribers, client)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Channels) GetEntityId() messages.MessageAddress {

	if me == nil {
		return messages.MessageAddress("")
	}

	return me.EntityId
}

