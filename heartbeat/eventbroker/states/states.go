package states

import (
	"sync"
	"time"
)

const (
	StateIdle          = "idle"
	StateUnknown       = "unknown"
	StateError         = "error"
	StateInitializing  = "initializing"
	StateInitialized   = "initialized"
	StateRegistering   = "registering"
	StateRegistered    = "registered"
	StateUnregistering = "unregistering"
	StateUnregistered  = "unregistered"
	StatePublishing    = "publishing"
	StatePublished     = "published"
	StateUnpublishing  = "unpublishing"
	StateUnpublished   = "unpublished"
	StateSubscribing   = "subscribing"
	StateSubscribed    = "subscribed"
	StateUnsubscribing = "unsubscribing"
	StateUnsubscribed  = "unsubscribed"
	StateStopped       = "stopped"
	StateStarting      = "starting"
	StateStarted       = "started"
	StateStopping      = "stopping"

	StateIndexIdle          = 0
	StateIndexUnknown       = iota
	StateIndexError         = iota
	StateIndexInitializing  = iota
	StateIndexInitialized   = iota
	StateIndexRegistering   = iota
	StateIndexRegistered    = iota
	StateIndexUnregistering = iota
	StateIndexUnregistered  = iota
	StateIndexPublishing    = iota
	StateIndexPublished     = iota
	StateIndexUnpublishing  = iota
	StateIndexUnpublished   = iota
	StateIndexSubscribing   = iota
	StateIndexSubscribed    = iota
	StateIndexUnsubscribing = iota
	StateIndexUnsubscribed  = iota
	StateIndexStopped       = iota
	StateIndexStarting      = iota
	StateIndexStarted       = iota
	StateIndexStopping      = iota
)
var StateName = map[int]Action{
	StateIndexIdle:          StateIdle,
	StateIndexUnknown:       StateUnknown,
	StateIndexError:         StateError,
	StateIndexInitializing:  StateInitializing,
	StateIndexInitialized:   StateInitialized,
	StateIndexRegistering:   StateRegistering,
	StateIndexRegistered:    StateRegistered,
	StateIndexUnregistering: StateUnregistering,
	StateIndexUnregistered:  StateUnregistered,
	StateIndexPublishing:    StatePublishing,
	StateIndexPublished:     StatePublished,
	StateIndexUnpublishing:  StateUnpublishing,
	StateIndexUnpublished:   StateUnpublished,
	StateIndexSubscribing:   StateSubscribing,
	StateIndexSubscribed:    StateSubscribed,
	StateIndexUnsubscribing: StateUnsubscribing,
	StateIndexUnsubscribed:  StateUnsubscribed,
	StateIndexStopped:       StateStopped,
	StateIndexStarting:      StateStarting,
	StateIndexStarted:       StateStarted,
	StateIndexStopping:      StateStopping,
}


const (
	ActionIdle        = "idle"
	ActionUnknown     = "unknown"
	ActionInitialize  = "init"
	ActionRegister    = "register"
	ActionUnregister  = "unregister"
	ActionPublish     = "publish"
	ActionUnpublish   = "unpublish"
	ActionSubscribe   = "subscribe"
	ActionUnsubscribe = "unsubscribe"
	ActionStop        = "stop"
	ActionStart       = "start"
	ActionStatus      = "status"
	ActionError       = "error"
	ActionGlob        = "*"

	ActionIndexIdle        = 0
	ActionIndexUnknown     = iota
	ActionIndexInitialize  = iota
	ActionIndexRegister    = iota
	ActionIndexUnregister  = iota
	ActionIndexPublish     = iota
	ActionIndexUnpublish   = iota
	ActionIndexSubscribe   = iota
	ActionIndexUnsubscribe = iota
	ActionIndexStop        = iota
	ActionIndexStart       = iota
	ActionIndexStatus      = iota
	ActionIndexError       = iota
	ActionIndexGlob        = iota
)
var ActionName = map[int]Action{
	ActionIndexIdle:        ActionIdle,
	ActionIndexUnknown:     ActionUnknown,
	ActionIndexInitialize:  ActionInitialize,
	ActionIndexRegister:    ActionRegister,
	ActionIndexUnregister:  ActionUnregister,
	ActionIndexPublish:     ActionPublish,
	ActionIndexUnpublish:   ActionUnpublish,
	ActionIndexSubscribe:   ActionSubscribe,
	ActionIndexUnsubscribe: ActionUnsubscribe,
	ActionIndexStop:        ActionStop,
	ActionIndexStart:       ActionStart,
	ActionIndexStatus:      ActionStatus,
	ActionIndexError:       ActionError,
	ActionIndexGlob:        ActionGlob,
}


