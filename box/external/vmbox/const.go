package vmbox

import (
	"bytes"
	"gearbox/eventbroker/channels"
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/ospaths"
	"gearbox/eventbroker/states"
	"gearbox/eventbroker/tasks"
	"sync"
	"time"
)


const (
	DefaultVmWaitTime = time.Millisecond * 1000
	DefaultRetries  = 12
	DefaultRetryDelay = time.Second * 10
	DefaultConsoleReadWait = time.Second * 5
	DefaultConsoleWaitDelay = time.Second
	DefaultBootWaitTime = time.Second * 60
	DefaultRunWaitTime = time.Second * 5

	DefaultConsoleOkString = "Gearbox Heartbeat"
	Basedir     = "/home/gearbox/projects"
	IconLogoPng = "app/dist/heartbeat/img/IconLogo.png"

	DefaultHoIp = "192.168.42.1"
	DefaultHoNetmask = "255.255.255.0"
	DefaultHoDhcpLowerIp = "192.168.42.10"
	DefaultHoDhcpUpperIp = "192.168.42.254"
)


type VmBox struct {
	EntityId       messages.MessageAddress
	EntityName     messages.MessageAddress
	EntityParent   messages.MessageAddress
	Boxname         string
	State           *states.Status
	Task            *tasks.Task
	Channels        *channels.Channels
	Releases        *Releases

	mutex           sync.RWMutex // Mutex control for map.
	channelHandler  *channels.Subscriber
	restartAttempts int
	waitTime        time.Duration
	vms             VmMap
	OsPaths         *ospaths.BasePaths
}
type Args VmBox

type VmMap map[messages.MessageAddress]*Vm
type Vm struct {
	EntityId        messages.MessageAddress
	EntityName      messages.MessageAddress
	EntityParent    *messages.MessageAddress
	State           *states.Status
	ApiState        *states.Status
	ChangeRequested bool
	IsManaged       bool
	Entry           *ServiceConfig

	mutex           sync.RWMutex // Mutex control for this struct.
	channels        *channels.Channels
	channelHandler  *channels.Subscriber
	osRelease       *Release
	osPaths         *ospaths.BasePaths
	// baseDir         *ospaths.Dir
}

type ServiceConfig struct {
	Name             messages.MessageAddress
	Version          string
	Console          Console
	Ssh              Ssh
	IconFile         *ospaths.File
	VmDir            *ospaths.Dir
	HostOnlyNic      HostOnlyNic

	retryMax         int
	retryDelay       time.Duration
	cmdStdout        bytes.Buffer
	cmdStderr        bytes.Buffer
	vmInfo           KeyValueMap
	vmNics           KeyValuesMap
}

type Console struct {
	Host      string
	Port      string
	ReadWait  time.Duration
	OkString  string
	WaitDelay time.Duration
	mutex     sync.RWMutex
}

type Ssh struct {
	Host string
	Port string
}

type HostOnlyNic struct {
	Name        string
	Index       int
	Ip          string
	Netmask     string
	DhcpLowerIp string
	DhcpUpperIp string
}


const (
	Package                    = "vmbox"
	InterfaceTypeVmBox    = "*" + Package + ".VmBox"
	InterfaceTypeService       = "*" + Package + ".Service"
	InterfaceTypeServiceConfig = "*" + Package + ".ServiceConfig"
)

