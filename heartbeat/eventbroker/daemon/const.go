package daemon

import (
	"fmt"
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
	DefaultRetryDelay = time.Second * 50
)


type Daemon struct {
	EntityId       messages.MessageAddress
	State          states.Status
	Task           *tasks.Task
	Channels       *channels.Channels
	ChannelHandler *channels.Subscriber
	Fluff          string

	mutex          sync.RWMutex	// Mutex control for this struct.
	daemons        ServicesMap
	osSupport      oss.OsSupporter
}
type Args Daemon


type Service struct {
	EntityId  messages.MessageAddress
	State     states.Status
	IsManaged bool
	JsonFile  string
	Entry     *CreateEntry

	mutex     sync.RWMutex	// Mutex control for map.
	channels  *channels.Channels
	instance  programInstance
}
//type ServicesMap struct {
//	instance map[messages.MessageAddress]*Service
//	mutex sync.RWMutex	// Mutex control for map.
//}
type ServicesMap  map[messages.MessageAddress]*Service

type CreateEntry struct {
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

type programInstance struct {
	exit    chan struct{}
	service service.Service
	cmd     *exec.Cmd
	*service.Config
}

func (me *CreateEntry) ToServiceType() *service.Config {

	return &me.Config
}

// Ensure we don't duplicate services.
func (me *Service) IsExisting(him CreateEntry) error {

	var err error

	switch {
		case me.Entry.Config.Name == him.Config.Name:
			err = me.EntityId.ProduceError("Daemon service Name:%s already exists", me.Entry.Config.Name)

		case me.Entry.Config.DisplayName == him.Config.DisplayName:
			err = me.EntityId.ProduceError("Daemon service DisplayName:%s already exists", me.Entry.DisplayName)

		case me.Entry.Config.Executable == him.Config.Executable:
			err = me.EntityId.ProduceError("Daemon service Executable:%s already exists", me.Entry.Config.Executable)

		case me.Entry.Url == him.Url:
			err = me.EntityId.ProduceError("Daemon service Url:%s already exists", me.Entry.Url)

		case (me.Entry.Host == him.Host) && (me.Entry.Port == him.Port):
			err = me.EntityId.ProduceError("Daemon service Host:%s:%s already exists", me.Entry.Host.String(), me.Entry.Port.String())
	}

	return err
}

// Ensure we don't duplicate services.
func (me *Daemon) IsExisting(him CreateEntry) error {

	var err error

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	for _, ce := range me.daemons {	// Managed by Mutex
		err = ce.IsExisting(him)
		if err != nil {
			break
		}
	}

	return err
}

//execCwd, _ := os.Getwd()
//if execCwd == "/" {
//execCwd = string(OsSupport.GetAdminRootDir())
//}
//_args.ServiceData.Path = execCwd

type ServiceUrl struct {
	*url.URL
}

func (j *ServiceUrl) UnmarshalJSON(b []byte) error {
	// Strip off the surrounding quotes and add a domain, one reason you might want a custom type

	u, err := url.Parse(fmt.Sprintf("%s", b[1:len(b)-1]))
	if err == nil {
		j.URL = u
	}

	return err
}

