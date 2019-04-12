package apimodels

import (
	"fmt"
	"gearbox/apimodeler"
	"gearbox/gearspec"
	"gearbox/only"
	"gearbox/status"
	"gearbox/types"
	"strings"
)

const GearspecType = "gearspec"

var NilGearspec = (*Gearspec)(nil)
var _ apimodeler.ApiItemer = NilGearspec

type GearspecMap map[types.Stackname]*Gearspec
type Gearspecs []*Gearspec

type Gearspec struct {
	GearspecId gearspec.Identifier   `json:"gearspec_id"`
	StackId    types.StackId         `json:"stack_id,omitempty"`
	Authority  types.AuthorityDomain `json:"authority,omitempty"`
	Stackname  types.Stackname       `json:"stackname,omitempty"`
	Role       types.StackRole       `json:"role,omitempty"`
	Revision   types.Revision        `json:"revision"`
}

func NewFromGearsGearspec(ctx *apimodeler.Context, gsgs *gearspec.Gearspec) (gs *Gearspec, sts status.Status) {
	return NewGearspec(gsgs), sts
}

func NewGearspec(gsgs *gearspec.Gearspec) *Gearspec {
	return &Gearspec{
		GearspecId: gsgs.GetIdentifier(),
		StackId:    gsgs.GetStackId(),
		Authority:  gsgs.Authority,
		Stackname:  gsgs.Stackname,
		Role:       gsgs.Role,
		Revision:   gsgs.Revision,
	}
}

func (me *Gearspec) GetItemLinkMap(*apimodeler.Context) (apimodeler.LinkMap, status.Status) {
	return apimodeler.LinkMap{}, nil
}

func (me *Gearspec) GetType() apimodeler.ItemType {
	return GearspecType
}

func (me *Gearspec) GetFullStackname() types.Stackname {
	return types.Stackname(me.GetId())
}

func (me *Gearspec) GetId() apimodeler.ItemId {
	return apimodeler.ItemId(me.GearspecId)
}

func (me *Gearspec) SetId(itemid apimodeler.ItemId) (sts status.Status) {
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

func (me *Gearspec) GetItem() (apimodeler.ApiItemer, status.Status) {
	return me, nil
}

func (me *Gearspec) GetRelatedItems(ctx *apimodeler.Context, item apimodeler.ApiItemer) (list apimodeler.List, sts status.Status) {
	return make(apimodeler.List, 0), sts
}
