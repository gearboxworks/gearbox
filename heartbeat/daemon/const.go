package daemon


// To be used with MyCaller()
const (
	callerStack0	= iota
	callerStack1	= iota
	CallerCurrent	= iota
	CallerParent	= iota
	CallerGrandParent = iota
	CallerGreatGrandParent = iota
)
// const CallerCurrent = 2
//const CallerParent = CallerCurrent + 1
//const CallerGrandParent = CallerParent + 1
//const CallerGrandParent = CallerParent + 1

