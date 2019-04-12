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

type ApiModeler interface {
	NameGetter
	BasepathGetter
	IdParamsGetter
	ListFiltersGetter
	ListGetter
	ListIdsGetter
	ListItemAdder
	ListItemUpdater
	ListItemDeleter
	ListItemGetter
	ListItemDetailsGetter
	ListFilterer
	ListLinkMapGetter
	ItemFilterer
	ItemCanAdder
	RelatedItemsGetter
}
type ItemCanAdder interface {
	CanAddItem(*Context) bool
}
type NameGetter interface {
	GetName() types.RouteName
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
	GetListIds(*Context, ...FilterPath) (ItemIds, status.Status)
}
type ListItemAdder interface {
	AddItem(*Context, ApiItemer) status.Status
}
type ListItemUpdater interface {
	UpdateItem(*Context, ApiItemer) status.Status
}
type ListItemDeleter interface {
	DeleteItem(*Context, ItemId) status.Status
}
type ListItemGetter interface {
	GetItem(*Context, ItemId) (ApiItemer, status.Status)
}
type ListItemDetailsGetter interface {
	GetItemDetails(*Context, ItemId) (ApiItemer, status.Status)
}
type ItemFilterer interface {
	FilterItem(ApiItemer, FilterPath) (ApiItemer, status.Status)
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
type RelatedItemsGetter interface {
	GetRelatedItems(ctx *Context, item ApiItemer) (list List, sts status.Status)
}

type ApiItemer interface {
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
	GetItem() (ApiItemer, status.Status)
}
type ItemIdSetter interface {
	SetId(ItemId) status.Status
}
type ItemLinkMapGetter interface {
	GetItemLinkMap(*Context) (LinkMap, status.Status)
}
type RootDocumenter interface {
	ResponseTypeGetter
	RootDocumentGetter
	IncludedSetter
}
type ResponseTypeGetter interface {
	GetResponseType() types.ResponseType
}
type RootDocumentGetter interface {
	GetRootDocument() interface{}
}
type IncludedSetter interface {
	SetIncluded(*Context, List) status.Status
}
