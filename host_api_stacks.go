package gearbox

import (
	"fmt"
	"gearbox/api"
	"gearbox/only"
	"gearbox/stat"
	"net/http"
)

const StacksResource api.ResourceName = "stacks"
const StackDetailsResource api.ResourceName = "stack-details"
const AuthorityStackDetailsResource api.ResourceName = "authority-stack-details"
const StackServicesResource api.ResourceName = "stack-services"
const AuthorityStackServicesResource api.ResourceName = "authority-stack-services"
const StackServiceResource api.ResourceName = "stack-service"
const AuthorityStackServiceResource api.ResourceName = "authority-stack-service"
const StackServiceOptionsResource api.ResourceName = "stack-service-options"
const AuthorityStackServiceOptionsResource api.ResourceName = "authority-stack-service-options"

const AuthorityResourceVar api.ResourceVarName = "authority"
const StackNameResourceVar api.ResourceVarName = "stack"

func (me *HostApi) addStackRoutes() {

	me.GET("/stacks", StacksResource, me.getStacksResponse, nil)
	me.GET("/stacks/:stack", StackDetailsResource, me.getStackResponse, me.getStackNameValues)
	me.GET("/stacks/:authority/:stack", AuthorityStackDetailsResource, me.getStackResponse, me.getStackNameValues)

	//me.GET("/stacks/:stack/services", StackServicesResource, me.getServicesResponse,nil)
	//me.GET("/stacks/:authority/:stack/services", AuthorityStackServicesResource, me.getServicesResponse,nil)
	//
	//me.GET("/stacks/:stack/services/:service", StackServiceResource, me.getServiceResponse,nil)
	//me.GET("/stacks/:authority/:stack/services/:service", AuthorityStackServiceResource, me.getServiceResponse,nil)
	//
	//me.GET("/stacks/:stack/services/:service/options", StackServiceOptionsResource, me.getServiceResponse,nil)
	//me.GET("/stacks/:authority/:stack/services/:service/options", AuthorityStackServiceOptionsResource, me.getServiceResponse,nil)
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

func (me *HostApi) getStackNameValues(...interface{}) (values []string, status stat.Status) {
	for range only.Once {
		var fsns StackNames
		fsns, status = GetFullStackNames(me.Gearbox)
		if status.IsError() {
			break
		}
		values = make([]string, len(fsns))
		for i, sn := range fsns {
			values[i] = string(sn)
		}
	}
	return values, status
}
