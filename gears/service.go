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

type ServiceMap map[gear.Identifier]*Service

type DefaultServiceMap map[gearspec.Identifier]*Service

type Services []*Service

func (me Services) ServiceIds() gear.Identifiers {
	services := make(gear.Identifiers, len(me))
	for i, s := range me {
		services[i] = s.ServiceId
	}
	return services
}

type Service struct {
	ServiceId gear.Identifier   `json:"service_id,omitempty"`
	OrgName   types.OrgName     `json:"org,omitempty"`
	Type      types.ServiceType `json:"type,omitempty"`
	Program   types.ProgramName `json:"program,omitempty"`
	Version   types.Version     `json:"version,omitempty"`
}
type ServiceArgs Service

func NewService(serviceId gear.Identifier, args ...*ServiceArgs) *Service {
	var _args *ServiceArgs
	if len(args) == 0 {
		_args = &ServiceArgs{}
	} else {
		_args = args[0]
	}
	_args.ServiceId = serviceId
	s := Service{}
	s = Service(*_args)
	return &s
}

func (me *Service) SetIdentifier(serviceId gear.Identifier) status.Status {
	me.ServiceId = serviceId
	return me.ApplyDefaults(me)
}

func (me *Service) Parse() (gid *gear.Gear, sts status.Status) {
	gid = gear.NewGear()
	for range only.Once {
		sts = gid.Parse(gear.Identifier(me.ServiceId))
		if status.IsError(sts) {
			break
		}
	}
	if is.Success(sts) {
		me.CaptureGearId(gid)
	}
	return gid, sts
}

func (me *Service) ApplyDefaults(defaults *Service) (sts status.Status) {
	var gid *gear.Gear
	for range only.Once {
		gid, sts = me.Parse()
		if status.IsError(sts) {
			break
		}
		if gid.OrgName == "" {
			if defaults != nil && defaults.OrgName != "" {
				gid.OrgName = defaults.OrgName
			} else if me.OrgName != "" {
				gid.OrgName = me.OrgName
			} else {
				gid.OrgName = global.DefaultOrgName
			}
		}
		if defaults == nil {
			break
		}
		if gid.Program == "" {
			gid.Program = defaults.Program
		}
		if gid.Type == "" {
			gid.Type = defaults.Type
		}
		if defaults.Version == "" {
			break
		}
		defaultVersion := version.NewVersion()
		sts = defaultVersion.Parse(defaults.Version)
		if is.Error(sts) {
			break
		}
		serviceVersion := gid.Version
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
		me.CaptureGearId(gid)
	}
	return sts
}

func (me *Service) CaptureGearId(gid *gear.Gear) {
	me.ServiceId = gear.Identifier(gid.GetIdentifier())
	me.OrgName = gid.OrgName
	me.Type = gid.Type
	me.Program = gid.Program
	me.Version = types.Version(gid.Version.String())
}
