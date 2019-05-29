package gbChannels

import (
	"fmt"
	"gearbox/box"
	"gearbox/global"
	"gearbox/heartbeat/daemon"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/help"
	"gearbox/only"
	oss "gearbox/os_support"
	"github.com/fhmq/hmq/broker"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/jinzhu/copier"
	"github.com/olebedev/emitter"
	"log"
	"os"
	"time"
)


func (me *Channels) New(OsSupport oss.OsSupporter, args ...Args) status.Status {

	var _args Args
	var sts status.Status

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		_args.OsSupport = OsSupport
		foo := box.Args{}
		err := copier.Copy(&foo, &_args)
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("unable to copy Go channels config").
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		if _args.EntityId == "" {
			_args.EntityId = "gearbox-channels"
		}

		_args.emitter = emitter.Emitter{}

		*me = Channels(_args)

		go func() {
			me.TestWaitForExit()
		}()

		me.TestSendAnExit()
	}

	return sts
}


func (me *Channels) TestWaitForExit() status.Status {

	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		log.Println("Starting.")
		for event := range me.emitter.On("/channel/exit") {
			if event.Args == nil {
				continue
			}

			foo := event.Args[0].(messages.Message)
			fmt.Printf("%s -> %s (%d): Event2(%s)\n", foo.Src, foo.Text, foo.Time.Unix(), event.OriginalTopic)
		}
		log.Println("Stopping.")
	}

	return sts
}


func (me *Channels) TestSendAnExit() status.Status {

	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		log.Println("Waiting .")
		time.Sleep(time.Second * 5)
		me.Emit("/channel/exit", "HELLO")

		log.Println("Shutting down.")
	}

	return sts
}


func (me *Channels) Start() status.Status {

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

			foo := event.Args[0].(messages.Message)
			fmt.Printf("%s -> %s (%d): Event2(%s)\n", foo.Src, foo.Text, foo.Time.Unix(), event.OriginalTopic)
		}
		fmt.Printf("GBevents - Poller(FINISHED)\n")

		sts = status.Success("%s GBevents - Poller exited.", global.Brandname)
	}
	me.Sts = sts
	status.Log(sts)

	return sts
}


func (me *Channels) Emit(topic string, args ...string)  {
	fmt.Printf("Emit\n")

	me.emits = me.emitter.Emit(topic, args)

	return
}


func (me *Channels) Listeners(topic string)  {
	fmt.Printf("Listeners\n")

	foo := me.emitter.Listeners(topic)[0]

	fmt.Printf("%v", foo)
	for f := range foo {
		fmt.Printf("[%s] - '%s' '%s' '%s'\n", f.Topic, f.OriginalTopic, f.Args, f.Flags)
	}

	return
}


func (me *Channels) Off(topic string, channels ...<-chan emitter.Event)  {
	fmt.Printf("Off\n")

	me.emitter.Off(topic, channels...)

	return
}


func (me *Channels) On(topic string, middleware ...func(emitter *emitter.Event))  {
	fmt.Printf("On\n")

	me.events = <-me.emitter.On(topic, middleware...)
	me.group.Add(me.emitter.On(topic))

	return
}


func (me *Channels) Once(topic string, middleware ...func(emitter *emitter.Event)) {
	fmt.Printf("Once\n")

	me.events = <-me.emitter.Once(topic, middleware...)
	// me.events.String(1)

	return
}


func (me *Channels) Topics() (topics []string) {
	fmt.Printf("Topics\n")

	topics = me.emitter.Topics()

	return
}


func (me *Channels) Use(pattern string, middleware ...func(emitter *emitter.Event))  {
	fmt.Printf("Use\n")

	me.emitter.Use(pattern, middleware...)

	return
}


