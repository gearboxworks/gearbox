package gearbox

import (
	"gearbox/only"
	"gearbox/stat"
	"strings"
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

func (me RoleMap) GetStackRoleMap(stackName StackName) (rm RoleMap, status stat.Status) {
	for range only.Once {
		rm = make(RoleMap, 0)
		for rs, r := range me {
			stackName, status = GetFullStackName(stackName)
			if status.IsError() {
				break
			}
			if !strings.HasPrefix(string(rs), string(stackName)) {
				continue
			}
			rm[rs] = r
		}
	}
	return rm, status
}

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

func (me *StackRole) FileRoleSpec() RoleSpec {
	spec := Spec{}
	spec = *me.Spec
	if spec.Authority == DefaultAuthority {
		spec.Authority = ""
	}
	return spec.GetSpec()
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

func (me *StackRole) Parse(name RoleSpec) (status stat.Status) {
	for range only.Once {
		me.RoleSpec = name
		if me.Spec == nil {
			me.Spec = &Spec{}
		}
		status := me.Spec.Parse(string(me.RoleSpec))
		if status.IsError() {
			break
		}
		status = stat.NewOkStatus("stack role '%s' successfully parsed", name)
	}
	return status
}
