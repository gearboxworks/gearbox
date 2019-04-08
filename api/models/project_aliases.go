package models

import (
	"gearbox/config"
	"gearbox/project"
	"gearbox/types"
)

const ProjectAliasesRoute types.RouteName = "project-aliases"
const ProjectAliasAdd types.RouteName = ProjectAliasesRoute + "-add"
const ProjectAliasUpdate types.RouteName = ProjectAliasesRoute + "-update"
const ProjectAliasDelete types.RouteName = ProjectAliasesRoute + "-delete"

type ProjectAliases struct {
	Aliases project.HostnameAliases
	Project *config.Project
}

func NewProjectAliases(aliases project.HostnameAliases, proj *config.Project) *ProjectAliases {
	return &ProjectAliases{
		Aliases: aliases,
		Project: proj,
	}
}

//func (me *ProjectAliases) GetApiSelfLink() api.UriTemplate {
//	url, _ := me.Project.GetApiUrl(ProjectAliasesRoute)
//	return url
//}
func (me *ProjectAliases) GetData() interface{} {
	return me.Aliases
}
