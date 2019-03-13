package gearbox

import (
	"fmt"
	"gearbox/api"
	"gearbox/dockerhub"
	"gearbox/only"
	"github.com/labstack/echo"
	"net/http"
)

func (me *HostApi) addStackRoutes() {
	_api := me.Api
	_api.GET("/stacks", "stacks", func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.Gearbox.Stacks)
	})

	_api.GET("/stacks/:stack", "stack-details", func(rt api.ResourceName, ctx echo.Context) error {
		response := me.getStackResponse(ctx)
		if _, ok := response.(*api.Status); !ok {
			response = response.(*Stack).CloneSansServices()
		}
		return me.jsonMarshalHandler(_api, ctx, rt, response)
	})

	_api.GET("/stacks/:stack/services", "stack-services", func(rt api.ResourceName, ctx echo.Context) error {
		response := me.getServicesResponse(ctx)
		return me.jsonMarshalHandler(_api, ctx, rt, response)
	})

	_api.GET("/stacks/:stack/services/:service", "stack-service", func(rt api.ResourceName, ctx echo.Context) error {
		response := me.getServiceResponse(ctx)
		return me.jsonMarshalHandler(_api, ctx, rt, response)
	})

	_api.GET("/stacks/:stack/services/:service/options", "stack-service-options", func(rt api.ResourceName, ctx echo.Context) error {
		response := me.getServiceResponse(ctx)
		response = me.Gearbox.RequestAvailableContainers(&dockerhub.ContainerQuery{})
		return me.jsonMarshalHandler(_api, ctx, rt, response)
	})
}

func (me *HostApi) getStackName(ctx echo.Context) StackName {
	return StackName(ctx.Param("stack"))
}
func (me *HostApi) getRoleName(ctx echo.Context) RoleName {
	return RoleName(ctx.Param("service"))
}

func (me *HostApi) getServicesResponse(ctx echo.Context) interface{} {
	var response interface{}
	for range only.Once {
		response = me.getStackResponse(ctx)
		if _, ok := response.(api.Status); ok {
			break
		}
		stack, ok := response.(*Stack)
		if !ok {
			response = &api.Status{
				StatusCode: http.StatusInternalServerError,
				Error: fmt.Errorf("unexpected: stack '%s' not found",
					me.getStackName(ctx),
				),
			}
			break
		}
		response = stack.GetRoleMap()
	}
	return response
}

func (me *HostApi) getServiceResponse(ctx echo.Context) interface{} {
	var response interface{}
	for range only.Once {
		response := me.getServicesResponse(ctx)
		if _, ok := response.(api.Status); ok {
			break
		}
		serviceMap, ok := response.(ServiceMap)
		if !ok {
			response = &api.Status{
				StatusCode: http.StatusInternalServerError,
				Error: fmt.Errorf("unexpected: service map for stack '%s' not found",
					me.getStackName(ctx),
				),
			}
			break
		}
		service, ok := serviceMap[me.getRoleName(ctx)]
		if !ok {
			response = &api.Status{
				StatusCode: http.StatusInternalServerError,
				Error: fmt.Errorf("unexpected: service map '%s' for stack '%s' not found",
					me.getRoleName(ctx),
					me.getStackName(ctx),
				),
			}
			break
		}
		response = service
	}
	return response
}

func (me *HostApi) getStackResponse(ctx echo.Context) interface{} {
	var response interface{}
	for range only.Once {
		sn := me.getStackName(ctx)
		var ok bool
		response, ok = me.Gearbox.Stacks[RoleName(sn)]
		if !ok {
			response = &api.Status{
				StatusCode: http.StatusNotFound,
				Error:      fmt.Errorf("'%s' is not a valid stack", sn),
			}
			break
		}
	}
	return response
}
