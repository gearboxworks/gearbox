package gearbox

import (
	"fmt"
	"gearbox/api"
	"gearbox/only"
	"gearbox/stat"
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

func getRoleSpec(rc *api.RequestContext) RoleSpec {
	return RoleSpec(rc.Param("service"))
}

func (me *HostApi) getStacksResponse(rc *api.RequestContext) (response interface{}) {
	response, status := me.Gearbox.GetStackMap()
	if status.IsError() {
		response = status
	}
	return response
}

func (me *HostApi) getServiceResponseOptions(rc *api.RequestContext) (response interface{}) {
	response, status := me.Gearbox.RequestAvailableContainers()
	if status.IsError() {
		response = status
	}
	return response
}

func (me *HostApi) getServiceResponse(rc *api.RequestContext) (response interface{}) {
	var status stat.Status
	for range only.Once {
		response := me.getServicesResponse(rc)
		if _, ok := response.(stat.Status); ok {
			break
		}
		roleMap, ok := response.(RoleMap)
		if !ok {
			status = stat.NewFailStatus(&stat.Args{
				Message: "getStackResponse() returned type other than Status or RoleMap",
			})
			break
		}
		role, ok := roleMap[getRoleSpec(rc)]
		if !ok {
			status = stat.NewFailStatus(&stat.Args{
				Message: fmt.Sprintf("role spec '%s' for stack '%s' not found",
					getRoleSpec(rc),
					me.getStackName(rc),
				),
				HttpStatus: http.StatusNotFound,
			})
			break
		}
		response = role
	}
	if status.IsError() {
		response = status
	}
	return response
}

func (me *HostApi) getServicesResponse(rc *api.RequestContext) (response interface{}) {
	var status stat.Status
	for range only.Once {
		response = me.getStackResponse(rc)
		if _, ok := response.(stat.Status); ok {
			break
		}
		stack, ok := response.(*Stack)
		if !ok {
			status = stat.NewFailStatus(&stat.Args{
				Message: "getStackResponse() returned type other than Status or Stack",
			})
			break
		}
		response = stack.GetRoleMap()
	}
	if status.IsError() {
		response = status
	}
	return response

}

func (me *HostApi) getStackResponse(rc *api.RequestContext) (response interface{}) {
	var status stat.Status
	for range only.Once {
		var sn StackName
		sn, status = getStackName(rc)
		if status.IsError() {
			break
		}
		var stack *Stack
		stack, status = FindNamedStack(me.Gearbox, sn)
		if status.IsError() {
			break
		}
		response = stack.LightweightClone()
	}
	if status.IsError() {
		response = status
	}
	return response
}
