package test

import (
	"gearbox/apiworks"
	"github.com/gearboxworks/go-status"
)

var NilTestableItem = (*TestableItem)(nil)
var _ apiworks.ItemModeler = NilTestableItem

type TestableItem struct {
	Id   apiworks.ItemId
	Type apiworks.ItemType
}

func (me *TestableItem) GetItemLinkMap(*apiworks.Context) (apiworks.LinkMap, status.Status) {
	panic("implement me")
}

func (me *TestableItem) GetRelatedItems(ctx *apiworks.Context) (list apiworks.List, sts status.Status) {
	panic("implement me")
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
