package heartbeat

import (
	"gearbox/help"
	"github.com/gearboxworks/go-status"
)

func (me *Heartbeat) EnsureNotNil() (sts status.Status) {
	if me == nil {
		sts = status.Fail().
			SetMessage("unexpected software error").
			SetAdditional("", ).
			SetData("").
			SetHelp(status.AllHelp, help.ContactSupportHelp())
	}
	return sts
}
