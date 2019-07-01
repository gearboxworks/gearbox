package eblog

import (
	"fmt"
)

const (
	callerStack0           = iota
	callerStack1           = iota
	CurrentCaller          = iota
	ParentCaller           = iota
	GrandParentCaller      = iota
	GreatGrandParentCaller = iota
)

type Callers []Caller
type Caller struct {
	File       string
	LineNumber int
	Function   string
}

// Determine the calling functions that called this function.
// IE: MyCaller's grand-parent.
func (me *Callers) String() string {

	var ret string

	if me == nil {
		return ""
	}

	for k, v := range *me {
		ret += fmt.Sprintf("[%d] %s %s:%d\n",
			k,
			v.File,
			v.Function,
			v.LineNumber,
		)
	}

	return ret
}
