package network

import (
	"fmt"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/heartbeat/eventbroker/tasks"
	oss "gearbox/os_support"
	"github.com/grandcat/zeroconf"
	"strconv"
	"sync"
	"time"
)


type ZeroConf struct {
	EntityId        messages.MessageAddress
	State           states.Status
	Task            *tasks.Task
	Channels        *channels.Channels

	mutex           sync.RWMutex // Mutex control for map.
	channelHandler	*channels.Subscriber
	restartAttempts int
	waitTime        time.Duration
	domain          string
	services        ServicesMap
	scannedServices ServicesMap
	osSupport       oss.OsSupporter
}
type Args ZeroConf

type Service struct {
	EntityId       messages.MessageAddress
	State          states.Status
	IsManaged      bool
	Entry          Entry

	mutex          sync.RWMutex // Mutex control for map.
	channels       *channels.Channels
	channelHandler *channels.Subscriber
	instance       *zeroconf.Server
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

