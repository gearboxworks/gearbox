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

type Error struct {
	StatusCode int   `json:"status_code"`
	Error      error `json:"error"`
}

// @TODO Add ?format=yes to pretty print JSON
func (me *Api) JsonMarshalHandler(ctx echo.Context, js interface{}) error {
	var err error
	for range only.Once {
		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ae, ok := js.(*Error)
		if ok {
			err = ctx.String(ae.StatusCode, ae.Error.Error())
			break
		}
		r := *me.Defaults.Clone()
		if err != nil {
			err = ctx.String(http.StatusInternalServerError, err.Error())
			break
		}
		r.Data = js
		r.Links["self"] = convertEchoPathToUriTemplatePath(ctx.Path())
		j, err := json.MarshalIndent(r, "", "   ")
		if err != nil {
			err = ctx.String(http.StatusInternalServerError, err.Error())
			break
		}
		err = ctx.String(http.StatusOK, string(j))
	}
	return err
}

func (me *Api) GET(path, name string, handler echo.HandlerFunc) *echo.Route {
	me.Defaults.Links[name] = convertEchoPathToUriTemplatePath(path)
	return me.Echo.GET(path, handler)
}
func (me *Api) POST(path, name string, handler echo.HandlerFunc) *echo.Route {
	me.Defaults.Links[name] = convertEchoPathToUriTemplatePath(path)
	return me.Echo.POST(path, handler)
}
func (me *Api) DELETE(path, name string, handler echo.HandlerFunc) *echo.Route {
	me.Defaults.Links[name] = convertEchoPathToUriTemplatePath(path)
	return me.Echo.DELETE(path, handler)
}
func (me *Api) PUT(path, name string, handler echo.HandlerFunc) *echo.Route {
	me.Defaults.Links[name] = convertEchoPathToUriTemplatePath(path)
	return me.Echo.PUT(path, handler)
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
