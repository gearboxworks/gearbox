package hypervisor

import (
	"bytes"
	"fmt"
)

var _ Logger = (*SimpleLog)(nil)

type SimpleLog struct {
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
}

func NewSimpleLog() *SimpleLog {
	cll := SimpleLog{}
	cll.Reset()
	return &cll
}

func (me *SimpleLog) Reset() {
	me.Stdout = &bytes.Buffer{}
	me.Stderr = &bytes.Buffer{}
}

func (me *SimpleLog) Error(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
}

func (me *SimpleLog) Debug(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
}

func (me *SimpleLog) GetStdout() *bytes.Buffer {
	return me.Stdout
}

func (me *SimpleLog) GetStderr() *bytes.Buffer {
	return me.Stderr
}

