package states

import (
	"gearbox/eventbroker/messages"
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
	StateGetting       = "getting"
	StatePutting       = "putting"
	StateUpdating      = "updating"

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
	StateIndexGetting       = iota
	StateIndexPutting       = iota
	StateIndexUpdating      = iota
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
	StateIndexGetting:       StateGetting,
	StateIndexPutting:       StatePutting,
	StateIndexUpdating:      StateUpdating,
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
	ActionAction      = "action"
	ActionGet         = "get"
	ActionPut         = "put"
	ActionUpdate      = "update"
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
	ActionIndexAction      = iota
	ActionIndexGet         = iota
	ActionIndexPut         = iota
	ActionIndexUpdate      = iota
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
	ActionIndexAction:      ActionAction,
	ActionIndexGet:         ActionGet,
	ActionIndexPut:         ActionPut,
	ActionIndexUpdate:      ActionUpdate,
	ActionIndexError:       ActionError,
	ActionIndexGlob:        ActionGlob,
}


type Status struct {
	EntityId   *messages.MessageAddress
	EntityName *messages.MessageAddress
	ParentId   *messages.MessageAddress
	Current    State
	Want       State
	Last       State
	LastWhen   time.Time
	Attempts   int
	Error      error
	Action     Action

	mutex      *sync.RWMutex // Mutex control for map.
}

const (
	Package					= "states"
	InterfaceTypeStatus		= "*" + Package + ".Status"
	InterfaceTypeError		= "*error"
)

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

