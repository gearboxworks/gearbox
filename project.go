package gearbox

import (
	"fmt"
	"gearbox/api"
	"gearbox/only"
	"gearbox/util"
	"net/http"
	"path/filepath"
	"strings"
)

type Projects []*Project

type Project struct {
	Hostname string   `json:"hostname"`
	Enabled  bool     `json:"enabled"`
	Basedir  string   `json:"basedir"`
	Notes    string   `json:"notes"`
	Path     string   `json:"path"`
	Gearbox  *Gearbox `json:"-"`
	*ProjectDetails
}

func NewProject(gb *Gearbox, path string) *Project {
	p := Project{}
	p.Renew(gb, path)
	return &p
}

func (me *Project) Renew(gb *Gearbox, path string) {
	me.Gearbox = gb
	me.Path = path
	if me.Hostname == "" {
		me.Hostname = me.GetHostname()
	}
	return
}

func (me *Project) GetAliases() (aliases Aliases) {
	if me.ProjectDetails != nil {
		aliases = me.Aliases
	}
	return aliases
}

func (me *Project) GetServiceMap() (svcmap ServiceMap) {
	if me.ProjectDetails != nil {
		svcmap = me.ServiceMap
	}
	return svcmap
}

func (me *Project) GetProjectDir() (dir string) {
	if me.ProjectDetails != nil {
		dir = filepath.Dir(me.ProjectDetails.Filepath)
	} else {
		fp, err := me.GetProjectFilepath()
		if err == nil {
			dir = filepath.Dir(fp)
		} else {
			msg := []byte(fmt.Sprintf("failed to get filepath for project '%s'", me.Hostname))
			_, _ = me.Gearbox.errorLog.Write(msg)
		}
	}
	return dir
}

func (me *Project) GetApiSelfLink(name ...api.ResourceName) string {
	var rn api.ResourceName
	if len(name) == 0 {
		rn = ProjectDetailsResource
	} else {
		rn = name[0]
	}
	return me.Gearbox.GetApiSelfLink(rn,
		api.UriTemplateVars{
			HostnameResourceVar: me.Hostname,
		},
	)
}

func (me *Project) ClearDetails() {
	me.ProjectDetails = nil
}

func (me *Project) MaybeLoadDetails() (status Status) {
	for range only.Once {
		if me.ProjectDetails != nil {
			break
		}
		status = me.LoadProjectDetails()
	}
	return status
}

func (me *Project) GetProjectFilepath() (fp string, err error) {
	return me.Gearbox.GetProjectFilepath(me.Path, me.Basedir)
}

func (me *Project) LoadProjectDetails() (status Status) {
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

func (me *Project) GetHostname() string {
	hostname := me.Path
	if !strings.Contains(hostname, ".") {
		hostname = fmt.Sprintf("%s.local", hostname)
	}
	return strings.ToLower(hostname)
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
