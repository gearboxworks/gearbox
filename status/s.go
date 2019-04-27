package status

import (
	"encoding/json"
	"fmt"
)

var _ Status = (*S)(nil)

type S struct {
	success    bool
	cause      error
	httpStatus int
	message    string
	data       interface{}
	help       HelpTypeMap
}

func (me *S) Json() []byte {
	js, _ := json.Marshal(&jsonS{
		Message: me.message,
		Help:    me.help[ApiHelp],
		Data:    me.data,
	})
	return js
}

func (me *S) Error() string {
	return me.message
}

func (me *S) Cause() error {
	return me.cause
}

func (me *S) Message() string {
	return me.message
}

func (me *S) SetSuccess(success bool) {
	me.success = success
}

func (me *S) SetMessage(msg string) {
	me.message = msg
}

func (me *S) SetHttpStatus(code int) {
	me.httpStatus = code
}

func (me *S) SetData(data interface{}) {
	me.data = data
}

func (me *S) GetData() (data interface{}) {
	return me.data
}

func (me *S) GetString() (s string, sts Status) {
	s, ok := me.GetData().(string)
	if !ok {
		sts = Fail(&Args{
			Message: fmt.Sprintf("string expected for Status.Data; contains type '%T' instead",
				me.data,
			),
		})
	}
	return s, sts
}

func (me *S) SetCause(err error) {
	me.cause = err
}

func (me *S) SetOtherHelp(help HelpTypeMap) {
	for t, h := range help {
		me.help[t] = h
	}
}
func (me *S) SetHelp(helptype HelpType, help string) {
	me.help[helptype] = help
}

func (me *S) GetHelp(helptype HelpType) (h string) {
	h, _ = me.help[helptype]
	return h
}

func (me *S) String() string {
	s := me.message
	if me.cause != nil && me.cause.Error() != me.message {
		s = fmt.Sprintf("%s: %s", s, me.cause.Error())
	}
	return s
}
func (me *S) HttpStatus() int {
	return me.httpStatus
}
func (me *S) IsSuccess() bool {
	return me.success
}
func (me *S) IsError() bool {
	return !me.success
}

func (me *S) GetFullError() (err error) {
	msg := me.message
	c := me.cause
	for {
		var ok bool
		c, ok = c.(error)
		if !ok {
			break
		}
		msg = fmt.Sprintf("%s; %s", c.Error(), msg)
		sts, ok := c.(Status)
		if !ok {
			break
		}
		c = sts.Cause()
	}
	return fmt.Errorf(msg)
}
