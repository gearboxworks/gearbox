package gears

import (
	"gearbox/gearspec"
	"gearbox/global"
	"gearbox/only"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
)

const MaxServicesPerRole = 10

type StackRoleMap map[gearspec.Identifier]*StackRole
type StackRoles []*StackRole
type StackRole struct {
	GearspecId gearspec.Identifier `json:"-"`
	Name       string              `json:"program,omitempty"`
	Label      string              `json:"label,omitempty"`
	Examples   []string            `json:"examples,omitempty"`
	Optional   bool                `json:"optional,omitempty"`
	Minimum    int                 `json:"min,omitempty"`
	Maximum    int                 `json:"max,omitempty"`
	*gearspec.Gearspec
}
type StackRoleArgs StackRole

func NewStackRole() *StackRole {
	return &StackRole{
		Gearspec: gearspec.NewGearspec(),
	}
}

func (me StackRoleMap) FilterByNamedStack(stackid types.StackId) (nsrm StackRoleMap, sts status.Status) {
	for range only.Once {
		gsi := gearspec.NewGearspec()
		sts = gsi.SetId(stackid)
		if is.Error(sts) {
			break
		}
		stackid = gsi.GetStackId()
		nsrm = make(StackRoleMap, 0)
		for i, r := range me {
			if r.GetStackId() != stackid {
				continue
			}
			nsrm[i] = r
		}
	}
	return nsrm, sts
}

func (me *StackRole) GetGearspecId() gearspec.Identifier {
	spec := gearspec.Gearspec{}
	spec = *me.Gearspec
	if spec.Authority == "" {
		spec.Authority = global.DefaultAuthority
	}
	return gearspec.Identifier(spec.String())
}

func (me *StackRole) Fixup(rolespec gearspec.Identifier) (sts status.Status) {
	for range only.Once {
		sts := me.Parse(rolespec)
		if status.IsError(sts) {
			break
		}
		if me.Minimum == 0 && !me.Optional {
			me.Minimum = 1
		}
		if me.Maximum == 0 {
			me.Maximum = MaxServicesPerRole
		}
	}
	return sts
}

func (me *StackRole) GetStackname() (name types.Stackname) {
	if me.Gearspec != nil {
		name = me.Gearspec.GetStackname()
	}
	return name
}

func (me *StackRole) NeedsParse() bool {
	return me.GearspecId == "" || me.Gearspec.Role == ""
}

func (me *StackRole) Parse(name gearspec.Identifier) (sts status.Status) {
	for range only.Once {
		me.GearspecId = name
		if me.Gearspec == nil {
			me.Gearspec = &gearspec.Gearspec{}
		}
		sts := me.Gearspec.Parse(gearspec.Identifier(me.GearspecId))
		if status.IsError(sts) {
			break
		}
		sts = status.Success("stack role '%s' successfully parsed", name)
	}
	return sts
}
