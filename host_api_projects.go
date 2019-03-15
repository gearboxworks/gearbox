package gearbox

import (
	"fmt"
	"gearbox/api"
	"gearbox/only"
	"gearbox/stat"
	"github.com/labstack/echo"
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

	me.GET("/projects", ProjectsResource, me.getProjectsResponse)
	me.GET("/projects/with-details", ProjectsWithDetailsResource, me.getProjectsResponse)
	me.GET("/projects/enabled", ProjectEnabledResource, me.getEnabledProjectsResponse)
	me.GET("/projects/disabled", ProjectDisabledResource, me.getDisabledProjectsResponse)
	me.GET("/projects/candidates", ProjectCandidatesResource, me.getCandidateProjectsResponse)

	me.GET("/projects/:hostname", ProjectDetailsResource, me.getProjectDetailsResponse)
	me.PUT("/projects/:hostname", "project-update", me.updateProjectDetails)
	me.POST("/projects/new", "project-add", me.addProjectDetails)
	me.DELETE("/projects/:hostname", "project-delete", me.deleteProjectDetails)

	me.GET("/projects/:hostname/aliases", ProjectAliasesResource, me.getProjectAliasesResponse)
	me.PUT("/projects/:hostname/aliases/:alias", ProjectAliasUpdate, me.updateProjectAlias)
	me.POST("/projects/:hostname/aliases/new", ProjectAliasAdd, me.addProjectAlias)
	me.DELETE("/projects/:hostname/aliases/:alias", ProjectAliasDelete, me.deleteProjectAlias)

	me.GET("/projects/:hostname/services", ProjectServicesResource, me.getProjectServicesResponse)
	me.PUT("/projects/:hostname/services/:service", ProjectServicesUpdate, me.updateProjectService)
	me.POST("/projects/:hostname/services/new", ProjectServicesAdd, me.addProjectService)
	me.DELETE("/projects/:hostname/services/:service", ProjectServicesDelete, me.deleteProjectService)

	me.GET("/projects/:hostname/stacks", ProjectStacksResource, me.getProjectStacksResponse)
	me.PUT("/projects/:hostname/stacks/:stack", ProjectStackUpdate, me.updateProjectStack)
	me.POST("/projects/:hostname/stacks/new", ProjectStackAdd, me.addProjectStack)
	me.DELETE("/projects/:hostname/stacks/:stack", ProjectStackDelete, me.deleteProjectStack)

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

func getStackName(rc *api.RequestContext) (sn StackName, status stat.Status) {
	for range only.Once {
		if rc.Context.Request().Method == echo.GET {
			sn = StackName(rc.Param("stack"))
			break
		}
		snr := StackNameRequest{}
		status := rc.UnmarshalFromRequest(&snr)
		if status.IsError() {
			status.Status = status
			status.Message = fmt.Sprintf("invalid request format for '%s' resource", rc.ResourceName)
			status.HttpStatus = http.StatusBadRequest
			status.ApiHelp = api.GetApiHelp("rc.ResourceName", "correct request format")
			break
		}
		sn = snr.StackName
	}
	if sn == "" {
		status = stat.NewStatus(&stat.Args{
			Message:    "stack name is empty",
			Help:       api.GetApiHelp(rc.ResourceName),
			HttpStatus: http.StatusBadRequest,
			Error:      stat.IsStatusError,
		})
	}
	return sn, status
}

type StackNameRequest struct {
	StackName StackName `json:"stack"`
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

//===[ Projects Aliases ]========================

func (me *HostApi) getProjectsResponse(rc *api.RequestContext) (response interface{}) {
	for range only.Once {
		prs := make(api.ListItemResponseMap, len(me.Config.Projects))
		withDetails := rc.ResourceName == ProjectsWithDetailsResource
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

func (me *HostApi) getCandidateProjectsResponse(rc *api.RequestContext) (response interface{}) {
	return me.Config.Candidates
}

func (me *HostApi) getEnabledProjectsResponse(rc *api.RequestContext) (response interface{}) {
	return me.Config.Projects.GetEnabled()
}

func (me *HostApi) getDisabledProjectsResponse(rc *api.RequestContext) (response interface{}) {
	return me.Config.Projects.GetDisabled()
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
