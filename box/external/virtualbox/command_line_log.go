package virtualbox

import (
	"bytes"
	"fmt"
)

var _ Logger = (*CommandLineLog)(nil)

type CommandLineLog struct {
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
}

func NewCommandLineLog() *CommandLineLog {
	cll := CommandLineLog{}
	cll.Reset()
	return &cll
}

func (me *CommandLineLog) Reset() {
	me.Stdout = &bytes.Buffer{}
	me.Stderr = &bytes.Buffer{}
}

func (me *CommandLineLog) Error(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
}

func (me *CommandLineLog) Debug(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
}

func (me *CommandLineLog) GetStdout() *bytes.Buffer {
	return me.Stdout
}

func (me *CommandLineLog) GetStderr() *bytes.Buffer {
	return me.Stderr
}
