package apimvc

import (
	"fmt"
	"gearbox/apiworks"
	"gearbox/gearbox"
	"gearbox/gears"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/only"
	"strings"
)

const NamedStackType ItemType = "stacks"

var NilNamedStackModel = (*NamedStackModel)(nil)
var _ ItemModeler = NilNamedStackModel

type NamedStackModelMap map[types.Stackname]*NamedStackModel
type NamedStackModels []*NamedStackModel

type NamedStackModel struct {
	Authority types.AuthorityDomain `json:"authority"`
	Stackname types.Stackname       `json:"stackname"`
	Members   StackMembers          `json:"members,omitempty"`
	Model
}

func (me *NamedStackModel) GetAttributeMap() apiworks.AttributeMap {
	panic("implement me")
}

func NewNamedStackModelFromGearsNamedStack(ctx *Context, gns *gears.NamedStack) (ns *NamedStackModel, sts Status) {
	for range only.Once {

		sms := make(StackMembers, len(gns.Gearspecs))

		for i, gs := range gns.Gearspecs {
			var sm *StackMember
			sm = NewStackMemberFromGearspec(ctx, gs)
			sms[i] = sm
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

func MakeGearboxStack(gb gearbox.Gearboxer, ns *NamedStackModel) (gbns *gears.NamedStack, sts Status) {
	gbns = gears.NewNamedStack(types.StackId(ns.GetId()))
	sts = gbns.Refresh(gb.GetGearRegistry())
	return gbns, sts
}
