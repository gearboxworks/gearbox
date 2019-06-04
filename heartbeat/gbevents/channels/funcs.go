package channels

import (
	"fmt"
	"github.com/getlantern/errors"
)

func (me *Channels) EnsureNotNil() error {

	var err error
	//var emptyChannelInstance channelsInstance

	switch {
		case me == nil:
			err = errors.New("channels instance is nil")
			fmt.Printf("FO\n")

		//case me.instance == emptyChannelInstance:
		//	err = errors.New("Funexpected software error")
		//	fmt.Printf("FO\n")
	}

	return err
}

func EnsureNotNil(me *Channels) error {

	var err error

	switch {
		case me == nil:
			err = errors.New("channels instance is nil")
		//case me.instance.emitter.Cap == nil:
		//	err = errors.New("channels instance is nil")
	}

	return err
}
