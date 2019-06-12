package states

import (
	"errors"
	"gearbox/heartbeat/eventbroker/only"
	"reflect"
)


func InterfaceToTypeStatus(i interface{}) (*Status, error) {

	var err error
	var zc *Status

	for range only.Once {
		if i == nil {
			err = errors.New("interface is nil, should be" + InterfaceTypeStatus)
			break
		}

		checkType := reflect.ValueOf(i)
		if checkType.Type().String() != InterfaceTypeStatus {
			err = errors.New("interface type not " + InterfaceTypeStatus)
			break
		}

		zc = i.(*Status)
	}

	return zc, err
}

func InterfaceToTypeError(i interface{}) (*error, error) {

	var err error
	var zc *error

	for range only.Once {
		if i == nil {
			err = errors.New("interface is nil, should be" + InterfaceTypeError)
			break
		}

		checkType := reflect.ValueOf(i)
		if checkType.Type().String() != InterfaceTypeError {
			err = errors.New("interface type not " + InterfaceTypeError)
			break
		}

		zc = i.(*error)
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
