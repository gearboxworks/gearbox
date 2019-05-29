package gbZeroConf

import (
	oss "gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"time"
)

type Client struct {
	EntityId  string
	OsSupport oss.OsSupporter
	Sts       status.Status

	WaitTime	time.Duration
	Domain		string
	Services
}
type Args Client

const (
	defaultEntityId	= "gearbox-zeroconf"
	defaultWaitTime = time.Second * 15
	defaultDomain	= "local."
)

type Service struct {
	Name ServiceName
	Type ServiceType
	Domain ServiceDomain
	Port int
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

