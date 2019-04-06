package modeler

import (
	"fmt"
	"gearbox/only"
	"gearbox/status"
	"gearbox/types"
	"github.com/labstack/echo"
	"net/http"
	"strings"
)

type ModelsMap map[types.Basepath]*Models
type Models struct {
	Self     Modeler
	Parent   *Models
	Children ModelsMap
}

func NewModels(self Modeler) *Models {
	return &Models{
		Self:     self,
		Children: make(ModelsMap, 0),
	}
}

func (me *Models) GetBasepath() types.Basepath {
	getter, ok := me.Self.(BasepathGetter)
	if !ok {
		panic("controller does not have GetBasepath()")
	}
	return getter.GetBasepath()
}

func (me *Models) GetIdFromUrl(ctx echo.Context) (id ItemId, sts status.Status) {
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

func (me *Models) GetIdTemplate() types.UrlTemplate {
	return types.UrlTemplate(me.GetIdParams().Slugify())
}

func (me *Models) GetIdParams() IdParams {
	getter, ok := me.Self.(IdParamsGetter)
	if !ok {
		panic("factory does not have GetIdParams()")
	}
	return getter.GetIdParams()
}

func (me *Models) GetResourceUrlTemplate() (ut types.UrlTemplate) {
	return types.UrlTemplate(fmt.Sprintf("%s/%s",
		string(me.GetBasepath()),
		string(me.GetIdTemplate()),
	))
}

func (me *Models) GetRouteNamePrefix() (rn string) {
	rn = strings.Trim(string(me.GetBasepath()), "/")
	if me.Parent != nil {
		prn := me.Parent.GetRouteNamePrefix()
		if prn != "" {
			rn = fmt.Sprintf("%s-%s", prn, rn)
		}
	}
	return rn
}
