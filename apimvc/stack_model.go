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

const StackType ItemType = "stacks"

var NilStackModel = (*StackModel)(nil)
var _ ItemModeler = NilStackModel

type StackModelMap map[types.Stackname]*StackModel
type StackModels []*StackModel

type StackModel struct {
	Authority types.AuthorityDomain `json:"authority"`
	Stackname types.Stackname       `json:"stackname"`
	Members   StackMembers          `json:"members,omitempty"`
	Model
}

func (me *StackModel) GetAttributeMap() apiworks.AttributeMap {
	panic("implement me")
}

func NewStackModelFromGearsStack(ctx *Context, gns *gears.Stack) (ns *StackModel, sts Status) {
	for range only.Once {

		sms := make(StackMembers, len(gns.Gearspecs))

		for i, gs := range gns.Gearspecs {
			var sm *StackMember
			sm = NewStackMemberFromGearspec(ctx, gs)
			sms[i] = sm
		}

		ns = &StackModel{
			Authority: gns.Authority,
			Stackname: gns.Stackname,
			Members:   sms,
		}
	}
	return ns, sts
}

func NewStackModel(ns *gears.Stack) *StackModel {
	return &StackModel{
		Authority: ns.Authority,
		Stackname: ns.Stackname,
		Members:   make(StackMembers, 0),
	}
}

func (me *StackModel) GetType() ItemType {
	return StackType
}

func (me *StackModel) GetFullStackname() types.Stackname {
	return types.Stackname(me.GetId())
}

func (me *StackModel) GetId() ItemId {
	return ItemId(fmt.Sprintf("%s/%s", me.Authority, me.Stackname))
}

func (me *StackModel) SetId(itemid ItemId) (sts Status) {
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

func MakeGearboxStack(gb gearbox.Gearboxer, ns *StackModel) (gbns *gears.Stack, sts Status) {
	gbns = gears.NewStack(types.StackId(ns.GetId()))
	sts = gbns.Refresh(gb.GetGearRegistry())
	return gbns, sts
}
