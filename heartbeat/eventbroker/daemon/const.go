package daemon

import (
	"gearbox/heartbeat/eventbroker/channels"
	"gearbox/heartbeat/eventbroker/messages"
	"gearbox/heartbeat/eventbroker/network"
	"gearbox/heartbeat/eventbroker/states"
	"gearbox/heartbeat/eventbroker/tasks"
	oss "gearbox/os_support"
	"github.com/kardianos/service"
	"net/url"
	"os/exec"
	"sync"
	"time"
)


const (
	DefaultEntityId = "eventbroker-daemon"

	defaultBaseDir = "eventbroker"
	defaultLogBaseDir = defaultBaseDir + "/logs"
	defaultEtcBaseDir = defaultBaseDir + "/etc"

	defaultJsonFile  = defaultEtcBaseDir + "/" + DefaultEntityId + ".json"
	defaultLogFile  = defaultLogBaseDir + "/"  + DefaultEntityId + ".log"
	defaultErrorLogFile  = defaultLogBaseDir + DefaultEntityId + "-error.log"
	DefaultJsonFiles  = defaultEtcBaseDir + "/daemons"

    defaultWaitTime = time.Millisecond * 2000
	defaultDomain   = "local"
	defaultRetries  = 12
	DefaultRetryDelay = time.Second * 20
)


type Daemon struct {
	EntityId       messages.MessageAddress
	State          states.Status
	Task           *tasks.Task
	Channels       *channels.Channels

	mutex          sync.RWMutex	// Mutex control for this struct.
	channelHandler *channels.Subscriber
	daemons        ServicesMap
	osSupport      oss.OsSupporter
}
type Args Daemon


type Service struct {
	EntityId        messages.MessageAddress
	State           states.Status
	IsManaged       bool
	Entry           *ServiceConfig
	JsonFile        JsonConfig

	mutex           sync.RWMutex // Mutex control for map.
	channels        *channels.Channels
	channelHandler  *channels.Subscriber
	instance        programInstance
}
type ServicesMap  map[messages.MessageAddress]*Service

type JsonConfig struct {
	Name        string
	LastModTime time.Time
}

type ServiceConfig struct {
	service.Config

	Stdout   string
	Stderr   string
	Stdin    string
	Url      string
	Env      []string
	Host     network.Host
	Port     network.Port
	MdnsType string
}

const (
	Package                    = "daemon"
	InterfaceTypeService       = "*" + Package + ".Service"
	InterfaceTypeServiceConfig = "*" + Package + ".ServiceConfig"
	InterfaceTypeDaemon        = "*" + Package + ".Daemon"
	InterfaceTypeError         = "error"
)

type programInstance struct {
	exit    chan struct{}
	service service.Service
	cmd     *exec.Cmd
	*service.Config
}

type ServiceUrl struct {
	*url.URL
}

