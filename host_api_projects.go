package gearbox

import (
	"gearbox/api"
	"gearbox/only"
	"github.com/labstack/echo"
)

const ProjectDetailsResource api.ResourceName = "project-details"
const HostnameResourceVar api.ResourceVarName = "hostname"
const ProjectsWithDetailsResource api.ResourceName = "projects-with-details"
const ProjectServicesResource api.ResourceName = "project-services"
const ProjectAliasesResource api.ResourceName = "project-aliases"

func getProjectHostname(ctx echo.Context) string {
	return ctx.Param("hostname")
}

func (me *HostApi) addProjectRoutes() {

	me.GET("/projects", "projects", me.getProjectsResponse)

	me.GET("/projects/with-details", ProjectsWithDetailsResource, me.getProjectsResponse)

	me.GET("/projects/:hostname", ProjectDetailsResource, me.getProjectResponse)

	me.GET("/projects/enabled", "projects-enabled", me.getEnabledProjectsResponse)

	me.GET("/projects/disabled", "projects-disabled", me.getDisabledProjectsResponse)

	me.GET("/projects/candidates", "project-candidates", me.getCandidateProjectsResponse)

	me.GET("/projects/:hostname/services", ProjectServicesResource, me.getProjectServicesResponse)

	me.GET("/projects/:hostname/aliases", ProjectAliasesResource, me.getProjectAliasesResponse)

	me.POST("/projects/:hostname/aliases/new", "project-alias-add", me.addBasedir)

	me.PUT("/projects/:hostname/aliases/:alias", "project-alias-update", me.updateBasedir)

	me.DELETE("/projects/:hostname/aliases/:alias", "project-alias-delete", me.deleteNamedBasedir)

	me.POST("/projects/:hostname/services/new", "project-service-add", me.addBasedir)

	me.PUT("/projects/:hostname/services/:service", "project-service-update", me.updateBasedir)

	me.DELETE("/projects/:hostname/services/:service", "project-service-delete", me.deleteNamedBasedir)

	me.POST("/projects/new", "project-add", me.Api.NotYetImplemented)

	me.POST("/projects/:hostname", "project-update", me.Api.NotYetImplemented)

	me.DELETE("/projects/:hostname", "project-delete", me.Api.NotYetImplemented)
}

func (me *HostApi) getCandidateProjectsResponse(rc *api.RequestContext) (response interface{}) {
	return me.Config.Projects.GetEnabled()
}

func (me *HostApi) getEnabledProjectsResponse(rc *api.RequestContext) (response interface{}) {
	return me.Config.Projects.GetEnabled()
}

func (me *HostApi) getDisabledProjectsResponse(rc *api.RequestContext) (response interface{}) {
	return me.Config.Projects.GetDisabled()
}

func (me *HostApi) getProjectsResponse(rc *api.RequestContext) (response interface{}) {
	for range only.Once {
		me.Gearbox.RequestType = rc.ResourceName
		prs := make(api.ListItemResponseMap, len(me.Config.Projects))
		withDetails := rc.ResourceName == "projects-with-details"
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

func (me *HostApi) getProjectAliasesResponse(rc *api.RequestContext) (response interface{}) {
	for range only.Once {
		me.Gearbox.RequestType = rc.ResourceName
		pr, status := me.Gearbox.GetProjectResponse(getProjectHostname(rc.Context))
		if status.IsError() {
			response = status
			break
		}
		response = NewProjectAliasesResponse(pr.Aliases, pr.Project)
	}
	return response
}

func (me *HostApi) getProjectServicesResponse(rc *api.RequestContext) (response interface{}) {
	for range only.Once {
		me.Gearbox.RequestType = rc.ResourceName
		pr, status := me.Gearbox.GetProjectResponse(getProjectHostname(rc.Context))
		if status.IsError() {
			response = status
			break
		}
		response = NewProjectServicesResponse(pr.ServiceMap, pr.Project)
	}
	return response
}

func (me *HostApi) getProjectResponse(rc *api.RequestContext) (response interface{}) {
	for range only.Once {
		me.Gearbox.RequestType = rc.ResourceName
		pr, status := me.Gearbox.GetProjectResponse(getProjectHostname(rc.Context))
		if status.IsError() {
			response = status
			break
		}
		response = pr
	}
	return response
}
