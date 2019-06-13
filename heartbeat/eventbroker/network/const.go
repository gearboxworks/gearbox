package network

import (
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/ospaths"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/heartbeat/eventbroker/tasks"
	"github.com/grandcat/zeroconf"
	"net/url"
	"strconv"
	"sync"
	"time"
)


const (
	// DefaultEntityId = "eventbroker-zeroconf"
	DefaultWaitTime = time.Millisecond * 1000
	DefaultDomain   = "local"
	DefaultRetries  = 12
	DefaultRetryDelay = time.Second * 10

	SelectRandomPort = "0"
)


type ZeroConf struct {
	EntityId        messages.MessageAddress
	Boxname         string
	State           *states.Status
	Task            *tasks.Task
	Channels        *channels.Channels

	mutex           sync.RWMutex // Mutex control for map.
	channelHandler	*channels.Subscriber
	restartAttempts int
	waitTime        time.Duration
	domain          string
	services        ServicesMap
	scannedServices ServicesMap
	OsPaths        *ospaths.BasePaths
}
type Args ZeroConf


type Service struct {
	EntityId       messages.MessageAddress
	EntityName     messages.MessageAddress
	EntityParent   *messages.MessageAddress
	State          *states.Status
	IsManaged      bool
	Entry          Entry

	mutex          sync.RWMutex // Mutex control for map.
	channels       *channels.Channels
	channelHandler *channels.Subscriber
	instance       *zeroconf.Server
	osPaths        *ospaths.BasePaths
}
type ServicesMap map[messages.MessageAddress]*Service
type Entry zeroconf.ServiceEntry

var DefaultText		= []string{"txtv=0", "lo=1", "la=2"}
var browseList		= []string{"_mqtt._udp", "_mqtt._tcp", "_nfs._udp", "_nfs._tcp"}


type ServiceConfig struct {
	EntityId  messages.MessageAddress `json:"entity_id"` //
	EntityName     string
	UrlString string                  `json:"urlstring"` //
	Url       *url.URL                `json:"url"`       //

	Name      Name                    `json:"name"`      // == Service.Entry.Instance
	Type      Type                    `json:"type"`      // == Service.Entry.Service
	Domain    Domain                  `json:"domain"`    // == Service.Entry.Domain
	Port      Port                    `json:"port"`      // == Service.Entry.Port
	Text      Text                    `json:"text"`      // == Service.Entry.Text
	TTL       uint32                  `json:"ttl"`       // == Service.Entry.TTL
}

const (
	Package                    = "network"
	InterfaceTypeZeroConf      = "*" + Package + ".ZeroConf"
	InterfaceTypeService       = "*" + Package + ".Service"
	InterfaceTypeServiceConfig = "*" + Package + ".ServiceConfig"
)

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


type Port string
func (me *Port) String() (string) {
	return string(*me)
}
func (me *Port) ToInt() int {

	p, _ := strconv.Atoi(me.String())

	return p
}
//type Port int
//func (me Port) String() (string) {
//
//	return fmt.Sprintf("%d", me)
//}
//func StringToPort(i string) Port {
//
//	p, _ := strconv.Atoi(i)
//
//	return Port(p)
//}


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

