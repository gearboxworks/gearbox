package api

import (
	"gearbox/apiworks"
	"gearbox/status"
	"gearbox/types"
)

type Apier interface {
	AddController(apiworks.ListController) status.Status
	SetParent(interface{})
	GetBaseUrl() types.UrlTemplate
	GetItemUrl(ctx *apiworks.Context, id apiworks.ItemModeler) (types.UrlTemplate, status.Status)
	Start()
	Stop()
	GetRootLinkMap(ctx *apiworks.Context) apiworks.LinkMap
	GetListLinkMap(ctx *apiworks.Context) apiworks.LinkMap
	GetItemLinkMap(ctx *apiworks.Context) apiworks.LinkMap
	GetCommonLinkMap(ctx *apiworks.Context) apiworks.LinkMap
	WireRoutes()
}
