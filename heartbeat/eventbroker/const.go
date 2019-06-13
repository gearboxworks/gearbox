package eventbroker

import (
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/daemon"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/mqttClient"
	"gearbox/heartbeat/eventbroker/network"
	"gearbox/heartbeat/eventbroker/ospaths"
	"gearbox/heartbeat/eventbroker/states"
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

	OsPaths        *ospaths.BasePaths
	channelHandler *channels.Subscriber
}
type Args EventBroker



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
type States map[messages.MessageAddress]*states.Status
type Callbacks map[messages.MessageAddress]interface{}
type Log struct {
	When  time.Time
	State states.Status
	//states.State
}
const LogSize = 128
type Logs []Log
type Services struct {
	States    States
	Callbacks Callbacks
	Logs      Logs

	mutex     sync.RWMutex	// Mutex control for map.
}



//type ServiceDataEntry struct {
//	When time.Time
//	states.Status
//}
//type ServiceDataLog []ServiceDataEntry
//type RegisterServices []ServiceData
//type RegisterServicesMap map[string]*ServiceData

