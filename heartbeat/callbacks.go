package heartbeat

import (
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/states"
)


func myCallback(i interface{}, state states.Status) error {

	var err error

	me := i.(*Heartbeat)

	if me == nil {
		return err
	}

	//fmt.Printf("Service %s moved to state '%s'\n", state.EntityName, state.Current)
	//if (state.Current == states.StateStarted) || (state.Current == states.StateStopped) {
	//	fmt.Printf("%s service %s %s\n", me.BoxInstance.Boxname, state.EntityName, state.Current)
	//}

	name := state.EntityName
	//fmt.Printf("%s\n", name)
	me.SetStateMenu(*name, state.Current)
	//me.SetControl(name, state.Current)


	//fmt.Printf("Service %s moved to state '%s'\n", state.EntityName, state.Current)

	//fmt.Printf("HELLO state: %s\n", PrintServiceState(&state))
	//fmt.Printf("EntityId:%s  Name:%s  Parent:%s  Action:%s  Want:%s  Current:%s  Last:%s  LastWhen:%v  Error:%v\n",
	//	state.EntityId.String(),
	//	state.EntityName.String(),
	//	state.ParentId.String(),
	//	state.Action.String(),
	//	state.Want.String(),
	//	state.Current.String(),
	//	state.Last.String(),
	//	state.LastWhen.Unix(),
	//	state.Error,
	//)


	return err
}


func (me *Heartbeat) SetStateMenu(m messages.MessageAddress, state states.State) {
	// This can clearly be refactored a LOT.

	if _, ok := me.menu[m]; !ok {
		return
	}

	if me.menu[m].MenuItem == nil {
		return
	}

	//if me.State.Unfsd.LastSts != nil {
	//	me.menu.unfsdStatusEntry.SetTooltip(me.State.Unfsd.LastSts.Message())
	//}

	mi := me.menu[m]
	switch state {
		case states.StateUnknown:
			mi.MenuItem.SetIcon(me.getIcon(IconError))

		case states.StateStopping:
			mi.MenuItem.SetIcon(me.getIcon(IconStopping))

		case states.StateStarting:
			mi.MenuItem.SetIcon(me.getIcon(IconStarting))

		case states.StateStarted:
			mi.MenuItem.SetIcon(me.getIcon(IconUp))

		case states.StateStopped:
			mi.MenuItem.SetIcon(me.getIcon(IconDown))

		default:
			mi.MenuItem.SetIcon(me.getIcon(IconWarning))
	}
	mi.MenuItem.SetTitle(mi.PrefixMenu + state.String())
	mi.MenuItem.SetTooltip(mi.PrefixToolTip + state.String())

	return
}


func (me *Heartbeat) SetControlMenu(m messages.MessageAddress, state states.State) {
	// This can clearly be refactored a LOT.

	if _, ok := me.menu[m]; !ok {
		return
	}

	if me.menu[m].MenuItem == nil {
		return
	}

	//if me.State.Unfsd.LastSts != nil {
	//	me.menu.unfsdStatusEntry.SetTooltip(me.State.Unfsd.LastSts.Message())
	//}

	mi := me.menu[m]
	switch state {
		case states.StateIdle:
			mi.MenuItem.Disabled()

		case states.StateUnknown:
			mi.MenuItem.Enable()

		case states.StateStopping:
			mi.MenuItem.Enable()

		case states.StateStarting:
			mi.MenuItem.Enable()

		case states.StateStarted:
			mi.MenuItem.Enable()

		case states.StateStopped:
			mi.MenuItem.Enable()

		default:
			mi.MenuItem.Enable()
	}
	mi.MenuItem.SetTitle(mi.PrefixMenu + state.String())
	mi.MenuItem.SetTooltip(mi.PrefixToolTip + state.String())

	return
}


//func (me *Heartbeat) SetMenuVM(state states.State) {
//	// This can clearly be refactored a LOT.
//
//	if me.menu.vmStatusEntry == nil {
//		return
//	}
//
//	//if me.State.Unfsd.LastSts != nil {
//	//	me.menu.unfsdStatusEntry.SetTooltip(me.State.Unfsd.LastSts.Message())
//	//}
//	switch state {
//		case states.StateUnknown:
//			me.menu.vmStatusEntry.SetIcon(me.getIcon(IconError))
//
//		case states.StateStopping:
//			me.menu.vmStatusEntry.SetIcon(me.getIcon(IconStopping))
//
//		case states.StateStarting:
//			me.menu.vmStatusEntry.SetIcon(me.getIcon(IconStarting))
//
//		case states.StateStarted:
//			me.menu.vmStatusEntry.SetIcon(me.getIcon(IconUp))
//
//		case states.StateStopped:
//			me.menu.vmStatusEntry.SetIcon(me.getIcon(IconDown))
//
//		default:
//			me.menu.vmStatusEntry.SetIcon(me.getIcon(IconWarning))
//	}
//	me.menu.vmStatusEntry.SetTitle("Gearbox OS: " + state.String())
//
//	return
//}
//
//
//func (me *Heartbeat) SetMenuAPI(state states.State) {
//	// This can clearly be refactored a LOT.
//
//	if me.menu.apiStatusEntry == nil {
//		return
//	}
//
//	//if me.State.Unfsd.LastSts != nil {
//	//	me.menu.unfsdStatusEntry.SetTooltip(me.State.Unfsd.LastSts.Message())
//	//}
//	switch state {
//		case states.StateUnknown:
//			me.menu.apiStatusEntry.SetIcon(me.getIcon(IconError))
//
//		case states.StateStopping:
//			me.menu.apiStatusEntry.SetIcon(me.getIcon(IconStopping))
//
//		case states.StateStarting:
//			me.menu.apiStatusEntry.SetIcon(me.getIcon(IconStarting))
//
//		case states.StateStarted:
//			me.menu.apiStatusEntry.SetIcon(me.getIcon(IconUp))
//
//		case states.StateStopped:
//			me.menu.apiStatusEntry.SetIcon(me.getIcon(IconDown))
//
//		default:
//			me.menu.apiStatusEntry.SetIcon(me.getIcon(IconWarning))
//	}
//	me.menu.apiStatusEntry.SetTitle("Gearbox API: " + state.String())
//
//	return
//}

