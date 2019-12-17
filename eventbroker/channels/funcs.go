package channels

import (
	"errors"
	"gearbox/eventbroker/msgs"
)

func (me *Channels) EnsureNotNil() error {
	var err error

	switch {
	case me == nil:
		err = errors.New("Channels instance is nil")
	case me.instance.emitter == nil:
		err = msgs.MakeError(me.EntityId, "instance.emitter is nil")
		//case me.instance.events == nil:
		//	err = msgs.MakeError(me.EntityId,"instance.events is nil")
		//case me.instance.emits == nil:
		//	err = msgs.MakeError(me.EntityId,"instance.emits is nil")
		//case me.instance.group == nil:
		//	err = msgs.MakeError(me.EntityId,"instance.group is nil")
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

func (me *Channels) EnsureSubscriberNotNil(client msgs.Address) error {

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	if _, ok := me.subscribers[client]; !ok { // Managed by Mutex
		return msgs.MakeError(me.EntityId, "subscriber doesn't exist")
	} else {
		return me.subscribers[client].EnsureNotNil() // Managed by Mutex
	}
}

func (me *Subscriber) EnsureNotNil() error {

	var err error

	switch {
	case me == nil:
		err = errors.New("subscriber is nil")

	case me.EntityId.EnsureNotEmpty() != nil:
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

// Old functions. Please keep for reference.
//func (me *Channels) off(topic msg.Topic, channels ...<-chan emitter.Event)  {
//	eblog.Debug(me.EntityId, "Off")
//
//	me.instance.emitter.Off(topic.String(), channels...)
//
//	return
//}
//
//
//func (me *Channels) on(topic msg.Topic, middleware ...func(emitter *emitter.Event)) <-chan emitter.Event {
//	eblog.Debug(me.EntityId, "On")
//
//	// me.instance.events = <-me.instance.emitter.On(topic.String(), middleware...)
//	// me.group.Add(me.instance.emitter.On(topic.String()))
//
//	return me.instance.emitter.On(topic.String(), middleware...)
//}
//
//
//func (me *Channels) once(topic msg.Topic, middleware ...func(emitter *emitter.Event)) <-chan emitter.Event {
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
