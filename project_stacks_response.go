package gearbox

import "gearbox/api"

type ProjectStacksResponse struct {
	StackNames `json:"stacks"`
	*Project   `json:"-"`
}

func NewProjectStacksResponse(p *Project) *ProjectStacksResponse {
	return &ProjectStacksResponse{
		StackNames: p.ServiceMap.GetStackNames(),
		Project:    p,
	}
}

func (me *ProjectStacksResponse) GetApiUrl() api.UriTemplate {
	url, _ := me.Project.GetApiUrl(ProjectStacksResource)
	return url
}

func (me *ProjectStacksResponse) GetResponseData() interface{} {
	return me.StackNames
}
