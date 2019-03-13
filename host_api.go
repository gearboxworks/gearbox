package gearbox

import (
	"fmt"
	"gearbox/api"
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

func GetApiDocsUrl(topic api.ResourceName) string {
	return fmt.Sprintf("%s/%s", ApiDocsBaseUrl, topic)
}

func GetApiHelp(topic api.ResourceName) string {
	return fmt.Sprintf("see %s", GetApiDocsUrl(topic))
}
