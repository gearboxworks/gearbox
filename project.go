package gearbox

import (
	"fmt"
	"github.com/projectcfg/projectcfg"
	"strings"
)

const ProjectFile = "project.json"

type ProjectMap map[string]*Project

type Projects []*Project

func (me ProjectMap) GetEnabled() Projects {
	enabled := make(Projects, 0)
	for _, p := range me {
		if !p.Enabled {
			continue
		}
		enabled = append(enabled, p)
	}
	return enabled
}
func (me ProjectMap) GetDisabled() Projects {
	disabled := make(Projects, 0)
	for _, p := range me {
		if p.Enabled {
			continue
		}
		disabled = append(disabled, p)
	}
	return disabled
}

func GetHostnameFromPath(path string) string {
	hostname := path
	if !strings.Contains(hostname, ".") {
		hostname = fmt.Sprintf("%s.local", hostname)
	}
	return strings.ToLower(hostname)
}

func (me ProjectMap) FindProject(basedir, path string) (p *Project) {
	var hn string
	var _p *Project
	for hn, _p = range me {
		if path == hn {
			p = _p
			break
		}
		if path != _p.Path {
			continue
		}
		if basedir != _p.BaseDir {
			continue
		}
		p = _p
		break
	}
	if p != nil {
		p.Hostname = hn
		p.Path = path
	}
	return p
}

type Project struct {
	Hostname string                 `json:"-"`
	Enabled  bool                   `json:"enabled"`
	BaseDir  string                 `json:"base_dir"`
	Notes    string                 `json:"notes"`
	Path     string                 `json:"path"`
	Config   *projectcfg.ProjectCfg `json:"-"`
}

func NewProject(path string) *Project {
	return &Project{
		Path:     path,
		Hostname: GetHostnameFromPath(path),
	}
}
