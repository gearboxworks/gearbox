package gearbox

import (
	"encoding/json"
	"fmt"
	"gearbox/api"
	"gearbox/dockerhub"
	"gearbox/only"
	"github.com/labstack/echo"
	"net/http"
)

const Port = "9999"

const ApiDocsBaseUrl = "https://docs.gearbox.works/api"

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
	_api.GET("/", "root", func(rt string, ctx echo.Context) error {
		defaults := apiResponseDefaults()
		defaults.Meta.RequestType = rt
		defaults.Links[rt] = "/"
		return me.jsonMarshalHandler(_api, ctx, rt, defaults)
	})

	me.addBaseDirRoutes()
	me.addProjectRoutes()
	me.addStackRoutes()

}

func (me *HostApi) jsonMarshalHandler(_api *api.Api, ctx echo.Context, requestType string, value interface{}) error {
	var apiError *api.Status
	for range only.Once {
		status, ok := value.(*Status)
		if ok && status.IsError() {
			apiError = &api.Status{
				Error:      status.Error,
				StatusCode: status.HttpStatus,
				Help:       status.ApiHelp,
			}
			apiError = _api.JsonMarshalHandler(ctx, requestType, apiError)
			break
		}
		if ok {
			ctx.Response().Status = status.HttpStatus
		}
		apiError = _api.JsonMarshalHandler(ctx, requestType, value)
	}
	return apiError.Error
}

func (me *HostApi) addStackRoutes() {
	_api := me.Api
	_api.GET("/stacks", "stacks", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.Gearbox.Stacks)
	})

	_api.GET("/stacks/:stack", "stack", func(rt string, ctx echo.Context) error {
		response := me.getStackResponse(ctx)
		if _, ok := response.(*api.Status); !ok {
			response = response.(*Stack).CloneSansMembers()
		}
		return me.jsonMarshalHandler(_api, ctx, rt, response)
	})

	_api.GET("/stacks/:stack/members", "stack-members", func(rt string, ctx echo.Context) error {
		response := me.getStackMembersResponse(ctx)
		return me.jsonMarshalHandler(_api, ctx, rt, response)
	})

	_api.GET("/stacks/:stack/members/:member", "stack-member", func(rt string, ctx echo.Context) error {
		response := me.getStackMemberResponse(ctx)
		return me.jsonMarshalHandler(_api, ctx, rt, response)
	})

	_api.GET("/stacks/:stack/members/:member/options", "stack-member-options", func(rt string, ctx echo.Context) error {
		response := me.getStackMemberResponse(ctx)
		response = me.Gearbox.RequestAvailableContainers(&dockerhub.ContainerQuery{})
		return me.jsonMarshalHandler(_api, ctx, rt, response)
	})
}

func (me *HostApi) addProjectRoutes() {
	_api := me.Api
	_api.GET("/projects", "projects", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.Config.Projects)
	})
	_api.GET("/projects/:hostname", "project", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.getProject(ctx, rt))
	})
	_api.GET("/projects/enabled", "projects-enabled", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.Config.Projects.GetEnabled())
	})
	_api.GET("/projects/disabled", "projects-disabled", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.Config.Projects.GetDisabled())
	})
	_api.GET("/projects/candidates", "project-candidates", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.Config.Candidates)
	})

	_api.POST("/projects/new", "project-add", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, &api.Status{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      fmt.Errorf("the 'project-add' method has not been implemented yet"),
		})
	})

	_api.POST("/projects/:hostname", "project-update", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, &api.Status{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      fmt.Errorf("the 'project-update' method has not been implemented yet"),
		})
	})

	_api.DELETE("/projects/:hostname", "project-delete", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, &api.Status{
			StatusCode: http.StatusMethodNotAllowed,
			Error:      fmt.Errorf("the 'project-delete' method has not been implemented yet"),
		})
	})
}
func (me *HostApi) addBaseDirRoutes() {
	_api := me.Api

	_api.GET("/basedirs", "basedirs", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.Config.GetHostBaseDirs())
	})

	_api.POST("/basedirs/new", "basedir-add", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.addBaseDir(me.Gearbox, ctx, rt))
	})

	_api.PUT("/basedirs/:nickname", "basedir-update", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.updateBaseDir(me.Gearbox, ctx, rt))
	})

	_api.DELETE("/basedirs/:nickname", "basedir-delete", func(rt string, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.deleteNamedBaseDir(me.Gearbox, ctx, rt))
	})

}

func readBaseDirFromRequest(name string, ctx echo.Context, bd *BaseDir) (status *Status) {
	for range only.Once {
		apiHelp := GetApiHelp(name)
		defer api.CloseRequestBody(ctx)
		b, err := api.ReadRequestBody(ctx)
		if err != nil {
			status = NewStatus(&StatusArgs{
				Message:    "could not read request body",
				HttpStatus: http.StatusUnprocessableEntity,
				ApiHelp:    apiHelp,
				Error:      err,
			})
			break
		}
		err = json.Unmarshal(b, bd)
		if err != nil {
			status = NewStatus(&StatusArgs{
				Message:    fmt.Sprintf("unexpected format for request body: '%s'", string(b)),
				HttpStatus: http.StatusUnprocessableEntity,
				ApiHelp:    apiHelp,
				Error:      err,
			})
			break
		}
		status = NewOkStatus("read %d bytes from body of '%s' request",
			len(b),
			name,
		)
	}
	return status
}

func (me *HostApi) addBaseDir(gb *Gearbox, ctx echo.Context, requestType string) (status *Status) {
	bd := BaseDir{}
	status = readBaseDirFromRequest(requestType, ctx, &bd)
	if !status.IsError() {
		me.Gearbox.RequestType = requestType
		status = me.Gearbox.AddBaseDir(bd.HostDir, bd.Nickname)
	}
	return status
}
func (me *HostApi) updateBaseDir(gb *Gearbox, ctx echo.Context, requestType string) (status *Status) {
	bd := BaseDir{}
	status = readBaseDirFromRequest(requestType, ctx, &bd)
	if !status.IsError() {
		me.Gearbox.RequestType = requestType
		status = me.Gearbox.UpdateBaseDir(bd.Nickname, bd.HostDir)
	}
	return status
}

func (me *HostApi) deleteNamedBaseDir(gb *Gearbox, ctx echo.Context, requestType string) (status *Status) {
	me.Gearbox.RequestType = requestType
	return me.Gearbox.DeleteNamedBaseDir(getBaseDirNickname(ctx))
}

func (me *HostApi) getProject(ctx echo.Context, requestType string) (response interface{}) {
	for range only.Once {
		me.Gearbox.RequestType = requestType
		pr, status := me.Gearbox.GetProjectResponse(getProjectHostname(ctx))
		if status.IsError() {
			response = status
			break
		}
		response = pr
	}
	return response
}

func getProjectHostname(ctx echo.Context) string {
	return ctx.Param("hostname")
}

func getBaseDirNickname(ctx echo.Context) string {
	return ctx.Param("nickname")
}

func GetApiDocsUrl(topic string) string {
	return fmt.Sprintf("%s/%s", ApiDocsBaseUrl, topic)
}
func GetApiHelp(topic string) string {
	return fmt.Sprintf("see %s", GetApiDocsUrl(topic))
}
