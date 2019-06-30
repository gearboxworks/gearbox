package msgs

import (
	"encoding/json"
	"errors"
	"github.com/gearboxworks/go-status/only"
)

type Text string

func (me *Text) EnsureNotNil() error {

	var err error

	switch {
	case me == nil:
		err = errors.New("message text is nil")

	case *me == "":
		err = errors.New("message text is empty")
	}

	return err
}

func (me *Text) String() string {

	return string(*me)
}

func (me *Text) ByteArray() []byte {

	return []byte(*me)
}

func (me *Text) ToMessage() (*Message, error) {

	var err error
	var ret Message

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = json.Unmarshal(me.ByteArray(), &ret)
		if err != nil {
			break
		}

		err = ret.Validate()
	}

	return &ret, err
}

func (me *Text) ToAddress() Address {
	return Address(me.String())
}
