package gearbox

import (
	"fmt"
	"gearbox/api"
	"gearbox/only"
	"net/http"
)

func (me *HostApi) addBasedirRoutes() {

	me.GET("/basedirs", "basedirs", me.getHostBasedirsResponse)

	me.POST("/basedirs/new", "basedir-add", me.addBasedir)

	me.PUT("/basedirs/:nickname", "basedir-update", me.updateBasedir)

	me.DELETE("/basedirs/:nickname", "basedir-delete", me.deleteNamedBasedir)

}

func (me *HostApi) getHostBasedirsResponse(rc *api.RequestContext) interface{} {
	return me.Config.GetHostBasedirs()
}

func (me *HostApi) addBasedir(rc *api.RequestContext) interface{} {
	var status Status
	for range only.Once {
		bd := Basedir{}
		err := rc.UnmarshalFromRequest(&bd)
		if err != nil {
			status = NewStatus(&StatusArgs{
				Failed:     true,
				Message:    "Unable to add basedir.",
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    fmt.Sprintf("verify that API request is in the correct format: %s", api.GetApiDocsUrl(rc.ResourceName)),
				Error:      err,
			})
			break
		}
		me.Gearbox.RequestType = rc.ResourceName
		status = me.Gearbox.AddBasedir(bd.HostDir, bd.Nickname)
	}
	return status
}
func (me *HostApi) updateBasedir(rc *api.RequestContext) interface{} {
	var status Status
	for range only.Once {
		bd := Basedir{}
		err := rc.UnmarshalFromRequest(&bd)
		if err != nil {
			status = NewStatus(&StatusArgs{
				Failed:     true,
				Message:    "Unable to update basedir.",
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    fmt.Sprintf("verify that API request is in the correct format: %s", api.GetApiDocsUrl(rc.ResourceName)),
				Error:      err,
			})
			break
		}
		me.Gearbox.RequestType = rc.ResourceName
		status = me.Gearbox.UpdateBasedir(bd.Nickname, bd.HostDir)
	}
	return status
}
func (me *HostApi) deleteNamedBasedir(rc *api.RequestContext) interface{} {
	me.Gearbox.RequestType = rc.ResourceName
	return me.Gearbox.DeleteNamedBasedir(getBasedirNickname(rc))
}

func getBasedirNickname(rc *api.RequestContext) string {
	return rc.Param("nickname")
}
