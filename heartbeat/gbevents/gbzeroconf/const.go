package gbZeroConf

import (
	"gearbox/heartbeat/gbevents/messages"
	oss "gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"time"
)

type Client struct {
	EntityId        messages.MessageAddress
	OsSupport       oss.OsSupporter
	Sts             status.Status
	State           bool
	RestartAttempts int
	restartCounter  int

	WaitTime time.Duration
	Domain   string
	Services
}
type Args Client

const (
	defaultEntityId = "gearbox-zeroconf"
	defaultWaitTime = time.Second * 15
	defaultDomain   = "local."
	defaultRetries  = 12
)

type Service struct {
	Name   ServiceName
	Type   ServiceType
	Domain ServiceDomain
	Port   int
}
type Services []Service

type ServiceName string

func (me *ServiceName) String() (string) {

	return string(*me)
}

type ServiceType string

func (me *ServiceType) String() (string) {

	return string(*me)
}

type ServiceDomain string

func (me *ServiceDomain) String() (string) {

	return string(*me)
}
