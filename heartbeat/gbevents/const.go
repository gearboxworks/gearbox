package gbevents

import (
	"fmt"
	"gearbox/heartbeat/gbevents/channels"
	"gearbox/heartbeat/gbevents/gbMqttBroker"
	"gearbox/heartbeat/gbevents/mqttClient"
	"gearbox/heartbeat/gbevents/network"
	oss "gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"github.com/olebedev/emitter"
)


type EventBroker struct {
	// EntityId messages.MessageAddress
	EntityId   string
	Boxname    string
	PidFile    string

	ZeroConf   network.ZeroConf
	MqttBroker gbMqttBroker.Mqtt
	MqttClient mqttClient.MqttClient
	Channels   channels.Channels
	StsEmitter status.Status

	OsSupport  oss.OsSupporter
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
	Create() status.Status
	Start() status.Status
	Stop() status.Status
	Restart() status.Status
	Status() status.Status
}


const (
	unknownState = "unknown"
	defaultPidFile = "gbevents.pid"
)


func (me ServiceState) String() string {
	fmt.Printf("String\n")
	return string(me)
}
