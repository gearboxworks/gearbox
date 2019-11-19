package daemon

import (
	"encoding/json"
	"fmt"
	"gearbox/eventbroker/channels"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/states"
	"github.com/gearboxworks/go-status/only"
)

////////////////////////////////////////////////////////////////////////////////
// Executed from a channel


// Non-exposed channel function that responds to an "stop" channel request.
func stopHandler(event *messages.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *Daemon

	for range only.Once {
		me, err = InterfaceToTypeDaemon(i)
		if err != nil {
			break
		}

		if event.Text.String() == "" {
			break
		}

		if event.Text.String() == entity.SelfEntityName {
			// Stop Daemon by default
			err = me.StopHandler()
		} else {
			// Stop of specific entity
			sc := me.IsExisting(messages.MessageAddress(event.Text))
			if sc != nil {
				err = sc.Stop()
			}
		}

		//err = me.StopHandler()
		//if err != nil {
		//	break
		//}

		eblog.Debug(me.EntityId, "stopHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return &err
}


// Non-exposed channel function that responds to an "start" channel request.
func startHandler(event *messages.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *Daemon

	for range only.Once {
		me, err = InterfaceToTypeDaemon(i)
		if err != nil {
			break
		}

		if event.Text.String() == "" {
			break
		}

		if event.Text.String() == entity.SelfEntityName {
			// Start Daemon by default
			err = me.StartHandler()
		} else {
			// Start of specific entity
			sc := me.IsExisting(messages.MessageAddress(event.Text))
			if sc != nil {
				err = sc.Start()
			}
		}

		//err = me.StartHandler()
		//if err != nil {
		//	break
		//}

		eblog.Debug(me.EntityId, "startHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return &err
}


// Non-exposed channel function that responds to a "status" channel request.
func statusHandler(event *messages.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *Daemon
	var ret *states.Status

	for range only.Once {
		me, err = InterfaceToTypeDaemon(i)
		if err != nil {
			break
		}

		if event.Text.String() == "" {
			// Get status of Daemon by default
			ret = me.State.GetStatus()
		} else {
			// Get status of specific sub
			sc := me.IsExisting(messages.MessageAddress(event.Text))
			if sc != nil {
				ret, err = sc.GetStatus()
			}
		}

		eblog.Debug(me.EntityId, "statusHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return ret
}


// Non-exposed channel function that responds to a "register" channel request.
func registerService(event *messages.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var me *Daemon
	var sc *Service
	var err error

	for range only.Once {
		me, err = InterfaceToTypeDaemon(i)
		if err != nil {
			break
		}

		//fmt.Printf("Rx: %v\n", event)

		ce := ServiceConfig{}
		err = json.Unmarshal(event.Text.ByteArray(), &ce)
		if err != nil {
			break
		}

		sc, err = me.Register(ce)
		if err != nil {
			break
		}
		if sc == nil {
			break
		}

		eblog.Debug(me.EntityId, "registered service by channel %s OK", sc.EntityId.String())
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return sc
}


// Non-exposed channel function that responds to an "unregister" channel request.
func unregisterService(event *messages.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var me *Daemon
	var err error

	for range only.Once {
		me, err = InterfaceToTypeDaemon(i)
		if err != nil {
			break
		}

		//fmt.Printf("MESSAGE Rx:\n[%v]\n", event.Text.String())

		// Use message element as the UUID.
		err = me.UnregisterByEntityId(event.Text.ToMessageAddress())
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "unregistered service by channel %s OK", event.Text.ToMessageAddress())
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return &err
}


// Non-exposed channel function that responds to a "get" channel request.
func getHandler(event *messages.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *Daemon
	var ret messages.SubTopics

	for range only.Once {
		me, err = InterfaceToTypeDaemon(i)
		if err != nil {
			break
		}
		fmt.Printf("ReturnType: %v\n", r)

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
	eblog.LogIfError(me.EntityId, err)

	return &ret
}


// Non-exposed channel function that responds to a "load" channel request.
func loadConfigHandler(event *messages.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *Daemon

	for range only.Once {
		me, err = InterfaceToTypeDaemon(i)
		if err != nil {
			break
		}

		err = me.LoadServiceFiles()

		eblog.Debug(me.EntityId, "loadConfigHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return &err
}

