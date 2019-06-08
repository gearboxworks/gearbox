package mqttClient

import (
	"errors"
	"fmt"
	"gearbox/only"
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
