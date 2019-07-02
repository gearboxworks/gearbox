package daemon

import (
	"errors"
	"fmt"
	"gearbox/eventbroker/msgs"
	"gearbox/eventbroker/network"
	"gearbox/eventbroker/states"
	"github.com/gearboxworks/go-status/only"
	"github.com/kardianos/service"
	"net/url"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"
)

func (me *Daemon) EnsureNotNil() error {

	var err error

	switch {
	case me == nil:
		err = errors.New("daemon instance is nil")
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
		err = errors.New("daemon servicesmap instance is nil")
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
		err = errors.New("daemon service instance is nil")
	case (me.IsManaged == true) && me.instance.cmd == nil:
		err = msgs.MakeError(me.EntityId, "service cmd instance nil")
	case (me.IsManaged == true) && me.instance.exit == nil:
		err = msgs.MakeError(me.EntityId, "service exit func is nil")
	case (me.IsManaged == true) && me.instance.service == nil:
		err = msgs.MakeError(me.EntityId, "service instance is nil")
	}

	return err
}
func EnsureServiceNotNil(me *Service) error {
	return me.EnsureNotNil()
}

func IsParentInit() bool {

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
	case <-tc: // @TODO For `<-tc` GoLand reports "Receiver may block because of a nil channel"
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

func (me *Daemon) GetId() msgs.Address {

	return me.EntityId
}

func InterfaceToTypeDaemon(i interface{}) (*Daemon, error) {

	var err error
	var me *Daemon

	for range only.Once {
		if i == nil {
			err = errors.New("interface is nil, should be" + InterfaceTypeDaemon)
			break
		}

		checkType := reflect.ValueOf(i)
		//fmt.Printf("InterfaceToTypeDaemon = %v\n", checkType.Type().String())
		if checkType.Type().String() != InterfaceTypeDaemon {
			err = errors.New("interface type not " + InterfaceTypeDaemon)
			break
		}

		me = i.(*Daemon)

		err = me.EnsureNotNil()
		if err != nil {
			break
		}
	}

	return me, err
}

func InterfaceToTypeService(i interface{}) (*Service, error) {

	var err error
	var s *Service

	for range only.Once {
		if i == nil {
			err = errors.New("interface is nil, should be" + InterfaceTypeService)
			break
		}

		checkType := reflect.ValueOf(i)
		if checkType.Type().String() != InterfaceTypeService {
			err = errors.New("interface type not " + InterfaceTypeService)
			break
		}

		s = i.(*Service)

		err = s.EnsureNotNil()
		if err != nil {
			break
		}
	}

	return s, err
}

func (me *ServiceConfig) Print() error {

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
		fmt.Printf("# Entry: %v", me.EntityId)
	}

	return err
}

func (me *programInstance) Start(s service.Service) error {
	panic("implement me")
}

func (me *programInstance) Stop(s service.Service) error {
	panic("implement me")
}

func (me *ServiceConfig) ToServiceType() *service.Config {

	return &me.Config
}

// Ensure we don't duplicate services.
func (me *Service) IsExisting(him ServiceConfig) error {

	var err error

	if (me.Entry.UrlPtr == nil) || (him.UrlPtr == nil) {
		// return errors.New("nil pointer")
		return err
	}

	switch {
	case me.Entry.Config.Name == him.Config.Name:
		err = msgs.MakeError(me.EntityId, "daemon service Name:%s already exists", me.Entry.Config.Name)

	case me.Entry.Config.DisplayName == him.Config.DisplayName:
		err = msgs.MakeError(me.EntityId, "daemon service DisplayName:%s already exists", me.Entry.DisplayName)

	case me.Entry.Config.Executable == him.Config.Executable:
		err = msgs.MakeError(me.EntityId, "daemon service Executable:%s already exists", me.Entry.Config.Executable)

	case me.Entry.UrlPtr == him.UrlPtr:
		err = msgs.MakeError(me.EntityId, "daemon service Url:%s already exists", me.Entry.UrlPtr)

	case me.Entry.Url == him.Url:
		err = msgs.MakeError(me.EntityId, "daemon service Url:%s already exists", me.Entry.UrlPtr)

	case (me.Entry.UrlPtr.Hostname() == him.UrlPtr.Hostname()) && (me.Entry.UrlPtr.Port() == him.UrlPtr.Port()):
		err = msgs.MakeError(me.EntityId, "daemon service Host:%s:%s already exists", me.Entry.UrlPtr.Hostname(), me.Entry.UrlPtr.Port())
	}

	return err
}

func (me *Daemon) HasFileChanged(fn string) (changed bool, err error) {

	var info os.FileInfo

	for range only.Once {
		jc := me.GetServiceFiles()
		if _, ok := jc[fn]; !ok {
			break
		}

		info, err = os.Stat(fn)
		if err != nil {
			break
		}

		if jc[fn] != info.ModTime() {
			jc[fn] = info.ModTime()
			changed = true
		}
	}

	return changed, err
}

func (me *Daemon) IsFileRegistered(fn string) (ok bool) {

	jc := me.GetServiceFiles()
	_, ok = jc[fn]

	return ok
}

func (me *Daemon) DoesFileExist(fn string) (exists bool, err error) {

	for range only.Once {
		_, err = os.Stat(fn)
		if err != nil {
			break
		}

		exists = true
	}

	return exists, err
}

func (me *Daemon) RemoveFileIfExist(fn string) (err error) {

	for range only.Once {
		_, err = os.Stat(fn)
		if err != nil {
			break
		}

		err := os.Remove(fn)
		if err != nil {
			break
		}
	}

	return err
}

// Ensure we don't duplicate services.
func (me *Service) IsRegistered() bool {

	var ret bool

	state, _ := me.Status(DontPublishState)
	switch state.Current {
	case states.StateIdle:
		fallthrough
	case states.StateUnknown:
		fallthrough
	case states.StateError:
		fallthrough
	case states.StateInitializing:
		fallthrough
	case states.StateInitialized:
		fallthrough
	case states.StateUnregistered:
		ret = false

	default:
		ret = true
	}

	return ret
}

func (j *ServiceUrl) UnmarshalJSON(b []byte) error {
	// Strip off the surrounding quotes and add a domain, one reason you might want a custom type
	u, err := url.Parse(fmt.Sprintf("%s", b[1:len(b)-1]))
	if err == nil {
		j.URL = u
	}

	return err
}

// Yup
func (me *Service) CreateMdnsEntry() (*network.ServiceConfig, error) {

	var err error
	var zc network.ServiceConfig

	for range only.Once {

		if me.Entry.MdnsType == "" {
			err = msgs.MakeError(me.EntityId, "MdnsType not set")
			break
		}

		foo := strings.ReplaceAll(me.Entry.Config.Name, ".", "_")
		zc = network.ServiceConfig{
			EntityId:   msgs.MakeAddress(),
			EntityName: me.EntityName,
			Name:       network.Name(strings.ToLower("_" + foo)),

			//@TODO Can we create a constant for "_%s._tcp"?  What would it's name be?
			//
			Type:   network.Type(fmt.Sprintf("_%s._tcp", me.Entry.MdnsType)),
			Domain: network.DefaultDomain,
			Port:   network.Port(me.Entry.UrlPtr.Port()),
		}
	}

	return &zc, err
}
