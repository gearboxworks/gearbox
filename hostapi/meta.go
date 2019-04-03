package hostapi

import (
	"gearbox/api"
	"gearbox/routes"
	"gearbox/status"
	"gearbox/status/is"
)

func (me *HostApi) addRoutes() (sts status.Status) {

	//me.GET___("/", api.LinksResource)
	//me.GET___("/meta/endpoints", api.MetaEndpointsResource, me.getMetaEndpointsResponse)
	//me.GET___("/meta/methods", api.MetaMethodsResource, me.getMetaMethodsResponse)

	sts = me.AddConnector(routes.NewProjectConnector(me.Gearbox))
	if is.Error(sts) {
		panic(sts.Message())
	}

	sts = me.AddConnector(routes.NewStackConnector(me.Gearbox))
	if is.Error(sts) {
		panic(sts.Message())
	}
	//me.addBasedirRoutes()
	//sts := me.AddConnector(e,routes.NewBasedirFactory(me.Gearbox))
	//if is.Error(sts) {
	//	panic(sts.Message())
	//}
	//me.addStackRoutes()
	//sts = me.AddConnector(e,routes.NewStackFactory(me.Gearbox))
	//if is.Error(sts) {
	//	panic(sts.Message())
	//}

	return

}
func (me *HostApi) getMetaEndpointsResponse(rc *api.RequestContext) interface{} {
	return me.Api.EndpointMap
}
func (me *HostApi) getMetaMethodsResponse(rc *api.RequestContext) interface{} {
	return me.Api.MethodMap
}
