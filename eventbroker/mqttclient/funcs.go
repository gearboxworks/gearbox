package mqttClient

import (
	"errors"
	"fmt"
	"gearbox/eventbroker/msgs"
	"github.com/gearboxworks/go-status/only"
	"reflect"
)

func (me *MqttClient) EnsureNotNil() error {

	var err error

	switch {
	case me == nil:
		err = errors.New("MqttClient instance is nil")
	//case me.instance.client == nil:
	//	err = msgs.MakeError(me.EntityId,"client instance nil")
	//case me.instance.token == nil:
	//	err = msgs.MakeError(me.EntityId,"client token nil")
	case me.instance.options == nil:
		err = msgs.MakeError(me.EntityId, "client options nil")
	}

	return err
}
func EnsureNotNil(me *MqttClient) error {
	return me.EnsureNotNil()
}

func (me *ServicesMap) EnsureNotNil() error {

	var err error

	switch {
	case me == nil:
		err = errors.New("MqttClient ServicesMap instance is nil")
	}

	return err
}
func EnsureServicesMapNotNil(me *ServicesMap) error {
	return me.EnsureNotNil()
}

func (me *Service) EnsureNotNil() error {
	var err error

	switch {
	case me == nil:
		err = errors.New("MqttClient Service instance is nil")
	case (me.IsManaged == true) && me.instance == nil:
		err = msgs.MakeError(me.EntityId, "service cmd instance nil")
	}

	return err
}
func EnsureServiceNotNil(me *Service) error {
	return me.EnsureNotNil()
}

// Ensure we don't duplicate services.
func (me *Service) IsExisting(him ServiceConfig) error {

	var err error

	switch {
	case me.Entry.Topic == him.Topic:
		err = msgs.MakeError(me.EntityId, "MqttClient service Topic:%s already exists", me.Entry.Topic)

	case me.Entry.Name == him.Name:
		err = msgs.MakeError(me.EntityId, "MqttClient service Name:%s already exists", me.Entry.Name)
	}

	return err
}

// Ensure we don't duplicate services.
func (me *ServicesMap) IsExisting(him ServiceConfig) error {

	var err error

	for _, ce := range *me {
		err = ce.IsExisting(him)
		if err != nil {
			break
		}
	}

	return err
}

func InterfaceToTypeMqttClient(i interface{}) (*MqttClient, error) {

	var err error
	var zc *MqttClient

	for range only.Once {
		if i == nil {
			err = errors.New("interface is nil, should be" + InterfaceTypeMqttClient)
			break
		}

		checkType := reflect.ValueOf(i)
		if checkType.Type().String() != InterfaceTypeMqttClient {
			err = errors.New("interface type not " + InterfaceTypeMqttClient)
			break
		}

		zc = i.(*MqttClient)

		err = zc.EnsureNotNil()
		if err != nil {
			break
		}
	}

	return zc, err
}

func InterfaceToTypeService(i interface{}) (*Service, error) {

	var err error
	var s *Service

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

		s = i.(*Service)

		err = s.EnsureNotNil()
		if err != nil {
			break
		}
	}

	return s, err
}

func (me *ServicesMap) Print() error {

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

func (me *Service) Print() error {

	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("software error")
			break
		}

		if (me.instance == nil) && (me.IsManaged == true) {
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
