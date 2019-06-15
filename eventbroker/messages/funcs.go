package messages

import (
	"errors"
	"gearbox/eventbroker/only"
	"reflect"
)

func InterfaceToTypeSubTopics(i interface{}) (*SubTopics, error) {

	var err error
	var zc *SubTopics

	for range only.Once {
		if i == nil {
			err = errors.New("interface is nil, should be" + InterfaceTypeSubTopics)
			break
		}

		checkType := reflect.ValueOf(i)
		if checkType.Type().String() != InterfaceTypeSubTopics {
			err = errors.New("interface type not " + InterfaceTypeSubTopics)
			break
		}

		zc = i.(*SubTopics)
	}

	return zc, err
}

