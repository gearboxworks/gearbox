package gearbox

import (
	"gearbox/only"
	"gearbox/util"
	"net/http"
)

const StackRoleHelpUrl = "https://docs.gearbox.works/roles/"

//
// Examples:
//
// 		wordpress/dbserver:1
//

type RoleName string
type RoleMap map[RoleName]*StackRole

type StackRole struct {
	Label       string    `json:"label,omitempty"`
	Role        RoleName  `json:"stack_role"`
	ShortLabel  string    `json:"short_label,omitempty"`
	Examples    []string  `json:"examples,omitempty"`
	StackName   StackName `json:"stack,omitempty"`
	ServiceType string    `json:"service_type,omitempty"`
	Optional    bool      `json:"optional,omitempty"`
	*Spec
}
type StackRoleArgs StackRole

func NewStackRole() *StackRole {
	return &StackRole{
		Spec: NewSpec(),
	}
}

func (me *StackRole) GetStackName() (name string) {
	if me.Spec != nil {
		name = me.Spec.GetStackName()
	}
	return name
}

func (me *StackRole) NeedsParse() bool {
	return me.Role == "" || me.Spec.Role == ""
}

func (me *StackRole) Parse(name RoleName) (status Status) {
	for range only.Once {
		me.Role = name
		err := me.Spec.Parse(string(me.Role))
		if err != nil {
			status = NewStatus(&StatusArgs{
				HelpfulError: err.(util.HelpfulError),
				HttpStatus:   http.StatusBadRequest,
			})
		}
		status = NewOkStatus("stack role '%s' successfully parsed", name)
	}
	return status
}
