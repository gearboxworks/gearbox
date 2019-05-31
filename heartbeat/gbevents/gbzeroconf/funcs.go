package gbZeroConf

import (
	"gearbox/help"
	"github.com/gearboxworks/go-status"
)

func (me *Client) EnsureNotNil() (sts status.Status) {
	if me == nil {
		sts = status.Fail().
			SetMessage("unexpected software error").
			SetAdditional("", ).
			SetData("").
			SetHelp(status.AllHelp, help.ContactSupportHelp())
	}

	return sts
}
