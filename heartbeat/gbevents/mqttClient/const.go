package mqttClient

import (
	"gearbox/heartbeat/gbevents/channels"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/heartbeat/gbevents/tasks"
	"gearbox/only"
	oss "gearbox/os_support"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"net/url"
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


var msgTemplate = messages.Message{
	Source: DefaultEntityId,
	Topic: messages.MessageTopic{
		Address: "",
		SubTopic: "",
	},
	Text: "",
}


type MqttClient struct {
	EntityId        messages.MessageAddress
	osSupport       oss.OsSupporter
	Error           error
	Task            *tasks.Task
	Channels        *channels.Channels
	ChannelHandler  *channels.Subscriber
	Server          *url.URL

	instance        clientInstance
	restartAttempts int
	waitTime        time.Duration
	domain          string
	services        ServicesMap
}
type Args MqttClient

type clientInstance struct {
	options *mqtt.ClientOptions
	client  mqtt.Client
	token   mqtt.Token
}

type msgCallback struct {
	Topic    Topic
	Function mqtt.MessageHandler
}

type Service struct {
	EntityId  uuid.UUID
	IsManaged bool
	instance  mqtt.Token
}
type ServicesArray []mqtt.Client
type ServicesMap map[uuid.UUID]*Service

type CreateEntry struct {
	Name   string		`json:"name"`	// == Service.Entry.Instance
	TTL    uint32   `json:"ttl"`	// == Service.Entry.TTL
}

type Topic string
func (me *Topic) String() (string) {

	return string(*me)
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
