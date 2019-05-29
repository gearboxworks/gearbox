package daemon

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"
)

func IsParentInit() (bool) {

	ppid := os.Getppid()
	if ppid == 1 {
		return true
	}

	return false
}


// This function will cause a Go() thread to sit and wait until
// a signal has been sent to the process.
// Very important for tidy up afterwards.
func WaitForSignal() os.Signal {
	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)

	signal.Notify(signalChan, os.Kill, os.Interrupt)
	s := <-signalChan
	signal.Stop(signalChan)

	return s
}


// This function will cause a Go() thread to sit and wait until
// a signal has been sent to the process.
// Very important for tidy up afterwards.
func SimpleWaitLoop(t string, i int, d time.Duration) {

	for iterate := 0; iterate < i; i++ {
		fmt.Printf("> Wait: %s\n", t)
		time.Sleep(d)
	}

	return
}


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

