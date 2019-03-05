package gearbox

import (
	"errors"
	"fmt"
	"gearbox/api"
	"gearbox/only"
	"net/http"
)

var StatusInstance = (*Status)(nil)
var _ api.SuccessInspector = StatusInstance

var IsStatusError = errors.New("")

type Status struct {
	Success    bool
	Message    string `json:"message,omitempty"`
	Help       string `json:"-"`
	ApiHelp    string `json:"api_help,omitempty"`
	CliHelp    string `json:"-"`
	HttpStatus int    `json:"-"`
	Error      error  `json:"-"`
}
type StatusArgs Status

func NewOkStatus(msg string, args ...interface{}) Status {
	return Status{
		Success:    true,
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
				Success:    false,
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
				Success:    false,
				HttpStatus: http.StatusInternalServerError,
				Message:    m,
				Error:      errors.New(m),
				Help:       ContactSupportHelp(),
			}
			break
		}
		status = Status(*args)
		if status.Error == IsStatusError {
			status.Error = errors.New(status.Message)
		}
		if !status.Success && status.Error == nil {
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
		status.Success = status.Error == nil
	}
	return status
}

//
// Call just because you need to return an HTTP response
//
func (me Status) Finalize() {
	if me.HttpStatus == 0 {
		me.Success = true
		me.HttpStatus = http.StatusOK
	}
}

func (me *Status) IsError() bool {
	return me.Error != nil
}

func (me *Status) IsSuccess() bool {
	return me.Success || me.HttpStatus == 0
}

func (me *Status) NotYetFinalized() bool {
	return me.HttpStatus == 0
}
