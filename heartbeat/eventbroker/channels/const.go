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
	EntityId    messages.MessageAddress
	State       states.Status

	subscribers Subscribers
	instance    channelsInstance
	osSupport   oss.OsSupporter
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
	Callbacks      Callbacks
	Arguments      Arguments
	Returns        Returns
	Executed       Executed
	mutexExecuted  sync.RWMutex
	mutexArguments sync.RWMutex
	mutexReturns   sync.RWMutex

	parentInstance *channelsInstance
}
type Callback func(event *messages.Message, you Argument) Return
type Callbacks map[messages.SubTopic]Callback
type Argument interface{}
type Arguments map[messages.SubTopic]Argument
type Return interface{}
type Returns map[messages.SubTopic]Return
type Executed map[messages.SubTopic]bool


// Mutex handling.
func (me *Subscriber) DeleteSubTopic(sub messages.SubTopic) {

	me.mutexExecuted.Lock()
	delete(me.Executed, sub)
	me.mutexExecuted.Unlock()

	me.mutexArguments.Lock()
	delete(me.Arguments, sub)
	me.mutexArguments.Unlock()

	me.mutexReturns.Lock()
	delete(me.Returns, sub)
	me.mutexReturns.Unlock()

	return
}


// Mutex handling.
func (me *Subscriber) GetExecuted(sub messages.SubTopic) bool {
	var r bool

	me.mutexExecuted.RLock()

	r = me.Executed[sub]

	me.mutexExecuted.RUnlock()

	return r
}

func (me *Subscriber) SetExecuted(sub messages.SubTopic, v bool) {
	me.mutexExecuted.Lock()

	me.Executed[sub] = v

	me.mutexExecuted.Unlock()

	return
}

func (me *Subscriber) ValidateExecuted(sub messages.SubTopic) bool {
	var r bool

	me.mutexExecuted.RLock()

	if _, ok := me.Executed[sub]; ok {
		r = true
	} else {
		r = false
	}

	me.mutexExecuted.RUnlock()

	return r
}


// Mutex handling.
func (me *Subscriber) GetArguments(sub messages.SubTopic) Argument {

	me.mutexArguments.RLock()
	v := me.Arguments[sub]
	me.mutexArguments.RUnlock()

	return v
}

func (me *Subscriber) SetArguments(sub messages.SubTopic, v Argument) {

	me.mutexArguments.Lock()
	me.Arguments[sub] = v
	me.mutexArguments.Unlock()

	return
}

func (me *Subscriber) ValidateArguments(sub messages.SubTopic) bool {
	var r bool

	me.mutexArguments.RLock()

	if _, ok := me.Arguments[sub]; ok {
		r = true
	} else {
		r = false
	}

	me.mutexArguments.RUnlock()

	return r
}


// Mutex handling.
func (me *Subscriber) GetReturns(sub messages.SubTopic) Return {

	me.mutexReturns.RLock()
	r := me.Returns[sub]
	me.mutexReturns.RUnlock()

	return r
}

func (me *Subscriber) SetReturns(sub messages.SubTopic, v Return) {

	me.mutexReturns.Lock()
	me.Returns[sub] = v
	me.mutexReturns.Unlock()

	return
}

func (me *Subscriber) ValidateReturns(sub messages.SubTopic) bool {
	var r bool

	me.mutexReturns.RLock()

	if _, ok := me.Returns[sub]; ok {
		r = true
	} else {
		r = false
	}

	me.mutexReturns.RUnlock()

	return r
}


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
