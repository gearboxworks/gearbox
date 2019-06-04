package network

import (
	"gearbox/app/logger"
	"gearbox/heartbeat/gbevents/channels"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/only"
)


////////////////////////////////////////////////////////////////////////////////
// Executed from a channel

// Non-exposed channel function that responds to a "status" channel request.
// Produces the status of the M-DNS handler via a channel.
func statusHandler(event *messages.Message, i channels.Argument) channels.Return {

	var err error
	var me *ZeroConf

	for range only.Once {
		me, err = InterfaceToTypeZeroConf(i)
		if err != nil {
			break
		}

		logger.Debug("ZeroConf %s handler status OK", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return err
}

// Non-exposed channel function that responds to an "stop" channel request.
// Causes the M-DNS handler task to stop via a channel.
func stopHandler(event *messages.Message, i channels.Argument) channels.Return {

	var err error
	var me *ZeroConf

	for range only.Once {
		me, err = InterfaceToTypeZeroConf(i)
		if err != nil {
			break
		}

		err = me.StopHandler()
		if err != nil {
			break
		}

		logger.Debug("ZeroConf %s handler stopped OK", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return err
}


// Non-exposed channel function that responds to an "start" channel request.
// Causes the M-DNS handler task to start via a channel.
func startHandler(event *messages.Message, i channels.Argument) channels.Return {

	var err error
	var me *ZeroConf

	for range only.Once {
		me, err = InterfaceToTypeZeroConf(i)
		if err != nil {
			break
		}

		err = me.StartHandler()
		if err != nil {
			break
		}

		logger.Debug("ZeroConf %s handler started OK", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return err
}
