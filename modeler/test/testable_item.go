package test

import (
	"gearbox/modeler"
	"gearbox/status"
)

var NilTestableItem = (*TestableItem)(nil)
var _ modeler.Item = NilTestableItem

type TestableItem struct {
	Id   modeler.ItemId
	Type modeler.ItemType
}

func (me *TestableItem) GetId() modeler.ItemId {
	return me.Id
}

func (me *TestableItem) GetType() modeler.ItemType {
	return me.Type
}

func (me *TestableItem) GetItem() (modeler.Item, status.Status) {
	return me, nil
}
