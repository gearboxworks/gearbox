package test

import (
	"gearbox/modeler"
	"gearbox/status"
	"gearbox/types"
	"github.com/labstack/echo"
)

const TestableModelBasepath = "/foo"

var TestableModelIdParams = modeler.IdParams{"foo", "bar"}

var NilTestableModel = (*TestableModel)(nil)
var _ modeler.Modeler = NilTestableModel

type TestableModel struct {
	Collection modeler.Collection
}

func NewTestableModel() *TestableModel {
	tm := TestableModel{
		Collection: make(modeler.Collection, 0),
	}
	coll := tm.Collection
	for id, typ := range testableItemData {
		coll = append(coll, &TestableItem{
			Id:   id,
			Type: typ,
		})
	}
	tm.Collection = coll
	return &tm
}

func (me *TestableModel) GetBasepath() types.Basepath {
	return TestableModelBasepath
}

func (me *TestableModel) GetIdParams() modeler.IdParams {
	return TestableModelIdParams
}

func (me *TestableModel) GetCollectionFilterMap() modeler.FilterMap {
	panic("implement me")
}

func (me *TestableModel) GetCollection(modeler.FilterPath) (modeler.Collection, status.Status) {
	return me.Collection, nil
}

func (me *TestableModel) GetCollectionIds() (modeler.ItemIds, status.Status) {
	panic("implement me")
}

func (me *TestableModel) AddItem(item modeler.Item) status.Status {
	me.Collection = append(me.Collection, item)
	return nil
}

func (me *TestableModel) UpdateItem(modeler.Item) status.Status {
	panic("implement me")
}

func (me *TestableModel) DeleteItem(itemid modeler.ItemId) (sts status.Status) {
	found := false
	for index, item := range me.Collection {
		if item.GetId() != itemid {
			continue
		}
		me.Collection = append(me.Collection[:index], me.Collection[index+1:]...)
		found = true
		break
	}
	if !found {
		sts = status.Fail(&status.Args{})
	}
	return sts
}

func (me *TestableModel) GetItem(modeler.ItemId, echo.Context) (modeler.Item, status.Status) {

	panic("implement me")
}

func (me *TestableModel) FilterItem(modeler.Item, modeler.FilterPath) (modeler.Item, status.Status) {
	panic("implement me")
}
