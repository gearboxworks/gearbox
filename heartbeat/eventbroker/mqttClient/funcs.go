package mqttClient

import (
	"errors"
	"fmt"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/only"
	"reflect"
)

func (me *MqttClient) EnsureNotNil() error {

	var err error

	switch {
		case me == nil:
			err = errors.New("MqttClient instance is nil")
		//case me.instance.client == nil:
		//	err = me.EntityId.ProduceError("client instance nil")
		//case me.instance.token == nil:
		//	err = me.EntityId.ProduceError("client token nil")
		case me.instance.options == nil:
			err = me.EntityId.ProduceError("client options nil")
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
			err = me.EntityId.ProduceError("service cmd instance nil")
	}

	return err
}
func EnsureServiceNotNil(me *Service) error {
	return me.EnsureNotNil()
}


// Ensure we don't duplicate services.
func (me *Service) IsExisting(him CreateEntry) error {

	var err error

	switch {
		case me.Entry.Topic == him.Topic:
			err = me.EntityId.ProduceError("MqttClient service Topic:%s already exists", me.Entry.Topic)

		case me.Entry.Name == him.Name:
			err = me.EntityId.ProduceError("MqttClient service Name:%s already exists", me.Entry.Name)
	}

	return err
}


// Ensure we don't duplicate services.
func (me *ServicesMap) IsExisting(him CreateEntry) error {

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
		err = channels.EnsureArgumentNotNil(i)
		if err != nil {
			break
		}

		checkType := reflect.ValueOf(i)
		if checkType.Type().String() != "*mqttClient.MqttClient" {
			err = errors.New("interface type not *mqttClient.MqttClient")
			break
		}

		zc = i.(*MqttClient)
		// zc = (i[0]).(*ZeroConf)
		// zc = i[0].(*ZeroConf)

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
		err = channels.EnsureArgumentNotNil(i)
		if err != nil {
			break
		}

		checkType := reflect.ValueOf(i)
		if checkType.Type().String() != "*mqttClient.Service" {
			err = errors.New("interface type not *mqttClient.MqttClient")
			break
		}

		s = i.(*Service)
		// zc = (i[0]).(*Service)
		// zc = i[0].(*Service)

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
