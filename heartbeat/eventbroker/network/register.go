package network

import (
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/only"
	"github.com/grandcat/zeroconf"
)


////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

// Register a service by method defined by a *CreateEntry structure.
func (me *ZeroConf) Register(s CreateEntry) (*Service, error) {

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
		for range only.Once {
			sc.State.SetNewAction(states.ActionRegister)
			sc.EntityId = messages.GenerateAddress()
			sc.IsManaged = true
			sc.channels = me.Channels
			channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)

			if s.Port == 0 {
				s.Port, err = GetFreePort()
				if err != nil {
					break
				}
			}

			if len(s.Text) == 0 {
				s.Text = []string{"txtv=0", "lo=1", "la=2"}
			}

			if s.Domain == "" {
				s.Domain = DefaultDomain
			}

			sc.instance, err = zeroconf.Register(
				s.Name.String(),
				s.Type.String(),
				s.Domain.String(),
				int(s.Port),
				s.Text,
				nil)
			if err != nil {
				err = me.EntityId.ProduceError("unable to register service")
				break
			}

			sc.Entry.Instance = s.Name.String()
			sc.Entry.Service = s.Type.String()
			sc.Entry.Domain = s.Domain.String()
			sc.Entry.Port = int(s.Port)
			sc.Entry.Text = s.Text

			me.services[sc.EntityId] = &sc

			eblog.Debug(me.EntityId, "registered service %s OK", sc.EntityId.String())
		}

		sc.State.SetNewState(states.StateRegistered, err)
		sc.channels.PublishCallerState(&sc.EntityId, &sc.State)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return &sc, err
}

// Register a service via a channel defined by a *CreateEntry structure and
// returns a *Service structure if successful.
func (me *ZeroConf) RegisterByChannel(caller messages.MessageAddress, s CreateEntry) (*Service, error) {

	var err error
	var sc *Service

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if s.Port == 0 {
			s.Port, err = GetFreePort()
			if err != nil {
				break
			}
		}

		if len(s.Text) == 0 {
			s.Text = DefaultText
		}

		if s.Domain == "" {
			s.Domain = DefaultDomain
		}

		reg := ConstructMdnsRegisterMessage(caller, me.EntityId, s)
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

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return sc, err
}


////////////////////////////////////////////////////////////////////////////////
// Executed from a channel.

// Non-exposed channel function that responds to a "register" channel request.
func registerService(event *messages.Message, i channels.Argument) channels.Return {

	var me *ZeroConf
	var sc *Service
	var err error

	for range only.Once {
		me, err = InterfaceToTypeZeroConf(i)
		if err != nil {
			break
		}

		var ce CreateEntry
		ce, err = DeconstructMdnsRegisterMessage(event)
		//err = json.Unmarshal(event.Text.ByteArray(), &ce)
		if err != nil {
			break
		}

		sc, err = me.Register(ce)
		if err != nil {
			break
		}

		eblog.Debug(me.EntityId, "registered service by channel %s OK", sc.EntityId.String())
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return sc
}

