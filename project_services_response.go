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

func (me *ProjectServicesResponse) GetApiUrl() string {
	url, _ := me.Project.GetApiUrl(ProjectServicesResource)
	return url
}
func (me *ProjectServicesResponse) GetResponseData() interface{} {
	return me.ServiceMap
}
