package models

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/config"
	"gearbox/gearbox"
	"gearbox/jsonapi"
	"gearbox/only"
	"gearbox/project"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"net/http"
	"reflect"
)

const HostnameIdParam apimodeler.IdParam = "hostname"

const ProjectsBasepath types.Basepath = "/projects"

const ProjectsWithDetailsFilter apimodeler.FilterPath = "/with-details"
const EnabledProjectsFilter apimodeler.FilterPath = "/enabled"
const DisabledProjectsFilter apimodeler.FilterPath = "/disabled"

var NilProjectModel = (*ProjectModel)(nil)
var _ apimodeler.Modeler = NilProjectModel

type ProjectModel struct {
	Gearbox gearbox.Gearboxer
}

func NewProjectModel(gb gearbox.Gearboxer) *ProjectModel {
	return &ProjectModel{
		Gearbox: gb,
	}
}

func (me *ProjectModel) Related() {
	return
}

func (me *ProjectModel) GetBasepath() types.Basepath {
	return ProjectsBasepath
}

func (me *ProjectModel) GetItemType() reflect.Kind {
	return reflect.Struct
}

func (me *ProjectModel) GetIdParams() apimodeler.IdParams {
	return apimodeler.IdParams{HostnameIdParam}
}

func (me *ProjectModel) GetCollection(ctx apimodeler.Contexter, filterPath ...apimodeler.FilterPath) (collection apimodeler.Collection, sts status.Status) {
	//var fp apimodeler.FilterPath
	//if len(filterPath) > 0 {
	//	fp = filterPath[0]
	//} else {
	//	fp = apimodeler.NoFilterPath
	//}
	for range only.Once {
		collection = make(apimodeler.Collection, 0)
		cpm, sts := me.Gearbox.GetConfig().GetProjectMap()
		if is.Error(sts) {
			break
		}
		for _, cp := range cpm {
			pp, sts := ConvertProject(cp)
			if is.Error(sts) {
				break
			}
			collection = append(collection, pp)
			if is.Error(sts) {
				break
			}
		}
	}
	return collection, sts
}

func (me *ProjectModel) FilterCollection(ctx apimodeler.Contexter, filterPath apimodeler.FilterPath) (collection apimodeler.Collection, sts status.Status) {
	for range only.Once {
		collection, sts := me.GetCollection(ctx, filterPath)
		if is.Error(sts) {
			break
		}
		for i, item := range collection {
			item, sts = me.FilterItem(item, filterPath)
			if is.Error(sts) {
				break
			}
			if item == nil {
				continue
			}
			collection[i] = item
		}
	}
	return collection, sts
}

func (me *ProjectModel) GetCollectionIds(ctx apimodeler.Contexter) (itemIds apimodeler.ItemIds, sts status.Status) {
	for range only.Once {
		gbpm, sts := me.getGearboxProjectMap()
		if is.Error(sts) {
			break
		}
		itemIds = make(apimodeler.ItemIds, len(gbpm))
		i := 0
		for _, gbp := range gbpm {
			itemIds[i] = apimodeler.ItemId(gbp.Hostname)
			i++
		}
	}
	return itemIds, sts
}

func (me *ProjectModel) AddItem(ctx apimodeler.Contexter, item apimodeler.Itemer) (sts status.Status) {
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

func (me *ProjectModel) UpdateItem(ctx apimodeler.Contexter, item apimodeler.Itemer) (sts status.Status) {
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

func (me *ProjectModel) DeleteItem(ctx apimodeler.Contexter, hostname apimodeler.ItemId) (sts status.Status) {
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

func (me *ProjectModel) GetItem(ctx apimodeler.Contexter, hostname apimodeler.ItemId) (collection apimodeler.Itemer, sts status.Status) {
	var p *Project
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
		p, sts = ConvertProject(cp)
		if is.Error(sts) {
			break
		}
		if ctx.Get(ja.ResponseTypeKey) == ja.DatasetResponse {
			sts = p.AddDetails()
			if is.Error(sts) {
				break
			}
		}
		sts = status.Success("project '%s' found", hostname)
	}
	return p, sts

}

func (me *ProjectModel) FilterItem(in apimodeler.Itemer, filterPath apimodeler.FilterPath) (out apimodeler.Itemer, sts status.Status) {
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

func (me *ProjectModel) FilterProject(in *Project, filterPath apimodeler.FilterPath) (out *Project, sts status.Status) {
	for range only.Once {
	}
	return out, sts
}

func (me *ProjectModel) GetFilterMap() apimodeler.FilterMap {
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

func (me *ProjectModel) getGearboxProjectMap() (pm project.Map, sts status.Status) {
	for range only.Once {
		pm, sts = me.Gearbox.GetProjectMap()
	}
	return pm, sts
}

func (me *ProjectModel) extractGearboxProject(ctx apimodeler.Contexter, item apimodeler.Itemer) (gbp *project.Project, collection apimodeler.Collection, sts status.Status) {
	var p *Project
	for range only.Once {
		collection, sts = me.GetCollection(ctx)
		if is.Error(sts) {
			break
		}
		p, sts = AssertProject(item)
		if is.Error(sts) {
			break
		}
		gbp, sts = MakeGearboxProject(me.Gearbox, p)
	}
	return gbp, collection, sts
}

func MakeGearboxProject(gb gearbox.Gearboxer, prj *Project) (pp *project.Project, sts status.Status) {
	for range only.Once {
		cp := config.NewProject(gb.GetConfig(), prj.Path)
		pp = project.NewProject(cp)
	}
	return pp, sts
}

func AssertProject(item apimodeler.Itemer) (p *Project, sts status.Status) {
	p, ok := item.(*Project)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item not a project: %v", item),
		})
	}
	return p, sts
}
