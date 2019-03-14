package api

import (
	"errors"
	"fmt"
	"gearbox/only"
	"gearbox/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"strings"
)

const Port = "9999"

type ResourceVarName string
type ResourceName string

const DocsBaseUrl = "https://docs.gearbox.works/api"

const SelfResource ResourceName = "self"
const LinksResource ResourceName = "links"

type Links map[ResourceName]string

type SelfLinkGetter interface {
	GetApiSelfLink() string
}

type Api struct {
	Echo     *echo.Echo
	Port     string
	Defaults *Response
}

func NewApi(echo *echo.Echo, defaults *Response) *Api {
	//
	//See https://echo.labstack.com/cookbook/cors
	//See https://flaviocopes.com/golang-enable-cors/
	//
	echo.Use(middleware.CORS())
	echo.HideBanner = true
	echo.HidePort = true
	return &Api{
		Echo:     echo,
		Defaults: defaults,
	}
}
func ContactSupportHelp() string {
	return "contact support"
}

func (me *Api) GetApiSelfLink(resourceType ResourceName) (url string, err error) {
	for range only.Once {
		if me.Defaults == nil {
			err = util.AddHelpToError(
				errors.New(fmt.Sprintf("the Defaults property is nil when accessing api for resource type '%s'",
					resourceType,
				)),
				ContactSupportHelp(),
			)
		}
		url, err = me.Defaults.GetApiSelfLink(resourceType)
	}
	return url, err
}

type SuccessInspector interface {
	IsSuccess() bool
}

type ResponseDataGetter interface {
	GetResponseData() interface{}
}

type HandlerFunc func(ctx *RequestContext) error

//func (me *Api) GET(path string, rc *RequestContext, handler HandlerFunc) *echo.Route {
//	me.Defaults.Meta.Resource = rc.ResourceName
//	me.Defaults.Links[rc.ResourceName] = convertEchoPathToUriTemplatePath(path)
//	return me.Echo.GET(path, rc.WrapHandler(handler))
//}
func (me *Api) GET(path string, name ResourceName, handler HandlerFunc) *echo.Route {
	rc := NewRequestContext(me, name)
	me.Defaults.Meta.Resource = name
	me.Defaults.Links[name] = convertEchoPathToUriTemplatePath(path)
	return me.Echo.GET(path, rc.WrapHandler(handler))
}
func (me *Api) POST(path string, name ResourceName, handler HandlerFunc) *echo.Route {
	rc := NewRequestContext(me, name)
	me.Defaults.Meta.Resource = name
	me.Defaults.Links[name] = convertEchoPathToUriTemplatePath(path)
	me.Echo.OPTIONS(path, optionsHandler)
	return me.Echo.POST(path, rc.WrapHandler(handler))
}
func (me *Api) DELETE(path string, name ResourceName, handler HandlerFunc) *echo.Route {
	rc := NewRequestContext(me, name)
	me.Defaults.Meta.Resource = name
	me.Defaults.Links[name] = convertEchoPathToUriTemplatePath(path)
	me.Echo.OPTIONS(path, optionsHandler)
	return me.Echo.DELETE(path, rc.WrapHandler(handler))
}
func (me *Api) PUT(path string, name ResourceName, handler HandlerFunc) *echo.Route {
	rc := NewRequestContext(me, name)
	me.Defaults.Meta.Resource = name
	me.Defaults.Links[name] = convertEchoPathToUriTemplatePath(path)
	me.Echo.OPTIONS(path, optionsHandler)
	return me.Echo.PUT(path, rc.WrapHandler(handler))
}
func optionsHandler(ctx echo.Context) error {
	return nil
}

func (me *Api) Start() {
	err := me.Echo.Start(":" + me.Port)
	if err != nil {
		util.Error(err)
	}
}

func (me *Api) Stop() {
	err := me.Echo.Close()
	if err != nil {
		util.Error(err)
	}
}

func convertEchoPathToUriTemplatePath(url string) string {
	parts := strings.Split(url, "/")
	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		if []byte(p)[0] != ':' {
			continue
		}
		parts[i] = fmt.Sprintf("{%s}", p[1:])
	}
	return strings.Join(parts, "/")
}

type UriTemplateVars map[ResourceVarName]string

func ExpandUriTemplate(template string, vars UriTemplateVars) string {
	url := template
	for vn, val := range vars {
		url = strings.Replace(url, fmt.Sprintf("{%s}", vn), val, -1)
	}
	return url
}

func GetApiDocsUrl(topic ResourceName) string {
	return fmt.Sprintf("%s/%s", DocsBaseUrl, topic)
}

func GetApiHelp(topic ResourceName) string {
	return fmt.Sprintf("see %s", GetApiDocsUrl(topic))
}

func (me *Api) NotYetImplemented(rc *RequestContext) interface{} {
	return &Status{
		StatusCode: http.StatusMethodNotAllowed,
		Error:      fmt.Errorf("the '%s' resource has not been implemented yet", rc.ResourceName),
	}
}
