package gears

import (
	"gearbox/gear"
	"gearbox/gearspec"
	"gearbox/global"
	"gearbox/only"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
	"gearbox/version"
)

type ServiceBag map[gearspec.Identifier]interface{}

type ServiceMap map[types.ServiceId]*Service

type DefaultServiceMap map[gearspec.Identifier]*Service

type Services []*Service

func (me Services) ServiceIds() types.ServiceIds {
	services := make(types.ServiceIds, len(me))
	for i, s := range me {
		services[i] = s.ServiceId
	}
	return services
}

type Service struct {
	ServiceId   types.ServiceId     `json:"service_id,omitempty"`
	Orgname     types.Orgname       `json:"org,omitempty"`
	ServiceType types.ServiceType   `json:"type,omitempty"`
	Program     types.ProgramName   `json:"program,omitempty"`
	Version     types.Version       `json:"version,omitempty"`
	GearspecId  gearspec.Identifier `json:"gearspec_id,omitempty"`
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

func (me *Service) GetIdentifier() (serviceId types.ServiceId) {
	return me.ServiceId
}

func (me *Service) SetIdentifier(serviceId types.ServiceId) status.Status {
	me.ServiceId = serviceId
	return me.ApplyDefaults(me)
}

func (me *Service) Parse(serviceId types.ServiceId) (sts status.Status) {
	return me.SetIdentifier(serviceId)
}

func (me *Service) CaptureGearId(g *gear.Gear) {
	me.ServiceId = types.ServiceId(g.GetIdentifier())
	me.Orgname = g.OrgName
	me.ServiceType = g.ServiceType
	me.Program = g.Program
	if g.Version == nil {
		me.Version = ""
	} else {
		me.Version = types.Version(g.Version.String())
	}
}

func (me *Service) GetGear(serviceId types.ServiceId) (g *gear.Gear, sts status.Status) {
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
		if g.ServiceType == "" {
			g.ServiceType = defaults.ServiceType
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
