package gearbox

import (
	"gearbox/api"
	"gearbox/only"
	"gearbox/stat"
	"net/http"
)

const ProjectsResource api.ResourceName = "projects"
const ProjectsWithDetailsResource api.ResourceName = ProjectsResource + "-with-details"
const ProjectDetailsResource api.ResourceName = "project-details"
const ProjectEnabledResource api.ResourceName = ProjectsResource + "-enabled"
const ProjectDisabledResource api.ResourceName = ProjectsResource + "-disabled"
const ProjectCandidatesResource api.ResourceName = "project-candidates"

const ProjectAliasesResource api.ResourceName = "project-aliases"
const ProjectAliasAdd api.ResourceName = ProjectAliasesResource + "-add"
const ProjectAliasUpdate api.ResourceName = ProjectAliasesResource + "-update"
const ProjectAliasDelete api.ResourceName = ProjectAliasesResource + "-delete"

const ProjectServicesResource api.ResourceName = "project-services"
const ProjectServicesAdd api.ResourceName = ProjectServicesResource + "-add"
const ProjectServicesUpdate api.ResourceName = ProjectServicesResource + "-update"
const ProjectServicesDelete api.ResourceName = ProjectServicesResource + "-delete"

const ProjectStacksResource api.ResourceName = "project-stacks"
const ProjectStackAdd api.ResourceName = ProjectStacksResource + "-add"
const ProjectStackUpdate api.ResourceName = ProjectStacksResource + "-update"
const ProjectStackDelete api.ResourceName = ProjectStacksResource + "-delete"

const HostnameResourceVar api.ResourceVarName = "hostname"

func getProjectHostname(rc *api.RequestContext) (hn string, status stat.Status) {
	for range only.Once {
		hn = rc.Param("hostname")
		if hn == "" {
			status = stat.NewStatus(&stat.Args{
				Message:    "hostname is empty",
				Help:       api.GetApiHelp("hostname"),
				HttpStatus: http.StatusBadRequest,
				Error:      stat.IsStatusError,
			})
			break
		}
	}
	return hn, status
}

func (me *HostApi) addProjectRoutes() {

	me.GET("/projects", ProjectsResource, me.getProjectsResponse, nil)
	me.GET("/projects/with-details", ProjectsWithDetailsResource, me.getProjectsResponse, nil)
	me.GET("/projects/enabled", ProjectEnabledResource, me.getEnabledProjectsResponse, nil)
	me.GET("/projects/disabled", ProjectDisabledResource, me.getDisabledProjectsResponse, nil)
	me.GET("/projects/candidates", ProjectCandidatesResource, me.getCandidateProjectsResponse, nil)

	me.GET("/projects/:hostname", ProjectDetailsResource, me.getProjectDetailsResponse, nil)
	me.PUT("/projects/:hostname", "project-update", me.updateProjectDetails, nil)
	me.POST("/projects/new", "project-add", me.addProjectDetails, nil)
	me.DELETE("/projects/:hostname", "project-delete", me.deleteProjectDetails, nil)

	me.GET("/projects/:hostname/aliases", ProjectAliasesResource, me.getProjectAliasesResponse, nil)
	me.PUT("/projects/:hostname/aliases/:alias", ProjectAliasUpdate, me.updateProjectAlias, nil)
	me.POST("/projects/:hostname/aliases/new", ProjectAliasAdd, me.addProjectAlias, nil)
	me.DELETE("/projects/:hostname/aliases/:alias", ProjectAliasDelete, me.deleteProjectAlias, nil)

	me.GET("/projects/:hostname/services", ProjectServicesResource, me.getProjectServicesResponse, nil)
	me.PUT("/projects/:hostname/services/:service", ProjectServicesUpdate, me.updateProjectService, nil)
	me.POST("/projects/:hostname/services/new", ProjectServicesAdd, me.addProjectService, nil)
	me.DELETE("/projects/:hostname/services/:service", ProjectServicesDelete, me.deleteProjectService, nil)

	me.GET("/projects/:hostname/stacks", ProjectStacksResource, me.getProjectStacksResponse, nil)
	me.PUT("/projects/:hostname/stacks/:stack", ProjectStackUpdate, me.updateProjectStack, nil)
	me.POST("/projects/:hostname/stacks/new", ProjectStackAdd, me.addProjectStack, nil)
	me.DELETE("/projects/:hostname/stacks/:stack", ProjectStackDelete, me.deleteProjectStack, nil)

}

//===[ Project Stacks ]========================

func (me *HostApi) getProjectStacksResponse(rc *api.RequestContext) (response interface{}) {
	var status stat.Status
	for range only.Once {
		hn, status := getProjectHostname(rc)
		if status.IsError() {
			break
		}
		p, status := me.Gearbox.FindProjectWithDetails(hn)
		if status.IsError() {
			break
		}
		response = NewProjectStacksResponse(p)
	}
	if status.IsError() {
		response = status
	}
	return response
}

