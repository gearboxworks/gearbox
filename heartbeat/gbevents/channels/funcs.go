package channels

import (
	"github.com/getlantern/errors"
)

func (me *Channels) EnsureNotNil() error {

	var err error

	if me == nil {
		err = errors.New("unexpected software error")
	}

	return err
}
