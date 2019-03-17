package gearbox

import (
	"gearbox/api"
	"gearbox/stat"
)

var _ api.UrlGetter = (*ProjectResponse)(nil)

type ProjectResponse struct {
	Hostname   string     `json:"hostname"`
	Enabled    bool       `json:"enabled"`
	Basedir    string     `json:"basedir"`
	Notes      string     `json:"notes"`
	Path       string     `json:"path"`
	ProjectDir string     `json:"project_dir"`
	Aliases    Aliases    `json:"aliases,omitempty"`
	ServiceMap ServiceMap `json:"stack,omitempty"`
	Project    *Project   `json:"-"`
}

func (me *ProjectResponse) GetApiUrl(name ...api.ResourceName) (url string, status stat.Status) {
	return me.Project.GetApiUrl(name...)
}

func NewProjectResponse(p *Project) *ProjectResponse {
	pr := ProjectResponse{
		Hostname: p.Hostname,
		Basedir:  p.Basedir,
		Notes:    p.Notes,
		Path:     p.Path,

		Aliases:    p.GetAliases(),
		ServiceMap: p.GetServiceMap(),
		ProjectDir: p.GetProjectDir(),
		Project:    p,
	}
	return &pr
}

func NewProjectListResponse(p *Project) *api.ListItemResponse {
	url, _ := p.GetApiUrl()
	return api.NewListItemResponse(url, NewProjectResponse(p))
}
