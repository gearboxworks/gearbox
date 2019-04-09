package apimodels

import (
	"gearbox/api/global"
	"gearbox/gear"
	"gearbox/gearspec"
	"gearbox/types"
)

type StackMemberMap map[gearspec.Identifier]*StackMember
type StackMembers []*StackMember
type StackMember struct {
	GearspecId       gearspec.Identifier
	OrgName          types.OrgName           `json:"orgname,omitempty"`
	DefaultServiceId gear.Identifier         `json:"default,omitempty"`
	Shareable        global.ShareableChoices `json:"shareable,omitempty"`
	ServiceIds       gear.Identifiers        `json:"services,omitempty"`
}
