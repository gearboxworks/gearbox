package gearbox

import (
	"fmt"
	"github.com/projectcfg/projectcfg"
	"strings"
)

const ProjectFile = "project.json"

type ProjectMap map[string]*Project

type Projects []*Project

type Project struct {
	Name     string                 `json:"name"`
	Hostname string                 `json:"hostname"`
	Enabled  bool                   `json:"enabled"`
	Group    int                    `json:"group"`
	Config   *projectcfg.ProjectCfg `json:"-"`
	Root     *string                `json:"-"`
}

func (me *Project) MakeHostname() string {
	hostname := me.Name
	if !strings.Contains(hostname, ".") {
		hostname = fmt.Sprintf("%s.local", hostname)
	}
	return hostname
}

func NewProject(name string, root *string) *Project {
	pr := Project{
		Root: root,
		Name: name,
	}
	pr.Hostname = pr.MakeHostname()
	return &pr
}

func (me Projects) GetEnabled() Projects {
	enabled := make(Projects, 0)
	for _, p := range me {
		if !p.Enabled {
			continue
		}
		enabled = append(enabled, p)
	}
	return enabled
}
func (me Projects) GetDisabled() Projects {
	disabled := make(Projects, 0)
	for _, p := range me {
		if p.Enabled {
			continue
		}
		disabled = append(disabled, p)
	}
	return disabled
}
