package gbevents

import (
	"fmt"
	"gearbox/heartbeat/gbevents/gbChannels"
	"gearbox/heartbeat/gbevents/gbMqttBroker"
	"gearbox/heartbeat/gbevents/gbZeroConf"
	oss "gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"github.com/olebedev/emitter"
)


type EventBroker struct {
	Identifier string
	Boxname    string
	PidFile    string

	ZeroConf gbZeroConf.Client
	MqttBroker gbMqttBroker.Mqtt
	Channels	gbChannels.Channels
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
type RegisterServicesMap map[string]ServiceData

type ServiceState string

type ServiceAction struct {
	State	 ServiceState
	CallBack interface{}
}


type Event emitter.Event
var _ Service = (*EventBroker)(nil)

var Instance Service

type Service interface {
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
