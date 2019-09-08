package vbmanage

import (
	"bufio"
	"bytes"
	"fmt"
	"gearbox/box/external/hypervisor"
	"github.com/gearboxworks/go-status/only"
	"strings"
	"unicode"
)

var _ hypervisor.CommandResulter = (*CommandResult)(nil)

type CommandResult struct {
	ExitCode string
	Command string
	StdoutBuffer *bytes.Buffer
	StderrBuffer *bytes.Buffer
	Error error
}

func NewCommandResult() *CommandResult {
	cll := CommandResult{}
	cll.Reset()
	return &cll
}

func (me *CommandResult) Reset() {
	me.StdoutBuffer = &bytes.Buffer{}
	me.StderrBuffer = &bytes.Buffer{}
}

type FileDescriptor = int
const Stdout FileDescriptor = 1
const Stderr FileDescriptor = 2

func (me *CommandResult) GetExitCode() string {
	return me.ExitCode
}

func (me *CommandResult) GetStdout() *bytes.Buffer {
	return me.StdoutBuffer
}

func (me *CommandResult) GetStderr() *bytes.Buffer {
	return me.StderrBuffer
}

func (me *CommandResult) String() string {
	s := fmt.Sprintf("ExitCode: %s, Stdout: %s, Stderr: %s",
		me.ExitCode,
		me.GetStdoutLn(),
		me.GetStderrLn(),
	)
	if me.Error != nil {
		s = fmt.Sprintf("[%s]: %v", s, me.Error)
	}
	return s
}

func (me *CommandResult) DecodeStdout(splitOn rune) (dr KeyValues, ok bool) {
	return me.Decode(Stdout,splitOn)
}

func (me *CommandResult) DecodeStderr(splitOn rune) (dr KeyValues, ok bool) {
	return me.Decode(Stderr,splitOn)
}

func (me *CommandResult) Decode(fd FileDescriptor, splitOn rune) (dr KeyValues, ok bool) {
	var b *bytes.Buffer
	var n string
	for range only.Once {
		switch fd {
		case Stdout:
			b = me.StdoutBuffer
		case Stderr:
			b = me.StderrBuffer
		}

		if b == nil {
			break
		}
		ok = false
		s := bufio.NewScanner(b)
		for s.Scan() {
			kv, _ok := me.extractKeyValue(s.Text(), splitOn)
			if !_ok {
				continue
			}
			ok = true
			dr = append(dr, kv)
		}

	}
	return dr,ok
}

func (me *CommandResult) extractKeyValue(s string, splitOn rune) (ld KeyValue, ok bool) {

	foundKey := false
	stripSpace := false

	// splitting string by space but considering quoted section
	items := strings.FieldsFunc(s, func(c rune) bool {
		switch {
		case c == splitOn:
			if foundKey {
				return false
			} else {
				foundKey = true
				stripSpace = true
				return true
			}

		case unicode.IsSpace(c) && stripSpace:
			return true

		default:
			stripSpace = false
			return false
		}
	})

	switch len(items) {
	case 1:
		ld.Key = ""
		ld.Value = ""
		ok = true
	case 2:
		ld.Key = items[0]
		ld.Value = items[1]
		ok = true
	}

	ld.Key =   strings.Trim(strings.Trim(ld.Key,   `"`), `{`)
	ld.Value = strings.Trim(strings.Trim(ld.Value, `"`), `{`)

	return
}

func (me *CommandResult) GetStderrLn() string {
	return collapseNewLines(me.StdoutBuffer.String())
}

func (me *CommandResult) GetStdoutLn() string {
	return collapseNewLines(me.StdoutBuffer.String())
}

func collapseNewLines(s string) string {
	return strings.ReplaceAll(s, "\n", " ")
}

