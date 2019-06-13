package mqttBroker

import (
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/heartbeat/eventbroker/tasks"
	oss "gearbox/os_support"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"net/url"
	"time"
)


const (
	DefaultEntityId = "eventbroker-mqttbroker"
	defaultWaitTime = time.Millisecond * 2000
	defaultDomain   = "local"
	DefaultRetries  = 12
	DefaultRetryDelay = time.Second * 3
	DefaultServer = "tcp://127.0.0.1:1883"
)


const mqttBrokerJson = `
{
	"Name": "com.gearbox.mqttbroker",
	"DisplayName": "Gearbox [MQTT broker]",
	"Description": "Provides network based communications channels",

	"Url": "tcp://{{.Host}}:{{.Port}}",

	"WorkingDirectory": "{{.GetUserHomeDir}}/.gearbox/admin/dist/eventbroker/mqttbroker",
	"Executable": "{{.GetUserHomeDir}}/.gearbox/admin/dist/eventbroker/mqttbroker/bin/mqttbroker",
	"Arguments": [
		"-i", "{{.GetUserHomeDir}}/.gearbox/admin/dist/eventbroker/mqttbroker/mqttbroker.pid",
		"--host", "{{.GetHost}}",
		"--port", "{{.GetPort}}",
		"--config", "{{.GetUserHomeDir}}/.gearbox/admin/dist/eventbroker/mqttbroker/etc/mqttbroker.cfg",
		"-d"
	],
	"Env": [
		"PROJECTS={{.GetUserHomeDir}}/Sites",
		"HOME={{.GetUserHomeDir}}",
		"HOMEDIR={{.GetUserHomeDir}}",
		"ADMINROOTDIR={{.GetAdminRootDir}}",
		"CACHEDIR={{.GetCacheDir}}",
		"BASEDIR={{.GetUserHomeDir}}",
		"USERCONFIGDIR={{.GetUserConfigDir}}"
	],

	"Option": {
		"KeepAlive": true,
		"RunAtLoad": false,
		"SessionCreate": false,
		"UserService": true
	},

	"Stderr": "{{.GetStderr}}",
	"Stdout": "{{.GetStdout}}"
}`

type MqttBroker struct {
	EntityId        messages.MessageAddress
	State           *states.Status
	Task            *tasks.Task
	Server          *url.URL
	Channels        *channels.Channels
	channelHandler  *channels.Subscriber

	osSupport       oss.OsSupporter
}
type Args MqttBroker

type CreateEntry struct {
	Name   string	`json:"name"`	// == Service.Entry.Instance
	Topic  Topic	`json:"topic"`
	TTL    uint32   `json:"ttl"`	// == Service.Entry.TTL
	Qos    byte		`json:"qos"`
	callback mqtt.MessageHandler
}

type Topic string
func (me *Topic) String() (string) {

	return string(*me)
}


//
//
//// Ensure we don't duplicate services.
//func (me *Service) IsExisting(him CreateEntry) error {
//
//	var err error
//
//	switch {
//		case me.Entry.Topic == him.Topic:
//			err = me.EntityId.ProduceError("MqttBroker service Topic:%s already exists", me.Entry.Topic)
//
//		case me.Entry.Name == him.Name:
//			err = me.EntityId.ProduceError("MqttBroker service Name:%s already exists", me.Entry.Name)
//	}
//
//	return err
//}
//
//// Ensure we don't duplicate services.
//func (me *ServicesMap) IsExisting(him CreateEntry) error {
//
//	var err error
//
//	for _, ce := range *me {
//		err = ce.IsExisting(him)
//		if err != nil {
//			break
//		}
//	}
//
//	return err
//}
//
//
//func InterfaceToTypeMqttBroker(i interface{}) (*MqttBroker, error) {
//
//	var err error
//	var zc *MqttBroker
//
//	for range only.Once {
//		err = channels.EnsureArgumentNotNil(i)
//		if err != nil {
//			break
//		}
//		zc = i.(*MqttBroker)
//		// zc = (i[0]).(*ZeroConf)
//		// zc = i[0].(*ZeroConf)
//
//		err = zc.EnsureNotNil()
//		if err != nil {
//			break
//		}
//	}
//
//	return zc, err
//}
//
//
//func InterfaceToTypeService(i interface{}) (*Service, error) {
//
//	var err error
//	var s *Service
//
//	for range only.Once {
//		err = channels.EnsureArgumentNotNil(i)
//		if err != nil {
//			break
//		}
//		s = i.(*Service)
//		// zc = (i[0]).(*Service)
//		// zc = i[0].(*Service)
//
//		err = s.EnsureNotNil()
//		if err != nil {
//			break
//		}
//	}
//
//	return s, err
//}
