package ab

import (
	"fmt"
	"gearbox/only"
	"gearbox/status"
	"github.com/labstack/echo"
	"net/http"
	"strings"
)

type IdParams []IdParam
type IdParam string

func (me IdParams) Slugify() UrlTemplate {
	ss := make([]string, len(me))
	for i, idp := range me {
		ss[i] = ":" + string(idp)
	}
	return UrlTemplate(strings.Join(ss, "/"))
}

type Basepath string

type ConnectionsMap map[Basepath]*Connections
type Connections struct {
	Self     Connector
	Parent   *Connections
	Children ConnectionsMap
}

func NewConnections(self Connector) *Connections {
	return &Connections{
		Self:     self,
		Children: make(ConnectionsMap, 0),
	}
}

type UrlTemplates []UrlTemplate
type UrlTemplate string

func (me *Connections) GetBasepath() Basepath {
	getter, ok := me.Self.(BasepathGetter)
	if !ok {
		panic("factory does not have GetBasepath()")
	}
	return getter.GetBasepath()
}

func (me *Connections) GetIdFromUrl(ctx echo.Context) (id ItemId, sts status.Status) {
	for range only.Once {
		params := me.GetIdParams()
		parts := make([]string, len(params))
		for i, idp := range params {
			val := ctx.Param(string(idp))
			if val == "" {
				sts = status.Fail(&status.Args{
					Message:    fmt.Sprintf("URL path segment for '%s' is empty", idp),
					HttpStatus: http.StatusBadRequest,
				})
				break
			}
			parts[i] = val
		}
		id = ItemId(strings.Join(parts, "/"))
	}
	return id, sts
}

func (me *Connections) GetIdTemplate() UrlTemplate {
	return UrlTemplate(me.GetIdParams().Slugify())
}

func (me *Connections) GetIdParams() IdParams {
	getter, ok := me.Self.(IdParamsGetter)
	if !ok {
		panic("factory does not have GetIdParams()")
	}
	return getter.GetIdParams()
}

func (me *Connections) GetResourceUrlTemplate() (ut UrlTemplate) {
	return UrlTemplate(fmt.Sprintf("%s/%s",
		string(me.GetBasepath()),
		string(me.GetIdTemplate()),
	))
}

func (me *Connections) GetRouteNamePrefix() (rn string) {
	rn = strings.Trim(string(me.GetBasepath()), "/")
	if me.Parent != nil {
		prn := me.Parent.GetRouteNamePrefix()
		if prn != "" {
			rn = fmt.Sprintf("%s-%s", prn, rn)
		}
	}
	return rn
}
