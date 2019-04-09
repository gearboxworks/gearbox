package models

import (
	"gearbox/apimodeler"
	"gearbox/gearbox"
	"gearbox/types"
)

const Rootname = "root"

var NilRootModel = (*RootModel)(nil)
var _ apimodeler.Modeler = NilRootModel

type RootModel struct {
	*apimodeler.BaseModel
	Gearbox gearbox.Gearboxer
}

func NewRootModel(gb gearbox.Gearboxer) *RootModel {
	return &RootModel{
		Gearbox:   gb,
		BaseModel: apimodeler.NewBaseModel(),
	}
}

func (me *RootModel) GetName() types.RouteName {
	return Rootname
}
