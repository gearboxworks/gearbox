package apimodels

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/gearbox"
	"gearbox/gears"
	"gearbox/gearspec"
	"gearbox/global"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"strings"
)

const NamedStackType = "stack"

var NilNamedStack = (*NamedStack)(nil)
var _ apimodeler.ApiItemer = NilNamedStack

func NewFromGearsNamedStack(ctx *apimodeler.Context, gns *gears.NamedStack) (ns *NamedStack, sts status.Status) {
	for range only.Once {
		var gsids gearspec.Identifiers
		if ctx.GetResponseType() == global.ItemResponse {
			gsids, sts = gns.GetGearspecIds()
		}
		if is.Error(sts) {
			break
		}
		ns = &NamedStack{
			Authority: gns.Authority,
			Stackname: gns.Stackname,
			Members:   gsids,
			//ServiceMap: gns.RoleServicesMap,
		}
	}
	return ns, sts
}

type NamedStackMap map[types.Stackname]*NamedStack
type NamedStacks []*NamedStack

type NamedStack struct {
	Authority  types.AuthorityDomain `json:"authority"`
	Stackname  types.Stackname       `json:"stackname"`
	Members    gearspec.Identifiers  `json:"members,omitempty"`
	ServiceMap interface{}           `json:"services,omitempty"`
}

func NewNamedStack(ns *gears.NamedStack) *NamedStack {
	return &NamedStack{
		Authority:  ns.Authority,
		Stackname:  ns.Stackname,
		Members:    make(gearspec.Identifiers, 0),
		ServiceMap: newServiceMap(ns.RoleServicesMap),
	}
}

func (me *NamedStack) GetItemLinkMap(*apimodeler.Context) (apimodeler.LinkMap, status.Status) {
	return apimodeler.LinkMap{}, nil
}

func (me *NamedStack) GetType() apimodeler.ItemType {
	return NamedStackType
}

func (me *NamedStack) GetFullStackname() types.Stackname {
	return types.Stackname(me.GetId())
}

func (me *NamedStack) GetId() apimodeler.ItemId {
	return apimodeler.ItemId(fmt.Sprintf("%s/%s", me.Authority, me.Stackname))
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
		me.Stackname = types.Stackname(parts[1])
	}
	return sts
}

func (me *NamedStack) GetItem() (apimodeler.ApiItemer, status.Status) {
	return me, nil
}

func MakeGearboxStack(gb gearbox.Gearboxer, ns *NamedStack) (gbns *gears.NamedStack, sts status.Status) {
	//	gbns = gears.NewNamedStack(gb.GetGears(), types.StackId(ns.GetId()))
	gbns = gears.NewNamedStack(types.StackId(ns.GetId()))
	sts = gbns.Refresh()
	return gbns, sts
}

func (me *NamedStack) GetRelatedItems(ctx *apimodeler.Context, itemid apimodeler.ItemId) (list apimodeler.List, sts status.Status) {
	list = make(apimodeler.List, 0)
	return list, sts
}
