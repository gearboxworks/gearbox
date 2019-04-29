package apiworks

import "gearbox/status"

type ItemModeler interface {
	ItemIdGetter
	ItemIdSetter
	ItemTypeGetter
	ItemGetter
	ItemLinkMapGetter
	RelatedItemsGetter
}

type RelatedItemsGetter interface {
	GetRelatedItems(ctx *Context) (list List, sts status.Status)
}

type ItemIdGetter interface {
	GetId() ItemId
}
type ItemTypeGetter interface {
	GetType() ItemType
}
type ItemGetter interface {
	GetItem() (ItemModeler, status.Status)
}
type ItemIdSetter interface {
	SetId(ItemId) status.Status
}
type ItemLinkMapGetter interface {
	GetItemLinkMap(*Context) (LinkMap, status.Status)
}
