package mqttBroker

import (
	"errors"
)

func (me *MqttBroker) EnsureNotNil() error {

	var err error

	switch {
		case me == nil:
			err = errors.New("MqttBroker instance is nil")
	}

	return err
}
func EnsureNotNil(me *MqttBroker) error {
	return me.EnsureNotNil()
}


//func (me *ServicesMap) EnsureNotNil() error {
//
//	var err error
//
//	switch {
//		case me == nil:
//			err = errors.New("MqttBroker ServicesMap instance is nil")
//	}
//
//	return err
//}
//func EnsureServicesMapNotNil(me *ServicesMap) error {
//	return me.EnsureNotNil()
//}
//
//
//func (me *Service) EnsureNotNil() error {
//	var err error
//
//	switch {
//		case me == nil:
//			err = errors.New("MqttBroker Service instance is nil")
//		case (me.IsManaged == true) && me.instance == nil:
//			err = me.EntityId.ProduceError("service cmd instance nil")
//	}
//
//	return err
//}
//func EnsureServiceNotNil(me *Service) error {
//	return me.EnsureNotNil()
//}
//
//
//func (me *ServicesMap) Print() error {
//
//	var err error
//
//	for range only.Once {
//		err = me.EnsureNotNil()
//		if err != nil {
//			break
//		}
//
//		for u, s := range *me {
//			fmt.Printf("# Entry: %s\n", u)
//			err = s.Print()
//			if err != nil {
//				break
//			}
//		}
//	}
//
//	return err
//}
//
//
//func (me *Service) Print() error {
//
//	var err error
//
//	for range only.Once {
//		if me == nil {
//			err = errors.New("software error")
//			break
//		}
//
//		if (me.instance == nil) && (me.IsManaged == true) {
//			fmt.Printf("# Entry(deleted): %v", me.EntityId)
//		} else {
//			fmt.Printf("# Entry: %v", me.EntityId)
//		}
//		//err = me.Entry.Print()
//		//if err != nil {
//		//	break
//		//}
//	}
//
//	return err
//}
