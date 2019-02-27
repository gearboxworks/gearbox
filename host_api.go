package gearbox

import (
	"fmt"
	"gearbox/api"
	"gearbox/dockerhub"
	"gearbox/only"
	"github.com/labstack/echo"
	"net/http"
)

const Port = "9999"

type HostApi struct {
	Config  *Config
	Api     *api.Api
	Gearbox *Gearbox
}

func apiResponseDefaults() *api.Response {
	return &api.Response{
		Meta: &api.ResponseMeta{
			Service: "GearBox API",
			Version: "0.1",
			DocsUrl: "https://docs.gearbox.works/api",
		},
		Links: make(api.Links, 0),
	}
}

func NewHostApi(gearbox *Gearbox) *HostApi {
	ha := &HostApi{
		Config:  gearbox.Config,
		Api:     api.NewApi(echo.New(), apiResponseDefaults()),
		Gearbox: gearbox,
	}
	ha.Api.Port = Port
	ha.addRoutes()
	return ha
}

func (me *HostApi) Url() string {
	return fmt.Sprintf("http://127.0.0.1:%s", me.Api.Port)
}

func (me *HostApi) Start() {
	me.Api.Start()
}

func (me *HostApi) Stop() {
	me.Api.Stop()
}

func (me *HostApi) getStackResponse(ctx echo.Context) interface{} {
	var response interface{}
	for range only.Once {
		sn := me.getStackName(ctx)
		var ok bool
		response, ok = me.Gearbox.Stacks[StackName(sn)]
		if !ok {
			response = &api.Error{
				StatusCode: http.StatusNotFound,
				Error:      fmt.Errorf("'%s' is not a valid stack", sn),
			}
			break
		}
	}
	return response
}

func (me *HostApi) getStackMembersResponse(ctx echo.Context) interface{} {
	var response interface{}
	for range only.Once {
		response = me.getStackResponse(ctx)
		if _, ok := response.(api.Error); ok {
			break
		}
		stack, ok := response.(*Stack)
		if !ok {
			response = &api.Error{
				StatusCode: http.StatusInternalServerError,
				Error: fmt.Errorf("unexpected: stack '%s' not found",
					me.getStackName(ctx),
				),
			}
			break
		}
		response = stack.GetMembers()
	}
	return response
}

func (me *HostApi) getStackMemberResponse(ctx echo.Context) interface{} {
	var response interface{}
	for range only.Once {
		response := me.getStackMembersResponse(ctx)
		if _, ok := response.(api.Error); ok {
			break
		}
		memberMap, ok := response.(StackMemberMap)
		if !ok {
			response = &api.Error{
				StatusCode: http.StatusInternalServerError,
				Error: fmt.Errorf("unexpected: member map for stack '%s' not found",
					me.getStackName(ctx),
				),
			}
			break
		}
		member, ok := memberMap[me.getStackMemberName(ctx)]
		if !ok {
			response = &api.Error{
				StatusCode: http.StatusInternalServerError,
				Error: fmt.Errorf("unexpected: member map '%s' for stack '%s' not found",
					me.getStackMemberName(ctx),
					me.getStackName(ctx),
				),
			}
			break
		}
		response = member
	}
	return response
}

func (me *HostApi) getStackName(ctx echo.Context) StackName {
	return StackName(ctx.Param("stack"))
}

func (me *HostApi) getStackMemberName(ctx echo.Context) StackMemberName {
	return StackMemberName(ctx.Param("member"))
}

func (me *HostApi) addRoutes() {

	_api := me.Api

	_api.GET("/", "root", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "{}")
	})
	_api.GET("/projects", "projects", func(ctx echo.Context) error {
		return _api.JsonMarshalHandler(ctx, me.Config.Projects)
	})
	_api.GET("/projects/:project", "project", func(ctx echo.Context) error {
		return _api.JsonMarshalHandler(ctx, me.Config.Projects)
	})
	_api.GET("/projects/enabled", "enabled-projects", func(ctx echo.Context) error {
		return _api.JsonMarshalHandler(ctx, me.Config.Projects.GetEnabled())
	})
	_api.GET("/projects/disabled", "disabled-projects", func(ctx echo.Context) error {
		return _api.JsonMarshalHandler(ctx, me.Config.Projects.GetDisabled())
	})
	_api.GET("/projects/candidates", "candidate-projects", func(ctx echo.Context) error {
		return _api.JsonMarshalHandler(ctx, me.Config.Candidates)
	})
	_api.GET("/stacks", "stacks", func(ctx echo.Context) error {
		return _api.JsonMarshalHandler(ctx, me.Gearbox.Stacks)
	})
	_api.GET("/stacks/:stack", "stack", func(ctx echo.Context) error {
		response := me.getStackResponse(ctx)
		if _, ok := response.(*api.Error); !ok {
			response = response.(*Stack).CloneSansMembers()
		}
		return _api.JsonMarshalHandler(ctx, response)
	})
	_api.GET("/stacks/:stack/members", "stack-members", func(ctx echo.Context) error {
		response := me.getStackMembersResponse(ctx)
		return _api.JsonMarshalHandler(ctx, response)
	})
	_api.GET("/stacks/:stack/members/:member", "stack-member", func(ctx echo.Context) error {
		response := me.getStackMemberResponse(ctx)
		return _api.JsonMarshalHandler(ctx, response)
	})

	_api.GET("/stacks/:stack/members/:member/options", "stack-member-options", func(ctx echo.Context) error {
		response := me.getStackMemberResponse(ctx)

		response = me.Gearbox.RequestAvailableContainers(&dockerhub.ContainerQuery{})
		return _api.JsonMarshalHandler(ctx, response)
	})

	_api.PUT("/projects/paths/:path_id", "update-projects-path", func(ctx echo.Context) error {
		return _api.JsonMarshalHandler(ctx, &api.Error{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      fmt.Errorf("the 'update-projects-path' method has not been implemented yet"),
		})
	})

	_api.POST("/projects/new", "add-project", func(ctx echo.Context) error {
		return _api.JsonMarshalHandler(ctx, &api.Error{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      fmt.Errorf("the 'add-project' method has not been implemented yet"),
		})
	})
	_api.POST("/projects/paths/new", "add-projects-path", func(ctx echo.Context) error {
		return _api.JsonMarshalHandler(ctx, &api.Error{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      fmt.Errorf("the 'add-projects-path' method has not been implemented yet"),
		})
	})

	_api.DELETE("/projects/paths/:path_id", "delete-projects-path", func(ctx echo.Context) error {
		return _api.JsonMarshalHandler(ctx, &api.Error{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      fmt.Errorf("the 'delete-projects-path' method has not been implemented yet"),
		})
	})

}
