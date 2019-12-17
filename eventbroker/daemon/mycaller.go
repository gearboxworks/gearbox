package daemon

import "runtime"

// To be used with MyCaller()
const (
	callerStack0           = iota
	callerStack1           = iota
	CallerCurrent          = iota
	CallerParent           = iota
	CallerGrandParent      = iota
	CallerGreatGrandParent = iota
)


// Determine the calling functions that called this function.
// IE: MyCaller's grand-parent.
func MyCaller(whichCaller int) (fileName string, lineNumber int) {

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

	_, _, lineNumber, _ = runtime.Caller(whichCaller)

	fileName = fun.Name()

	return
}

