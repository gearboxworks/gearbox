package vmbox


////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

//// Unsubscribe a service by method defined by a UUID reference.
//func (me *VmBox) UnsubscribeByUuid(client messages.MessageAddress) error {
//
//	var err error
//
//	for range only.Once {
//		err = me.EnsureNotNil()
//		if err != nil {
//			break
//		}
//
//		err = me.services[client].EnsureNotNil()
//		if err != nil {
//			break
//		}
//
//		me.services[client].State.SetNewAction(states.ActionStop)	// Was states.ActionUnsubscribe
//		me.services[client].channels.PublishState(me.State)
//
//		// Do something
//
//		me.services[client].State.SetNewState(states.StateStopped, err)	// Was states.StateUnsubscribed
//		me.services[client].channels.PublishState(me.services[client].State)
//
//		err = me.DeleteEntity(client)
//		if err != nil {
//			break
//		}
//
//		//me.Channels.PublishSpecificState(&client, states.State(states.StateUnsubscribed))
//		eblog.Debug(me.EntityId, "unregistered service %s OK", client.String())
//	}
//
//	me.Channels.PublishState(me.State)
//	eblog.LogIfNil(me, err)
//	eblog.LogIfError(me.EntityId, err)
//
//	return err
//}
//
//// Unsubscribe a service via a channel defined by a UUID reference.
//func (me *VmBox) UnsubscribeByChannel(caller messages.MessageAddress, u messages.MessageAddress) error {
//
//	var err error
//
//	for range only.Once {
//		err = me.EnsureNotNil()
//		if err != nil {
//			break
//		}
//
//		//unreg := me.EntityId.Construct(me.EntityId, states.ActionUnsubscribe, messages.MessageText(u.String()))
//		unreg := caller.ConstructMessage(me.EntityId, states.ActionUnsubscribe, messages.MessageText(u.String()))
//		err = me.Channels.Publish(unreg)
//		if err != nil {
//			break
//		}
//
//		eblog.Debug(me.EntityId, "unsubscribed service by channel %s OK", u.String())
//	}
//
//	eblog.LogIfNil(me, err)
//	eblog.LogIfError(me.EntityId, err)
//
//	return err
//}
//
//// Register a service by method defined by a *CreateTopic structure.
//func (me *VmBox) Subscribe(ce ServiceConfig) (*Service, error) {
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
//		sc.EntityId = *messages.GenerateAddress()
//		sc.EntityName = messages.MessageAddress(ce.Name)
//		sc.EntityParent = &me.EntityId
//		sc.State = states.New(&sc.EntityId, &sc.EntityName, me.EntityId)
//		sc.State.SetNewAction(states.ActionSubscribe)
//		sc.IsManaged = true
//		sc.channels = me.Channels
//		sc.channels.PublishState(sc.State)
//
//		err = me.AddEntity(sc.EntityId, &sc)
//		if err != nil {
//			break
//		}
//
//		sc.State.SetNewState(states.StateSubscribed, err)
//		sc.channels.PublishState(sc.State)
//		eblog.Debug(me.EntityId, "subscribed %s OK", sc.EntityId.String())
//	}
//
//	me.Channels.PublishState(me.State)
//	eblog.LogIfNil(me, err)
//	eblog.LogIfError(me.EntityId, err)
//
//	return &sc, err
//}

