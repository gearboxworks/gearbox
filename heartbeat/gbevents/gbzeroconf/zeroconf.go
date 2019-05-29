package gbZeroConf

import (
	"context"
	"fmt"
	"gearbox/box"
	"gearbox/help"
	"gearbox/only"
	oss "gearbox/os_support"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/grandcat/zeroconf"
	"github.com/jinzhu/copier"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (me *Client) New(OsSupport oss.OsSupporter, args ...Args) status.Status {

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
				SetMessage("unable to copy MQTT client config").
				SetAdditional("", ).
				SetData("").
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
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

		status.Success("GBevents - zeroconf(INIT)").Log()

		*me = Client(_args)
	}

	return sts
}


func (me *Client) Browse(s string) status.Status {

	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
			break
		}

		fmt.Printf("GBevents - ZeroConf(STARTED)\n")
		status.Success("GBevents - ZeroConf(STARTED)").Log()
		resolver, err := zeroconf.NewResolver(nil)
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("failed to initialize zeroconf resolver").
				SetAdditional("").
				SetData(s).
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		entries := make(chan *zeroconf.ServiceEntry)
		go func(results <-chan *zeroconf.ServiceEntry) {
			for entry := range results {
				log.Println(entry)
			}
			log.Println("No more entries.")
		}(entries)

		ctx, cancel := context.WithTimeout(context.Background(), me.WaitTime)
		defer cancel()
		err = resolver.Browse(ctx, s, me.Domain, entries)
		if err != nil {
			sts = status.Wrap(err).
				SetMessage("failed to browse zeroconf network").
				SetAdditional("").
				SetData(s).
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}

		<-ctx.Done()
		fmt.Printf("GBevents - ZeroConf(STOPPED)\n")
		status.Success("GBevents - ZeroConf(STOPPED)").Log()
	}

	return sts
}


func (me *Client) Register(s Service) status.Status {

	var sts status.Status

	for range only.Once {
		sts = EnsureNotNil(me)
		if is.Error(sts) {
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
			sts = status.Wrap(err).
				SetMessage("unable to register zeroconf service").
				SetAdditional("").
				SetData(s).
				SetCause(err).
				SetHelp(status.AllHelp, help.ContactSupportHelp())
			break
		}
		defer server.Shutdown()
		log.Println("Published service:")
		log.Println("- Name:", s.Name)
		log.Println("- Type:", s.Type)
		log.Println("- Domain:", s.Domain)
		log.Println("- Port:", s.Port)

		// Clean exit.
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		// Timeout timer.
		var tc <-chan time.Time
		if me.WaitTime > 0 {
			tc = time.After(me.WaitTime)
		}

		select {
			case <-sig:
				// Exit by user
			case <-tc:
				// Exit by timeout
		}

		log.Println("Shutting down.")
	}

	return sts
}

