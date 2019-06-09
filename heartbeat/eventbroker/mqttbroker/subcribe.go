package mqttBroker


////////////////////////////////////////////////////////////////////////////////
// Executed as a method.
//
//// Register a service by method defined by a *CreateTopic structure.
//func (me *MqttBroker) Subscribe(ce CreateEntry) (*Service, error) {
//
//	var err error
//	var sc Service
//
//	for range only.Once {
//		err = me.EnsureNotNil()
//		if err != nil {
//			break
//		}
//
//		err = me.services.IsExisting(ce)
//		if err != nil {
//			break
//		}
//
//		// Create new client entry.
//		for range only.Once {
//			me.State.SetNewAction(states.ActionSubscribe)
//			sc.EntityId = messages.GenerateAddress()
//			sc.IsManaged = true
//			sc.channels = me.Channels
//			channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
//
//			if ce.callback == nil {
//				ce.callback = defaultCallback
//			}
//
//			if ce.Topic.String() == "" {
//				// Nope, not gonna do it.
//				err = me.EntityId.ProduceError("empty topic")
//				break
//			}
//
//			sc.instance = me.instance.client.Subscribe(ce.Topic.String(), ce.Qos, ce.callback)
//			if sc.instance == nil {
//				err = me.EntityId.ProduceError("unable to subscribe")
//				break
//			}
//
//			me.services[sc.EntityId] = &sc
//
//			eblog.Debug(me.EntityId, "registered service %s OK", sc.EntityId.String())
//		}
//
//		sc.State.SetNewState(states.StateSubscribed, err)
//		sc.channels.PublishCallerState(&sc.EntityId, &sc.State)
//	}
//
//	eblog.LogIfNil(me, err)
//	eblog.LogIfError(me.EntityId, err)
//
//	return &sc, err
//}
//
//
//// Subscribe a service via a channel defined by a *CreateTopic structure and
//// returns a *Service structure if successful.
//func (me *MqttBroker) SubscribeByChannel(caller messages.MessageAddress, s Topic) (*Service, error) {
//
//	var err error
//	var sc *Service
//
//	for range only.Once {
//		err = me.EnsureNotNil()
//		if err != nil {
//			break
//		}
//
//		if s == "" {
//			err = me.EntityId.ProduceError("unable to subscribe")
//			break
//		}
//
//		reg := caller.ConstructMessage(me.EntityId, states.ActionSubscribe, messages.MessageText(s))
//		err = me.Channels.Publish(reg)
//		if err != nil {
//			break
//		}
//
//		rs, err := me.Channels.GetCallbackReturn(reg, 100)
//		if err != nil {
//			break
//		}
//
//		sc, err = InterfaceToTypeService(rs)	// sc = rs.(*Service)
//		if err != nil {
//			break
//		}
//
//		eblog.Debug(me.EntityId, "subscribed by channel %s OK", sc.EntityId.String())
//	}
//
//	eblog.LogIfNil(me, err)
//	eblog.LogIfError(me.EntityId, err)
//
//	return sc, err
//}
//
//
//////////////////////////////////////////////////////////////////////////////////
//// Executed from a channel.
//
//// Non-exposed channel function that responds to a "register" channel request.
//func subscribeTopic(event *messages.Message, i channels.Argument) channels.Return {
//
//	var me *MqttBroker
//	var sc *Service
//	var err error
//
//	for range only.Once {
//		me, err = InterfaceToTypeMqttBroker(i)
//		if err != nil {
//			break
//		}
//
//		//fmt.Printf("Rx: %v\n", event)
//
//		var ce CreateEntry
//		err = json.Unmarshal(event.Text.ByteArray(), &ce)
//
//		sc, err = me.Subscribe(ce)
//		if err != nil {
//			break
//		}
//
//		eblog.Debug(me.EntityId, "subscribed by channel %s OK", sc.EntityId.String())
//	}
//
//	eblog.LogIfNil(me, err)
//	eblog.LogIfError(me.EntityId, err)
//
//	return sc
//}
//
//
//
//func foo2(client mqtt.Client, msg mqtt.Message) {
//	fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
//}
//
//
//func defaultCallback(client mqtt.Client, msg mqtt.Message) {
//
//	fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
//}
//
