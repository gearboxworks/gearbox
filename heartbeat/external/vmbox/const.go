package vmbox

import (
	"bytes"
	"gearbox/eventbroker/channels"
	"gearbox/eventbroker/messages"
	"gearbox/eventbroker/ospaths"
	"gearbox/eventbroker/states"
	"gearbox/eventbroker/tasks"
	"github.com/gearboxworks/go-status"
	"sync"
	"time"
)


const (
	DefaultVmWaitTime = time.Millisecond * 1000
	DefaultRetries  = 12
	DefaultRetryDelay = time.Second * 10
)


type VmBox struct {
	EntityId       messages.MessageAddress
	EntityName     messages.MessageAddress
	EntityParent   messages.MessageAddress
	Boxname         string
	State           *states.Status
	Task            *tasks.Task
	Channels        *channels.Channels
	Releases        Releases

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
	EntityId       messages.MessageAddress
	EntityName     messages.MessageAddress
	EntityParent   *messages.MessageAddress
	State          *states.Status
	IsManaged      bool
	Entry          *ServiceConfig

	mutex          sync.RWMutex // Mutex control for this struct.
	channels       *channels.Channels
	channelHandler *channels.Subscriber
	osPaths        *ospaths.BasePaths
}

type ServiceConfig struct {
	Name           messages.MessageAddress
	Version        string
	//IsoFile        string
	ConsolePort    string
	SshPort        string

	retryMax       int
	retryDelay     time.Duration
	baseDir        *ospaths.Dir
	osRelease      *Release
	cmdStdout      bytes.Buffer
	cmdStderr      bytes.Buffer
	vmInfo         KeyValueMap
	vmNics         KeyValuesMap
}

const (
	Package                    = "vmBox"
	InterfaceTypeVmBox    = "*" + Package + ".VmBox"
	InterfaceTypeService       = "*" + Package + ".Service"
	InterfaceTypeServiceConfig = "*" + Package + ".ServiceConfig"
)

const (
	Basedir     = "/home/gearbox/projects"
)

const (
	VmStateInit			= ""
	VmStateNotPresent	= "not present"
	VmStatePowerOff 	= "poweroff"	// Valid VM state return from listvm
	VmStatePaused 		= "paused"		// Valid VM state return from listvm
	VmStateSaved 		= "saved"		// Valid VM state return from listvm
	VmStateRunning  	= "running"		// Valid VM state return from listvm
	VmStateStarting		= "starting"
	VmStateStopping		= "stopping"
	VmStateDontCare		= "dont care"
	VmStateUnknown		= "unknown"

	OkState      = "ok"
	NotOkState   = "nok"
)


const (
	IconLogoPng = "heartbeat/img/IconLogo.png"
)


type BoxEntity struct {
	Name         string
	CurrentState string
	WantState    string
}

type BoxState struct {
	VM      BoxEntity
	API     BoxEntity
	LastSts status.Status
}

type VmDisplayState struct {
	Name			string

	VmIconState		string
	VmTitleState	string
	VmHintState		string

	ApiIconState	string
	ApiTitleState	string
	ApiHintState	string

	Title			string
	Hint			string
	Sts	    		status.Status
}

