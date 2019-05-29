package gbMqttBroker

import (
	"fmt"
	"gearbox/box"
	"gearbox/heartbeat/daemon"
	"gearbox/help"
	"gearbox/only"
	oss "gearbox/os_support"
	"github.com/fhmq/hmq/broker"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/jinzhu/copier"
	"net/url"
	"os"
)


func (me *Broker) New(OsSupport oss.OsSupporter, args ...Args) status.Status {

	var _args Args
	var sts status.Status

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		_args.OsSupport = OsSupport
		foo := box.Args{}
		err := copier.Copy(&foo, &_args)
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("unable to copy MQTT broker config").
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		if _args.EntityId == "" {
			_args.EntityId = "gearbox-mqtt-broker"
		}

		_args.Config, err = broker.ConfigureConfig(os.Args[1:])
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("unable to parse MQTT client config").
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}
		_args.Config = broker.DefaultConfig

		_args.Server, err = url.Parse(fmt.Sprintf("tcp://%s:%s/", _args.Config.Host, _args.Config.Port))
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("unable to parse MQTT broker url").
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}
		//_args.Config.Host = me.Server.Host

		_args.instance, err = broker.NewBroker(_args.Config)
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("error starting MQTT broker").
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		// _args.Server = fmt.Sprintf("tcp://%s:%s/", _args.Config.Host, _args.Config.Port)

		*me = Broker(_args)
	}

	return sts
}


func (me *Broker) Start() status.Status {

	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		status.Success("GBevents - MQTT broker(STARTED)").Log()
		me.instance.Start()

		s := daemon.WaitForSignal()

		status.Success("GBevents - MQTT broker(STOPPED)").Log()
		sts = status.Success("MQTT broker exited with signal %v.", s)
	}
	// status.Log(sts.SetLogTo(status.DebugLog))
	me.Sts = sts
	status.Log(sts)

	return sts
}
