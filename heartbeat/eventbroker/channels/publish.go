package channels

import (
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/only"
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

		me.State.Want = states.StatePublished

		//fmt.Printf("MESSAGE Tx:[%v]\n", msg)

		// eblog.Debug("Publish(%s) =>\tmsg.CreateTopic():%v\tme.instance.emitter:%v", msg.Topic.String(), msg, me.instance.emitter)
		me.instance.emits = me.instance.emitter.Emit(msg.Topic.String(), msg)
		if me.instance.emits == nil {
			err = me.EntityId.ProduceError("failed to send channel message")
			break
		}


		eblog.Debug("channel MSG:'%s' DATA:'%s'", msg.Topic.String(), msg.Text.String())
		//fmt.Printf(">>> MSG: %s DATA: %s", msg.Topic.String(), msg.Text.String())
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

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
	}

	return err
}


func (me *Channels) GetCallbackReturn(msg messages.Message, waitForExecute int) (Return, error) {

	var r Return
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
		if _, ok := me.subscribers[client].Returns[subtopic]; !ok {
			err = me.EntityId.ProduceError("channel subscriber return not defined")
			break
		}

		for loop := 0; (me.subscribers[client].Executed[subtopic] == false) && (loop < waitForExecute); loop++ {
			// Wait if we are asked to.
			time.Sleep(time.Millisecond * 10)
		}

		if me.subscribers[client].Executed[subtopic] == false {
			err = me.EntityId.ProduceError("no response from channel")
			break
		}

		r = me.subscribers[client].Returns[subtopic]
		me.subscribers[client].Executed[subtopic] = false
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
	}

	return r, err
}


func (me *Channels) SetCallbackReturnToNil(msg messages.Message) (error) {

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
		if _, ok := me.subscribers[client].Returns[subtopic]; ok {
			me.subscribers[client].Returns[subtopic] = nil
		}
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
	}

	return err
}


// Send channel message on state changes only.
func PublishCallerState(me *Channels, caller *messages.MessageAddress, state *states.Status) {

	switch {
		case me == nil:

		case caller == nil:

		case state == nil:

		case state.Error != nil:
			_ = me.Publish(caller.ConstructMessage(*caller, states.ActionError, messages.MessageText(state.Error.Error())))

		case state.Current != state.Want:
			fallthrough
		case state.Current != state.Last:
			_ = me.Publish(caller.ConstructMessage(*caller, states.ActionStatus, messages.MessageText(state.Current)))
	}

	return
}


func (me *Channels) PublishCallerState(caller *messages.MessageAddress, state *states.Status) {

	PublishCallerState(me, caller, state)

	return
}


func (me *Channels) PublishSpecificCallerState(caller *messages.MessageAddress, state states.State) {

	for range only.Once {
		switch {
			case me == nil:
			case state == "":
			case caller == nil:
		}

		PublishCallerState(me, caller, &states.Status{Current: state})
	}

	return
}

