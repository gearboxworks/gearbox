package gbMqttBroker

import (
	"gearbox/heartbeat/daemon/tasks"
	"gearbox/heartbeat/gbevents/messages"
	oss "gearbox/os_support"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fhmq/hmq/broker"
	"github.com/gearboxworks/go-status"
	"net/url"
)

type Mqtt struct {
	EntityId      messages.MessageAddress
	OsSupport     oss.OsSupporter
	ServerURL     *url.URL
	BrokerWorkers int
	Broker        Broker
	Client        Client
}
type Args Mqtt

type Broker struct {
	EntityId        messages.MessageAddress
	Task            tasks.Task
	Sts             status.Status
	RestartAttempts int
	restartCounter  int
	config          *broker.Config
	instance        *broker.Broker
}

type Client struct {
	EntityId        messages.MessageAddress
	Task            tasks.Task
	Sts             status.Status
	RestartAttempts int
	restartCounter  int
	config          *mqtt.ClientOptions
	instance        mqtt.Client
	token           mqtt.Token
}

type msgCallback struct {
	Topic    messages.Topic
	Function mqtt.MessageHandler
}

const (
	DefaultEntityId  = "gearbox-mqtt"
	DefaultBrokerEntityId  = DefaultEntityId + "-broker"
	DefaultClientEntityId  = DefaultEntityId + "-client"
	DefaultServerUrl = "tcp://0.0.0.0:1883"
	DefaultWorkers   = 1024
	DefaultRetries   = 12
)
