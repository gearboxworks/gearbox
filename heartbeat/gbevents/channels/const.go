package channels

import (
	"fmt"
	"errors"
	"gearbox/heartbeat/gbevents/messages"
	oss "gearbox/os_support"
	"github.com/olebedev/emitter"
)

type Channels struct {
	EntityId    messages.MessageAddress
	OsSupport   oss.OsSupporter
	Error       error
	Subscribers Subscribers
	instance    channelsInstance
}
type Args Channels

const DefaultEntityId = "gearbox-channels"

type channelsInstance struct {
	emitter   emitter.Emitter
	events    emitter.Event
	emits     chan struct{}
	group     emitter.Group
}

type Event emitter.Event

type Subscribers Subscriber
type Subscriber map[messages.MessageAddress]SubTopics

type SubTopics struct {
	Address   messages.MessageAddress
	Callbacks Callbacks
	Interfaces Interfaces
	instance *channelsInstance
}
type Callback func(event *messages.Message, you interface{}) error
type Callbacks map[messages.SubTopic]Callback
type Interface interface{}
type Interfaces map[messages.SubTopic]Interface


func (me *SubTopics) EnsureNotNil() error {

	var err error

	switch {
		case me == nil:
			err = errors.New("subtopic is nil")

		case me.Address.IsNil() == true:
			err = errors.New("subtopic address is nil")

		case me.Callbacks == nil:
			err = errors.New("subtopic callbacks is nil")

		default:
			// err = errors.New("subtopic not nil")
	}

	return err
}

func (topics *SubTopics) List() {

	fmt.Printf("# SubTopics: %v\n", topics)
}

func (subs *Subscribers) List () {

	fmt.Printf("# Subscribers: %v\n", subs)
}

var IsEmptySubScribers = Subscribers{}
var IsEmptySubScriber = Subscriber{}
var IsEmptySubTopics = SubTopics{}
var IsEmptyCallbacks = Callbacks{}
