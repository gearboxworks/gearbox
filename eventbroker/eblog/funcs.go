package eblog

import (
	"gearbox/eventbroker/msgs"
	"github.com/sirupsen/logrus"
	"reflect"
	"runtime"
)

// Determine the calling functions that called this function.
// IE: MyCaller's grand-parent.
func MyCallers(whichCaller int, howMany int) *Callers {

	if whichCaller == 0 {
		whichCaller = ParentCaller
	}

	if howMany == 0 {
		howMany = 2
	}

	pc := make([]uintptr, howMany)
	count := runtime.Callers(whichCaller, pc)

	callers := make(Callers, count)
	for i, d := range pc {
		e := runtime.FuncForPC(d)
		callers[i].Function = e.Name()
		callers[i].File, callers[i].LineNumber = e.FileLine(d)
	}

	return &callers
}

// Determine the calling functions that called this function.
// IE: MyCaller's grand-parent.
func MyCallerExtended(whichCaller int) (fileName string, lineNumber int) {

	fileName = "unknown"

	if whichCaller == 0 {
		whichCaller = ParentCaller
	}

	// we get the callers as uintptrs - but we just need 1
	fpcs := make([]uintptr, 1)

	// skip 3 levels to get to the caller of whoever called Caller()
	n := runtime.Callers(whichCaller, fpcs)
	if n == 0 {
		return // Proper error handling would be better here.
	}

	// get the info of the actual function that's in the pointer
	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return
	}

	_, fileName, lineNumber, _ = runtime.Caller(whichCaller)

	fileName = fun.Name()

	return
}

// Determine the calling functions that called this function.
// IE: MyCaller's grand-parent.
func MyCaller(whichCaller int) (string, int) {

	pc, _, _, _ := runtime.Caller(whichCaller)
	e := runtime.FuncForPC(pc)
	fn := e.Name()
	_, ln := e.FileLine(pc)

	return fn, ln
}

// Check for a nil type.
func IsNil(i interface{}) bool {
	defer func() { recover() }()
	if i == nil || reflect.ValueOf(i).IsNil() {
		// It's a nil type.
		return true
	} else {
		// It's not a nil type.
		return false
	}
}

// Check for a nil type.
func IsNotNil(i interface{}) bool {
	defer func() { recover() }()
	if i == nil || reflect.ValueOf(i).IsNil() {
		// It's a nil type.
		return false
	} else {
		// It's not a nil type.
		return true
	}
}

// Check for a nil type.
func LogIfNil(i interface{}, format ...interface{}) bool {

	var ret bool

	switch {
	case reflect.ValueOf(i).String() == "":

	case i == nil:
		fallthrough
	case reflect.ValueOf(i).IsNil():
		ret = true
		localLogger.printLog(logrus.ErrorLevel, "nil interface")
	}

	return ret
}

// Check for a nil type or err and log.
func LogIfError(err error, format ...interface{}) bool {

	var ret bool

	if err != nil {
		ret = true

		if len(format) == 0 {
			localLogger.printLog(logrus.ErrorLevel, "%v", err)
		} else {
			localLogger.printLog(logrus.ErrorLevel, format[0].(string), format[1:]...)
		}
	}

	return ret
}

func Debug(client msgs.Address, s string, args ...interface{}) {
	if localLogger == nil {
		return
	}
	localLogger.printLog(logrus.DebugLevel, s, args...)
}
