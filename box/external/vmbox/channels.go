package vmbox

import (
	"gearbox/eventbroker/channels"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/only"
	"gearbox/eventbroker/states"
)


////////////////////////////////////////////////////////////////////////////////
// Executed from a channel

// Non-exposed channel function that responds to an "stop" channel request.
func stopHandler(event *messages.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *VmBox

	for range only.Once {
		me, err = InterfaceToTypeVmBox(i)
		if err != nil {
			break
		}

		if event.Text.String() == "" {
			break
		}

		if event.Text.String() == entity.SelfEntityName {
			// Stop Daemon by default
			err = me.Stop()
		} else {
			// Stop of specific entity
			sc := me.IsExisting(messages.MessageAddress(event.Text))
			if sc != nil {
				err = sc.Stop()
			}
		}

		eblog.Debug(me.EntityId, "requested service stop via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return &err
}


// Non-exposed channel function that responds to an "start" channel request.
func startHandler(event *messages.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *VmBox

	for range only.Once {
		me, err = InterfaceToTypeVmBox(i)
		if err != nil {
			break
		}

		if event.Text.String() == "" {
			break
		}

		if event.Text.String() == entity.SelfEntityName {
			// Start Daemon by default
			err = me.Start()
		} else {
			// Start of specific entity
			sc := me.IsExisting(messages.MessageAddress(event.Text))
			if sc != nil {
				err = sc.Start()
			}
		}

		eblog.Debug(me.EntityId, "requested service start via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return &err
}


// Non-exposed channel function that responds to a "status" channel request.
func statusHandler(event *messages.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *VmBox
	var ret *states.Status

	for range only.Once {
		me, err = InterfaceToTypeVmBox(i)
		if err != nil {
			break
		}

		if event.Text.String() == "" {
			break
		}

		if event.Text.String() == entity.SelfEntityName {
			// Get status of Daemon by default
			ret = me.State.GetStatus()
		} else {
			// Get status of specific entity
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


// Non-exposed channel function that responds to a "update" channel request.
func updateHandler(event *messages.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *VmBox

	for range only.Once {
		me, err = InterfaceToTypeVmBox(i)
		if err != nil {
			break
		}

		if event.Text.String() == "" {
			break
		}

		//if event.Text.String() == entity.SelfEntityName {
		//	// Get status of Daemon by default
		//	ret = me.State.GetStatus()
		//} else {
		//	// Get status of specific entity
		//	sc := me.IsExisting(messages.MessageAddress(event.Text))
		//	if sc != nil {
		//		ret, err = sc.GetStatus()
		//	}
		//}

		if !me.Releases.Selected.IsDownloading {
			err = me.Releases.Selected.GetIso()
		}

		eblog.Debug(me.EntityId, "updateHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return &err
}


// Non-exposed channel function that responds to a "update" channel request.
func createHandler(event *messages.Message, i channels.Argument, r channels.ReturnType) channels.Return {

	var err error
	var me *VmBox

	for range only.Once {
		me, err = InterfaceToTypeVmBox(i)
		if err != nil {
			break
		}

		if event.Text.String() == "" {
			break
		}

		if event.Text.String() == entity.SelfEntityName {
			// Get status of Daemon by default
		} else {
			// Get status of specific entity
			sc := me.IsExisting(messages.MessageAddress(event.Text))
			if sc == nil {
				_, err = me.New(ServiceConfig{
					Name: messages.MessageAddress(event.Text),
					Version: "latest",
				})
			}
		}

		//if !me.Releases.Selected.IsDownloading {
		//	err = me.Releases.Selected.GetIso()
		//}

		eblog.Debug(me.EntityId, "createHandler() via channel")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return &err
}

