package test

import (
	"gearbox/apimodeler"
	"gearbox/status"
)

var NilTestableItem = (*TestableItem)(nil)
var _ apimodeler.ApiItemer = NilTestableItem

type TestableItem struct {
	Id   apimodeler.ItemId
	Type apimodeler.ItemType
}

func (me *TestableItem) SetId(itemid apimodeler.ItemId) status.Status {
	me.Id = itemid
	return nil
}

func (me *TestableItem) GetId() apimodeler.ItemId {
	return me.Id
}

func (me *TestableItem) GetType() apimodeler.ItemType {
	return me.Type
}

func (me *TestableItem) GetItem() (apimodeler.ApiItemer, status.Status) {
	return me, nil
}
