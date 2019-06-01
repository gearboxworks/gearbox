package network

import (
	"gearbox/heartbeat/gbevents/channels"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/heartbeat/gbevents/tasks"
	oss "gearbox/os_support"
	"github.com/grandcat/zeroconf"
	"time"
)

type Client struct {
	EntityId        messages.MessageAddress
	OsSupport       oss.OsSupporter
	Error           error
	Task            *tasks.Task
	Channels        *channels.Channels
	ChannelHandler	channels.SubTopics
	RestartAttempts int
	WaitTime        time.Duration
	Domain          string
	Services
}
type Args Client

type ServiceEntry zeroconf.ServiceEntry
type ServiceEntries []ServiceEntry


const (
	defaultEntityId = "gearbox-zeroconf"
	defaultWaitTime = time.Second * 1000
	defaultDomain   = "local."
	DefaultRetries  = 12
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
