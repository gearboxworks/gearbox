package hostapi

import (
	"fmt"
	"gearbox/api"
	"gearbox/config"
	"gearbox/only"
	"gearbox/status"
	"gearbox/types"
	"net/http"
)

const BasedirsResource api.RouteName = "basedirs"
const BasedirsAdd api.RouteName = "basedir-add"
const BasedirsUpdate api.RouteName = "basedir-update"
const BasedirsDelete api.RouteName = "basedir-delete"

func (me *HostApi) addBasedirRoutes() {

	//me.GET___("/basedirs", BasedirsResource, me.getHostBasedirsResponse, nil)
	//me.PUT___("/basedirs/:nickname", "basedir-update", me.updateBasedir, nil)
	//me.POST__("/basedirs/new", BasedirsAdd, me.addBasedir, nil)
	//me.DELETE("/basedirs/:nickname", BasedirsDelete, me.deleteNamedBasedir, nil)
	//
	//me.Relate(BasedirsResource, &api.Related{
	//	List:   BasedirsResource,
	//	New:    BasedirsAdd,
	//	Update: BasedirsUpdate,
	//	Delete: BasedirsDelete,
	//})

}

func (me *HostApi) getHostBasedirsResponse(rc *api.RequestContext) interface{} {
	return me.Config.GetHostBasedirs()
}

func (me *HostApi) addBasedir(rc *api.RequestContext) interface{} {
	var sts status.Status
	for range only.Once {
		bd := config.Basedir{}
		sts = rc.UnmarshalFromRequest(&bd)
		if status.IsError(sts) {
			sts = status.Wrap(sts, &status.Args{
				Message:    "unable to add basedir.",
				HttpStatus: http.StatusBadRequest,
				ApiHelp: fmt.Sprintf("verify that API request is in the correct format: %s",
					api.GetApiDocsUrl(rc.RouteName),
				),
			})
			break
		}
		me.Gearbox.SetRouteName(rc.RouteName)
		sts = me.Gearbox.AddBasedir(bd.HostDir, bd.Nickname)
	}
	return sts
}
func (me *HostApi) updateBasedir(rc *api.RequestContext) interface{} {
	var sts status.Status
	for range only.Once {
		bd := config.Basedir{}
		sts = rc.UnmarshalFromRequest(&bd)
		if status.IsError(sts) {
			sts = status.Wrap(sts, &status.Args{
				Message:    "unable to update basedir.",
				HttpStatus: http.StatusBadRequest,
				ApiHelp: fmt.Sprintf("verify that API request is in the correct format: %s",
					api.GetApiDocsUrl(rc.RouteName),
				),
			})
			break
		}
		me.Gearbox.SetRouteName(rc.RouteName)
		sts = me.Gearbox.UpdateBasedir(bd.Nickname, bd.HostDir)
	}
	return sts
}
func (me *HostApi) deleteNamedBasedir(rc *api.RequestContext) interface{} {
	me.Gearbox.SetRouteName(rc.RouteName)
	return me.Gearbox.DeleteNamedBasedir(getBasedirNickname(rc))
}

func getBasedirNickname(rc *api.RequestContext) types.Nickname {
	return types.Nickname(rc.Param("nickname"))
}
