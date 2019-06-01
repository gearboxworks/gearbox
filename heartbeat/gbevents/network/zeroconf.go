package network

import (
	"context"
	"errors"
	"fmt"
	"gearbox/box"
	"gearbox/heartbeat/daemon"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/heartbeat/gbevents/tasks"
	"gearbox/only"
	oss "gearbox/os_support"
	"github.com/grandcat/zeroconf"
	"github.com/jinzhu/copier"
	"log"
	"time"
)

func (me *Client) New(OsSupport oss.OsSupporter, args ...Args) error {

	var _args Args
	var err error

	for range only.Once {

		if len(args) > 0 {
			_args = args[0]
		}

		_args.OsSupport = OsSupport
		foo := box.Args{}
		err = copier.Copy(&foo, &_args)
		if err != nil {
			err = errors.New("unable to copy MQTT client config")
			break
		}

		if _args.EntityId == "" {
			_args.EntityId = defaultEntityId
		}

		if _args.Domain == "" {
			_args.Domain = defaultDomain
		}

		if _args.WaitTime == 0 {
			_args.WaitTime = defaultWaitTime
		}

		*me = Client(_args)
		messages.Debug("GBevents - ZeroConf init (%s).", me.EntityId.String())
	}

	if err != nil {
		messages.Debug("Error: ", err)
	}

	// Save last state.
	me.Error = err
	return err
}


func (me *Client) StartHandler() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.Task, err = tasks.StartTask(initZeroConf, startZeroConf, monitorZeroConf, stopZeroConf, me)
		if err != nil {
			break
		}

		time.Sleep(time.Second * 5)

		scanZeroConf := messages.Message{
			Topic: messages.Topic{
				Address: me.EntityId,
				SubTopic: "scan",
			},
			Text: "_gearbox._tcp",
		}
		_ = me.Channels.Publish(scanZeroConf)
		time.Sleep(time.Second * 5)
		_ = me.Channels.Publish(scanZeroConf)
		time.Sleep(time.Second * 5)


		fmt.Printf("Sleeping...\n")
		time.Sleep(time.Hour * 2000)

		port, _ := GetFreePort()
		fmt.Printf("FreePort == %d\n", port)
		me.Register(Service{
			Name: "HELO",
			Type: "_gearbox._tcp",
			Domain: "local",
			Port: port,
		})


		fmt.Printf("Timeout...\n")
		time.Sleep(time.Second * 20)
		fmt.Printf("Stopping...\n")
		me.Task.Stop()


		time.Sleep(time.Hour * 24)

		messages.Debug("started zeroconf handler for %s", me.EntityId.String())
	}

	if err != nil {
		messages.Debug("Error: ", err)
	}

	// Save last state.
	me.Error = err
	return err
}


func (me *Client) StopHandler() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		err = me.Task.Stop()
		if err != nil {
			break
		}

		messages.Debug("stopped zeroconf handler for %s", me.EntityId.String())
	}

	if err != nil {
		messages.Debug("Error: ", err)
	}

	// Save last state.
	me.Error = err
	return err
}



////////////////////////////////////////////////////////////////////////////////
// Executed as a task.
func initZeroConf(task *tasks.Task, i ...interface{}) error {

	var err error

	for range only.Once {
		//me := (i[0]).(*Client)
		if i == nil {
			err = errors.New("software error")
			break
		}
		me := i[0].(*Client)

		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		task.RetryLimit = DefaultRetries
		task.RetryDelay = time.Second * 5

		me.ChannelHandler, err = me.Channels.StartHandler(me.EntityId)
		if err != nil {
			break
		}
		err = me.ChannelHandler.Subscribe(messages.SubTopic("register"), registerService, me)
		if err != nil {
			break
		}
		err = me.ChannelHandler.Subscribe(messages.SubTopic("unregister"), unregisterService, me)
		if err != nil {
			break
		}
		err = me.ChannelHandler.Subscribe(messages.SubTopic("status"), statusService, me)
		if err != nil {
			break
		}
		err = me.ChannelHandler.Subscribe(messages.SubTopic("scan"), scanServices, me)
		if err != nil {
			break
		}

		messages.Debug("ZeroConf %s initialized OK", me.EntityId.String())

		err = nil
	}

	if err != nil {
		messages.Debug("Error: ", err)
	}

	return err
}


// Executed as a task.
func startZeroConf(task *tasks.Task, i ...interface{}) error {

	var err error

	for range only.Once {
		if i == nil {
			err = errors.New("software error")
			break
		}
		me := i[0].(*Client)

		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		// Already started as part of initZeroConf().

		messages.Debug("ZeroConf %s started OK", me.EntityId.String())

		err = nil
	}

	if err != nil {
		messages.Debug("Error: ", err)
	}

	return err
}


