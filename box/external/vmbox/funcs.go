package vmbox

import (
	"errors"
	"fmt"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/states"
	"github.com/gearboxworks/go-status/only"
	"reflect"
)


func (me *VmBox) EnsureNotNil() error {

	var err error

	switch {
		case me == nil:
			err = errors.New("VmBox instance is nil")
	}

	return err
}
func EnsureNotNil(me *VmBox) error {
	return me.EnsureNotNil()
}


func (me *VmMap) EnsureNotNil() error {

	var err error

	switch {
		case me == nil:
			err = errors.New("VmBox ServicesMap instance is nil")
	}

	return err
}
func EnsureServicesMapNotNil(me *VmMap) error {
	return me.EnsureNotNil()
}


func (me *Vm) EnsureNotNil() error {
	var err error

	switch {
		case me == nil:
			err = errors.New("VmBox Service instance is nil")
	}

	return err
}
func EnsureServiceNotNil(me *Vm) error {
	return me.EnsureNotNil()
}


func (me *ServiceConfig) EnsureNotNil() error {
	var err error

	switch {
	case me == nil:
		err = errors.New("VmBox Service instance is nil")
	}

	return err
}
func EnsureServiceConfigNotNil(me *ServiceConfig) error {
	return me.EnsureNotNil()
}


// Ensure we don't duplicate services.
func (me *VmBox) IsExisting(client messages.MessageAddress) *Vm {

	var ret *Vm

	for _, v := range me.vms {
		if v.EntityName == client {
			ret = v
			break
		}

		if v.EntityId == client {
			ret = v
			break
		}
	}

	//if _, ok := me.vms[client]; ok {
	//	ret = me.vms[client]
	//}

	return ret
}


//// Ensure we don't duplicate services.
//func (me *Vm) IsExisting(him ServiceConfig) error {
//
//	var err error
//
//	switch {
//		case me.Entry.Name == him.Name:
//			err = me.EntityId.ProduceError("VmBox service Name:%s already exists", me.Entry.Name)
//	}
//
//	return err
//}
//
//
//// Ensure we don't duplicate services.
//func (me *VmBox) IsExisting(him ServiceConfig) error {
//
//	var err error
//
//	for _, ce := range me.vms {
//		err = ce.IsExisting(him)
//		if err != nil {
//			break
//		}
//	}
//
//	return err
//}


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


func (me *VmMap) Print() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for u, s := range *me {
			fmt.Printf("# Entry: %s\n", u)
			err = s.Print()
			if err != nil {
				break
			}
		}
	}

	return err
}


func (me *Vm) Print() error {

	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("software error")
			break
		}

		if me.IsManaged == true {
			fmt.Printf("# Entry(deleted): %v", me.EntityId)
		} else {
			fmt.Printf("# Entry: %v", me.EntityId)
		}
		//err = me.Entry.Print()
		//if err != nil {
		//	break
		//}
	}

	return err
}


func ConstructVmMessage(me messages.MessageAddress, to messages.MessageAddress, a states.Action) messages.Message {

	var err error
	var msgTemplate messages.Message

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		msgTemplate = messages.Message{
			Source: me,
			Topic: messages.MessageTopic{
				Address:  entity.VmBoxEntityName,
				SubTopic: messages.SubTopic(a),
			},
			Text: messages.MessageText(to),
		}
	}

	return msgTemplate
}

