package states

import "time"


func (me *Status) SetNewState(new State, err error) bool {

	var ok bool
	me.mutex.Lock()
	defer me.mutex.Unlock()

	if err == nil {
		if me.Current != new {
			me.Last = me.Current
			me.LastWhen = time.Now()
			me.Current = new

			ok = true
		}
	}
	me.Error = err

	return ok
}


func (me *Status) SetNewAction(a Action) bool {

	var ok bool
	me.mutex.Lock()
	defer me.mutex.Unlock()

	me.Action = a
	switch me.Action {
	case ActionRegister:
		me.Current = StateRegistering
		me.Want = StateRegistered
	case ActionUnregister:
		me.Current = StateUnregistering
		me.Want = StateUnregistered
	case ActionPublish:
		me.Current = StatePublishing
		me.Want = StatePublished
	case ActionUnpublish:
		me.Current = StateUnpublishing
		me.Want = StateUnpublished
	case ActionSubscribe:
		me.Current = StateSubscribing
		me.Want = StateSubscribed
	case ActionUnsubscribe:
		me.Current = StateUnsubscribing
		me.Want = StateUnsubscribed


	case ActionInitialize:
		me.Current = StateInitializing
		me.Want = StateInitialized
	case ActionStop:
		me.Current = StateStopping
		me.Want = StateStopped
	case ActionStart:
		me.Current = StateStarting
		me.Want = StateStarted
	}

	return ok
}


func (me *Status) GetStatus() *Status {

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	return me
}

func (me *Status) HasChangedState() bool {

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	if me.Last != me.Current {
		return true
	} else {
		return false
	}
}

func (me *Status) ExpectingNewState() bool {

	me.mutex.RLock()
	defer me.mutex.RUnlock()

	if me.Want != me.Current {
		return true
	} else {
		return false
	}
}

func (me *Status) GetCurrent() State {

	me.mutex.RLock()
	defer me.mutex.RUnlock()
	return me.Current
}

func (me *Status) GetWant() State {

	me.mutex.RLock()
	defer me.mutex.RUnlock()
	return me.Want
}

func (me *Status) GetLast() State {

	me.mutex.RLock()
	defer me.mutex.RUnlock()
	return me.Last
}

func (me *Status) GetLastWhen() time.Time {

	me.mutex.RLock()
	defer me.mutex.RUnlock()
	return me.LastWhen
}

func (me *Status) GetAttempts() int {

	me.mutex.RLock()
	defer me.mutex.RUnlock()
	return me.Attempts
}

func (me *Status) GetError() error {

	me.mutex.RLock()
	defer me.mutex.RUnlock()
	return me.Error
}

func (me *Status) GetAction() Action {

	me.mutex.RLock()
	defer me.mutex.RUnlock()
	return me.Action
}

func (me *Status) SetCurrent(s State) {

	me.mutex.Lock()
	defer me.mutex.Unlock()
	me.Current = s
}

func (me *Status) SetWant(s State) {

	me.mutex.Lock()
	defer me.mutex.Unlock()
	me.Want = s
}

func (me *Status) SetLast(s State) {

	me.mutex.Lock()
	defer me.mutex.Unlock()
	me.Last = s
}

func (me *Status) SetLastWhen(t time.Time) {

	me.mutex.Lock()
	defer me.mutex.Unlock()
	me.LastWhen = t
}

func (me *Status) SetAttempts(a int) {

	me.mutex.Lock()
	defer me.mutex.Unlock()
	me.Attempts = a
}

func (me *Status) AddAttempts() {

	me.mutex.Lock()
	defer me.mutex.Unlock()
	me.Attempts++
}

func (me *Status) ZeroAttempts() {

	me.mutex.Lock()
	defer me.mutex.Unlock()
	me.Attempts = 0
}

func (me *Status) SetError(e error) {

	me.mutex.Lock()
	defer me.mutex.Unlock()
	me.Error = e
}

func (me *Status) SetaAction(a Action) {

	me.mutex.Lock()
	defer me.mutex.Unlock()
	me.Action = a
}
