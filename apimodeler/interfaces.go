package apimodeler

import (
	"gearbox/status"
	"gearbox/types"
	"net/http"
)

type Contexter interface {
	ParamGetter
	KeyValueGetter
	KeyValueSetter
	RequestGetter
}

type RequestGetter interface {
	Request() *http.Request
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
	ListFiltersGetter
	ListGetter
	ListIdsGetter
	ListItemAdder
	ListItemUpdater
	ListItemDeleter
	ListItemGetter
	ListFilterer
	ListLinkMapGetter
	ItemFilterer
}

type BasepathGetter interface {
	GetBasepath() types.Basepath
}
type IdParamsGetter interface {
	GetIdParams() IdParams
}
type ListGetter interface {
	GetList(*Context, ...FilterPath) (List, status.Status)
}
type ListIdsGetter interface {
	GetListIds(*Context) (ItemIds, status.Status)
}
type ListItemAdder interface {
	AddItem(*Context, Itemer) status.Status
}
type ListItemUpdater interface {
	UpdateItem(*Context, Itemer) status.Status
}
type ListItemDeleter interface {
	DeleteItem(*Context, ItemId) status.Status
}
type ListItemGetter interface {
	GetItem(*Context, ItemId) (Itemer, status.Status)
}
type ItemFilterer interface {
	FilterItem(Itemer, FilterPath) (Itemer, status.Status)
}
type ListFilterer interface {
	FilterList(*Context, FilterPath) (List, status.Status)
}
type ListFiltersGetter interface {
	GetFilterMap() FilterMap
}
type ItemIdsGetter interface {
	GetIds(*Context) ItemIds
}
type ListLinkMapGetter interface {
	GetListLinkMap(*Context, ...FilterPath) (LinkMap, status.Status)
}

type Itemer interface {
	ItemIdGetter
	ItemIdSetter
	ItemTypeGetter
	ItemGetter
	ItemLinkMapGetter
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
type ItemLinkMapGetter interface {
	GetItemLinkMap(*Context) (LinkMap, status.Status)
}

type ResponseTypeGetter interface {
	GetResponseType() types.ResponseType
}
type RootDocumentGetter interface {
	ResponseTypeGetter
	GetRootDocument() interface{}
}
