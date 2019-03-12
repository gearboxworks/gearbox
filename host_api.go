package gearbox

import (
	"encoding/json"
	"fmt"
	"gearbox/api"
	"gearbox/dockerhub"
	"gearbox/only"
	"gearbox/util"
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

func (me *HostApi) GetApiSelfLink(resourceType api.ResourceName) (url string, status Status) {
	for range only.Once {
		if me.Api == nil {
			status = NewStatus(&StatusArgs{
				HttpStatus: http.StatusInternalServerError,
				Help:       ContactSupportHelp(),
				Message: fmt.Sprintf("accessing host api when internal api property is nil for resource type '%s'",
					resourceType,
				),
			})
		}
		var err error
		url, err = me.Api.GetApiSelfLink(resourceType)
		if err != nil {
			status = NewStatus(&StatusArgs{
				HttpStatus:   http.StatusInternalServerError,
				HelpfulError: err.(util.HelpfulError),
				Message: fmt.Sprintf("the Api property is nil when accessing host api for resource type '%s'",
					resourceType,
				),
			})
		}
	}
	return url, status
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

func (me *HostApi) getStackName(ctx echo.Context) StackName {
	return StackName(ctx.Param("stack"))
}

func (me *HostApi) getRoleName(ctx echo.Context) RoleName {
	return RoleName(ctx.Param("service"))
}

func (me *HostApi) addRoutes() {

	_api := me.Api
	_api.GET("/", api.LinksResource, func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, nil)
	})

	me.addBasedirRoutes()
	me.addProjectRoutes()
	me.addStackRoutes()

}

func (me *HostApi) jsonMarshalHandler(_api *api.Api, ctx echo.Context, requestType api.ResourceName, value interface{}) error {
	var apiError *api.Status
	for range only.Once {
		status, ok := value.(Status)
		if ok && status.IsError() {
			apiError = &api.Status{
				Error:      status.Error,
				StatusCode: status.HttpStatus,
				Help:       status.ApiHelp,
			}
			apiError = _api.JsonMarshalHandler(ctx, requestType, apiError)
			break
		}
		status.Finalize()
		if ok {
			ctx.Response().Status = status.HttpStatus
		}
		apiError = _api.JsonMarshalHandler(ctx, requestType, value)
	}
	return apiError.Error
}

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

func (me *HostApi) addBasedirRoutes() {
	_api := me.Api

	_api.GET("/basedirs", "basedirs", func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.Config.GetHostBasedirs())
	})

	_api.POST("/basedirs/new", "basedir-add", func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.addBasedir(me.Gearbox, ctx, rt))
	})

	_api.PUT("/basedirs/:nickname", "basedir-update", func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.updateBasedir(me.Gearbox, ctx, rt))
	})

	_api.DELETE("/basedirs/:nickname", "basedir-delete", func(rt api.ResourceName, ctx echo.Context) error {
		return me.jsonMarshalHandler(_api, ctx, rt, me.deleteNamedBasedir(me.Gearbox, ctx, rt))
	})

}

func readBasedirFromRequest(name api.ResourceName, ctx echo.Context, bd *Basedir) (status Status) {
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

func (me *HostApi) addBasedir(gb *Gearbox, ctx echo.Context, requestType api.ResourceName) (status Status) {
	bd := Basedir{}
	status = readBasedirFromRequest(requestType, ctx, &bd)
	if !status.IsError() {
		me.Gearbox.RequestType = requestType
		status = me.Gearbox.AddBasedir(bd.HostDir, bd.Nickname)
	}
	return status
}
func (me *HostApi) updateBasedir(gb *Gearbox, ctx echo.Context, requestType api.ResourceName) (status Status) {
	bd := Basedir{}
	status = readBasedirFromRequest(requestType, ctx, &bd)
	if !status.IsError() {
		me.Gearbox.RequestType = requestType
		status = me.Gearbox.UpdateBasedir(bd.Nickname, bd.HostDir)
	}
	return status
}

func (me *HostApi) deleteNamedBasedir(gb *Gearbox, ctx echo.Context, requestType api.ResourceName) (status Status) {
	me.Gearbox.RequestType = requestType
	return me.Gearbox.DeleteNamedBasedir(getBasedirNickname(ctx))
}

func getBasedirNickname(ctx echo.Context) string {
	return ctx.Param("nickname")
}

func GetApiDocsUrl(topic api.ResourceName) string {
	return fmt.Sprintf("%s/%s", ApiDocsBaseUrl, topic)
}
func GetApiHelp(topic api.ResourceName) string {
	return fmt.Sprintf("see %s", GetApiDocsUrl(topic))
}
