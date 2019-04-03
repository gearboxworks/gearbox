package routes

import (
	"gearbox/api"
	"gearbox/config"
	"gearbox/service"
)

const ProjectServicesResource api.RouteName = "project-services"
const ProjectServicesAdd api.RouteName = ProjectServicesResource + "-add"
const ProjectServicesUpdate api.RouteName = ProjectServicesResource + "-update"
const ProjectServicesDelete api.RouteName = ProjectServicesResource + "-delete"

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
