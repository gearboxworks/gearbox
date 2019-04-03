package util

import (
	"encoding/json"
	"fmt"
	"gearbox/only"
)

type index map[string]int
type keys []string
type values []interface{}

type OrderedMap struct {
	index  index
	keys   keys
	values values
}

func NewOrderedMap() *OrderedMap {
	om := OrderedMap{}
	om.Renew()
	return &om
}
func (me *OrderedMap) Renew() {
	me.index = make(index, 0)
	me.keys = make(keys, 0)
	me.values = make(values, 0)
}

func (me *OrderedMap) Get(key string) (value interface{}, exists bool) {
	var pos int
	pos, exists = me.index[key]
	if exists {
		value = me.values[pos]
	}
	return value, exists
}

func (me *OrderedMap) Set(key string, value interface{}) (exists bool) {
	pos, exists := me.index[key]
	if exists {
		me.values[pos] = value
	} else {
		me.keys = append(me.keys, key)
		me.values = append(me.values, value)
		me.index[key] = len(me.keys) - 1
	}
	return exists
}

func (me *OrderedMap) Delete(key string) (value interface{}, exists bool) {
	pos, exists := me.index[key]
	if exists {
		value = me.values[pos]
		me.keys = append(me.keys[:pos], me.keys[pos+1:]...)
		me.values = append(me.values[:pos], me.values[pos+1:]...)
		delete(me.index, key)
		for i := len(me.keys) - 1; i >= pos; i-- {
			me.index[me.keys[i]]--
		}
	}
	return value, exists
}

func (me *OrderedMap) Len() int {
	return len(me.keys)
}

func (me *OrderedMap) Keys() []string {
	return me.keys
}

func (me *OrderedMap) Values() []interface{} {
	return me.values
}

func (me *OrderedMap) Map() map[string]interface{} {
	m := make(map[string]interface{}, len(me.keys))
	for k, i := range me.index {
		m[k] = me.values[i]
	}
	return m
}

func (me *OrderedMap) UnmarshalJSON(b []byte) (err error) {
	for range only.Once {
		data := make(map[string]json.RawMessage, 0)
		err = json.Unmarshal(b, &data)
		if err != nil {
			break
		}
		me.Renew()
		for k, v := range data {
			me.Set(k, v)
		}
	}
	return err
}

func (me OrderedMap) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("{%s}", "")
	return []byte(s), nil
}
