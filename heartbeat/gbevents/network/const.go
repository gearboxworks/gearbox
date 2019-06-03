package network

import (
	"gearbox/heartbeat/gbevents/channels"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/heartbeat/gbevents/tasks"
	"gearbox/only"
	oss "gearbox/os_support"
	"github.com/google/uuid"
	"github.com/grandcat/zeroconf"
	"time"
)


const (
	defaultEntityId = "gearbox-zeroconf"
	defaultWaitTime = time.Millisecond * 2000
	defaultDomain   = "local"
	defaultRetries  = 12
	DefaultRetryDelay = time.Second * 5
)

var defaultText		= []string{"txtv=0", "lo=1", "la=2"}
var browseList		= []string{"_mqtt._udp", "_mqtt._tcp", "_nfs._udp", "_nfs._tcp"}

type ZeroConf struct {
	EntityId        messages.MessageAddress
	osSupport       oss.OsSupporter
	Error           error
	Task            *tasks.Task
	Channels        *channels.Channels
	ChannelHandler	*channels.Subscriber
	restartAttempts int
	waitTime        time.Duration
	domain          string
	services        ServicesMap
	scannedServices ServicesMap
}
type Args ZeroConf

type Service struct {
	EntityId  uuid.UUID
	IsManaged bool
	Entry     Entry
	instance  *zeroconf.Server
}
type Entry zeroconf.ServiceEntry
//type ServicesArray []Entry
type ServicesMap map[uuid.UUID]*Service

type CreateEntry struct {
	Name   Name		`json:"name"`	// == Service.Entry.Instance
	Type   Type		`json:"type"`	// == Service.Entry.Service
	Domain Domain	`json:"domain"`	// == Service.Entry.Domain
	Port   int		`json:"port"`	// == Service.Entry.Port
	Text   Text		`json:"text"`	// == Service.Entry.Text
	TTL    uint32   `json:"ttl"`	// == Service.Entry.TTL
}

/*
func (me *Uuid) String() (string) {

	var f uuid.UUID

	f, _ = uuid.Parse(fmt.Sprintf("%s", me))
	return f.String() // .Sprintf("%s", *me)
}
*/

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
