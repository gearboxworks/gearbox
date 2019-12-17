package network

import (
	"fmt"
	"gearbox/eventbroker/channels"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/states"
	"github.com/gearboxworks/go-status/only"
)

////////////////////////////////////////////////////////////////////////////////
// Executed from a channel

// Non-exposed channel function that responds to an "stop" channel request.
func stopHandler(event *msgs.Message, i channels.Argument, r channels.ReturnType) channels.Return {

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

		eblog.Debug(me.EntityId, "requested service stop via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return &err
}

// Non-exposed channel function that responds to an "start" channel request.
func startHandler(event *msgs.Message, i channels.Argument, r channels.ReturnType) channels.Return {

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

		eblog.Debug(me.EntityId, "requested service start via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return &err
}

// Non-exposed channel function that responds to a "status" channel request.
func statusHandler(event *msgs.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *ZeroConf
	var ret *states.Status

	for range only.Once {
		me, err = InterfaceToTypeZeroConf(i)
		if err != nil {
			break
		}

		if event.Text.String() == "" {
			// Get status of Daemon by default
			ret = me.State.GetStatus()
		} else {
			// Get status of specific sub
			sc := me.IsExisting(msgs.Address(event.Text))
			if sc != nil {
				ret, err = sc.GetStatus()
			}
		}

		eblog.Debug(me.EntityId, "statusHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return ret
}

// Non-exposed channel function that responds to a "register" channel request.
func registerService(event *msgs.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var me *ZeroConf
	var ret *Service
	var err error

	for range only.Once {
		me, err = InterfaceToTypeZeroConf(i)
		if err != nil {
			break
		}

		var ce ServiceConfig
		ce, err = DeconstructMdnsMessage(event)
		//err = json.Unmarshal(event.Text.ByteArray(), &ce)
		if err != nil {
			err = msgs.MakeError(me.EntityId, "cannot deconstruct MDNS message with error '%v'", err)
			break
		}

		ret, err = me.Register(ce)
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "registered service by channel %s OK", ret.EntityId.String())
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return ret
}

// Non-exposed channel function that responds to an "unregister" channel request.
func unregisterService(event *msgs.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var me *ZeroConf
	var err error

	for range only.Once {
		me, err = InterfaceToTypeZeroConf(i)
		if err != nil {
			break
		}

		// Use message element as the UUID.
		err = me.UnregisterByEntityId(event.Text.ToAddress())
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "unregistered service by channel %s OK", event.Text.ToAddress())
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return &err
}

// Non-exposed channel function that responds to a "get" channel request.
func getHandler(event *msgs.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *ZeroConf
	var ret msgs.SubTopics

	for range only.Once {
		me, err = InterfaceToTypeZeroConf(i)
		if err != nil {
			break
		}

		switch event.Text.String() {
		case "topics":
			ret = me.channelHandler.GetTopics()
		case "topics/subs":
			ret = me.channelHandler.GetTopics()
		}

		fmt.Printf("topics: %v\n", ret)

		eblog.Debug(me.EntityId, "topicsHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return &ret
}

// Non-exposed channel function that responds to a "scan" channel request.
func scanServices(event *msgs.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var me *ZeroConf
	var err error

	for range only.Once {
		me, err = InterfaceToTypeZeroConf(i)
		if err != nil {
			break
		}

		_, err = me.Browse(event.Text.String(), me.domain)
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "service scan completed")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(err)

	return &err
}
