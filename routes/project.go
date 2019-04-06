package routes

import (
	"gearbox/modeler"
	"gearbox/only"
	"gearbox/project"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
)

const ProjectTypeName = "project"

var ProjectInstance = (*Project)(nil)
var _ modeler.Item = ProjectInstance

type Project struct {
	Hostname   modeler.ItemId          `json:"hostname"`
	Enabled    bool                    `json:"enabled"`
	Basedir    types.Nickname          `json:"basedir"`
	Notes      string                  `json:"notes"`
	Path       types.RelativePath      `json:"path"`
	ProjectDir types.AbsoluteDir       `json:"project_dir"`
	Aliases    project.HostnameAliases `json:"aliases,omitempty"`
	Services   Services                `json:"stack,omitempty"`
}

func NewProject(hostname modeler.ItemId) *Project {
	return &Project{
		Hostname: hostname,
	}
}

func (me *Project) GetType() modeler.ItemType {
	return ProjectTypeName
}

func (me *Project) GetId() modeler.ItemId {
	return me.Hostname
}

func (me *Project) GetItem() (modeler.Item, status.Status) {
	return me, nil
}

func ConvertProject(pp *project.Project) (p *Project, sts status.Status) {
	for range only.Once {
		pd, sts := pp.GetDir()
		if is.Error(sts) {
			break
		}
		p = &Project{
			Hostname: modeler.ItemId(pp.Hostname),
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
