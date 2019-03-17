package gearbox

import "gearbox/api"

func (me *HostApi) addRoutes() {

	me.GET("/", api.LinksResource, nil)

	me.GET("/meta/endpoints", api.MetaEndpointsResource, me.getMetaEndpointsResponse)

	me.GET("/meta/methods", api.MetaMethodsResource, me.getMetaMethodsResponse)

	me.addBasedirRoutes()
	me.addProjectRoutes()
	me.addStackRoutes()

}
func (me *HostApi) getMetaEndpointsResponse(rc *api.RequestContext) interface{} {
	return me.Api.EndpointMap
}
func (me *HostApi) getMetaMethodsResponse(rc *api.RequestContext) interface{} {
	return me.Api.MethodMap
}
