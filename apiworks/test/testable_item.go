package test

import (
	"gearbox/apiworks"
	"gearbox/status"
)

var NilTestableItem = (*TestableItem)(nil)
var _ apiworks.ItemModeler = NilTestableItem

type TestableItem struct {
	Id   apiworks.ItemId
	Type apiworks.ItemType
}

func (me *TestableItem) SetId(itemid apiworks.ItemId) status.Status {
	me.Id = itemid
	return nil
}

func (me *TestableItem) GetId() apiworks.ItemId {
	return me.Id
}

func (me *TestableItem) GetType() apiworks.ItemType {
	return me.Type
}

func (me *TestableItem) GetItem() (apiworks.ItemModeler, status.Status) {
	return me, nil
}
