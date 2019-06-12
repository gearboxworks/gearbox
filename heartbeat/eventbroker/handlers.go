package eventbroker

import (
	"fmt"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/entity"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/only"
	"gearbox/heartbeat/eventbroker/states"
	"time"
)


func (me *EventBroker) StartChannelHandler() error {

	var err error

	for range only.Once {
		me.State.SetNewAction(states.ActionInitialize)

		me.channelHandler, err = me.Channels.StartClientHandler(entity.BroadcastEntityName)
		if err != nil {
			break
		}

		err = me.channelHandler.Subscribe(states.ActionStop, stopHandler, me, states.InterfaceTypeError)
		if err != nil {
			break
		}
		err = me.channelHandler.Subscribe(states.ActionStart, startHandler, me, states.InterfaceTypeError)
		if err != nil {
			break
		}
		err = me.channelHandler.Subscribe(states.ActionStatus, statusHandler, me, states.InterfaceTypeStatus)
		if err != nil {
			break
		}
		err = me.channelHandler.Subscribe("get", getHandler, me, states.InterfaceTypeStatus)
		if err != nil {
			break
		}

		me.State.SetNewState(states.StateInitialized, err)
		me.Channels.PublishState(&me.EntityId, &me.State)
		eblog.Debug(me.EntityId, "task handler init completed OK")

	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


func (me *EventBroker) StopChannelHandler() error {

	var err error

	for range only.Once {
		me.State.SetNewAction(states.ActionStop)
		me.Channels.PublishState(&me.EntityId, &me.State)

		err = me.channelHandler.StopHandler()
		if err != nil {
			break
		}

		me.State.SetNewState(states.StateStopped, err)
		me.Channels.PublishState(&me.EntityId, &me.State)
		eblog.Debug(me.EntityId, "task handler stopped OK")
	}

	eblog.LogIfError(me.EntityId, err)
	eblog.LogIfNil(me, err)

	return err
}


// Non-exposed channel function that responds to a "stop" channel request.
func getHandler(event *messages.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *EventBroker
	var ret *states.Status

	for range only.Once {
		me, err = InterfaceToTypeEventBroker(i)
		if err != nil {
			break
		}

		fmt.Printf("Get event: %s\n", event.String())

		var msg *messages.Message
		msg, err = event.Text.ToMessage()
		fmt.Printf("%d msg == %v, err == %v\n", time.Now().Unix(), msg, err)
		if err != nil {
			break
		}

		switch msg.Topic.SubTopic {
			case states.ActionStatus:
				fmt.Printf("Republish status request message: %v\n", msg.String())
				var ir channels.Return
				ir, err = me.Channels.PublishAndWaitForReturn(*msg, 200)
				fmt.Printf("%d ir == %v, err == %v\n", time.Now().Unix(), ir, err)
				if err != nil {
					break
				}

				ret, err = states.InterfaceToTypeStatus(ir)
				fmt.Printf("%d ret == %v, err == %v\n", time.Now().Unix(), ret, err)
				if err == nil {
					fmt.Printf("%d status after: %v\n", time.Now().Unix(), ret.GetError())
					// fmt.Printf("%d gbevents after: %v (%v)\n", time.Now().Unix(), f.GetError(), me.Daemon.Fluff)
				} else {
					fmt.Printf("%d status after: is nil!\n", time.Now().Unix())
				}
		}

		//unreg := me.EntityId.ConstructMessage(event.Text.ToMessageAddress(), states.ActionUnregister, messages.MessageText(u.String()))
		//err = me.Channels.Publish(unreg)
		//if err != nil {
		//	break
		//}

		//fmt.Printf("Event(%s) Time:%d Src:%s Text:%s\n", event.Topic.String(), event.Time.Convert().Unix(), event.Source.String(), event.Text.String())

		eblog.Debug(me.EntityId, "statusHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return ret
}


// Non-exposed channel function that responds to a "status" channel request.
func statusHandler(event *messages.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *EventBroker
	var ret *states.Status

	for range only.Once {
		me, err = InterfaceToTypeEventBroker(i)
		if err != nil {
			break
		}

		//fmt.Printf("Status event: %s\n", event.String())
		//fmt.Printf("Event(%s) Time:%d Src:%s Text:%s\n", event.Topic.String(), event.Time.Convert().Unix(), event.Source.String(), event.Text.String())

		//msg := messages.Message{
		//	Source: me.EntityId,
		//	Topic: messages.MessageTopic{
		//		Address: event.Source,
		//		SubTopic: "status",
		//	},
		//	Text: "",
		//}
		//fmt.Printf("\n\n%d gbevents before: %v\n", time.Now().Unix(), me.Daemon.State.GetError())
		//i, _ := me.Channels.PublishAndWaitForReturn(msg, 400)
		//f, err := states.InterfaceToTypeStatus(i)
		//if err == nil {
		//	fmt.Printf("%d gbevents after: %v\n", time.Now().Unix(), f)
		//} else {
		//	fmt.Printf("%d gbevents after: is nil!\n", time.Now().Unix())
		//}

		eblog.Debug(me.EntityId, "statusHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return ret
}


// Non-exposed channel function that responds to a "stop" channel request.
func stopHandler(event *messages.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *EventBroker
	var ret *states.Status

	for range only.Once {
		me, err = InterfaceToTypeEventBroker(i)
		if err != nil {
			break
		}

		fmt.Printf("Stop event: %s\n", event.String())

		//t := event.Text.ToMessage()
		//fmt.Printf("Translate: %v\n", t)

		//unreg := me.EntityId.ConstructMessage(event.Text.ToMessageAddress(), states.ActionUnregister, messages.MessageText(u.String()))
		//err = me.Channels.Publish(unreg)
		//if err != nil {
		//	break
		//}

		//fmt.Printf("Event(%s) Time:%d Src:%s Text:%s\n", event.Topic.String(), event.Time.Convert().Unix(), event.Source.String(), event.Text.String())

		eblog.Debug(me.EntityId, "statusHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return ret
}


// Non-exposed channel function that responds to a "start" channel request.
func startHandler(event *messages.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *EventBroker
	var ret *states.Status

	for range only.Once {
		me, err = InterfaceToTypeEventBroker(i)
		if err != nil {
			break
		}

		fmt.Printf("Start event: %s\n", event.String())
		//fmt.Printf("Event(%s) Time:%d Src:%s Text:%s\n", event.Topic.String(), event.Time.Convert().Unix(), event.Source.String(), event.Text.String())

		eblog.Debug(me.EntityId, "statusHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return ret
}

