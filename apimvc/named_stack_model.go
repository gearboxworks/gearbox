package apimvc

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/gearbox"
	"gearbox/gears"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"strings"
)

const NamedStackType apimodeler.ItemType = "stack"

var NilNamedStackModel = (*NamedStackModel)(nil)
var _ apimodeler.ItemModeler = NilNamedStackModel

type NamedStackModelMap map[types.Stackname]*NamedStackModel
type NamedStackModels []*NamedStackModel

type NamedStackModel struct {
	Authority types.AuthorityDomain `json:"authority"`
	Stackname types.Stackname       `json:"stackname"`
	Members   StackMembers          `json:"members,omitempty"`
}

func NewNamedStackModelFromGearsNamedStack(ctx *apimodeler.Context, gns *gears.NamedStack) (ns *NamedStackModel, sts status.Status) {
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

func (me *NamedStackModel) GetItemLinkMap(*apimodeler.Context) (apimodeler.LinkMap, status.Status) {
	return apimodeler.LinkMap{}, nil
}

func (me *NamedStackModel) GetType() apimodeler.ItemType {
	return NamedStackType
}

func (me *NamedStackModel) GetFullStackname() types.Stackname {
	return types.Stackname(me.GetId())
}

func (me *NamedStackModel) GetId() apimodeler.ItemId {
	return apimodeler.ItemId(fmt.Sprintf("%s/%s", me.Authority, me.Stackname))
}

func (me *NamedStackModel) SetStackId(itemid apimodeler.ItemId) (sts status.Status) {
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

func (me *NamedStackModel) GetItem() (apimodeler.ItemModeler, status.Status) {
	return me, nil
}

func MakeGearboxStack(gb gearbox.Gearboxer, ns *NamedStackModel) (gbns *gears.NamedStack, sts status.Status) {
	//	gbns = gears.NewNamedStackModel(gb.GetGears(), types.StackId(ns.GetId()))
	gbns = gears.NewNamedStack(types.StackId(ns.GetId()))
	sts = gbns.Refresh(gb.GetGears())
	return gbns, sts
}

func (me *NamedStackModel) GetRelatedItems(ctx *apimodeler.Context) (list apimodeler.List, sts status.Status) {
	//for range only.Once {
	//	list = make(apimodeler.List, 0)
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
