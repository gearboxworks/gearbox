package network

import (
	"context"
	"errors"
	"fmt"
	"gearbox/app/logger"
	"gearbox/heartbeat/gbevents/channels"
	"gearbox/heartbeat/gbevents/messages"
	"gearbox/only"
	"github.com/google/uuid"
	"github.com/grandcat/zeroconf"
)


////////////////////////////////////////////////////////////////////////////////
// Executed as a method.

// Browses the M-DNS broadcast network for registered services.
func (me *ZeroConf) Browse(s string, d string) (ServicesMap, error) {

	var err error
	found := make(ServicesMap)

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		logger.Debug("GBevents - ZeroConf scan started (%s).", me.EntityId.String())
		resolver, err := zeroconf.NewResolver(nil)
		if err != nil {
			err = errors.New("failed to initialize zeroconf resolver")
			break
		}

		// fmt.Printf("Looking for: service:%s domain:%s\n", s, d)
		entries := make(chan *zeroconf.ServiceEntry)
		go func(results <-chan *zeroconf.ServiceEntry) {
			for entry := range results {
				u := uuid.New()
				//fmt.Printf("Found: %v\n", *entry)
				found[u] = &Service{
						EntityId: u,
						Entry: Entry(*entry),
				}
			}
			// fmt.Println("No more entries.")
		}(entries)

		ctx, cancel := context.WithTimeout(context.Background(), me.waitTime)
		defer cancel()
		err = resolver.Browse(ctx, s, d, entries)
		if err != nil {
			err = errors.New("failed to browse zeroconf network")
			break
		}

		<-ctx.Done()

		logger.Debug("GBevents - ZeroConf scan finished (%s).", me.EntityId.String())
	}

	if err != nil {
		logger.Debug("Error: %v", err)
	}

	return found, err
}

// Returns a *Service reference based on a given CreateEntry structure.
func (me *ZeroConf) FindService(e CreateEntry) (*Service, error) {

	var err error
	var u uuid.UUID
	var sc *Service

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		// First find locally defined services.
		fmt.Print("Checking me.services.findServiceInMap\n")
		sc, err = me.services.findServiceInMap(e)
		if err != nil {
			break
		}
		if sc != nil {
			logger.Debug("ZeroConf %s found managed service %s OK", me.EntityId.String(), u.String())
			break
		}

		// Then look on the network.
		fmt.Print("Checking me.scannedServices.findServiceInMap\n")
		sc, err = me.scannedServices.findServiceInMap(e)
		if err != nil {
			break
		}
		if sc != nil {
			logger.Debug("ZeroConf %s found network service %s OK", me.EntityId.String(), u.String())
			break
		}

		logger.Debug("ZeroConf no service found")
	}

	return sc, err
}

// Returns a *Service reference based on a given CreateEntry structure.
func (me ServicesMap) findServiceInMap(e CreateEntry) (*Service, error) {

	var err error
	var ok bool
	var u uuid.UUID
	var sc *Service

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		fmt.Printf("findServiceInMap: %v\n", e)
		fmt.Printf("findServiceInMap: %v\n", me)
		for u = range me {
			//me[u].Print()

			ok, err = me[u].compareService(e)
			if err != nil {
				break
			}

			if ok == true {
				sc = me[u]
				break
			}
		}
	}

	return sc, err
}

// Returns a *Service reference based on a given CreateEntry structure.
func (me *Service) compareService(e CreateEntry) (bool, error) {

	var err error
	var found bool

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		fmt.Printf("compareService: %v\n", e)
		switch {
			// Search for exact service definition.
			case (me.Entry.Instance == e.Name.String()) &&
				(me.Entry.Service == e.Type.String()) &&
				(me.Entry.Domain == e.Domain.String()) &&
				(me.Entry.Port == e.Port):
				found = true
				fmt.Printf("ZeroConf matched service %s OK\n", me.EntityId.String())
				break

			// Search just by name without port.
			case (me.Entry.Instance == e.Name.String()) &&
				(me.Entry.Service == e.Type.String()) &&
				(me.Entry.Domain == e.Domain.String()):
				found = true
				fmt.Printf("ZeroConf matched service %s OK\n", me.EntityId.String())
				break
		}

		logger.Debug("ZeroConf matched service %s OK", me.EntityId.String())
	}

	return found, err
}

// Non-exposed function allowing registered *Services to be updated
// with additional information, (such as IP addresses & TTL).
func (me *ZeroConf) updateRegisteredServices() error {

	var err error
	var ok bool
	var added, deleted int

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for u, _ := range me.services {
			if _, ok = me.services[u]; !ok {
				// Shouldn't ever see this, but hey, might as well be anal about it.
				logger.Debug("Deleting zeroconf entry %v.", u)
				delete(me.services, u)
				deleted++
				continue
			}
			if (me.services[u].instance == nil) && (me.services[u].IsManaged == true) {
				// If we are managing this service locally and the instance has been removed,
				// then delete.
				logger.Debug("Deleting zeroconf entry %v.", u)
				delete(me.services, u)
				deleted++
				continue
			}

			for su, c := range me.scannedServices {
				ok, err = me.services[u].Entry.UpdateService(c.Entry)
				if ok {
					// Remove from scannedServices to avoid checking on next iteration.
					added++
					delete(me.scannedServices, su)

					//scanResults[i] = scanResults[len(scanResults)-1]	// Copy last element to index i.
					//scanResults[len(scanResults)-1] = Entry{}			// Erase last element (write zero value).
					//scanResults = scanResults[:len(scanResults)-1]		// Truncate slice.
				}
			}
		}

		logger.Debug("ZeroConf updated registered services - added:%d, deleted:%d.", added, deleted)
	}

	return err
}


////////////////////////////////////////////////////////////////////////////////
// Executed from a channel

// Non-exposed channel function that responds to a "scan" channel request.
func scanServices(event *messages.Message, i channels.Argument) channels.Return {

	var me *ZeroConf
	var err error

	for range only.Once {
		me, err = InterfaceToTypeZeroConf(i)
		if err != nil {
			break
		}

		_, err = me.Browse(event.Text.String(), me.domain)
		if err != nil {
			break
		}

		logger.Debug("ZeroConf %s service scan OK", me.EntityId.String())
		err = nil
	}

	if err != nil {
		logger.Debug("Error: ", err)
	}

	return err
}
