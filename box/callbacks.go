package box

import (
	"fmt"
	"gearbox/box/external/vmbox"
	"gearbox/eventbroker/ebutil"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/states"
)

func myCallback(i interface{}, state states.Status) error {

	var err error

	me := i.(*Box)

	if me == nil {
		return err
	}

	switch state.EntityName {
	case entity.VmEntityName:
		var s int
		err = me.VmBox.EnsureNotNil()
		if err != nil {
			ebutil.LogError(fmt.Errorf("myCallBack[state=%s] error: %s",
				entity.VmEntityName,
				err.Error(),
			))
			// @TODO Should this continue on to set state, or set
			//       a differnt state and break from here?

		} else {
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
		m := me.menu[menuVmUpdate]
		mi := m.MenuItem
		if state.Current != "100%" {
			// @TODO Can we set '100%' as a constant named to indicate why it is formatted as it is?
			mi.SetTitle(fmt.Sprintf("%s-%s", m.PrefixMenu, state.Current))
			mi.SetTooltip(fmt.Sprintf("%s-%s", m.PrefixToolTip, state.Current))
		} else {
			mi.SetTitle(m.PrefixMenu)
			mi.SetTooltip(m.PrefixToolTip)
		}
		m.MenuItem = mi
		me.menu[menuVmUpdate] = m
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
