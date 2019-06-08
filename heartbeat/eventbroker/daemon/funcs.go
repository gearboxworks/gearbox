package daemon

import (
	"errors"
	"fmt"
	"gearbox/global"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/network"
	"gearbox/only"
	"github.com/kardianos/service"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)


func (me *Daemon) EnsureNotNil() error {

	var err error

	switch {
		case me == nil:
			err = errors.New("Daemon instance is nil")
	}

	return err
}
func EnsureNotNil(me *Daemon) error {
	return me.EnsureNotNil()
}


func (me *ServicesMap) EnsureNotNil() error {

	var err error

	switch {
		case me == nil:
			err = errors.New("Daemon ServicesMap instance is nil")
	}

	return err
}
func EnsureServicesMapNotNil(me *ServicesMap) error {
	return me.EnsureNotNil()
}


func (me *Service) EnsureNotNil() error {
	var err error

	switch {
		case me == nil:
			err = errors.New("Daemon Service instance is nil")
		case (me.IsManaged == true) && me.instance.cmd == nil:
			err = me.EntityId.ProduceError("service cmd instance nil")
		case (me.IsManaged == true) && me.instance.exit == nil:
			err = me.EntityId.ProduceError("service exit func is nil")
		case (me.IsManaged == true) && me.instance.service == nil:
			err = me.EntityId.ProduceError("service instance is nil")
	}

	return err
}
func EnsureServiceNotNil(me *Service) error {
	return me.EnsureNotNil()
}


func IsParentInit() (bool) {

	ppid := os.Getppid()
	if ppid == 1 {
		return true
	}

	return false
}


// This function will cause a Go() thread to sit and wait until
// a signal has been sent to the process.
// Very important for tidy up afterwards.
// func WaitForSignal(name string) os.Signal {
func WaitForSignal() os.Signal {

	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)

	signal.Notify(signalChan, os.Kill, os.Interrupt)
	s := <-signalChan
	signal.Stop(signalChan)

	return s
}


// Wait for an ever increasing period of time - a very simple retry back-off system.
// This is used with processes that die too quickly and will ensure that retries don't hammer the system.
func WaitDelay(retry int) {

	// First time wait for 100mS
	// Second time wait for 200mS
	// And so on...
	time.Sleep(time.Millisecond * 100 * time.Duration(retry))
}


// This function will cause a Go() thread to sit and wait until
// a signal has been sent to the process.
// Very important for tidy up afterwards.
func WaitForTimeout(wt time.Duration) bool {

	var exitState bool

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// Timeout timer.
	var tc <-chan time.Time
	if wt > 0 {
		tc = time.After(wt)
	}

	select {
		case <-sig:
			exitState = false
			// Exit by user
		case <-tc:
			exitState = true
			// Exit by timeout
	}

	return exitState
}


// This function will cause a Go() thread to sit and wait until
// a signal has been sent to the process.
// Very important for tidy up afterwards.
func SimpleWaitLoop(t string, i int, d time.Duration) {

	for iterate := 0; iterate < i; i++ {
		fmt.Printf("> Wait: %s\n", t)
		time.Sleep(d)
	}

	return
}


func (me *Daemon) GetId() messages.MessageAddress {

	return me.EntityId
}


func InterfaceToTypeDaemon(i interface{}) (*Daemon, error) {

	var err error
	var me *Daemon

	for range only.Once {
		err = channels.EnsureArgumentNotNil(i)
		if err != nil {
			//break
		}
		me = i.(*Daemon)

		err = me.EnsureNotNil()
		if err != nil {
			break
		}
	}

	return me, err
}


func (me *CreateEntry) Print() error {

	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("software error")
			break
		}

//		sn,_ := me.ServiceName()
//		sin,_ := me.ServiceInstanceName()
//		stn,_ := me.ServiceTypeName()
//
//		//
//		fmt.Printf(` me.Instance = %v
// me.Service = %v
// me.Domain = %v
// me.Port = %v
// me.Text = %v
// me.AddrIPv4 = %v
// me.AddrIPv6 = %v
// me.HostName = %v
// me.TTL = %v
// me.ServiceName() = %v
// me.ServiceInstanceName() = %v
// me.ServiceTypeName() = %v
//`,
//			me.Instance,
//			me.Service,
//			me.Domain,
//			me.Port,
//			me.Text,
//			me.AddrIPv4,
//			me.AddrIPv6,
//			me.HostName,
//			me.TTL,
//			sn,
//			sin,
//			stn,
//		)
	}

	return err
}


func InterfaceToTypeService(i interface{}) (*Service, error) {

	var err error
	var s *Service

	for range only.Once {
		err = channels.EnsureArgumentNotNil(i)
		if err != nil {
			break
		}

		s = i.(*Service)
		// zc = (i[0]).(*Service)
		// zc = i[0].(*Service)

		err = s.EnsureNotNil()
		if err != nil {
			break
		}
	}

	return s, err
}


func (me *ServicesMap) Print() error {

	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		for u, s := range *me {
			fmt.Printf("# Entry: %s\n", u)
			err = s.Print()
			if err != nil {
				break
			}
		}
	}

	return err
}


func (me *Service) Print() error {

	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("software error")
			break
		}

		//if (me.instance == nil) && (me.IsManaged == true) {
		//	fmt.Printf("# Entry(deleted): %v", me.EntityId)
		//} else {
			fmt.Printf("# Entry: %v", me.EntityId)
		//}

		//err = me.Entry.Print()
		if err != nil {
			break
		}
	}

	return err
}


func (me *CreateEntry) IsTheSame(e CreateEntry) (bool, error) {

	var same bool
	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("software error")
			same = false
			break
		}

		//if (me.Instance == e.Instance) &&
		//	(trimDot(me.Service) == trimDot(e.Service)) &&
		//	(trimDot(me.Domain) == trimDot(e.Domain)) &&
		//	(me.Port == e.Port) {
		//	same = true
		//}
	}

	return same, err
}


func (me *programInstance) Start(s service.Service) error {
	panic("implement me")
}

func (me *programInstance) Stop(s service.Service) error {
	panic("implement me")
}


// Yup
func (me *Service) CreateMdnsEntry() (*network.CreateEntry, error) {

	var err error
	var zc network.CreateEntry

	for range only.Once {

		if me.Entry.MdnsType == "" {
			err = me.EntityId.ProduceError("MdnsType not set")
			break
		}

		zc = network.CreateEntry{
			Name:   network.Name(strings.ToLower("_" + global.Brandname + "-" + me.Entry.Name)),
			Type:   network.Type(fmt.Sprintf("_%s._tcp", me.Entry.MdnsType)),
			Domain: network.DefaultDomain,
			Port:   me.Entry.Port,
		}
	}

	return &zc, err
}
