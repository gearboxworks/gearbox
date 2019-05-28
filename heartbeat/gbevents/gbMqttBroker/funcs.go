package gbMqttBroker

import (
	"gearbox/help"
	"github.com/gearboxworks/go-status"
)

func EnsureNotNil(bx *Broker) (sts status.Status) {
	if bx == nil {
		sts = status.Fail().
			SetMessage("unexpected software error").
			SetAdditional("", ).
			SetData("").
			SetHelp(status.AllHelp, help.ContactSupportHelp())
	}
	return sts
}
