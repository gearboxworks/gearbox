package modeler

import (
	"gearbox/status"
	"gearbox/types"
)

type Modeler interface {
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
	GetBasepath() types.Basepath
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
	AddItem(Item) status.Status
}
type CollectionItemUpdater interface {
	UpdateItem(Item) status.Status
}
type CollectionItemDeleter interface {
	DeleteItem(ItemId) status.Status
}
type CollectionItemGetter interface {
	GetItem(ItemId) (Item, status.Status)
}
type ItemFilterer interface {
	FilterItem(Item, FilterPath) (Item, status.Status)
}
type CollectionFiltersGetter interface {
	GetCollectionFilterMap() FilterMap
}

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
