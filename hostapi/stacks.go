package hostapi

import (
	"fmt"
	"gearbox"
	"gearbox/api"
	"gearbox/gears"
	"gearbox/gearspecid"
	"gearbox/hostapi/routes"
	"gearbox/only"
	"gearbox/status"
	"gearbox/types"
	"net/http"
	"strings"
)

const StacksResource api.RouteName = "stacks"

const StackDetailsResource api.RouteName = "stack-details"
const StackDetailsAdd api.RouteName = "stack-details-add"
const StackDetailsUpdate api.RouteName = "stack-details-update"
const StackDetailsDelete api.RouteName = "stack-details-delete"

const StackServicesResource api.RouteName = "stack-services"
const StackServiceResource api.RouteName = "stack-service"
const StackServiceOptionsResource api.RouteName = "stack-service-options"

func (me *HostApi) addStackRoutes() {

	//me.GET___("/stacks", StacksResource, me.getStacksResponse, nil, me.getStackNameValues)
	//me.GET___("/stacks/:authority/:stack", StackDetailsResource, me.getStackResponse, me.getStackNameValues)
	//
	////me.GET___("/stacks/:authority/:stack/services", StackServicesResource, me.getServicesResponse,nil)
	////me.GET___("/stacks/:authority/:stack/services/:service", StackServiceResource, me.getServiceResponse,nil)
	////me.GET___("/stacks/:stack/services/:service/options", StackServiceOptionsResource, me.getServiceResponse,nil)
	//
	//me.Relate(StacksResource, &api.Related{
	//	List:   StacksResource,
	//	Item:   StackDetailsResource,
	//	Update: StackDetailsUpdate,
	//	New:    StackDetailsAdd,
	//	Delete: StackDetailsDelete,
	//})

}

func (me *HostApi) getStackName(rc *api.RequestContext) types.Stackname {
	return types.Stackname(rc.Param("stack"))
}

func getRoleSpec(rc *api.RequestContext) gsid.Identifier {
	return gsid.Identifier(rc.Param("service"))
}

func (me *HostApi) getStacksResponse(rc *api.RequestContext) (response interface{}) {
	response, sts := me.Gearbox.GetNamedStackMap()
	if status.IsError(sts) {
		response = sts
	}
	return response
}

func (me *HostApi) getServiceResponseOptions(rc *api.RequestContext) (response interface{}) {
	response, sts := me.Gearbox.RequestAvailableContainers()
	if status.IsError(sts) {
		response = sts
	}
	return response
}

func (me *HostApi) getServiceResponse(rc *api.RequestContext) (response interface{}) {
	var sts status.Status
	for range only.Once {
		response := me.getServicesResponse(rc)
		if _, ok := response.(status.Status); ok {
			break
		}
		roleMap, ok := response.(gearbox.RoleMap)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: "getStackResponse() returned type other than Status or StackRoleMap",
			})
			break
		}
		role, ok := roleMap[getRoleSpec(rc)]
		if !ok {
			sts = status.Fail(&status.Args{
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
	if status.IsError(sts) {
		response = sts
	}
	return response
}

func (me *HostApi) getServicesResponse(rc *api.RequestContext) (response interface{}) {
	var sts status.Status
	for range only.Once {
		response = me.getStackResponse(rc)
		if _, ok := response.(status.Status); ok {
			break
		}
		stack, ok := response.(*gears.NamedStack)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: "getStackResponse() returned type other than Status or NamedStack",
			})
			break
		}
		response = stack.RoleMap
	}
	if status.IsError(sts) {
		response = sts
	}
	return response

}

func (me *HostApi) getStackResponse(rc *api.RequestContext) (response interface{}) {
	var sts status.Status
	for range only.Once {
		var nsid types.StackId
		nsid, sts = routes.GetNamedStackIdFromRequest(rc)
		if status.IsError(sts) {
			break
		}
		var stack *gears.NamedStack
		stack, sts = me.Gearbox.FindNamedStack(nsid)
		if status.IsError(sts) {
			break
		}
		response = routes.NewNamedStack(stack)
		//response = stack.LightweightClone()
	}
	if status.IsError(sts) {
		response = sts
	}
	return response
}

func (me *HostApi) getStackNameValues(...interface{}) (values api.ValuesFuncValues, sts status.Status) {
	for range only.Once {
		var nsids types.StackIds
		nsids, sts = me.Gearbox.GetGears().GetNamedStackIds()
		if status.IsError(sts) {
			break
		}
		values = make(api.ValuesFuncValues, 2)
		values[0] = make(api.ValueFuncVarsValues, len(nsids))
		values[1] = make(api.ValueFuncVarsValues, len(nsids))
		for i, fns := range nsids {
			parts := strings.Split(string(fns), "/")
			if len(parts) != 2 {
				sts = status.Fail(&status.Args{
					Message: fmt.Sprintf("stack name should have two (2) segments: '%s'", fns),
				})
				break
			}
			values[0][i] = parts[0]
			values[1][i] = parts[1]
		}
	}
	return values, sts
}
