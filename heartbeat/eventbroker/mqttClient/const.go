package mqttClient

import (
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/heartbeat/eventbroker/tasks"
	"gearbox/only"
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
	Server          *url.URL
	Channels        *channels.Channels
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
	EntityId  messages.MessageAddress
	State     states.Status
	IsManaged bool
	Entry     *CreateEntry

	mutex     sync.RWMutex	// Mutex control for this struct.
	channels  *channels.Channels
	instance  mqtt.Token
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


// Ensure we don't duplicate services.
func (me *Service) IsExisting(him CreateEntry) error {

	var err error

	switch {
		case me.Entry.Topic == him.Topic:
			err = me.EntityId.ProduceError("MqttClient service Topic:%s already exists", me.Entry.Topic)

		case me.Entry.Name == him.Name:
			err = me.EntityId.ProduceError("MqttClient service Name:%s already exists", me.Entry.Name)
	}

	return err
}

// Ensure we don't duplicate services.
func (me *ServicesMap) IsExisting(him CreateEntry) error {

	var err error

	for _, ce := range *me {
		err = ce.IsExisting(him)
		if err != nil {
			break
		}
	}

	return err
}


func InterfaceToTypeMqttClient(i interface{}) (*MqttClient, error) {

	var err error
	var zc *MqttClient

	for range only.Once {
		err = channels.EnsureArgumentNotNil(i)
		if err != nil {
			break
		}
		zc = i.(*MqttClient)
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
