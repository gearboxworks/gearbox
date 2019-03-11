package gearbox

import (
	"fmt"
	"gearbox/api"
	"github.com/labstack/echo"
	"net/http"
)

func (me *HostApi) addProjectRoutes() {
	_api := me.Api
	_api.GET("/projects", "projects", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.Config.Projects)
	})
	_api.GET("/projects/:hostname", "project-details", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.getProject(ctx, rt))
	})
	_api.GET("/projects/enabled", "projects-enabled", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.Config.Projects.GetEnabled())
	})
	_api.GET("/projects/disabled", "projects-disabled", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.Config.Projects.GetDisabled())
	})
	_api.GET("/projects/candidates", "project-candidates", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.Config.Candidates)
	})

	_api.POST("/projects/new", "project-add", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, &api.Status{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      fmt.Errorf("the 'project-add' method has not been implemented yet"),
		})
	})

	_api.POST("/projects/:hostname", "project-update", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, &api.Status{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      fmt.Errorf("the 'project-update' method has not been implemented yet"),
		})
	})

	_api.DELETE("/projects/:hostname", "project-delete", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, &api.Status{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      fmt.Errorf("the 'project-delete' method has not been implemented yet"),
		})
	})
}
