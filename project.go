package gearbox

import (
	"fmt"
	"gearbox/only"
	"github.com/projectcfg/projectcfg"
	"net/http"
	"strings"
)

const ProjectFile = "project.json"

type ProjectMap map[string]*Project

type Projects []*Project

type Project struct {
	Hostname string                 `json:"-"`
	Enabled  bool                   `json:"enabled"`
	BaseDir  string                 `json:"base_dir"`
	Notes    string                 `json:"notes"`
	Path     string                 `json:"path"`
	Config   *projectcfg.ProjectCfg `json:"-"`
	Gearbox  *Gearbox               `json:"-"`
}

type ProjectResponse struct {
	Hostname string `json:"hostname"`
	Enabled  bool   `json:"enabled"`
	BaseDir  string `json:"base_dir"`
	Notes    string `json:"notes"`
	FullPath string `json:"full_path"`
	projectcfg.ProjectCfg
}

func NewProject(gb *Gearbox, path string) *Project {
	return &Project{
		Gearbox:  gb,
		Path:     path,
		Hostname: GetHostnameFromPath(path),
	}
}

func (me *Project) Fullpath() (fp string) {
	fp, _ = ExpandBaseDirPath(me.Gearbox, me.BaseDir, me.Path)
	return fp
}

func (me ProjectMap) GetProjectResponse(gb *Gearbox, hostname string) (pr *ProjectResponse, status *Status) {
	for range only.Once {
		var p *Project
		p, status = me.GetProject(gb, hostname)
		if status.IsError() {
			break
		}
		var fp string
		fp, status = ExpandBaseDirPath(gb, p.BaseDir, p.Path)
		if status.IsError() {
			break
		}
		pr = &ProjectResponse{
			Hostname: p.Hostname,
			BaseDir:  p.BaseDir,
			Notes:    p.Notes,
			FullPath: fp,
		}
		status = NewOkStatus()
	}
	return pr, status
}

func (me ProjectMap) GetProject(gb *Gearbox, hostname string) (p *Project, status *Status) {
	var ok bool
	p, ok = me[hostname]
	if ok {
		status = NewOkStatus()
	} else {
		status = NewStatus(&StatusArgs{
			Success:    false,
			Message:    fmt.Sprintf("project hostname '%s' does not exist", hostname),
			HttpStatus: http.StatusBadRequest,
			ApiHelp:    GetApiDocsUrl(gb.RequestType),
		})
	}
	return p, status
}

func (me ProjectMap) ProjectExists(hostname string) (ok bool) {
	_, ok = me[hostname]
	return ok
}

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

func ValidateProjectHostname(hostname string, args ...*validateArgs) (status *Status) {
	for range only.Once {
		var apiHelp string
		var _args *validateArgs
		if len(args) == 0 {
			_args = &validateArgs{}
		} else {
			_args = args[0]
		}
		if _args.ApiHelpUrl != "" {
			apiHelp = fmt.Sprintf("see %s", _args.ApiHelpUrl)
		}

		if _args.MustNotBeEmpty && hostname == "" {
			status = NewStatus(&StatusArgs{
				Success:    false,
				Message:    "project hostname is empty",
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    apiHelp,
			})
			break
		}
		hnExists := _args.Gearbox.ProjectExists(hostname)
		if _args.MustExist && !hnExists {
			status = NewStatus(&StatusArgs{
				Success:    false,
				Message:    fmt.Sprintf("no project exists with hostname '%s'", hostname),
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    apiHelp,
			})
			break
		}
		if _args.MustNotExist && hnExists {
			status = NewStatus(&StatusArgs{
				Success:    false,
				Message:    fmt.Sprintf("project hostname '%s' already exists", hostname),
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    apiHelp,
			})
			break
		}
		status = NewOkStatus()

	}
	return status
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
