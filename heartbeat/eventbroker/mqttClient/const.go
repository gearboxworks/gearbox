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
	defaultWaitTime = time.Millisecond * 1000
	defaultDomain   = "local"
	DefaultRetries  = 12
	DefaultRetryDelay = time.Second * 10
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

type Service struct {
	EntityId       messages.MessageAddress
	State          states.Status
	IsManaged      bool
	Entry          *ServiceConfig

	mutex          sync.RWMutex // Mutex control for this struct.
	channels       *channels.Channels
	channelHandler *channels.Subscriber
	instance       mqtt.Token
}
type ServicesMap map[messages.MessageAddress]*Service
//type ServicesArray []mqtt.Client

type ServiceConfig struct {
	Name   string	`json:"name"`	// == Service.Entry.Instance
	Topic  Topic	`json:"topic"`
	TTL    uint32   `json:"ttl"`	// == Service.Entry.TTL
	Qos    byte		`json:"qos"`

	callback mqtt.MessageHandler
}

const (
	Package                    = "mqttClient"
	InterfaceTypeMqttClient    = "*" + Package + ".MqttClient"
	InterfaceTypeService       = "*" + Package + ".Service"
	InterfaceTypeServiceConfig = "*" + Package + ".ServiceConfig"
	InterfaceTypeError         = "error"
)

type Topic string
func (me *Topic) String() (string) {

	return string(*me)
}

//type msgCallback struct {
//	Topic    Topic
//	Function mqtt.MessageHandler
//}

