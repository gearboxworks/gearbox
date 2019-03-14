package util

import (
	"fmt"
	"os"
	"strings"
)

//func NewHelpfulError(msg, help string) HelpfulError {
//	return HelpfulError{
//		error: errors.New(msg),
//		Help:  help,
//	}
//}

type HelpfulError struct {
	ErrorObj error
	Help     string
}

func AddHelpToError(err error, help string) HelpfulError {
	return HelpfulError{
		ErrorObj: err,
		Help:     help,
	}
}

func (me HelpfulError) IsNil() bool {
	return me.ErrorObj == nil
}

func (me HelpfulError) Error() string {
	return me.ErrorObj.Error()
}

func Error(msg interface{}, args ...interface{}) {
	var err error
	_msg, ok := msg.(string)
	if !ok {
		err, ok = msg.(error)
		if ok {
			_msg = err.Error()
		}
	}
	if !ok {
		panic(err)
	}
	if len(args) > 0 {
		_msg = fmt.Sprintf(_msg, args...)
	}
	parts := strings.Split(_msg, " ")
	_msg = strings.Title(parts[0])
	if len(parts) > 1 {
		_msg += " " + strings.Join(parts[1:], " ")
	}
	fmt.Println(_msg + ".")
	os.Exit(1)
}
