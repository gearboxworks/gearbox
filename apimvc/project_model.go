package apimvc

import (
	"gearbox/apiworks"
	"gearbox/config"
	"gearbox/gearspec"
	"gearbox/only"
	"gearbox/project"
	"gearbox/service"
	"gearbox/types"
	"github.com/gearboxworks/go-status/is"
)

const ProjectModelType ItemType = "projects"

var NilProjectModel = (*ProjectModel)(nil)
var _ ItemModeler = NilProjectModel

type ProjectModelMap map[types.Hostname]*ProjectModel
type ProjectModels []*ProjectModel
type ProjectModel struct {
	Hostname      types.Hostname          `json:"hostname"`
	Enabled       bool                    `json:"enabled"`
	Basedir       types.Nickname          `json:"basedir"`
	Notes         string                  `json:"notes"`
	Path          types.Path              `json:"path"`
	ProjectDir    types.Dir               `json:"project_dir"`
	Filepath      types.Filepath          `json:"filepath"`
	Aliases       project.HostnameAliases `json:"aliases,omitempty"`
	Stack         ProjectStackItems       `json:"stack,omitempty"`
	ConfigProject *config.Project         `json:"-"`
	Model
}

func (me *ProjectModel) GetAttributeMap() apiworks.AttributeMap {
	panic("implement me")
}

func NewModelFromConfigProject(cp *config.Project) (p *ProjectModel, sts Status) {
	for range only.Once {
		pd, sts := cp.GetDir()
		if is.Error(sts) {
			break
		}
		fp, sts := cp.GetFilepath()
		if is.Error(sts) {
			break
		}
		p = &ProjectModel{
			Hostname:      cp.Hostname,
			Basedir:       cp.Basedir,
			Notes:         cp.Notes,
			Path:          cp.Path,
			Filepath:      fp,
			ProjectDir:    pd,
			ConfigProject: cp,
		}
	}
	return p, sts
}

func NewProjectModel(hostname ItemId) *ProjectModel {
	return &ProjectModel{
		Hostname: types.Hostname(hostname),
	}
}

func (me *ProjectModel) GetType() ItemType {
	return ProjectModelType
}

func (me *ProjectModel) GetId() ItemId {
	return ItemId(me.Hostname)
}

func (me *ProjectModel) SetId(hostname ItemId) Status {
	me.Hostname = types.Hostname(hostname)
	return nil
}

func (me *ProjectModel) AddDetails(ctx *Context) (sts Status) {
	for range only.Once {
		pp := project.NewProject(me.ConfigProject)
		sts = pp.Load()
		if is.Error(sts) {
			break
		}
		me.Aliases = pp.Aliases
		me.Filepath = pp.Filepath
		var sms ServiceModels
		sms, sts = GetServiceModelsFromServiceServicerMap(ctx, pp.GetServicerMap())
		if is.Error(sts) {
			break
		}
		me.Stack = make(ProjectStackItems, len(sms))
		for i, sm := range sms {
			me.Stack[i] = NewProjectStackItemFromServiceModel(sm)
		}
	}
	return sts
}

func (me *ProjectModel) GetRelatedItems(ctx *Context) (list List, sts Status) {
	for range only.Once {
		list = make(List, 0)
		for _, s := range me.Stack {
			gsgs := gearspec.NewGearspec()
			sts = gsgs.Parse(s.GearspecId)
			if is.Error(sts) {
				break
			}
			gsm, sts := NewGearspecModelFromGearspecGearspec(ctx, gsgs)
			if is.Error(sts) {
				break
			}
			list = append(list, gsm)

			ss := service.NewService()
			sts = ss.Parse(s.ServiceId)
			if is.Error(sts) {
				break
			}
			sm, sts := NewModelFromServiceServicer(ctx, ss)
			if is.Error(sts) {
				break
			}
			//sm.GearspecId = gsm.GearspecId
			list = append(list, sm)
		}
	}
	return list, sts
}

