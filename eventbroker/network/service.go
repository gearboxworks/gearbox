package network

import (
	"gearbox/eventbroker/channels"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/osdirs"
	"gearbox/eventbroker/states"
	"github.com/grandcat/zeroconf"
	"sync"
)

type ServicesMap map[msgs.Address]*Service
type Service struct {
	EntityId     msgs.Address
	EntityName   msgs.Address
	EntityParent *msgs.Address
	State        *states.Status
	IsManaged    bool
	Entry        Entry

	mutex          sync.RWMutex // Mutex control for map.
	channels       *channels.Channels
	channelHandler *channels.Subscriber
	instance       *zeroconf.Server
	osPaths        *osdirs.BaseDirs
}
