package apiworks

import (
	"gearbox/types"
	"github.com/gearboxworks/go-status"
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
	ListFilterer
	ListLinkMapGetter
	NilItemGetter
	ItemFilterer
	ItemCanAdder
	ParentGetter
	RelatedFieldsGetter
	LinkAdder
	LinksAdder
	RootObjectGetter
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
	AddItem(*Context, ItemModeler) (ItemModeler, status.Status)
}
type ListItemUpdater interface {
	UpdateItem(*Context, ItemModeler) (ItemModeler, status.Status)
}
type ListItemDeleter interface {
	DeleteItem(*Context, ItemId) status.Status
}
type ListItemGetter interface {
	GetItem(*Context, ItemId) (ItemModeler, status.Status)
}
type NilItemGetter interface {
	GetNilItem(ctx *Context) ItemModeler
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
type MetaAdder interface {
	AddMeta(Metaname, MetaValue)
}
type RootObjectGetter interface {
	GetRootObject() interface{}
}
