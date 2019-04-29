package api

import (
	"gearbox/apiworks"
	"gearbox/status"
	"gearbox/types"
)

type (
	HttpHeaderValue   = apiworks.HttpHeaderValue
	BasepathGetter    = apiworks.BasepathGetter
	ContextArgs       = apiworks.ContextArgs
	Context           = apiworks.Context
	Contexter         = apiworks.Contexter
	ItemId            = apiworks.ItemId
	ControllerMap     = apiworks.ControllerMap
	ItemLinkMapGetter = apiworks.ItemLinkMapGetter
	ItemModeler       = apiworks.ItemModeler
	Link              = apiworks.Link
	LinkMap           = apiworks.LinkMap
	List              = apiworks.List
	ListController    = apiworks.ListController
	ListGetter        = apiworks.ListGetter
	Metaname          = apiworks.Metaname
	RelType           = apiworks.RelType
	Status            = status.Status
)

type Apier interface {
	AddController(ListController) Status
	SetParent(interface{})
	GetBaseUrl() types.UrlTemplate
	GetItemUrl(ctx *Context, id ItemModeler) (types.UrlTemplate, Status)
	Start()
	Stop()
	GetRootLinkMap(ctx *Context) LinkMap
	GetListLinkMap(ctx *Context) LinkMap
	GetItemLinkMap(ctx *Context) LinkMap
	GetCommonLinkMap(ctx *Context) LinkMap
	WireRoutes()
}
