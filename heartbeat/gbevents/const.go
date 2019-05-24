package gbevents

import (
	"fmt"
	oss "gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"github.com/olebedev/emitter"
	"net/url"
	"time"
)


type ServiceEvents struct {
	Identifier string
	Boxname    string
	PidFile    string
	StsMqtt    status.Status
	StsEmitter status.Status

	mqttServer url.URL
	emitter    emitter.Emitter
	events     emitter.Event
	emits      chan struct{}
	group      emitter.Group

	OsSupport  oss.OsSupporter
}
type Args ServiceEvents

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
var _ Service = (*ServiceEvents)(nil)

var Instance Service

type Service interface {
	Create() status.Status
	Start() status.Status
	Stop() status.Status
	Restart() status.Status
	Status() status.Status
}

type Message struct {
	Src string
	Time time.Time
	Text string
}


const (
	unknownState = "unknown"
	defaultPidFile = "gbevents.pid"
)


func (me ServiceState) String() string {
	fmt.Printf("String\n")
	return string(me)
}