//var ProjectsInstance = (Projects)(nil)
//var _ ab.ItemCollection = ProjectsInstance
//
//
//func (me Projects) GetFilterMap() ab.FilterMap {
//	return GetProjectFilterMap()
//}
//
//func (me Projects) GetCollectionItemIds() (ab.ItemIds, Status) {
//	itemIds := make(ab.ItemIds, len(me))
//	for i, p := range me {
//		itemIds[i] = p.GetIdentifier()
//	}
//	return itemIds, nil
//}
//
//func (me Projects) AddItem(item ab.ItemInstance) (collection ab.ItemCollection, sts Status) {
//	found := false
//	collection = me
//	var project *Project
//	for range only.Once {
//		item, sts = me.GetItem(item.GetIdentifier())
//		if !is.Error(sts) {
//			sts = status.Fail(&status.Args{
//				Cause:      sts,
//				Message:    fmt.Sprintf("project '%s' already exists", item.GetIdentifier()),
//				HttpStatus: http.StatusConflict,
//			})
//		}
//		project, sts = AssertProject(item)
//		if is.Error(sts) {
//			break
//		}
//		found = true
//		break
//	}
//	if !found {
//		collection = append(me, project)
//	}
//	return collection, sts
//}
//
//func (me Projects) UpdateItem(item ab.ItemInstance) (collection ab.ItemCollection, sts Status) {
//	updated := false
//	collection = me
//	for range only.Once {
//		item, sts = me.GetItem(item.GetIdentifier())
//		if is.Error(sts) {
//			break
//		}
//		project, sts := AssertProject(item)
//		if is.Error(sts) {
//			break
//		}
//		collection, sts = collection.AddItem(project)
//		if is.Error(sts) {
//			break
//		}
//		updated = true
//	}
//	if !updated {
//		sts = status.Fail(&status.Args{
//			Message:    fmt.Sprintf("project '%s' not found", item.GetIdentifier()),
//			HttpStatus: http.StatusNotFound,
//		})
//	}
//	return collection, nil
//}
//
//func (me Projects) DeleteItem(id ab.ItemId) (collection ab.ItemCollection, sts Status) {
//	deleted := false
//	collection = me
//	for i, p := range me {
//		if id != p.GetIdentifier() {
//			continue
//		}
//		collection = append(me[:i], me[i+1:]...)
//		deleted = true
//		break
//	}
//	if !deleted {
//		sts = status.Fail(&status.Args{
//			Message:    fmt.Sprintf("project '%s' not found", id),
//			HttpStatus: http.StatusNotFound,
//		})
//	}
//	return collection, nil
//}
//
//func (me Projects) GetItem(id ab.ItemId) (item ab.ItemInstance, sts Status) {
//	found := false
//	for _, p := range me {
//		if id != p.GetIdentifier() {
//			continue
//		}
//		item = p
//		found = true
//		break
//	}
//	if !found {
//		sts = status.Fail(&status.Args{
//			Message:    fmt.Sprintf("project '%s' not found", id),
//			HttpStatus: http.StatusNotFound,
//		})
//	}
//	return item, sts
//}
//
//func (me Projects) GetItemCollection(filterPath ab.FilterPath) (collection ab.ItemCollection, sts Status) {
//	collection = make(Projects, len(me))
//	for _, p := range me {
//		p, sts = FilterProject(p, filterPath)
//		if p == nil {
//			continue
//		}
//		if is.Error(sts) {
//			break
//		}
//		collection, sts = collection.AddItem(p)
//		if is.Error(sts) {
//			break
//		}
//	}
//	return collection, sts
//}
//
//func (me Projects) GetCollectionSlice(filterPath ab.FilterPath) (slice ab.ItemInstances, sts Status) {
//	slice = make(ab.ItemInstances, len(me))
//	for i, p := range me {
//		p, sts = FilterProject(p, filterPath)
//		if p == nil {
//			continue
//		}
//		if is.Error(sts) {
//			break
//		}
//		slice[i] = p
//	}
//	return slice, sts
//}
//

//func (me *Api) addProjectRoutes() {
//	gb,ok := me.Gearbox.(gearbox.Gearboxer)
//	if !ok {
//
//	}
//	sts := me.AddController(apimvc.NewProjectController(gb))
//	if is.Error(sts) {
//		panic(sts.Message())
//	}
//}

