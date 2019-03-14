package gearbox

import (
	"fmt"
	"gearbox/api"
	"gearbox/dockerhub"
	"gearbox/only"
	"net/http"
)

func (me *HostApi) addStackRoutes() {

	me.GET("/stacks", "stacks", me.getStacksResponse)
	me.GET("/stacks/:stack", "stack-details", me.getStackResponse)
	me.GET("/stacks/:stack/services", "stack-services", me.getServicesResponse)
	me.GET("/stacks/:stack/services/:service", "stack-service", me.getServiceResponse)
	me.GET("/stacks/:stack/services/:service/options", "stack-service-options", me.getServiceResponse)
}

func (me *HostApi) getStackName(rc *api.RequestContext) StackName {
	return StackName(rc.Param("stack"))
}

func (me *HostApi) getRoleSpec(rc *api.RequestContext) RoleSpec {
	return RoleSpec(rc.Param("service"))
}

func (me *HostApi) getStacksResponse(rc *api.RequestContext) interface{} {
	return me.Gearbox.Stacks
}

func (me *HostApi) getServicesResponse(rc *api.RequestContext) interface{} {
	var response interface{}
	for range only.Once {
		response = me.getStackResponse(rc)
		if _, ok := response.(api.Status); ok {
			break
		}
		stack, ok := response.(*Stack)
		if !ok {
			response = &api.Status{
				StatusCode: http.StatusInternalServerError,
				Error: fmt.Errorf("unexpected: stack '%s' not found",
					me.getStackName(rc),
				),
			}
			break
		}
		response = stack.GetRoleMap()
	}
	return response

}

func (me *HostApi) getServiceResponseOptions(rc *api.RequestContext) interface{} {
	return me.Gearbox.RequestAvailableContainers(&dockerhub.ContainerQuery{})
}

func (me *HostApi) getServiceResponse(rc *api.RequestContext) interface{} {
	var response interface{}
	for range only.Once {
		response := me.getServicesResponse(rc)
		if _, ok := response.(api.Status); ok {
			break
		}
		serviceMap, ok := response.(ServiceMap)
		if !ok {
			response = &api.Status{
				StatusCode: http.StatusInternalServerError,
				Error: fmt.Errorf("unexpected: service map for stack '%s' not found",
					me.getStackName(rc),
				),
			}
			break
		}
		service, ok := serviceMap[me.getRoleSpec(rc)]
		if !ok {
			response = &api.Status{
				StatusCode: http.StatusInternalServerError,
				Error: fmt.Errorf("unexpected: service map '%s' for stack '%s' not found",
					me.getRoleSpec(rc),
					me.getStackName(rc),
				),
			}
			break
		}
		response = service
	}
	return response
}

func (me *HostApi) getStackResponse(rc *api.RequestContext) interface{} {
	var response interface{}
	for range only.Once {
		sn := me.getStackName(rc)
		var ok bool
		response, ok = me.Gearbox.Stacks[RoleSpec(sn)]
		if !ok {
			response = &api.Status{
				StatusCode: http.StatusNotFound,
				Error:      fmt.Errorf("'%s' is not a valid stack", sn),
			}
			break
		}
		if _, ok := response.(*api.Status); !ok {
			response = response.(*Stack).CloneSansServices()
		}
	}
	return response
}
