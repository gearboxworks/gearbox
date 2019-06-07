package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/heartbeat/eventbroker/tasks"
	"gearbox/only"
	oss "gearbox/os_support"
	"github.com/grandcat/zeroconf"
	"strconv"
	"time"
)


type ZeroConf struct {
	EntityId        messages.MessageAddress
	State           states.Status
	Task            *tasks.Task
	Channels        *channels.Channels
	ChannelHandler	*channels.Subscriber

	restartAttempts int
	waitTime        time.Duration
	domain          string
	services        ServicesMap
	scannedServices ServicesMap
	osSupport       oss.OsSupporter
}
type Args ZeroConf

type Service struct {
	EntityId  messages.MessageAddress
	State     states.Status
	IsManaged bool
	Entry     Entry

	channels  *channels.Channels
	instance  *zeroconf.Server
}
type Entry zeroconf.ServiceEntry
type ServicesMap map[messages.MessageAddress]*Service


const (
	DefaultEntityId = "eventbroker-zeroconf"
	DefaultWaitTime = time.Millisecond * 2000
	DefaultDomain   = "local"
	DefaultRetries  = 12
	DefaultRetryDelay = time.Second * 5
)

var DefaultText		= []string{"txtv=0", "lo=1", "la=2"}
var browseList		= []string{"_mqtt._udp", "_mqtt._tcp", "_nfs._udp", "_nfs._tcp"}


type CreateEntry struct {
	EntityId messages.MessageAddress `json:"entity_id"` //
	Url      string         `json:"url"`       //

	Name     Name           `json:"name"`      // == Service.Entry.Instance
	Type     Type           `json:"type"`      // == Service.Entry.Service
	Domain   Domain         `json:"domain"`    // == Service.Entry.Domain
	Port     Port           `json:"port"`      // == Service.Entry.Port
	Text     Text           `json:"text"`      // == Service.Entry.Text
	TTL      uint32         `json:"ttl"`       // == Service.Entry.TTL
}


var msgTemplate = messages.Message{
	Source: "",
	Topic: messages.MessageTopic{
		Address: "",
		SubTopic: "",
	},
	Text: "",
}


/*
	Servers                 []*url.URL
	ClientID                string
	Username                string
	Password                string

*/


type Port int
func (me Port) String() (string) {

	return fmt.Sprintf("%d", me)
}
func StringToPort(i string) Port {

	p, _ := strconv.Atoi(i)

	return Port(p)
}

type Host string
func (me *Host) String() (string) {

	return string(*me)
}

type Name string
func (me *Name) String() (string) {

	return string(*me)
}

type Type string
func (me *Type) String() (string) {

	return string(*me)
}

type Domain string
func (me *Domain) String() (string) {

	return string(*me)
}

type Text []string
func (me *Text) String() ([]string) {

	return []string(*me)
}


func ConstructMdnsRegisterMessage(me messages.MessageAddress, to messages.MessageAddress, s CreateEntry) messages.Message {

	var err error
	var msgTemplate messages.Message
	var j []byte

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		j, err = json.Marshal(s)
		if err != nil {
			break
		}

		msgTemplate = messages.Message{
			Source: me,
			Topic: messages.MessageTopic{
				Address:  to,
				SubTopic: states.ActionRegister,
			},
			Text: messages.MessageText(j),
		}
	}

	return msgTemplate
}


func DeconstructMdnsRegisterMessage(event *messages.Message) (CreateEntry, error) {

	var err error
	var ce CreateEntry

	for range only.Once {
		//err = ce.EnsureNotNil()
		if event == nil {
			err = errors.New("message is nil")
			break
		}

		err = json.Unmarshal(event.Text.ByteArray(), &ce)
		if err != nil {
			break
		}
	}

	return ce, err
}


func ConstructMdnsUnregisterMessage(me messages.MessageAddress, to messages.MessageAddress, s CreateEntry) messages.Message {

	var err error
	var msgTemplate messages.Message
	var j []byte

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		j, err = json.Marshal(s)
		if err != nil {
			break
		}

		msgTemplate = messages.Message{
			Source: me,
			Topic: messages.MessageTopic{
				Address:  to,
				SubTopic: states.ActionUnregister,
			},
			Text: messages.MessageText(j),
		}
	}

	return msgTemplate
}


func InterfaceToTypeZeroConf(i interface{}) (*ZeroConf, error) {

	var err error
	var zc *ZeroConf

	for range only.Once {
		err = channels.EnsureArgumentNotNil(i)
		if err != nil {
			break
		}
		zc = i.(*ZeroConf)
		// zc = (i[0]).(*ZeroConf)
		// zc = i[0].(*ZeroConf)

		err = zc.EnsureNotNil()
		if err != nil {
			break
		}
	}

	return zc, err
}


func InterfaceToTypeService(i interface{}) (*Service, error) {

	var err error
	var s *Service

	for range only.Once {
		err = channels.EnsureArgumentNotNil(i)
		if err != nil {
			break
		}
		s = i.(*Service)
		// zc = (i[0]).(*Service)
		// zc = i[0].(*Service)

		err = s.EnsureNotNil()
		if err != nil {
			break
		}
	}

	return s, err
}
