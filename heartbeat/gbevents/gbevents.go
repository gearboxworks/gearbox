package gbevents

import (
	"fmt"
	"gearbox/box"
	"gearbox/global"
	"gearbox/heartbeat/gbevents/gbChannels"
	"gearbox/heartbeat/gbevents/gbMqttBroker"
	"gearbox/heartbeat/gbevents/gbMqttClient"
	"gearbox/help"
	"gearbox/only"
	oss "gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/jinzhu/copier"
	"path/filepath"
)


func New(OsSupport oss.OsSupporter, args ...Args) (*EventBroker, status.Status) {

	var _args Args
	var sts status.Status
	se := &EventBroker{}

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		if _args.Boxname == "" {
			_args.Boxname = global.Brandname
		}

		_args.OsSupport = OsSupport
		foo := box.Args{}
		err := copier.Copy(&foo, &_args)
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("unable to configure event handler").
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		sts = _args.Channels.New(OsSupport, gbChannels.Args{})
		if is.Error(sts) {
			break
		}

		sts = _args.MqttBroker.New(OsSupport, gbMqttBroker.Args{})
		if is.Error(sts) {
			break
		}

		sts = _args.MqttClient.New(OsSupport, gbMqttClient.Args{})
		if is.Error(sts) {
			break
		}

		_args.MqttClient.Server.Host = fmt.Sprintf("%s:%d", _args.MqttBroker.Config.Host, _args.MqttBroker.Config.Port)

		_args.PidFile = filepath.FromSlash(fmt.Sprintf("%s/%s", _args.OsSupport.GetAdminRootDir(), defaultPidFile))

		*se = EventBroker(_args)

		sts = status.Success("created new event broker")
	}
	// status.Log(sts)

	return se, sts
}


func (me *EventBroker) RegisterService(topic string, args ...ServiceData) {
	fmt.Printf("RegisterService\n")

	// .

	return
}


func (me *EventBroker) Create() status.Status {
	fmt.Printf("(me *EventBroker) CreateService() status.Status\n")

	return nil
}


func (me *EventBroker) Start() status.Status {

	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		sts = status.Success("me.EventBroker.Start() - DEBUG")
		break

		// Start the inter-thread channels service.
		go func() {
			sts := me.Channels.Start()
			sts.Log()
		}()

		// Start the inter-process service.
		go func() {
			sts := me.MqttBroker.Start()
			sts.Log()
		}()

		// Start the inter-process service client.
		go func() {
			sts := me.MqttClient.Start()
			sts.Log()
		}()

		sts = status.Success("started event broker")
	}

	return sts
}


func (me *EventBroker) Stop() status.Status {
	fmt.Printf("(me *EventBroker) StopService() status.Status\n")

	return nil
}


func (me *EventBroker) Restart() status.Status {
	fmt.Printf("(me *EventBroker) RestartService() status.Status\n")

	return nil
}


func (me *EventBroker) Status() status.Status {
	fmt.Printf("(me *EventBroker) ServiceStatus() status.Status\n")

	return nil
}
