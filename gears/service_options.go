package gears

import (
	"gearbox/gearspecid"
	"gearbox/only"
	"gearbox/service"
	"gearbox/status"
	"gearbox/status/is"
	"gearbox/types"
)

// Example from gears.json/services
var _ = `{
  "gearbox.works/lamp/apache": {
	 "org": "gearboxworks",
	 "default": "apache",
	 "options": [
		"apache:2.4"
	 ]
  },
  "gearbox.works/lamp/mysql": {
	 "org": "gearboxworks",
	 "default": "mysql",
	 "options": [
		"mysql:5.5",
		"mysql:5.6",
		"mysql:5.7",
		"mysql:8.0"
	 ]
  },
  "gearbox.works/lamp/php": {
	 "org": "gearboxworks",
	 "default": "php",
	 "options": [
		"php:5.2",
		"php:5.6",
		"php:7.0",
		"php:7.1",
		"php:7.2"
	 ]
  }
}`

type ServiceOptionsMap map[gsid.Identifier]*ServiceOptions

type ServiceOptions struct {
	NamedStackId   types.StackId       `json:"-"`
	OrgName        types.OrgName       `json:"org,omitempty"`
	Default        service.Identifier  `json:"default"`
	Shareable      ShareableChoices    `json:"shareable"`
	ServiceIds     service.Identifiers `json:"options,omitempty"`
	DefaultService *Service            `json:"-"`
	Services       Services            `json:"-"`
}

func NewServiceOptions(nsid types.StackId) *ServiceOptions {
	return &ServiceOptions{
		NamedStackId: nsid,
	}
}

func (me ServiceOptionsMap) FilterForNamedStack(stackid types.StackId) (nsrm ServiceOptionsMap, sts status.Status) {
	for range only.Once {
		gsi := gsid.NewGearspecId()
		sts = gsi.Parse(gsid.Identifier(stackid))
		if is.Error(sts) {
			break
		}
		stackid = types.StackId(gsi.String())
		nsrm = make(ServiceOptionsMap, 0)
		for i, so := range me {
			if so.NamedStackId != stackid {
				continue
			}
			nsrm[i] = so
		}
	}
	return nsrm, sts
}

func (me *ServiceOptions) Fixup(nsid types.StackId) (sts status.Status) {
	for range only.Once {
		gsi := gsid.NewGearspecId()
		sts = gsi.Parse(gsid.Identifier(me.NamedStackId))
		if is.Error(sts) {
			break
		}
		me.NamedStackId = gsi.GetStackId()
		if me.Default != "" {
			me.DefaultService, sts = me.FixupService(me.Default)
			if is.Error(sts) {
				break
			}
			if me.DefaultService.OrgName == "" {
				me.DefaultService.OrgName = me.OrgName
			}
		}
		me.Default = ""
		me.Services = make(Services, len(me.ServiceIds))
		for i, sid := range me.ServiceIds {
			me.Services[i], sts = me.FixupService(sid)
			if is.Error(sts) {
				break
			}
		}
		me.ServiceIds = nil
	}
	return sts
}

func (me *ServiceOptions) FixupService(serviceId service.Identifier) (s *Service, sts status.Status) {
	for range only.Once {
		s = NewService(serviceId)
		if me.DefaultService != nil {
			sts = s.ApplyDefaults(me.DefaultService)
			break
		}
		me.DefaultService = NewService(serviceId)
		_, sts = me.DefaultService.Parse()
		if is.Error(sts) {
			break
		}
		sts = s.ApplyDefaults(me.DefaultService)
	}
	return s, sts
}
