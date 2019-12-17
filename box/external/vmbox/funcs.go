package vmbox

import (
	"errors"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/states"
	"github.com/gearboxworks/go-status/only"
	"reflect"
)

func EnsureNotNil(me *VmBox) error {
	return me.EnsureNotNil()
}

func InterfaceToTypeVmBox(i interface{}) (*VmBox, error) {

	var err error
	var zc *VmBox

	for range only.Once {
		if i == nil {
			err = errors.New("interface is nil, should be" + InterfaceTypeVmBox)
			break
		}

		checkType := reflect.ValueOf(i)
		if checkType.Type().String() != InterfaceTypeVmBox {
			err = errors.New("interface type not " + InterfaceTypeVmBox)
			break
		}

		zc = i.(*VmBox)

		err = zc.EnsureNotNil()
		if err != nil {
			break
		}
	}

	return zc, err
}

func InterfaceToTypeService(i interface{}) (*Vm, error) {

	var err error
	var s *Vm

	for range only.Once {
		if i == nil {
			err = errors.New("interface is nil, should be" + InterfaceTypeService)
			break
		}

		checkType := reflect.ValueOf(i)
		if checkType.Type().String() != InterfaceTypeService {
			err = errors.New("interface type not " + InterfaceTypeService)
			break
		}

		s = i.(*Vm)

		err = s.EnsureNotNil()
		if err != nil {
			break
		}
	}

	return s, err
}

func ConstructVmMessage(me msgs.Address, to msgs.Address, a states.Action) msgs.Message {

	var err error
	var msgTemplate msgs.Message

	for range only.Once {
		err = me.EnsureNotEmpty()
		if err != nil {
			break
		}

		msgTemplate = msgs.Message{
			Source: me,
			Topic: msgs.Topic{
				Address:  entity.VmBoxEntityName,
				SubTopic: msgs.SubTopic(a),
			},
			Text: msgs.Text(to),
		}
	}

	return msgTemplate
}
