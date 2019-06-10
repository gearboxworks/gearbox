package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/only"
	"net"
	"reflect"
	"strings"
)


func (me *ZeroConf) EnsureNotNil() error {
	var err error

	switch {
		case me == nil:
			err = errors.New("ZeroConf instance is nil")
	}

	return err
}
func EnsureNotNil(me *ZeroConf) error {
	return me.EnsureNotNil()
}


func (me *ServicesMap) EnsureNotNil() error {

	var err error

	switch {
		case me == nil:
			err = errors.New("ZeroConf ServicesMap instance is nil")
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
			err = errors.New("ZeroConf Service instance is nil")
		case (me.instance == nil) && (me.IsManaged == true):
			err = me.EntityId.ProduceError("service instance is nil")
	}

	return err
}
func EnsureServicesNotNil(me *Service) error {
	return me.EnsureNotNil()
}


//func (me *ServicesArray) Print() error {
//
//	var err error
//
//	for range only.Once {
//		if me == nil {
//			err = errors.New("software error")
//			break
//		}
//
//		for i, e := range *me {
//			fmt.Printf("# Entry: #%d\n", i)
//			err = e.Print()
//			if err != nil {
//				break
//			}
//		}
//	}
//
//	return err
//}


// Ensure we don't duplicate services.
func (me *Service) IsExisting(him CreateEntry) error {

	var err error

	// @TODO - Need to check to see if this service has already been registered.
	//switch {
	//	case strconv.Itoa(me.Entry.Port) == him.Port.String():
	//		err = me.EntityId.ProduceError("service HostName:%s already exists", me.Entry.HostName)
	//
	//	case me.Entry.HostName == him:
	//		err = me.EntityId.ProduceError("service Name:%s already exists", me.Entry.Name)
	//}

	return err
}


// Ensure we don't duplicate services.
func (me *ServicesMap) IsExisting(him CreateEntry) error {

	var err error

	for _, ce := range *me {
		err = ce.IsExisting(him)
		if err != nil {
			break
		}
	}

	return err
}


func ConstructMdnsRegisterMessage(me messages.MessageAddress, to messages.MessageAddress, s CreateEntry) messages.Message {

	var err error
	var msgTemplate messages.Message
	var j []byte

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		j, err = json.Marshal(s)
		if err != nil {
			break
		}

		msgTemplate = messages.Message{
			Source: me,
			Topic: messages.MessageTopic{
				Address:  to,
				SubTopic: states.ActionRegister,
			},
			Text: messages.MessageText(j),
		}
	}

	return msgTemplate
}


func DeconstructMdnsRegisterMessage(event *messages.Message) (CreateEntry, error) {

	var err error
	var ce CreateEntry

	for range only.Once {
		//err = ce.EnsureNotNil()
		if event == nil {
			err = errors.New("message is nil")
			break
		}

		err = json.Unmarshal(event.Text.ByteArray(), &ce)
		if err != nil {
			break
		}
	}

	return ce, err
}


func ConstructMdnsUnregisterMessage(me messages.MessageAddress, to messages.MessageAddress, s CreateEntry) messages.Message {

	var err error
	var msgTemplate messages.Message
	var j []byte

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		j, err = json.Marshal(s)
		if err != nil {
			break
		}

		msgTemplate = messages.Message{
			Source: me,
			Topic: messages.MessageTopic{
				Address:  to,
				SubTopic: states.ActionUnregister,
			},
			Text: messages.MessageText(j),
		}
	}

	return msgTemplate
}


func InterfaceToTypeZeroConf(i interface{}) (*ZeroConf, error) {

	var err error
	var zc *ZeroConf

	for range only.Once {
		err = channels.EnsureArgumentNotNil(i)
		if err != nil {
			break
		}

		checkType := reflect.ValueOf(i)
		if checkType.Type().String() != "*network.ZeroConf" {
			err = errors.New("interface type not *network.ZeroConf")
			break
		}

		zc = i.(*ZeroConf)
		// zc = (i[0]).(*ZeroConf)
		// zc = i[0].(*ZeroConf)

		err = zc.EnsureNotNil()
		if err != nil {
			break
		}
	}

	return zc, err
}


