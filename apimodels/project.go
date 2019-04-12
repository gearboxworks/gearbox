package apimodels

import (
	"gearbox/apimodeler"
	"gearbox/config"
	"gearbox/only"
	"gearbox/project"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
)

const ProjectTypeName apimodeler.ItemType = "project"

var ProjectInstance = (*Project)(nil)
var _ apimodeler.ApiItemer = ProjectInstance

type Project struct {
	Hostname      types.Hostname          `json:"hostname"`
	Enabled       bool                    `json:"enabled"`
	Basedir       types.Nickname          `json:"basedir"`
	Notes         string                  `json:"notes"`
	Path          types.RelativePath      `json:"path"`
	ProjectDir    types.AbsoluteDir       `json:"project_dir"`
	Filepath      types.AbsoluteFilepath  `json:"filepath"`
	Aliases       project.HostnameAliases `json:"aliases,omitempty"`
	Services      Services                `json:"stack,omitempty"`
	ConfigProject *config.Project         `json:"-"`
}

func NewProject(hostname apimodeler.ItemId) *Project {
	return &Project{
		Hostname: types.Hostname(hostname),
	}
}

func (me *Project) GetType() apimodeler.ItemType {
	return ProjectTypeName
}

func (me *Project) GetId() apimodeler.ItemId {
	return apimodeler.ItemId(me.Hostname)
}

func (me *Project) SetId(hostname apimodeler.ItemId) status.Status {
	me.Hostname = types.Hostname(hostname)
	return nil
}

func (me *Project) GetItem() (apimodeler.ApiItemer, status.Status) {
	return me, nil
}

func (me *Project) GetItemLinkMap(*apimodeler.Context) (lm apimodeler.LinkMap, sts status.Status) {
	return apimodeler.LinkMap{
		//apimodeler.RelatedRelType: apimodeler.Link("https://example.com"),
	}, sts
}

func (me *Project) GetRelatedItems(ctx *apimodeler.Context, itemid apimodeler.ItemId) (list apimodeler.List, sts status.Status) {
	list = make(apimodeler.List, 0)
	return list, sts
}

func (me *Project) AddDetails(ctx *apimodeler.Context) (sts status.Status) {
	for range only.Once {
		pp := project.NewProject(me.ConfigProject)
		sts = pp.Load()
		if is.Error(sts) {
			break
		}
		me.Aliases = pp.Aliases
		me.Filepath = pp.Filepath
		me.Services, sts = GetFromServiceStackMap(ctx, pp.GetServiceMap())
		if is.Error(sts) {
			break
		}
	}
	return sts
}

func NewFromConfigProject(cp *config.Project) (p *Project, sts status.Status) {
	for range only.Once {
		pd, sts := cp.GetDir()
		if is.Error(sts) {
			break
		}
		p = &Project{
			Hostname:      cp.Hostname,
			Basedir:       cp.Basedir,
			Notes:         cp.Notes,
			Path:          cp.Path,
			ProjectDir:    pd,
			ConfigProject: cp,
		}
	}
	return p, sts
}
