package states

import (
	"time"
)

func (me *Status) SetNewState(new State, err error) bool {

	var ok bool

	if me.Current != new {
		me.Last = me.Current
		me.LastWhen = time.Now()
		me.Current = new
		ok = true
	}
	me.Error = err

	return ok
}


func (me *Status) SetNewWantState(new State) {

	me.Want = new
}

