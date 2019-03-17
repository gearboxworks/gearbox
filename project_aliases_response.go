package gearbox

import "gearbox/api"

type ProjectAliasesResponse struct {
	Aliases Aliases
	Project *Project
}

func NewProjectAliasesResponse(aliases Aliases, proj *Project) *ProjectAliasesResponse {
	return &ProjectAliasesResponse{
		Aliases: aliases,
		Project: proj,
	}
}

func (me *ProjectAliasesResponse) GetApiSelfLink() api.UriTemplate {
	url, _ := me.Project.GetApiUrl(ProjectAliasesResource)
	return url
}
func (me *ProjectAliasesResponse) GetResponseData() interface{} {
	return me.Aliases
}
