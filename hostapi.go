package gearbox

import (
	"gearbox/api"
	"gearbox/status"
)

type HostApi interface {
	SetGearbox(Gearboxer)
	GetBaseUrl() api.UriTemplate
	GetMethodMap() api.MethodMap
	//GetUriTemplateVars(api.RouteName, interface{}, int) (api.UriTemplateVars, status.Status)
	GetUrl(api.RouteName, api.UriTemplateVars) (api.UriTemplate, status.Status)
	GetUrlPath(api.RouteName, api.UriTemplateVars) (api.UriTemplate, status.Status)
	Start()
	Stop()
}