// Executed as a task.
func monitorZeroConf(task *tasks.Task, i ...interface{}) error {

	var err error
	var entries ServiceEntries

	for range only.Once {
		if i == nil {
			err = errors.New("software error")
			break
		}
		me := i[0].(*Client)

		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		entries, err = me.Browse("_workstation._tcp")
		if err != nil {
			break
		}

		entries.Print()

		messages.Debug("ZeroConf %s status", me.EntityId.String())

		err = nil
	}

	if err != nil {
		messages.Debug("Error: ", err)
	}

	return err
}


// Executed as a task.
func stopZeroConf(task *tasks.Task, i ...interface{}) error {

	var err error

	for range only.Once {
		if i == nil {
			err = errors.New("software error")
			break
		}
		me := i[0].(*Client)

		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.ChannelHandler.StopHandler()

		messages.Debug("ZeroConf %s stopped", me.EntityId.String())

		err = nil
	}

	if err != nil {
		messages.Debug("Error: ", err)
	}

	return err
}


////////////////////////////////////////////////////////////////////////////////
// Executed from a channel
func scanServices(event *messages.Message, i interface{}) error {

	var err error

	for range only.Once {
		//if i == nil {
		//	err = errors.New("software error")
		//	break
		//}
		me := i.(*Client)

		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		me.Browse(event.Text.String())
		messages.Debug("ZeroConf %s service scan OK", me.EntityId.String())

		err = nil
	}

	if err != nil {
		messages.Debug("Error: ", err)
	}

	return err
}


// Executed from a channel
func statusService(event *messages.Message, i interface{}) error {

	var err error

	for range only.Once {
		if i == nil {
			err = errors.New("software error")
			break
		}
		me := i.(*Client)

		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		messages.Debug("ZeroConf %s service status OK", me.EntityId.String())

		err = nil
	}

	if err != nil {
		messages.Debug("Error: ", err)
	}

	return err
}


// Executed from a channel
func registerService(event *messages.Message, i interface{}) error {

	var err error

	for range only.Once {
		if i == nil {
			err = errors.New("software error")
			break
		}
		me := i.(*Client)

		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		messages.Debug("ZeroConf %s registered service OK", me.EntityId.String())

		err = nil
	}

	if err != nil {
		messages.Debug("Error: ", err)
	}

	return err
}


// Executed from a channel
func unregisterService(event *messages.Message, i interface{}) error {

	var err error

	for range only.Once {
		if i == nil {
			err = errors.New("software error")
			break
		}
		me := i.(*Client)

		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		messages.Debug("ZeroConf %s unregistered service OK", me.EntityId.String())

		err = nil
	}

	if err != nil {
		messages.Debug("Error: ", err)
	}

	return err
}



func (me *Client) Browse(s string) (ServiceEntries, error) {

	var err error
	var results2 ServiceEntries

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		messages.Debug("GBevents - ZeroConf scan started (%s).", me.EntityId.String())
		resolver, err := zeroconf.NewResolver(nil)
		if err != nil {
			err = errors.New("failed to initialize zeroconf resolver")
			break
		}

		entries := make(chan *zeroconf.ServiceEntry)
		go func(results <-chan *zeroconf.ServiceEntry) {
			for entry := range results {
				fmt.Println(entry)
				results2 = append(results2, ServiceEntry(*entry))
			}
			// fmt.Println("No more entries.")
		}(entries)
		fmt.Println("No more entries.")
		fmt.Printf("No more entries: %v\n", entries)

		ctx, cancel := context.WithTimeout(context.Background(), me.WaitTime)
		defer cancel()
		err = resolver.Browse(ctx, s, me.Domain, entries)
		if err != nil {
			err = errors.New("failed to browse zeroconf network")
			break
		}

		<-ctx.Done()

		// fmt.Println(results2)

		messages.Debug("GBevents - ZeroConf scan stopped (%s).", me.EntityId.String())
	}

	return ServiceEntries(results2), err
}


func (me *Client) Register(s Service) error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		server, err := zeroconf.Register(
			s.Name.String(),
			s.Type.String(),
			s.Domain.String(),
			s.Port,
			[]string{"txtv=0", "lo=1", "la=2"},
			nil)
		if err != nil {
			err = errors.New("unable to register zeroconf service")
			break
		}
		defer server.Shutdown()
		log.Println("Published service:")
		log.Println("- Name:", s.Name)
		log.Println("- Type:", s.Type)
		log.Println("- Domain:", s.Domain)
		log.Println("- Port:", s.Port)

		daemon.WaitForTimeout(me.WaitTime)

		log.Println("Shutting down.")
	}

	return err
}

