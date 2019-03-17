package gearbox

import (
	"fmt"
	"gearbox/api"
	"gearbox/only"
	"gearbox/stat"
	"github.com/labstack/echo"
)

type HostApi struct {
	Config  Config
	Api     *api.Api
	Gearbox Gearbox
}

func apiResponseDefaults() *api.Response {
	return &api.Response{
		Meta: api.ResponseMeta{
			Service: HostAPIServiceName,
			Version: HostAPIVersion,
			DocsUrl: HostAPIDocsUrl,
		},
		Links: make(api.Links, 0),
	}
}

func NewHostApi(gearbox Gearbox) *HostApi {
	ha := &HostApi{
		Config:  gearbox.GetConfig(),
		Api:     api.NewApi(echo.New(), apiResponseDefaults()),
		Gearbox: gearbox,
	}
	ha.Api.Port = HostApiPort
	ha.addRoutes()
	return ha
}

func (me *HostApi) GetBaseUrl() (url string) {
	return me.Api.GetBaseUrl()
}

func (me *HostApi) GetUrl(name api.ResourceName, vars api.UriTemplateVars) (url string, status stat.Status) {
	return me.Api.GetUrl(name, vars)
}

func (me *HostApi) GetUrlPathTemplate(name api.ResourceName) (url string, status stat.Status) {
	for range only.Once {
		if me.Api == nil {
			status = stat.NewStatus(&stat.Args{
				Message: fmt.Sprintf("accessing host api when internal api property is nil for resource type '%s'",
					name,
				),
				Help: stat.ContactSupportHelp(),
			})
			break
		}
		url, status = me.Api.GetUrlPathTemplate(name)
		if status.IsError() {
			break
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

type HandlerFunc func(rc *api.RequestContext) interface{}

func (me *HostApi) GET(path string, name api.ResourceName, handler HandlerFunc) *echo.Route {
	return me.Api.GET(path, name, func(rc *api.RequestContext) (err error) {
		me.Gearbox.SetResourceName(rc.ResourceName)
		if handler != nil {
			err = me.jsonMarshalHandler(rc, handler(rc))
		} else {
			err = me.jsonMarshalHandler(rc, nil)
		}
		return err
	})
}

func (me *HostApi) POST(path string, name api.ResourceName, handler HandlerFunc) *echo.Route {
	return me.Api.POST(path, name, func(rc *api.RequestContext) error {
		return me.jsonMarshalHandler(rc, handler(rc))
	})
}

func (me *HostApi) PUT(path string, name api.ResourceName, handler HandlerFunc) *echo.Route {
	return me.Api.PUT(path, name, func(rc *api.RequestContext) error {
		return me.jsonMarshalHandler(rc, handler(rc))
	})
}

func (me *HostApi) DELETE(path string, name api.ResourceName, handler HandlerFunc) *echo.Route {
	return me.Api.DELETE(path, name, func(rc *api.RequestContext) error {
		return me.jsonMarshalHandler(rc, handler(rc))
	})
}

func (me *HostApi) jsonMarshalHandler(rc *api.RequestContext, value interface{}) error {
	var status stat.Status
	for range only.Once {
		status, ok := value.(stat.Status)
		if ok {
			status.Finalize()
			rc.Context.Response().Status = status.HttpStatus
			break
		}
		status = rc.JsonMarshalHandler(value)
	}
	return status
}