func (me *Channels) Test() {
	/*
		argh := ServiceEvents{
			Name: "mqtt",
			State: "running",
			Action: ServiceAction{
				State: "running",
				CallBack: hello,
			},
		}

		argh.CreateService()
		foo := ServiceEvents{}
		foo.State.String()
	*/


	fmt.Printf("DEBUG - 1\n")
	e1 := &emitter.Emitter{}
	e2 := &emitter.Emitter{}
	e3 := &emitter.Emitter{}


	go func() {
		fmt.Printf("poller1()\n")
		poller1(e1, e2, e3)
	}()

	go func() {
		fmt.Printf("poller2()\n")
		poller2(e1)
	}()

	go func() {
		fmt.Printf("poller3()\n")
		poller3(e2)
	}()


	go func() {
		time.Sleep(time.Millisecond * 1500)

		e1.Emit("change", messages.Message{Src: "1", Time: time.Now(), Text: "A"})
		time.Sleep(time.Millisecond * 100)

		e1.Emit("nope", messages.Message{Src: "1", Time: time.Now(), Text: "B"})
		time.Sleep(time.Millisecond * 100)

		e1.Emit("change", messages.Message{Src: "1", Time: time.Now(), Text: "C"})
		time.Sleep(time.Millisecond * 100)

		e1.Emit("change", messages.Message{Src: "1", Time: time.Now(), Text: "D"})

		time.Sleep(time.Second * 10)
		e1.Off("*") // unsubscribe any listeners
	}()


	fmt.Printf("DEBUG - 2\n")
	go func() {
		time.Sleep(time.Second * 2)

		e1.Emit("change", messages.Message{Src: "2", Time: time.Now(), Text: "A"})
		time.Sleep(time.Millisecond * 100)

		e1.Emit("*", messages.Message{Text: "B", Time: time.Now(), Src: "2"})
		time.Sleep(time.Millisecond * 100)

		e1.Emit("change", messages.Message{Src: "2", Time: time.Now(), Text: "C"})

		time.Sleep(time.Second * 10)
		e2.Off("*") // unsubscribe any listeners
	}()


	fmt.Printf("DEBUG - 3\n")
	go func() {
		time.Sleep(time.Millisecond * 1100)

		e3.Emit("change", messages.Message{Src: "3", Time: time.Now(), Text: "A"})
		time.Sleep(time.Millisecond * 100)

		e3.Emit("*", messages.Message{Src: "3", Time: time.Now(), Text: "B"})
		time.Sleep(time.Millisecond * 100)

		e3.Emit("change", messages.Message{Src: "3", Time: time.Now(), Text: "C"})

		time.Sleep(time.Second * 10)
		e3.Off("*") // unsubscribe any listeners
	}()


	// listener channel was closed
	fmt.Printf("DEBUG - 4\n")

	time.Sleep(time.Hour)

	fmt.Printf("DEBUG - ENTRY\n")
	config, err := broker.ConfigureConfig(os.Args[1:])
	if err != nil {
		log.Fatal("configure broker config error: ", err)
	}

	b, err := broker.NewBroker(config)
	if err != nil {
		log.Fatal("New Broker error: ", err)
	}
	b.Start()
	fmt.Printf("DEBUG - START\n")

	s := daemon.WaitForSignal()
	log.Println("signal received, broker closed.", s)

}


func poller1(e1 *emitter.Emitter, e2 *emitter.Emitter, e3 *emitter.Emitter) {
	g := &emitter.Group{Cap: 1}
	g.Add(e1.On("*"), e2.On("*"), e3.On("*"))
	for event := range g.On() {
		if event.Args == nil {
			continue
		}

		foo := event.Args[0].(messages.Message)
		fmt.Printf("%s -> %s (%d): Event1(%s)\n", foo.Src, foo.Text, foo.Time.Unix(), event.OriginalTopic)
	}
	fmt.Printf("Event1(FINISHED)\n")
}

func poller2(e *emitter.Emitter) {
	for event := range e.On("*") {
		if event.Args == nil {
			continue
		}

		foo := event.Args[0].(messages.Message)
		fmt.Printf("%s -> %s (%d): Event2(%s)\n", foo.Src, foo.Text, foo.Time.Unix(), event.OriginalTopic)
	}
	fmt.Printf("Event2(FINISHED)\n")
}

func poller3(e *emitter.Emitter) {
	for event := range e.On("*") {
		if event.Args == nil {
			continue
		}

		foo := event.Args[0].(messages.Message)
		fmt.Printf("%s -> %s (%d): Event3(%s)\n", foo.Src, foo.Text, foo.Time.Unix(), event.OriginalTopic)
		//		go func() {
		//e.Emit("change", "...")
		//		}()
	}
	fmt.Printf("Event3(FINISHED)\n")
}


func hello() {
	fmt.Printf("CALLBACK\n")
}
