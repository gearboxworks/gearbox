package api

import (
	"fmt"
	"gearbox/only"
	"gearbox/stat"
	"gearbox/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"strings"
)

type ResourceVarName string
type ResourceName string

const DocsBaseUrl = "https://docs.gearbox.works/api"

const SelfResource ResourceName = "self"
const LinksResource ResourceName = "links"

type Links map[ResourceName]string

type UrlGetter interface {
	GetApiUrl(...ResourceName) (string, stat.Status)
}

type MethodName string

type MethodMap map[MethodName]EndpointMap
type EndpointMap map[ResourceName]Endpoint
type Endpoint string

type Api struct {
	Echo      *echo.Echo
	Port      string
	Defaults  *Response
	MethodMap MethodMap
}

func NewApi(echo *echo.Echo, defaults *Response) *Api {
	//
	//See https://echo.labstack.com/cookbook/cors
	//See https://flaviocopes.com/golang-enable-cors/
	//
	echo.Use(middleware.CORS())
	echo.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	echo.HideBanner = true
	echo.HidePort = true
	return &Api{
		Echo:     echo,
		Defaults: defaults,
		MethodMap: MethodMap{
			http.MethodGet:     make(EndpointMap, 0),
			http.MethodPut:     make(EndpointMap, 0),
			http.MethodPost:    make(EndpointMap, 0),
			http.MethodDelete:  make(EndpointMap, 0),
			http.MethodOptions: make(EndpointMap, 0),
		},
	}
}

func (me *Api) GetBaseUrl() (url string) {
	return fmt.Sprintf("http://127.0.0.1:%s", me.Port)
}

func (me *Api) GetUrl(name ResourceName, vars UriTemplateVars) (url string, status stat.Status) {
	var path string
	for range only.Once {
		path, status = me.GetUrlPathTemplate(name)
		if status.IsError() {
			break
		}
		path = convertEchoTemplateToUriTemplate(path)
		path = ExpandUriTemplate(path, vars)
	}
	return fmt.Sprintf("http://127.0.0.1:%s/%s", me.Port, path), status
}

func (me *Api) GetUrlPathTemplate(name ResourceName) (url string, status stat.Status) {
	for range only.Once {
		if me.Defaults == nil {
			status = stat.NewFailedStatus(&stat.Args{
				Message: fmt.Sprintf("the Defaults property is nil when accessing api for resource type '%s'",
					name,
				),
				Help:  ContactSupportHelp(),
				Error: stat.IsStatusError,
			})
			break
		}
		url, status = me.Defaults.GetUrlPathTemplate(name)
	}
	return url, status
}

type ResponseDataGetter interface {
	GetResponseData() interface{}
}

type HandlerFunc func(ctx *RequestContext) error

//func (me *Api) GET(path string, rc *RequestContext, handler HandlerFunc) *echo.Route {
//	me.Defaults.Meta.Resource = rc.ResourceName
//	me.Defaults.Links[rc.ResourceName] = convertEchoTemplateToUriTemplate(path)
//	return me.Echo.GET(path, rc.WrapHandler(handler))
//}

func (me *Api) GET(path string, name ResourceName, handler HandlerFunc) *echo.Route {
	rc := NewRequestContext(me, name)
	me.Defaults.Meta.Resource = name
	uriTemplate := convertEchoTemplateToUriTemplate(path)
	me.Defaults.Links[name] = uriTemplate
	me.MethodMap[http.MethodGet][name] = Endpoint(uriTemplate)
	return me.Echo.GET(path, rc.WrapHandler(handler))
}
func (me *Api) PUT(path string, name ResourceName, handler HandlerFunc) *echo.Route {
	rc := NewRequestContext(me, name)
	me.Defaults.Meta.Resource = name
	uriTemplate := convertEchoTemplateToUriTemplate(path)
	me.Defaults.Links[name] = uriTemplate
	me.MethodMap[http.MethodPut][name] = Endpoint(uriTemplate)
	me.MethodMap[http.MethodOptions][name] = Endpoint(uriTemplate)
	me.Echo.OPTIONS(path, optionsHandler)
	return me.Echo.PUT(path, rc.WrapHandler(handler))
}
func (me *Api) POST(path string, name ResourceName, handler HandlerFunc) *echo.Route {
	rc := NewRequestContext(me, name)
	me.Defaults.Meta.Resource = name
	uriTemplate := convertEchoTemplateToUriTemplate(path)
	me.Defaults.Links[name] = uriTemplate
	me.MethodMap[http.MethodPost][name] = Endpoint(uriTemplate)
	me.MethodMap[http.MethodOptions][name] = Endpoint(uriTemplate)
	me.Echo.OPTIONS(path, optionsHandler)
	return me.Echo.POST(path, rc.WrapHandler(handler))
}
func (me *Api) DELETE(path string, name ResourceName, handler HandlerFunc) *echo.Route {
	rc := NewRequestContext(me, name)
	me.Defaults.Meta.Resource = name
	uriTemplate := convertEchoTemplateToUriTemplate(path)
	me.Defaults.Links[name] = uriTemplate
	me.MethodMap[http.MethodDelete][name] = Endpoint(uriTemplate)
	me.MethodMap[http.MethodOptions][name] = Endpoint(uriTemplate)
	me.Echo.OPTIONS(path, optionsHandler)
	return me.Echo.DELETE(path, rc.WrapHandler(handler))
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

func (me *Api) NotYetImplemented(rc *RequestContext) interface{} {
	return stat.NewFailedStatus(&stat.Args{
		Message:    fmt.Sprintf("the '%s' resource has not been implemented yet", rc.ResourceName),
		HttpStatus: http.StatusMethodNotAllowed,
	})
}

func convertEchoTemplateToUriTemplate(url string) string {
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

func GetApiHelp(topic ResourceName, more ...string) string {
	var _more string
	if len(more) > 0 {
		_more = " " + more[0]
	}
	return fmt.Sprintf("see API docs for%s: %s", _more, GetApiDocsUrl(topic))
}

func ContactSupportHelp() string {
	return "contact support"
}
