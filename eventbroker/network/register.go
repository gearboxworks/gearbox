package network

import (
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/states"
	"github.com/gearboxworks/go-status/only"
	"github.com/grandcat/zeroconf"
)

////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

// Register a service by method defined by a *CreateEntry structure.
func (me *ZeroConf) Register(s ServiceConfig) (*Service, error) {

	var err error
	var sc Service

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = me.services.IsExisting(s)
		if err != nil {
			break
		}

		// Create new service entry.
		if s.EntityId != "" {
			sc.EntityId = s.EntityId
		} else {
			sc.EntityId = msgs.MakeAddress()
		}
		if s.EntityName != "" {
			sc.EntityName = msgs.Address(s.EntityName)
		} else {
			if s.Name != "" {
				sc.EntityName = msgs.Address(s.Name)
			} else {
				sc.EntityName = sc.EntityId
			}
		}
		sc.EntityParent = &me.EntityId
		sc.State = states.New(sc.EntityId, sc.EntityName, me.EntityId)
		sc.State.SetNewAction(states.ActionStart) // Was states.ActionRegister
		sc.IsManaged = true
		sc.channels = me.Channels
		sc.channels.PublishState(sc.State)

		err = s.Port.IfZeroFindFreePort()
		if err != nil {
			break
		}

		if len(s.Text) == 0 {
			s.Text = []string{"txtv=0", "lo=1", "la=2"}
		}

		if s.Domain == "" {
			s.Domain = DefaultDomain
		}

		// Consider using RegisterProxy() instead.
		sc.instance, err = zeroconf.Register(
			s.Name,
			s.Type,
			s.Domain,
			s.Port.ToInt(),
			s.Text,
			nil)
		if err != nil {
			err = msgs.MakeError(me.EntityId, "unable to register service")
			break
		}

		sc.Entry.Instance = s.Name
		sc.Entry.Service = s.Type
		sc.Entry.Domain = s.Domain
		sc.Entry.Port = s.Port.ToInt()
		sc.Entry.Text = s.Text
		sc.Entry.TTL = s.TTL

		err = me.AddEntity(sc.EntityId, &sc)
		if err != nil {
			break
		}

		sc.State.SetNewState(states.StateStarted, err) // Was states.StateRegistered
		sc.channels.PublishState(sc.State)
		eblog.Debug(me.EntityId, "registered service %s OK", sc.EntityId.String())
	}

	me.Channels.PublishState(me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return &sc, err
}

// Register a service via a channel defined by a *CreateEntry structure and
// returns a *Service structure if successful.
func (me *ZeroConf) RegisterByChannel(caller msgs.Address, s ServiceConfig) (*Service, error) {

	var err error
	var sc *Service

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = s.Port.IfZeroFindFreePort()
		if err != nil {
			break
		}

		if len(s.Text) == 0 {
			s.Text = DefaultText
		}

		if s.Domain == "" {
			s.Domain = DefaultDomain
		}

		reg := ConstructMdnsMessage(caller, me.EntityId, s, states.ActionRegister)
		err = me.Channels.Publish(reg)
		if err != nil {
			break
		}

		rs, err := me.Channels.GetCallbackReturn(reg, 100)
		if err != nil {
			break
		}

		sc, err = InterfaceToTypeService(rs) // sc = rs.(*Service)
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "registered service by channel %s OK", sc.EntityId.String())
	}

	me.Channels.PublishState(me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return sc, err
}
