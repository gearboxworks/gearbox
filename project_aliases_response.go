package gearbox

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

func (me *ProjectAliasesResponse) GetApiSelfLink() string {
	url, _ := me.Project.GetApiUrl(ProjectAliasesResource)
	return url
}
func (me *ProjectAliasesResponse) GetResponseData() interface{} {
	return me.Aliases
}
