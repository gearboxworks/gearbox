package api

//
//See https://echo.labstack.com/cookbook/cors
//See https://flaviocopes.com/golang-enable-cors/
//
import (
	"fmt"
	"gearbox/only"
	"gearbox/status"
	"gearbox/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"strings"
)

const RequestContextKey = "api-request-context"
const BaseUrlPattern = "http://127.0.0.1:%s"

type Api struct {
	Echo          *echo.Echo
	Port          string
	Defaults      *Response
	EndpointMap   EndpointMap
	MethodMap     MethodMap
	ValuesFuncMap ValuesFuncMap
	RelatedMap    ResourcesMap
}

func NewApi(e *echo.Echo, defaults *Response) *Api {
	a := Api{
		Defaults:    defaults,
		EndpointMap: make(EndpointMap, 0),
		MethodMap: MethodMap{
			http.MethodGet:     make(ResourceMap, 0),
			http.MethodPut:     make(ResourceMap, 0),
			http.MethodPost:    make(ResourceMap, 0),
			http.MethodDelete:  make(ResourceMap, 0),
			http.MethodOptions: make(ResourceMap, 0),
		},
		ValuesFuncMap: make(ValuesFuncMap, 0),
		RelatedMap:    make(ResourcesMap, 0),
	}

	e.Use(middleware.CORS())
	e.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		_, ok := c.Get(RequestContextKey).(*RequestContext)
		if !ok {
			e.DefaultHTTPErrorHandler(err, c)
		}
	}
	a.Echo = e
	return &a
}

func (me *Api) GET(path UriTemplate, name RouteName, valuesFunc ValuesFunc, handler HandlerFunc) *echo.Route {
	r := me.Echo.GET(string(path), me.GetRequestContext(http.MethodGet, path, name, valuesFunc).WrapHandler(handler))
	r.Name = string(name)
	return r
}
func (me *Api) PUT(path UriTemplate, name RouteName, valuesFunc ValuesFunc, handler HandlerFunc) *echo.Route {
	r := me.Echo.PUT(string(path), me.GetRequestContext(http.MethodPut, path, name, valuesFunc).WrapHandler(handler))
	r.Name = string(name)
	return r
}

func (me *Api) POST(path UriTemplate, name RouteName, valuesFunc ValuesFunc, handler HandlerFunc) *echo.Route {
	r := me.Echo.POST(string(path), me.GetRequestContext(http.MethodPost, path, name, valuesFunc).WrapHandler(handler))
	r.Name = string(name)
	return r
}
func (me *Api) DELETE(path UriTemplate, name RouteName, valuesFunc ValuesFunc, handler HandlerFunc) *echo.Route {
	r := me.Echo.DELETE(string(path), me.GetRequestContext(http.MethodDelete, path, name, valuesFunc).WrapHandler(handler))
	r.Name = string(name)
	return r
}

func (me *Api) GetRequestContext(method Method, path UriTemplate, name RouteName, valuesFunc ValuesFunc) *RequestContext {
	me.ValuesFuncMap[name] = valuesFunc
	me.Defaults.Meta.RouteName = name
	uriTemplate := path.Convert()
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

func (me *Api) GetBaseUrl() (url UriTemplate) {
	return UriTemplate(fmt.Sprintf(BaseUrlPattern, me.Port))
}

func (me *Api) GetUrlPath(name RouteName, vars UriTemplateVars) (urlpath UriTemplate, sts status.Status) {
	var path UriTemplate
	for range only.Once {
		path, sts = me.GetUrlPathTemplate(name)
		if status.IsError(sts) {
			break
		}
		if len(vars) == 0 {
			break
		}
		path = path.Expand(vars)
	}
	urlpath = UriTemplate(strings.TrimRight(string(path), "/"))
	return urlpath, sts
}

func (me *Api) GetUrl(name RouteName, vars UriTemplateVars) (url UriTemplate, sts status.Status) {
	var path UriTemplate
	for range only.Once {
		path, sts = me.GetUrlPath(name, vars)
		if status.IsError(sts) {
			break
		}
	}
	url = UriTemplate(fmt.Sprintf("%s/%s",
		me.GetBaseUrl(),
		strings.TrimLeft(string(path), "/"),
	))
	return url, sts
}

func (me *Api) GetUriTemplateVars(name RouteName, values interface{}, index int) (utvars UriTemplateVars, sts status.Status) {
	for range only.Once {
		var vars TemplateVars
		vars, sts = me.GetUrlVars(name)
		if status.IsError(sts) {
			break
		}
		if len(vars) == 0 {
			break
		}
		vals, ok := values.(ValuesFuncValues)
		if !ok {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("values for '%s' does not support %d URL params, e.g. %s",
					name,
					len(vars),
					util.OxfordComma(vars.Values()),
				),
			})
			break
		}
		utvars = make(UriTemplateVars, len(vars))
		for i, v := range vars {
			utvars[i] = &UriTemplateVar{
				Name:  ResourceVarName(v),
				Value: string(vals[i][index]),
			}
		}
	}
	return utvars, sts
}

func (me *Api) GetUrlVars(name RouteName) (vars TemplateVars, sts status.Status) {
	for range only.Once {
		var ut UriTemplate
		ut, sts = me.GetUrlPathTemplate(name)
		if status.IsError(sts) {
			break
		}
		parts := strings.Split(string(ut), "/")
		vars = make(TemplateVars, 0)
		for _, s := range parts {
			if len(s) > 0 && s[0] == '{' && s[len(s)-1] == '}' {
				vars = append(vars, TemplateVar(s[1:len(s)-1]))
			}
		}
	}
	return vars, sts
}

func (me *Api) GetUrlPathTemplate(name RouteName) (url UriTemplate, sts status.Status) {
	for range only.Once {
		if me.Defaults == nil {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("the Defaults property is nil when accessing api for resource type '%s'",
					name,
				),
				Help: ContactSupportHelp(),
			})
			break
		}
		url, sts = me.Defaults.GetUrlPathTemplate(name)
	}
	return url, sts
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
	return status.Fail(&status.Args{
		Message:    fmt.Sprintf("the '%s' resource has not been implemented yet", rc.RouteName),
		HttpStatus: http.StatusMethodNotAllowed,
	})
}
