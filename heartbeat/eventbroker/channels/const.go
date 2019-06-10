package channels

import (
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	oss "gearbox/os_support"
	"github.com/olebedev/emitter"
	"sync"
)


const (
	DefaultEntityId   = "eventbroker-channels"
)


type Channels struct {
	EntityId messages.MessageAddress
	State    states.Status

	mutex            sync.RWMutex	// Mutex control for this struct.
	subscribers      Subscribers
	instance         channelsInstance
	osSupport        oss.OsSupporter
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
	EntityId       messages.MessageAddress
	State          states.Status
	IsManaged      bool

	mutex          sync.RWMutex // Mutex control for this struct.
	topics         References
	parentInstance *channelsInstance
}
type Subscribers map[messages.MessageAddress]*Subscriber

type Reference struct {
	Callback   Callback
	Argument   Argument
	Return     Return
	Executed   bool
	ReturnType ReturnType

	mutex      sync.RWMutex // Mutex control for this struct.
}
type References map[messages.SubTopic]*Reference

type Callback func(event *messages.Message, you Argument) Return

type Argument interface{}

type Return interface{}

type Executed map[messages.SubTopic]bool

type 	ReturnType     string


var IsEmptySubScribers = Subscribers{}
var IsEmptySubScriber = Subscriber{}
//var IsEmptySubTopics = SubTopics{}
//var IsEmptyCallbacks = Callbacks{}

