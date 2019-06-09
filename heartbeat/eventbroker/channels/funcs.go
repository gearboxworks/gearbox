package channels

import (
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"errors"
	"github.com/olebedev/emitter"
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


func (me *Channels) off(topic messages.MessageTopic, channels ...<-chan emitter.Event)  {
	eblog.Debug(me.EntityId, "Off")

	me.instance.emitter.Off(topic.String(), channels...)

	return
}


func (me *Channels) on(topic messages.MessageTopic, middleware ...func(emitter *emitter.Event)) <-chan emitter.Event {
	eblog.Debug(me.EntityId, "On")

	// me.instance.events = <-me.instance.emitter.On(topic.String(), middleware...)
	// me.group.Add(me.instance.emitter.On(topic.String()))

	return me.instance.emitter.On(topic.String(), middleware...)
}


func (me *Channels) once(topic messages.MessageTopic, middleware ...func(emitter *emitter.Event)) <-chan emitter.Event {
	eblog.Debug(me.EntityId, "Once")

	// me.instance.events = <-me.instance.emitter.Once(topic.String(), middleware...)
	// me.instance.events.String(1)

	return me.instance.emitter.Once(topic.String(), middleware...)
}


func (me *Channels) use(pattern string, middleware ...func(emitter *emitter.Event))  {
	eblog.Debug(me.EntityId, "Use")

	me.instance.emitter.Use(pattern, middleware...)

	return
}

