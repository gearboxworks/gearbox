package gearbox

import (
	"fmt"
	"gearbox/api"
	"gearbox/only"
	"github.com/labstack/echo"
	"net/http"
)

const ProjectDetailsResource api.ResourceName = "project-details"
const HostnameResourceVar api.ResourceVarName = "hostname"
const ProjectsWithDetailsResource api.ResourceName = "projects-with-details"
const ProjectServicesResource api.ResourceName = "projects-services"
const ProjectAliasesResource api.ResourceName = "projects-aliases"

func getProjectHostname(ctx echo.Context) string {
	return ctx.Param("hostname")
}

func (me *HostApi) addProjectRoutes() {
	_api := me.Api
	_api.GET("/projects", "projects", func(rt api.ResourceName, ctx echo.Context) error {

		return me.jsonMarshalHandler(_api, ctx, rt, me.getProjectsResponse(ctx, rt))
	})
	_api.GET("/projects/with-details", ProjectsWithDetailsResource, func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.getProjectsResponse(ctx, rt))
	})
	_api.GET("/projects/:hostname", ProjectDetailsResource, func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.getProjectResponse(ctx, rt))
	})
	_api.GET("/projects/enabled", "projects-enabled", func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.Config.Projects.GetEnabled())
	})
	_api.GET("/projects/disabled", "projects-disabled", func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.Config.Projects.GetDisabled())
	})
	_api.GET("/projects/candidates", "project-candidates", func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.Config.Candidates)
	})

	_api.GET("/projects/:hostname/services", ProjectServicesResource, func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.getProjectServicesResponse(ctx, rt))
	})

	_api.GET("/projects/:hostname/aliases", ProjectAliasesResource, func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.getProjectAliasesResponse(ctx, rt))
	})

	_api.POST("/projects/:hostname/services/new", "project-service-add", func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.addBasedir(me.Gearbox, ctx, rt))
	})

	_api.PUT("/projects/:hostname/services/:service", "project-service-update", func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.updateBasedir(me.Gearbox, ctx, rt))
	})

	_api.DELETE("/projects/:hostname/services/:service", "project-service-delete", func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.deleteNamedBasedir(me.Gearbox, ctx, rt))
	})

	_api.POST("/projects/new", "project-add", func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, &api.Status{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      fmt.Errorf("the 'project-add' method has not been implemented yet"),
		})
	})

	_api.POST("/projects/:hostname", "project-update", func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, &api.Status{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      fmt.Errorf("the 'project-update' method has not been implemented yet"),
		})
	})

	_api.DELETE("/projects/:hostname", "project-delete", func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, &api.Status{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      fmt.Errorf("the 'project-delete' method has not been implemented yet"),
		})
	})
}

func (me *HostApi) getProjectsResponse(ctx echo.Context, requestType api.ResourceName) (response interface{}) {
	for range only.Once {
		me.Gearbox.RequestType = requestType
		prs := make(api.ListItemResponseMap, len(me.Config.Projects))
		withDetails := requestType == "projects-with-details"
		for _, p := range me.Config.Projects {
			if withDetails {
				p.MaybeLoadDetails()
			} else {
				p.ClearDetails()
			}
			prs[p.Hostname] = NewProjectListResponse(p)
		}
		response = prs
	}
	return response
}

func (me *HostApi) getProjectAliasesResponse(ctx echo.Context, requestType api.ResourceName) (response interface{}) {
	for range only.Once {
		me.Gearbox.RequestType = requestType
		pr, status := me.Gearbox.GetProjectResponse(getProjectHostname(ctx))
		if status.IsError() {
			response = status
			break
		}
		response = NewProjectAliasesResponse(pr.Aliases, pr.Project)
	}
	return response
}

func (me *HostApi) getProjectServicesResponse(ctx echo.Context, requestType api.ResourceName) (response interface{}) {
	for range only.Once {
		me.Gearbox.RequestType = requestType
		pr, status := me.Gearbox.GetProjectResponse(getProjectHostname(ctx))
		if status.IsError() {
			response = status
			break
		}
		response = NewProjectServicesResponse(pr.ServiceMap, pr.Project)
	}
	return response
}

func (me *HostApi) getProjectResponse(ctx echo.Context, requestType api.ResourceName) (response interface{}) {
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
