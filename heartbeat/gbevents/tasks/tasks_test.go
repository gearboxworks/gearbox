package tasks_test

import (
	"errors"
	"fmt"
	"gearbox/heartbeat/gbevents/tasksnts/tasks"
	"github.com/cheekybits/is"
	"testing"
	"time"
)




func StartTest() error {

	task, err := StartTask(initMe, startMe, monitorMe, stopMe, loopCount)

	time.Sleep(time.Second * 10)
	task.Stop()

	return err
}

type Foo int
var loopCount Foo
func (me *Foo) Test() error {
	var err error

	switch {
	case me.Get() >= 11:
		err = errors.New("FAIL!")
	case me.Get() >= 8:
		err = nil
	case me.Get() >= 5:
		err = errors.New("FAIL!")
	case me.Get() >= 2:
		err = nil
	case me.Get() >= 1:
		err = errors.New("FAIL!")
	case me.Get() >= 0:
		err = nil
	}
	me.Add()
	fmt.Printf("return[%v]: %v\n", me.Get(), err)

	return err
}
func (me *Foo) Get() Foo {
	return loopCount
}
func (me *Foo) Add() {
	loopCount++
}
func (me *Foo) Minus() {
	loopCount--
}
func (me *Foo) Zero() {
	loopCount = 0
}


func initMe(task *Task, i ...interface{}) error {

	var err error
	//me := (*i[0]).(int)
	me := (i[0]).(Foo)
	fmt.Printf("initMe[%v]:	retry:%v me:%v \n", task.RunState, task.RetryCounter, me.Get())
	task.RetryLimit = 2
	me.Zero()

	err = me.Test()
	return err
}

func startMe(task *Task, i ...interface{}) error {

	var err error
	me := (i[0]).(Foo)
	fmt.Printf("startMe[%v]:	retry:%v me:%v \n", task.RunState, task.RetryCounter, me.Get())

	err = me.Test()
	return err
}

func monitorMe(task *Task, i ...interface{}) error {

	var err error
	me := (i[0]).(Foo)
	fmt.Printf("monitorMe[%v]:	retry:%v me:%v \n", task.RunState, task.RetryCounter, me.Get())

	err = me.Test()
	return err
}

func stopMe(task *Task, i ...interface{}) error {

	var err error
	me := (i[0]).(Foo)
	fmt.Printf("stopMe[%v]:	retry:%v me:%v \n", task.RunState, task.RetryCounter, me.Get())

	err = me.Test()
	return err
}


func TestRun(t *testing.T) {
	is := is.New(t)
	var ticker []time.Time
	task := tasks.Go(func(shouldStop tasks.S) error {
		for {
			ticker = append(ticker, time.Now())
			time.Sleep(100 * time.Millisecond)
			if shouldStop() {
				break
			}
		}
		return nil
	})
	is.Equal(true, task.Running())
	time.Sleep(1 * time.Second)
	task.Stop()
	select {
	case <-task.StopChan():
	case <-time.After(2 * time.Second):
		is.Fail("timed out")
	}
	is.Equal(false, task.Running())
	is.Equal(10, len(ticker))
}

func TestRunErr(t *testing.T) {
	is := is.New(t)

	err := errors.New("something went wrong")
	task := tasks.Go(func(shouldStop tasks.S) error {
		return err
	})

	time.Sleep(100 * time.Millisecond)
	is.Equal(false, task.Running())
	is.Equal(err, task.Err())

	task.Stop()
	select {
	case <-task.StopChan():
	case <-time.After(2 * time.Second):
		is.Fail("timed out")
	}

}
