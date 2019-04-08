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
	Start()
	Stop()
}
