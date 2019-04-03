package config

import (
	"errors"
	"fmt"
	"gearbox/status"
	"gearbox/types"
	"net/http"
)

type ProjectMap map[types.Hostname]*Project

func (me ProjectMap) GetProject(cfg Configer, hostname types.Hostname) (p *Project, sts status.Status) {
	var ok bool
	p, ok = me[hostname]
	if ok {
		// The next two
		p.Config = cfg
		p.Hostname = hostname
		sts = status.Success("got project '%s'", hostname)
	} else {
		msg := fmt.Sprintf("project hostname '%s' does not exist", hostname)
		sts = status.Wrap(errors.New(msg), &status.Args{
			Message:    msg,
			HttpStatus: http.StatusBadRequest,
		})
	}
	return p, sts
}

func (me ProjectMap) ProjectExists(hostname types.Hostname) (ok bool) {
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

func (me ProjectMap) FindProject(hostname types.Hostname) (p *Project, sts status.Status) {
	var ok bool
	p, ok = me[hostname]
	if !ok {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("hostname '%s' not found", hostname),
		})
	}
	return p, sts
}

func (me ProjectMap) FindProjectByPath(basedir types.Nickname, path types.RelativePath) (p *Project, sts status.Status) {
	var hn types.Hostname
	var _p *Project
	for hn, _p = range me {
		if path == types.RelativePath(hn) {
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
	if p == nil {
		sts = status.Fail(&status.Args{
			Message: fmt.Sprintf("project path '%s' not found", path),
		})
	} else {

		p.Hostname = hn
		p.Path = path
	}
	return p, sts
}
