package gearbox

import (
	"errors"
	"fmt"
	"gearbox/api"
	"gearbox/only"
	"gearbox/util"
	"net/http"
)

var StatusInstance = (*Status)(nil)
var _ api.SuccessInspector = StatusInstance

var IsStatusError = errors.New("")

type Status struct {
	Failed     bool
	Message    string `json:"message,omitempty"`
	Help       string `json:"-"`
	ApiHelp    string `json:"api_help,omitempty"`
	CliHelp    string `json:"-"`
	HttpStatus int    `json:"-"`
	Error      error  `json:"-"`
}

type StatusArgs struct {
	Failed     bool
	Message    string
	Help       string
	ApiHelp    string
	CliHelp    string
	HttpStatus int
	Error      error
	util.HelpfulError
}

func NewOkStatus(msg string, args ...interface{}) Status {
	return Status{
		Failed:     false,
		Message:    fmt.Sprintf(msg, args...),
		HttpStatus: http.StatusOK,
	}
}
func ContactSupportHelp() string {
	return "contact support"
}

func NewSuccessStatus(code int, msg ...string) (status Status) {
	for range only.Once {
		if len(msg) == 0 {
			m := fmt.Sprintf("NewSuccessStatus(%d) called with no msg parameter",
				code,
			)
			status = NewStatus(&StatusArgs{
				Failed:     true,
				HttpStatus: http.StatusInternalServerError,
				Message:    m,
				Error:      errors.New(m),
				Help:       ContactSupportHelp(),
			})
			break
		}
		msg = msg[1:]
		is := make([]interface{}, len(msg))
		for i, m := range msg {
			is[i] = m
		}
		status := NewOkStatus(msg[0], is...)
		status.HttpStatus = code
	}
	return status
}

func NewStatus(args *StatusArgs) (status Status) {
	for range only.Once {
		if args.HttpStatus == 0 {
			m := fmt.Sprintf("NewStatus() called with no HttpStatus for %s",
				args.Message,
			)
			status = Status{
				Failed:     true,
				HttpStatus: http.StatusInternalServerError,
				Message:    m,
				Error:      errors.New(m),
				Help:       ContactSupportHelp(),
			}
			break
		}

		status = Status{
			Failed:     args.Failed,
			Message:    args.Message,
			Help:       args.Help,
			ApiHelp:    args.ApiHelp,
			CliHelp:    args.CliHelp,
			HttpStatus: args.HttpStatus,
			Error:      args.Error,
		}

		if status.Error == nil {
			status.Error = args.HelpfulError
		}
		if status.Help == "" && args.HelpfulError.Help != "" {
			status.Help = args.HelpfulError.Help
		}
		if status.Error == IsStatusError {
			status.Error = errors.New(status.Message)
		}
		if status.Failed && status.Error == nil {
			status.Error = errors.New(status.Message)
		}
		if status.Help != "" {
			if status.ApiHelp == "" {
				status.ApiHelp = status.Help
			}
			if status.CliHelp == "" {
				status.CliHelp = status.Help
			}
		}
		status.Failed = status.Error != nil
	}
	return status
}

//
// Call just because you need to return an HTTP response
//
func (me Status) Finalize() {
	if me.HttpStatus == 0 {
		me.Failed = false
		me.HttpStatus = http.StatusOK
	}
}

func (me *Status) IsError() bool {
	return me.Error != nil
}

func (me *Status) IsSuccess() bool {
	return !me.Failed || me.HttpStatus == 0
}

func (me *Status) NotYetFinalized() bool {
	return me.HttpStatus == 0
}
