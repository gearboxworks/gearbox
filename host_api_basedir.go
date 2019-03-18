package gearbox

import (
	"fmt"
	"gearbox/api"
	"gearbox/only"
	"gearbox/stat"
	"net/http"
)

func (me *HostApi) addBasedirRoutes() {

	me.GET("/basedirs", "basedirs", me.getHostBasedirsResponse, nil)
	me.PUT("/basedirs/:nickname", "basedir-update", me.updateBasedir, nil)
	me.POST("/basedirs/new", "basedir-add", me.addBasedir, nil)
	me.DELETE("/basedirs/:nickname", "basedir-delete", me.deleteNamedBasedir, nil)

}

func (me *HostApi) getHostBasedirsResponse(rc *api.RequestContext) interface{} {
	return me.Config.GetHostBasedirs()
}

func (me *HostApi) addBasedir(rc *api.RequestContext) interface{} {
	var status stat.Status
	for range only.Once {
		bd := Basedir{}
		status = rc.UnmarshalFromRequest(&bd)
		if status.IsError() {
			status.PriorStatus = status.String()
			status.Message = "Unable to add basedir."
			status.HttpStatus = http.StatusBadRequest
			status.ApiHelp = fmt.Sprintf("verify that API request is in the correct format: %s",
				api.GetApiDocsUrl(rc.ResourceName),
			)
			break
		}
		me.Gearbox.SetResourceName(rc.ResourceName)
		status = me.Gearbox.AddBasedir(bd.HostDir, bd.Nickname)
	}
	return status
}
func (me *HostApi) updateBasedir(rc *api.RequestContext) interface{} {
	var status stat.Status
	for range only.Once {
		bd := Basedir{}
		status = rc.UnmarshalFromRequest(&bd)
		if status.IsError() {
			status.PriorStatus = status.String()
			status.Message = "Unable to update basedir."
			status.HttpStatus = http.StatusBadRequest
			status.ApiHelp = fmt.Sprintf("verify that API request is in the correct format: %s",
				api.GetApiDocsUrl(rc.ResourceName),
			)
			break
		}
		me.Gearbox.SetResourceName(rc.ResourceName)
		status = me.Gearbox.UpdateBasedir(bd.Nickname, bd.HostDir)
	}
	return status
}
func (me *HostApi) deleteNamedBasedir(rc *api.RequestContext) interface{} {
	me.Gearbox.SetResourceName(rc.ResourceName)
	return me.Gearbox.DeleteNamedBasedir(getBasedirNickname(rc))
}

func getBasedirNickname(rc *api.RequestContext) string {
	return rc.Param("nickname")
}
