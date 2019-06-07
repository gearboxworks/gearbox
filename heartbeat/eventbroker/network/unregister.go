package network

import (
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/only"
)


////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

// Unregister a service by *Service method reference.
func (me *Service) Unregister() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.instance.Shutdown()
		eblog.Debug("ZeroConf unregister service (%s).", me.EntityId)
		me.instance = nil
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
	}

	return err
}

// Unregister a service by method defined by a UUID reference.
func (me *ZeroConf) UnregisterByUuid(u messages.MessageAddress) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if _, ok := me.services[u]; !ok {
			err = me.EntityId.ProduceError("no service defined")
			break
		}

		err = me.services[u].EnsureNotNil()
		if err != nil {
			break
		}

		me.services[u].instance.Shutdown()
		eblog.Debug("ZeroConf %s unregister via UUID (%s).", me.EntityId.String(), u.String())
		delete(me.services, u)
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
	}

	return err
}

// Unregister a service via a channel defined by a UUID reference.
func (me *ZeroConf) UnregisterByChannel(caller messages.MessageAddress, u messages.MessageAddress) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		//unreg := me.EntityId.Construct(me.EntityId, states.ActionUnregister, messages.MessageText(u.String()))
		unreg := caller.ConstructMessage(me.EntityId, states.ActionUnregister, messages.MessageText(u.String()))
		err = me.Channels.Publish(unreg)
		if err != nil {
			break
		}

		eblog.Debug("ZeroConf %s unregister via channel (%s).", me.EntityId.String(), u.String())
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
	}

	return err
}


////////////////////////////////////////////////////////////////////////////////
// Executed from a channel.

// Non-exposed channel function that responds to an "unregister" channel request.
func unregisterService(event *messages.Message, i channels.Argument) channels.Return {

	var me *ZeroConf
	var err error

	for range only.Once {
		me, err = InterfaceToTypeZeroConf(i)
		if err != nil {
			break
		}

		//fmt.Printf("MESSAGE Rx:\n[%v]\n", event.Text.String())

		// Use message element as the UUID.
		err = me.UnregisterByUuid(event.Text.ToUuid())
		if err != nil {
			break
		}

		eblog.Debug("ZeroConf %s unregistered service OK", me.EntityId.String())
	}

	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
	}

	return err
}

