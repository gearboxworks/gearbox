package channels

import (
	"errors"
	"fmt"
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

type Subscribers map[messages.MessageAddress]*Subscriber

type Subscriber struct {
	EntityId       messages.MessageAddress
	State          states.Status
	mutex          sync.RWMutex	// Mutex control for this struct.

	//Callbacks      Callbacks
	//Arguments      Arguments
	//Returns        Returns
	//Executed       Executed
	//mutexExecuted  sync.RWMutex	// Mutex control for map.
	//mutexArguments sync.RWMutex	// Mutex control for map.
	//mutexReturns   sync.RWMutex	// Mutex control for map.

	topics       References

	parentInstance *channelsInstance
}
type Callback func(event *messages.Message, you Argument) Return
type Callbacks map[messages.SubTopic]Callback
type Argument interface{}
type Arguments map[messages.SubTopic]Argument
type Return interface{}
type Returns map[messages.SubTopic]Return
type Executed map[messages.SubTopic]bool

type Reference struct {
	Callback Callback
	Argument Argument
	Return Return
	Executed bool

	mutex          sync.RWMutex	// Mutex control for this struct.
}
type References map[messages.SubTopic]*Reference

func EnsureArgumentNotNil(me Argument) error {

	var err error

	switch {
		case me == nil:
			err = errors.New("channel argument is nil")

		default:
			// err = errors.New("subscriber not nil")
	}

	return err
}


func (me *Channels) EnsureSubscriberNotNil(u messages.MessageAddress) error {

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	if _, ok := me.subscribers[u]; !ok {	// Managed by Mutex
		return me.EntityId.ProduceError("subscriber doesn't exist")
	} else {
		return me.subscribers[u].EnsureNotNil()      // Managed by Mutex
	}
}


func (me *Subscriber) EnsureNotNil() error {

	var err error

	switch {
		case me == nil:
			err = errors.New("subscriber is nil")

		case me.EntityId.EnsureNotNil() != nil:
			err = errors.New("subscriber address is nil")
		//
		//case me.Callbacks == nil:
		//	err = errors.New("subscriber callbacks is nil")
		//
		//case me.Returns == nil:
		//	err = errors.New("subscriber returns is nil")

		default:
			// err = errors.New("subscriber not nil")
	}

	return err
}

func (topics *Subscriber) List() {

	fmt.Printf("# SubTopics created for this entity: %v\n", topics)
}

func (subs *Subscribers) List() {

	fmt.Printf("# Subscribers: %v\n", subs)
}

var IsEmptySubScribers = Subscribers{}
var IsEmptySubScriber = Subscriber{}
//var IsEmptySubTopics = SubTopics{}
var IsEmptyCallbacks = Callbacks{}
