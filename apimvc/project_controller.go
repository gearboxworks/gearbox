package apimvc

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/config"
	"gearbox/gearbox"
	"gearbox/gearspec"
	"gearbox/only"
	"gearbox/project"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"net/http"
	"reflect"
)

const HostnameIdParam apimodeler.IdParam = "hostname"

const ProjectsName types.RouteName = "projects"
const ProjectsBasepath types.Basepath = "/projects"

const ProjectsWithDetailsFilter apimodeler.FilterPath = "/with-details"
const EnabledProjectsFilter apimodeler.FilterPath = "/enabled"
const DisabledProjectsFilter apimodeler.FilterPath = "/disabled"

var NilProjectController = (*ProjectController)(nil)
var _ apimodeler.ApiController = NilProjectController

type ProjectController struct {
	apimodeler.Controller
	Gearbox gearbox.Gearboxer
}

func (me *ProjectController) GetName() types.RouteName {
	return ProjectsName
}

func (me *ProjectController) GetListLinkMap(*apimodeler.Context, ...apimodeler.FilterPath) (lm apimodeler.LinkMap, sts status.Status) {
	return apimodeler.LinkMap{
		//apimodeler.RelatedRelType: apimodeler.Link("http://example.org/"),
	}, sts
}

func NewProjectController(gb gearbox.Gearboxer) *ProjectController {
	return &ProjectController{
		Gearbox: gb,
	}
}

func (me *ProjectController) Related() {
	return
}

func (me *ProjectController) GetBasepath() types.Basepath {
	return ProjectsBasepath
}

func (me *ProjectController) GetItemType() reflect.Kind {
	return reflect.Struct
}

func (me *ProjectController) GetIdParams() apimodeler.IdParams {
	return apimodeler.IdParams{HostnameIdParam}
}

func (me *ProjectController) GetList(ctx *apimodeler.Context, filterPath ...apimodeler.FilterPath) (list apimodeler.List, sts status.Status) {
	//var fp apimodeler.FilterPath
	//if len(filterPath) > 0 {
	//	fp = filterPath[0]
	//} else {
	//	fp = apimodeler.NoFilterPath
	//}
	for range only.Once {
		list = make(apimodeler.List, 0)
		cpm, sts := me.Gearbox.GetConfig().GetProjectMap()
		if is.Error(sts) {
			break
		}
		for _, cp := range cpm {
			pp, sts := NewModelFromConfigProject(cp)
			if is.Error(sts) {
				break
			}
			list = append(list, pp)
			if is.Error(sts) {
				break
			}
		}
	}
	return list, sts
}

func (me *ProjectController) FilterList(ctx *apimodeler.Context, filterPath apimodeler.FilterPath) (list apimodeler.List, sts status.Status) {
	for range only.Once {
		list, sts := me.GetList(ctx, filterPath)
		if is.Error(sts) {
			break
		}
		for i, item := range list {
			item, sts = me.FilterItem(item, filterPath)
			if is.Error(sts) {
				break
			}
			if item == nil {
				continue
			}
			list[i] = item
		}
	}
	return list, sts
}

func (me *ProjectController) GetListIds(ctx *apimodeler.Context, filterPath ...apimodeler.FilterPath) (itemids apimodeler.ItemIds, sts status.Status) {
	for range only.Once {
		gbpm, sts := me.getGearboxProjectMap()
		if is.Error(sts) {
			break
		}
		itemids = make(apimodeler.ItemIds, len(gbpm))
		i := 0
		for _, gbp := range gbpm {
			itemids[i] = apimodeler.ItemId(gbp.Hostname)
			i++
		}
	}
	return itemids, sts
}

func (me *ProjectController) AddItem(ctx *apimodeler.Context, item apimodeler.Itemer) (sts status.Status) {
	for range only.Once {
		var pp *project.Project
		pp, _, sts = me.extractGearboxProject(ctx, item)
		if status.IsError(sts) {
			break
		}
		sts = me.Gearbox.AddProject(pp)
		if status.IsError(sts) {
			break
		}
		sts = status.Success("project '%s' added", pp.Hostname)
		sts.SetHttpStatus(http.StatusCreated)
	}
	return sts
}

func (me *ProjectController) UpdateItem(ctx *apimodeler.Context, item apimodeler.Itemer) (sts status.Status) {
	for range only.Once {
		var pp *project.Project
		pp, _, sts = me.extractGearboxProject(ctx, item)
		if status.IsError(sts) {
			break
		}
		sts = me.Gearbox.UpdateProject(pp)
		if status.IsError(sts) {
			break
		}
		sts = status.Success("project '%s' updated", item.GetId())
		sts.SetHttpStatus(http.StatusNoContent)
	}
	return sts

}

func (me *ProjectController) DeleteItem(ctx *apimodeler.Context, hostname apimodeler.ItemId) (sts status.Status) {
	for range only.Once {
		sts := me.Gearbox.DeleteProject(types.Hostname(hostname))
		if status.IsError(sts) {
			break
		}
		sts = status.Success("project '%s' found", hostname)
		sts.SetHttpStatus(http.StatusNoContent)
	}
	return sts
}

