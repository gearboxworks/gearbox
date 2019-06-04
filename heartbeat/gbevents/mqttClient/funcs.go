package mqttClient

import (
	"errors"
	"fmt"
	"gearbox/only"
)

func (me *MqttClient) EnsureNotNil() error {

	var err error

	if me == nil {
		err = errors.New("MQTT client instance nil")
	}

	return err
}

func EnsureNotNil(me *MqttClient) error {

	var err error

	if me == nil {
		err = errors.New("MQTT client instance nil")
	}

	return err
}



func (me *ServicesMap) EnsureNotNil() error {
	var err error

	if me == nil {
		err = errors.New("unexpected software error")
	}

	return err
}


func (me *Service) EnsureNotNil() error {
	var err error

	if me == nil {
		err = errors.New("no zeroconf service defined")
	}

	if (me.instance == nil) && (me.IsManaged == true) {
		err = errors.New("no zeroconf service instance defined")
	}

	return err
}


func (me *ServicesArray) Print() error {

	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("software error")
			break
		}

		for i, _ := range *me {
			fmt.Printf("# Entry: #%d\n", i)
			//err = e.Print()
			//if err != nil {
			//	break
			//}
		}
	}

	return err
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
