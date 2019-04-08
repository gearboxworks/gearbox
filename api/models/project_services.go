package models

import (
	"gearbox/config"
	"gearbox/service"
	"gearbox/types"
)

const ProjectServicesResource types.RouteName = "project-services"
const ProjectServicesAdd types.RouteName = ProjectServicesResource + "-add"
const ProjectServicesUpdate types.RouteName = ProjectServicesResource + "-update"
const ProjectServicesDelete types.RouteName = ProjectServicesResource + "-delete"

type ProjectServices struct {
	ServiceMap service.StackMap
	Project    *config.Project
}

func NewProjectServices(svcmap service.StackMap, proj *config.Project) *ProjectServices {
	return &ProjectServices{
		ServiceMap: svcmap,
		Project:    proj,
	}
}

//func (me *Services) GetApiUrl() api.UriTemplate {
//	url, _ := me.Project.GetApiUrl(ProjectServicesResource)
//	return url
//}
func (me *ProjectServices) GetData() interface{} {
	return me.ServiceMap
}
