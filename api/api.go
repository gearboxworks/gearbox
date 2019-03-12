package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"gearbox/only"
	"gearbox/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	//	"github.com/labstack/echo/middleware"
	"net/http"
	"strings"
)

const Port = "9999"

type ResourceVarName string
type ResourceName string

const SelfResource ResourceName = "self"
const LinksResource ResourceName = "links"

type Links map[ResourceName]string

type ListItemResponseMap map[string]*ListItemResponse

type ListItemResponse struct {
	Links Links       `json:"links"`
	Data  interface{} `json:"data"`
}

func NewListItemResponse(link string, data interface{}) *ListItemResponse {
	return &ListItemResponse{
		Links: Links{
			SelfResource: link,
		},
		Data: data,
	}
}

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

type ResourceType string

type ResponseMeta struct {
	Version  string       `json:"version"`
	Service  string       `json:"service"`
	DocsUrl  string       `json:"docs_url"`
	Resource ResourceName `json:"resource"`
}

type Response struct {
	Success    bool         `json:"success"`
	StatusCode int          `json:"status_code"`
	Meta       ResponseMeta `json:"meta"`
	Links      Links        `json:"links"`
	Data       interface{}  `json:"data,omitempty"`
}

func (me *Response) GetApiSelfLink(resourceType ResourceName) (url string, err error) {
	for range only.Once {
		var ok bool
		url, ok = me.Links[resourceType]
		if !ok {
			err = util.AddHelpToError(
				errors.New(fmt.Sprintf("no '%s' in resource links", resourceType)),
				ContactSupportHelp(),
			)
		}
	}
	return url, err
}

func (me *Response) Clone() *Response {
	r := Response{}
	for range only.Once {
		b, err := json.Marshal(me)
		if err != nil {
			break
		}
		_ = json.Unmarshal(b, &r)
	}
	return &r
}

type Status struct {
	Success    bool
	StatusCode int
	Help       string
	Error      error
}
type StatusResponse struct {
	Help  string `json:"help"`
	Error string `json:"error"`
}

func (me *Status) ToResponse() *StatusResponse {
	return &StatusResponse{
		Help:  me.Help,
		Error: me.Error.Error(),
	}
}

type SuccessInspector interface {
	IsSuccess() bool
}

type ResponseDataGetter interface {
	GetResponseData() interface{}
}

// @TODO Add ?format=yes to pretty print JSON
func (me *Api) JsonMarshalHandler(ctx echo.Context, requestType ResourceName, js interface{}) (status *Status) {
	var err error
	for range only.Once {
		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		var ok bool
		success := true
		status, ok = js.(*Status)
		if ok {
			success = false
			js = status.ToResponse()
			ctx.Response().Status = status.StatusCode
		}
		httpStatus := ctx.Response().Status
		r := *me.Defaults.Clone()
		if rdg, ok := js.(ResponseDataGetter); !ok {
			r.Data = js
		} else {
			r.Data = rdg.(ResponseDataGetter).GetResponseData()
		}
		path := ctx.Path()
		if path != "/" {
			r.Links = make(Links, 0)
		}
		if slg, ok := js.(SelfLinkGetter); ok {
			r.Links[SelfResource] = slg.GetApiSelfLink()
		} else {
			r.Links[SelfResource] = convertEchoPathToUriTemplatePath(path)
		}
		r.Meta.DocsUrl = fmt.Sprintf("%s/%s", r.Meta.DocsUrl, string(requestType))
		r.Meta.Resource = requestType
		r.StatusCode = httpStatus
		r.Success = success
		if si, ok := js.(SuccessInspector); ok {
			r.Success = si.IsSuccess()
		}
		var j []byte
		j, err = json.MarshalIndent(r, "", "   ")
		if err != nil {
			break
		}
		err = ctx.String(httpStatus, string(j))
		status = &Status{StatusCode: httpStatus}
	}
	if status == nil && err != nil {
		status = &Status{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}
	if status.Error != nil {
		b, _ := json.Marshal(status.ToResponse())
		_ = ctx.String(status.StatusCode, string(b))
	}
	return status
}

func GetCurrentRequestType() string {
	return actionNameStack.get()
}

type stack []string

var actionNameStack stack

func init() {
	actionNameStack = make(stack, 0)
}
func (me stack) push(name string) {
	actionNameStack = append(me, name)
}
func (me stack) get() string {
	if len(me) == 0 {
		return ""
	}
	return me[len(me)-1]
}
func (me stack) pop() string {
	l := len(me)
	if l == 0 {
		log.Fatal("attempt to pop from api.actionNameStack when empty.")
	}
	pop := me[l-1]
	actionNameStack = me[:l-1]
	return pop
}

func passRequestType(requestType ResourceName, next HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return next(requestType, ctx)
	}
}

type HandlerFunc func(requestType ResourceName, ctx echo.Context) error

func (me *Api) GET(path string, requestType ResourceName, handler HandlerFunc) *echo.Route {
	me.Defaults.Meta.Resource = requestType
	me.Defaults.Links[requestType] = convertEchoPathToUriTemplatePath(path)
	return me.Echo.GET(path, passRequestType(requestType, handler))
}
func (me *Api) POST(path string, requestType ResourceName, handler HandlerFunc) *echo.Route {
	me.Defaults.Meta.Resource = requestType
	me.Defaults.Links[requestType] = convertEchoPathToUriTemplatePath(path)
	me.Echo.OPTIONS(path, optionsHandler)
	return me.Echo.POST(path, passRequestType(requestType, handler))
}
func (me *Api) DELETE(path string, requestType ResourceName, handler HandlerFunc) *echo.Route {
	me.Defaults.Meta.Resource = requestType
	me.Defaults.Links[requestType] = convertEchoPathToUriTemplatePath(path)
	me.Echo.OPTIONS(path, optionsHandler)
	return me.Echo.DELETE(path, passRequestType(requestType, handler))
}
func (me *Api) PUT(path string, requestType ResourceName, handler HandlerFunc) *echo.Route {
	me.Defaults.Meta.Resource = requestType
	me.Defaults.Links[requestType] = convertEchoPathToUriTemplatePath(path)
	me.Echo.OPTIONS(path, optionsHandler)
	return me.Echo.PUT(path, passRequestType(requestType, handler))
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

func ReadRequestBody(ctx echo.Context) ([]byte, error) {
	return ioutil.ReadAll(ctx.Request().Body)
}
func CloseRequestBody(ctx echo.Context) {
	_ = ctx.Request().Body.Close()
}

type UriTemplateVars map[ResourceVarName]string

func ExpandUriTemplate(template string, vars UriTemplateVars) string {
	url := template
	for vn, val := range vars {
		url = strings.Replace(url, fmt.Sprintf("{%s}", vn), val, -1)
	}
	return url
}
