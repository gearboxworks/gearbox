package apimvc

import (
	"gearbox/apimodeler"
	"gearbox/gears"
	"gearbox/gearspec"
	"gearbox/global"
	"gearbox/service"
	"gearbox/types"
)

func NewStackMemberFromGearsStackRoleServices(ctx *apimodeler.Context, grss *gears.RoleServices, gsr *gears.StackRole) (rss *StackMember) {
	var dsi service.Identifier
	if grss.DefaultService != nil {
		dsi = grss.DefaultService.GetIdentifier()
	}
	return &StackMember{
		GearspecId:       gsr.GearspecId,
		Role:             gsr.GetRole(),
		Revision:         gsr.GetRevision(),
		DefaultServiceId: dsi,
		Shareable:        grss.Shareable,
		ServiceIds:       grss.Services.ServiceIds(),
	}
}

type StackMembers []*StackMember
type StackMember struct {
	GearspecId       gearspec.Identifier     `json:"gearspec_id"`
	Role             types.StackRole         `json:"role,omitempty"`
	Revision         types.Revision          `json:"revision,omitempty"`
	DefaultServiceId service.Identifier      `json:"default_service,omitempty"`
	Shareable        global.ShareableChoices `json:"shareable,omitempty"`
	ServiceIds       service.Identifiers     `json:"services,omitempty"`
}
