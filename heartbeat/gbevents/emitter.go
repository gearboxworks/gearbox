package gbevents

import (
	"fmt"
	"gearbox/global"
	"gearbox/only"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/olebedev/emitter"
)


func (me *ServiceEvents) startPoller() status.Status {

	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		fmt.Printf("GBevents - Poller(STARTED)\n")
		for event := range me.emitter.On("*") {
			if event.Args == nil {
				continue
			}

			foo := event.Args[0].(Message)
			fmt.Printf("%s -> %s (%d): Event2(%s)\n", foo.Src, foo.Text, foo.Time.Unix(), event.OriginalTopic)
		}
		fmt.Printf("GBevents - Poller(FINISHED)\n")

		sts = status.Success("%s GBevents - Poller exited.", global.Brandname)
	}

	return sts
}


func (me *ServiceEvents) Emit(topic string, args ...string)  {
	fmt.Printf("Emit\n")

	me.emits = me.emitter.Emit(topic, args)

	return
}


func (me *ServiceEvents) Listeners(topic string)  {
	fmt.Printf("Listeners\n")

	foo := me.emitter.Listeners(topic)[0]

	fmt.Printf("%v", foo)
	for f := range foo {
		fmt.Printf("[%s] - '%s' '%s' '%s'\n", f.Topic, f.OriginalTopic, f.Args, f.Flags)
	}

	return
}


func (me *ServiceEvents) Off(topic string, channels ...<-chan emitter.Event)  {
	fmt.Printf("Off\n")

	me.emitter.Off(topic, channels...)

	return
}


func (me *ServiceEvents) On(topic string, middleware ...func(emitter *emitter.Event))  {
	fmt.Printf("On\n")

	me.events = <-me.emitter.On(topic, middleware...)
	me.group.Add(me.emitter.On(topic))

	return
}


func (me *ServiceEvents) Once(topic string, middleware ...func(emitter *emitter.Event)) {
	fmt.Printf("Once\n")

	me.events = <-me.emitter.Once(topic, middleware...)
	// me.events.String(1)

	return
}


func (me *ServiceEvents) Topics() (topics []string) {
	fmt.Printf("Topics\n")

	topics = me.emitter.Topics()

	return
}


func (me *ServiceEvents) Use(pattern string, middleware ...func(emitter *emitter.Event))  {
	fmt.Printf("Use\n")

	me.emitter.Use(pattern, middleware...)

	return
}