func InterfaceToTypeService(i interface{}) (*Service, error) {

	var err error
	var s *Service

	for range only.Once {
		err = channels.EnsureArgumentNotNil(i)
		if err != nil {
			break
		}

		checkType := reflect.ValueOf(i)
		if checkType.Type().String() != "*network.Service" {
			err = errors.New("interface type not *network.Service")
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

		if (me.instance == nil) && (me.IsManaged == true) {
			fmt.Printf("# Entry(deleted): %v", me.EntityId)
		} else {
			fmt.Printf("# Entry: %v", me.EntityId)
		}
		err = me.Entry.Print()
		if err != nil {
			break
		}
	}

	return err
}


func (me *Entry) IsTheSame(e Entry) (bool, error) {

	var same bool
	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("software error")
			same = false
			break
		}

		if (me.Instance == e.Instance) &&
			(trimDot(me.Service) == trimDot(e.Service)) &&
			(trimDot(me.Domain) == trimDot(e.Domain)) &&
			(me.Port == e.Port) {
			same = true
		}
	}

	return same, err
}


func (me *Entry) UpdateService(e Entry) (bool, error) {

	var same bool
	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("software error")
			break
		}

		ok, err := me.IsTheSame(e)
		if err != nil {
			break
		}
		if ok {
			me.Text = e.Text
			me.AddrIPv4 = e.AddrIPv4
			me.AddrIPv6 = e.AddrIPv6
			me.HostName = e.HostName
			me.TTL = e.TTL
			same = true
		}

		// fmt.Printf("DEBUG1:\n%v\n", zeroconf.NewServiceEntry(me.Instance, me.Service, me.Domain))
		// fmt.Printf("DEBUG2:\n%v\n", me)
	}

	return same, err
}


// Replace zeroconf.ServiceEntry.ServiceName() function with our own.
func (me *Entry) ServiceName() (string, error) {

	var sn string
	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("software error")
			break
		}

		sn = fmt.Sprintf("%s.%s.", trimDot(me.Service), trimDot(me.Domain))
	}

	return sn, err
}


// Replace zeroconf.ServiceEntry.ServiceInstanceName() function with our own.
func (me *Entry) ServiceInstanceName() (string, error) {

	var sin string
	var sn string
	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("software error")
			break
		}

		sn, err = me.ServiceName()
		if err != nil {
			err = errors.New("software error")
			break
		}

		sin = fmt.Sprintf("%s.%s", trimDot(me.Instance), sn)
	}

	return sin, err
}


// Replace zeroconf.ServiceEntry.ServiceTypeName() function with our own.
func (me *Entry) ServiceTypeName() (string, error) {

	var sn string
	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("software error")
			break
		}

		sn = fmt.Sprintf("_services._dns-sd._udp.%s.", trimDot(me.Domain))
	}

	return sn, err
}


// trimDot is used to trim the dots from the start or end of a string
func trimDot(s string) string {
	return strings.Trim(s, ".")
}


func (me *Entry) Print() error {

	var err error

	for range only.Once {
		if me == nil {
			err = errors.New("software error")
			break
		}

		sn,_ := me.ServiceName()
		sin,_ := me.ServiceInstanceName()
		stn,_ := me.ServiceTypeName()

		//
		fmt.Printf(` me.Instance = %v
 me.Service = %v
 me.Domain = %v
 me.Port = %v
 me.Text = %v
 me.AddrIPv4 = %v
 me.AddrIPv6 = %v
 me.HostName = %v
 me.TTL = %v
 me.ServiceName() = %v
 me.ServiceInstanceName() = %v
 me.ServiceTypeName() = %v
`,
			me.Instance,
			me.Service,
			me.Domain,
			me.Port,
			me.Text,
			me.AddrIPv4,
			me.AddrIPv6,
			me.HostName,
			me.TTL,
			sn,
			sin,
			stn,
		)
	}

	return err
}


// GetFreePort asks the kernel for a free open port that is ready to use.
func GetFreePort() (Port, error) {

	var port int

	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()

	port = l.Addr().(*net.TCPAddr).Port
	eblog.Debug(DefaultEntityId, "Foung a free port on == %d", port)

	return Port(port), nil
}


// GetFreePort asks the kernel for free open ports that are ready to use.
func GetFreePorts(count int) ([]int, error) {
	var ports []int
	for i := 0; i < count; i++ {
		addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
		if err != nil {
			return nil, err
		}

		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			return nil, err
		}
		defer l.Close()
		ports = append(ports, l.Addr().(*net.TCPAddr).Port)
	}
	return ports, nil
}

