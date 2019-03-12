package gearbox

type ProjectServicesResponse struct {
	ServiceMap ServiceMap
	Project    *Project
}

func NewProjectServicesResponse(svcmap ServiceMap, proj *Project) *ProjectServicesResponse {
	return &ProjectServicesResponse{
		ServiceMap: svcmap,
		Project:    proj,
	}
}

func (me *ProjectServicesResponse) GetApiSelfLink() string {
	return me.Project.GetApiSelfLink(ProjectServicesResource)
}
func (me *ProjectServicesResponse) GetResponseData() interface{} {
	return me.ServiceMap
}
