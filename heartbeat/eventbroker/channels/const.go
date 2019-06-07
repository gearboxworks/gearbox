package channels

import (
	"errors"
	"fmt"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	oss "gearbox/os_support"
	"github.com/olebedev/emitter"
)


const (
	DefaultEntityId   = "eventbroker-channels"
)


type Channels struct {
	EntityId    messages.MessageAddress
	State       states.Status

	subscribers Subscribers
	instance    channelsInstance
	osSupport   oss.OsSupporter
}
type Args Channels

type channelsInstance struct {
	emitter emitter.Emitter
	events  emitter.Event
	emits   chan struct{}
	group   emitter.Group
}

type Event emitter.Event

type Subscribers map[messages.MessageAddress]*Subscriber

type Subscriber struct {
	EntityId       messages.MessageAddress
	State          states.Status
	Callbacks      Callbacks
	Arguments      Arguments
	Returns        Returns
	Executed       Executed
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
	Callback
	Argument
	Return
	Executed bool
}
type References map[messages.SubTopic]Reference

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

func (me *Subscriber) EnsureNotNil() error {

	var err error

	switch {
		case me == nil:
			err = errors.New("subscriber is nil")

		case me.EntityId.EnsureNotNil() != nil:
			err = errors.New("subscriber address is nil")

		case me.Callbacks == nil:
			err = errors.New("subscriber callbacks is nil")

		case me.Returns == nil:
			err = errors.New("subscriber returns is nil")

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
