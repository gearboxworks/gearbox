package box

import (
	"github.com/gearboxworks/go-status"
)

func (me *Box) EnsureNotNil() (sts status.Status) {
	if me == nil {
		sts = status.Fail().SetMessage("unexpected software error")
	}
	return sts
}
