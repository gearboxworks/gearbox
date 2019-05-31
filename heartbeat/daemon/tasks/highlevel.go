package tasks

import (
	"github.com/google/uuid"
	"time"
)

type TaskFunc func(me *Task, you ...interface{}) error
type Tasks map[uuid.UUID]*Task
var tasks = make(Tasks)


func StartTask(initFunc TaskFunc, runFunc TaskFunc, endFunc TaskFunc, ref ...interface{}) (*Task, error) {

	var err error
	var task *Task

	task, err = goRoutine(func(shouldStop stopFunc, task *Task) error {
		var errorInside error

		defer func() {
			errorInside = endFunc(task, ref...)
		}()

		errorInside = initFunc(task, ref...)
		if errorInside == nil {
			for {
				task.RunCounter++

				errorInside = runFunc(task, ref...)

				// Should we stop or not?
				if shouldStop() {
					break
				}
			}
		}

		return errorInside
	})

	if (err == nil) && (task != nil) {
		tasks[task.ID] = task
		tasks[task.ID].InitFunc = initFunc
		tasks[task.ID].RunFunc = runFunc
		tasks[task.ID].EndFunc = endFunc
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

	delete(tasks, me.ID)

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
