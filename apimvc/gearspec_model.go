package apimvc

import (
	"fmt"
	"gearbox/apiworks"
	"gearbox/gearspec"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/only"
	"strings"
)

const GearspecModelType = "gearspecs"

var NilGearspecModel = (*GearspecModel)(nil)
var _ ItemModeler = NilGearspecModel

type GearspecModelMap map[types.Stackname]*GearspecModel
type GearspecModels []*GearspecModel

type GearspecModel struct {
	GearspecId gearspec.Identifier   `json:"-"`
	StackId    types.StackId         `json:"stack_id,omitempty"`
	Authority  types.AuthorityDomain `json:"authority,omitempty"`
	Stackname  types.Stackname       `json:"stackname,omitempty"`
	Specname   types.Specname        `json:"role,omitempty"`
	Revision   types.Revision        `json:"revision"`
	Model
}

func (me *GearspecModel) GetAttributeMap() apiworks.AttributeMap {
	panic("implement me")
}

func NewGearspecModelFromGearspecGearspec(ctx *Context, gsgs *gearspec.Gearspec) (gsm *GearspecModel, sts Status) {
	return &GearspecModel{
		GearspecId: gsgs.GetIdentifier(),
		StackId:    gsgs.GetStackId(),
		Authority:  gsgs.AuthorityDomain,
		Stackname:  gsgs.Stackname,
		Specname:   gsgs.Specname,
		Revision:   gsgs.Revision,
	}, sts
}

func NewGearspecModel() *GearspecModel {
	return &GearspecModel{}
}

func (me *GearspecModel) GetType() ItemType {
	return GearspecModelType
}

func (me *GearspecModel) GetFullStackname() types.Stackname {
	return types.Stackname(me.GetId())
}

func (me *GearspecModel) GetId() ItemId {
	return ItemId(me.GearspecId)
}

func (me *GearspecModel) SetId(itemid ItemId) (sts Status) {
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