//me.AddController("/projects", &ab.ResourceArgs{
//Program:       apimvc.ProjectsRoute,
//IdParams:   ab.IdParams{"hostname"},
//Item:       apimvc.NilProjectModel,
//Factory: apimvc.NewProjectCollection(gb),
//Children: ab.ControllerMap{
//	ab.NewResource("/aliases", &ab.ResourceArgs{
//		Program:       apimvc.ProjectAliasesRoute,
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
//	"gearbox/apimvc"
//	"github.com/gearboxworks/go-status"
//	"net/http"
//)
//
//func getProjectHostname(rc *api.RequestContext) (hn gearbox.Hostname, sts Status) {
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
//func (me *Api) addProjectRoutes() (gb *gearbox.Gearbox) {
//
//	//me.GET___("/projects", apimvc.ProjectsRoute, me.getProjectsResponse)
//	//me.POST__("/projects/new", apimvc.ProjectDetailsAdd, me.addProjectDetails)
//	//
//	//me.GET___("/projects/with-details", apimvc.ProjectsWithDetails, me.getProjectsResponse)
//	//me.GET___("/projects/enabled", apimvc.EnabledProjects, me.getEnabledProjectsResponse)
//	//me.GET___("/projects/disabled", apimvc.DisabledProjects, me.getDisabledProjectsResponse)
//	//me.GET___("/projects/candidates", apimvc.CandidateProjects, me.getCandidateProjectsResponse)
//	//
//	//me.GET___("/projects/:hostname", apimvc.ProjectDetails, me.getProjectDetailsResponse, me.getProjectHostnameValues)
//	//me.PUT___("/projects/:hostname", apimvc.ProjectDetailsUpdate, me.updateProjectDetails, me.getProjectHostnameValues)
//	//me.DELETE("/projects/:hostname", apimvc.ProjectDetailsDelete, me.deleteProjectDetails, me.getProjectHostnameValues)
//	//
//	//me.GET___("/projects/:hostname/aliases", apimvc.ProjectAliasesRoute, me.getProjectAliasesResponse, me.getProjectHostnameValues)
//	//me.POST__("/projects/:hostname/aliases/new", apimvc.ProjectAliasAdd, me.addProjectAlias)
//	//me.PUT___("/projects/:hostname/aliases/:alias", apimvc.ProjectAliasUpdate, me.updateProjectAlias)
//	//me.DELETE("/projects/:hostname/aliases/:alias", apimvc.ProjectAliasDelete, me.deleteProjectAlias)
//	//
//	//me.Relate(apimvc.ProjectsRoute, &api.RelatedField{
//	//	List: apimvc.ProjectsRoute,
//	//	Others: api.RouteNameMap{
//	//		apimvc.ProjectsWithDetails: api.ListResource,
//	//		apimvc.EnabledProjects:     api.ListResource,
//	//		apimvc.DisabledProjects:    api.ListResource,
//	//		apimvc.CandidateProjects:   api.ListResource,
//	//	},
//	//})
//	//
//	//me.Relate(apimvc.ProjectDetails, &api.RelatedField{
//	//	Item:   apimvc.ProjectDetails,
//	//	List:   apimvc.ProjectsRoute,
//	//	New:    apimvc.ProjectDetailsAdd,
//	//	Update: apimvc.ProjectDetailsUpdate,
//	//	Delete: apimvc.ProjectDetailsDelete,
//	//	Others: api.RouteNameMap{
//	//		apimvc.ProjectsWithDetails: api.ItemResource,
//	//	},
//	//})
//	//
//	//me.Relate(apimvc.ProjectAliasesRoute, &api.RelatedField{
//	//	List:   apimvc.ProjectAliasesRoute,
//	//	New:    apimvc.ProjectAliasAdd,
//	//	Update: apimvc.ProjectAliasUpdate,
//	//	Delete: apimvc.ProjectAliasDelete,
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
//func (me *Api) getProjectHostnameValues(...interface{}) (values api.ValuesFuncValues, sts Status) {
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
//func (me *Api) getProjectStacksResponse(rc *api.RequestContext) (response interface{}) {
//	var sts Status
//	for range only.Once {
//		hn, sts := getProjectHostname(rc)
//		if status.IsError(sts) {
//			break
//		}
//		p, sts := me.Gearbox.GetProjects(hn)
//		if status.IsError(sts) {
//			break
//		}
//		response = apimvc.NewProjectStacksResponse(p)
//	}
//	if status.IsError(sts) {
//		response = sts
//	}
//	return response
//}
//
//func (me *Api) addProjectStack(rc *api.RequestContext) (response interface{}) {
//	var sts Status
//	for range only.Once {
//		var hn gearbox.Hostname
//		hn, sts = getProjectHostname(rc)
//		if status.IsError(sts) {
//			break
//		}
//		var sn gearbox.Stackname
//		sn, sts = apimvc.GetStackname(rc)
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
//func (me *Api) updateProjectStack(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//func (me *Api) deleteProjectStack(rc *api.RequestContext) (response interface{}) {
//	return me.Api.NotYetImplemented(rc)
//}
//
//func (me *Api) getProjectServicesResponse(rc *api.RequestContext) (response interface{}) {
//	var sts Status
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
//		response = apimvc.NewProjectServices(p.Stack, p)
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
//	var sts Status
//	for range only.Once {
//		var hn gearbox.Hostname
//		hn, sts = getProjectHostname(rc)
//		if status.IsError(sts) {
//			break
//		}
//		if hn == "" {
//			sts = status.Fail(&status.Args{
//				Message: fmt.Sprintf("hostname not found for apimvc name '%s'", rc.RouteName),
//			})
//			break
//		}
//		var p *gearbox.Project
//		p, sts = me.Gearbox.GetProjects(hn)
//		if status.IsError(sts) {
//			break
//		}
//		response = apimvc.NewProjectAliases(p.Aliases, p)
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
//	var sts Status
//	var prs api.ListItemResponseMap
//	for range only.Once {
//		var pm gearbox.ProjectMap
//		pm, sts = me.Gearbox.GetProjectMap()
//		if status.IsError(sts) {
//			break
//		}
//		prs = make(api.ListItemResponseMap, len(pm))
//		withDetails := rc.RouteName == apimvc.ProjectsWithDetailsFilter
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
//			prs[string(p.Hostname)] = apimvc.NewProjectList(apimvc.ConvertConfigProject(p))
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
//func (me *Api) getDisabledProjectsResponse(rc *api.RequestContext) (response interface{}) {
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
//func (me *Api) getProjectDetailsResponse(rc *api.RequestContext) (response interface{}) {
//	var sts Status
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
//		response = apimvc.ConvertConfigProject(p)
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
//func GetProjectHostnames(gb gearbox.Gearboxer) (hns Hostnames, sts Status) {
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
