package mqttClient

import (
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/heartbeat/eventbroker/tasks"
	oss "gearbox/os_support"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"net/url"
	"sync"
	"time"
)


const (
	DefaultEntityId = "eventbroker-mqttclient"
	defaultWaitTime = time.Millisecond * 2000
	defaultDomain   = "local"
	DefaultRetries  = 12
	DefaultRetryDelay = time.Second * 3
	DefaultServer = "tcp://127.0.0.1:1883"
)


type MqttClient struct {
	EntityId        messages.MessageAddress
	State           states.Status
	Task            *tasks.Task
	Channels        *channels.Channels
	Server          *url.URL

	mutex           sync.RWMutex // Mutex control for map.
	channelHandler  *channels.Subscriber
	instance        clientInstance
	restartAttempts int
	waitTime        time.Duration
	domain          string
	services        ServicesMap
	osSupport       oss.OsSupporter
}
type Args MqttClient

type clientInstance struct {
	options *mqtt.ClientOptions
	client  mqtt.Client
	token   mqtt.Token
}

//type msgCallback struct {
//	Topic    Topic
//	Function mqtt.MessageHandler
//}

type Service struct {
	EntityId       messages.MessageAddress
	State          states.Status
	IsManaged      bool
	Entry          *CreateEntry

	mutex          sync.RWMutex // Mutex control for this struct.
	channels       *channels.Channels
	channelHandler *channels.Subscriber
	instance       mqtt.Token
}
type ServicesArray []mqtt.Client
type ServicesMap map[messages.MessageAddress]*Service

type CreateEntry struct {
	Name   string	`json:"name"`	// == Service.Entry.Instance
	Topic  Topic	`json:"topic"`
	TTL    uint32   `json:"ttl"`	// == Service.Entry.TTL
	Qos    byte		`json:"qos"`
	callback mqtt.MessageHandler
}

type Topic string
func (me *Topic) String() (string) {

	return string(*me)
}

