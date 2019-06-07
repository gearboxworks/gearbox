package channels

import (
	"fmt"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"github.com/getlantern/errors"
	"github.com/olebedev/emitter"
)

func (me *Channels) EnsureNotNil() error {

	var err error
	//var emptyChannelInstance channelsInstance

	switch {
		case me == nil:
			err = errors.New("channels instance is nil")
			fmt.Printf("FO\n")

		//case me.instance == emptyChannelInstance:
		//	err = errors.New("Funexpected software error")
		//	fmt.Printf("FO\n")
	}

	return err
}

func EnsureNotNil(me *Channels) error {

	var err error

	switch {
		case me == nil:
			err = errors.New("channels instance is nil")
		//case me.instance.emitter.Cap == nil:
		//	err = errors.New("channels instance is nil")
	}

	return err
}


func (me *Channels) off(topic messages.MessageTopic, channels ...<-chan emitter.Event)  {
	eblog.Debug("Off")

	me.instance.emitter.Off(topic.String(), channels...)

	return
}


func (me *Channels) on(topic messages.MessageTopic, middleware ...func(emitter *emitter.Event)) <-chan emitter.Event {
	eblog.Debug("On")

	// me.instance.events = <-me.instance.emitter.On(topic.String(), middleware...)
	// me.group.Add(me.instance.emitter.On(topic.String()))

	return me.instance.emitter.On(topic.String(), middleware...)
}


func (me *Channels) once(topic messages.MessageTopic, middleware ...func(emitter *emitter.Event)) <-chan emitter.Event {
	eblog.Debug("Once")

	// me.instance.events = <-me.instance.emitter.Once(topic.String(), middleware...)
	// me.instance.events.String(1)

	return me.instance.emitter.Once(topic.String(), middleware...)
}


func (me *Channels) use(pattern string, middleware ...func(emitter *emitter.Event))  {
	eblog.Debug("Use")

	me.instance.emitter.Use(pattern, middleware...)

	return
}

