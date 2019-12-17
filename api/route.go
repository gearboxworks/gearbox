package api

import (
	"fmt"
	"gearbox/types"
	"github.com/gedex/inflector"
	"github.com/labstack/echo"
)

type Route struct {
	*echo.Route
}

func (me *Route) SetName(name types.RouteName) *Route {
	me.Route.Name = string(name)
	return me
}
func (me *Route) SetSingularName(pattern, prefix string) *Route {
	me.Route.Name = fmt.Sprintf(pattern, inflector.Singularize(prefix))
	return me
}
func (me *Route) SetPluralName(pattern, prefix string) *Route {
	me.Route.Name = fmt.Sprintf(pattern, inflector.Pluralize(prefix))
	return me
}
