package logger

import (
	"github.com/gearboxworks/go-status"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)


func (me *Logger) printLog(level logrus.Level, fileName string, lineNumber int, textString string, opt ...interface{}) (returnCode bool) {

	returnCode = true

	textString = strings.TrimSuffix(textString, "\n")
	fields := logrus.Fields{
//		"tenant_name": notifyConfig.TenantName,
//		"tenant_guid": notifyConfig.TenantGUID,
//		"process_guid": notifyConfig.ProcessGUID,
		"filename": fileName,
		"line": lineNumber}

	switch {
		case level == DebugLevel:
			me.logrusInstance.WithFields(fields).Debugf(textString, opt...)

		case level == InfoLevel:
			me.logrusInstance.WithFields(fields).Infof(textString, opt...)

		case level == WarnLevel:
			me.logrusInstance.WithFields(fields).Warnf(textString, opt...)

		case level == ErrorLevel:
			me.logrusInstance.WithFields(fields).Errorf(textString, opt...)

		case level == FatalLevel:
			me.logrusInstance.WithFields(fields).Fatalf(textString, opt...)

		case level == PanicLevel:
			me.logrusInstance.WithFields(fields).Panicf(textString, opt...)
	}

	returnCode = false

	return
}

func Debug(format string, a ...interface{}) {
	//fn, ln := daemon.MyCaller(daemon.CallerParent)
	//status.Success(fmt.Sprintf("DEBUG '%s':[%d] ", fn, ln) + format, a).Log()
	status.Success("DEBUG: " + format, a...).Log()
}

// To be used with MyCaller()
const (
	callerStack0	= iota
	callerStack1	= iota
	CallerCurrent	= iota
	CallerParent	= iota
	CallerGrandParent = iota
	CallerGreatGrandParent = iota
)

//// Determine the calling functions that called this function.
//// IE: MyCaller's grand-parent.
//func MyCaller(whichCaller int) (fileName string, lineNumber int) {
//
//	fileName = "unknown"
//
//	if whichCaller == 0 {
//		whichCaller = CallerParent
//	}
//
//	// we get the callers as uintptrs - but we just need 1
//	fpcs := make([]uintptr, 1)
//
//	// skip 3 levels to get to the caller of whoever called Caller()
//	n := runtime.Callers(whichCaller, fpcs)
//	if n == 0 {
//		return          // Proper error handling would be better here.
//	}
//
//	// get the info of the actual function that's in the pointer
//	fun := runtime.FuncForPC(fpcs[0]-1)
//	if fun == nil {
//		return
//	}
//
//	_, _, lineNumber, _ = runtime.Caller(whichCaller)
//
//	fileName = fun.Name()
//
//	return
//}

// Determine the calling functions that called this function.
// IE: MyCaller's grand-parent.
func MyCaller(whichCaller int) (string, int) {

	if whichCaller == 0 {
		whichCaller = CallerParent
	}

	pc, _, _, _ := runtime.Caller(whichCaller)
	e := runtime.FuncForPC(pc)
	fn := e.Name()
	_, ln := e.FileLine(pc)

	return fn, ln
}


type Caller struct {
	File string
	LineNumber int
	Function string
}
type Callers []Caller

// Determine the calling functions that called this function.
// IE: MyCaller's grand-parent.
func MyCallers(whichCaller int, howMany int) (*Callers) {


	if whichCaller == 0 {
		whichCaller = CallerParent
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

