package daemon

import (
	"fmt"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/only"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// Executed from a channel


// Non-exposed channel function that responds to an "stop" channel request.
func stopHandler(event *messages.Message, i channels.Argument) channels.Return {

	var err error
	var me *Daemon

	for range only.Once {
		me, err = InterfaceToTypeDaemon(i)
		if err != nil {
			break
		}

		err = me.StopHandler()
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "stopHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


// Non-exposed channel function that responds to an "start" channel request.
func startHandler(event *messages.Message, i channels.Argument) channels.Return {

	var err error
	var me *Daemon

	for range only.Once {
		me, err = InterfaceToTypeDaemon(i)
		if err != nil {
			break
		}

		err = me.StartHandler()
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "startHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}


// Non-exposed channel function that responds to a "status" channel request.
func statusHandler(event *messages.Message, i channels.Argument) channels.Return {

	var err error
	var me *Daemon
	var ret *states.Status

	for range only.Once {
		me, err = InterfaceToTypeDaemon(i)
		if err != nil {
			break
		}

		fmt.Printf("%d getHandler records: %v\n", time.Now().Unix(), me.State.GetError())
		//me.Fluff = fmt.Sprintf("%d", time.Now().Unix())

		//me.State.SetError(errors.New("YES - " + me.Fluff))

		if event.Text.String() == "" {
			// Get status of Daemon by default
			ret = me.State.GetStatus()
		} else {
			// Get status of specific sub
			me.daemons.Print()
		}

		eblog.Debug(me.EntityId, "statusHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return ret
}


// Non-exposed channel function that responds to a "load" channel request.
func loadConfigHandler(event *messages.Message, i channels.Argument) channels.Return {

	var err error
	var me *Daemon

	for range only.Once {
		me, err = InterfaceToTypeDaemon(i)
		if err != nil {
			break
		}

		err = me.LoadFiles()

		eblog.Debug(me.EntityId, "loadConfigHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return &err
}


// Non-exposed channel function that responds to a "status" channel request.
// Produces the status of the M-DNS handler via a channel.
func getHandler(event *messages.Message, i channels.Argument) channels.Return {

	var err error
	var me *Daemon
	var ret messages.SubTopics

	for range only.Once {
		me, err = InterfaceToTypeDaemon(i)
		if err != nil {
			break
		}

		switch event.Text.String() {
			case "topics":
				ret = me.channelHandler.GetTopics()
			case "1":
				ret = me.channelHandler.GetTopics()
			case "2":
				ret = me.channelHandler.GetTopics()
		}

		fmt.Printf("topics: %v\n", ret)

		eblog.Debug(me.EntityId, "topicsHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return &ret
}

