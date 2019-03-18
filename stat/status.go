package stat

import (
	"errors"
	"fmt"
	"gearbox/only"
	"net/http"
)

var Instance = (*Status)(nil)
var _ SuccessInspector = Instance

var IsStatusError = errors.New("")

type SuccessInspector interface {
	IsSuccess() bool
}

type Status struct {
	Failed      bool
	Message     string      `json:"message,omitempty"`
	Help        string      `json:"-"`
	ApiHelp     string      `json:"help,omitempty"`
	CliHelp     string      `json:"-"`
	HttpStatus  int         `json:"-"`
	Err         error       `json:"-"`
	PriorStatus string      `json:"-"`
	Data        interface{} `json:"data,omitempty"`
}

type Args struct {
	Failed     bool
	Message    string
	Help       string
	ApiHelp    string
	CliHelp    string
	HttpStatus int
	Error      error
	Data       interface{}
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
			status = NewStatus(&Args{
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

func NewFailStatus(args *Args) (status Status) {
	args.Failed = true
	return NewStatus(args)
}

func NewErrorStatus(err error, args *Args) (status Status) {
	args.Error = err
	return NewFailStatus(args)
}

func NewStatus(args *Args) (status Status) {
	for range only.Once {

		status = Status{
			Failed:     args.Failed,
			Message:    args.Message,
			Help:       args.Help,
			ApiHelp:    args.ApiHelp,
			CliHelp:    args.CliHelp,
			HttpStatus: args.HttpStatus,
			Err:        args.Error,
		}

		if status.Err == IsStatusError {
			status.Err = errors.New(status.Message)
		}

		if status.Failed && status.Err == nil {
			status.Err = errors.New(status.Message)
		}

		if status.HttpStatus == 0 {
			status.HttpStatus = http.StatusInternalServerError
		}

		if status.Help == "" {
			status.Help = ContactSupportHelp()
		}

		if status.ApiHelp == "" {
			status.ApiHelp = status.Help
		}

		if status.CliHelp == "" {
			status.CliHelp = status.Help
		}

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

func (me Status) IsError() bool {
	return me.Err != nil
}

func (me Status) IsSuccess() bool {
	return !me.Failed || me.HttpStatus == 0
}

func (me Status) NotYetFinalized() bool {
	return me.HttpStatus == 0
}

func (me Status) String() string {
	s := me.Message
	if me.Err != nil && me.Err.Error() != me.Message {
		s = fmt.Sprintf("%s: %s", s, me.Err.Error())
	}
	return s
}

func (me Status) Error() string {
	return me.Message
}
