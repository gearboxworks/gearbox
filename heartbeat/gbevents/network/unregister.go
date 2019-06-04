package network

import (
	"errors"
	"gearbox/app/logger"
	"gearbox/heartbeat/gbevents/channels"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/only"
	"github.com/google/uuid"
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
		logger.Debug("ZeroConf unregister service (%s).", me.EntityId)
		me.instance = nil
	}

	return err
}

// Unregister a service by method defined by a UUID reference.
func (me *ZeroConf) UnregisterByUuid(u uuid.UUID) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if _, ok := me.services[u]; !ok {
			err = errors.New("no service defined")
			break
		}

		err = me.services[u].EnsureNotNil()
		if err != nil {
			break
		}

		me.services[u].instance.Shutdown()
		logger.Debug("ZeroConf %s unregister via UUID (%s).", me.EntityId.String(), u.String())
		delete(me.services, u)
	}

	return err
}

// Unregister a service via a channel defined by a UUID reference.
func (me *ZeroConf) UnregisterByChannel(caller messages.MessageAddress, u uuid.UUID) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		//unreg := me.EntityId.Construct(me.EntityId, messages.SubTopicUnregister, messages.MessageText(u.String()))
		unreg := caller.ConstructMessage(me.EntityId, messages.SubTopicUnregister, messages.MessageText(u.String()))
		err = me.Channels.Publish(unreg)
		if err != nil {
			break
		}

		logger.Debug("ZeroConf %s unregister via channel (%s).", me.EntityId.String(), u.String())
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

		// Use message element as the UUID.
		var u uuid.UUID
		//fmt.Printf("MESSAGE Rx:\n[%v]\n", event.Text.String())
		u, err = uuid.Parse(event.Text.String())
		if err != nil {
			logger.Debug("ZeroConf %s unregistered service OK", me.EntityId.String())
			break
		}

		err = me.UnregisterByUuid(u)
		if err != nil {
			break
		}

		logger.Debug("ZeroConf %s unregistered service OK", me.EntityId.String())
		err = nil
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return err
}

