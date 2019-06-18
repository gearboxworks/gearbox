package external


import (
	"sync"
)


const (
	StateInit     = 0
	StateUnknown  = 1
	StateDontCare = 2
	StatePowerOff = 3
	StateStarting = 4
	StateRunning  = 5
	StateStopping = 6
	StateLoaded   = 7
	StateUnloaded = 8
)


type TaskFunc func(me *Task, you ...interface{}) error
type Tasks map[string]*Task
type Task struct {
	Func       TaskFunc
	JsonFile   string
	JsonString string
	lock       sync.RWMutex
}


//