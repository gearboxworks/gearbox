package hypervisor

import (
	"bytes"
	"gearbox/box/external/vbmanage"
	"gearbox/box/external/virtualbox"
	"sync"
	"time"
)


type CommandResulter interface {
	GetExitCode() string
	GetStdout() *bytes.Buffer
	GetStderr() *bytes.Buffer
}

type Consoler interface {
	GetHost() string
	GetPort() string
	GetReadWait() time.Duration
	GetOkString() string
	GetWaitDelay() time.Duration
	GetMutex() sync.RWMutex
}

type Configer interface {
	Destroy() error
}

type VirtualMachiner interface {
	GetId() string
	GetName() string
	GetVmDir() string
	GetIconFile() string
	SetInfo(kvm vbmanage.KeyValueMap)
	GetInfo() vbmanage.KeyValueMap
	GetInfoValue(string) string
	GetNics() vbmanage.KeyValueMap
	SetNics(kvm vbmanage.KeyValuesMap)
	GetNic() *virtualbox.HostOnlyNic
	SetNic(*virtualbox.HostOnlyNic)
	GetConsole() Consoler
	GetSsh() SecureSheller
	GetReleaser() Releaser
	GetRetryDelay() time.Duration
	GetRetryMax() int

	// @TODO This is not in the right place, but it's here to allow refactoring
	DestroyConfig() error
}

type Logger interface {
	Error(string, ...interface{})
	Debug(string, ...interface{})
	GetStdout() *bytes.Buffer
	GetStderr() *bytes.Buffer
}

type Diskers []Disker
type Disker interface {
	GetName() string
	GetFormat() string
	GetSize() string
}

type SecureSheller interface {
	GetHost() string
	GetPort() string
}

type Releaser interface {
	GetFilepath() string
}
