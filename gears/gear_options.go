package gears

import (
	"gearbox/gearspec"
	"gearbox/service"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/gearboxworks/go-status/only"
)

type GearOptions []*GearOption

type GearOption struct {
	Orgname    types.Orgname        `json:"orgname,omitempty"`
	Options    service.Identifiers  `json:"options"`
	Default    service.Identifier   `json:"default"`
	Implements gearspec.Identifiers `json:"implements,omitempty"`

	Gears       Gears              `json:"-"`
	Gearspecs   gearspec.Gearspecs `json:"-"`
	DefaultGear *Gear              `json:"-"`
}

func (me GearOptions) FilterForNamedStack(stackid types.StackId) (sos GearOptions, sts Status) {
	for range only.Once {
		// The next 4 lines just validates the stack ID.
		// Should probably great a more explicit func to do that.
		gsi := gearspec.NewGearspec()
		sts = gsi.SetStackId(stackid)
		if is.Error(sts) {
			break
		}
		stackid = gsi.GetStackId()
		sos = make(GearOptions, 0)
		for _, so := range me {
			for _, s := range so.Gears {
				if s.GetStackId() != stackid {
					continue
				}
				sos = append(sos, so)
			}
		}
	}
	return sos, sts
}

func (me *GearOption) Fixup() (sts status.Status) {
	for range only.Once {
		if me.Default != ZeroString {
			me.DefaultGear, sts = me.FixupGear(me.Default)
			if is.Error(sts) {
				break
			}
		}
		me.Default = ZeroString

		me.Gears = make(Gears, 0)
		s := NewGear()
		for _, o := range me.Options {
			sts = s.Parse(o)
			if is.Error(sts) {
				break
			}
			s, sts = me.FixupGear(s.GearId)
			if is.Error(sts) {
				break
			}
			me.Gears = append(me.Gears, s)
		}
		me.Options = nil

		me.Gearspecs = make(gearspec.Gearspecs, 0)
		for _, gsid := range me.Implements {
			gs := gearspec.NewGearspec()
			sts = gs.Parse(gsid)
			if is.Error(sts) {
				break
			}
			me.Gearspecs = append(me.Gearspecs, gs)
		}
		me.Implements = nil
	}
	return sts
}

//
func (me *GearOption) FixupGear(gearid service.Identifier) (g *Gear, sts status.Status) {
	g = NewGear()
	sts = g.Parse(gearid)
	return g, sts
}
