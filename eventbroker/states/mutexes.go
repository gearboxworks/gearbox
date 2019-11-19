package states

import (
	"fmt"
	"github.com/gearboxworks/go-status/only"
	"time"
)


func (me *Status) SetNewState(new State, err error) bool {

	var ok bool

	if new == "" {
		return ok
	}

	if me.mutex != nil {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}

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

	if me.mutex != nil {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}

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


func (me *Status) String() string {

	var ret string
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if me.mutex != nil {
			me.mutex.RLock()
			defer me.mutex.RUnlock()
		}

		ret = fmt.Sprintf(`EntityId:%s  Name:%s  Parent:%s  Action:%s  Want:%s  Current:%s  Last:%s  LastWhen:%v  Error:%v`,
			me.EntityId.String(),
			me.EntityName.String(),
			me.ParentId.String(),
			me.Action.String(),
			me.Want.String(),
			me.Current.String(),
			me.Last.String(),
			me.LastWhen.Unix(),
			me.Error,
		)
	}

	return ret
}


func (me *Status) ShortString() string {

	var ret string
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if me.mutex != nil {
			me.mutex.RLock()
			defer me.mutex.RUnlock()
		}

		ret = fmt.Sprintf("Name:%s\tCurrent:%s\tParent:%s\tAction:%s\tLast:%s\tLastWhen:%v\tError:%v",
			me.EntityName.String(),
			me.Current.String(),
			me.ParentId.String(),
			me.Action.String(),
			me.Last.String(),
			me.LastWhen.Unix(),
			me.Error,
		)
	}

	return ret
}


func (me *Status) IsTheSame(other Status) bool {

	var ret bool
	var err error

	for range only.Once {
		err = me.EnsureNotNil()
		if err != nil {
			break
		}

		if me.mutex != nil {
			me.mutex.RLock()
			defer me.mutex.RUnlock()
		}

		if me.EntityId == nil {
			break
		}
		if me.EntityName == nil {
			break
		}
		if me.ParentId == nil {
			break
		}

		if other.EntityId == nil {
			break
		}
		if other.EntityName == nil {
			break
		}
		if other.ParentId == nil {
			break
		}

		if *me.EntityId != *other.EntityId {
			break
		}
		if *me.EntityName != *other.EntityName {
			break
		}
		if *me.ParentId != *other.ParentId {
			break
		}
		if me.Current != other.Current {
			break
		}
		if me.Want != other.Want {
			break
		}
		if me.Last != other.Last {
			break
		}
		if me.LastWhen != other.LastWhen {
			break
		}
		if me.Attempts != other.Attempts {
			break
		}
		if me.Error != other.Error {
			break
		}
		if me.Action != other.Action {
			break
		}

		ret = true
	}

	return ret
}


func StatusAsString(me *Status) string {

	return me.String()
}


func (me *Status) GetStatus() *Status {

	if me.mutex != nil {
		me.mutex.RLock()
		defer me.mutex.RUnlock()
	}

	return me
}

func (me *Status) HasChangedState() bool {

	if me.mutex != nil {
		me.mutex.RLock()
		defer me.mutex.RUnlock()
	}

	if me.Last != me.Current {
		return true
	} else {
		return false
	}
}

func (me *Status) ExpectingNewState() bool {

	if me.mutex != nil {
		me.mutex.RLock()
		defer me.mutex.RUnlock()
	}

	if me.Want != me.Current {
		return true
	} else {
		return false
	}
}

func (me *Status) GetCurrent() State {

	if me.mutex != nil {
		me.mutex.RLock()
		defer me.mutex.RUnlock()
	}

	return me.Current
}

func (me *Status) GetWant() State {

	if me.mutex != nil {
		me.mutex.RLock()
		defer me.mutex.RUnlock()
	}

	return me.Want
}

func (me *Status) GetLast() State {

	if me.mutex != nil {
		me.mutex.RLock()
		defer me.mutex.RUnlock()
	}

	return me.Last
}

func (me *Status) GetLastWhen() time.Time {

	if me.mutex != nil {
		me.mutex.RLock()
		defer me.mutex.RUnlock()
	}

	return me.LastWhen
}

func (me *Status) GetAttempts() int {

	if me.mutex != nil {
		me.mutex.RLock()
		defer me.mutex.RUnlock()
	}

	return me.Attempts
}

func (me *Status) GetError() error {

	if me.mutex != nil {
		me.mutex.RLock()
		defer me.mutex.RUnlock()
	}

	return me.Error
}

func (me *Status) GetAction() Action {

	if me.mutex != nil {
		me.mutex.RLock()
		defer me.mutex.RUnlock()
	}

	return me.Action
}

func (me *Status) SetCurrent(s State) {

	if me.mutex != nil {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}

	me.Current = s
}

func (me *Status) SetWant(s State) {

	if me.mutex != nil {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}

	me.Want = s
}

func (me *Status) SetLast(s State) {

	if me.mutex != nil {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}

	me.Last = s
}

func (me *Status) SetLastWhen(t time.Time) {

	if me.mutex != nil {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}

	me.LastWhen = t
}

func (me *Status) SetAttempts(a int) {

	if me.mutex != nil {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}

	me.Attempts = a
}

func (me *Status) AddAttempts() {

	if me.mutex != nil {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}

	me.Attempts++
}

func (me *Status) ZeroAttempts() {

	if me.mutex != nil {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}

	me.Attempts = 0
}

func (me *Status) SetError(e error) {

	if me.mutex != nil {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}

	me.Error = e
}

func (me *Status) SetaAction(a Action) {

	if me.mutex != nil {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}

	me.Action = a
}