func (me *ProjectController) GetItem(ctx *apimodeler.Context, hostname apimodeler.ItemId) (list apimodeler.Itemer, sts status.Status) {
	var p *ProjectModel
	for range only.Once {
		cp, sts := me.Gearbox.GetConfig().FindProject(types.Hostname(hostname))
		if is.Error(sts) {
			break
		}
		if cp == nil {
			sts = status.Fail(&status.Args{
				Message:    fmt.Sprintf("project '%s' not found", hostname),
				HttpStatus: http.StatusNotFound,
			})
			break
		}
		p, sts = NewModelFromConfigProject(cp)
		if is.Error(sts) {
			break
		}
		sts = status.Success("project '%s' found", hostname)
	}
	return p, sts
}

func (me *ProjectController) GetItemDetails(ctx *apimodeler.Context, itemid apimodeler.ItemId) (item apimodeler.Itemer, sts status.Status) {
	for range only.Once {
		item, sts = me.GetItem(ctx, itemid)
		if is.Error(sts) {
			break
		}
		p, ok := item.(*ProjectModel)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("item '%s' not a project.Project", itemid),
			})
			break
		}
		sts = p.AddDetails(ctx)
		if is.Error(sts) {
			break
		}
	}
	return item, sts
}

func (me *ProjectController) GetRelatedItems(ctx *apimodeler.Context, item apimodeler.Itemer) (list apimodeler.List, sts status.Status) {
	var p *ProjectModel
	for range only.Once {
		var ok bool
		p, ok = item.(*ProjectModel)
		if !ok {
			sts = status.OurBad("item '%s' does not implement *apimodeler.Project", item.GetId())
			return
		}
		list = make(apimodeler.List, 0)
		for _, s := range p.Services {
			gsgs := gearspec.NewGearspec()
			sts = gsgs.Parse(s.GearspecId)
			if is.Error(sts) {
				break
			}
			gs, sts := NewModelFromGearspecGearspec(ctx, gsgs)
			if is.Error(sts) {
				break
			}
			list = append(list, gs)
		}
		noop(p)
	}
	return list, nil
}

func (me *ProjectController) FilterItem(in apimodeler.Itemer, filterPath apimodeler.FilterPath) (out apimodeler.Itemer, sts status.Status) {
	for range only.Once {
		if filterPath == apimodeler.NoFilterPath {
			out = in
			break
		}
		fm := me.GetFilterMap()
		f, ok := fm[filterPath]
		if !ok {
			sts = status.Fail(&status.Args{
				Message:    fmt.Sprintf("filter '%s' not found", filterPath),
				HttpStatus: http.StatusBadRequest,
			})
			break
		}
		out, sts = AssertProject(f.ItemFilter(in))
	}
	return out, sts
}

func (me *ProjectController) FilterProject(in *ProjectModel, filterPath apimodeler.FilterPath) (out *ProjectModel, sts status.Status) {
	for range only.Once {
	}
	return out, sts
}

func (me *ProjectController) GetFilterMap() apimodeler.FilterMap {
	return apimodeler.FilterMap{
		ProjectsWithDetailsFilter: apimodeler.Filter{
			Label: "Projects with Details",
			Path:  ProjectsWithDetailsFilter,
			ItemFilter: func(item apimodeler.Itemer) apimodeler.Itemer {
				panic(fmt.Sprintf("%s not yet implemented", ProjectsWithDetailsFilter))
				return nil
			},
		},
		EnabledProjectsFilter: apimodeler.Filter{
			Label: "Enabled Projects",
			Path:  EnabledProjectsFilter,
			ItemFilter: func(item apimodeler.Itemer) apimodeler.Itemer {
				p, sts := AssertProject(item)
				if is.Success(sts) && p.Enabled {
					return item
				}
				return nil
			},
		},
		DisabledProjectsFilter: apimodeler.Filter{
			Label: "Disabled Projects",
			Path:  DisabledProjectsFilter,
			ItemFilter: func(item apimodeler.Itemer) apimodeler.Itemer {
				p, sts := AssertProject(item)
				if is.Success(sts) && !p.Enabled {
					return item
				}
				return nil
			},
		},
	}
}

func (me *ProjectController) getGearboxProjectMap() (pm project.Map, sts status.Status) {
	for range only.Once {
		pm, sts = me.Gearbox.GetProjectMap()
	}
	return pm, sts
}

func (me *ProjectController) extractGearboxProject(ctx *apimodeler.Context, item apimodeler.Itemer) (gbp *project.Project, list apimodeler.List, sts status.Status) {
	var p *ProjectModel
	for range only.Once {
		list, sts = me.GetList(ctx)
		if is.Error(sts) {
			break
		}
		p, sts = AssertProject(item)
		if is.Error(sts) {
			break
		}
		gbp, sts = MakeGearboxProject(me.Gearbox, p)
	}
	return gbp, list, sts
}

func MakeGearboxProject(gb gearbox.Gearboxer, prj *ProjectModel) (pp *project.Project, sts status.Status) {
	for range only.Once {
		cp := config.NewProject(gb.GetConfig(), prj.Path)
		pp = project.NewProject(cp)
	}
	return pp, sts
}

func AssertProject(item apimodeler.Itemer) (p *ProjectModel, sts status.Status) {
	p, ok := item.(*ProjectModel)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item not a project: %v", item),
		})
	}
	return p, sts
}
