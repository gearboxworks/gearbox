package gearbox

import (
	"encoding/json"
	"fmt"
	"gearbox/api"
	"gearbox/only"
	"github.com/labstack/echo"
	"net/http"
)

func (me *HostApi) addBasedir(ctx echo.Context, requestType api.ResourceName) (status Status) {
	bd := Basedir{}
	status = readBasedirFromRequest(requestType, ctx, &bd)
	if !status.IsError() {
		me.Gearbox.RequestType = requestType
		status = me.Gearbox.AddBasedir(bd.HostDir, bd.Nickname)
	}
	return status
}
func (me *HostApi) updateBasedir(ctx echo.Context, requestType api.ResourceName) (status Status) {
	bd := Basedir{}
	status = readBasedirFromRequest(requestType, ctx, &bd)
	if !status.IsError() {
		me.Gearbox.RequestType = requestType
		status = me.Gearbox.UpdateBasedir(bd.Nickname, bd.HostDir)
	}
	return status
}
func (me *HostApi) deleteNamedBasedir(ctx echo.Context, requestType api.ResourceName) (status Status) {
	me.Gearbox.RequestType = requestType
	return me.Gearbox.DeleteNamedBasedir(getBasedirNickname(ctx))
}

func (me *HostApi) addBasedirRoutes() {
	_api := me.Api

	_api.GET("/basedirs", "basedirs", func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.Config.GetHostBasedirs())
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

func readBasedirFromRequest(name api.ResourceName, ctx echo.Context, bd *Basedir) (status Status) {
	for range only.Once {
		apiHelp := GetApiHelp(name)
		defer api.CloseRequestBody(ctx)
		b, err := api.ReadRequestBody(ctx)
		if err != nil {
			status = NewStatus(&StatusArgs{
				Message:    "could not read request body",
				HttpStatus: http.StatusUnprocessableEntity,
				ApiHelp:    apiHelp,
				Error:      err,
			})
			break
		}
		err = json.Unmarshal(b, bd)
		if err != nil {
			status = NewStatus(&StatusArgs{
				Message:    fmt.Sprintf("unexpected format for request body: '%s'", string(b)),
				HttpStatus: http.StatusUnprocessableEntity,
				ApiHelp:    apiHelp,
				Error:      err,
			})
			break
		}
		status = NewOkStatus("read %d bytes from body of '%s' request",
			len(b),
			name,
		)
	}
	return status
}

func getBasedirNickname(ctx echo.Context) string {
	return ctx.Param("nickname")
}
