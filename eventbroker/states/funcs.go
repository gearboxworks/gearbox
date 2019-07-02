package states

import (
	"encoding/json"
	"errors"
	"gearbox/eventbroker/msgs"
	"github.com/gearboxworks/go-status/only"
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
		if zc == nil {
			err = errors.New("interface type not " + InterfaceTypeStatus + " is nil")
			break
		}
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

func EnsureNotNil(me *Status) error {
	return me.EnsureNotNil()
}

func FromMessageText(me msgs.Text) (*Status, error) {

	var err error
	var ret Status

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = json.Unmarshal([]byte(me), &ret)
		if err != nil {
			break
		}

		//err = ret.Validate()
	}

	return &ret, err
}
