package api

import (
	"encoding/json"
	"fmt"
	"gearbox/only"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"

	//	"github.com/labstack/echo/middleware"
	"github.com/projectcfg/projectcfg/util"
	"net/http"
	"strings"
)

const Port = "9999"

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
	return &Api{
		Echo:     echo,
		Defaults: defaults,
	}
}

type ResponseMeta struct {
	Version string `json:"version"`
	Service string `json:"service"`
	DocsUrl string `json:"docs_url"`
}

type Links map[string]string

type Response struct {
	Success    bool         `json:"success,omitempty"`
	StatusCode int          `json:"status_code"`
	Meta       ResponseMeta `json:"meta"`
	Links      Links        `json:"links"`
	Data       interface{}  `json:"data,omitempty"`
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
	StatusCode int
	Help       string
	Error      error
}
type jsonStatus struct {
	StatusCode int    `json:"status_code"`
	Help       string `json:"help"`
	Error      string `json:"error"`
}

func (me *Status) ToJson() string {
	js := &jsonStatus{
		StatusCode: me.StatusCode,
		Help:       me.Help,
		Error:      me.Error.Error(),
	}
	j, err := json.Marshal(js)
	if err != nil {
		j = []byte(`{"status_code": 500, "error":"multiple errors occurred"`)
	}
	return string(j)
}

type SuccessInspector interface {
	IsSuccess() bool
}

// @TODO Add ?format=yes to pretty print JSON
func (me *Api) JsonMarshalHandler(ctx echo.Context, js interface{}) (status *Status) {
	var err error
	for range only.Once {
		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		httpStatus := ctx.Response().Status
		var ok bool
		status, ok = js.(*Status)
		if ok {
			break
		}
		r := *me.Defaults.Clone()
		if err != nil {
			break
		}
		r.Data = js
		r.Links["self"] = convertEchoPathToUriTemplatePath(ctx.Path())
		r.StatusCode = httpStatus
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
		_ = ctx.String(status.StatusCode, status.ToJson())
	}
	return status
}

func GetCurrentActionName() string {
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

func pushPopActionName(name string, next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		actionNameStack.push(name)
		response := next(ctx)
		actionNameStack.pop()
		return response
	}
}

func (me *Api) GET(path, name string, handler echo.HandlerFunc) *echo.Route {
	me.Defaults.Links[name] = convertEchoPathToUriTemplatePath(path)
	return me.Echo.GET(path, pushPopActionName(name, handler))
}
func (me *Api) POST(path, name string, handler echo.HandlerFunc) *echo.Route {
	me.Defaults.Links[name] = convertEchoPathToUriTemplatePath(path)
	me.Echo.OPTIONS(path, optionsHandler)
	return me.Echo.POST(path, pushPopActionName(name, handler))
}
func (me *Api) DELETE(path, name string, handler echo.HandlerFunc) *echo.Route {
	me.Defaults.Links[name] = convertEchoPathToUriTemplatePath(path)
	me.Echo.OPTIONS(path, optionsHandler)
	return me.Echo.DELETE(path, pushPopActionName(name, handler))
}
func (me *Api) PUT(path, name string, handler echo.HandlerFunc) *echo.Route {
	me.Defaults.Links[name] = convertEchoPathToUriTemplatePath(path)
	me.Echo.OPTIONS(path, optionsHandler)
	return me.Echo.PUT(path, pushPopActionName(name, handler))
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
