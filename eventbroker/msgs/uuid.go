package msgs

import (
	"errors"
	"github.com/google/uuid"
)

type Uuid uuid.UUID

func (me *Uuid) String() string {
	if me == nil {
		return ""
	}
	return uuid.UUID(*me).String()
}

func (me *Uuid) ToMessageText() Text {
	return Text(me.String())
}

func (me *Uuid) EnsureNotNil() (err error) {
	if me == nil {
		err = errors.New("message address uuid is nil")
	}
	return err
}
