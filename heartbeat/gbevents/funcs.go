package gbevents

import (
	"fmt"
	"gearbox/box"
	"gearbox/global"
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
	"os/signal"
	"path/filepath"
	"time"
)


func New(OsSupport oss.OsSupporter, args ...Args) (*ServiceEvents, status.Status) {

	var _args Args
	var sts status.Status
	se := &ServiceEvents{}

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		_args.OsSupport = OsSupport
		foo := box.Args{}
		err := copier.Copy(&foo, &_args)
		if err != nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("%s GBevents - Struct copy failed.", global.Brandname),
				Help:    help.ContactSupportHelp(), // @TODO need better support here
				Data:    err,
			})
			break
		}

		if _args.Boxname == "" {
			_args.Boxname = global.Brandname
		}

		//execPath, _ := os.Executable()
		//execCwd := string(_args.OsSupport.GetAdminRootDir()) + "/heartbeat" // os.Getwd()

		_args.PidFile = filepath.FromSlash(fmt.Sprintf("%s/%s", _args.OsSupport.GetAdminRootDir(), defaultPidFile))

		*se = ServiceEvents(_args)
	}

	return se, sts
}


func (me *ServiceEvents) StartEventServer() status.Status {

	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		go func() {
			me.StsMqtt = me.startMqttBroker()
		}()

		go func() {
			me.StsEmitter = me.startPoller()
		}()

		go func() {
			me.StsEmitter = me.startMqttClient()
		}()

		sts = status.Success("%s GBevents - Server started.", global.Brandname)
	}

	return sts
}


func (me *ServiceEvents) RegisterService(topic string, args ...ServiceData) {
	fmt.Printf("RegisterService\n")

	// .

	return
}


func (me *ServiceEvents) Create() status.Status {
	fmt.Printf("(me *ServiceEvents) CreateService() status.Status\n")

	return nil
}

func (me *ServiceEvents) Start() status.Status {
	fmt.Printf("(me *ServiceEvents) StartService() status.Status\n")

	return nil
}


func (me *ServiceEvents) Stop() status.Status {
	fmt.Printf("(me *ServiceEvents) StopService() status.Status\n")

	return nil
}


func (me *ServiceEvents) Restart() status.Status {
	fmt.Printf("(me *ServiceEvents) RestartService() status.Status\n")

	return nil
}


func (me *ServiceEvents) Status() status.Status {
	fmt.Printf("(me *ServiceEvents) ServiceStatus() status.Status\n")

	return nil
}


func (me *ServiceEvents) Test() {
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

		e1.Emit("change", Message{Src: "1", Time: time.Now(), Text: "A"})
		time.Sleep(time.Millisecond * 100)

		e1.Emit("nope", Message{Src: "1", Time: time.Now(), Text: "B"})
		time.Sleep(time.Millisecond * 100)

		e1.Emit("change", Message{Src: "1", Time: time.Now(), Text: "C"})
		time.Sleep(time.Millisecond * 100)

		e1.Emit("change", Message{Src: "1", Time: time.Now(), Text: "D"})

		time.Sleep(time.Second * 10)
		e1.Off("*") // unsubscribe any listeners
	}()


	fmt.Printf("DEBUG - 2\n")
	go func() {
		time.Sleep(time.Second * 2)

		e1.Emit("change", Message{Src: "2", Time: time.Now(), Text: "A"})
		time.Sleep(time.Millisecond * 100)

		e1.Emit("*", Message{Text: "B", Time: time.Now(), Src: "2"})
		time.Sleep(time.Millisecond * 100)

		e1.Emit("change", Message{Src: "2", Time: time.Now(), Text: "C"})

		time.Sleep(time.Second * 10)
		e2.Off("*") // unsubscribe any listeners
	}()


	fmt.Printf("DEBUG - 3\n")
	go func() {
		time.Sleep(time.Millisecond * 1100)

		e3.Emit("change", Message{Src: "3", Time: time.Now(), Text: "A"})
		time.Sleep(time.Millisecond * 100)

		e3.Emit("*", Message{Src: "3", Time: time.Now(), Text: "B"})
		time.Sleep(time.Millisecond * 100)

		e3.Emit("change", Message{Src: "3", Time: time.Now(), Text: "C"})

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

	s := waitForSignal()
	log.Println("signal received, broker closed.", s)

}


func poller1(e1 *emitter.Emitter, e2 *emitter.Emitter, e3 *emitter.Emitter) {
	g := &emitter.Group{Cap: 1}
	g.Add(e1.On("*"), e2.On("*"), e3.On("*"))
	for event := range g.On() {
		if event.Args == nil {
			continue
		}

		foo := event.Args[0].(Message)
		fmt.Printf("%s -> %s (%d): Event1(%s)\n", foo.Src, foo.Text, foo.Time.Unix(), event.OriginalTopic)
	}
	fmt.Printf("Event1(FINISHED)\n")
}

func poller2(e *emitter.Emitter) {
	for event := range e.On("*") {
		if event.Args == nil {
			continue
		}

		foo := event.Args[0].(Message)
		fmt.Printf("%s -> %s (%d): Event2(%s)\n", foo.Src, foo.Text, foo.Time.Unix(), event.OriginalTopic)
	}
	fmt.Printf("Event2(FINISHED)\n")
}

func poller3(e *emitter.Emitter) {
	for event := range e.On("*") {
		if event.Args == nil {
			continue
		}

		foo := event.Args[0].(Message)
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


func EnsureNotNil(bx *ServiceEvents) (sts status.Status) {
	if bx == nil {
		sts = status.Fail(&status.Args{
			Message: "unexpected error",
			Help:    help.ContactSupportHelp(), // @TODO need better support here
			Data:    unknownState,
		})
	}
	return sts
}


func waitForSignal() os.Signal {
	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)
	signal.Notify(signalChan, os.Kill, os.Interrupt)
	s := <-signalChan
	signal.Stop(signalChan)
	return s
}

