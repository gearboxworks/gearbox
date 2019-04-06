package test

import (
	"gearbox/apimodeler"
	"gearbox/status"
)

var NilTestableItem = (*TestableItem)(nil)
var _ apimodeler.Itemer = NilTestableItem

type TestableItem struct {
	Id   apimodeler.ItemId
	Type apimodeler.ItemType
}

func (me *TestableItem) SetId(itemid apimodeler.ItemId) {
	me.Id = itemid
}

func (me *TestableItem) SetType(typ apimodeler.ItemType) {
	me.Type = typ
}

func (me *TestableItem) GetId() apimodeler.ItemId {
	return me.Id
}

func (me *TestableItem) GetType() apimodeler.ItemType {
	return me.Type
}

func (me *TestableItem) GetItem() (apimodeler.Itemer, status.Status) {
	return me, nil
}
