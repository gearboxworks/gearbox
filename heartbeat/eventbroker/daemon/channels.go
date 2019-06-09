package daemon

import (
	"errors"
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

// Non-exposed channel function that responds to a "status" channel request.
// Produces the status of the M-DNS handler via a channel.
func getHandler(event *messages.Message, i channels.Argument) channels.Return {

	var err error
	var me *Daemon
	var response *states.Status

	for range only.Once {
		me, err = InterfaceToTypeDaemon(i)
		if err != nil {
			break
		}

		//foo := reflect.ValueOf(me)
		//fmt.Printf("Reflect %s, %s	", foo.Type(), foo.String())
		//fmt.Printf("getHandler(%s)\n", me.Fluff)
		fmt.Printf("%d getHandler records: %v\n", time.Now().Unix(), me.State.GetError())
		me.Fluff = fmt.Sprintf("%d", time.Now().Unix())

		switch event.Text.String() {
			case "status":
				//response = event.Text.String()
		}

		me.State.SetError(errors.New("YES - " + me.Fluff))
		response = me.State.GetFullState()

		eblog.Debug(me.EntityId, "requested service status via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return response
}


// Non-exposed channel function that responds to an "stop" channel request.
// Causes the M-DNS handler task to stop via a channel.
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

		eblog.Debug(me.EntityId, "requested service start via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}

