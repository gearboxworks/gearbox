package tasks

import (
	"github.com/google/uuid"
	"sync"
	"time"
)


const (
	TaskIdle = State("idle")
	TaskInitializing = State("initializing")
	TaskInitialized = State("initialized")
	TaskStopped = State("stopped")
	TaskStarting = State("starting")
	TaskStarted = State("started")
	TaskStopping = State("stopping")
)


// S is a function that will return true if the
// goroutine should stop executing.
type exitFunc func() bool
type StateIndex int
type State string
type Uuid uuid.UUID


// Task represents an interruptable goroutine.
type Task struct {
	id           Uuid
	runState     State
	runLock      bool
	retryCounter int
	retryLimit   int
	retryDelay   time.Duration
	initFunc     TaskFunc
	startFunc    TaskFunc
	monitorFunc  TaskFunc
	stopFunc     TaskFunc

	lock          sync.RWMutex
	stopChan      chan struct{}
	shouldStop    bool
	running       bool
	err           error
}

type TaskFunc func(me *Task, you ...interface{}) error
type Tasks map[Uuid]*Task
var tasks = make(Tasks)
