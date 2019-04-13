package apimvc

import (
	"gearbox/gear"
	"gearbox/gearspec"
	"gearbox/global"
	"gearbox/types"
)

type StackMemberMap map[gearspec.Identifier]*StackMember
type StackMembers []*StackMember
type StackMember struct {
	GearspecId       gearspec.Identifier
	OrgName          types.Orgname           `json:"orgname,omitempty"`
	DefaultServiceId gear.Identifier         `json:"default,omitempty"`
	Shareable        global.ShareableChoices `json:"shareable,omitempty"`
	ServiceIds       gear.Identifiers        `json:"services,omitempty"`
}
