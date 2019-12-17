package apimvc

import (
	"gearbox/apiworks"
	"gearbox/gearbox"
	"gearbox/types"
)

const Rootname = "root"

var NilRootController = (*RootController)(nil)
var _ apiworks.ListController = NilRootController

type RootController struct {
	*apiworks.Controller
	Gearbox gearbox.Gearboxer
}

func NewRootController(gb gearbox.Gearboxer) *RootController {
	return &RootController{
		Gearbox:    gb,
		Controller: apiworks.NewController(),
	}
}

func (me *RootController) GetName() types.RouteName {
	return Rootname
}

func (me *RootController) GetNilItem(ctx *Context) ItemModeler {
	panic("not yet implemented")
}
