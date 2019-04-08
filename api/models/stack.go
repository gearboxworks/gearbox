package models

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/gearbox"
	"gearbox/gears"
	"gearbox/only"
	"gearbox/status"
	"gearbox/types"
	"strings"
)

const StackTypeName = "stack"

var StackInstance = (*NamedStack)(nil)
var _ apimodeler.Itemer = StackInstance

func (me *NamedStack) GetType() apimodeler.ItemType {
	return StackTypeName
}

func (me *NamedStack) GetFullStackname() types.Stackname {
	return types.Stackname(me.GetId())
}

func (me *NamedStack) GetId() apimodeler.ItemId {
	return apimodeler.ItemId(fmt.Sprintf("%s/%s", me.Authority, me.StackName))
}

func (me *NamedStack) SetId(itemid apimodeler.ItemId) (sts status.Status) {
	for range only.Once {
		parts := strings.Split(string(itemid), "/")
		if len(parts) < 2 {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("stack ID '%s' missing '/'", itemid),
			})
			break
		} else if len(parts) > 2 {
			sts = status.Fail(&status.Args{
				Message: fmt.Sprintf("stack ID '%s' has too many '/'", itemid),
			})
			break
		}
		me.Authority = types.AuthorityDomain(parts[0])
		me.StackName = types.Stackname(parts[1])
	}
	return sts
}

func (me *NamedStack) GetItem() (apimodeler.Itemer, status.Status) {
	return me, nil
}

func MakeGearboxStack(gb gearbox.Gearboxer, ns *NamedStack) (gbns *gears.NamedStack, sts status.Status) {
	gbns = gears.NewNamedStack(gb.GetGears(), types.StackId(ns.GetId()))
	sts = gbns.Refresh()
	return gbns, sts
}
