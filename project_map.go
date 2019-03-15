package gearbox

import (
	"errors"
	"fmt"
	"gearbox/api"
	"gearbox/only"
	"gearbox/stat"
	"net/http"
)

type ProjectMap map[string]*Project

func (me ProjectMap) FindProjectWithDetails(gb *Gearbox, hostname string) (p *Project, status stat.Status) {
	for range only.Once {
		p, status = me.GetProject(gb, hostname)
		if status.IsError() {
			break
		}
		status = p.LoadProjectDetails()
		if status.IsError() {
			break
		}
		status = stat.NewOkStatus("got project '%s'", hostname)
	}
	return p, status
}

func (me ProjectMap) GetProject(gb *Gearbox, hostname string) (p *Project, status stat.Status) {
	var ok bool
	p, ok = me[hostname]
	if ok {
		// The next two
		p.Gearbox = gb
		p.Hostname = hostname
		status = stat.NewOkStatus("got project '%s'", hostname)
	} else {
		msg := fmt.Sprintf("project hostname '%s' does not exist", hostname)
		status = stat.NewStatus(&stat.Args{
			Failed:     true,
			Message:    msg,
			HttpStatus: http.StatusBadRequest,
			ApiHelp:    api.GetApiDocsUrl(gb.RequestType),
			Error:      errors.New(msg),
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
