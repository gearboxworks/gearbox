package eventbroker

import (
	"errors"
	"fmt"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/entity"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/only"
	"gearbox/heartbeat/eventbroker/states"
	"sort"
	"sync"
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
// This expects the message text to contain an embedded status request message.
// Thus exposes entity status' to the outside.
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
	var state *states.Status
	var sc *Service
	var ok bool


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

		state, err = states.FromMessageText(event.Text)
		if err != nil {
			fmt.Printf("Error %v - %s\n", err, event.String())
			break
		}


		// Create callback reference if not already present.
		sc, ok, err = me.AttachCallback(*state.EntityName, nil, nil)
		if err != nil {
			break
		}

		if !ok {
			// If it was already there, only update if there's a change.
			ok = sc.IsTheSame(*state)
			if ok {
				break
			}
		}

		err = sc.updateState(*state)
		if err != nil {
			break
		}

		err = sc.processCallback(*state)
		if err != nil {
			break
		}

		//fmt.Printf(">> %s is at state '%s'\n", state.EntityId.String(), state.Current.String())

		eblog.Debug(me.EntityId, "statusHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return state
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


func (me Services) Exists(client messages.MessageAddress) (*Service, error) {

	var err error
	var ret *Service

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = client.EnsureNotNil()
		if err != nil {
			break
		}

		if _, ok := me[client]; !ok {
			break
		}

		err = me[client].EnsureNotNil()
		if err != nil {
			break
		}

		ret = me[client]
	}

	return ret, err
}


func (me Services) LookFor(client messages.MessageAddress) (*Service, error) {

	var err error
	var ret *Service

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = client.EnsureNotNil()
		if err != nil {
			break
		}

		var keys messages.MessageAddress
		for keys, ret = range me {
			if keys == client {
				break
			}

			if *ret.State.EntityId == client {
				break
			}

			if *ret.State.EntityName == client {
				break
			}
		}

		//err = me.Exists(keys)
	}

	return ret, err
}


func (me *EventBroker) AttachCallback(client messages.MessageAddress, cb Callback, args interface{}) (*Service, bool, error) {

	var err error
	var ok bool
	var ret *Service

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = client.EnsureNotNil()
		if err != nil {
			break
		}

		ret, err = me.Services.Exists(client)
		if err != nil {
			break
		}
		if ret != nil {
			// If map entry exists, just update cb.
			//
			//me[client].mutex.Lock()
			//defer me[client].mutex.Unlock()
			//
			//me[client].Callback = cb
			break
		}

		ret = &Service{
			State: states.New(&client, &client, entity.BroadcastEntityName),
			Callback: cb,
			Args: args,
			Logs: make(Logs, 0),
			mutex: sync.RWMutex{},
		}

		me.Services[client] = ret
		ok = true
	}

	return ret, ok, err
}


func (me Services) PrintStates() error {

	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("status.Status is nil")
			break
		}

		//me.mutex.RLock()
		//defer me.mutex.RUnlock()
		var keys []string
		for k := range me {
			keys = append(keys, k.String())
		}
		sort.Strings(keys)

		for i, e := range keys {
			fmt.Printf("%d %s\n", i, me[messages.MessageAddress(e)].PrintState())
		}
	}

	return err
}


func (me Services) EnsureNotNil() error {

	var err error

	if me == nil {
		err = errors.New("services is nil")
	}

	return err
}



func (me *Service) updateState(state states.Status) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.mutex.Lock()
		defer me.mutex.Unlock()

		me.State = &state

		//me.Logs = append(me.Logs, Log{
		//	State: state,
		//	When: time.Now(),
		//})
	}

	return err
}


func (me *Service) IsTheSame(state states.Status) bool {

	var err error
	var ok bool

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if state.EntityId == nil {
			err = errors.New("status.EntityId is nil")
			break
		}

		if me.State.Current != state.Current {
			break
		}
		if me.State.Want != state.Want {
			break
		}
		if me.State.Last != state.Last {
			break
		}
		if me.State.LastWhen != state.LastWhen {
			break
		}
		if me.State.Attempts != state.Attempts {
			break
		}
		if me.State.Error != state.Error {
			break
		}
		if me.State.Action != state.Action {
			break
		}

		ok = true
	}

	return ok
}


func (me *Service) processCallback(state states.Status) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if me.Callback == nil {
			// Don't fire, if we don't have a CB defined.
			break
		}


		// Ensure we process in the correct order.
		me.mutex.Lock()
		err = me.Callback(me.Args, state)
		me.mutex.Unlock()
	}

	return err
}


func (me *Service) PrintState() string {

	var err error
	var ret string

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.mutex.RLock()
		defer me.mutex.RUnlock()
		fmt.Printf("%s", me.State.String())
	}

	return ret
}


func (me *Service) EnsureNotNil() error {

	var err error

	if me == nil {
		err = errors.New("service is nil")
	}

	return err
}



func (me Callback) EnsureNotNil() error {

	var err error

	if me == nil {
		err = errors.New("callback is nil")
	}

	return err
}


//func (me Services) DeleteState(client messages.SubTopic) error {
//
//	var err error
//
//	me.mutex.Lock()
//	defer me.mutex.Unlock()
//
//	for range only.Once {
//		_, ok := me.topics[client] // Managed by Mutex
//		if !ok {
//			err = me.EntityId.ProduceError("service doesn't exist")
//			break
//		}
//
//		delete(me.topics, client) // Managed by Mutex
//	}
//
//	return err
//}

