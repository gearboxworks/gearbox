package test

import (
	"gearbox/apimodeler"
	"gearbox/status"
)

var NilTestableItem = (*TestableItem)(nil)
var _ apimodeler.ItemModeler = NilTestableItem

type TestableItem struct {
	Id   apimodeler.ItemId
	Type apimodeler.ItemType
}

func (me *TestableItem) SetStackId(itemid apimodeler.ItemId) status.Status {
	me.Id = itemid
	return nil
}

func (me *TestableItem) GetId() apimodeler.ItemId {
	return me.Id
}

func (me *TestableItem) GetType() apimodeler.ItemType {
	return me.Type
}

func (me *TestableItem) GetItem() (apimodeler.ItemModeler, status.Status) {
	return me, nil
}
