package box

import (
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/states"
	"gearbox/box/external/vmbox"
)


func myCallback(i interface{}, state states.Status) error {

	var err error

	me := i.(*Box)

	if me == nil {
		return err
	}

	//fmt.Printf("Service %s moved to state '%s'\n", state.EntityName, state.Current)
	//if (state.Current == states.StateStarted) || (state.Current == states.StateStopped) {
	//	fmt.Printf("%s service %s %s\n", me.BoxInstance.Boxname, state.EntityName, state.Current)
	//}

	name := state.EntityName
	//fmt.Printf("%s => %s\n", name, state.String())
	//me.SetStateMenu(*name, state.Current)
	//me.SetControlMenu(*name, state.Current)


	switch *name {
		case entity.VmEntityName:
			var s int
			err = me.VmBox.EnsureNotNil()
			if err == nil {
				s, _ = me.VmBox.Releases.Selected.IsIsoFilePresent()
				if s != vmbox.IsoFileDownloaded {
					state.Current = states.StateUpdating
				}
			}
			me.SetStateMenu(entity.VmEntityName, state.Current)
			me.SetControlMenu(entity.VmEntityName, state.Current)

		case entity.ApiEntityName:
			me.SetStateMenu(entity.ApiEntityName, state.Current)

		case entity.UnfsdEntityName:
			me.SetStateMenu(entity.UnfsdEntityName, state.Current)
			err = me.NfsExports.ReadExport()


		case menuVmUpdate:
			if state.Current.String() != "100%" {
				me.menu[menuVmUpdate].MenuItem.SetTitle(me.menu[menuVmUpdate].PrefixMenu + " - " + state.Current.String())
				me.menu[menuVmUpdate].MenuItem.SetTooltip(me.menu[menuVmUpdate].PrefixToolTip + " - " + state.Current.String())
			} else {
				me.menu[menuVmUpdate].MenuItem.SetTitle(me.menu[menuVmUpdate].PrefixMenu)
				me.menu[menuVmUpdate].MenuItem.SetTooltip(me.menu[menuVmUpdate].PrefixToolTip)
			}
	}

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


//func (me *Box) SetMenuVM(state states.State) {
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
//func (me *Box) SetMenuAPI(state states.State) {
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

