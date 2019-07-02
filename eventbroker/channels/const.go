package channels

import (
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/osdirs"
	"gearbox/eventbroker/states"
	"github.com/olebedev/emitter"
	"sync"
)

const (
// DefaultEntityId   = "eventbroker-channels"
)

type Channels struct {
	EntityId msgs.Address
	Boxname  string
	State    *states.Status

	mutex       sync.RWMutex // Mutex control for this struct.
	subscribers Subscribers
	instance    channelsInstance
	OsDirs      *osdirs.BaseDirs
}
type Args Channels

type channelsInstance struct {
	emitter *emitter.Emitter
	events  emitter.Event
	emits   chan struct{}
	//group   *emitter.Group
}

type Event emitter.Event

type Subscriber struct {
	EntityId     msgs.Address
	EntityName   msgs.Address
	EntityParent *msgs.Address
	State        *states.Status
	IsManaged    bool

	mutex          sync.RWMutex // Mutex control for this struct.
	topics         References
	parentInstance *channelsInstance
}
type Subscribers map[msgs.Address]*Subscriber

type Reference struct {
	Callback   Callback
	Argument   Argument
	Return     Return
	Executed   bool
	ReturnType ReturnType

	mutex sync.RWMutex // Mutex control for this struct.
}
type References map[msgs.SubTopic]*Reference

type Callback func(event *msgs.Message, you Argument, ret ReturnType) Return

type Argument interface{}

type Return interface{}

type Executed map[msgs.SubTopic]bool

type ReturnType string

var IsEmptySubScribers = Subscribers{}
var IsEmptySubScriber = Subscriber{}

//var IsEmptySubTopics = SubTopics{}
//var IsEmptyCallbacks = Callbacks{}
