package util

import (
	"fmt"
	"os"
	"strings"
)

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
