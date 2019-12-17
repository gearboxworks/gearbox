package vmbox

import (
	"gearbox/box/external/virtualbox"
	"sync"
	"time"
)

var _ virtualbox.Consoler = (*Console)(nil)

type Console struct {
	Host      string
	Port      string
	ReadWait  time.Duration
	OkString  string
	WaitDelay time.Duration
	mutex     sync.RWMutex
}

func (me *Console) GetHost() string {
	return me.Host
}

func (me *Console) GetPort() string {
	return me.Port
}

func (me *Console) GetReadWait() time.Duration {
	return me.ReadWait
}

func (me *Console) GetOkString() string {
	return me.OkString
}

func (me *Console) GetWaitDelay() time.Duration {
	return me.WaitDelay
}

func (me *Console) GetMutex() sync.RWMutex {
	return me.mutex
}
