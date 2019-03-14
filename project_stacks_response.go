package gearbox

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

func (me *ProjectStacksResponse) GetApiSelfLink() string {
	return me.Project.GetApiSelfLink(ProjectStacksResource)
}

func (me *ProjectStacksResponse) GetResponseData() interface{} {
	return me.StackNames
}
