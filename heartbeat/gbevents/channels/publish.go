package channels

import (
	"errors"
	"gearbox/app/logger"
	"gearbox/heartbeat/gbevents/messages"
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

		if msg.Source.IsNil() {
			msg.Source = me.entityId
		}

		if msg.Topic.Address.IsNil() {
			err = errors.New("no destination for channel message")
			break
		}

		//fmt.Printf("MESSAGE Tx:[%v]\n", msg)

		// logger.Debug("Publish(%s) =>\tmsg.CreateTopic():%v\tme.instance.emitter:%v", msg.Topic.String(), msg, me.instance.emitter)
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
		logger.Debug("Error: %v", err)
	}

	// Save last state.
	me.Error = err
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
			err = errors.New("unknown channel subscriber")
			break
		}

		subtopic := msg.Topic.SubTopic
		if _, ok := me.subscribers[client].Returns[subtopic]; !ok {
			err = errors.New("channel subscriber return not defined")
			break
		}

		for loop := 0; (me.subscribers[client].Executed[subtopic] == false) && (loop < waitForExecute); loop++ {
			// Wait if we are asked to.
			time.Sleep(time.Millisecond * 10)
		}

		if me.subscribers[client].Executed[subtopic] == false {
			err = errors.New("no response from channel")
			break
		}

		r = me.subscribers[client].Returns[subtopic]
		me.subscribers[client].Executed[subtopic] = false
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
			err = errors.New("unknown channel subscriber")
			break
		}

		subtopic := msg.Topic.SubTopic
		if _, ok := me.subscribers[client].Returns[subtopic]; ok {
			me.subscribers[client].Returns[subtopic] = nil
		}
	}

	return err
}

