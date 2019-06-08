package mqttClient

import (
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/only"
)


////////////////////////////////////////////////////////////////////////////////
// Executed from a channel

// Non-exposed channel function that responds to a "status" channel request.
// Produces the status of the M-DNS handler via a channel.
func statusHandler(event *messages.Message, i channels.Argument) channels.Return {

	var err error
	var me *MqttClient

	for range only.Once {
		me, err = InterfaceToTypeMqttClient(i)
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "requested service status via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


// Non-exposed channel function that responds to an "stop" channel request.
// Causes the M-DNS handler task to stop via a channel.
func stopHandler(event *messages.Message, i channels.Argument) channels.Return {

	var err error
	var me *MqttClient

	for range only.Once {
		me, err = InterfaceToTypeMqttClient(i)
		if err != nil {
			break
		}

		err = me.StopHandler()
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "requested service stop via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


// Non-exposed channel function that responds to an "start" channel request.
// Causes the M-DNS handler task to start via a channel.
func startHandler(event *messages.Message, i channels.Argument) channels.Return {

	var err error
	var me *MqttClient

	for range only.Once {
		me, err = InterfaceToTypeMqttClient(i)
		if err != nil {
			break
		}

		err = me.StartHandler()
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "requested service start via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}
