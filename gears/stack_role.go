package gears

import (
	"gearbox/gearspecid"
	"gearbox/global"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
)

const MaxServicesPerRole = 10

type StackRoleMap map[gsid.Identifier]*StackRole
type StackRoles []*StackRole
type StackRole struct {
	GearspecId gsid.Identifier `json:"-"`
	Name       string          `json:"program,omitempty"`
	Label      string          `json:"label,omitempty"`
	Examples   []string        `json:"examples,omitempty"`
	Optional   bool            `json:"optional,omitempty"`
	Minimum    int             `json:"min,omitempty"`
	Maximum    int             `json:"max,omitempty"`
	*gsid.Id
}
type StackRoleArgs StackRole

func NewStackRole() *StackRole {
	return &StackRole{
		Id: gsid.NewGearspecId(),
	}
}

func (me StackRoleMap) FilterByNamedStack(stackid types.StackId) (nsrm StackRoleMap, sts status.Status) {
	for range only.Once {
		gsi := gsid.NewGearspecId()
		sts = gsi.SetStackId(stackid)
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

func (me *StackRole) GetGearspecId() gsid.Identifier {
	spec := gsid.Id{}
	spec = *me.Id
	if spec.Authority == "" {
		spec.Authority = global.DefaultAuthority
	}
	return gsid.Identifier(spec.String())
}

func (me *StackRole) Fixup(rolespec gsid.Identifier) (sts status.Status) {
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
	if me.Id != nil {
		name = me.Id.GetStackname()
	}
	return name
}

func (me *StackRole) NeedsParse() bool {
	return me.GearspecId == "" || me.Id.Role == ""
}

func (me *StackRole) Parse(name gsid.Identifier) (sts status.Status) {
	for range only.Once {
		me.GearspecId = name
		if me.Id == nil {
			me.Id = &gsid.Id{}
		}
		sts := me.Id.Parse(gsid.Identifier(me.GearspecId))
		if status.IsError(sts) {
			break
		}
		sts = status.Success("stack role '%s' successfully parsed", name)
	}
	return sts
}
