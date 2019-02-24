package gearbox

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/projectcfg/projectcfg/util"
	"net/http"
)

const Port = "9999"

type Api struct {
	Port   string
	Config *Config
	Echo   *echo.Echo
}

type ResponseMeta struct {
	Version string `json:"version"`
	Name    string `json:"name"`
	Docs    string `json:"docs"`
}

type Response struct {
	Meta  ResponseMeta      `json:"meta"`
	Links map[string]string `json:"links"`
	Data  interface{}       `json:"data"`
}

const ApiResponse = `{
	"meta":{
		"version":"0.1",
		"name":"GearBox API",
		"docs":"https://docs.gearbox.works/api"
	},
	"links":{
		"root":          	  "/",
		"projects":           "/projects",
		"enabled-projects":   "/projects/enabled",
		"disabled-projects":  "/projects/disabled",
		"candidate-projects": "/projects/candidates",
		"stacks":             "/stacks",
		"stack":              "/stacks/{stack}",
		"services":           "/services/{service_id}"
	}
}`

func NewApi(conf *Config) *Api {
	e := echo.New()

	api := &Api{
		Port:   Port,
		Config: conf,
		Echo:   e,
	}
	e.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "{}")
	})
	e.GET("/projects", func(ctx echo.Context) error {
		return jsonMarshalHandler(api, ctx, api.Config.Projects)
	})
	e.GET("/projects/:project", func(ctx echo.Context) error {
		return jsonMarshalHandler(api, ctx, api.Config.Projects)
	})
	e.GET("/projects/enabled", func(ctx echo.Context) error {
		return jsonMarshalHandler(api, ctx, api.Config.Projects.GetEnabled())
	})
	e.GET("/projects/disabled", func(ctx echo.Context) error {
		return jsonMarshalHandler(api, ctx, api.Config.Projects.GetDisabled())
	})
	e.GET("/projects/candidates", func(ctx echo.Context) error {
		return jsonMarshalHandler(api, ctx, api.Config.Candidates)
	})
	e.POST("/projects/new", func(ctx echo.Context) error {
		return jsonMarshalHandler(api, ctx, api.Config.Candidates)
	})

	e.GET("/stacks", func(ctx echo.Context) error {
		return jsonMarshalHandler(api, ctx, []string{
			"WordPress",
			"Drupal",
			"Joomla",
		})
	})

	e.GET("/stacks/:stack", func(ctx echo.Context) error {
		return jsonMarshalHandler(api, ctx, []string{
			"wordpress/dbserver",
			"wordpress/webserver",
			"wordpress/processvm",
			"wordpress/cacheserver",
		})
	})

	e.GET("/services/:service_id", func(ctx echo.Context) error {
		return jsonMarshalHandler(api, ctx, []string{
			"gearbox/mariadb:5.5",
			"gearbox/mariadb:10.0",
			"gearbox/mariadb:10.1",
			"gearbox/mariadb:10.2",
			"gearbox/mariadb:10.3",
			"gearbox/mysql:5.5",
			"gearbox/mysql:5.6",
			"gearbox/mysql:5.7",
			"gearbox/mysql:8.0",
		})
	})
	e.GET("/whatever", func(ctx echo.Context) error {
		return jsonMarshalHandler(api, ctx, map[string]string{
			"foo": "bar",
			"bar": "baz",
		})
	})

	return api
}

func (me *Api) Start() {
	err := me.Echo.Start(":" + me.Port)
	if err != nil {
		util.Error(err)
	}
}

// @TODO Add ?format=yes to pretty print JSON
func jsonMarshalHandler(api *Api, ctx echo.Context, js interface{}) error {
	r := &Response{}
	err := json.Unmarshal([]byte(ApiResponse), &r)
	if err != nil {
		panic(err)
	}
	r.Data = js
	r.Links["self"] = ctx.Path()
	j, err := json.MarshalIndent(r, "", "   ")
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	return ctx.String(http.StatusOK, string(j))
}
