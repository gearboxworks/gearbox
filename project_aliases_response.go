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
	return me.Project.GetApiSelfLink(ProjectAliasesResource)
}
func (me *ProjectAliasesResponse) GetResponseData() interface{} {
	return me.Aliases
}
