package gearbox

import (
	"gearbox/api"
	"gearbox/only"
)

const ProjectsResource api.ResourceName = "projects"
const ProjectDetailsResource api.ResourceName = "project-details"
const ProjectEnabledResource api.ResourceName = "project-enabled"
const ProjectDisabledResource api.ResourceName = "project-disabled"
const ProjectCandidatesResource api.ResourceName = "project-candidates"
const ProjectsWithDetailsResource api.ResourceName = "projects-with-details"

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

func getProjectHostname(rc *api.RequestContext) string {
	return rc.Param("hostname")
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
	me.PUT("/projects/:hostname/stacks/:alias", ProjectStackUpdate, me.updateProjectStack)
	me.POST("/projects/:hostname/stacks/new", ProjectStackAdd, me.addProjectStack)
	me.DELETE("/projects/:hostname/stacks/:alias", ProjectStackDelete, me.deleteProjectStack)

}

//===[ Project Details ]========================

func (me *HostApi) getProjectDetailsResponse(rc *api.RequestContext) (response interface{}) {
	for range only.Once {
		me.Gearbox.RequestType = rc.ResourceName
		p, status := me.Gearbox.FindProjectWithDetails(getProjectHostname(rc))
		if status.IsError() {
			response = status
			break
		}
		response = NewProjectResponse(p)
	}
	return response
}
func (me *HostApi) addProjectDetails(rc *api.RequestContext) (response interface{}) {
	return nil
}
func (me *HostApi) updateProjectDetails(rc *api.RequestContext) (response interface{}) {
	return nil
}
func (me *HostApi) deleteProjectDetails(rc *api.RequestContext) (response interface{}) {
	return nil
}

//===[ Project Stacks ]========================

func (me *HostApi) getProjectStacksResponse(rc *api.RequestContext) (response interface{}) {
	for range only.Once {
		me.Gearbox.RequestType = rc.ResourceName
		p, status := me.Gearbox.FindProjectWithDetails(getProjectHostname(rc))
		if status.IsError() {
			response = status
			break
		}
		response = NewProjectStacksResponse(p)
	}
	return response
}
func (me *HostApi) addProjectStack(rc *api.RequestContext) (response interface{}) {
	return nil
}
func (me *HostApi) updateProjectStack(rc *api.RequestContext) (response interface{}) {
	return nil
}
func (me *HostApi) deleteProjectStack(rc *api.RequestContext) (response interface{}) {
	return nil
}

//===[ Project Services ]========================

func (me *HostApi) getProjectServicesResponse(rc *api.RequestContext) (response interface{}) {
	for range only.Once {
		me.Gearbox.RequestType = rc.ResourceName
		p, status := me.Gearbox.FindProjectWithDetails(getProjectHostname(rc))
		if status.IsError() {
			response = status
			break
		}
		response = NewProjectServicesResponse(p.ServiceMap, p)
	}
	return response
}
func (me *HostApi) addProjectService(rc *api.RequestContext) (response interface{}) {
	return nil
}
func (me *HostApi) updateProjectService(rc *api.RequestContext) (response interface{}) {
	return nil
}
func (me *HostApi) deleteProjectService(rc *api.RequestContext) (response interface{}) {
	return nil
}

//===[ Project Aliases ]========================

func (me *HostApi) getProjectAliasesResponse(rc *api.RequestContext) (response interface{}) {
	for range only.Once {
		me.Gearbox.RequestType = rc.ResourceName
		p, status := me.Gearbox.FindProjectWithDetails(getProjectHostname(rc))
		if status.IsError() {
			response = status
			break
		}
		response = NewProjectAliasesResponse(p.Aliases, p)
	}
	return response
}

func (me *HostApi) addProjectAlias(rc *api.RequestContext) (response interface{}) {
	return nil
}
func (me *HostApi) updateProjectAlias(rc *api.RequestContext) (response interface{}) {
	return nil
}
func (me *HostApi) deleteProjectAlias(rc *api.RequestContext) (response interface{}) {
	return nil
}

//===[ Projects Aliases ]========================

func (me *HostApi) getProjectsResponse(rc *api.RequestContext) (response interface{}) {
	for range only.Once {
		me.Gearbox.RequestType = rc.ResourceName
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
	return me.Config.Projects.GetEnabled()
}

func (me *HostApi) getEnabledProjectsResponse(rc *api.RequestContext) (response interface{}) {
	return me.Config.Projects.GetEnabled()
}

func (me *HostApi) getDisabledProjectsResponse(rc *api.RequestContext) (response interface{}) {
	return me.Config.Projects.GetDisabled()
}
