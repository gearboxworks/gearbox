package tasks

import (
	"errors"
	"gearbox/eventbroker/eblog"
	"github.com/google/uuid"
	"time"
)


func StartTask(initFunc TaskFunc, startFunc TaskFunc, monitorFunc TaskFunc, stopFunc TaskFunc, ref ...interface{}) (*Task, error) {

	var err error
	var task *Task

	task, err = goRoutine(func(shouldStop exitFunc, task *Task) error {
		var errorInside error

		if task.runLock == true {
			return errorInside
		}
		task.runLock = true
		task.runState = TaskIdle
		task.retryCounter = 0

		defer func() {
			if stopFunc != nil {
				task.runState = TaskStopping
				errorInside = stopFunc(task, ref...)
				if errorInside != nil {
					task.runState = TaskStopped
				}
			}
		}()

		if initFunc != nil {
			task.runState = TaskInitializing
			errorInside = initFunc(task, ref...)
			if errorInside == nil {
				task.runState = TaskInitialized
			}
		}

		if errorInside == nil {
			for {
				// Execute the run function.
				errorInside = monitorFunc(task, ref...)
				if errorInside != nil {
					// Keep track of the number of failed restarts.
					task.retryCounter++

					task.runState = TaskStarting
					errorInside = startFunc(task, ref...)
					if errorInside == nil {
						task.runState = TaskStarted
						task.retryCounter = 0
					}
				}

				// Exit conditions.
				if task.retryLimit == 0 {
					// No limit set, keep going forever.
				} else if task.retryCounter >= task.retryLimit {
					// Reached the maximum number of retries, abort.
					break
				}

				if shouldStop() {
					// Have we been told to stop?
					break
				}

				// Sleep if we want to.
				if task.retryDelay > 0 {
					time.Sleep(task.retryDelay)
				}

				if shouldStop() {
					// Have we been told to stop?
					break
				}
			}
		}

		task.runLock = false
		return errorInside
	})

	if (err == nil) && (task != nil) {
		tasks[task.id] = task
		task.initFunc = initFunc
		task.startFunc = startFunc
		task.monitorFunc = monitorFunc
		task.stopFunc = stopFunc
	}

	if eblog.LogIfError(eblog.SkipNilCheck, err) {
		// Save last state.
		// me.State.Error = err
	}

	return task, err
}


func ListTasks() (Tasks, error) {

	return tasks, nil
}


func EmptyTask(task *Task, i ...interface{}) error {
	return nil
}


// Mirrors tasks.stop() with extra enhancements.
// EG: Wait for the process to actually stop
func (me *Task) Stop() error {

	if me == nil {
		return errors.New("non-existant task")
	}

	me.stop()
	select {
		case <-me.StopChan():
			// task successfully stopped
		case <-time.After(10 * time.Second):
			// task didn't stop in time
	}

	delete(tasks, me.id)

	return me.Err()
}


// Mirrors tasks.StopChan() with extra enhancements.
// This will wait indefinitely until a task has stopped.
func (me *Task) WaitUntilStopped() bool {

	if me == nil {
		return false
	}

	wait := me.StopChan()
	if wait != nil {
		return true
	}

	return true
}


// Mirrors tasks.Running().
func (me *Task) IsRunning() bool {

	if me == nil {
		return false
	}

	return me.Running()
}


func (me *Task) GetRetryLimit() int {

	if me == nil {
		return 0
	}

	return me.retryLimit
}


func (me *Task) SetRetryLimit(v int) error {

	if me == nil {
		return errors.New("unexpected software error")
	}

	me.retryLimit = v

	return nil
}


func (me *Task) GetRetryCounter() int {

	if me == nil {
		return 0
	}

	return me.retryCounter
}


func (me *Task) SetRetryCounter(v int) error {

	if me == nil {
		return errors.New("unexpected software error")
	}

	me.retryCounter = v

	return nil
}


func (me *Task) GetRetryDelay() time.Duration {

	if me == nil {
		return 0
	}

	return me.retryDelay
}


func (me *Task) SetRetryDelay(v time.Duration) error {

	if me == nil {
		return errors.New("unexpected software error")
	}

	me.retryDelay = v

	return nil
}


func (me *Task) GetState() State {

	if me == nil {
		return "unknown"
	}

	return me.runState
}


func (me *Task) GetId() Uuid {

	empty := Uuid(uuid.UUID{})

	if me == nil {
		return empty
	}

	return me.id
}
