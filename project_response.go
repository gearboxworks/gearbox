package gearbox

import (
	"gearbox/api"
)

var _ api.SelfLinkGetter = (*ProjectResponse)(nil)

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

func (me *ProjectResponse) GetApiSelfLink() string {
	return me.Project.GetApiSelfLink()
}

func NewProjectResponse(p *Project) *ProjectResponse {
	return &ProjectResponse{
		Hostname: p.Hostname,
		Basedir:  p.Basedir,
		Notes:    p.Notes,
		Path:     p.Path,

		Aliases:    p.GetAliases(),
		ServiceMap: p.GetServiceMap(),
		ProjectDir: p.GetProjectDir(),
		Project:    p,
	}
}

func NewProjectListResponse(p *Project) *api.ListItemResponse {
	return api.NewListItemResponse(p.GetApiSelfLink(), NewProjectResponse(p))
}
