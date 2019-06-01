package network

import (
	"fmt"
	"errors"
	"net"
)

func (me *Client) EnsureNotNil() error {
	var err error

	if me == nil {
		err = errors.New("unexpected software error")
	}

	return err
}


func (me *ServiceEntries) Print() {

	for i, e := range *me {
		fmt.Printf("- Entry #%d\n", i)
		e.Print()
	}
}


func (me *ServiceEntry) Print() {

	//
	fmt.Printf(`
me.Text = %v
me.Port = %v
me.Domain = %v
me.Instance = %v
me.Service = %v
me.AddrIPv4 = %v
me.AddrIPv6 = %v
me.HostName = %v
me.TTL = %v
me.ServiceRecord = %v
me.ServiceName() = %v
me.ServiceInstanceName() = %v
me.ServiceTypeName() = %v
`,
		me.Text,
		me.Port,
		me.Domain,
		me.Instance,
		me.Service,
		me.AddrIPv4,
		me.AddrIPv6,
		me.HostName,
		me.TTL,
		me.ServiceRecord,
		me.ServiceName(),
		me.ServiceInstanceName(),
		me.ServiceTypeName(),
		)


	return
}


// GetFreePort asks the kernel for a free open port that is ready to use.
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
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
