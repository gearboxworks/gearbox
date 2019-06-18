package tasks

//
// This package is intended to be moved into a standalone package with minimal dependencies.
//

import (
	"github.com/getlantern/errors"
	"github.com/google/uuid"
	"time"
)


// Go executes the function in a goroutine and returns a
// Task capable of stopping the execution.
func goRoutine(fn func(exitFunc, *Task) error) (*Task, error) {

	u, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	id := Uuid(u)

	t := &Task{
		id:           id,
		runState:     TaskIdle,
		runLock:      false,
		retryCounter: 0,
		retryLimit:   0,
		retryDelay:   time.Second,
		initFunc:     EmptyTask,
		startFunc:    EmptyTask,
		monitorFunc:  EmptyTask,
		stopFunc:     EmptyTask,

		stopChan:      make(chan struct{}),
		running:       true,
		//err: nil,
	}

	go func() {
		// call the target function
		err := fn(func() bool {
			// this is the shouldStop() function available to the
			// target function
			t.lock.RLock()
			shouldStop := t.shouldStop
			t.lock.RUnlock()
			return shouldStop
		}, t)
		// stopped
		t.lock.Lock()
		t.err = err
		t.running = false

		close(t.stopChan)
		t.lock.Unlock()
	}()

	return t, err
}


// Stop tells the goroutine to stop.
func (t *Task) stop() {

	if t == nil {
		return
	}

	// When task is stopped from a different go-routine other than the one
	// that actually started it.
	t.lock.Lock()
	t.shouldStop = true
	t.lock.Unlock()
}


// StopChan gets the stop channel for this task.
// Reading from this channel will block while the task is running, and will
// unblock once the task has stopped (because the channel gets closed).
func (t *Task) StopChan() <-chan struct{} {
	return t.stopChan
}


// Running gets whether the goroutine is
// running or not.
func (t *Task) Running() bool {

	if t == nil {
		return false
	}

	t.lock.RLock()
	running := t.running
	t.lock.RUnlock()

	return running
}


// Err gets the error returned by the goroutine.
func (t *Task) Err() error {

	if t == nil {
		return errors.New("No task")
	}

	t.lock.RLock()
	err := t.err
	t.lock.RUnlock()

	return err
}

