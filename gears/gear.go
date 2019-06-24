package gears

import (
	"gearbox/gear"
	"gearbox/gearspec"
	"gearbox/global"
	"gearbox/service"
	"gearbox/types"
	"gearbox/version"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/gearboxworks/go-status/only"
)

type GearBag map[gearspec.Identifier]interface{}

type GearMap map[service.Identifier]*Gear

type DefaultGearMap map[gearspec.Identifier]*Gear

type Gears []*Gear

func (me Gears) GetGearIds() service.Identifiers {
	gs := make(service.Identifiers, len(me))
	for i, s := range me {
		gs[i] = s.GearId
	}
	return gs
}

type Gear struct {
	GearId   service.Identifier `json:"gear_id,omitempty"`
	GearType types.GearType     `json:"gear_type"`
	Orgname  types.Orgname      `json:"org,omitempty"`
	Program  types.ProgramName  `json:"program,omitempty"`
	Version  types.Version      `json:"version,omitempty"`
}
type GearArgs Gear

func NewGear(args ...*GearArgs) *Gear {
	if len(args) == 0 {
		args = []*GearArgs{{}}
	}
	var g Gear
	g = Gear(*args[0])
	if g.GearType == ZeroString {
		g.GearType = ServiceGear
	}
	return &g
}

func (me *Gear) Clone() *Gear {
	_s := Gear{}
	_s = *me
	return &_s
}

func (me *Gear) GetStackId() (sid types.StackId) {
	panic("Not yet implemented")
}

func (me *Gear) GetIdentifier() (serviceId service.Identifier) {
	return me.GearId
}

func (me *Gear) SetIdentifier(serviceId service.Identifier) status.Status {
	me.GearId = serviceId
	return me.ApplyDefaults(me)
}

func (me *Gear) Parse(serviceId service.Identifier) (sts status.Status) {
	return me.SetIdentifier(serviceId)
}

func (me *Gear) CaptureGearId(g *gear.Gear) {
	me.GearId = service.Identifier(g.GetIdentifier())
	me.Orgname = g.OrgName
	me.Program = g.Program
	if g.Version == nil {
		me.Version = ZeroString
	} else {
		me.Version = types.Version(g.Version.String())
	}
}

func (me *Gear) MakeGear(serviceId service.Identifier) (g *gear.Gear, sts status.Status) {
	g = gear.NewGear()
	sts = g.Parse(gear.Identifier(serviceId))
	return g, sts
}

func (me *Gear) ApplyDefaults(defaults *Gear) (sts status.Status) {
	var g *gear.Gear
	for range only.Once {
		if defaults == nil {
			g = gear.NewGear()
		} else {
			g, sts = me.MakeGear(defaults.GetIdentifier())
			if status.IsError(sts) {
				break
			}
		}
		if g.OrgName == ZeroString {
			if defaults != nil && defaults.Orgname != ZeroString {
				g.OrgName = defaults.Orgname
			} else if me.Orgname != ZeroString {
				g.OrgName = me.Orgname
			} else {
				g.OrgName = global.DefaultOrgName
			}
		}
		if defaults == nil {
			break
		}
		if g.Program == ZeroString {
			g.Program = defaults.Program
		}
		if defaults.Version == ZeroString {
			break
		}
		defaultVersion := version.NewVersion()
		sts = defaultVersion.Parse(defaults.Version)
		if is.Error(sts) {
			break
		}
		serviceVersion := g.Version
		if serviceVersion.Revision == ZeroString {
			serviceVersion.Revision = defaultVersion.Revision
			if serviceVersion.Patch == ZeroString {
				serviceVersion.Patch = defaultVersion.Patch
				if serviceVersion.Minor == ZeroString {
					serviceVersion.Minor = defaultVersion.Minor
					if serviceVersion.Major == ZeroString {
						serviceVersion.Major = defaultVersion.Major
					}
				}
			}
		}
	}
	if is.Success(sts) {
		me.CaptureGearId(g)
	}
	return sts
}
