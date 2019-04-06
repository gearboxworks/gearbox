package modeler

import (
	"gearbox/status"
	"gearbox/status/is"
)

type Collection []Item

func (Collection) ContainsResource() {}

func (me Collection) GetIds() (ItemIds, status.Status) {
	itemIds := make(ItemIds, len(me))
	for i, item := range me {
		itemIds[i] = item.GetId()
	}
	return itemIds, nil
}

func (me Collection) GetItems() (Collection, status.Status) {
	return me, nil
}

func GetCollectionSlice(collection Collection, sts status.Status) (Collection, status.Status) {
	var slice = make(Collection, len(collection))
	if is.Success(sts) {
		for _, item := range collection {
			slice = append(slice, item)
		}
	}
	return slice, sts
}
