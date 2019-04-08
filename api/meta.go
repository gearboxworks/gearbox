package api

//func (me *Api) addRoutes() (sts status.Status) {
//
//	//me.GET___("/", api.LinksResource)
//	//me.GET___("/meta/endpoints", api.MetaEndpointsResource, me.getMetaEndpointsResponse)
//	//me.GET___("/meta/methods", api.MetaMethodsResource, me.getMetaMethodsResponse)
//
//	sts = me.AddModels(models.NewProjectModel(me.Parent))
//	if is.Error(sts) {
//		panic(sts.Message())
//	}
//
//	sts = me.AddModels(models.NewStackConnector(me.Parent))
//	if is.Error(sts) {
//		panic(sts.Message())
//	}
//	//me.addBasedirRoutes()
//	//sts := me.AddModels(e,models.NewBasedirFactory(me.Parent))
//	//if is.Error(sts) {
//	//	panic(sts.Message())
//	//}
//	//me.addStackRoutes()
//	//sts = me.AddModels(e,models.NewStackFactory(me.Parent))
//	//if is.Error(sts) {
//	//	panic(sts.Message())
//	//}
//
//	return
//
//}
//func (me *Api) getMetaEndpointsResponse(rc *RequestContext) interface{} {
//	return me.Api.EndpointMap
//}
//func (me *Api) getMetaMethodsResponse(rc *RequestContext) interface{} {
//	return me.Api.MethodMap
//}
