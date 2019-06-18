package gbMqttBroker

import (
	"gearbox/help"
	"github.com/gearboxworks/go-status"
)

func (me *Mqtt) EnsureNotNil() (sts status.Status) {
	if me == nil {
		sts = status.Fail().
			SetMessage("unexpected software error").
			SetAdditional("", ).
			SetData("").
			SetHelp(status.AllHelp, help.ContactSupportHelp())
	}
	return sts
}


func (me *Mqtt) EnsureNotNil2() (sts status.Status) {
	if me == nil {
		sts = status.Fail().
			SetMessage("unexpected software error").
			SetAdditional("", ).
			SetData("").
			SetHelp(status.AllHelp, help.ContactSupportHelp())
	}
	return sts
}
