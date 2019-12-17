package channels

import (
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/states"
	"github.com/gearboxworks/go-status/only"
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

		if _args.OsDirs == nil {
			err = msgs.MakeError(me.EntityId, "ospaths is nil")
			break
		}

		if _args.EntityId == "" {
			_args.EntityId = entity.ChannelEntityName
		}
		_args.State = states.New(_args.EntityId, _args.EntityId, entity.SelfEntityName)

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
	eblog.LogIfError(err)

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
	eblog.LogIfError(err)

	return err
}

func (me *Channels) StopClientHandler(client msgs.Address) {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		topicStop := msgs.NewTopic(client, states.ActionStop)

		eblog.Debug(me.EntityId, "StopHandler:'%s'", topicStop.String())
		me.instance.emitter.Off(topicStop.String())
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return
}

func (me *Subscriber) StopHandler() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		topicStop := msgs.NewTopic(me.EntityId, states.ActionStop)

		eblog.Debug(me.EntityId, "StopHandler:'%s'", topicStop.String())
		me.parentInstance.emitter.Off(topicStop.String())
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return nil
}

func (me *Channels) StartClientHandler(client msgs.Address) (*Subscriber, error) {

	var err error
	var sub Subscriber

	for range only.Once {
		err = EnsureNotNil(me)
		if err != nil {
			break
		}

		err = client.EnsureNotEmpty()
		if err != nil {
			break
		}

		if me.subscribers == nil {
			me.subscribers = make(Subscribers)
		}

		sub = Subscriber{
			EntityId:     client,
			EntityName:   client,
			EntityParent: &me.EntityId,
			State:        states.New(client, client, me.EntityId),
			IsManaged:    true,

			topics:         make(References),
			mutex:          sync.RWMutex{},
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
	eblog.LogIfError(err)

	return &sub, err
}

func (me *Channels) rxHandler(client msgs.Address) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		child := 0

		eblog.Debug(me.EntityId, "channels handler started '%s'.", client.String())
		topicGlob := msgs.NewGlobTopic(client).String()
		topicExit := msgs.NewTopic(client, states.ActionStop).String()

		for me.instance.events = range me.instance.emitter.On(topicGlob) {
			if me.instance.events.Args == nil {
				eblog.Debug(me.EntityId, "channels handler saw zero args")
				continue
			}

			// Only one message ever sent.
			msg := me.instance.events.Args[0].(msgs.Message)

			eblog.Debug(me.EntityId, "Event(%s) Time:%d Src:%s Text:%s",
				msg.Topic.String(),
				msg.Time.Convert().Unix(),
				msg.Source.String(),
				msg.Text.String(),
			)
			if me.instance.events.OriginalTopic == topicExit {
				eblog.Debug(me.EntityId, "EXIT TIME: %s => %s", me.instance.events.OriginalTopic, topicGlob)
				me.instance.emitter.Off(topicGlob)
			}

			// Always replace topic with the correct one. Never trust calling entity.
			msg.Topic = msgs.StringToTopic(me.instance.events.OriginalTopic)

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
					sub.SetExecuted(topic, false)
					if ret != nil {
						r := callback(&msg, args, retType)
						sub.SetReturns(topic, r)
					} else {
						_ = callback(&msg, args, retType)
					}
					sub.SetExecuted(topic, true)

				}(child)
				child++
			}
		}

		eblog.Debug(me.EntityId, "channels handler stopped '%s'.", client.String())

		// Remove client from map.
		delete(me.subscribers, client)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return err
}

func (me *Channels) GetEntityId() msgs.Address {

	if me == nil {
		return msgs.Address("")
	}

	return me.EntityId
}
