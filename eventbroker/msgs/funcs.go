package msgs

import (
	"errors"
	"fmt"
	"github.com/gearboxworks/go-status/only"
	"github.com/google/uuid"
	"reflect"
)

// @TODO "Make" is shorter than "Generate" and means the same
//
func MakeAddress() Address {
	return Address(uuid.New().String())
}

func MakeError(me Address, msg string, args ...interface{}) error {
	if me == "" {
		return errors.New(fmt.Sprintf(msg, args...))
	} else {
		return errors.New(fmt.Sprintf(me.String()+": "+msg, args...))
	}
}

func SprintfTopic(address Address, topic SubTopic) string {

	return fmt.Sprintf(TopicPattern, address.String(), topic.String())
}

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
