package test

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
)

var NilTestableModel = (*TestableModel)(nil)
var _ apimodeler.Modeler = NilTestableModel

type TestableModel struct {
	Collection apimodeler.Collection
}

func NewTestableModel() *TestableModel {
	tm := TestableModel{
		Collection: make(apimodeler.Collection, 0),
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
	return testableModelBasepath
}

func (me *TestableModel) GetIdParams() apimodeler.IdParams {
	return testableModelIdParams
}

func (me *TestableModel) GetFilterMap() apimodeler.FilterMap {
	return apimodeler.FilterMap{
		FrobinatorsFilter: apimodeler.Filter{
			Label: "Items of Frobinator type",
			Path:  FrobinatorsFilter,
			CollectionFilter: func(coll apimodeler.Collection) apimodeler.Collection {
				newcoll := make(apimodeler.Collection, 0)
				for _, item := range coll {
					if item.GetType() != FrobinatorType {
						continue
					}
					newcoll = append(newcoll, item)
				}
				return newcoll
			},
		},
		UnicornFilter: apimodeler.Filter{
			Label: "Items of Unicorn type",
			Path:  UnicornFilter,
			ItemFilter: func(item apimodeler.Itemer) apimodeler.Itemer {
				if item.GetType() != UnicornType {
					return nil
				}
				return item
			},
		},
	}
}

func (me *TestableModel) GetCollection(filter ...apimodeler.FilterPath) (coll apimodeler.Collection, sts status.Status) {
	return me.Collection, sts
}

func (me *TestableModel) AddItem(item apimodeler.Itemer) status.Status {
	me.Collection = append(me.Collection, item)
	return nil
}

func (me *TestableModel) DeleteItem(itemid apimodeler.ItemId) (sts status.Status) {
	for range only.Once {
		index, sts := me.getItemIndex(itemid)
		if is.Error(sts) {
			break
		}
		me.Collection = append(me.Collection[:index], me.Collection[index+1:]...)
	}
	return sts
}

func (me *TestableModel) GetItem(itemid apimodeler.ItemId) (item apimodeler.Itemer, sts status.Status) {
	for range only.Once {
		index, sts := me.getItemIndex(itemid)
		if is.Error(sts) {
			break
		}
		item = me.Collection[index]
	}
	return item, sts
}

func (me *TestableModel) getItemIndex(itemid apimodeler.ItemId) (index int, sts status.Status) {
	found := false
	var item apimodeler.Itemer
	for index, item = range me.Collection {
		if item.GetId() != itemid {
			continue
		}
		found = true
		break
	}
	if !found {
		index = -1
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item '%s' not found", itemid),
		})
	}
	return index, sts
}

func (me *TestableModel) GetCollectionIds() (apimodeler.ItemIds, status.Status) {
	cids := make(apimodeler.ItemIds, len(me.Collection))
	for i, c := range me.Collection {
		cids[i] = c.GetId()
	}
	return cids, nil
}

func (me *TestableModel) UpdateItem(item apimodeler.Itemer) (sts status.Status) {
	for range only.Once {
		var index int
		index, sts = me.getItemIndex(item.GetId())
		if is.Error(sts) {
			break
		}
		me.Collection[index] = item
	}
	return sts
}

func (me *TestableModel) FilterItem(item apimodeler.Itemer, filter apimodeler.FilterPath) (_item apimodeler.Itemer, sts status.Status) {
	for range only.Once {
		_item = item
		f, sts := me.getFilter(filter)
		if is.Error(sts) {
			break
		}
		if f.ItemFilter == nil {
			break
		}
		_item = f.ItemFilter(item)
	}
	return _item, sts
}

func (me *TestableModel) FilterCollection(filter apimodeler.FilterPath) (coll apimodeler.Collection, sts status.Status) {
	for range only.Once {
		coll = me.Collection
		fm := me.GetFilterMap()
		f, ok := fm[filter]
		if !ok {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("filter '%s' not found ", filter),
			})
			break
		}
		if f.CollectionFilter == nil {
			break
		}
		coll = f.CollectionFilter(coll)
	}
	return coll, sts
}

func (me *TestableModel) getFilter(filter apimodeler.FilterPath) (f apimodeler.Filter, sts status.Status) {
	fm := me.GetFilterMap()
	f, ok := fm[filter]
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("filter '%s' not found ", filter),
		})
	}
	return f, sts
}
