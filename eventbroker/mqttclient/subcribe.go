package mqttClient

import (
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/states"
	"github.com/gearboxworks/go-status/only"
)

////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

// Register a service by method defined by a *NewTopic structure.
func (me *MqttClient) Subscribe(ce ServiceConfig) (*Service, error) {

	var err error
	var sc Service

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = me.services.IsExisting(ce)
		if err != nil {
			break
		}

		// Create new client entry.
		sc.EntityId = msgs.MakeAddress()
		sc.EntityName = msgs.Address(ce.Name)
		sc.EntityParent = &me.EntityId
		sc.State = states.New(sc.EntityId, sc.EntityName, me.EntityId)
		sc.State.SetNewAction(states.ActionSubscribe)
		sc.IsManaged = true
		sc.channels = me.Channels
		sc.channels.PublishState(sc.State)

		if ce.callback == nil {
			ce.callback = defaultCallback
		}

		if ce.Topic.String() == "" {
			// Nope, not gonna do it.
			err = msgs.MakeError(me.EntityId, "empty topic")
			break
		}

		sc.instance = me.instance.client.Subscribe(ce.Topic.String(), ce.Qos, ce.callback)
		if sc.instance == nil {
			err = msgs.MakeError(me.EntityId, "unable to subscribe")
			break
		}

		err = me.AddEntity(sc.EntityId, &sc)
		if err != nil {
			break
		}

		sc.State.SetNewState(states.StateSubscribed, err)
		sc.channels.PublishState(sc.State)
		eblog.Debug(me.EntityId, "subscribed %s OK", sc.EntityId.String())
	}

	me.Channels.PublishState(me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return &sc, err
}

// Subscribe a service via a channel defined by a *NewTopic structure and
// returns a *Service structure if successful.
func (me *MqttClient) SubscribeByChannel(caller msgs.Address, s Topic) (*Service, error) {

	var err error
	var sc *Service

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if s == "" {
			err = msgs.MakeError(me.EntityId, "unable to subscribe")
			break
		}

		reg := caller.MakeMessage(me.EntityId, states.ActionSubscribe, msgs.Text(s))
		err = me.Channels.Publish(reg)
		if err != nil {
			break
		}

		rs, err := me.Channels.GetCallbackReturn(reg, 100)
		if err != nil {
			break
		}

		sc, err = InterfaceToTypeService(rs) // sc = rs.(*Service)
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "subscribed by channel %s OK", sc.EntityId.String())
	}

	me.Channels.PublishState(me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return sc, err
}
