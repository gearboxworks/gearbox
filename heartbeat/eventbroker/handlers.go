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
		me.Channels.PublishState(me.State)
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
		me.Channels.PublishState(me.State)

		err = me.channelHandler.StopHandler()
		if err != nil {
			break
		}

		me.State.SetNewState(states.StateStopped, err)
		me.Channels.PublishState(me.State)
		eblog.Debug(me.EntityId, "task handler stopped OK")
	}

	eblog.LogIfError(me.EntityId, err)
	eblog.LogIfNil(me, err)

	return err
}


// Non-exposed channel function that responds to a "get" channel request.
// Wraps a status request into another entity.
func getHandler(event *messages.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *EventBroker
	var ret *states.Status

	for range only.Once {
		me, err = InterfaceToTypeEventBroker(i)
		if err != nil {
			break
		}

		//fmt.Printf("getHandler: %s\n", event.String())

		var msg *messages.Message
		msg, err = event.Text.ToMessage()
		if err == nil {
			//fmt.Printf("%d: msg == %s\n", time.Now().Unix(), msg.String())
		} else {
			fmt.Printf("getHandler %d: msg == %v / err == %v\n", time.Now().Unix(), msg, err)
			break
		}

		switch msg.Topic.SubTopic {
			case states.ActionStatus:
				//fmt.Printf("%d: Republish status request message: %s\n", time.Now().Unix(), msg.String())
				var ir channels.Return
				ir, err = me.Channels.PublishAndWaitForReturn(*msg, 200)
				if err == nil {
					//fmt.Printf("%d: OK - ir == %v\n", time.Now().Unix(), ir)
				} else {
					fmt.Printf("getHandler %d: ER - ir == %v /  err == %v\n", time.Now().Unix(), ir, err)
					break
				}

				ret, err = states.InterfaceToTypeStatus(ir)
				if err == nil {
					//fmt.Printf("%d: OK - ret == %s\n", time.Now().Unix(), ret.String())
				} else {
					fmt.Printf("getHandler %d: ER - ret == nil / err == %v\n", time.Now().Unix(), err)
				}
		}

		//unreg := me.EntityId.ConstructMessage(event.Text.ToMessageAddress(), states.ActionUnregister, messages.MessageText(u.String()))
		//err = me.Channels.Publish(unreg)
		//if err != nil {
		//	break
		//}

		//fmt.Printf("Event(%s) Time:%d Src:%s Text:%s\n", event.Topic.String(), event.Time.Convert().Unix(), event.Source.String(), event.Text.String())

		eblog.Debug(me.EntityId, "getHandler() via channel")
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

		//fmt.Printf("statusHandler: %s\n", event.String())
		//fmt.Printf("Event(%s) Time:%d Src:%s Text:%s\n", event.Topic.String(), event.Time.Convert().Unix(), event.Source.String(), event.Text.String())

		if event.Topic.Address.String() == "" {
			break
		}

		if event.Text.String() == "" {
			break
		}

		ret, err = states.FromMessageText(event.Text)
		if err != nil {
			fmt.Printf("Error %v - %s\n", err, event.String())
			break
		}

		//fmt.Printf("Rec: %s\n", ret.String())
		err = me.Services.AddState(*ret)
		if err != nil {
			break
		}
		//fmt.Printf(">> %s is at state '%s'\n", ret.EntityId.String(), ret.Current.String())

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

