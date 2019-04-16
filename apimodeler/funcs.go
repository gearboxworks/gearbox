package apimodeler

import (
	"fmt"
	"gearbox/only"
	"gearbox/status"
	"gearbox/types"
	"net/http"
	"reflect"
	"strings"
)

func GetBasepath(me ListController) types.Basepath {
	getter, ok := me.(BasepathGetter)
	if !ok {
		panic("API controller does not have GetBasepath()")
	}
	return getter.GetBasepath()
}

func GetIdParams(me ListController) IdParams {
	getter, ok := me.(IdParamsGetter)
	if !ok {
		panic("API controller does not have GetIdParams()")
	}
	return getter.GetIdParams()
}

func GetResourceUrlTemplate(controller ListController) (ut types.UrlTemplate) {
	for range only.Once {
		bp := controller.GetBasepath()
		idt := GetIdTemplate(controller)
		if idt == "" {
			ut = types.UrlTemplate(bp)
			break
		}
		ut = types.UrlTemplate(fmt.Sprintf("%s/%s", bp, idt))
	}
	return ut
}

func GetRouteNamePrefix(controller ListController) types.RouteName {
	rn := strings.Trim(string(GetBasepath(controller)), "/")
	if !reflect.ValueOf(controller).IsNil() {
		prn := GetRouteNamePrefix(controller.GetParent())
		if prn != "" {
			rn = fmt.Sprintf("%s-%s", prn, rn)
		}
	}
	return types.RouteName(rn)
}

func GetIdFromUrl(ctx *Context, controller ListController) (id ItemId, sts status.Status) {
	for range only.Once {
		params := controller.GetIdParams()
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

func GetIdTemplate(controller ListController) types.UrlTemplate {
	return types.UrlTemplate(controller.GetIdParams().Slugify())
}

func noop(i ...interface{}) {}
