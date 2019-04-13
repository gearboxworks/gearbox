package apimvc

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/gearspec"
	"gearbox/only"
	"gearbox/status"
	"gearbox/types"
	"strings"
)

const GearspecModelType = "gearspec"

var NilGearspecModel = (*GearspecModel)(nil)
var _ apimodeler.ItemModeler = NilGearspecModel

type GearspecModelMap map[types.Stackname]*GearspecModel
type GearspecModels []*GearspecModel

type GearspecModel struct {
	GearspecId gearspec.Identifier   `json:"gearspec_id"`
	StackId    types.StackId         `json:"stack_id,omitempty"`
	Authority  types.AuthorityDomain `json:"authority,omitempty"`
	Stackname  types.Stackname       `json:"stackname,omitempty"`
	Role       types.StackRole       `json:"role,omitempty"`
	Revision   types.Revision        `json:"revision"`
}

func NewModelFromGearspecGearspec(ctx *apimodeler.Context, gsgs *gearspec.Gearspec) (gs *GearspecModel, sts status.Status) {
	return &GearspecModel{
		GearspecId: gsgs.GetIdentifier(),
		StackId:    gsgs.GetStackId(),
		Authority:  gsgs.Authority,
		Stackname:  gsgs.Stackname,
		Role:       gsgs.Role,
		Revision:   gsgs.Revision,
	}, sts
}

func NewGearspec() *GearspecModel {
	return &GearspecModel{}
}

func (me *GearspecModel) GetItemLinkMap(*apimodeler.Context) (apimodeler.LinkMap, status.Status) {
	return apimodeler.LinkMap{}, nil
}

func (me *GearspecModel) GetType() apimodeler.ItemType {
	return GearspecModelType
}

func (me *GearspecModel) GetFullStackname() types.Stackname {
	return types.Stackname(me.GetId())
}

func (me *GearspecModel) GetId() apimodeler.ItemId {
	return apimodeler.ItemId(me.GearspecId)
}

func (me *GearspecModel) SetId(itemid apimodeler.ItemId) (sts status.Status) {
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

func (me *GearspecModel) GetItem() (apimodeler.ItemModeler, status.Status) {
	return me, nil
}

func (me *GearspecModel) GetRelatedItems(ctx *apimodeler.Context) (list apimodeler.List, sts status.Status) {
	return make(apimodeler.List, 0), sts
}
