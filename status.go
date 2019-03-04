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
	Success    bool   `json:"-"`
	Message    string `json:"message,omitempty"`
	Help       string `json:"-"`
	ApiHelp    string `json:"api_help,omitempty"`
	CliHelp    string `json:"-"`
	HttpStatus int    `json:"-"`
	Error      error  `json:"-"`
}
type StatusArgs Status

func NewOkStatus(msg string, args ...interface{}) *Status {
	return &Status{
		Success:    true,
		Message:    fmt.Sprintf(msg, args...),
		HttpStatus: http.StatusOK,
	}
}

func NewSuccessStatus(code int, msg ...string) (status *Status) {
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
				Help:       "contact support",
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

func NewStatus(args *StatusArgs) *Status {
	s := Status(*args)
	if s.Error == IsStatusError {
		s.Error = errors.New(s.Message)
	}
	if !s.Success && s.Error == nil {
		s.Error = errors.New(s.Message)
	}
	if s.Help != "" {
		if s.ApiHelp == "" {
			s.ApiHelp = s.Help
		}
		if s.CliHelp == "" {
			s.CliHelp = s.Help
		}
	}
	s.Success = s.Error == nil
	return &s
}

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
	return me.Success
}
