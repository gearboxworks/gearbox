package network

import (
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/only"
	"github.com/google/uuid"
	"github.com/grandcat/zeroconf"
)


////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

// Register a service by method defined by a *CreateEntry structure.
func (me *ZeroConf) Register(s CreateEntry) (*Service, error) {

	var err error
	var u uuid.UUID
	var sc Service

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for range only.Once {
			// @TODO - Need to check to see if this service has already been registered.

			sc.State.SetNewWantState(states.StateRegistered)

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

			sc.EntityId = messages.GenerateAddress()
			sc.Entry.Instance = s.Name.String()
			sc.Entry.Service = s.Type.String()
			sc.Entry.Domain = s.Domain.String()
			sc.Entry.Port = int(s.Port)
			sc.Entry.Text = s.Text
			sc.IsManaged = true
			sc.channels = me.Channels
			sc.State.SetNewState(states.StateRegistered)

			me.services[sc.EntityId] = &sc
			eblog.Debug("ZeroConf %s registered service %s OK", me.EntityId.String(), u.String())
		}

		// Save last state.
		sc.State.Error = err
		channels.PublishCallerState(me.services[sc.EntityId].channels, &me.services[sc.EntityId].EntityId, &me.services[sc.EntityId].State)
	}
	eblog.LogIfError(&me, err)

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

		eblog.Debug("ZeroConf %s registered service %s via channel", me.EntityId.String(), sc.EntityId.String())
	}
	eblog.LogIfError(&me, err)

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

		//fmt.Printf("Rx: %v\n", event)

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

		eblog.Debug("ZeroConf %s registered service %s via channel", me.EntityId.String(), sc.EntityId.String())
	}
	eblog.LogIfError(&me, err)

	return sc
}

