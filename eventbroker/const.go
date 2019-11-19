package eventbroker

import (
	"gearbox/eventbroker/channels"
	"gearbox/eventbroker/daemon"
	"gearbox/eventbroker/eblog"
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/mqttClient"
	"gearbox/eventbroker/network"
	"gearbox/eventbroker/ospaths"
	"gearbox/eventbroker/states"
	"github.com/olebedev/emitter"
	"sync"
	"time"
)

const (
	Package                  = "eventbroker"
	InterfaceTypeEventBroker = "*" + Package + ".EventBroker"
	DefaultEntityName        = "eventbroker"
)


type Event emitter.Event
var _ EventService = (*EventBroker)(nil)
var Instance EventService
type EventService interface {
	Create() error
	Start() error
	Stop() error
	Restart() error
	Status() error
}


type EventBroker struct {
	EntityId       messages.MessageAddress
	Boxname        string
	SubBaseDir     string
	State          *states.Status

	Channels       channels.Channels
	ZeroConf       network.ZeroConf
	Daemon         daemon.Daemon
	MqttClient     mqttClient.MqttClient

	Services       Services
	Logger         *eblog.Logger

	OsPaths        *ospaths.BasePaths
	channelHandler *channels.Subscriber
}
type Args EventBroker


type Callback func(args interface{}, state states.Status) error

//type States map[messages.MessageAddress]states.Status
//type Callbacks map[messages.MessageAddress]Callback
//type CallbackLocks map[messages.MessageAddress]sync.RWMutex // Mutex control for map
type Log struct {
	When  time.Time
	State states.Status
	//states.State
}
const LogSize = 128
type Logs []Log
type Service struct {
	State         *states.Status
	Callback      Callback
	Args          interface{}
	Logs          Logs

	mutex         sync.RWMutex	// Mutex control for map.
}
type Services map[messages.MessageAddress]*Service


// type Callback func(state states.Status) error
//type States map[messages.MessageAddress]states.Status
//type Callbacks map[messages.MessageAddress]Callback
//type CallbackLocks map[messages.MessageAddress]sync.RWMutex // Mutex control for map
//type Log struct {
//	When  time.Time
//	State states.Status
//	//states.State
//}
//const LogSize = 128
//type Logs []Log
//type Services struct {
//	States        States
//	Callbacks     Callbacks
//	CallbackLocks CallbackLocks
//	Logs          Logs
//
//	mutex         sync.RWMutex	// Mutex control for map.
//}
//type ServiceAction struct {
//	State	 *states.Status
//	CallBack interface{}
//}
//
//func (me *ServiceAction) String() string {
//	return ""
//}
//type Entities map[messages.MessageAddress]*Entity
//type EntityLog []Entities
//type Entity struct {
//	State	 states.Status
//	//State	 states.State
//	CallBack interface{}
//}
//type ServiceDataEntry struct {
//	When time.Time
//	states.Status
//}
//type ServiceDataLog []ServiceDataEntry
//type RegisterServices []ServiceData
//type RegisterServicesMap map[string]*ServiceData

