package gbChannels

import (
	oss "gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"github.com/olebedev/emitter"
)

type Channels struct {
	EntityId  string
	OsSupport  oss.OsSupporter
	Sts        status.Status

	emitter   emitter.Emitter
	events    emitter.Event
	emits     chan struct{}
	group     emitter.Group
}
type Args Channels
