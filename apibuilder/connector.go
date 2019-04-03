package ab

import (
	"gearbox/status"
	"gearbox/status/is"
	"github.com/labstack/echo"
)

type Connector interface {
	BasepathGetter
	IdParamsGetter
	CollectionFiltersGetter
	CollectionGetter
	CollectionIdsGetter
	CollectionItemAdder
	CollectionItemUpdater
	CollectionItemDeleter
	CollectionItemGetter
	ItemFilterer
}

type BasepathGetter interface {
	GetBasepath() Basepath
}
type IdParamsGetter interface {
	GetIdParams() IdParams
}
type CollectionGetter interface {
	GetCollection(FilterPath) (Collection, status.Status)
}
type CollectionIdsGetter interface {
	GetCollectionIds() (ItemIds, status.Status)
}
type CollectionItemAdder interface {
	AddItem(Item) (Collection, status.Status)
}
type CollectionItemUpdater interface {
	UpdateItem(Item) (Collection, status.Status)
}
type CollectionItemDeleter interface {
	DeleteItem(ItemId) (Collection, status.Status)
}
type CollectionItemGetter interface {
	GetItem(ItemId, echo.Context) (Item, status.Status)
}
type ItemFilterer interface {
	FilterItem(Item, FilterPath) (Item, status.Status)
}
type CollectionFiltersGetter interface {
	GetCollectionFilterMap() FilterMap
}

type Collection []Item

func GetCollectionSlice(collection Collection, sts status.Status) (Collection, status.Status) {
	var slice = make(Collection, len(collection))
	if is.Success(sts) {
		for _, item := range collection {
			slice = append(slice, item)
		}
	}
	return slice, sts
}

func (me Collection) GetIds() (ItemIds, status.Status) {
	itemIds := make(ItemIds, len(me))
	for i, item := range me {
		itemIds[i] = item.GetId()
	}
	return itemIds, nil
}
func (Collection) ContainsResource() {}
func (me Collection) GetItems() (Collection, status.Status) {
	return me, nil
}

type Item interface {
	ItemIdGetter
	ItemTypeGetter
	ItemGetter
}

type ItemIds []ItemId
type ItemId string
type ItemType string

type ItemIdGetter interface {
	GetId() ItemId
}
type ItemTypeGetter interface {
	GetType() ItemType
}
type ItemGetter interface {
	GetItem() (Item, status.Status)
}
type ItemIdsGetter interface {
	GetIds() ItemIds
}
type ItemsGetter interface {
	GetItems() (Collection, status.Status)
}
type ItemIdSetter interface {
	SetId(ItemId)
}
