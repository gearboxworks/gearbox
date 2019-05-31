package gbMqttClient

import (
	"gearbox/heartbeat/gbevents/messages"
	oss "gearbox/os_support"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gearboxworks/go-status"
	"net/url"
)

type Client struct {
	EntityId  messages.MessageAddress
	OsSupport oss.OsSupporter
	Sts       status.Status
	Config    *mqtt.ClientOptions
	Server *url.URL

	client   mqtt.Client
	token    mqtt.Token
}
type Args Client

const DefaultEntityId = "gearbox-mqtt-client"

type msgCallback struct {
	Topic    messages.Topic
	Function mqtt.MessageHandler
}
