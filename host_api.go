package gearbox

import (
	"encoding/json"
	"fmt"
	"gearbox/api"
	"gearbox/dockerhub"
	"gearbox/only"
	"github.com/labstack/echo"
	"io/ioutil"
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
		Meta: api.ResponseMeta{
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
			response = &api.Status{
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
		response = stack.GetMembers()
	}
	return response
}

func (me *HostApi) getStackMemberResponse(ctx echo.Context) interface{} {
	var response interface{}
	for range only.Once {
		response := me.getStackMembersResponse(ctx)
		if _, ok := response.(api.Status); ok {
			break
		}
		memberMap, ok := response.(StackMemberMap)
		if !ok {
			response = &api.Status{
				StatusCode: http.StatusInternalServerError,
				Error: fmt.Errorf("unexpected: member map for stack '%s' not found",
					me.getStackName(ctx),
				),
			}
			break
		}
		member, ok := memberMap[me.getStackMemberName(ctx)]
		if !ok {
			response = &api.Status{
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

	me.addBaseDirRoutes()
	me.addProjectRoutes()
	me.addStackRoutes()

}

func (me *HostApi) jsonMarshalHandler(_api *api.Api, ctx echo.Context, value interface{}) error {
	var apiError *api.Status
	for range only.Once {
		s, ok := value.(*Status)
		if ok && s.IsError() {
			apiError = &api.Status{
				Error:      s.Error,
				StatusCode: s.HttpStatus,
				Help:       s.ApiHelp,
			}
			_ = ctx.String(apiError.StatusCode, apiError.ToJson())
			break
		}
		if ok {
			ctx.Response().Status = s.HttpStatus
		}
		apiError = _api.JsonMarshalHandler(ctx, value)
	}
	return apiError.Error
}

func (me *HostApi) addStackRoutes() {
	_api := me.Api
	_api.GET("/stacks", "stacks", func(ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, me.Gearbox.Stacks)
	})

	_api.GET("/stacks/:stack", "stack", func(ctx echo.Context) error {
		response := me.getStackResponse(ctx)
		if _, ok := response.(*api.Status); !ok {
			response = response.(*Stack).CloneSansMembers()
		}
		return me.jsonMarshalHandler(_api, ctx, response)
	})

	_api.GET("/stacks/:stack/members", "stack-members", func(ctx echo.Context) error {
		response := me.getStackMembersResponse(ctx)
		return me.jsonMarshalHandler(_api, ctx, response)
	})

	_api.GET("/stacks/:stack/members/:member", "stack-member", func(ctx echo.Context) error {
		response := me.getStackMemberResponse(ctx)
		return me.jsonMarshalHandler(_api, ctx, response)
	})

	_api.GET("/stacks/:stack/members/:member/options", "stack-member-options", func(ctx echo.Context) error {
		response := me.getStackMemberResponse(ctx)
		response = me.Gearbox.RequestAvailableContainers(&dockerhub.ContainerQuery{})
		return me.jsonMarshalHandler(_api, ctx, response)
	})
}

func (me *HostApi) addProjectRoutes() {
	_api := me.Api
	_api.GET("/projects", "projects", func(ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, me.Config.Projects)
	})
	_api.GET("/projects/:project", "project", func(ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, me.Config.Projects)
	})
	_api.GET("/projects/enabled", "projects-enabled", func(ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, me.Config.Projects.GetEnabled())
	})
	_api.GET("/projects/disabled", "projects-disabled", func(ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, me.Config.Projects.GetDisabled())
	})
	_api.GET("/projects/candidates", "project-candidates", func(ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, me.Config.Candidates)
	})

	_api.POST("/projects/new", "project-add", func(ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, &api.Status{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      fmt.Errorf("the 'project-add' method has not been implemented yet"),
		})
	})

	_api.POST("/projects/:project", "project-update", func(ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, &api.Status{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      fmt.Errorf("the 'project-update' method has not been implemented yet"),
		})
	})

	_api.DELETE("/basedirs/:project", "project-delete", func(ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, &api.Status{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      fmt.Errorf("the 'project-delete' method has not been implemented yet"),
		})
	})
}

func closeBody(ctx echo.Context) {
	_ = ctx.Request().Body.Close()
}

func (me *HostApi) addBaseDirRoutes() {
	_api := me.Api

	_api.GET("/basedirs", "basedirs", func(ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, me.Config.GetHostBaseDirs())
	})

	_api.POST("/basedirs/new", "basedir-add", func(ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, me.addBaseDir(ctx))
	})

	_api.GET("/basedirs/:nickname", "basedir-delete", func(ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, me.deleteBaseDir(ctx))
	})

	_api.PUT("/basedirs/:nickname", "basedir-update", func(ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, &api.Status{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      fmt.Errorf("the 'update-basedir' method has not been implemented yet"),
		})
	})

}

func getBaseDirNickname(ctx echo.Context) string {
	return ctx.Param("nickname")
}

func validateBaseDirNickename(nickname string) (status *Status) {
	for range only.Once {
		if nickname == "" {
			status = NewStatus(&StatusArgs{
				Success:    false,
				Message:    "basedir nickname is empty",
				HttpStatus: http.StatusBadRequest,
				ApiHelp:    "see https://docs.gearbox.works/api/delete-basedir",
			})
			break
		}
		status = NewOkStatus()
	}
	return status
}

func (me *HostApi) deleteBaseDir(ctx echo.Context) (status *Status) {
	for range only.Once {
		nickname := getBaseDirNickname(ctx)
		status = validateBaseDirNickename(nickname)
		if status.IsError() {
			break
		}
		status = NewStatus(&StatusArgs{
			Success:    false,
			Message:    "not yet implemented",
			HttpStatus: http.StatusNotImplemented,
		})
	}
	return status
}

func (me *HostApi) addBaseDir(ctx echo.Context) (status *Status) {
	for range only.Once {
		bd := BaseDir{}
		defer closeBody(ctx)
		b, err := ioutil.ReadAll(ctx.Request().Body)
		if err != nil {
			status = NewStatus(&StatusArgs{
				Message:    "could not read request body",
				HttpStatus: http.StatusUnprocessableEntity,
				ApiHelp:    "see https://docs.api.gearbox.works",
				Error:      err,
			})
			break
		}
		err = json.Unmarshal(b, &bd)
		if err != nil {
			status = NewStatus(&StatusArgs{
				Message:    fmt.Sprintf("unexpected format for request body: '%s'", string(b)),
				HttpStatus: http.StatusUnprocessableEntity,
				ApiHelp:    "see https://docs.api.gearbox.works/add-basedir",
				Error:      err,
			})
			break
		}
		status = me.Gearbox.AddBaseDir(bd.HostDir, bd.Nickname)
		if !status.IsError() {
			status = NewSuccessStatus(http.StatusCreated)
		}
	}
	return status
}
