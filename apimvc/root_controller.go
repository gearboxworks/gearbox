package apimvc

import (
	"gearbox/apimodeler"
	"gearbox/gearbox"
	"gearbox/types"
)

const Rootname = "root"

var NilRootController = (*RootController)(nil)
var _ apimodeler.ListController = NilRootController

type RootController struct {
	*apimodeler.Controller
	Gearbox gearbox.Gearboxer
}

func NewRootController(gb gearbox.Gearboxer) *RootController {
	return &RootController{
		Gearbox:    gb,
		Controller: apimodeler.NewController(),
	}
}

func (me *RootController) GetName() types.RouteName {
	return Rootname
}
