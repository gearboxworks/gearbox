package gearbox

import (
	"fmt"
	"gearbox/api"
	"gearbox/only"
	"github.com/labstack/echo"
	"net/http"
)

func (me *HostApi) addProjectRoutes() {
	_api := me.Api
	_api.GET("/projects", "projects", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.getProjects(ctx, rt))
	})
	_api.GET("/projects/with-details", "projects-with-details", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.getProjects(ctx, rt))
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

func (me *HostApi) getProjects(ctx echo.Context, requestType string) (response interface{}) {
	for range only.Once {
		me.Gearbox.RequestType = requestType
		projects := make(ProjectMap, len(me.Config.Projects))
		withDetails := requestType == "projects-with-details"
		for _, p := range me.Config.Projects {
			if withDetails {
				p.MaybeLoadDetails()
			} else {
				p.ClearDetails()
			}
			projects[p.Hostname] = p
		}
		response = projects
	}
	return response
}

func (me *HostApi) getProject(ctx echo.Context, requestType string) (response interface{}) {
	for range only.Once {
		me.Gearbox.RequestType = requestType
		pr, status := me.Gearbox.GetProjectResponse(getProjectHostname(ctx))
		if status.IsError() {
			response = status
			break
		}
		response = pr
	}
	return response
}

func getProjectHostname(ctx echo.Context) string {
	return ctx.Param("hostname")
}
