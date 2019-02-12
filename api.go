package gearbox

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/projectcfg/projectcfg/util"
	"net/http"
)

const Port = "19970"

type Api struct {
	Config *Config
	Echo   *echo.Echo
}

func NewApi(conf *Config) *Api {
	e := echo.New()

	api := &Api{
		Config: conf,
		Echo:   e,
	}
	e.GET("/", func(ctx echo.Context) error {
		return rootHandler(api, ctx)
	})
	e.GET("/projects", func(ctx echo.Context) error {
		return projectsHandler(api, ctx)
	})

	return api
}

func (me *Api) Run() {
	err := me.Echo.Start(":" + Port)
	if err != nil {
		util.Error(err)
	}
}

func rootHandler(api *Api, context echo.Context) error {
	return context.String(http.StatusOK, "Hello. Thanks for playing.")
}

func projectsHandler(api *Api, context echo.Context) error {
	j, err := json.MarshalIndent(api.Config.Projects, "", "   ")
	if err != nil {
		return context.String(http.StatusInternalServerError, err.Error())
	}
	return context.String(http.StatusOK, string(j))
}
