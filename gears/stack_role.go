package gears

import (
	"fmt"
	"gearbox/gearspec"
	"gearbox/global"
	"gearbox/only"
	"gearbox/service"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
)

const MaxServicesPerRole = 10

type StackRoleMap map[gearspec.Identifier]*StackRole
type StackRoles []*StackRole
type StackRole struct {
	GearspecId       gearspec.Identifier    `json:"gearspec_id"`
	AuthorityDomain  types.AuthorityDomain  `json:"authority"`
	Stackname        types.Stackname        `json:"stackname"`
	Role             types.StackRole        `json:"role"`
	Revision         types.Revision         `json:"revision,omitempty"`
	Label            string                 `json:"label,omitempty"`
	Described        string                 `json:"described,omitempty"`
	Examples         []string               `json:"examples,omitempty"`
	Optional         bool                   `json:"optional,omitempty"`
	Shareable        global.ShareableChoice `json:"shareable"`
	DefaultServiceId service.Identifier     `json:"default,omitempty"`
	Minimum          int                    `json:"min,omitempty"`
	Maximum          int                    `json:"max,omitempty"`
	stackid          types.StackId
}
type StackRoleArgs StackRole

func NewStackRole() *StackRole {
	return &StackRole{}
}

func (me StackRoles) FilterByNamedStack(stackid types.StackId) (nsrs StackRoles, sts status.Status) {
	for range only.Once {
		gs := gearspec.NewGearspec()
		sts = gs.SetStackId(stackid)
		if is.Error(sts) {
			break
		}
		stackid = gs.GetStackId()
		nsrs = make(StackRoles, 0)
		for i, r := range me {
			if r.GetStackId() != stackid {
				continue
			}
			nsrs[i] = r
		}
	}
	return nsrs, sts
}

func (me *StackRole) GetDefaultService() (ds *Service) {
	for range only.Once {
		if me.DefaultServiceId == "" {
			break
		}
		ds = NewService()
		sts := ds.Parse(me.DefaultServiceId)
		if is.Error(sts) {
			sts.Log()
		}
	}
	return ds
}

func (me *StackRole) GetStackId() types.StackId {
	for range only.Once {
		if me.stackid != "" {
			break
		}
		if me.AuthorityDomain == "" {
			me.AuthorityDomain = global.DefaultAuthorityDomain
		}
		me.stackid = types.StackId(fmt.Sprintf("%s/%s", me.AuthorityDomain, me.Stackname))
	}
	return me.stackid
}

func (me *StackRole) GetGearspecId() gearspec.Identifier {
	for range only.Once {
		if me.GearspecId != "" {
			break
		}
		gs := gearspec.NewGearspec()
		gs.AuthorityDomain = me.AuthorityDomain
		gs.Stackname = me.Stackname
		gs.Role = me.Role
		gs.Revision = me.Revision
		me.GearspecId = gs.GetIdentifier()
	}
	return me.GearspecId
}

func (me *StackRole) String() string {
	return fmt.Sprintf("%#v", me)
}

func (me *StackRole) Parse(gsi gearspec.Identifier) (sts status.Status) {
	for range only.Once {
		gs := gearspec.NewGearspec()
		sts = gs.Parse(gsi)
		if is.Error(sts) {
			break
		}
		me.AuthorityDomain = gs.AuthorityDomain
		me.Stackname = gs.Stackname
		me.Role = gs.Role
		me.Revision = gs.Revision
		me.stackid = gs.GetStackId()
		me.GearspecId = gs.Identifier
		sts = status.Success("stack role '%s' successfully parsed", gsi)
	}
	return sts
}

func (me *StackRole) Fixup() (sts status.Status) {
	for range only.Once {
		save := *me
		if me.GearspecId != "" {
			sts = me.Parse(me.GearspecId)
			if status.IsError(sts) {
				break
			}
		}
		if me.Stackname == "" {
			sts = status.Fail().SetMessage("stackname cannot be null for stack role: %s", me)
			break
		}
		if me.Role == "" {
			sts = status.Fail().SetMessage("role cannot be null for stack role: %s", me)
			break
		}
		if save.Stackname != "" {
			me.Stackname = save.Stackname
		}
		if save.Role != "" {
			me.Role = save.Role
		}
		if me.stackid == "" {
			me.stackid = me.GetStackId()
		}
		if me.GearspecId == "" {
			me.GearspecId = me.GetGearspecId()
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
