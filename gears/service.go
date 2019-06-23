package gears

import (
	"gearbox/gear"
	"gearbox/gearspec"
	"gearbox/global"
	"gearbox/only"
	"gearbox/service"
	"gearbox/types"
	"gearbox/version"
	"github.com/gearboxworks/go-status"
	"github.com/gearboxworks/go-status/is"
)

type ServiceBag map[gearspec.Identifier]interface{}

type ServiceMap map[service.Identifier]*Service

type DefaultServiceMap map[gearspec.Identifier]*Service

type Services []*Service

func (me Services) ServiceIds() service.Identifiers {
	services := make(service.Identifiers, len(me))
	for i, s := range me {
		services[i] = s.ServiceId
	}
	return services
}

type Service struct {
	ServiceId service.Identifier `json:"service_id,omitempty"`
	Orgname   types.Orgname      `json:"org,omitempty"`
	Program   types.ProgramName  `json:"program,omitempty"`
	Version   types.Version      `json:"version,omitempty"`
	//GearspecId  gearspec.Identifier  `json:"gearspec_id,omitempty"`
}
type ServiceArgs Service

func NewService() *Service {
	return &Service{}
}

func (me *Service) Clone() *Service {
	_s := Service{}
	_s = *me
	return &_s
}

func (me *Service) GetStackId() (sid types.StackId) {
	panic("Not yet implemented")
}

func (me *Service) GetIdentifier() (serviceId service.Identifier) {
	return me.ServiceId
}

func (me *Service) SetIdentifier(serviceId service.Identifier) status.Status {
	me.ServiceId = serviceId
	return me.ApplyDefaults(me)
}

func (me *Service) Parse(serviceId service.Identifier) (sts status.Status) {
	return me.SetIdentifier(serviceId)
}

func (me *Service) CaptureGearId(g *gear.Gear) {
	me.ServiceId = service.Identifier(g.GetIdentifier())
	me.Orgname = g.OrgName
	me.Program = g.Program
	if g.Version == nil {
		me.Version = ""
	} else {
		me.Version = types.Version(g.Version.String())
	}
}

func (me *Service) GetGear(serviceId service.Identifier) (g *gear.Gear, sts status.Status) {
	g = gear.NewGear()
	sts = g.Parse(gear.Identifier(serviceId))
	return g, sts
}

func (me *Service) ApplyDefaults(defaults *Service) (sts status.Status) {
	var g *gear.Gear
	for range only.Once {
		if defaults == nil {
			g = gear.NewGear()
		} else {
			g, sts = me.GetGear(defaults.GetIdentifier())
			if status.IsError(sts) {
				break
			}
		}
		if g.OrgName == "" {
			if defaults != nil && defaults.Orgname != "" {
				g.OrgName = defaults.Orgname
			} else if me.Orgname != "" {
				g.OrgName = me.Orgname
			} else {
				g.OrgName = global.DefaultOrgName
			}
		}
		if defaults == nil {
			break
		}
		if g.Program == "" {
			g.Program = defaults.Program
		}
		if defaults.Version == "" {
			break
		}
		defaultVersion := version.NewVersion()
		sts = defaultVersion.Parse(defaults.Version)
		if is.Error(sts) {
			break
		}
		serviceVersion := g.Version
		if serviceVersion.Revision == "" {
			serviceVersion.Revision = defaultVersion.Revision
			if serviceVersion.Patch == "" {
				serviceVersion.Patch = defaultVersion.Patch
				if serviceVersion.Minor == "" {
					serviceVersion.Minor = defaultVersion.Minor
					if serviceVersion.Major == "" {
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
