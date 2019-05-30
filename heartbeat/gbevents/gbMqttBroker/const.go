package gbMqttBroker

import (
	"gearbox/heartbeat/gbevents/messages"
	oss "gearbox/os_support"
	"github.com/fhmq/hmq/broker"
	"github.com/gearboxworks/go-status"
	"net/url"
)

type Broker struct {
	EntityId  messages.MessageAddress
	OsSupport oss.OsSupporter
	Sts       status.Status
	Server    *url.URL

	Config   *broker.Config
	instance *broker.Broker
}

type Args Broker
