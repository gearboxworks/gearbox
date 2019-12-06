package states

import (
	"encoding/json"
	"errors"
	"fmt"
	"gearbox/eventbroker/entity"
	"gearbox/eventbroker/msgs"
	"github.com/gearboxworks/go-status/only"
	"sync"
	"time"
)

type Status struct {
	EntityId   msgs.Address
	EntityName msgs.Address
	ParentId   msgs.Address
	Current    State
	Want       State
	Last       State
	LastWhen   time.Time
	Attempts   int
	Error      error
	Action     Action

	mutex *sync.RWMutex // Mutex control for map.
}

func New(client msgs.Address, name msgs.Address, parent msgs.Address) *Status {

	var ret Status

	if name == "" || name == entity.SelfEntityName {
		name = client
	}

	ret = Status{
		EntityId:   client,
		EntityName: name,
		ParentId:   parent,
		Action:     ActionIdle,
		Want:       StateIdle,
		Current:    StateIdle,
		Last:       StateIdle,
		LastWhen:   time.Now(),
		Error:      nil,
		mutex:      &sync.RWMutex{},
	}

	return &ret
}

func (me *Status) EnsureNotNil() error {

	var err error

	switch {
	case me == nil:
		err = errors.New("status.Status is nil")

	//case me.mutex == nil:
	//	err = errors.New("status.mutex is nil")

	case me.EntityId == "":
		err = errors.New("status.EntityId is empty")
	}

	return err
}

func (me *Status) ToMessageText() msgs.Text {

	j := []byte("{}")

	for range only.Once {

		err := me.EnsureNotNil()
		if err != nil {
			break
		}

		j, err = json.Marshal(me)
		if err != nil {

			break
		}
	}

	return msgs.Text(j)
}

///////////[ Mutexs ]//////////////

func (me *Status) GetStatus() *Status {

	if me.mutex != nil {
		me.mutex.RLock()
		defer me.mutex.RUnlock()
	}

	return me
}

func (me *Status) SetState(new State) bool {

	var changed bool

	if new == "" {
		return false
	}

	if me.mutex != nil {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}

	if me.Current != new {
		me.Last = me.Current
		me.LastWhen = time.Now()
		me.Current = new
		changed = true
	}

	return changed
}

func (me *Status) SetNewState(new State, err error) bool {

	var changed bool

	if new == "" {
		return changed
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
			changed = true
		}
	}
	me.Error = err

	return changed
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
			me.Action,
			me.Want,
			me.Current,
			me.Last,
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
			me.Current,
			me.ParentId.String(),
			me.Action,
			me.Last,
			me.LastWhen.Unix(),
			me.Error,
		)
	}

	return ret
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

func (me *Status) IsEqualTo(other Status) bool {

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

		// @TODO For the next six ifs, shouldn't these throw an error?
		//
		if me.EntityId == "" {
			break
		}
		if me.EntityName == "" {
			break
		}
		if me.ParentId == "" {
			break
		}

		if other.EntityId == "" {
			break
		}
		if other.EntityName == "" {
			break
		}
		if other.ParentId == "" {
			break
		}

		if me.EntityId != other.EntityId {
			break
		}
		if me.EntityName != other.EntityName {
			break
		}
		if me.ParentId != other.ParentId {
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

func (me *Status) ZeroAttempts() {
	me.SetAttempts(0)
}

func (me *Status) GetCurrent() State {

	if me.mutex != nil {
		me.mutex.RLock()
		defer me.mutex.RUnlock()
	}

	return me.Current
}

func (me *Status) SetCurrent(s State) {

	if me.mutex != nil {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}

	me.Current = s
}

func (me *Status) GetLastWhen() time.Time {

	if me.mutex != nil {
		me.mutex.RLock()
		defer me.mutex.RUnlock()
	}

	return me.LastWhen
}

func (me *Status) SetLastWhen(t time.Time) {

	if me.mutex != nil {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}

	me.LastWhen = t
}

func (me *Status) GetAction() Action {

	if me.mutex != nil {
		me.mutex.RLock()
		defer me.mutex.RUnlock()
	}

	return me.Action
}

func (me *Status) SetAction(a Action) {

	if me.mutex != nil {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}

	me.Action = a
}

func (me *Status) GetLast() State {

	if me.mutex != nil {
		me.mutex.RLock()
		defer me.mutex.RUnlock()
	}

	return me.Last
}

func (me *Status) SetLast(s State) {

	if me.mutex != nil {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}

	me.Last = s
}

func (me *Status) GetError() error {

	if me.mutex != nil {
		me.mutex.RLock()
		defer me.mutex.RUnlock()
	}

	return me.Error
}

func (me *Status) SetError(e error) {

	if me.mutex != nil {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}

	me.Error = e
}

func (me *Status) GetWant() State {

	if me.mutex != nil {
		me.mutex.RLock()
		defer me.mutex.RUnlock()
	}

	return me.Want
}

func (me *Status) SetWant(s State) {

	if me.mutex != nil {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}

	me.Want = s
}

func (me *Status) GetAttempts() int {

	if me.mutex != nil {
		me.mutex.RLock()
		defer me.mutex.RUnlock()
	}

	return me.Attempts
}

func (me *Status) SetAttempts(a int) {

	if me.mutex != nil {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}

	me.Attempts = a
}

func (me *Status) AddAttempt() {

	if me.mutex != nil {
		me.mutex.Lock()
		defer me.mutex.Unlock()
	}

	me.Attempts++
}
