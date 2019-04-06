package apimodeler

import (
	"gearbox/status"
	"gearbox/types"
)

type Contexter interface {
	ParamGetter
}
type ParamGetter interface {
	Param(string) string
}

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
	CollectionFilterer
	ItemFilterer
}

type BasepathGetter interface {
	GetBasepath() types.Basepath
}
type IdParamsGetter interface {
	GetIdParams() IdParams
}
type CollectionGetter interface {
	GetCollection() (Collection, status.Status)
}
type CollectionIdsGetter interface {
	GetCollectionIds() (ItemIds, status.Status)
}
type CollectionItemAdder interface {
	AddItem(Itemer) status.Status
}
type CollectionItemUpdater interface {
	UpdateItem(Itemer) status.Status
}
type CollectionItemDeleter interface {
	DeleteItem(ItemId) status.Status
}
type CollectionItemGetter interface {
	GetItem(ItemId) (Itemer, status.Status)
}
type ItemFilterer interface {
	FilterItem(Itemer, FilterPath) (Itemer, status.Status)
}
type CollectionFilterer interface {
	FilterCollection(FilterPath) (Collection, status.Status)
}
type CollectionFiltersGetter interface {
	GetFilterMap() FilterMap
}

type ItemIdsGetter interface {
	GetIds() ItemIds
}
type ItemsGetter interface {
	GetItems() (Collection, status.Status)
}

type Itemer interface {
	ItemIdGetter
	ItemIdSetter
	ItemTypeGetter
	ItemGetter
}

type ItemIdGetter interface {
	GetId() ItemId
}
type ItemTypeGetter interface {
	GetType() ItemType
}
type ItemGetter interface {
	GetItem() (Itemer, status.Status)
}
type ItemIdSetter interface {
	SetId(ItemId) status.Status
}
