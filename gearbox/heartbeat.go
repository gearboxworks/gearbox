package gearbox

import (
	"gearbox/heartbeat"
	"gearbox/only"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
)



func (me *Gearbox) HeartbeatDaemon(args heartbeat.Args) (sts status.Status) {

	for range only.Once {
		var bx *heartbeat.Heartbeat

		bx, sts = heartbeat.New(me.GetOsSupport(), args)
		if is.Error(sts) {
			break
		}

		sts = bx.HeartbeatDaemon()
		if is.Error(sts) {
			break
		}
	}
	//status.Log(sts)

	return sts
}

func (me *Gearbox) StartHeartbeat(args heartbeat.Args) (sts status.Status) {

	for range only.Once {
		var bx *heartbeat.Heartbeat

		bx, sts = heartbeat.New(me.GetOsSupport(), args)
		if is.Error(sts) {
			break
		}

		sts = bx.StartHeartbeat()
		if is.Error(sts) {
			break
		}
	}
	status.Log(sts)

	return sts
}

func (me *Gearbox) StopHeartbeat(args heartbeat.Args) (sts status.Status) {

	for range only.Once {
		var bx *heartbeat.Heartbeat

		bx, sts = heartbeat.New(me.GetOsSupport(), args)
		if is.Error(sts) {
			break
		}

		sts = bx.StopHeartbeat()
		if is.Error(sts) {
			break
		}
	}
	status.Log(sts)

	return sts
}

func (me *Gearbox) RestartHeartbeat(args heartbeat.Args) (sts status.Status) {

	for range only.Once {
		var bx *heartbeat.Heartbeat

		bx, sts = heartbeat.New(me.GetOsSupport(), args)
		if is.Error(sts) {
			break
		}

		sts = bx.RestartHeartbeat()
		if is.Error(sts) {
			break
		}
	}
	status.Log(sts)

	return sts
}

func (me *Gearbox) PrintHeartbeatStatus(args heartbeat.Args) (sts status.Status) {

	for range only.Once {
		var bx *heartbeat.Heartbeat

		bx, sts = heartbeat.New(me.GetOsSupport(), args)
		if is.Error(sts) {
			break
		}

		sts := bx.GetState()
		if is.Error(sts) {
			break
		}

//		var state string
//		state, sts = sts.GetString()
//		meaning := heartbeat.GetStateMeaning(heartbeat.State(state))
//		if meaning == "" {
//			fmt.Println(heartbeat.GetStateMeaning(heartbeat.UnknownState))
//			break
//		}
//		fmt.Println(meaning)
	}
	status.Log(sts)

	return sts
}