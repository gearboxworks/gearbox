package channels

import (
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/entity"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/only"
	"gearbox/heartbeat/eventbroker/states"
	"time"
)


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

		if msg.Source.EnsureNotNil() != nil {
			msg.Source = me.EntityId
		}

		if msg.Topic.Address.EnsureNotNil() != nil {
			err = me.EntityId.ProduceError("no destination for channel message")
			break
		}

		// eblog.Debug(me.EntityId, "Publish(%s) =>\tmsg.CreateTopic():%v\tme.instance.emitter:%v", msg.Topic.String(), msg, me.instance.emitter)
		me.instance.emits = me.instance.emitter.Emit(msg.Topic.String(), msg)
		if me.instance.emits == nil {
			err = me.EntityId.ProduceError("failed to send channel message")
			break
		}

		eblog.Debug(me.EntityId, "Channel time:%d src:%s topic:%s msg:%s", msg.Time.Unix(), msg.Source.String(), msg.Topic.String(), msg.Text.String())
		/*
			select {
				case <-me.emits:
					// err = me.EntityId.ProduceError("channel message sent OK")

				case <-time.After(time.Second * 10):
					err = me.EntityId.ProduceError("timeout sending channel message")
					close(me.emits)
			}
		*/
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Channels) GetCallbackReturn(msg messages.Message, waitForExecute int) (Return, error) {

	var ret Return
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		client := msg.Topic.Address
		if _, ok := me.subscribers[client]; !ok {
			err = me.EntityId.ProduceError("unknown channel subscriber")
			break
		}

		subtopic := msg.Topic.SubTopic
		// MUTEX if _, ok := me.subscribers[client].Returns[subtopic]; !ok {
		err, _, _, _, _ = me.subscribers[client].GetTopic(subtopic)
		if err != nil {
			break
		}

		// MUTEX for loop := 0; (me.subscribers[client].Executed[subtopic] == false) && (loop < waitForExecute); loop++ {
		for loop := 0; (me.subscribers[client].GetExecuted(subtopic) == false) && (loop < waitForExecute); loop++ {
			// Wait if we are asked to.
			time.Sleep(time.Millisecond * 10)
		}

		// MUTEX if me.subscribers[client].Executed[subtopic] == false {
		if me.subscribers[client].GetExecuted(subtopic) == false {
			err = me.EntityId.ProduceError("no response from channel")
			break
		}

		// MUTEX ret = me.subscribers[client].Returns[subtopic]
		ret = me.subscribers[client].GetReturns(subtopic)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return ret, err
}


func (me *Channels) SetCallbackReturnToNil(msg messages.Message) (error) {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = me.EnsureSubscriberNotNil(msg.Topic.Address)
		if err != nil {
			break
		}

		err, _, _, _, _ = me.subscribers[msg.Topic.Address].GetTopic(msg.Topic.SubTopic)
		if err != nil {
			break
		}

		// MUTEX me.subscribers[client].Returns[subtopic] = nil
		me.subscribers[msg.Topic.Address].SetReturns(msg.Topic.SubTopic, nil)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *Channels) PublishAndWaitForReturn(msg messages.Message, waitForExecute int) (Return, error) {

	var err error
	var ret Return

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = me.Publish(msg)
		if err != nil {
			break
		}

		ret, err = me.GetCallbackReturn(msg, waitForExecute)
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "message returned by channel OK")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return ret, err
}


// Send channel message on state changes only.
func PublishState(me *Channels, caller *messages.MessageAddress, state *states.Status) {

	switch {
		case me == nil:

		case caller == nil:

		case state == nil:

		case state.GetError() != nil:
			msg := caller.ConstructMessage(entity.BroadcastEntityName, states.ActionError, messages.MessageText(state.GetError().Error()))
			//fmt.Printf("ERROR: %s\n", msg.String())
			_ = me.Publish(msg)

		case state.ExpectingNewState():
			fallthrough
		case state.HasChangedState():
			msg := caller.ConstructMessage(entity.BroadcastEntityName, states.ActionStatus, messages.MessageText(state.GetCurrent()))
			//fmt.Printf("EXPECTING: %s\n", msg.String())
			_ = me.Publish(msg)
	}

	return
}


func (me *Channels) PublishState(caller *messages.MessageAddress, state *states.Status) {

	PublishState(me, caller, state)

	return
}


func (me *Channels) PublishSpecificState(caller *messages.MessageAddress, state states.State) {

	switch {
		case me == nil:
		case state == "":
		case caller == nil:
			return
	}

	PublishState(me, caller, &states.Status{Current: state})

	return
}


// Send channel message on state changes only.
func PublishCallerState(me *Channels, state *states.Status) {

	switch {
		case me == nil:

		case state.EntityId == nil:

		case state == nil:

		case state.GetError() != nil:
			// _ = me.Publish(caller.ConstructMessage(*caller, states.ActionError, messages.MessageText(state.GetError().Error())))
			msg := state.EntityId.ConstructMessage(entity.BroadcastEntityName, states.ActionError, messages.MessageText(state.GetError().Error()))
			//fmt.Printf("ERROR: %s\n", msg.String())
			_ = me.Publish(msg)

		case state.ExpectingNewState():
			fallthrough
		case state.HasChangedState():
			// _ = me.Publish(caller.ConstructMessage(*caller, states.ActionStatus, messages.MessageText(state.GetCurrent())))
			msg := state.EntityId.ConstructMessage(entity.BroadcastEntityName, states.ActionStatus, messages.MessageText(state.GetCurrent()))
			//fmt.Printf("EXPECTING: %s\n", msg.String())
			_ = me.Publish(msg)
	}

	return
}


func (me *Channels) PublishCallerState(state *states.Status) {

	PublishCallerState(me, state)

	return
}

