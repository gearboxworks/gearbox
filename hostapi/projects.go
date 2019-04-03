package hostapi

import (
	"gearbox/routes"
	"gearbox/status/is"
)

func (me *HostApi) addProjectRoutes() {
	sts := me.AddConnector(routes.NewProjectConnector(me.Gearbox))
	if is.Error(sts) {
		panic(sts.Message())
	}
}

//me.AddConnector("/projects", &ab.ResourceArgs{
//Name:       routes.ProjectsRoute,
//IdParams:   ab.IdParams{"hostname"},
//Item:       routes.ProjectInstance,
//Factory: routes.NewProjectCollection(gb),
//Children: ab.ConnectionsMap{
//	ab.NewResource("/aliases", &ab.ResourceArgs{
//		Name:       routes.ProjectAliasesRoute,
//		IdParams:   ab.IdParams{"alias"},
//		ItemType:   reflect.String,
//	}),
//},
//})

//import (
//	"fmt"
//	"gearbox/gearbox"
//	"gearbox/api"
//	"gearbox/only"
//	"gearbox/routes"
//	"gearbox/status"
//	"net/http"
//)
//
//func getProjectHostname(rc *api.RequestContext) (hn gearbox.Hostname, sts status.Status) {
//	for range only.Once {
//		hn = gearbox.Hostname(rc.Param("hostname"))
//		if hn == "" {
//			sts = status.Fail(&status.Args{
//				Message:    "hostname is empty",
//				Help:       api.GetApiHelp("hostname"),
//				HttpStatus: http.StatusBadRequest,
//			})
//			break
//		}
//	}
//	return hn, sts
//}
//func (me *HostApi) addProjectRoutes() (gb *gearbox.Gearbox) {
//
//	//me.GET___("/projects", routes.ProjectsRoute, me.getProjectsResponse)
//	//me.POST__("/projects/new", routes.ProjectDetailsAdd, me.addProjectDetails)
//	//
//	//me.GET___("/projects/with-details", routes.ProjectsWithDetails, me.getProjectsResponse)
//	//me.GET___("/projects/enabled", routes.EnabledProjects, me.getEnabledProjectsResponse)
//	//me.GET___("/projects/disabled", routes.DisabledProjects, me.getDisabledProjectsResponse)
//	//me.GET___("/projects/candidates", routes.CandidateProjects, me.getCandidateProjectsResponse)
//	//
//	//me.GET___("/projects/:hostname", routes.ProjectDetails, me.getProjectDetailsResponse, me.getProjectHostnameValues)
//	//me.PUT___("/projects/:hostname", routes.ProjectDetailsUpdate, me.updateProjectDetails, me.getProjectHostnameValues)
//	//me.DELETE("/projects/:hostname", routes.ProjectDetailsDelete, me.deleteProjectDetails, me.getProjectHostnameValues)
//	//
//	//me.GET___("/projects/:hostname/aliases", routes.ProjectAliasesRoute, me.getProjectAliasesResponse, me.getProjectHostnameValues)
//	//me.POST__("/projects/:hostname/aliases/new", routes.ProjectAliasAdd, me.addProjectAlias)
//	//me.PUT___("/projects/:hostname/aliases/:alias", routes.ProjectAliasUpdate, me.updateProjectAlias)
//	//me.DELETE("/projects/:hostname/aliases/:alias", routes.ProjectAliasDelete, me.deleteProjectAlias)
//	//
//	//me.Relate(routes.ProjectsRoute, &api.Related{
//	//	List: routes.ProjectsRoute,
//	//	Others: api.RouteNameMap{
//	//		routes.ProjectsWithDetails: api.ListResource,
//	//		routes.EnabledProjects:     api.ListResource,
//	//		routes.DisabledProjects:    api.ListResource,
//	//		routes.CandidateProjects:   api.ListResource,
//	//	},
//	//})
//	//
//	//me.Relate(routes.ProjectDetails, &api.Related{
//	//	Item:   routes.ProjectDetails,
//	//	List:   routes.ProjectsRoute,
//	//	New:    routes.ProjectDetailsAdd,
//	//	Update: routes.ProjectDetailsUpdate,
//	//	Delete: routes.ProjectDetailsDelete,
//	//	Others: api.RouteNameMap{
//	//		routes.ProjectsWithDetails: api.ItemResource,
//	//	},
//	//})
//	//
//	//me.Relate(routes.ProjectAliasesRoute, &api.Related{
//	//	List:   routes.ProjectAliasesRoute,
//	//	New:    routes.ProjectAliasAdd,
//	//	Update: routes.ProjectAliasUpdate,
//	//	Delete: routes.ProjectAliasDelete,
//	//})
//
//	//me.GET___("/projects/:hostname/services", resp.ProjectServicesResource, me.getProjectServicesResponse)
//	//me.PUT___("/projects/:hostname/services/:service", resp.ProjectServicesUpdate, me.updateProjectService)
//	//me.POST__("/projects/:hostname/services/new", resp.ProjectServicesAdd, me.addProjectService, nil)
//	//me.DELETE("/projects/:hostname/services/:service", resp.ProjectServicesDelete, me.deleteProjectService)
//	//
//	//me.GET___("/projects/:hostname/stacks", resp.ProjectStacksResource, me.getProjectStacksResponse)
//	//me.PUT___("/projects/:hostname/stacks/:stack", resp.ProjectStackUpdate, me.updateProjectStack)
//	//me.POST__("/projects/:hostname/stacks/new", resp.ProjectStackAdd, me.addProjectStack)
//	//me.DELETE("/projects/:hostname/stacks/:stack", resp.ProjectStackDelete, me.deleteProjectStack)
//
//}
//
//func (me *HostApi) getProjectHostnameValues(...interface{}) (values api.ValuesFuncValues, sts status.Status) {
//	for range only.Once {
//		var hns gearbox.Hostnames
//		hns, sts = gearbox.GetProjectHostnames(me.Gearbox)
//		if status.IsError(sts) {
//			break
//		}
//		values = make(api.ValuesFuncValues, 1)
//		values[0] = make(api.ValueFuncVarsValues, len(hns))
//		for i, sn := range hns {
//			values[0][i] = string(sn)
//		}
//	}
//	return values, sts
//}
//
//func (me *HostApi) getProjectStacksResponse(rc *api.RequestContext) (response interface{}) {
//	var sts status.Status
//	for range only.Once {
//		hn, sts := getProjectHostname(rc)
//		if status.IsError(sts) {
//			break
//		}
//		p, sts := me.Gearbox.GetProjects(hn)
//		if status.IsError(sts) {
//			break
//		}
//		response = routes.NewProjectStacksResponse(p)
//	}
//	if status.IsError(sts) {
//		response = sts
//	}
//	return response
//}
//
//func (me *HostApi) addProjectStack(rc *api.RequestContext) (response interface{}) {
//	var sts status.Status
//	for range only.Once {
//		var hn gearbox.Hostname
//		hn, sts = getProjectHostname(rc)
//		if status.IsError(sts) {
//			break
//		}
//		var sn gearbox.Stackname
//		sn, sts = routes.GetStackname(rc)
//		if status.IsError(sts) {
//			break
//		}
//		sts = me.Gearbox.AddNamedStackToProject(sn, hn)
//		if status.IsError(sts) {
//			break
//		}
//	}
//	if status.IsError(sts) {
//		response = sts
//	}
//	return response
//}
//
//func (me *HostApi) updateProjectStack(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//func (me *HostApi) deleteProjectStack(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//
//func (me *HostApi) getProjectServicesResponse(rc *api.RequestContext) (response interface{}) {
//	var sts status.Status
//	for range only.Once {
//		var hn gearbox.Hostname
//		hn, sts := getProjectHostname(rc)
//		if status.IsError(sts) {
//			break
//		}
//		p, sts := me.Gearbox.GetProjects(hn)
//		if status.IsError(sts) {
//			response = sts
//			break
//		}
//		response = routes.NewProjectServices(p.Stack, p)
//	}
//	if status.IsError(sts) {
//		response = sts
//	}
//	return response
//}
//func (me *HostApi) addProjectService(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//func (me *HostApi) updateProjectService(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//func (me *HostApi) deleteProjectService(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//
//func (me *HostApi) getProjectAliasesResponse(rc *api.RequestContext) (response interface{}) {
//	var sts status.Status
//	for range only.Once {
//		var hn gearbox.Hostname
//		hn, sts = getProjectHostname(rc)
//		if status.IsError(sts) {
//			break
//		}
//		if hn == "" {
//			sts = status.Fail(&status.Args{
//				Message: fmt.Sprintf("hostname not found for routes name '%s'", rc.RouteName),
//			})
//			break
//		}
//		var p *gearbox.Project
//		p, sts = me.Gearbox.GetProjects(hn)
//		if status.IsError(sts) {
//			break
//		}
//		response = routes.NewProjectAliases(p.Aliases, p)
//	}
//	if status.IsError(sts) {
//		response = sts
//	}
//	return response
//}
//
//func (me *HostApi) addProjectAlias(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//func (me *HostApi) updateProjectAlias(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//func (me *HostApi) deleteProjectAlias(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//
//func (me *HostApi) getProjectsResponse(rc *api.RequestContext) (response interface{}) {
//	var sts status.Status
//	var prs api.ListItemResponseMap
//	for range only.Once {
//		var pm gearbox.ProjectMap
//		pm, sts = me.Gearbox.GetProjectMap()
//		if status.IsError(sts) {
//			break
//		}
//		prs = make(api.ListItemResponseMap, len(pm))
//		withDetails := rc.RouteName == routes.ProjectsWithDetailsFilter
//		for _, p := range pm {
//			if withDetails {
//				sts = p.MaybeLoadDetails()
//				if status.IsError(sts) {
//					response = sts
//					break
//				}
//			} else {
//				p.ClearDetails()
//			}
//			prs[string(p.Hostname)] = routes.NewProjectList(routes.ConvertConfigProject(p))
//		}
//	}
//	if !status.IsError(sts) {
//		response = prs
//	}
//	return response
//}
//
//func (me *HostApi) getCandidateProjectsResponse(rc *api.RequestContext) (response interface{}) {
//	return me.Config.GetCandidates()
//}
//
//func (me *HostApi) getEnabledProjectsResponse(rc *api.RequestContext) (response interface{}) {
//	for range only.Once {
//		pm, sts := me.Gearbox.GetProjectMap()
//		if status.IsError(sts) {
//			response = sts
//			break
//		}
//		response = pm.GetEnabled()
//	}
//	return response
//}
//
//func (me *HostApi) getDisabledProjectsResponse(rc *api.RequestContext) (response interface{}) {
//	for range only.Once {
//		pm, sts := me.Gearbox.GetProjectMap()
//		if status.IsError(sts) {
//			response = sts
//			break
//		}
//		response = pm.GetDisabled()
//	}
//	return response
//}
//
//func (me *HostApi) getProjectDetailsResponse(rc *api.RequestContext) (response interface{}) {
//	var sts status.Status
//	for range only.Once {
//		var hn gearbox.Hostname
//		hn, sts := getProjectHostname(rc)
//		if status.IsError(sts) {
//			break
//		}
//		p, sts := me.Gearbox.GetProjects(hn)
//		if status.IsError(sts) {
//			break
//		}
//		response = routes.ConvertConfigProject(p)
//	}
//	if status.IsError(sts) {
//		response = sts
//	}
//
//	return response
//}
//func (me *HostApi) addProjectDetails(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//func (me *HostApi) updateProjectDetails(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//func (me *HostApi) deleteProjectDetails(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//
//func GetProjectHostnames(gb gearbox.Gearboxer) (hns Hostnames, sts status.Status) {
//	for range only.Once {
//		pm, sts := gb.GetProjectMap()
//		if status.IsError(sts) {
//			break
//		}
//		hns = make(Hostnames, len(pm))
//		i := 0
//		for hn := range pm {
//			hns[i] = hn
//			i++
//		}
//	}
//	if !status.IsError(sts) {
//		sts = status.Success("hostnames gotten")
//	}
//	return hns, sts
//}
//
//
