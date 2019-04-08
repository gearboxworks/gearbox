package api

import (
	"gearbox/apimodeler"
	"gearbox/status"
	"gearbox/types"
)

type Apier interface {
	AddModels(apimodeler.Modeler) status.Status
	SetParent(interface{})
	GetBaseUrl() types.UrlTemplate
	GetItemUrl(ctx *apimodeler.Context, id apimodeler.Itemer) (types.UrlTemplate, status.Status)
	Start()
	Stop()
}
