package eventbroker

import (
	"fmt"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/daemon"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/mqttClient"
	"gearbox/heartbeat/eventbroker/network"
	"gearbox/heartbeat/eventbroker/ospaths"
	"gearbox/heartbeat/eventbroker/states"
	"github.com/olebedev/emitter"
)

const (
	unknownState = "unknown"
	DefaultEntityName = "eventbroker"
	defaultPidFile = "gbevents.pid"
	DefaultBaseDir = "dist/eventbroker"
)


type EventBroker struct {
	EntityId       messages.MessageAddress
	Boxname        string
	SubBaseDir     string
	State          states.Status

	Channels       channels.Channels
	ZeroConf       network.ZeroConf
	Daemon         daemon.Daemon
	MqttClient     mqttClient.MqttClient

	Entities       Entities

	OsPaths        *ospaths.BasePaths
	channelHandler *channels.Subscriber
}
type Args EventBroker

type ServiceData struct {
	Name	string
	State	ServiceState
	Action  ServiceAction
}
type RegisterServices []ServiceData
type RegisterServicesMap map[string]*ServiceData

type ServiceState string

type ServiceAction struct {
	State	 ServiceState
	CallBack interface{}
}

type Entity struct {
	Src   messages.MessageAddress
	State *states.Status
	StateString  states.State
}
type Entities map[messages.MessageAddress]*Entity


const (
	Package                  = "eventbroker"
	InterfaceTypeEventBroker = "*" + Package + ".EventBroker"
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


func (me ServiceState) String() string {
	fmt.Printf("String\n")
	return string(me)
}
