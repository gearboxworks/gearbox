package gearbox

import (
	"fmt"
	"gearbox/only"
	"gearbox/util"
	"net/http"
	"path/filepath"
	"strings"
)

type ProjectMap map[string]*Project

type Projects []*Project

type Project struct {
	Hostname string   `json:"hostname"`
	Scope    string   `json:"scope"`
	Enabled  bool     `json:"enabled"`
	Basedir  string   `json:"basedir"`
	Notes    string   `json:"notes"`
	Path     string   `json:"path"`
	Gearbox  *Gearbox `json:"-"`
	*ProjectDetails
}

func NewProject(gb *Gearbox, path string) *Project {
	return &Project{
		Gearbox:  gb,
		Path:     path,
		Hostname: GetHostnameFromPath(path),
	}
}

type ProjectResponse struct {
	Hostname   string     `json:"hostname"`
	Scope      string     `json:"scope"`
	Enabled    bool       `json:"enabled"`
	Basedir    string     `json:"basedir"`
	Notes      string     `json:"notes"`
	Path       string     `json:"path"`
	ProjectDir string     `json:"project_dir"`
	Aliases    Aliases    `json:"aliases"`
	ServiceMap ServiceMap `json:"stack"`
}

func NewProjectResponse(p *Project) *ProjectResponse {
	return &ProjectResponse{
		Hostname:   p.Hostname,
		Scope:      p.Scope,
		Basedir:    p.Basedir,
		Notes:      p.Notes,
		Path:       p.Path,
		Aliases:    p.Aliases,
		ServiceMap: p.ServiceMap,
		ProjectDir: filepath.Dir(p.ProjectDetails.Filepath),
	}
}

func (me *Project) ClearDetails() {
	me.ProjectDetails = nil
}

func (me *Project) MaybeLoadDetails() (status Status) {
	for range only.Once {
		if me.ProjectDetails != nil {
			break
		}
		status = me.LoadProjectFile()
	}
	return status
}

func (me *Project) GetProjectFilepath() (fp string, err error) {
	return me.Gearbox.GetProjectFilepath(me.Path, me.Basedir)
}

func (me *Project) LoadProjectFile() (status Status) {
	var err error
	for range only.Once {
		var fp string
		fp, err = me.GetProjectFilepath()
		if err != nil {
			status = NewStatus(&StatusArgs{
				Failed:       true,
				HttpStatus:   http.StatusInternalServerError,
				HelpfulError: err.(util.HelpfulError),
			})
			break
		}
		var j []byte
		j, err = util.ReadBytes(fp)
		if status.IsError() {
			break
		}
		pf := NewProjectFile(fp)
		if len(j) > 0 {
			status = pf.Unmarshal(j)
		}
		status = pf.FixupStack()
		if status.IsError() {
			break
		}
		me.ProjectDetails = pf.ExportProjectDetails()
	}
	return status
}

func (me *Project) Fullpath() (fp string) {
	fp, _ = ExpandHostBasedirPath(me.Gearbox, me.Basedir, me.Path)
	return fp
}

func (me ProjectMap) GetProjectResponse(gb *Gearbox, hostname string) (pr *ProjectResponse, status Status) {
	for range only.Once {
		var p *Project
		p, status = me.GetProject(gb, hostname)
		if status.IsError() {
			break
		}
		status = p.LoadProjectFile()
		if status.IsError() {
			break
		}
		pr = NewProjectResponse(p)
		status = NewOkStatus("got response for project '%s'", hostname)
	}
	return pr, status
}

func (me ProjectMap) GetProject(gb *Gearbox, hostname string) (p *Project, status Status) {
	var ok bool
	p, ok = me[hostname]
	if ok {
		// The next two
		p.Gearbox = gb
		p.Hostname = hostname
		status = NewOkStatus("got project '%s'", hostname)
	} else {
		status = NewStatus(&StatusArgs{
			Failed:     true,
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

func ValidateProjectHostname(hostname string, args ...*validateArgs) (status Status) {
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
				Failed:     true,
				Message:    "project hostname is empty",
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    apiHelp,
			})
			break
		}
		hnExists := _args.Gearbox.ProjectExists(hostname)
		if _args.MustExist && !hnExists {
			status = NewStatus(&StatusArgs{
				Failed:     true,
				Message:    fmt.Sprintf("no project exists with hostname '%s'", hostname),
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    apiHelp,
			})
			break
		}
		if _args.MustNotExist && hnExists {
			status = NewStatus(&StatusArgs{
				Failed:     true,
				Message:    fmt.Sprintf("project hostname '%s' already exists", hostname),
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    apiHelp,
			})
			break
		}
		status = NewOkStatus("validated project hostname '%s'", hostname)

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
		if basedir != _p.Basedir {
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
