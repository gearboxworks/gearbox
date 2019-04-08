package apimodeler

import (
	"gearbox/status"
	"gearbox/types"
)

type Contexter interface {
	ParamGetter
	KeyValueGetter
	KeyValueSetter
}
type ParamGetter interface {
	Param(string) string
}
type KeyValueGetter interface {
	Get(string) interface{}
}
type KeyValueSetter interface {
	Set(string, interface{})
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
	GetCollection(Contexter, ...FilterPath) (Collection, status.Status)
}
type CollectionIdsGetter interface {
	GetCollectionIds(Contexter) (ItemIds, status.Status)
}
type CollectionItemAdder interface {
	AddItem(Contexter, Itemer) status.Status
}
type CollectionItemUpdater interface {
	UpdateItem(Contexter, Itemer) status.Status
}
type CollectionItemDeleter interface {
	DeleteItem(Contexter, ItemId) status.Status
}
type CollectionItemGetter interface {
	GetItem(Contexter, ItemId) (Itemer, status.Status)
}
type ItemFilterer interface {
	FilterItem(Itemer, FilterPath) (Itemer, status.Status)
}
type CollectionFilterer interface {
	FilterCollection(Contexter, FilterPath) (Collection, status.Status)
}
type CollectionFiltersGetter interface {
	GetFilterMap() FilterMap
}
type ItemIdsGetter interface {
	GetIds(Contexter) ItemIds
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

type Critera interface{}
