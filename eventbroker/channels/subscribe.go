package channels

import (
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/only"
	"gearbox/eventbroker/states"
	"sync"
)


func (me *Channels) Subscribe(client messages.MessageTopic, callback Callback, argInterface Argument, retType ReturnType) (*Subscriber, error) {

	var err error
	var sub Subscriber

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = client.EnsureNotNil()
		if err != nil {
			break
		}

		if callback == nil {
			err = me.EntityId.ProduceError("callback function is empty")
			break
		}

		if retType == "" {
			err = me.EntityId.ProduceError("return type is empty")
			break
		}

		if me.subscribers == nil {
			me.subscribers = make(Subscribers)
		}

		if _, ok := me.subscribers[client.Address]; !ok {
			addr := client.Address
			sub = Subscriber{
				EntityId:  addr,
				EntityName: addr,
				EntityParent: &me.EntityId,
				State: states.New(&addr, &addr, me.EntityId),
				IsManaged: true,

				topics: make(References),
				mutex: sync.RWMutex{},
				parentInstance: &me.instance,
			}
			me.subscribers[client.Address] = &sub
		}

		me.subscribers[client.Address].State.SetNewAction(states.ActionSubscribe)

		me.subscribers[client.Address].AddTopic(client.SubTopic, callback, argInterface, retType)

		me.subscribers[client.Address].State.SetNewState(states.StateSubscribed, err)
		eblog.Debug(me.EntityId, "channel subscriber: %s", messages.SprintfTopic(client.Address, client.SubTopic))
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return &sub, err
}


func (me *Subscriber) Subscribe(topic messages.SubTopic, callback Callback, argInterface Argument, retType ReturnType) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for range only.Once {
			err = topic.EnsureNotNil()
			if err != nil {
				break
			}

			if callback == nil {
				err = me.EntityId.ProduceError("callback function is empty")
				break
			}

			if retType == "" {
				err = me.EntityId.ProduceError("return type is empty")
				break
			}

			me.State.SetNewAction(states.ActionSubscribe)

			me.AddTopic(topic, callback, argInterface, retType)

			me.State.SetNewState(states.StateSubscribed, err)
			eblog.Debug(me.EntityId, "channel new subscriber: %s", messages.SprintfTopic(me.EntityId, topic))
		}
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}

