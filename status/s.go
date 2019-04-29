package status

import (
	"encoding/json"
	"fmt"
)

var _ Status = (*S)(nil)

type S struct {
	success    bool
	cause      error
	httpstatus int
	message    string
	details    string
	data       interface{}
	help       HelpTypeMap
	errorcode  int
}

func (me *S) Json() []byte {
	js, _ := json.Marshal(&jsonS{
		Message: me.message,
		Help:    *me.help[ApiHelp],
		Data:    me.data,
	})
	return js
}

func (me *S) Error() string {
	return me.message
}

func (me *S) IsSuccess() bool {
	return me.success
}

func (me *S) IsError() bool {
	return !me.success
}

func (me *S) Cause() error {
	return me.cause
}

func (me *S) Message() string {
	return me.message
}

func (me *S) Detail() string {
	return me.details
}

func (me *S) Help() string {
	return me.GetHelp(AllHelp)
}

func (me *S) Data() (data interface{}) {
	return me.data
}

func (me *S) HttpStatus() int {
	return me.httpstatus
}

func (me *S) ErrorCode() int {
	return me.errorcode
}

func (me *S) FullError() (err error) {
	msg := me.message
	c := me.cause
	for {
		var ok bool
		c, ok = c.(error)
		if !ok {
			break
		}
		s := c.Error()
		if s != msg {
			msg = fmt.Sprintf("%s; %s", s, msg)
		}
		sts, ok := c.(Status)
		if !ok {
			break
		}
		c = sts.Cause()
	}
	return fmt.Errorf(msg)
}

func (me *S) SetSuccess(success bool) Status {
	me.success = success
	return me
}

func (me *S) SetMessage(msg string, args ...interface{}) Status {
	me.message = fmt.Sprintf(msg, args...)
	return me
}

func (me *S) SetDetail(details string, args ...interface{}) Status {
	me.details = fmt.Sprintf(details, args...)
	return me
}

func (me *S) SetHttpStatus(httpstatus int) Status {
	me.httpstatus = httpstatus
	return me
}

func (me *S) SetErrorCode(code int) Status {
	me.errorcode = code
	return me
}

func (me *S) SetData(data interface{}) Status {
	me.data = data
	return me
}

func (me *S) SetCause(err error) Status {
	me.cause = err
	return me
}

func (me *S) SetOtherHelp(help HelpTypeMap) Status {
	for t, h := range help {
		me.help[t] = h
	}
	return me
}
func (me *S) SetHelp(helptype HelpType, help string, args ...interface{}) Status {
	if len(args) > 0 {
		help = fmt.Sprintf(help, args...)
	}
	me.help[helptype] = &help
	if helptype == AllHelp {
		for t := range me.help {
			me.help[t] = &help
		}
	}
	return me
}

func (me *S) GetString() (s string, sts Status) {
	s, ok := me.Data().(string)
	if !ok {
		sts = Fail(&Args{
			Message: fmt.Sprintf("string expected for Status.Data; contains type '%T' instead",
				me.data,
			),
		})
	}
	return s, sts
}

func (me *S) GetHelp(helptype HelpType) string {
	h, _ := me.help[helptype]
	return *h
}

func (me *S) String() string {
	s := me.message
	if me.cause != nil && me.cause.Error() != me.message {
		s = fmt.Sprintf("%s: %s", s, me.cause.Error())
	}
	return s
}
