package messages

import (
	"errors"
	"gearbox/only"
	"reflect"
)

func InterfaceToTypeSubTopics(i interface{}) (*SubTopics, error) {

	var err error
	var zc *SubTopics

	for range only.Once {
		if i == nil {
			err = errors.New("status.Status is nil")
			break
		}

		checkType := reflect.ValueOf(i)
		//fmt.Printf("InterfaceToTypeSubTopics = %v\n", checkType.Type().String())
		if checkType.Type().String() != InterfaceTypeSubTopics {
			err = errors.New("interface type not " + InterfaceTypeSubTopics)
			break
		}

		zc = i.(*SubTopics)
	}

	return zc, err
}

