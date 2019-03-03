package api

import (
	"encoding/json"
	"fmt"
	"gearbox/only"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

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
	StatusCode int          `json:"status_code"`
	Meta       ResponseMeta `json:"meta"`
	Links      Links        `json:"links"`
	Data       interface{}  `json:"data"`
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
	StatusCode int    `json:"status_code"`
	Help       string `json:"help"`
	Error      error  `json:"error"`
}

func (me *Status) ToJson() string {
	j, err := json.Marshal(me)
	if err != nil {
		j = []byte(`{"error":"Multiple errors occurred"`)
	}
	return string(j)
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

func (me *Api) GET(path, name string, handler echo.HandlerFunc) *echo.Route {
	me.Defaults.Links[name] = convertEchoPathToUriTemplatePath(path)
	return me.Echo.GET(path, handler)
}
func (me *Api) POST(path, name string, handler echo.HandlerFunc) *echo.Route {
	me.Defaults.Links[name] = convertEchoPathToUriTemplatePath(path)
	me.Echo.OPTIONS(path, optionsHandler)
	return me.Echo.POST(path, handler)
}
func (me *Api) DELETE(path, name string, handler echo.HandlerFunc) *echo.Route {
	me.Defaults.Links[name] = convertEchoPathToUriTemplatePath(path)
	me.Echo.OPTIONS(path, optionsHandler)
	return me.Echo.DELETE(path, handler)
}
func (me *Api) PUT(path, name string, handler echo.HandlerFunc) *echo.Route {
	me.Defaults.Links[name] = convertEchoPathToUriTemplatePath(path)
	me.Echo.OPTIONS(path, optionsHandler)
	return me.Echo.PUT(path, handler)
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
