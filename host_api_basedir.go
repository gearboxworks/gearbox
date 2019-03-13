package gearbox

import (
	"fmt"
	"gearbox/api"
	"gearbox/only"
	"github.com/labstack/echo"
	"net/http"
)

func (me *HostApi) addBasedirRoutes() {
	_api := me.Api

	_api.GET("/basedirs", "basedirs", func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.getHostBasedirs(ctx, rt))
	})

	_api.POST("/basedirs/new", "basedir-add", func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.addBasedir(ctx, rt))
	})

	_api.PUT("/basedirs/:nickname", "basedir-update", func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.updateBasedir(ctx, rt))
	})

	_api.DELETE("/basedirs/:nickname", "basedir-delete", func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.deleteNamedBasedir(ctx, rt))
	})

}

func (me *HostApi) getHostBasedirs(ctx echo.Context, resource api.ResourceName) interface{} {
	return me.Config.GetHostBasedirs()
}

func (me *HostApi) addBasedir(ctx echo.Context, resource api.ResourceName) interface{} {
	var status Status
	for range only.Once {
		bd := Basedir{}
		err := api.UnmarshalFromRequest(resource, ctx, &bd)
		if err != nil {
			status = NewStatus(&StatusArgs{
				Failed:     true,
				Message:    "Unable to add basedir.",
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    fmt.Sprintf("verify that API request is in the correct format: %s", api.GetApiDocsUrl(resource)),
				Error:      err,
			})
			break
		}
		me.Gearbox.RequestType = resource
		status = me.Gearbox.AddBasedir(bd.HostDir, bd.Nickname)
	}
	return status
}
func (me *HostApi) updateBasedir(ctx echo.Context, resource api.ResourceName) interface{} {
	var status Status
	for range only.Once {
		bd := Basedir{}
		err := api.UnmarshalFromRequest(resource, ctx, &bd)
		if err != nil {
			status = NewStatus(&StatusArgs{
				Failed:     true,
				Message:    "Unable to update basedir.",
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    fmt.Sprintf("verify that API request is in the correct format: %s", api.GetApiDocsUrl(resource)),
				Error:      err,
			})
			break
		}
		me.Gearbox.RequestType = resource
		status = me.Gearbox.UpdateBasedir(bd.Nickname, bd.HostDir)
	}
	return status
}
func (me *HostApi) deleteNamedBasedir(ctx echo.Context, resource api.ResourceName) interface{} {
	me.Gearbox.RequestType = resource
	return me.Gearbox.DeleteNamedBasedir(getBasedirNickname(ctx))
}

func getBasedirNickname(ctx echo.Context) string {
	return ctx.Param("nickname")
}
