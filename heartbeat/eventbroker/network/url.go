package network

import (
	"fmt"
	"gearbox/heartbeat/eventbroker/eblog"
	"gearbox/heartbeat/eventbroker/entity"
	"gearbox/heartbeat/eventbroker/only"
	"net"
	"net/url"
	"strconv"
)


func ParseUrl(us string) (*url.URL, error) {

	var uri *url.URL
	var err error

	for range only.Once {
		// u, err = url.Parse(fmt.Sprintf("tcp://%s:%d", mqttService.Entry.HostName, mqttService.Entry.Port))

		uri, err = url.Parse(us)
		if err != nil {
			break
		}
		if uri == nil {
			break
		}

		var p Port
		if uri.Port() == "0" {
			p, err = GetFreePort()
			if err != nil {
				break
			}
			//uri.Host = uri.Host + ":" + p.String()
		}

		if uri.Scheme == "" {
			uri.Scheme = "tcp"
		}

		uri, err = url.Parse(fmt.Sprintf("%s://%s:%s", uri.Scheme, uri.Hostname(), p.String()))
		if err != nil {
			break
		}
		if uri == nil {
			break
		}

	}

	return uri, err
}


func (port *Port) IfZeroFindFreePort() (err error) {

	switch *port {
		case "0":
			fallthrough
		case "":
			*port, err = GetFreePort()
	}

	return err
}


// GetFreePort asks the kernel for a free open port that is ready to use.
func GetFreePort() (Port, error) {

	var port int

	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "0", err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "0", err
	}
	defer l.Close()

	port = l.Addr().(*net.TCPAddr).Port
	eblog.Debug(entity.NetworkEntityName, "Foung a free port on == %d", port)

	return Port(strconv.Itoa(port)), nil
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

