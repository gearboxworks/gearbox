package network

import (
	"encoding/json"
	"errors"
	"gearbox/app/logger"
	"gearbox/heartbeat/gbevents/channels"
	"gearbox/heartbeat/gbevents/messages"
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
			s.Domain = defaultDomain
		}

		sc.instance, err = zeroconf.Register(
			s.Name.String(),
			s.Type.String(),
			s.Domain.String(),
			s.Port,
			s.Text,
			nil)
		if err != nil {
			err = errors.New("unable to register zeroconf service")
			break
		}

		u, err = uuid.NewUUID()
		if err != nil {
			sc.instance.Shutdown()
			return &sc, err
		}

		sc.EntityId = u
		sc.Entry.Instance = s.Name.String()
		sc.Entry.Service = s.Type.String()
		sc.Entry.Domain = s.Domain.String()
		sc.Entry.Port = s.Port
		sc.Entry.Text = s.Text
		sc.IsManaged = true

		me.services[u] = &sc

		logger.Debug("ZeroConf %s registered service %s OK", me.EntityId.String(), u.String())
	}

	return &sc, err
}

// Register a service via a channel defined by a *CreateEntry structure and
// returns a *Service structure if successful.
func (me *ZeroConf) RegisterByChannel(caller messages.MessageAddress, s CreateEntry) (*Service, error) {

	var err error
	var j []byte
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
			s.Text = defaultText
		}

		if s.Domain == "" {
			s.Domain = defaultDomain
		}

		j, err = json.Marshal(s)
		if err != nil {
			break
		}

		reg := messages.Message{
			Source: caller,
			Topic: messages.Topic{
				Address: me.EntityId,
				SubTopic: "register",
			},
			Text: messages.MessageText(j),
		}

		err = me.Channels.Publish(reg)
		if err != nil {
			break
		}

		rs, err := me.Channels.GetCallbackReturn(reg, 100)
		if err != nil {
			break
		}

		sc, err = InterfaceToTypeService(rs)	// sc = rs.(*Service)
		if err != nil {
			break
		}

		logger.Debug("ZeroConf %s registered service %s via channel", me.EntityId.String(), sc.EntityId.String())
	}

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

		ce := CreateEntry{}
		err = json.Unmarshal(event.Text.ByteArray(), &ce)
		if err != nil {
			break
		}

		sc, err = me.Register(ce)
		if err != nil {
			break
		}
		// logger.Debug("Service: %v", sc)

		err = nil
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return sc
}

