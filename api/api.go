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
const MetaEndpointsResource ResourceName = "meta-endpoints"
const MetaMethodsResource ResourceName = "meta-methods"

type Links map[ResourceName]UriTemplate

type UrlGetter interface {
	GetApiUrl(...ResourceName) (UriTemplate, stat.Status)
}

type ResponseDataGetter interface {
	GetResponseData() interface{}
}

type HandlerFunc func(ctx *RequestContext) error

func optionsHandler(ctx echo.Context) error {
	return nil
}

type UriTemplate string

type EndpointMap map[ResourceName]*Endpoint
type Endpoint struct {
	UriTemplate UriTemplate `json:"uri_template"`
	Methods     Methods     `json:"methods"`
}

type ResourceMap map[ResourceName]UriTemplate
type MethodMap map[Method]ResourceMap
type Methods []Method
type Method string

type Api struct {
	Echo        *echo.Echo
	Port        string
	Defaults    *Response
	EndpointMap EndpointMap
	MethodMap   MethodMap
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
		Echo:        echo,
		Defaults:    defaults,
		EndpointMap: make(EndpointMap, 0),
		MethodMap: MethodMap{
			http.MethodGet:     make(ResourceMap, 0),
			http.MethodPut:     make(ResourceMap, 0),
			http.MethodPost:    make(ResourceMap, 0),
			http.MethodDelete:  make(ResourceMap, 0),
			http.MethodOptions: make(ResourceMap, 0),
		},
	}
}

func (me *Api) GET(path UriTemplate, name ResourceName, handler HandlerFunc) *echo.Route {
	return me.Echo.GET(string(path), me.GetRequestContext(http.MethodGet, name, path).WrapHandler(handler))
}
func (me *Api) PUT(path UriTemplate, name ResourceName, handler HandlerFunc) *echo.Route {
	return me.Echo.PUT(string(path), me.GetRequestContext(http.MethodPut, name, path).WrapHandler(handler))
}

func (me *Api) POST(path UriTemplate, name ResourceName, handler HandlerFunc) *echo.Route {
	return me.Echo.POST(string(path), me.GetRequestContext(http.MethodPost, name, path).WrapHandler(handler))
}
func (me *Api) DELETE(path UriTemplate, name ResourceName, handler HandlerFunc) *echo.Route {
	return me.Echo.DELETE(string(path), me.GetRequestContext(http.MethodDelete, name, path).WrapHandler(handler))
}

func (me *Api) GetRequestContext(method Method, name ResourceName, path UriTemplate) *RequestContext {
	me.Defaults.Meta.Resource = name
	uriTemplate := convertEchoTemplateToUriTemplate(path)
	me.Defaults.Links[name] = uriTemplate
	me.MethodMap[method][name] = uriTemplate
	me.MethodMap[http.MethodOptions][name] = uriTemplate

	if me.EndpointMap[name] == nil {
		me.EndpointMap[name] = &Endpoint{
			UriTemplate: uriTemplate,
			Methods:     make(Methods, 0),
		}
	}
	me.EndpointMap[name].Methods = append(
		me.EndpointMap[name].Methods,
		method,
	)
	if method != http.MethodGet {
		me.Echo.OPTIONS(string(path), optionsHandler)
		me.EndpointMap[name].Methods = append(
			me.EndpointMap[name].Methods,
			http.MethodOptions,
		)
	}
	return NewRequestContext(me, name)
}

func (me *Api) GetBaseUrl() (url string) {
	return fmt.Sprintf("http://127.0.0.1:%s", me.Port)
}

func (me *Api) GetUrl(name ResourceName, vars UriTemplateVars) (url UriTemplate, status stat.Status) {
	var path UriTemplate
	for range only.Once {
		path, status = me.GetUrlPathTemplate(name)
		if status.IsError() {
			break
		}
		path = convertEchoTemplateToUriTemplate(path)
		path = ExpandUriTemplate(path, vars)
	}
	url = UriTemplate(fmt.Sprintf("http://127.0.0.1:%s/%s", me.Port, path))
	return url, status
}

func (me *Api) GetUrlPathTemplate(name ResourceName) (url UriTemplate, status stat.Status) {
	for range only.Once {
		if me.Defaults == nil {
			status = stat.NewFailStatus(&stat.Args{
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
	return stat.NewFailStatus(&stat.Args{
		Message:    fmt.Sprintf("the '%s' resource has not been implemented yet", rc.ResourceName),
		HttpStatus: http.StatusMethodNotAllowed,
	})
}

func convertEchoTemplateToUriTemplate(url UriTemplate) UriTemplate {
	parts := strings.Split(string(url), "/")
	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		if []byte(p)[0] != ':' {
			continue
		}
		parts[i] = fmt.Sprintf("{%s}", p[1:])
	}
	return UriTemplate(strings.Join(parts, "/"))
}

type UriTemplateVars map[ResourceVarName]string

func ExpandUriTemplate(template UriTemplate, vars UriTemplateVars) UriTemplate {
	url := template
	for vn, val := range vars {
		url = UriTemplate(strings.Replace(string(url), fmt.Sprintf("{%s}", vn), val, -1))
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
