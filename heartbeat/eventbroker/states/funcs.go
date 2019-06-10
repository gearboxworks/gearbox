package states

import (
	"errors"
	"gearbox/only"
	"reflect"
)


func InterfaceToTypeStatus(i interface{}) (*Status, error) {

	var err error
	var zc *Status

	for range only.Once {
		if i == nil {
			err = errors.New("status.Status is nil")
			break
		}

		checkType := reflect.ValueOf(i)
		if checkType.Type().String() != "*states.Status" {
			err = errors.New("interface type not *status.Status")
			// fmt.Printf("%v\n", f.GetError())
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

		checkType := reflect.ValueOf(i)
		if checkType.Type().String() != "*error" {
			err = errors.New("interface type not *error")
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
