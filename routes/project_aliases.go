package routes

import (
	"gearbox/api"
	"gearbox/config"
	"gearbox/project"
)

const ProjectAliasesRoute api.RouteName = "project-aliases"
const ProjectAliasAdd api.RouteName = ProjectAliasesRoute + "-add"
const ProjectAliasUpdate api.RouteName = ProjectAliasesRoute + "-update"
const ProjectAliasDelete api.RouteName = ProjectAliasesRoute + "-delete"

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
