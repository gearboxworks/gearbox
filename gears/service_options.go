package gears

import (
	"gearbox/gearspec"
	"gearbox/service"
	"gearbox/types"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
	"github.com/gearboxworks/go-status/only"
)

type ServiceOptions []*ServiceOption

type ServiceOption struct {
	Orgname    types.Orgname        `json:"orgname,omitempty"`
	Options    service.Identifiers  `json:"options"`
	Default    service.Identifier   `json:"default"`
	Implements gearspec.Identifiers `json:"implements,omitempty"`

	Services       Services           `json:"-"`
	Gearspecs      gearspec.Gearspecs `json:"-"`
	DefaultService *Service           `json:"-"`
}

func (me ServiceOptions) FilterForNamedStack(stackid types.StackId) (sos ServiceOptions, sts Status) {
	for range only.Once {
		// The next 4 lines just validates the stack ID.
		// Should probably great a more explicit func to do that.
		gsi := gearspec.NewGearspec()
		sts = gsi.SetStackId(stackid)
		if is.Error(sts) {
			break
		}
		stackid = gsi.GetStackId()
		sos = make(ServiceOptions, 0)
		for _, so := range me {
			for _, s := range so.Services {
				if s.GetStackId() != stackid {
					continue
				}
				sos = append(sos, so)
			}
		}
	}
	return sos, sts
}

func (me *ServiceOption) Fixup() (sts status.Status) {
	for range only.Once {
		if me.Default != "" {
			me.DefaultService, sts = me.FixupService(me.Default)
			if is.Error(sts) {
				break
			}
		}
		me.Default = ""

		me.Services = make(Services, 0)
		s := NewService()
		for _, o := range me.Options {
			sts = s.Parse(o)
			if is.Error(sts) {
				break
			}
			s, sts = me.FixupService(s.ServiceId)
			if is.Error(sts) {
				break
			}
			me.Services = append(me.Services, s)
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

func (me *ServiceOption) FixupService(serviceid service.Identifier) (s *Service, sts status.Status) {
	s = NewService()
	sts = s.Parse(serviceid)
	return s, sts
}
