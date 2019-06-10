package network

import (
	"context"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/only"
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

		eblog.Debug(me.EntityId, "service scan started")
		resolver, err := zeroconf.NewResolver(nil)
		if err != nil {
			err = me.EntityId.ProduceError("failed to initialize scan resolver")
			break
		}

		// fmt.Printf("Looking for: service:%s domain:%s\n", s, d)
		entries := make(chan *zeroconf.ServiceEntry)
		go func(results <-chan *zeroconf.ServiceEntry) {
			for entry := range results {
				u := messages.GenerateAddress()
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
			err = me.EntityId.ProduceError("failed to scan network")
			break
		}

		<-ctx.Done()

		eblog.Debug(me.EntityId, "service scan completed")
	}

	channels.PublishCallerState(me.Channels, &me.EntityId, &me.State)
	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return found, err
}

// Returns a *Service reference based on a given CreateEntry structure.
func (me *ZeroConf) FindService(e ServiceConfig) (*Service, error) {

	var err error
	var sc *Service

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		// First find locally defined services.
		sc, err = me.services.findServiceInMap(e)
		if err != nil {
			break
		}
		if sc != nil {
			eblog.Debug(me.EntityId, "found managed service %s", sc.EntityId.String())
			break
		}

		// Then look on the network.
		sc, err = me.scannedServices.findServiceInMap(e)
		if err != nil {
			break
		}
		if sc != nil {
			eblog.Debug(me.EntityId, "found network service %s", sc.EntityId.String())
			break
		}

		eblog.Debug(me.EntityId, "no services found on network")
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return sc, err
}

// Returns a *Service reference based on a given CreateEntry structure.
func (me ServicesMap) findServiceInMap(e ServiceConfig) (*Service, error) {

	var err error
	var ok bool
	var sc *Service

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for u := range me {
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
func (me *Service) compareService(e ServiceConfig) (bool, error) {

	var err error
	var found bool

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		switch {
			// Search for exact service definition.
			case (me.Entry.Instance == e.Name.String()) &&
				(me.Entry.Service == e.Type.String()) &&
				(me.Entry.Domain == e.Domain.String()) &&
				(me.Entry.Port == int(e.Port)):
				found = true
				break

			// Search just by name without port.
			case (me.Entry.Instance == e.Name.String()) &&
				(me.Entry.Service == e.Type.String()) &&
				(me.Entry.Domain == e.Domain.String()):
				found = true
				break
		}

		eblog.Debug(me.EntityId, "matched service %s to %s", me.EntityId.String(), e.EntityId.String())
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
				eblog.Debug(me.EntityId, "deleting entry %s", u.String())
				delete(me.services, u)
				deleted++
				continue
			}
			if (me.services[u].instance == nil) && (me.services[u].IsManaged == true) {
				// If we are managing this service locally and the instance has been removed,
				// then delete.
				eblog.Debug(me.EntityId, "deleting entry %s", u.String())
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

		eblog.Debug(me.EntityId, "updated registered services - added:%d, deleted:%d.", added, deleted)
	}

	eblog.LogIfNil(me, err)
	eblog.LogIfError(me.EntityId, err)

	return err
}

