package api

import (
	"gearbox/apimodeler"
	"gearbox/status"
	"gearbox/types"
)

type Apier interface {
	AddController(apimodeler.ApiController) status.Status
	SetParent(interface{})
	GetBaseUrl() types.UrlTemplate
	GetItemUrl(ctx *apimodeler.Context, id apimodeler.Itemer) (types.UrlTemplate, status.Status)
	Start()
	Stop()
	GetRootLinkMap(ctx *apimodeler.Context) apimodeler.LinkMap
	GetListLinkMap(ctx *apimodeler.Context) apimodeler.LinkMap
	GetItemLinkMap(ctx *apimodeler.Context) apimodeler.LinkMap
	GetCommonLinkMap(ctx *apimodeler.Context) apimodeler.LinkMap
}
