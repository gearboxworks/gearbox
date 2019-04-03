package routes

import (
	"gearbox/apibuilder"
	"gearbox/only"
	"gearbox/project"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
)

const ProjectTypeName = "project"

var ProjectInstance = (*Project)(nil)
var _ ab.Item = ProjectInstance

type Project struct {
	Hostname   ab.ItemId               `json:"hostname"`
	Enabled    bool                    `json:"enabled"`
	Basedir    types.Nickname          `json:"basedir"`
	Notes      string                  `json:"notes"`
	Path       types.RelativePath      `json:"path"`
	ProjectDir types.AbsoluteDir       `json:"project_dir"`
	Aliases    project.HostnameAliases `json:"aliases,omitempty"`
	Services   Services                `json:"stack,omitempty"`
}

func NewProject(hostname ab.ItemId) *Project {
	return &Project{
		Hostname: hostname,
	}
}

func (me *Project) GetType() ab.ItemType {
	return ProjectTypeName
}

func (me *Project) GetId() ab.ItemId {
	return me.Hostname
}

func (me *Project) GetItem() (ab.Item, status.Status) {
	return me, nil
}

func ConvertProject(pp *project.Project) (p *Project, sts status.Status) {
	for range only.Once {
		pd, sts := pp.GetDir()
		if is.Error(sts) {
			break
		}
		p = &Project{
			Hostname: ab.ItemId(pp.Hostname),
			Basedir:  pp.Basedir,
			Notes:    pp.Notes,
			Path:     pp.Path,

			//Stack: ConvertServiceMap(p.GetServiceMap()),
			Services: ConvertServices(pp.GetServiceMap()),

			Aliases:    pp.GetAliases(),
			ProjectDir: pd,
		}
	}
	return p, sts
}
