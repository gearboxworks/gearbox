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
		_msg := msg[0]
		var params []interface{}
		if len(msg) > 1 {
			params = make([]interface{}, len(msg)-1)
			for i, param := range msg[1:] {
				params[i] = param
			}
		}
		status = NewOkStatus(_msg, params...)
		status.HttpStatus = code
	}
	return status
}

func NewStatusFromHelpfulError(helpfulError error, args *StatusArgs) (status Status) {
	var help string
	var err error
	if he, ok := helpfulError.(util.HelpfulError); ok {
		help = he.Help
		err = he.ErrorObj
	}
	return NewStatus(&StatusArgs{
		Failed:     true,
		Message:    err.Error(),
		HttpStatus: args.HttpStatus,
		Help:       help,
		Error:      err,
	})
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

		if status.Error == nil && !args.HelpfulError.IsNil() {
			status.Error = args.HelpfulError
		}
		if status.Error == IsStatusError {
			status.Error = errors.New(status.Message)
		}
		he, ok := status.Error.(util.HelpfulError)
		if status.Error == nil && ok && he.IsNil() {
			status.Error = errors.New(status.Message)
		}
		if status.Error == nil && status.Failed {
			status.Error = errors.New(status.Message)
		}

		if status.Help == "" && args.HelpfulError.Help != "" {
			status.Help = args.HelpfulError.Help
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
