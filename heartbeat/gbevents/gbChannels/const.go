package gbChannels

import (
	"fmt"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/help"
	oss "gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"github.com/olebedev/emitter"
)

type Channels struct {
	EntityId    messages.MessageAddress
	OsSupport   oss.OsSupporter
	Sts         status.Status
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
	instance *channelsInstance
}
type Callback func(event *messages.Message) status.Status
type Callbacks map[messages.SubTopic]Callback


func (me *SubTopics) EnsureNotNil() status.Status {

	var sts status.Status

	switch {
		case me == nil:
			sts = status.Warn("").
				SetMessage("subtopic is nil").
				SetAdditional("", ).
				SetData("").
				SetHelp(status.AllHelp, help.ContactSupportHelp())

		case me.Address.IsNil() == true:
			sts = status.Warn("").
				SetMessage("subtopic address is nil").
				SetAdditional("", ).
				SetData("").
				SetHelp(status.AllHelp, help.ContactSupportHelp())

		case me.Callbacks == nil:
			sts = status.Warn("").
				SetMessage("subtopic callbacks is nil").
				SetAdditional("", ).
				SetData("").
				SetHelp(status.AllHelp, help.ContactSupportHelp())

		default:
			sts = status.Success("subtopic not nil")
	}

	return sts
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
