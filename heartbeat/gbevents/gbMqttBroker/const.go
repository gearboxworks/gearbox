package gbMqttBroker

import (
	oss "gearbox/os_support"
	"github.com/fhmq/hmq/broker"
	"github.com/gearboxworks/go-status"
)

type Broker struct {
	OsSupport oss.OsSupporter
	Sts       status.Status

	Config   *broker.Config
	instance *broker.Broker
}

type Args Broker
