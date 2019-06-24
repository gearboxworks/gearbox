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
	Implements gearspec.Identifiers `json:"implements"`

	Gearspecs   Gearspecs `json:"-"`
	DefaultGear *Gear     `json:"-"`
}

func (me GearOptions) FilterForNamedStack(stackid types.StackId) (sos GearOptions, sts Status) {
	for range only.Once {
		if is.Error(sts) {
			break
		}
		//ns := NewNamedStack(stackid)
		//sos = make(GearOptions, 0)
		//for _, so := range me {
		//	gs := gearspec.NewGearspec()
		//	for _, g := range so.Gears {
		//		sts = gs.Parse(gsid)
		//		if is.Error(sts) {
		//			sts.Log()
		//			continue
		//		}
		//		if gs.GetStackId() != ns.GetIdentifier() {
		//			continue
		//		}
		//		sos = append(sos, so)
		//	}
		//}
	}
	return sos, sts
}

func (me *GearOption) Fixup(gr *GearRegistry) (sts status.Status) {
	for range only.Once {
		if me.Default != ZeroString {
			me.DefaultGear, sts = me.FixupGear(me.Default)
			if is.Error(sts) {
				break
			}
		}
		tgs := NewGearspec()
		gsm := gr.Gearspecs.GetMap()
		gm := make(GearMap, 0)
		me.Gearspecs = make(Gearspecs, 0)
		tg := NewGear()
		for _, gsid := range me.Implements {
			sts = tgs.Parse(gsid)
			if is.Error(sts) {
				status.Fail().SetMessage("unable to parse gearspec '%s'", gsid).Log()
				continue
			}
			gs, ok := gsm[tgs.GetIdentifier()]
			if !ok {
				status.Fail().SetMessage("gearspec '%s' not found", gsid).Log()
				continue
			}
			gs.Gears = make(Gears, 0)
			for _, gid := range me.Options {
				sts = tg.Parse(gid)
				if is.Error(sts) {
					sts.Log()
					continue
				}
				g, ok := gm[tg.GearId]
				if !ok {
					g, sts = me.FixupGear(tg.GearId)
					if is.Error(sts) {
						sts.Log()
						continue
					}
					gm[g.GearId] = g
				}
				gs.Gears = append(gs.Gears, g)
			}
			me.Gearspecs = append(me.Gearspecs, gs)
		}

		me.Implements = nil
		me.Options = nil
		me.Default = ZeroString

	}
	return sts
}

//
func (me *GearOption) FixupGear(gearid service.Identifier) (g *Gear, sts status.Status) {
	g = NewGear()
	sts = g.Parse(gearid)
	return g, sts
}
