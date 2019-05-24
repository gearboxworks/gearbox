package gbevents

import (
	"fmt"
	"gearbox/global"
	"gearbox/help"
	"gearbox/only"
	"github.com/fhmq/hmq/broker"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"os"
)


func (me *ServiceEvents) startMqttBroker() status.Status {

	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		fmt.Printf("GBevents - MQTT broker(STARTED)\n")
		config, err := broker.ConfigureConfig(os.Args[1:])
		if err != nil {
			sts = status.Fail(&status.Args{
				Message: "GBevents - MQTT broker config error",
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		b, err := broker.NewBroker(config)
		if err != nil {
			sts = status.Fail(&status.Args{
				Message: "GBevents - MQTT broker error",
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}
		b.Start()

		s := waitForSignal()

		fmt.Printf("GBevents - MQTT broker(FINISHED)\n")
		sts = status.Success("%s GBevents - Poller exited with signal %v.", global.Brandname, s)
	}

	return sts
}
