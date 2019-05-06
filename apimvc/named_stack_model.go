package apimvc

import (
	"fmt"
	"gearbox/gearbox"
	"gearbox/gears"
	"gearbox/only"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"strings"
)

const NamedStackType ItemType = "stack"

var NilNamedStackModel = (*NamedStackModel)(nil)
var _ ItemModeler = NilNamedStackModel

type NamedStackModelMap map[types.Stackname]*NamedStackModel
type NamedStackModels []*NamedStackModel

type NamedStackModel struct {
	Authority types.AuthorityDomain `json:"authority"`
	Stackname types.Stackname       `json:"stackname"`
	Members   StackMembers          `json:"members,omitempty"`
}

func NewNamedStackModelFromGearsNamedStack(ctx *Context, gns *gears.NamedStack) (ns *NamedStackModel, sts Status) {
	for range only.Once {

		gsom, sts := gns.GetServiceOptionMap()
		if is.Error(sts) {
			break
		}
		gsrm, sts := gns.GetRoleMap()
		if is.Error(sts) {
			break
		}
		sms := make(StackMembers, len(gsom))
		i := 0
		for gs, gso := range gsom {
			sr, ok := gsrm[gs]
			if !ok {
				sr = &gears.StackRole{}
			}
			var sm *StackMember
			sm = NewStackMemberFromGearsStackRoleServices(ctx, gso, sr)
			sms[i] = sm
			i++
		}
		ns = &NamedStackModel{
			Authority: gns.Authority,
			Stackname: gns.Stackname,
			Members:   sms,
		}
	}
	return ns, sts
}

func NewNamedStackModel(ns *gears.NamedStack) *NamedStackModel {
	return &NamedStackModel{
		Authority: ns.Authority,
		Stackname: ns.Stackname,
		Members:   make(StackMembers, 0),
	}
}

func (me *NamedStackModel) GetItemLinkMap(*Context) (LinkMap, Status) {
	return LinkMap{}, nil
}

func (me *NamedStackModel) GetType() ItemType {
	return NamedStackType
}

func (me *NamedStackModel) GetFullStackname() types.Stackname {
	return types.Stackname(me.GetId())
}

func (me *NamedStackModel) GetId() ItemId {
	return ItemId(fmt.Sprintf("%s/%s", me.Authority, me.Stackname))
}

func (me *NamedStackModel) SetId(itemid ItemId) (sts Status) {
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

func (me *NamedStackModel) GetItem() (ItemModeler, Status) {
	return me, nil
}

func MakeGearboxStack(gb gearbox.Gearboxer, ns *NamedStackModel) (gbns *gears.NamedStack, sts Status) {
	//	gbns = gears.NewNamedStackModel(gb.GetGears(), types.StackId(ns.GetId()))
	gbns = gears.NewNamedStack(types.StackId(ns.GetId()))
	sts = gbns.Refresh(gb.GetGears())
	return gbns, sts
}

func (me *NamedStackModel) GetRelatedItems(ctx *Context) (list List, sts Status) {
	//for range only.Once {
	//	list = make(List, 0)
	//	for _, si := range me.ProjectStackItems {
	//		gsgs := gearspec.NewGearspec()
	//		sts = gsgs.Parse(si.GearspecId)
	//		if is.Error(sts) {
	//			break
	//		}
	//		gsm, sts := NewGearspecModelFromGearspecGearspec(ctx, gsgs)
	//		if is.Error(sts) {
	//			break
	//		}
	//		list = append(list, gsm)
	//
	//		ss := service.NewService()
	//		sts = ss.Parse(si.ServiceId)
	//		if is.Error(sts) {
	//			break
	//		}
	//		sm, sts := NewModelFromServiceServicer(ctx, ss)
	//		if is.Error(sts) {
	//			break
	//		}
	//		sm.GearspecId = gsm.GearspecId
	//		list = append(list, sm)
	//	}
	//}
	return list, sts

}
