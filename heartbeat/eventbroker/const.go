package eventbroker

import (
	"fmt"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/daemon"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/mqttClient"
	"gearbox/heartbeat/eventbroker/network"
	"gearbox/heartbeat/eventbroker/states"
	oss "gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"github.com/olebedev/emitter"
)

const (
	unknownState = "unknown"
	DefaultEntityName = "eventbroker"
	defaultPidFile = "gbevents.pid"
)


var ServiceMqtt = network.CreateEntry{
	Name:   "_gearbox-mqtt",
	Type:   "_mqtt._tcp",
	Domain: "local",
}


type EventBroker struct {
	EntityId   messages.MessageAddress
	Boxname    string
	PidFile    string
	State      states.Status

	Channels   channels.Channels
	ZeroConf   network.ZeroConf
	Daemon     daemon.Daemon
	MqttClient mqttClient.MqttClient
	StsEmitter status.Status

	osSupport  oss.OsSupporter
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
