package eblog

import (
	"fmt"
	"gearbox/heartbeat/eventbroker/messages"
	"github.com/gearboxworks/go-status"
	"reflect"
	"runtime"
	"strconv"
)


// To be used with MyCaller()
const (
	callerStack0	= iota
	callerStack1	= iota
	CallerCurrent	= iota
	CallerParent	= iota
	CallerGrandParent = iota
	CallerGreatGrandParent = iota
)

// Determine the calling functions that called this function.
// IE: MyCaller's grand-parent.
func MyCallerExtended(whichCaller int) (fileName string, lineNumber int) {

	fileName = "unknown"

	if whichCaller == 0 {
		whichCaller = CallerParent
	}

	// we get the callers as uintptrs - but we just need 1
	fpcs := make([]uintptr, 1)

	// skip 3 levels to get to the caller of whoever called Caller()
	n := runtime.Callers(whichCaller, fpcs)
	if n == 0 {
		return          // Proper error handling would be better here.
	}

	// get the info of the actual function that's in the pointer
	fun := runtime.FuncForPC(fpcs[0]-1)
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

// Determine the calling functions that called this function.
// IE: MyCaller's grand-parent.
func (me *Callers) Print() string {

	var ret string

	if me == nil {
		return ""
	}

	for k, v := range *me {
		ret += fmt.Sprintf("[%d] %s %s:%d\n", k, v.File, v.Function, v.LineNumber)
	}

	return ret
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

	rn := true

	switch {
		case reflect.ValueOf(i).String() == "":

		case i == nil:
			fallthrough
		case reflect.ValueOf(i).IsNil():
			rn = false

			callers := " NIL:["
			// Fetch last two callers.
			for _, d := range *MyCallers(CallerParent, howMany) {
				callers += " <- " + d.Function + ":" + strconv.Itoa(d.LineNumber)
			}
			callers += "] "

			status.Success("nil interface" + callers).Log()
		}

	return rn
}


func Debug(client messages.MessageAddress, format string, a ...interface{}) {
	//fn, ln := daemon.MyCaller(daemon.CallerParent)
	//status.Success(fmt.Sprintf("DEBUG '%s':[%d] ", fn, ln) + format, a).Log()
	status.Success(string(client) + ": " + format, a...).Log()
}


const SkipNilCheck = ""
const howMany = 2
// Check for a nil type or err and log.
func LogIfError(address messages.MessageAddress, err error, format ...interface{}) bool {

	rn := true
	if err != nil {
		callers := "%s ERROR:["
		// callers := address.String() + " ERROR:%s ["
		// Fetch last two callers.
		for _, d := range *MyCallers(CallerParent, howMany) {
			callers += " <- " + d.Function + ":" + strconv.Itoa(d.LineNumber)
		}
		callers += "] "

		if len(format) == 0 {
			status.Success(callers, err).Log()
			//fmt.Printf(callers + "\n", err)
		} else {
			status.Success(format[0].(string) + callers, format[1:]...).Log()
			//fmt.Printf(format[0].(string) + callers + "\n", format[1:]...)
		}
	}

	return rn
}

/*
	if eblog.LogIfError(me, err) {
		// Save last state.
		me.State.Error = err
	}
*/

