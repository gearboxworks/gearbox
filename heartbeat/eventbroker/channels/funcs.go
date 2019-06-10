package channels

import (
	"errors"
	"gearbox/only"
	"gearbox/heartbeat/eventbroker/messages"
	"reflect"
)


func (me *Channels) EnsureNotNil() error {
	var err error

	switch {
		case me == nil:
			err = errors.New("Channels instance is nil")
		case me.instance.emitter == nil:
			err = me.EntityId.ProduceError("instance.emitter is nil")
		//case me.instance.events == nil:
		//	err = me.EntityId.ProduceError("instance.events is nil")
		//case me.instance.emits == nil:
		//	err = me.EntityId.ProduceError("instance.emits is nil")
		//case me.instance.group == nil:
		//	err = me.EntityId.ProduceError("instance.group is nil")
	}

	return err
}
func EnsureNotNil(me *Channels) error {
	return me.EnsureNotNil()
}


func EnsureArgumentNotNil(me Argument) error {

	var err error

	switch {
	case me == nil:
		err = errors.New("channel argument is nil")

	default:
		// err = errors.New("subscriber not nil")
	}

	return err
}


func (me *Channels) EnsureSubscriberNotNil(entity messages.MessageAddress) error {

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	if _, ok := me.subscribers[entity]; !ok {	// Managed by Mutex
		return me.EntityId.ProduceError("subscriber doesn't exist")
	} else {
		return me.subscribers[entity].EnsureNotNil()      // Managed by Mutex
	}
}


func (me *Subscriber) EnsureNotNil() error {

	var err error

	switch {
	case me == nil:
		err = errors.New("subscriber is nil")

	case me.EntityId.EnsureNotNil() != nil:
		err = errors.New("subscriber address is nil")
	//
	//case me.Callbacks == nil:
	//	err = errors.New("subscriber callbacks is nil")
	//
	//case me.Returns == nil:
	//	err = errors.New("subscriber returns is nil")

	default:
		// err = errors.New("subscriber not nil")
	}

	return err
}


func ReturnInterfaceType(i interface{}) (string, error) {

	var err error
	var ret string

	for range only.Once {
		err = EnsureArgumentNotNil(i)
		if err != nil {
			break
		}

		checkType := reflect.ValueOf(i)
		ret = checkType.Type().String()
		if ret == "" {
			err = errors.New("interface type is nil")
		}
	}

	return ret, err
}


// Old functions. Please keep for reference.
//func (me *Channels) off(topic messages.MessageTopic, channels ...<-chan emitter.Event)  {
//	eblog.Debug(me.EntityId, "Off")
//
//	me.instance.emitter.Off(topic.String(), channels...)
//
//	return
//}
//
//
//func (me *Channels) on(topic messages.MessageTopic, middleware ...func(emitter *emitter.Event)) <-chan emitter.Event {
//	eblog.Debug(me.EntityId, "On")
//
//	// me.instance.events = <-me.instance.emitter.On(topic.String(), middleware...)
//	// me.group.Add(me.instance.emitter.On(topic.String()))
//
//	return me.instance.emitter.On(topic.String(), middleware...)
//}
//
//
//func (me *Channels) once(topic messages.MessageTopic, middleware ...func(emitter *emitter.Event)) <-chan emitter.Event {
//	eblog.Debug(me.EntityId, "Once")
//
//	// me.instance.events = <-me.instance.emitter.Once(topic.String(), middleware...)
//	// me.instance.events.String(1)
//
//	return me.instance.emitter.Once(topic.String(), middleware...)
//}
//
//
//func (me *Channels) use(pattern string, middleware ...func(emitter *emitter.Event))  {
//	eblog.Debug(me.EntityId, "Use")
//
//	me.instance.emitter.Use(pattern, middleware...)
//
//	return
//}

