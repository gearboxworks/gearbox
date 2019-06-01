package tasks

import (
	"time"
)


type TaskFunc func(me *Task, you ...interface{}) error
type Tasks map[Uuid]*Task
var tasks = make(Tasks)


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
		task.RetryCounter = 0

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
					task.RetryCounter++

					task.runState = TaskStarting
					errorInside = startFunc(task, ref...)
					if errorInside == nil {
						task.runState = TaskStarted
						task.RetryCounter = 0
					}
				}

				// Exit conditions.
				if task.RetryLimit == 0 {
					// No limit set, keep going forever.
				} else if task.RetryCounter >= task.RetryLimit {
					// Reached the maximum number of retries, abort.
					break
				}

				if shouldStop() {
					// Have we been told to stop?
					break
				}

				// Sleep if we want to.
				if task.RetryDelay > 0 {
					time.Sleep(task.RetryDelay)
				}
			}
		}

		task.runLock = false
		return errorInside
	})

	if (err == nil) && (task != nil) {
		tasks[task.id] = task
		task.InitFunc = initFunc
		task.StartFunc = startFunc
		task.MonitorFunc = monitorFunc
		task.StopFunc = stopFunc
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
