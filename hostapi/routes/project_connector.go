package routes

import (
	"fmt"
	"gearbox"
	"gearbox/apibuilder"
	"gearbox/config"
	"gearbox/only"
	"gearbox/project"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"github.com/labstack/echo"
	"net/http"
	"reflect"
)

const HostnameIdParam ab.IdParam = "hostname"

const ProjectsWithDetailsFilter ab.FilterPath = "/with-details"
const EnabledProjectsFilter ab.FilterPath = "/enabled"
const DisabledProjectsFilter ab.FilterPath = "/disabled"

var NilProjectConnector = (*ProjectConnector)(nil)
var _ ab.Connector = NilProjectConnector

type ProjectConnector struct {
	Gearbox gearbox.Gearboxer
}

func NewProjectConnector(gb gearbox.Gearboxer) *ProjectConnector {
	return &ProjectConnector{
		Gearbox: gb,
	}
}

func (me *ProjectConnector) Related() {
	return
}

func (me *ProjectConnector) GetBasepath() ab.Basepath {
	return "/projects"
}

func (me *ProjectConnector) GetItemType() reflect.Kind {
	return reflect.Struct
}

func (me *ProjectConnector) GetIdParams() ab.IdParams {
	return ab.IdParams{HostnameIdParam}
}

func (me *ProjectConnector) GetCollection(filterPath ab.FilterPath) (collection ab.Collection, sts status.Status) {
	for range only.Once {
		collection = make(ab.Collection, 0)
		gbpm, sts := me.getGearboxProjectMap()
		if is.Error(sts) {
			break
		}
		for _, gbp := range gbpm {
			var rp *Project
			pp, sts := ConvertProject(gbp)
			if is.Error(sts) {
				break
			}
			rp, sts = FilterProject(pp, filterPath)
			if is.Error(sts) {
				break
			}
			if rp == nil {
				continue
			}
			collection = append(collection, rp)
			if is.Error(sts) {
				break
			}
		}
	}
	return collection, sts
}

func (me *ProjectConnector) GetCollectionIds() (itemIds ab.ItemIds, sts status.Status) {
	for range only.Once {
		gbpm, sts := me.getGearboxProjectMap()
		if is.Error(sts) {
			break
		}
		itemIds = make(ab.ItemIds, len(gbpm))
		i := 0
		for _, gbp := range gbpm {
			itemIds[i] = ab.ItemId(gbp.Hostname)
			i++
		}
	}
	return itemIds, sts
}

func (me *ProjectConnector) AddItem(item ab.Item) (collection ab.Collection, sts status.Status) {
	for range only.Once {
		var pp *project.Project
		pp, collection, sts = me.extractGearboxProject(item)
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
	return collection, sts
}

func (me *ProjectConnector) UpdateItem(item ab.Item) (collection ab.Collection, sts status.Status) {
	for range only.Once {
		var pp *project.Project
		pp, collection, sts = me.extractGearboxProject(item)
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
	return collection, sts

}

func (me *ProjectConnector) DeleteItem(hostname ab.ItemId) (collection ab.Collection, sts status.Status) {
	for range only.Once {
		sts := me.Gearbox.DeleteProject(types.Hostname(hostname))
		if status.IsError(sts) {
			break
		}
		sts = status.Success("project '%s' found", hostname)
		sts.SetHttpStatus(http.StatusNoContent)
	}
	return collection, sts
}

func (me *ProjectConnector) GetItem(hostname ab.ItemId, ctx echo.Context) (collection ab.Item, sts status.Status) {
	var p *Project
	for range only.Once {
		gbp, sts := me.Gearbox.FindProject(types.Hostname(hostname))
		if is.Error(sts) {
			break
		}
		if gbp == nil {
			sts = status.Fail(&status.Args{
				Message:    fmt.Sprintf("project '%s' not found", hostname),
				HttpStatus: http.StatusNotFound,
			})
			break
		}
		p, sts = ConvertProject(gbp)
		if is.Error(sts) {
			break
		}
		sts = status.Success("project '%s' found", hostname)
	}
	return p, sts

}

func (me *ProjectConnector) FilterItem(in ab.Item, filterPath ab.FilterPath) (out ab.Item, sts status.Status) {
	for range only.Once {
		var p *Project
		p, sts = AssertProject(in)
		if is.Error(sts) {
			break
		}
		out, sts = FilterProject(p, filterPath)
	}
	return out, sts
}

func (me *ProjectConnector) GetCollectionFilterMap() ab.FilterMap {
	return GetProjectFilterMap()
}

func (me *ProjectConnector) getGearboxProjectMap() (pm project.Map, sts status.Status) {
	for range only.Once {
		pm, sts = me.Gearbox.GetProjectMap()
	}
	return pm, sts
}

func (me *ProjectConnector) extractGearboxProject(item ab.Item) (gbp *project.Project, collection ab.Collection, sts status.Status) {
	var p *Project
	for range only.Once {
		collection, sts = me.GetCollection(ab.NoFilterPath)
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

func GetProjectFilterMap() ab.FilterMap {
	return ab.FilterMap{
		ProjectsWithDetailsFilter: ab.Filter{
			Name: "Projects with Details",
			Path: ProjectsWithDetailsFilter,
			Filter: func(item ab.Item) ab.Item {
				panic(fmt.Sprintf("%s not yet implemented", ProjectsWithDetailsFilter))
				return nil
			},
		},
		EnabledProjectsFilter: ab.Filter{
			Name: "Enabled Projects",
			Path: EnabledProjectsFilter,
			Filter: func(item ab.Item) ab.Item {
				p, sts := AssertProject(item)
				if is.Success(sts) && p.Enabled {
					return item
				}
				return nil
			},
		},
		DisabledProjectsFilter: ab.Filter{
			Name: "Disabled Projects",
			Path: DisabledProjectsFilter,
			Filter: func(item ab.Item) ab.Item {
				p, sts := AssertProject(item)
				if is.Success(sts) && !p.Enabled {
					return item
				}
				return nil
			},
		},
	}
}

func FilterProject(in *Project, filterPath ab.FilterPath) (out *Project, sts status.Status) {
	for range only.Once {
		if filterPath == ab.NoFilterPath {
			out = in
			break
		}
		fm := GetProjectFilterMap()
		f, ok := fm[filterPath]
		if !ok {
			sts = status.Fail(&status.Args{
				Message:    fmt.Sprintf("filter '%s' not found", filterPath),
				HttpStatus: http.StatusBadRequest,
			})
			break
		}
		out, sts = AssertProject(f.Filter(in))
	}
	return out, sts
}

func AssertProject(item ab.Item) (p *Project, sts status.Status) {
	p, ok := item.(*Project)
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("item not a project: %v", item),
		})
	}
	return p, sts
}