func (me *HostApi) addProjectStack(rc *api.RequestContext) (response interface{}) {
	var status stat.Status
	for range only.Once {
		var hn string
		hn, status = getProjectHostname(rc)
		if status.IsError() {
			break
		}
		var sn StackName
		sn, status = getStackName(rc)
		if status.IsError() {
			break
		}
		status = me.Gearbox.AddNamedStackToProject(sn, hn)
		if status.IsError() {
			break
		}
	}
	if status.IsError() {
		response = status
	}
	return response
}

func (me *HostApi) updateProjectStack(rc *api.RequestContext) (response interface{}) {
	return me.Api.NotYetImplemented(rc)
}
func (me *HostApi) deleteProjectStack(rc *api.RequestContext) (response interface{}) {
	return me.Api.NotYetImplemented(rc)
}

//===[ Project Services ]========================

func (me *HostApi) getProjectServicesResponse(rc *api.RequestContext) (response interface{}) {
	var status stat.Status
	for range only.Once {
		var hn string
		hn, status := getProjectHostname(rc)
		if status.IsError() {
			break
		}
		p, status := me.Gearbox.FindProjectWithDetails(hn)
		if status.IsError() {
			response = status
			break
		}
		response = NewProjectServicesResponse(p.ServiceMap, p)
	}
	if status.IsError() {
		response = status
	}
	return response
}
func (me *HostApi) addProjectService(rc *api.RequestContext) (response interface{}) {
	return me.Api.NotYetImplemented(rc)
}
func (me *HostApi) updateProjectService(rc *api.RequestContext) (response interface{}) {
	return me.Api.NotYetImplemented(rc)
}
func (me *HostApi) deleteProjectService(rc *api.RequestContext) (response interface{}) {
	return me.Api.NotYetImplemented(rc)
}

//===[ Project Aliases ]========================

func (me *HostApi) getProjectAliasesResponse(rc *api.RequestContext) (response interface{}) {
	var status stat.Status
	for range only.Once {
		var hn string
		hn, status := getProjectHostname(rc)
		if status.IsError() {
			break
		}
		p, status := me.Gearbox.FindProjectWithDetails(hn)
		if status.IsError() {
			break
		}
		response = NewProjectAliasesResponse(p.Aliases, p)
	}
	if status.IsError() {
		response = status
	}
	return response
}

func (me *HostApi) addProjectAlias(rc *api.RequestContext) (response interface{}) {
	return me.Api.NotYetImplemented(rc)
}
func (me *HostApi) updateProjectAlias(rc *api.RequestContext) (response interface{}) {
	return me.Api.NotYetImplemented(rc)
}
func (me *HostApi) deleteProjectAlias(rc *api.RequestContext) (response interface{}) {
	return me.Api.NotYetImplemented(rc)
}

//===[ ProjectMap Aliases ]========================

func (me *HostApi) getProjectsResponse(rc *api.RequestContext) (response interface{}) {
	for range only.Once {
		var status stat.Status
		pm := me.Config.GetProjectMap()
		prs := make(api.ListItemResponseMap, len(pm))
		withDetails := rc.ResourceName == ProjectsWithDetailsResource
		for _, p := range pm {
			if withDetails {
				status = p.MaybeLoadDetails()
				if status.IsError() {
					response = status
					break
				}
			} else {
				p.ClearDetails()
			}
			prs[p.Hostname] = NewProjectListResponse(p)
		}
		if !status.IsError() {
			response = prs
			break
		}

	}
	return response
}

func (me *HostApi) getCandidateProjectsResponse(rc *api.RequestContext) (response interface{}) {
	return me.Config.GetCandidates()
}

func (me *HostApi) getEnabledProjectsResponse(rc *api.RequestContext) (response interface{}) {
	return me.Config.GetProjectMap().GetEnabled()
}

func (me *HostApi) getDisabledProjectsResponse(rc *api.RequestContext) (response interface{}) {
	return me.Config.GetProjectMap().GetDisabled()
}

//===[ Project Details ]========================

func (me *HostApi) getProjectDetailsResponse(rc *api.RequestContext) (response interface{}) {
	var status stat.Status
	for range only.Once {
		var hn string
		hn, status := getProjectHostname(rc)
		if status.IsError() {
			break
		}
		p, status := me.Gearbox.FindProjectWithDetails(hn)
		if status.IsError() {
			break
		}
		response = NewProjectResponse(p)
	}
	if status.IsError() {
		response = status
	}

	return response
}
func (me *HostApi) addProjectDetails(rc *api.RequestContext) (response interface{}) {
	return me.Api.NotYetImplemented(rc)
}
func (me *HostApi) updateProjectDetails(rc *api.RequestContext) (response interface{}) {
	return me.Api.NotYetImplemented(rc)
}
func (me *HostApi) deleteProjectDetails(rc *api.RequestContext) (response interface{}) {
	return me.Api.NotYetImplemented(rc)
}
