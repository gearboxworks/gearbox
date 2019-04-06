package routes

import (
	"fmt"
	"gearbox/gearbox"
	"gearbox/gears"
	"gearbox/modeler"
	"gearbox/status"
	"gearbox/types"
)

const StackTypeName = "stack"

var StackInstance = (*NamedStack)(nil)
var _ modeler.Item = StackInstance

func (me *NamedStack) GetType() modeler.ItemType {
	return StackTypeName
}

func (me *NamedStack) GetFullStackname() types.Stackname {
	return types.Stackname(me.GetId())
}

func (me *NamedStack) GetId() modeler.ItemId {
	return modeler.ItemId(fmt.Sprintf("%s/%s", me.Authority, me.StackName))
}

func (me *NamedStack) GetItem() (modeler.Item, status.Status) {
	return me, nil
}

func MakeGearboxStack(gb gearbox.Gearboxer, ns *NamedStack) (gbns *gears.NamedStack, sts status.Status) {
	gbns = gears.NewNamedStack(gb.GetGears(), types.StackId(ns.GetId()))
	sts = gbns.Refresh()
	return gbns, sts
}
