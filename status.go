package gearbox

import (
	"errors"
	"net/http"
)

var IsStatusError = errors.New("")

type Status struct {
	Success    bool   `json:"success"`
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

func NewSuccessStatus(code int) *Status {
	return &Status{
		Success:    true,
		HttpStatus: code,
	}
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
