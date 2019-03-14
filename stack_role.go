package gearbox

import (
	"gearbox/only"
	"gearbox/util"
	"net/http"
)

const StackRoleHelpUrl = "https://docs.gearbox.works/roles/"
const MaxServicesPerRole = 10

//
// Examples:
//
// 		wordpress/dbserver:1
//

type RoleSpec string
type RoleMap map[RoleSpec]*StackRole

type StackRole struct {
	RoleSpec   RoleSpec `json:"-"` //  `json:"rolespec"`
	Label      string   `json:"label,omitempty"`
	ShortLabel string   `json:"short_label,omitempty"`
	Examples   []string `json:"examples,omitempty"`
	Optional   bool     `json:"optional,omitempty"`
	Minimum    int      `json:"min,omitempty"`
	Maximum    int      `json:"max,omitempty"`
	//Authority   AuthorityDomain `json:"authority,omitempty"`
	//StackName   StackName       `json:"stack,omitempty"`
	//ServiceType string          `json:"type,omitempty"`
	*Spec
}
type StackRoleArgs StackRole

func NewStackRole() *StackRole {
	return &StackRole{
		Spec: NewSpec(),
	}
}

func (me *StackRole) Fixup(rolespec RoleSpec) {
	me.Parse(rolespec)
	if me.Minimum == 0 && !me.Optional {
		me.Minimum = 1
	}
	if me.Maximum == 0 {
		me.Maximum = MaxServicesPerRole
	}
}

func (me *StackRole) GetStackName() (name StackName) {
	if me.Spec != nil {
		name = me.Spec.GetStackName()
	}
	return name
}

func (me *StackRole) NeedsParse() bool {
	return me.RoleSpec == "" || me.Spec.ServiceType == ""
}

func (me *StackRole) Parse(name RoleSpec) (status Status) {
	for range only.Once {
		me.RoleSpec = name
		if me.Spec == nil {
			me.Spec = &Spec{}
		}
		err := me.Spec.Parse(string(me.RoleSpec))
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
