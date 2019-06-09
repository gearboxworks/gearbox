package states

import (
	"errors"
	"gearbox/only"
)


func InterfaceToTypeStatus(i interface{}) (*Status, error) {

	var err error
	var zc *Status

	for range only.Once {
		if i == nil {
			err = errors.New("status.Status is nil")
			break
		}
		zc = i.(*Status)
		// zc = (i[0]).(*ZeroConf)
		// zc = i[0].(*ZeroConf)
	}

	return zc, err
}


func InterfaceToTypeError(i interface{}) (*error, error) {

	var err error
	var zc *error

	for range only.Once {
		if i == nil {
			err = errors.New("error is nil")
			break
		}
		zc = i.(*error)
		// zc = (i[0]).(*ZeroConf)
		// zc = i[0].(*ZeroConf)
	}

	return zc, err
}


func EnsureArgumentNotNil(me *Status) error {

	var err error

	switch {
		case me == nil:
			err = errors.New("status.Status is nil")

		default:
			// err = errors.New("subscriber not nil")
	}

	return err
}
