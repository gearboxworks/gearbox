package gbMqttClient

import (
	"gearbox/heartbeat/gbevents/messages"
	oss "gearbox/os_support"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gearboxworks/go-status"
	"net/url"
)

type Client struct {
	EntityId  string
	OsSupport oss.OsSupporter
	Sts       status.Status
	Config    *mqtt.ClientOptions
	Server *url.URL

	client   mqtt.Client
	token    mqtt.Token
}
type Args Client

type msgCallback struct {
	Topic    messages.Topic
	Function mqtt.MessageHandler
}
