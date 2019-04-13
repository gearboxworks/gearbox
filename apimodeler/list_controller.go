package apimodeler

import (
	"gearbox/status"
	"gearbox/types"
)

type ListController interface {
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
	ParentGetter
	RelatedFieldsGetter
	LinkAdder
	LinksAdder
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
	AddItem(*Context, ItemModeler) status.Status
}
type ListItemUpdater interface {
	UpdateItem(*Context, ItemModeler) status.Status
}
type ListItemDeleter interface {
	DeleteItem(*Context, ItemId) status.Status
}
type ListItemGetter interface {
	GetItem(*Context, ItemId) (ItemModeler, status.Status)
}
type ListItemDetailsGetter interface {
	GetItemDetails(*Context, ItemId) (ItemModeler, status.Status)
}
type ItemFilterer interface {
	FilterItem(ItemModeler, FilterPath) (ItemModeler, status.Status)
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
type ParentGetter interface {
	GetParent() ListController
}
type RelatedFieldsGetter interface {
	GetRelatedFields() RelatedFields
}
type LinkAdder interface {
	AddLink(rel RelType, link LinkImplementor)
}
type LinksAdder interface {
	AddLinks(links LinkMap)
}
