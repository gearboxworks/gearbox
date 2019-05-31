package daemon

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
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
// func WaitForSignal(name string) os.Signal {
func WaitForSignal() os.Signal {

	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)

	signal.Notify(signalChan, os.Kill, os.Interrupt)
	s := <-signalChan
	signal.Stop(signalChan)

	return s
}


// Wait for an ever increasing period of time - a very simple retry back-off system.
// This is used with processes that die too quickly and will ensure that retries don't hammer the system.
func WaitDelay(retry int) {

	// First time wait for 100mS
	// Second time wait for 200mS
	// And so on...
	time.Sleep(time.Millisecond * 100 * time.Duration(retry))
}


// This function will cause a Go() thread to sit and wait until
// a signal has been sent to the process.
// Very important for tidy up afterwards.
func WaitForTimeout(wt time.Duration) bool {

	var exitState bool

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// Timeout timer.
	var tc <-chan time.Time
	if wt > 0 {
		tc = time.After(wt)
	}

	select {
		case <-sig:
			exitState = false
			// Exit by user
		case <-tc:
			exitState = true
			// Exit by timeout
	}

	return exitState
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