type Status struct {
	Current  State
	Want     State
	Last     State
	LastWhen time.Time
	Attempts int
	Error    error
	Action   Action

	mutex    sync.RWMutex	// Mutex control for map.
}


type Action string
func (me Action) String() string {
	return string(me)
}
func (me Action) Index() string {
	return string(me)
}


type State string
func (me State) String() string {
	return string(me)
}


// An FSM of sorts.
//func (me *Status) NextAction() Action {
//
//	var intended Action
//
//	for range only.Once {
//		//err = me.EnsureNotNil()
//		//if err != nil {
//		//	break
//		//}
//
//		switch me.GetWant() {
//			case ActionUnregister:
//				switch me.GetCurrent() {
//					case StateStarted:
//						intended = ActionStop
//					case StateStopped:
//						intended = ActionUnregister
//				}
//
//			case ActionRegister:
//				switch me.GetCurrent() {
//					case StateStarted:
//						intended = ActionStop
//					case StateStopped:
//						intended = ActionUnregister
//				}
//
//
//		}
//	}
//
//	return intended
//}
//
//type NextStates []State
//
//var TakeAction = map[State]NextStates{
//	// [Action][State] = dosomething
//	// [ActionUnregister][
//	//ActionIndexIdle:		{Current: StateIdle,			Want: StateError, Action: ActionPublish},
//	//ActionIndexUnknown:		{Current: StateUnknown,			Want: StateError, Action: ActionPublish},
//	//ActionIndexRegister:	{Current: StateRegistered,		Want: StateError, Action: ActionPublish},
//	//ActionIndexUnregister:	{Current: StateUnregistered,	Want: StateError, Action: ActionPublish},
//	//ActionIndexPublish:		{Current: StatePublished,		Want: StateError, Action: ActionPublish},
//	//ActionIndexUnpublish:	{Current: StateUnpublished,		Want: StateError, Action: ActionPublish},
//	//ActionIndexSubscribe:	{Current: StateSubscribed,		Want: StateError, Action: ActionPublish},
//	//ActionIndexUnsubscribe:	{Current: StateUnsubscribed,	Want: StateError, Action: ActionPublish},
//	//ActionIndexStop:		{Current: StateStopped,			Want: StateError, Action: ActionPublish},
//	//ActionIndexStart:		{Current: StateStarted,			Want: StateError, Action: ActionPublish},
//	//ActionIndexError:		{Current: StateError,			Want: StateError, Action: ActionPublish},
//}
//
//
////var TakeAction = map[int]Status{
////	ActionIndexIdle:		{Current: StateIdle,			Want: StateError, Action: ActionPublish},
////	ActionIndexUnknown:		{Current: StateUnknown,			Want: StateError, Action: ActionPublish},
////	ActionIndexRegister:	{Current: StateRegistered,		Want: StateError, Action: ActionPublish},
////	ActionIndexUnregister:	{Current: StateUnregistered,	Want: StateError, Action: ActionPublish},
////	ActionIndexPublish:		{Current: StatePublished,		Want: StateError, Action: ActionPublish},
////	ActionIndexUnpublish:	{Current: StateUnpublished,		Want: StateError, Action: ActionPublish},
////	ActionIndexSubscribe:	{Current: StateSubscribed,		Want: StateError, Action: ActionPublish},
////	ActionIndexUnsubscribe:	{Current: StateUnsubscribed,	Want: StateError, Action: ActionPublish},
////	ActionIndexStop:		{Current: StateStopped,			Want: StateError, Action: ActionPublish},
////	ActionIndexStart:		{Current: StateStarted,			Want: StateError, Action: ActionPublish},
////	ActionIndexError:		{Current: StateError,			Want: StateError, Action: ActionPublish},
////}

