package api

//func (me *Api) addProjectRoutes() {
//	gb,ok := me.Parent.(gearbox.Gearboxer)
//	if !ok {
//
//	}
//	sts := me.AddModels(apimodels.NewProjectModel(gb))
//	if is.Error(sts) {
//		panic(sts.Message())
//	}
//}

//me.AddModels("/projects", &ab.ResourceArgs{
//Program:       apimodels.ProjectsRoute,
//IdParams:   ab.IdParams{"hostname"},
//Item:       apimodels.ProjectInstance,
//Factory: apimodels.NewProjectCollection(gb),
//Children: ab.ModelsMap{
//	ab.NewResource("/aliases", &ab.ResourceArgs{
//		Program:       apimodels.ProjectAliasesRoute,
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
//	"gearbox/apimodels"
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
//				Help:       help.GetApiHelp("hostname"),
//				HttpStatus: http.StatusBadRequest,
//			})
//			break
//		}
//	}
//	return hn, sts
//}
//func (me *Api) addProjectRoutes() (gb *gearbox.Parent) {
//
//	//me.GET___("/projects", apimodels.ProjectsRoute, me.getProjectsResponse)
//	//me.POST__("/projects/new", apimodels.ProjectDetailsAdd, me.addProjectDetails)
//	//
//	//me.GET___("/projects/with-details", apimodels.ProjectsWithDetails, me.getProjectsResponse)
//	//me.GET___("/projects/enabled", apimodels.EnabledProjects, me.getEnabledProjectsResponse)
//	//me.GET___("/projects/disabled", apimodels.DisabledProjects, me.getDisabledProjectsResponse)
//	//me.GET___("/projects/candidates", apimodels.CandidateProjects, me.getCandidateProjectsResponse)
//	//
//	//me.GET___("/projects/:hostname", apimodels.ProjectDetails, me.getProjectDetailsResponse, me.getProjectHostnameValues)
//	//me.PUT___("/projects/:hostname", apimodels.ProjectDetailsUpdate, me.updateProjectDetails, me.getProjectHostnameValues)
//	//me.DELETE("/projects/:hostname", apimodels.ProjectDetailsDelete, me.deleteProjectDetails, me.getProjectHostnameValues)
//	//
//	//me.GET___("/projects/:hostname/aliases", apimodels.ProjectAliasesRoute, me.getProjectAliasesResponse, me.getProjectHostnameValues)
//	//me.POST__("/projects/:hostname/aliases/new", apimodels.ProjectAliasAdd, me.addProjectAlias)
//	//me.PUT___("/projects/:hostname/aliases/:alias", apimodels.ProjectAliasUpdate, me.updateProjectAlias)
//	//me.DELETE("/projects/:hostname/aliases/:alias", apimodels.ProjectAliasDelete, me.deleteProjectAlias)
//	//
//	//me.Relate(apimodels.ProjectsRoute, &api.Related{
//	//	List: apimodels.ProjectsRoute,
//	//	Others: api.RouteNameMap{
//	//		apimodels.ProjectsWithDetails: api.ListResource,
//	//		apimodels.EnabledProjects:     api.ListResource,
//	//		apimodels.DisabledProjects:    api.ListResource,
//	//		apimodels.CandidateProjects:   api.ListResource,
//	//	},
//	//})
//	//
//	//me.Relate(apimodels.ProjectDetails, &api.Related{
//	//	Item:   apimodels.ProjectDetails,
//	//	List:   apimodels.ProjectsRoute,
//	//	New:    apimodels.ProjectDetailsAdd,
//	//	Update: apimodels.ProjectDetailsUpdate,
//	//	Delete: apimodels.ProjectDetailsDelete,
//	//	Others: api.RouteNameMap{
//	//		apimodels.ProjectsWithDetails: api.ItemResource,
//	//	},
//	//})
//	//
//	//me.Relate(apimodels.ProjectAliasesRoute, &api.Related{
//	//	List:   apimodels.ProjectAliasesRoute,
//	//	New:    apimodels.ProjectAliasAdd,
//	//	Update: apimodels.ProjectAliasUpdate,
//	//	Delete: apimodels.ProjectAliasDelete,
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
//func (me *Api) getProjectHostnameValues(...interface{}) (values api.ValuesFuncValues, sts status.Status) {
//	for range only.Once {
//		var hns gearbox.Hostnames
//		hns, sts = gearbox.GetProjectHostnames(me.Parent)
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
//func (me *Api) getProjectStacksResponse(rc *api.RequestContext) (response interface{}) {
//	var sts status.Status
//	for range only.Once {
//		hn, sts := getProjectHostname(rc)
//		if status.IsError(sts) {
//			break
//		}
//		p, sts := me.Parent.GetProjects(hn)
//		if status.IsError(sts) {
//			break
//		}
//		response = apimodels.NewProjectStacksResponse(p)
//	}
//	if status.IsError(sts) {
//		response = sts
//	}
//	return response
//}
//
//func (me *Api) addProjectStack(rc *api.RequestContext) (response interface{}) {
//	var sts status.Status
//	for range only.Once {
//		var hn gearbox.Hostname
//		hn, sts = getProjectHostname(rc)
//		if status.IsError(sts) {
//			break
//		}
//		var sn gearbox.Stackname
//		sn, sts = apimodels.GetStackname(rc)
//		if status.IsError(sts) {
//			break
//		}
//		sts = me.Parent.AddNamedStackToProject(sn, hn)
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
//func (me *Api) updateProjectStack(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//func (me *Api) deleteProjectStack(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//
//func (me *Api) getProjectServicesResponse(rc *api.RequestContext) (response interface{}) {
//	var sts status.Status
//	for range only.Once {
//		var hn gearbox.Hostname
//		hn, sts := getProjectHostname(rc)
//		if status.IsError(sts) {
//			break
//		}
//		p, sts := me.Parent.GetProjects(hn)
//		if status.IsError(sts) {
//			response = sts
//			break
//		}
//		response = apimodels.NewProjectServices(p.Stack, p)
//	}
//	if status.IsError(sts) {
//		response = sts
//	}
//	return response
//}
//func (me *Api) addProjectService(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//func (me *Api) updateProjectService(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//func (me *Api) deleteProjectService(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//
//func (me *Api) getProjectAliasesResponse(rc *api.RequestContext) (response interface{}) {
//	var sts status.Status
//	for range only.Once {
//		var hn gearbox.Hostname
//		hn, sts = getProjectHostname(rc)
//		if status.IsError(sts) {
//			break
//		}
//		if hn == "" {
//			sts = status.Fail(&status.Args{
//				Message: fmt.Sprintf("hostname not found for apimodels name '%s'", rc.RouteName),
//			})
//			break
//		}
//		var p *gearbox.Project
//		p, sts = me.Parent.GetProjects(hn)
//		if status.IsError(sts) {
//			break
//		}
//		response = apimodels.NewProjectAliases(p.Aliases, p)
//	}
//	if status.IsError(sts) {
//		response = sts
//	}
//	return response
//}
//
//func (me *Api) addProjectAlias(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//func (me *Api) updateProjectAlias(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//func (me *Api) deleteProjectAlias(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//
//func (me *Api) getProjectsResponse(rc *api.RequestContext) (response interface{}) {
//	var sts status.Status
//	var prs api.ListItemResponseMap
//	for range only.Once {
//		var pm gearbox.ProjectMap
//		pm, sts = me.Parent.GetProjectMap()
//		if status.IsError(sts) {
//			break
//		}
//		prs = make(api.ListItemResponseMap, len(pm))
//		withDetails := rc.RouteName == apimodels.ProjectsWithDetailsFilter
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
//			prs[string(p.Hostname)] = apimodels.NewProjectList(apimodels.ConvertConfigProject(p))
//		}
//	}
//	if !status.IsError(sts) {
//		response = prs
//	}
//	return response
//}
//
//func (me *Api) getCandidateProjectsResponse(rc *api.RequestContext) (response interface{}) {
//	return me.Config.GetCandidates()
//}
//
//func (me *Api) getEnabledProjectsResponse(rc *api.RequestContext) (response interface{}) {
//	for range only.Once {
//		pm, sts := me.Parent.GetProjectMap()
//		if status.IsError(sts) {
//			response = sts
//			break
//		}
//		response = pm.GetEnabled()
//	}
//	return response
//}
//
//func (me *Api) getDisabledProjectsResponse(rc *api.RequestContext) (response interface{}) {
//	for range only.Once {
//		pm, sts := me.Parent.GetProjectMap()
//		if status.IsError(sts) {
//			response = sts
//			break
//		}
//		response = pm.GetDisabled()
//	}
//	return response
//}
//
//func (me *Api) getProjectDetailsResponse(rc *api.RequestContext) (response interface{}) {
//	var sts status.Status
//	for range only.Once {
//		var hn gearbox.Hostname
//		hn, sts := getProjectHostname(rc)
//		if status.IsError(sts) {
//			break
//		}
//		p, sts := me.Parent.GetProjects(hn)
//		if status.IsError(sts) {
//			break
//		}
//		response = apimodels.ConvertConfigProject(p)
//	}
//	if status.IsError(sts) {
//		response = sts
//	}
//
//	return response
//}
//func (me *Api) addProjectDetails(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//func (me *Api) updateProjectDetails(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//func (me *Api) deleteProjectDetails(rc *api.RequestContext) (response interface{}) {
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
