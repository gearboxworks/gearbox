package gearbox

import (
	"errors"
	"gearbox/api"
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

func NewOkStatus(msg ...string) *Status {
	var _msg string
	if len(msg) > 0 {
		_msg = msg[0]
	}
	return &Status{
		Success:    true,
		Message:    _msg,
		HttpStatus: http.StatusOK,
	}
}

func (me Status) Finalize() {
	if me.HttpStatus == 0 {
		me.Success = true
		me.HttpStatus = http.StatusOK
	}
}

func NewSuccessStatus(code int, msg ...string) *Status {
	s := NewOkStatus(msg...)
	s.HttpStatus = code
	return s
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

func (me *Status) IsError() bool {
	return me.Error != nil
}

func (me *Status) IsSuccess() bool {
	return me.Success
}
