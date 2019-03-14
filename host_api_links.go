package gearbox

import "gearbox/api"

func (me *HostApi) addRoutes() {

	me.GET("/", api.LinksResource, nil)

	me.addBasedirRoutes()
	me.addProjectRoutes()
	me.addStackRoutes()

}
