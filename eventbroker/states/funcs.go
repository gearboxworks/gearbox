package states

import (
	"encoding/json"
	"errors"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/msgs"
	"github.com/gearboxworks/go-status/only"
	"reflect"
	"sync"
	"time"
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

func New(client msgs.Address, name msgs.Address, parent msgs.Address) *Status {

	var ret Status

	if name == "" || name == entity.SelfEntityName {
		name = client
	}

	ret = Status{
		EntityId:   &client,
		EntityName: &name,
		ParentId:   &parent,
		Action:     ActionIdle,
		Want:       StateIdle,
		Current:    StateIdle,
		Last:       StateIdle,
		LastWhen:   time.Now(),
		Error:      nil,
		mutex:      &sync.RWMutex{},
	}

	return &ret
}

func EnsureNotNil(me *Status) error {
	return me.EnsureNotNil()
}

func (me *Status) EnsureNotNil() error {

	var err error

	switch {
	case me == nil:
		err = errors.New("status.Status is nil")

	//case me.mutex == nil:
	//	err = errors.New("status.mutex is nil")

	case me.EntityId == nil:
		err = errors.New("status.EntityId is nil")
	}

	return err
}

func (me *Status) ToMessageText() msgs.Text {

	var err error
	var j []byte

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			j = []byte("{}")
			break
		}

		j, err = json.Marshal(me)
		if err != nil {
			break
		}
	}

	return msgs.Text(j)
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
