package api

//func (me *Api) addProjectRoutes() {
//	gb,ok := me.Parent.(gearbox.Gearboxer)
//	if !ok {
//
//	}
//	sts := me.AddModels(models.NewProjectModel(gb))
//	if is.Error(sts) {
//		panic(sts.Message())
//	}
//}

//me.AddModels("/projects", &ab.ResourceArgs{
//Program:       models.ProjectsRoute,
//IdParams:   ab.IdParams{"hostname"},
//Item:       models.ProjectInstance,
//Factory: models.NewProjectCollection(gb),
//Children: ab.ModelsMap{
//	ab.NewResource("/aliases", &ab.ResourceArgs{
//		Program:       models.ProjectAliasesRoute,
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
//	"gearbox/api/models"
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
//	//me.GET___("/projects", models.ProjectsRoute, me.getProjectsResponse)
//	//me.POST__("/projects/new", models.ProjectDetailsAdd, me.addProjectDetails)
//	//
//	//me.GET___("/projects/with-details", models.ProjectsWithDetails, me.getProjectsResponse)
//	//me.GET___("/projects/enabled", models.EnabledProjects, me.getEnabledProjectsResponse)
//	//me.GET___("/projects/disabled", models.DisabledProjects, me.getDisabledProjectsResponse)
//	//me.GET___("/projects/candidates", models.CandidateProjects, me.getCandidateProjectsResponse)
//	//
//	//me.GET___("/projects/:hostname", models.ProjectDetails, me.getProjectDetailsResponse, me.getProjectHostnameValues)
//	//me.PUT___("/projects/:hostname", models.ProjectDetailsUpdate, me.updateProjectDetails, me.getProjectHostnameValues)
//	//me.DELETE("/projects/:hostname", models.ProjectDetailsDelete, me.deleteProjectDetails, me.getProjectHostnameValues)
//	//
//	//me.GET___("/projects/:hostname/aliases", models.ProjectAliasesRoute, me.getProjectAliasesResponse, me.getProjectHostnameValues)
//	//me.POST__("/projects/:hostname/aliases/new", models.ProjectAliasAdd, me.addProjectAlias)
//	//me.PUT___("/projects/:hostname/aliases/:alias", models.ProjectAliasUpdate, me.updateProjectAlias)
//	//me.DELETE("/projects/:hostname/aliases/:alias", models.ProjectAliasDelete, me.deleteProjectAlias)
//	//
//	//me.Relate(models.ProjectsRoute, &api.Related{
//	//	List: models.ProjectsRoute,
//	//	Others: api.RouteNameMap{
//	//		models.ProjectsWithDetails: api.ListResource,
//	//		models.EnabledProjects:     api.ListResource,
//	//		models.DisabledProjects:    api.ListResource,
//	//		models.CandidateProjects:   api.ListResource,
//	//	},
//	//})
//	//
//	//me.Relate(models.ProjectDetails, &api.Related{
//	//	Item:   models.ProjectDetails,
//	//	List:   models.ProjectsRoute,
//	//	New:    models.ProjectDetailsAdd,
//	//	Update: models.ProjectDetailsUpdate,
//	//	Delete: models.ProjectDetailsDelete,
//	//	Others: api.RouteNameMap{
//	//		models.ProjectsWithDetails: api.ItemResource,
//	//	},
//	//})
//	//
//	//me.Relate(models.ProjectAliasesRoute, &api.Related{
//	//	List:   models.ProjectAliasesRoute,
//	//	New:    models.ProjectAliasAdd,
//	//	Update: models.ProjectAliasUpdate,
//	//	Delete: models.ProjectAliasDelete,
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
//		response = models.NewProjectStacksResponse(p)
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
//		sn, sts = models.GetStackname(rc)
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
//		response = models.NewProjectServices(p.Stack, p)
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
//				Message: fmt.Sprintf("hostname not found for models name '%s'", rc.RouteName),
//			})
//			break
//		}
//		var p *gearbox.Project
//		p, sts = me.Parent.GetProjects(hn)
//		if status.IsError(sts) {
//			break
//		}
//		response = models.NewProjectAliases(p.Aliases, p)
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
//		withDetails := rc.RouteName == models.ProjectsWithDetailsFilter
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
//			prs[string(p.Hostname)] = models.NewProjectList(models.ConvertConfigProject(p))
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
//		response = models.ConvertConfigProject(p)
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
